package api

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/pires/go-proxyproto"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

type ApiServer struct {
	ctx          context.Context
	db           *sql.DB
	grpcListener net.Listener
	grpcServer   *grpc.Server
	log          *zap.Logger
	registrant   *registrant.Registrant
	service      message_api.ReplicationApiServer
	Wg           sync.WaitGroup
}

func NewAPIServer(
	ctx context.Context,
	writerDB *sql.DB,
	log *zap.Logger,
	port int,
	registrant *registrant.Registrant,
	enableReflection bool,
) (*ApiServer, error) {
	grpcListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		return nil, err
	}
	s := &ApiServer{
		ctx: ctx,
		db:  writerDB,
		grpcListener: &proxyproto.Listener{
			Listener:          grpcListener,
			ReadHeaderTimeout: 10 * time.Second,
		},
		log:        log.Named("api"),
		registrant: registrant,
		Wg:         sync.WaitGroup{},
	}

	// TODO: Add interceptors

	options := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 5 * time.Minute,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			MinTime:             15 * time.Second,
		}),
		// grpc.MaxRecvMsgSize(s.Config.Options.MaxMsgSize),
	}

	s.grpcServer = grpc.NewServer(options...)

	if enableReflection {
		// Register reflection service on gRPC server.
		reflection.Register(s.grpcServer)
		s.log.Info("enabling gRPC Server Reflection")
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s.grpcServer, healthcheck)

	replicationService, err := NewReplicationApiService(ctx, log, registrant, writerDB)
	if err != nil {
		return nil, err
	}
	s.service = replicationService
	message_api.RegisterReplicationApiServer(s.grpcServer, s.service)

	tracing.GoPanicWrap(s.ctx, &s.Wg, "grpc", func(ctx context.Context) {
		s.log.Info("serving grpc", zap.String("address", s.grpcListener.Addr().String()))
		if err = s.grpcServer.Serve(s.grpcListener); err != nil &&
			!isErrUseOfClosedConnection(err) {
			s.log.Error("serving grpc", zap.Error(err))
		}
		log.Info("grpc thread has exited")
	})

	return s, nil
}

func (s *ApiServer) Addr() net.Addr {
	return s.grpcListener.Addr()
}

func (s *ApiServer) DialGRPCTest(ctx context.Context) (*grpc.ClientConn, error) {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", s.grpcListener.Addr().String())
	println(dialAddr)
	return grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
}

func (s *ApiServer) DialGRPC(dialAddr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
}

func (s *ApiServer) gracefulShutdown(timeout time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	// Attempt to use GracefulStop up until the timeout
	go func() {
		defer cancel()
		s.grpcServer.GracefulStop()
	}()
	go func() {
		defer cancel()
		<-time.NewTimer(timeout).C
		s.log.Info("Graceful shutdown timed out. Stopping...")
		s.grpcServer.Stop()
	}()
	<-ctx.Done()
}

func (s *ApiServer) Close() {
	s.log.Info("closing")
	if s.grpcServer != nil {
		s.gracefulShutdown(1 * time.Second)
	}
	if s.grpcListener != nil {
		if err := s.grpcListener.Close(); err != nil && !isErrUseOfClosedConnection(err) {
			s.log.Error("closing grpc listener", zap.Error(err))
		}
		s.grpcListener = nil
	}
	s.Wg.Wait()
	s.log.Info("closed")
}

func isErrUseOfClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

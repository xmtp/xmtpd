package api

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

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
	log          *zap.Logger
	registrant   *registrant.Registrant
	service      message_api.ReplicationApiServer
	wg           sync.WaitGroup
}

func NewAPIServer(
	ctx context.Context,
	writerDB *sql.DB,
	log *zap.Logger,
	registrant *registrant.Registrant,
	port int,
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
		wg:         sync.WaitGroup{},
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
	grpcServer := grpc.NewServer(options...)

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	replicationService, err := NewReplicationApiService(ctx, log, registrant, writerDB)
	if err != nil {
		return nil, err
	}
	s.service = replicationService

	tracing.GoPanicWrap(s.ctx, &s.wg, "grpc", func(ctx context.Context) {
		s.log.Info("serving grpc", zap.String("address", s.grpcListener.Addr().String()))
		err := grpcServer.Serve(s.grpcListener)
		if err != nil && !isErrUseOfClosedConnection(err) {
			s.log.Error("serving grpc", zap.Error(err))
		}
	})

	return s, nil
}

func (s *ApiServer) Addr() net.Addr {
	return s.grpcListener.Addr()
}

func (s *ApiServer) Close() {
	s.log.Info("closing")

	if s.grpcListener != nil {
		err := s.grpcListener.Close()
		if err != nil {
			s.log.Error("closing grpc listener", zap.Error(err))
		}
		s.grpcListener = nil
	}

	s.wg.Wait()
	s.log.Info("closed")
}

func isErrUseOfClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

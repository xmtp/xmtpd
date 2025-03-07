package api

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"

	"google.golang.org/grpc/reflection"

	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pires/go-proxyproto"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

var (
	prometheusOnce sync.Once
)

type RegistrationFunc func(server *grpc.Server) error

type ApiServer struct {
	ctx          context.Context
	grpcListener net.Listener
	httpListener net.Listener
	grpcServer   *grpc.Server
	log          *zap.Logger
	wg           sync.WaitGroup
}

func NewAPIServer(
	ctx context.Context,
	log *zap.Logger,
	listenAddress string,
	httpListenAddress string,
	enableReflection bool,
	registrationFunc RegistrationFunc,
	httpRegistrationFunc HttpRegistrationFunc,
	jwtVerifier authn.JWTVerifier,
) (*ApiServer, error) {
	grpcListener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	httpListener, err := net.Listen("tcp", httpListenAddress)
	if err != nil {
		return nil, err
	}

	s := &ApiServer{
		ctx: ctx,
		grpcListener: &proxyproto.Listener{
			Listener:          grpcListener,
			ReadHeaderTimeout: 10 * time.Second,
		},
		httpListener: httpListener,
		log:          log.Named("api"),
		wg:           sync.WaitGroup{},
	}
	s.log.Info("Creating API server")

	prometheusOnce.Do(func() {
		prometheus.EnableHandlingTimeHistogram()
	})

	loggingInterceptor, err := server.NewLoggingInterceptor(log)
	if err != nil {
		return nil, err
	}

	unary := []grpc.UnaryServerInterceptor{
		prometheus.UnaryServerInterceptor,
	}
	stream := []grpc.StreamServerInterceptor{
		prometheus.StreamServerInterceptor,
	}

	if jwtVerifier != nil {
		interceptor := server.NewAuthInterceptor(jwtVerifier, log)
		unary = append(unary, interceptor.Unary())
		stream = append(stream, interceptor.Stream())
	}

	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unary...),
		grpc.ChainStreamInterceptor(stream...),
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 5 * time.Minute,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			MinTime:             15 * time.Second,
		}),
		grpc.ChainUnaryInterceptor(loggingInterceptor.Unary()),
		grpc.ChainStreamInterceptor(loggingInterceptor.Stream()),

		// grpc.MaxRecvMsgSize(s.Config.Options.MaxMsgSize),
	}

	s.grpcServer = grpc.NewServer(options...)
	if err := registrationFunc(s.grpcServer); err != nil {
		return nil, err
	}

	if enableReflection {
		// Register reflection service on gRPC server.
		reflection.Register(s.grpcServer)
		s.log.Info("enabling gRPC Server Reflection")
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s.grpcServer, healthcheck)

	tracing.GoPanicWrap(s.ctx, &s.wg, "grpc", func(ctx context.Context) {
		s.log.Info("serving grpc", zap.String("address", s.grpcListener.Addr().String()))
		if err = s.grpcServer.Serve(s.grpcListener); err != nil &&
			!isErrUseOfClosedConnection(err) {
			s.log.Error("serving grpc", zap.Error(err))
		}
	})

	if err := s.startHTTP(ctx, log, httpRegistrationFunc); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ApiServer) DialGRPC(ctx context.Context) (*grpc.ClientConn, error) {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", s.grpcListener.Addr().String())
	return grpc.NewClient(dialAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (s *ApiServer) Addr() net.Addr {
	return s.grpcListener.Addr()
}

func (s *ApiServer) HttpAddr() net.Addr {
	return s.httpListener.Addr()
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
		s.log.Debug("Graceful shutdown timed out. Stopping...")
		s.grpcServer.Stop()
	}()

	<-ctx.Done()
}

func (s *ApiServer) Close(timeout time.Duration) {
	s.log.Debug("closing")
	if s.grpcServer != nil {
		if timeout != 0 {
			s.gracefulShutdown(timeout)
		} else {
			s.grpcServer.Stop()
		}
	}
	if s.grpcListener != nil {
		if err := s.grpcListener.Close(); err != nil && !isErrUseOfClosedConnection(err) {
			s.log.Error("Error while closing grpc listener", zap.Error(err))
		}
		s.grpcListener = nil
	}

	if s.httpListener != nil {
		err := s.httpListener.Close()
		if err != nil {
			s.log.Error("Error while closing http listener", zap.Error(err))
		}
		s.httpListener = nil
	}

	s.wg.Wait()
	s.log.Debug("closed")
}

func isErrUseOfClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

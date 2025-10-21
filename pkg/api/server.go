// Package api implements the API server.
package api

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/pires/go-proxyproto"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type RegistrationFunc func(server *grpc.Server) error

type APIServerConfig struct {
	Ctx                context.Context
	Logger             *zap.Logger
	GRPCListener       net.Listener
	EnableReflection   bool
	RegistrationFunc   RegistrationFunc
	PromRegistry       *prometheus.Registry
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

type APIServerOption func(*APIServerConfig)

func WithContext(ctx context.Context) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Ctx = ctx }
}

func WithLogger(logger *zap.Logger) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Logger = logger }
}

func WithGRPCListener(listener net.Listener) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.GRPCListener = listener }
}

func WithReflection(enabled bool) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.EnableReflection = enabled }
}

func WithRegistrationFunc(fn RegistrationFunc) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.RegistrationFunc = fn }
}

func WithPrometheusRegistry(reg *prometheus.Registry) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.PromRegistry = reg }
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.UnaryInterceptors = append(cfg.UnaryInterceptors, interceptors...) }
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.StreamInterceptors = append(cfg.StreamInterceptors, interceptors...) }
}

type APIServer struct {
	ctx          context.Context
	grpcListener net.Listener
	grpcServer   *grpc.Server
	logger       *zap.Logger
	wg           sync.WaitGroup
}

func NewAPIServer(opts ...APIServerOption) (*APIServer, error) {
	cfg := &APIServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if cfg.GRPCListener == nil {
		return nil, fmt.Errorf("GRPCListener is required")
	}

	if cfg.RegistrationFunc == nil {
		return nil, fmt.Errorf("grpc registration function is required")
	}

	s := &APIServer{
		ctx: cfg.Ctx,
		grpcListener: &proxyproto.Listener{
			Listener:          cfg.GRPCListener,
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger: cfg.Logger.Named(utils.APILoggerName),
	}

	s.logger.Info("creating api server")

	loggingInterceptor, err := server.NewLoggingInterceptor(s.logger)
	if err != nil {
		return nil, err
	}

	openConnectionsInterceptor, err := server.NewOpenConnectionsInterceptor()
	if err != nil {
		return nil, err
	}

	unary := []grpc.UnaryServerInterceptor{
		openConnectionsInterceptor.Unary(),
		loggingInterceptor.Unary(),
	}
	stream := []grpc.StreamServerInterceptor{
		openConnectionsInterceptor.Stream(),
		loggingInterceptor.Stream(),
	}

	// Add any additional interceptors from config
	unary = append(unary, cfg.UnaryInterceptors...)
	stream = append(stream, cfg.StreamInterceptors...)

	if cfg.PromRegistry != nil {
		srvMetrics := grpcprom.NewServerMetrics(
			grpcprom.WithServerHandlingTimeHistogram(
				grpcprom.WithHistogramBuckets(
					[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
				),
			),
		)
		cfg.PromRegistry.MustRegister(srvMetrics)
		unary = append([]grpc.UnaryServerInterceptor{srvMetrics.UnaryServerInterceptor()}, unary...)
		stream = append(
			[]grpc.StreamServerInterceptor{srvMetrics.StreamServerInterceptor()},
			stream...)
	}

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(unary...),
		grpc.ChainStreamInterceptor(stream...),
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{Time: 5 * time.Minute}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			MinTime:             15 * time.Second,
		}),
	)

	if err := cfg.RegistrationFunc(s.grpcServer); err != nil {
		return nil, err
	}

	if cfg.EnableReflection {
		reflection.Register(s.grpcServer)
		s.logger.Info("enabling gRPC Server Reflection")
	}

	healthgrpc.RegisterHealthServer(s.grpcServer, health.NewServer())

	tracing.GoPanicWrap(s.ctx, &s.wg, "grpc", func(ctx context.Context) {
		s.logger.Info("serving grpc", utils.AddressField(s.grpcListener.Addr().String()))
		if err := s.grpcServer.Serve(s.grpcListener); err != nil &&
			!isErrUseOfClosedConnection(err) {
			s.logger.Error("serving grpc", zap.Error(err))
		}
	})

	return s, nil
}

func (s *APIServer) DialGRPC(ctx context.Context) (*grpc.ClientConn, error) {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", s.grpcListener.Addr().String())
	return grpc.NewClient(dialAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (s *APIServer) Addr() net.Addr {
	return s.grpcListener.Addr()
}

func (s *APIServer) gracefulShutdown(timeout time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	// Attempt to use GracefulStop up until the timeout
	go func() {
		defer cancel()
		s.grpcServer.GracefulStop()
	}()
	go func() {
		defer cancel()
		<-time.NewTimer(timeout).C
		s.logger.Debug("graceful shutdown timed out, stopping")
		s.grpcServer.Stop()
	}()

	<-ctx.Done()
}

func (s *APIServer) Close(timeout time.Duration) {
	s.logger.Debug("closing")
	if s.grpcServer != nil {
		if timeout != 0 {
			s.gracefulShutdown(timeout)
		} else {
			s.grpcServer.Stop()
		}
	}
	if s.grpcListener != nil {
		_ = s.grpcListener.Close()
	}

	s.wg.Wait()
	s.logger.Debug("closed")
}

func isErrUseOfClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

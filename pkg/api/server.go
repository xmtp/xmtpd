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
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type RegistrationFunc func(server *grpc.Server) error

type ApiServerConfig struct {
	Ctx                  context.Context
	Log                  *zap.Logger
	ListenAddress        string
	HTTPListenAddress    string
	EnableReflection     bool
	RegistrationFunc     RegistrationFunc
	HTTPRegistrationFunc HttpRegistrationFunc
	JWTVerifier          authn.JWTVerifier
	PromRegistry         *prometheus.Registry
}

type ApiServerOption func(*ApiServerConfig)

func WithContext(ctx context.Context) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.Ctx = ctx }
}

func WithLogger(log *zap.Logger) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.Log = log }
}

func WithListenAddress(addr string) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.ListenAddress = addr }
}

func WithHTTPListenAddress(addr string) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.HTTPListenAddress = addr }
}

func WithReflection(enabled bool) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.EnableReflection = enabled }
}

func WithRegistrationFunc(fn RegistrationFunc) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.RegistrationFunc = fn }
}

func WithHTTPRegistrationFunc(fn HttpRegistrationFunc) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.HTTPRegistrationFunc = fn }
}

func WithJWTVerifier(verifier authn.JWTVerifier) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.JWTVerifier = verifier }
}

func WithPrometheusRegistry(reg *prometheus.Registry) ApiServerOption {
	return func(cfg *ApiServerConfig) { cfg.PromRegistry = reg }
}

type ApiServer struct {
	ctx          context.Context
	grpcListener net.Listener
	httpListener net.Listener
	grpcServer   *grpc.Server
	log          *zap.Logger
	wg           sync.WaitGroup
}

func NewAPIServer(opts ...ApiServerOption) (*ApiServer, error) {
	cfg := &ApiServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	if cfg.Log == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if cfg.PromRegistry == nil {
		return nil, fmt.Errorf("prometheus registry is required")
	}

	if cfg.ListenAddress == "" || cfg.HTTPListenAddress == "" {
		return nil, fmt.Errorf("both listenAddress and httpListenAddress are required")
	}
	if cfg.RegistrationFunc == nil {
		return nil, fmt.Errorf("grpc registration function is required")
	}
	if cfg.HTTPRegistrationFunc == nil {
		return nil, fmt.Errorf("http registration function is required")
	}

	grpcListener, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		return nil, err
	}
	httpListener, err := net.Listen("tcp", cfg.HTTPListenAddress)
	if err != nil {
		return nil, err
	}

	s := &ApiServer{
		ctx: cfg.Ctx,
		grpcListener: &proxyproto.Listener{
			Listener:          grpcListener,
			ReadHeaderTimeout: 10 * time.Second,
		},
		httpListener: httpListener,
		log:          cfg.Log.Named("api"),
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(
				[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			),
		),
	)
	cfg.PromRegistry.MustRegister(srvMetrics)

	loggingInterceptor, err := server.NewLoggingInterceptor(cfg.Log)
	if err != nil {
		return nil, err
	}
	openConnectionsInterceptor, err := server.NewOpenConnectionsInterceptor()
	if err != nil {
		return nil, err
	}

	unary := []grpc.UnaryServerInterceptor{
		srvMetrics.UnaryServerInterceptor(),
		openConnectionsInterceptor.Unary(),
		loggingInterceptor.Unary(),
	}
	stream := []grpc.StreamServerInterceptor{
		srvMetrics.StreamServerInterceptor(),
		openConnectionsInterceptor.Stream(),
		loggingInterceptor.Stream(),
	}

	if cfg.JWTVerifier != nil {
		authInterceptor := server.NewAuthInterceptor(cfg.JWTVerifier, cfg.Log)
		unary = append(unary, authInterceptor.Unary())
		stream = append(stream, authInterceptor.Stream())
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
		s.log.Info("enabling gRPC Server Reflection")
	}

	healthgrpc.RegisterHealthServer(s.grpcServer, health.NewServer())

	tracing.GoPanicWrap(s.ctx, &s.wg, "grpc", func(ctx context.Context) {
		s.log.Info("serving grpc", zap.String("address", s.grpcListener.Addr().String()))
		if err := s.grpcServer.Serve(s.grpcListener); err != nil &&
			!isErrUseOfClosedConnection(err) {
			s.log.Error("serving grpc", zap.Error(err))
		}
	})

	if err := s.startHTTP(cfg.Ctx, cfg.Log, cfg.HTTPRegistrationFunc); err != nil {
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

// Package api implements the API server.
package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type APIServerConfig struct {
	Ctx                context.Context
	Logger             *zap.Logger
	PromRegistry       *prometheus.Registry
	RegistrationFunc   RegistrationFunc
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	Port               int
	EnableReflection   bool
}

type APIServerOption func(*APIServerConfig)

type RegistrationFunc func(mux *http.ServeMux) error

func WithContext(ctx context.Context) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Ctx = ctx }
}

func WithLogger(logger *zap.Logger) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Logger = logger }
}

func WithPort(port int) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Port = port }
}

func WithPrometheusRegistry(reg *prometheus.Registry) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.PromRegistry = reg }
}

func WithReflection(enabled bool) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.EnableReflection = enabled }
}

func WithRegistrationFunc(registrationFunc RegistrationFunc) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.RegistrationFunc = registrationFunc }
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.StreamInterceptors = append(cfg.StreamInterceptors, interceptors...) }
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.UnaryInterceptors = append(cfg.UnaryInterceptors, interceptors...) }
}

type APIServer struct {
	ctx        context.Context
	wg         sync.WaitGroup
	httpServer *http.Server
	logger     *zap.Logger
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

	if cfg.Port == 0 {
		return nil, fmt.Errorf("port is required")
	}

	if cfg.RegistrationFunc == nil {
		return nil, fmt.Errorf("registration function is required")
	}

	svc := &APIServer{
		ctx:    cfg.Ctx,
		logger: cfg.Logger.Named(utils.APILoggerName),
	}

	// Create a new HTTP mux to serve the API handlers.
	mux := http.NewServeMux()

	// Wrap the handler with h2c to support HTTP/2 Cleartext for gRPC reflection.
	// This is required for gRPC reflection to work with HTTP/2, and tools such as grpcurl.
	h2cHandler := h2c.NewHandler(mux, &http2.Server{})

	svc.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: h2cHandler,
	}

	svc.logger.Info("creating api server")

	loggingInterceptor, err := server.NewLoggingInterceptor(svc.logger)
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

	// TODO: Fix!
	// Extend server interceptors properly.
	if cfg.PromRegistry != nil {
		srvMetrics := grpcprom.NewServerMetrics(
			grpcprom.WithServerHandlingTimeHistogram(
				grpcprom.WithHistogramBuckets(
					[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
				),
			),
		)
		cfg.PromRegistry.MustRegister(srvMetrics)
		// Prepend metrics interceptors to the chain
		unary = append(
			[]grpc.UnaryServerInterceptor{srvMetrics.UnaryServerInterceptor()},
			unary...,
		)
		stream = append(
			[]grpc.StreamServerInterceptor{srvMetrics.StreamServerInterceptor()},
			stream...,
		)
	}

	// Note: The interceptor chains (unary, stream) are currently not used as the gRPC server
	// implementation is commented out. These may be reactivated in the future or removed.
	_ = unary
	_ = stream

	// svc.grpcServer = grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(unary...),
	// 	grpc.ChainStreamInterceptor(stream...),
	// 	grpc.Creds(insecure.NewCredentials()),
	// 	grpc.KeepaliveParams(keepalive.ServerParameters{Time: 5 * time.Minute}),
	// 	grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
	// 		PermitWithoutStream: true,
	// 		MinTime:             15 * time.Second,
	// 	}),
	// )

	if err := svc.registerBaseAPIServerHandlers(mux, cfg.EnableReflection); err != nil {
		return nil, err
	}

	if err := cfg.RegistrationFunc(mux); err != nil {
		return nil, err
	}

	return svc, nil
}

func (svc *APIServer) Start() {
	svc.logger.Info("starting api server", zap.String("address", svc.httpServer.Addr))

	tracing.GoPanicWrap(svc.ctx, &svc.wg, "api-server", func(ctx context.Context) {
		if err := svc.httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			svc.logger.Fatal("error serving api server", zap.Error(err))
		}
	})
}

func (svc *APIServer) Addr() string {
	return svc.httpServer.Addr
}

func (svc *APIServer) Close() {
	svc.logger.Info("stopping api server")

	// Create a context with timeout for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the HTTP server.
	if err := svc.httpServer.Shutdown(shutdownCtx); err != nil {
		svc.logger.Error("error shutting down api server", zap.Error(err))
	}

	// Wait for the goroutine to finish.
	svc.wg.Wait()

	svc.logger.Info("api server stopped")
}

func (svc *APIServer) registerBaseAPIServerHandlers(
	mux *http.ServeMux,
	enableReflection bool,
) error {
	svc.registerHealthHandler(mux)

	if enableReflection {
		svc.registerReflectionHandlerV1(mux)
		svc.registerReflectionHandlerV1Alpha(mux)
	}

	return nil
}

func (svc *APIServer) registerHealthHandler(mux *http.ServeMux) {
	healthChecker := grpchealth.NewStaticChecker(
		message_apiconnect.ReplicationApiName,
		metadata_apiconnect.MetadataApiName,
	)

	path, handler := grpchealth.NewHandler(healthChecker)

	mux.Handle(path, handler)

	svc.logger.Info("health handler registered")
}

// TODO: Fix!
// Implement dynamic reflector.
func reflector() *grpcreflect.Reflector {
	return grpcreflect.NewStaticReflector(
		grpchealth.HealthV1ServiceName,
	)
}

func (svc *APIServer) registerReflectionHandlerV1(mux *http.ServeMux) {
	reflector := reflector()

	path, handler := grpcreflect.NewHandlerV1(reflector)

	mux.Handle(path, handler)

	svc.logger.Info("reflection handler v1 registered")
}

func (svc *APIServer) registerReflectionHandlerV1Alpha(mux *http.ServeMux) {
	reflector := reflector()

	path, handler := grpcreflect.NewHandlerV1Alpha(reflector)

	mux.Handle(path, handler)

	svc.logger.Info("reflection handler v1 alpha registered")
}

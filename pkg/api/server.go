// Package api implements the API server.
package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	interceptors "github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type APIServerConfig struct {
	Ctx              context.Context
	Logger           *zap.Logger
	PromRegistry     *prometheus.Registry
	RegistrationFunc RegistrationFunc
	Listener         net.Listener
	EnableReflection bool
}

type APIServerOption func(*APIServerConfig)

type RegistrationFunc func(mux *http.ServeMux, interceptors ...connect.Interceptor) (servicePaths []string, err error)

func WithContext(ctx context.Context) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Ctx = ctx }
}

func WithListener(listener net.Listener) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Listener = listener }
}

func WithLogger(logger *zap.Logger) APIServerOption {
	return func(cfg *APIServerConfig) { cfg.Logger = logger }
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

type APIServer struct {
	ctx        context.Context
	wg         sync.WaitGroup
	logger     *zap.Logger
	httpServer *http.Server
	listener   net.Listener
}

func NewAPIServer(opts ...APIServerOption) (*APIServer, error) {
	cfg := &APIServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Ctx == nil {
		return nil, errors.New("context is required")
	}

	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}

	if cfg.Listener == nil {
		return nil, errors.New("listener is required")
	}

	if cfg.RegistrationFunc == nil {
		return nil, errors.New("registration function is required")
	}

	svc := &APIServer{
		ctx:      cfg.Ctx,
		logger:   cfg.Logger.Named(utils.APILoggerName),
		listener: cfg.Listener,
	}

	// Create a new HTTP mux to serve the API handlers.
	mux := http.NewServeMux()

	// Wrap the handler with h2c to support HTTP/2 Cleartext for gRPC reflection.
	// This is required for gRPC reflection to work with HTTP/2, and tools such as grpcurl.
	h2cHandler := h2c.NewHandler(mux, &http2.Server{
		IdleTimeout: 5 * time.Minute,
	})

	// TODO: Do we need more timeouts?
	svc.httpServer = &http.Server{
		Handler:     handleCORS(h2cHandler),
		IdleTimeout: 5 * time.Minute,
	}

	svc.logger.Info("creating api server")

	// Create protocol validation interceptor.
	protocolValidationInterceptor := interceptors.NewProtocolValidationInterceptor()

	// Create server side interceptors.
	openConnInterceptor, err := interceptors.NewOpenConnectionsInterceptor()
	if err != nil {
		return nil, err
	}

	loggingInterceptor, err := interceptors.NewLoggingInterceptor(svc.logger)
	if err != nil {
		return nil, err
	}

	grpcMetricsInterceptor := interceptors.NewGRPCMetricsInterceptor()

	// Register services.
	servicePaths, err := cfg.RegistrationFunc(
		mux,
		protocolValidationInterceptor,
		grpcMetricsInterceptor,
		openConnInterceptor,
		loggingInterceptor,
	)
	if err != nil {
		return nil, err
	}

	// Register health handler.
	svc.registerHealthHandler(mux, servicePaths)

	// Register reflection handlers.
	if cfg.EnableReflection {
		svc.registerReflectionHandlers(mux, servicePaths)
	}

	return svc, nil
}

func (svc *APIServer) Start() {
	tracing.GoPanicWrap(svc.ctx, &svc.wg, "api-server", func(ctx context.Context) {
		if err := svc.httpServer.Serve(svc.listener); err != nil &&
			err != http.ErrServerClosed {
			svc.logger.Fatal("error serving api server", zap.Error(err))
		}
	})

	svc.logger.Info("started api server", zap.String("address", svc.Addr()))
}

func (svc *APIServer) Addr() string {
	return svc.listener.Addr().String()
}

func (svc *APIServer) Close(timeout time.Duration) {
	svc.logger.Info("closing")

	// Create a context with timeout for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Gracefully shutdown the HTTP server.
	if err := svc.httpServer.Shutdown(shutdownCtx); err != nil {
		svc.logger.Error("error shutting down api server", zap.Error(err))
	}

	// Wait for the goroutine to finish.
	svc.wg.Wait()

	svc.logger.Info("closed")
}

func (svc *APIServer) registerHealthHandler(
	mux *http.ServeMux,
	servicePaths []string,
) {
	healthChecker := grpchealth.NewStaticChecker(
		append(servicePaths, grpchealth.HealthV1ServiceName)...,
	)

	path, handler := grpchealth.NewHandler(healthChecker)

	mux.Handle(path, handler)

	svc.logger.Info("health handler registered")
}

func (svc *APIServer) registerReflectionHandlers(mux *http.ServeMux, servicePaths []string) {
	reflector := grpcreflect.NewStaticReflector(
		append(servicePaths, grpchealth.HealthV1ServiceName)...,
	)

	pathV1, handlerV1 := grpcreflect.NewHandlerV1(reflector)

	mux.Handle(pathV1, handlerV1)

	svc.logger.Info("reflection handler v1 registered")

	pathV1Alpha, handlerV1Alpha := grpcreflect.NewHandlerV1Alpha(reflector)

	mux.Handle(pathV1Alpha, handlerV1Alpha)

	svc.logger.Info("reflection handler v1 alpha registered")
}

func handleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Expose gRPC/gRPC-Web status headers to browsers.
		w.Header().Set("Access-Control-Expose-Headers",
			"grpc-status,grpc-message,grpc-status-details-bin")

		// Cache CORS preflight requests for 24 hours.
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle OPTIONS preflight requests.
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			handleCORSOptionsPreflight(w, r)
			return
		}

		// Serve the request.
		h.ServeHTTP(w, r)
	})
}

func handleCORSOptionsPreflight(w http.ResponseWriter, _ *http.Request) {
	// Allowed headers.
	headers := []string{
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Client-Version",
		"X-App-Version",
		"Baggage",
		"DNT",
		"Sec-CH-UA",
		"Sec-CH-UA-Mobile",
		"Sec-CH-UA-Platform",
		"x-grpc-web",
		"grpc-timeout",
		"Sentry-Trace",
		"User-Agent",
		"x-libxmtp-version",
		"x-app-version",
	}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	// Allowed methods.
	methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

	// No content is the correct response for an OPTIONS preflight request.
	w.WriteHeader(http.StatusNoContent)
}

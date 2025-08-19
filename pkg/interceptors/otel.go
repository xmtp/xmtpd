// Package interceptors provides gRPC interceptor wrappers with OpenTelemetry support
package interceptors

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/stats"
)

// OTelServerInterceptor provides gRPC server instrumentation with OpenTelemetry tracing
type OTelServerInterceptor struct {
	logger  *zap.Logger
	handler stats.Handler
}

// NewOTelServerInterceptor creates a new server interceptor with OTEL tracing
func NewOTelServerInterceptor(logger *zap.Logger) *OTelServerInterceptor {
	return &OTelServerInterceptor{
		logger:  logger.Named("otel-server-interceptor"),
		handler: otelgrpc.NewServerHandler(),
	}
}

// Handler returns the stats.Handler for gRPC server instrumentation
// Use this with grpc.StatsHandler() when creating your gRPC server
func (i *OTelServerInterceptor) Handler() stats.Handler {
	return i.handler
}

// OTelClientInterceptor provides gRPC client instrumentation with OpenTelemetry tracing
type OTelClientInterceptor struct {
	logger  *zap.Logger
	handler stats.Handler
}

// NewOTelClientInterceptor creates a new client interceptor with OTEL tracing
func NewOTelClientInterceptor(logger *zap.Logger) *OTelClientInterceptor {
	return &OTelClientInterceptor{
		logger:  logger.Named("otel-client-interceptor"),
		handler: otelgrpc.NewClientHandler(),
	}
}

// Handler returns the stats.Handler for gRPC client instrumentation
// Use this with grpc.WithStatsHandler() when creating your gRPC client connection
func (i *OTelClientInterceptor) Handler() stats.Handler {
	return i.handler
}

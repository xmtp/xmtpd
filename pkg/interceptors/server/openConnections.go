package server

import (
	"context"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"google.golang.org/grpc"
)

// OpenConnectionsInterceptor reports open connections for unary and stream RPCs.
type OpenConnectionsInterceptor struct {
}

// NewOpenConnectionsInterceptor creates a new instance of OpenConnectionsInterceptor.
func NewOpenConnectionsInterceptor() (*OpenConnectionsInterceptor, error) {
	return &OpenConnectionsInterceptor{}, nil
}

// Unary intercepts unary RPC calls to log errors.
func (i *OpenConnectionsInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		oc := metrics.NewApiOpenConnection("unary", info.FullMethod)
		defer oc.Close()
		return handler(ctx, req) // Call the actual RPC handler
	}
}

// Stream intercepts stream RPC calls to log errors.
func (i *OpenConnectionsInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		oc := metrics.NewApiOpenConnection("stream", info.FullMethod)
		defer oc.Close()
		return handler(srv, ss) // Call the actual stream handler
	}
}

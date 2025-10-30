package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"google.golang.org/grpc"
)

// OpenConnectionsInterceptor reports open connections for unary and stream RPCs.
type OpenConnectionsInterceptor struct{}

// Connect-go interceptor interface implementation.
var _ connect.Interceptor = (*OpenConnectionsInterceptor)(nil)

// NewOpenConnectionsInterceptor creates a new instance of OpenConnectionsInterceptor.
func NewOpenConnectionsInterceptor() (*OpenConnectionsInterceptor, error) {
	return &OpenConnectionsInterceptor{}, nil
}

/* gRPC interceptors */

func (i *OpenConnectionsInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		oc := metrics.NewAPIOpenConnection("unary", info.FullMethod)
		defer oc.Close()
		return handler(ctx, req)
	}
}

func (i *OpenConnectionsInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		oc := metrics.NewAPIOpenConnection("stream", info.FullMethod)
		defer oc.Close()
		return handler(srv, stream)
	}
}

/* Connect-go interceptors */

func (i *OpenConnectionsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		oc := metrics.NewAPIOpenConnection("unary", req.Spec().Procedure)
		defer oc.Close()
		return next(ctx, req)
	}
}

// WrapStreamingClient is a no-op. Interface requirement.
func (i *OpenConnectionsInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *OpenConnectionsInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		oc := metrics.NewAPIOpenConnection("stream", conn.Spec().Procedure)
		defer oc.Close()
		return next(ctx, conn)
	}
}

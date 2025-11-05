package server

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// ProtocolValidationInterceptor enforces protocol restrictions by rejecting
// the Connect protocol and allowing only gRPC and gRPC-Web.
type ProtocolValidationInterceptor struct{}

// Compile-time check that ProtocolValidationInterceptor implements connect.Interceptor.
var _ connect.Interceptor = (*ProtocolValidationInterceptor)(nil)

var errUnsupportedProtocol = errors.New(
	"Connect-RPC protocol not supported, use gRPC or gRPC-Web",
)

func NewProtocolValidationInterceptor() *ProtocolValidationInterceptor {
	return &ProtocolValidationInterceptor{}
}

func (i *ProtocolValidationInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := validateProtocol(req.Peer().Protocol); err != nil {
			return nil, err
		}

		return next(ctx, req)
	}
}

// WrapStreamingClient is a no-op. Interface requirement.
func (i *ProtocolValidationInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *ProtocolValidationInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		if err := validateProtocol(conn.Peer().Protocol); err != nil {
			return err
		}

		return next(ctx, conn)
	}
}

func validateProtocol(protocol string) error {
	if protocol == connect.ProtocolConnect {
		return connect.NewError(
			connect.CodeFailedPrecondition,
			errUnsupportedProtocol,
		)
	}

	return nil
}

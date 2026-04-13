package server

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/constants"
)

// RequireNodeAuthInterceptor rejects requests that do not carry a verified
// node JWT. It is intended to be wrapped around ReplicationApi only, gating
// all methods behind Tier 0 authentication.
//
// The upstream ServerAuthInterceptor sets VerifiedNodeRequestCtxKey when a
// valid node JWT is present. This interceptor checks that key and returns
// Unauthenticated if it is missing or false.
type RequireNodeAuthInterceptor struct {
	logger *zap.Logger
}

var _ connect.Interceptor = (*RequireNodeAuthInterceptor)(nil)

func NewRequireNodeAuthInterceptor(logger *zap.Logger) *RequireNodeAuthInterceptor {
	return &RequireNodeAuthInterceptor{
		logger: logger.Named("require-node-auth"),
	}
}

func (i *RequireNodeAuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if !isVerifiedNode(ctx) {
			i.logger.Warn(
				"unauthenticated replication request rejected",
				zap.String("procedure", req.Spec().Procedure),
			)
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("node authentication required"),
			)
		}
		return next(ctx, req)
	}
}

func (i *RequireNodeAuthInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *RequireNodeAuthInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if !isVerifiedNode(ctx) {
			i.logger.Warn(
				"unauthenticated replication stream rejected",
				zap.String("procedure", conn.Spec().Procedure),
			)
			return connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("node authentication required"),
			)
		}
		return next(ctx, conn)
	}
}

func isVerifiedNode(ctx context.Context) bool {
	v, ok := ctx.Value(constants.VerifiedNodeRequestCtxKey{}).(bool)
	return ok && v
}

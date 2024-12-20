package server

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// wrappedServerStream allows us to modify the context of the stream
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// AuthInterceptor validates JWT tokens from other nodes
type AuthInterceptor struct {
	verifier authn.JWTVerifier
	logger   *zap.Logger
}

func NewAuthInterceptor(verifier authn.JWTVerifier, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		verifier: verifier,
		logger:   logger,
	}
}

// extractToken gets the JWT token from the request metadata
func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get(constants.NODE_AUTHORIZATION_HEADER_NAME)
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing auth token")
	}

	if len(values) > 1 {
		return "", status.Error(codes.Unauthenticated, "multiple auth tokens provided")
	}

	return values[0], nil
}

// Unary returns a grpc.UnaryServerInterceptor that validates JWT tokens
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, err := extractToken(ctx)
		if err != nil {
			i.logger.Debug("failed to find auth token. Allowing request to proceed", zap.Error(err))
			return handler(ctx, req)
		}

		if err := i.verifier.Verify(token); err != nil {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"invalid auth token: %v",
				err,
			)
		}
		ctx = context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)

		return handler(ctx, req)
	}
}

// Stream returns a grpc.StreamServerInterceptor that validates JWT tokens
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		token, err := extractToken(stream.Context())
		if err != nil {
			i.logger.Debug("failed to find auth token. Allowing request to proceed", zap.Error(err))
			return handler(srv, stream)
		}

		if err := i.verifier.Verify(token); err != nil {
			return status.Errorf(
				codes.Unauthenticated,
				"invalid auth token: %v",
				err,
			)
		}

		stream = &wrappedServerStream{
			ServerStream: stream,
			ctx: context.WithValue(
				stream.Context(),
				constants.VerifiedNodeRequestCtxKey{},
				true,
			),
		}

		return handler(srv, stream)
	}
}

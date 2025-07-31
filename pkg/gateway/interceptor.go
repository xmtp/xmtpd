package gateway

import (
	"context"
	"strings"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type identityCtxKey struct{}

type gatewayWrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *gatewayWrappedServerStream) Context() context.Context {
	return w.ctx
}

type GatewayInterceptor struct {
	identityFn  IdentityFn
	authorizers []AuthorizePublishFn
	logger      *zap.Logger
}

func NewGatewayInterceptor(
	logger *zap.Logger,
	identityFn IdentityFn,
	authorizers []AuthorizePublishFn,
) *GatewayInterceptor {
	return &GatewayInterceptor{
		logger:      logger,
		identityFn:  identityFn,
		authorizers: authorizers,
	}
}

func (i *GatewayInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		identity, err := i.identityFn(ctx)
		if err != nil {
			i.logger.Error("failed to get identity", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to get identity: %v", err)
		}

		ctx = context.WithValue(ctx, identityCtxKey{}, identity)

		if strings.HasSuffix(info.FullMethod, "PublishClientEnvelopes") {
			publishReq, ok := req.(*payer_api.PublishClientEnvelopesRequest)
			if !ok {
				return nil, status.Error(codes.Internal, "invalid request type")
			}

			// Create a summary of the request
			summary := PublishRequestSummary{
				TotalEnvelopes: len(publishReq.Envelopes),
				// TODO: Calculate cost estimates
			}

			for _, authorizer := range i.authorizers {
				authorized, err := authorizer(ctx, identity, summary)
				if err != nil {
					i.logger.Error("authorization error",
						zap.Error(err),
						zap.String("identity", identity.Identity),
						zap.String("identity_kind", string(identity.Kind)))
					return nil, status.Errorf(codes.Internal, "authorization error: %v", err)
				}
				if !authorized {
					i.logger.Warn("unauthorized publish request",
						zap.String("identity", identity.Identity),
						zap.String("identity_kind", string(identity.Kind)))
					return nil, status.Error(codes.PermissionDenied, "unauthorized")
				}
			}
		}

		return handler(ctx, req)
	}
}

func (i *GatewayInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		identity, err := i.identityFn(stream.Context())
		if err != nil {
			i.logger.Error("failed to get identity", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to get identity: %v", err)
		}

		stream = &gatewayWrappedServerStream{
			ServerStream: stream,
			ctx:          context.WithValue(stream.Context(), identityCtxKey{}, identity),
		}

		return handler(srv, stream)
	}
}

func GetIdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(identityCtxKey{}).(Identity)
	return identity, ok
}

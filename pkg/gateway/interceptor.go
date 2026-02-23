package gateway

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"go.uber.org/zap"
)

const (
	identityField     = "identity"
	identityKindField = "identity_kind"
)

type identityCtxKey struct{}

// GatewayInterceptor is the server-side interceptor for the gateway API.
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

func (i *GatewayInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		identity, err := i.identityFn(req.Header(), req.Peer().Addr)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("failed to get identity: %w", err),
			)
		}

		ctx = context.WithValue(ctx, identityCtxKey{}, identity)

		if strings.HasSuffix(req.Spec().Procedure, "PublishClientEnvelopes") {
			publishReq, ok := req.Any().(*payer_api.PublishClientEnvelopesRequest)
			if !ok {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("invalid request type"),
				)
			}

			// Create a summary of the request
			summary := PublishRequestSummary{
				TotalEnvelopes: len(publishReq.GetEnvelopes()),
				// TODO: Calculate cost estimates
			}

			for _, authorizer := range i.authorizers {
				authorized, err := authorizer(ctx, identity, summary)
				if err != nil {
					var rlError GatewayServiceError
					if errors.As(err, &rlError) {
						i.logger.Info("request rejected", zap.Error(err))
						return nil, returnRetryAfterError(rlError)
					}

					i.logger.Warn("authorization error",
						zap.Error(err),
						zap.String(identityField, identity.Identity),
						zap.String(identityKindField, string(identity.Kind)))
					return nil, connect.NewError(
						connect.CodeInternal,
						fmt.Errorf("authorization error: %w", err),
					)
				}

				if !authorized {
					i.logger.Warn("unauthorized publish request",
						zap.String(identityField, identity.Identity),
						zap.String(identityKindField, string(identity.Kind)))
					return nil, connect.NewError(
						connect.CodePermissionDenied,
						errors.New("unauthorized"),
					)
				}
			}
		}

		return next(ctx, req)
	}
}

// WrapStreamingClient is a no-op. Interface requirement.
func (i *GatewayInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *GatewayInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		identity, err := i.identityFn(conn.RequestHeader(), conn.Peer().Addr)
		if err != nil {
			i.logger.Error("failed to get identity", zap.Error(err))
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("failed to get identity: %w", err),
			)
		}

		// TODO: Check if the identity is authorized to publish to the gateway.

		ctx = context.WithValue(ctx, identityCtxKey{}, identity)

		return next(ctx, conn)
	}
}

func returnRetryAfterError(rlError GatewayServiceError) error {
	connectErr := connect.NewError(rlError.Code(), rlError)

	retryAfter := rlError.RetryAfter()
	if retryAfter != nil {
		retryAfterValue := strconv.Itoa(int(retryAfter.Seconds()))
		connectErr.Meta().Set("Retry-After", retryAfterValue)
	}

	return connectErr
}

func GetIdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(identityCtxKey{}).(Identity)
	return identity, ok
}

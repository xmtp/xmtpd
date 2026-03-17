// Package server implements the server authentication interceptors.
// It validates JWT tokens from other nodes.
package server

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// TODO(borja): Next PR - Fail requests if the token is not valid.

// ServerAuthInterceptor validates JWT tokens from other nodes
type ServerAuthInterceptor struct {
	verifier authn.JWTVerifier
	logger   *zap.Logger
}

var _ connect.Interceptor = (*ServerAuthInterceptor)(nil)

// NewServerAuthInterceptor creates a new ServerAuthInterceptor.
func NewServerAuthInterceptor(
	verifier authn.JWTVerifier,
	logger *zap.Logger,
) *ServerAuthInterceptor {
	return &ServerAuthInterceptor{
		verifier: verifier,
		logger:   logger,
	}
}

func (i *ServerAuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token := req.Header().Get(constants.NodeAuthorizationHeaderName)
		// If token is missing, allow the request to proceed without authentication.
		// Handlers must check VerifiedNodeRequestCtxKey if authentication is required.
		if token == "" {
			return next(ctx, req)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			logFields := []zap.Field{
				zap.String("procedure", req.Spec().Procedure),
				zap.String("protocol", req.Peer().Protocol),
				zap.Error(err),
			}
			if id := tryExtractNodeIDFromToken(token); id != 0 {
				logFields = append(logFields, utils.OriginatorIDField(id))
			}
			i.logger.Error("JWT verification failed", logFields...)

			// Do not expose too much information to the client (e.g. wrapped errors)
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("invalid auth token"),
			)
		}
		defer cancel()

		i.connectLogIncomingAddress(nodeID)

		ctx = context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)

		return next(ctx, req)
	}
}

// WrapStreamingClient is a no-op for server interceptors.
// It's only implemented to satisfy the connect.Interceptor interface.
// This method is never called on the server side.
func (i *ServerAuthInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *ServerAuthInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		token := conn.RequestHeader().Get(constants.NodeAuthorizationHeaderName)
		// If token is missing, allow the request to proceed without authentication.
		// Handlers must check VerifiedNodeRequestCtxKey if authentication is required.
		if token == "" {
			return next(ctx, conn)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			logFields := []zap.Field{
				zap.String("procedure", conn.Spec().Procedure),
				zap.String("protocol", conn.Peer().Protocol),
				zap.Error(err),
			}
			if id := tryExtractNodeIDFromToken(token); id != 0 {
				logFields = append(logFields, utils.OriginatorIDField(id))
			}
			i.logger.Error("JWT verification failed", logFields...)

			// Do not expose too much information to the client (e.g. wrapped errors)
			return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid auth token"))
		}
		defer cancel()

		i.connectLogIncomingAddress(nodeID)

		ctx = context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)

		return next(ctx, conn)
	}
}

/* Connect-go interceptors helpers */

func (i *ServerAuthInterceptor) connectLogIncomingAddress(nodeID uint32) {
	if i.logger.Core().Enabled(zap.DebugLevel) {
		i.logger.Debug("incoming connection", utils.OriginatorIDField(nodeID))
	}
}

// tryExtractNodeIDFromToken parses the JWT subject claim without signature verification.
// Returns 0 if the subject cannot be extracted or parsed as a node ID.
func tryExtractNodeIDFromToken(tokenString string) uint32 {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0
	}
	parsed, err := strconv.ParseInt(subject, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(parsed)
}

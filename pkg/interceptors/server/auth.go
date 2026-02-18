// Package server implements the server authentication interceptors.
// It validates JWT tokens from other nodes and logs the incoming address.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"connectrpc.com/connect"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var defaultConfig = ServerAuthConfig{
	RequireToken: false,
	DNSLookup:    true,
}

// TODO(borja): Next PR - Fail requests if the token is not valid.

const (
	dnsNameField       = "dns_name"
	clientAddressField = "client_address"
)

// ServerAuthInterceptor validates JWT tokens from other nodes
type ServerAuthInterceptor struct {
	cfg ServerAuthConfig

	verifier authn.JWTVerifier
	logger   *zap.Logger
}

var _ connect.Interceptor = (*ServerAuthInterceptor)(nil)

type ServerAuthConfig struct {
	// Requests without a token should be rejected.
	RequireToken bool

	// Do not perform DNS lookup for logging requests.
	DNSLookup bool
}

type ServerAuthOption func(*ServerAuthConfig)

func RequireToken(b bool) ServerAuthOption {
	return func(cfg *ServerAuthConfig) {
		cfg.RequireToken = b
	}
}

func DoDNSLookup(b bool) ServerAuthOption {
	return func(cfg *ServerAuthConfig) {
		cfg.DNSLookup = b
	}
}

// NewServerAuthInterceptor creates a new ServerAuthInterceptor.
func NewServerAuthInterceptor(
	logger *zap.Logger,
	verifier authn.JWTVerifier,
	opts ...ServerAuthOption,
) *ServerAuthInterceptor {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &ServerAuthInterceptor{
		cfg:      cfg,
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

			if i.cfg.RequireToken {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("missing auth token"),
				)
			}

			return next(ctx, req)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("invalid auth token: %w", err),
			)
		}
		defer cancel()

		i.connectLogIncomingAddress(req.Peer().Addr, nodeID)

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

			if i.cfg.RequireToken {
				return connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("missing auth token"),
				)
			}

			return next(ctx, conn)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			return connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("invalid auth token: %w", err),
			)
		}
		defer cancel()

		i.connectLogIncomingAddress(conn.Peer().Addr, nodeID)

		ctx = context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)

		return next(ctx, conn)
	}
}

/* Connect-go interceptors helpers */

func (i *ServerAuthInterceptor) connectLogIncomingAddress(
	addr string,
	nodeID uint32,
) {
	if !i.logger.Core().Enabled(zap.DebugLevel) {
		return
	}

	// Do not do costly DNS lookup if not necessary.
	if !i.cfg.DNSLookup {
		i.logger.Debug("incoming connection",
			zap.String(clientAddressField, addr),
			utils.OriginatorIDField(nodeID))

		return
	}

	// TODO: Potentially cache these values.
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// Do nothing.
		return
	}

	dnsName, err := net.LookupAddr(host)
	if err != nil || len(dnsName) == 0 {
		dnsName = []string{unknownDNSName}
	}

	i.logger.Debug(
		"incoming connection",
		zap.String(clientAddressField, addr),
		zap.String(dnsNameField, dnsName[0]),
		utils.OriginatorIDField(nodeID),
	)
}

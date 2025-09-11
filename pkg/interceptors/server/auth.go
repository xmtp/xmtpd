// Package server implements the server for the interceptors package.
package server

import (
	"context"
	"net"

	"google.golang.org/grpc/peer"

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

	values := md.Get(constants.NodeAuthorizationHeaderName)
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing auth token")
	}

	if len(values) > 1 {
		return "", status.Error(codes.Unauthenticated, "multiple auth tokens provided")
	}

	return values[0], nil
}

func (i *AuthInterceptor) logIncomingAddressIfAvailable(ctx context.Context, nodeID uint32) {
	if i.logger.Core().Enabled(zap.DebugLevel) {
		if p, ok := peer.FromContext(ctx); ok {
			clientAddr := p.Addr.String()
			var dnsName []string
			// Attempt to resolve the DNS name
			host, _, err := net.SplitHostPort(clientAddr)
			if err == nil {
				dnsName, err = net.LookupAddr(host)
				if err != nil || len(dnsName) == 0 {
					dnsName = []string{"Unknown"}
				}
			} else {
				dnsName = []string{"Unknown"}
			}
			i.logger.Debug(
				"Incoming connection",
				zap.String("client_addr", clientAddr),
				zap.String("dns_name", dnsName[0]),
				zap.Uint32("node_id", nodeID),
			)
		}
	}
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
			return handler(ctx, req)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"invalid auth token: %v",
				err,
			)
		}
		defer cancel()

		i.logIncomingAddressIfAvailable(ctx, nodeID)

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
			return handler(srv, stream)
		}

		nodeID, cancel, err := i.verifier.Verify(token)
		if err != nil {
			return status.Errorf(
				codes.Unauthenticated,
				"invalid auth token: %v",
				err,
			)
		}
		defer cancel()

		i.logIncomingAddressIfAvailable(stream.Context(), nodeID)

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

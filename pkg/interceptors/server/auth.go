// Package server implements the server authentication interceptors.
// It validates JWT tokens from other nodes and logs the incoming address.
package server

import (
	"context"
	"fmt"
	"net"

	"connectrpc.com/connect"
	"google.golang.org/grpc/peer"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TODO(borja): Next PR - Fail requests if the token is not valid.

const (
	dnsNameField       = "dns_name"
	clientAddressField = "client_address"
)

// wrappedServerStream allows us to modify the context of the stream
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// ServerAuthInterceptor validates JWT tokens from other nodes
type ServerAuthInterceptor struct {
	verifier authn.JWTVerifier
	logger   *zap.Logger
}

// NewServerAuthInterceptor creates a new ServerAuthInterceptor.
// Supports gRPC and connect-go interceptors.
func NewServerAuthInterceptor(
	verifier authn.JWTVerifier,
	logger *zap.Logger,
) *ServerAuthInterceptor {
	return &ServerAuthInterceptor{
		verifier: verifier,
		logger:   logger,
	}
}

/* gRPC interceptors */

// Unary returns a gRPC unary server interceptor that validates JWT tokens.
func (i *ServerAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, err := grpcTokenFromContext(ctx)
		// If token extraction fails, allow the request to proceed without authentication.
		// Handlers must check VerifiedNodeRequestCtxKey if authentication is required.
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

		i.grpcLogIncomingAddress(ctx, nodeID)

		ctx = context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)

		return handler(ctx, req)
	}
}

// Stream returns a gRPC stream server interceptor that validates JWT tokens.
func (i *ServerAuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		token, err := grpcTokenFromContext(stream.Context())
		// If token extraction fails, allow the request to proceed without authentication.
		// Handlers must check VerifiedNodeRequestCtxKey if authentication is required.
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

		i.grpcLogIncomingAddress(stream.Context(), nodeID)

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

/* gRPC interceptors helpers */

// grpcTokenFromContext gets the JWT token from the request metadata.
// This method is used by gRPC interceptors, as setting/getting the token
// via the context is not possible with connect-go.
func grpcTokenFromContext(ctx context.Context) (string, error) {
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

// grpcLogIncomingAddress logs the incoming address if available.
// This method is used by gRPC interceptors, as setting/getting the address
// via the context is not possible with connect-go.
func (i *ServerAuthInterceptor) grpcLogIncomingAddress(ctx context.Context, nodeID uint32) {
	if i.logger.Core().Enabled(zap.DebugLevel) {
		if p, ok := peer.FromContext(ctx); ok {
			clientAddr := p.Addr.String()
			var dnsName []string
			// Attempt to resolve the DNS name
			host, _, err := net.SplitHostPort(clientAddr)
			if err == nil {
				dnsName, err = net.LookupAddr(host)
				if err != nil || len(dnsName) == 0 {
					dnsName = []string{unknownDNSName}
				}
			} else {
				dnsName = []string{unknownDNSName}
			}
			i.logger.Debug(
				"incoming connection",
				zap.String(clientAddressField, clientAddr),
				zap.String(dnsNameField, dnsName[0]),
				utils.OriginatorIDField(nodeID),
			)
		}
	}
}

/* Connect-go interceptors */

var _ connect.Interceptor = (*ServerAuthInterceptor)(nil)

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
	if i.logger.Core().Enabled(zap.DebugLevel) {
		host, _, err := net.SplitHostPort(addr)
		if err == nil {
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
	}
}

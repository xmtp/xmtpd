// Package client implements the authentication interceptors for the client.
// The auth interceptors are used mostly for gRPC clients.
// We also have support for connect-go clients for future proofing.
package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ClientAuthInterceptor is a struct for holding the token and adding it to each request.
type ClientAuthInterceptor struct {
	tokenFactory authn.TokenFactory
	targetNodeID uint32
	currentToken *authn.Token
	mu           sync.RWMutex
}

// Ensure AuthInterceptor implements connect.Interceptor for client-side usage.
var _ connect.Interceptor = (*ClientAuthInterceptor)(nil)

// NewClientAuthInterceptor creates a new AuthInterceptor.
func NewClientAuthInterceptor(
	tokenFactory authn.TokenFactory,
	targetNodeID uint32,
) *ClientAuthInterceptor {
	return &ClientAuthInterceptor{
		tokenFactory: tokenFactory,
		targetNodeID: targetNodeID,
	}
}

func (i *ClientAuthInterceptor) getToken() (*authn.Token, error) {
	// Fast path: check with read lock if current token is still valid.
	i.mu.RLock()
	if i.currentToken != nil &&
		i.currentToken.ExpiresAt.After(time.Now().Add(authn.MaxClockSkew)) {
		token := i.currentToken
		i.mu.RUnlock()
		return token, nil
	}
	i.mu.RUnlock()

	// Slow path: acquire write lock to create new token.
	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have created it).
	if i.currentToken != nil &&
		i.currentToken.ExpiresAt.After(time.Now().Add(authn.MaxClockSkew)) {
		return i.currentToken, nil
	}

	token, err := i.tokenFactory.CreateToken(i.targetNodeID)
	if err != nil {
		return nil, err
	}

	i.currentToken = token
	return token, nil
}

/* gRPC interceptors */

func (i *ClientAuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		token, err := i.getToken()
		if err != nil {
			return status.Errorf(
				codes.Unauthenticated,
				"failed to get token: %v",
				err,
			)
		}

		// Create the metadata with the token.
		md := metadata.Pairs(constants.NodeAuthorizationHeaderName, token.SignedString)

		// Attach metadata to the outgoing context.
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Proceed with the request.
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *ClientAuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		token, err := i.getToken()
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "failed to get token: %v", err)
		}

		// Create the metadata with the token.
		md := metadata.Pairs(constants.NodeAuthorizationHeaderName, token.SignedString)

		// Attach the metadata to the outgoing context.
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Proceed with the stream.
		return streamer(ctx, desc, cc, method, opts...)
	}
}

/* Connect-go interceptors */

func (i *ClientAuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token, err := i.getToken()
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("failed to get token: %w", err),
			)
		}

		// Set the auth header.
		req.Header().Set(constants.NodeAuthorizationHeaderName, token.SignedString)

		// Call the next handler.
		return next(ctx, req)
	}
}

func (i *ClientAuthInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		token, err := i.getToken()
		if err != nil {
			// If token generation fails, return a connection that will fail on any operation.
			// We still need to call next() to maintain the connection lifecycle,
			// but wrap it to fail immediately.
			return &streamingAuthInterceptorFailure{
				StreamingClientConn: next(ctx, spec),
				err:                 err,
			}
		}

		// Establish the connection and set the auth header.
		conn := next(ctx, spec)

		conn.RequestHeader().Set(constants.NodeAuthorizationHeaderName, token.SignedString)

		return conn
	}
}

// WrapStreamingHandler is a no-op for client interceptors.
// It's only implemented to satisfy the connect.Interceptor interface.
// This method is never called on the client side.
func (i *ClientAuthInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return next
}

// streamingAuthInterceptorFailure is a wrapper around a streaming client connection
// that returns an authentication error on Send() or Receive() operations.
type streamingAuthInterceptorFailure struct {
	connect.StreamingClientConn
	err error
}

func (s *streamingAuthInterceptorFailure) Send(msg any) error {
	return connect.NewError(
		connect.CodeUnauthenticated,
		fmt.Errorf("failed to get authentication token: %w", s.err),
	)
}

func (s *streamingAuthInterceptorFailure) Receive(msg any) error {
	return connect.NewError(
		connect.CodeUnauthenticated,
		fmt.Errorf("failed to get authentication token: %w", s.err),
	)
}

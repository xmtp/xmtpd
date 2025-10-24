// Package client implements the client for the interceptors package.
package client

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a struct for holding the token and adding it to each request.
type AuthInterceptor struct {
	tokenFactory authn.TokenFactory
	targetNodeID uint32
	currentToken *authn.Token
}

// NewAuthInterceptor creates a new AuthInterceptor.
// It's compliant with the grpc.UnaryClientInterceptor and grpc.StreamClientInterceptor interfaces.
// Do not use in connect-go clients.
func NewAuthInterceptor(
	tokenFactory authn.TokenFactory,
	targetNodeID uint32,
) *AuthInterceptor {
	return &AuthInterceptor{
		tokenFactory: tokenFactory,
		targetNodeID: targetNodeID,
	}
}

func (i *AuthInterceptor) getToken() (*authn.Token, error) {
	// If we have a token that is not expired (or nearing expiry) then return it
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

// Unary method to intercept requests and inject the token into headers.
func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
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
			return status.Errorf(codes.Unauthenticated, "failed to get token: %v", err)
		}
		// Create the metadata with the token
		md := metadata.Pairs(constants.NodeAuthorizationHeaderName, token.SignedString)
		// Attach metadata to the outgoing context
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Proceed with the request
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
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
		// Create the metadata with the token
		md := metadata.Pairs(constants.NodeAuthorizationHeaderName, token.SignedString)
		// Attach the metadata to the outgoing context
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Proceed with the stream
		return streamer(ctx, desc, cc, method, opts...)
	}
}

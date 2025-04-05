package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/mocks/authn"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryInterceptor(t *testing.T) {
	mockVerifier := authn.NewMockJWTVerifier(t)
	logger := zaptest.NewLogger(t)
	interceptor := NewAuthInterceptor(mockVerifier, logger)

	tests := []struct {
		name             string
		setupContext     func() context.Context
		setupVerifier    func()
		wantError        error
		wantVerifiedNode bool
	}{
		{
			name: "valid token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{
					constants.NODE_AUTHORIZATION_HEADER_NAME: "valid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().Verify("valid_token").Return(uint32(0), func() {}, nil)
			},
			wantError:        nil,
			wantVerifiedNode: true,
		},
		{
			name: "missing metadata",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupVerifier:    func() {},
			wantError:        nil,
			wantVerifiedNode: false,
		},
		{
			name: "missing token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupVerifier:    func() {},
			wantError:        nil,
			wantVerifiedNode: false,
		},
		{
			name: "invalid token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{
					constants.NODE_AUTHORIZATION_HEADER_NAME: "invalid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().
					Verify("invalid_token").
					Return(uint32(0), func() {}, errors.New("invalid signature"))
			},
			wantError: status.Error(
				codes.Unauthenticated,
				"invalid auth token: invalid signature",
			),
			wantVerifiedNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupVerifier()

			ctx := tt.setupContext()
			var handlerCtx context.Context
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				handlerCtx = ctx
				return "ok", nil
			}

			_, err := interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{}, handler)

			if tt.wantError != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				isVerified, hasContextValue := handlerCtx.Value(constants.VerifiedNodeRequestCtxKey{}).(bool)
				if tt.wantVerifiedNode {
					require.True(t, isVerified)
				} else {
					require.False(t, hasContextValue)
				}
			}
		})
	}
}

func TestStreamInterceptor(t *testing.T) {
	mockVerifier := authn.NewMockJWTVerifier(t)
	logger := zaptest.NewLogger(t)
	interceptor := NewAuthInterceptor(mockVerifier, logger)

	tests := []struct {
		name             string
		setupContext     func() context.Context
		setupVerifier    func()
		wantError        error
		wantVerifiedNode bool
	}{
		{
			name: "valid token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{
					constants.NODE_AUTHORIZATION_HEADER_NAME: "valid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().Verify("valid_token").Return(uint32(0), func() {}, nil)
			},
			wantError:        nil,
			wantVerifiedNode: true,
		},
		{
			name: "missing metadata",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupVerifier:    func() {},
			wantError:        nil,
			wantVerifiedNode: false,
		},
		{
			name: "invalid token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{
					constants.NODE_AUTHORIZATION_HEADER_NAME: "invalid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().
					Verify("invalid_token").
					Return(uint32(0), func() {}, errors.New("invalid signature"))
			},
			wantError: status.Error(
				codes.Unauthenticated,
				"invalid auth token: invalid signature",
			),
			wantVerifiedNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupVerifier()

			ctx := tt.setupContext()
			var handlerStream grpc.ServerStream
			stream := &mockServerStreamWithCtx{ctx: ctx}
			handler := func(srv interface{}, stream grpc.ServerStream) error {
				handlerStream = stream
				return nil
			}

			err := interceptor.Stream()(nil, stream, &grpc.StreamServerInfo{}, handler)

			if tt.wantError != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				isVerified, hasContextValue := handlerStream.Context().Value(constants.VerifiedNodeRequestCtxKey{}).(bool)
				if tt.wantVerifiedNode {
					require.True(t, isVerified)
				} else {
					require.False(t, hasContextValue)
				}
			}
		})
	}
}

type mockServerStreamWithCtx struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *mockServerStreamWithCtx) Context() context.Context {
	return s.ctx
}

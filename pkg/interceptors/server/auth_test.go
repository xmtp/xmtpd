package server

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/mocks/authn"
	"go.uber.org/zap/zaptest"
)

type mockConnectRequestAuthInterceptor struct {
	connect.AnyRequest
	header http.Header
	peer   connect.Peer
}

func (m *mockConnectRequestAuthInterceptor) Header() http.Header {
	return m.header
}

func (m *mockConnectRequestAuthInterceptor) Peer() connect.Peer {
	return m.peer
}

type mockStreamingConnAuthInterceptor struct {
	connect.StreamingHandlerConn
	header http.Header
	peer   connect.Peer
}

func (m *mockStreamingConnAuthInterceptor) RequestHeader() http.Header {
	return m.header
}

func (m *mockStreamingConnAuthInterceptor) Peer() connect.Peer {
	return m.peer
}

func TestUnaryInterceptor(t *testing.T) {
	mockVerifier := authn.NewMockJWTVerifier(t)
	logger := zaptest.NewLogger(t)
	interceptor := NewServerAuthInterceptor(logger, mockVerifier)

	tests := []struct {
		name             string
		setupRequest     func() connect.AnyRequest
		setupVerifier    func()
		wantError        bool
		wantVerifiedNode bool
	}{
		{
			name: "valid token",
			setupRequest: func() connect.AnyRequest {
				header := http.Header{}
				header.Set(constants.NodeAuthorizationHeaderName, "valid_token")
				return &mockConnectRequestAuthInterceptor{
					header: header,
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().Verify("valid_token").Return(uint32(0), func() {}, nil)
			},
			wantError:        false,
			wantVerifiedNode: true,
		},
		{
			name: "missing token",
			setupRequest: func() connect.AnyRequest {
				return &mockConnectRequestAuthInterceptor{
					header: http.Header{},
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier:    func() {},
			wantError:        false,
			wantVerifiedNode: false,
		},
		{
			name: "invalid token",
			setupRequest: func() connect.AnyRequest {
				header := http.Header{}
				header.Set(constants.NodeAuthorizationHeaderName, "invalid_token")
				return &mockConnectRequestAuthInterceptor{
					header: header,
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().
					Verify("invalid_token").
					Return(uint32(0), func() {}, errors.New("invalid signature"))
			},
			wantError:        true,
			wantVerifiedNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupVerifier()

			req := tt.setupRequest()
			var handlerCtx context.Context
			next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
				handlerCtx = ctx
				return nil, nil
			}

			wrappedUnary := interceptor.WrapUnary(next)
			_, err := wrappedUnary(context.Background(), req)

			if tt.wantError {
				require.Error(t, err)
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
	interceptor := NewServerAuthInterceptor(logger, mockVerifier)

	tests := []struct {
		name             string
		setupConn        func() connect.StreamingHandlerConn
		setupVerifier    func()
		wantError        bool
		wantVerifiedNode bool
	}{
		{
			name: "valid token",
			setupConn: func() connect.StreamingHandlerConn {
				header := http.Header{}
				header.Set(constants.NodeAuthorizationHeaderName, "valid_token")
				return &mockStreamingConnAuthInterceptor{
					header: header,
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().Verify("valid_token").Return(uint32(0), func() {}, nil)
			},
			wantError:        false,
			wantVerifiedNode: true,
		},
		{
			name: "missing token",
			setupConn: func() connect.StreamingHandlerConn {
				return &mockStreamingConnAuthInterceptor{
					header: http.Header{},
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier:    func() {},
			wantError:        false,
			wantVerifiedNode: false,
		},
		{
			name: "invalid token",
			setupConn: func() connect.StreamingHandlerConn {
				header := http.Header{}
				header.Set(constants.NodeAuthorizationHeaderName, "invalid_token")
				return &mockStreamingConnAuthInterceptor{
					header: header,
					peer:   connect.Peer{Addr: "127.0.0.1:1234"},
				}
			},
			setupVerifier: func() {
				mockVerifier.EXPECT().
					Verify("invalid_token").
					Return(uint32(0), func() {}, errors.New("invalid signature"))
			},
			wantError:        true,
			wantVerifiedNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupVerifier()

			conn := tt.setupConn()
			var handlerCtx context.Context
			next := func(ctx context.Context, c connect.StreamingHandlerConn) error {
				handlerCtx = ctx
				return nil
			}

			wrappedStream := interceptor.WrapStreamingHandler(next)
			err := wrappedStream(context.Background(), conn)

			if tt.wantError {
				require.Error(t, err)
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

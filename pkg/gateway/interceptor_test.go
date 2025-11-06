package gateway

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"go.uber.org/zap"
)

type mockConnectRequest struct {
	connect.AnyRequest
	header    http.Header
	peer      connect.Peer
	procedure string
	msg       any
}

func (m *mockConnectRequest) Header() http.Header {
	return m.header
}

func (m *mockConnectRequest) Peer() connect.Peer {
	return m.peer
}

func (m *mockConnectRequest) Spec() connect.Spec {
	return connect.Spec{Procedure: m.procedure}
}

func (m *mockConnectRequest) Any() any {
	return m.msg
}

func TestGatewayInterceptor(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name              string
		setupIdentityFn   func() IdentityFn
		setupAuthorizers  func() []AuthorizePublishFn
		setupRequest      func() *mockConnectRequest
		wantError         bool
		wantAuthzCalled   bool
		wantHandlerCalled bool
		checkContext      func(t *testing.T, ctx context.Context)
	}{
		{
			name: "Identity Injection",
			setupIdentityFn: func() IdentityFn {
				return func(headers http.Header, peer string) (Identity, error) {
					return Identity{
						Kind:     identityKindIP,
						Identity: "192.168.1.1",
					}, nil
				}
			},
			setupAuthorizers: func() []AuthorizePublishFn {
				return nil
			},
			setupRequest: func() *mockConnectRequest {
				return &mockConnectRequest{
					header:    http.Header{},
					peer:      connect.Peer{Addr: "127.0.0.1:1234"},
					procedure: "/xmtp.xmtpv4.payer_api.PayerApi/GetNodes",
					msg:       &payer_api.GetNodesRequest{},
				}
			},
			wantError:         false,
			wantAuthzCalled:   false,
			wantHandlerCalled: true,
			checkContext: func(t *testing.T, ctx context.Context) {
				identity, ok := IdentityFromContext(ctx)
				assert.True(t, ok)
				assert.Equal(t, Identity{
					Kind:     identityKindIP,
					Identity: "192.168.1.1",
				}, identity)
			},
		},
		{
			name: "Authorization Success",
			setupIdentityFn: func() IdentityFn {
				return func(headers http.Header, peer string) (Identity, error) {
					return Identity{
						Kind:     identityKindIP,
						Identity: "127.0.0.1",
					}, nil
				}
			},
			setupAuthorizers: func() []AuthorizePublishFn {
				return []AuthorizePublishFn{
					func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
						assert.Equal(t, "127.0.0.1", id.Identity)
						assert.Equal(t, 2, req.TotalEnvelopes)
						return true, nil
					},
				}
			},
			setupRequest: func() *mockConnectRequest {
				return &mockConnectRequest{
					header:    http.Header{},
					peer:      connect.Peer{Addr: "127.0.0.1:1234"},
					procedure: "/xmtp.xmtpv4.payer_api.PayerApi/PublishClientEnvelopes",
					msg: &payer_api.PublishClientEnvelopesRequest{
						Envelopes: []*envelopes.ClientEnvelope{{}, {}},
					},
				}
			},
			wantError:         false,
			wantAuthzCalled:   true,
			wantHandlerCalled: true,
			checkContext:      func(t *testing.T, ctx context.Context) {},
		},
		{
			name: "Authorization Denied",
			setupIdentityFn: func() IdentityFn {
				return func(headers http.Header, peer string) (Identity, error) {
					return Identity{
						Kind:     identityKindIP,
						Identity: "192.168.1.1",
					}, nil
				}
			},
			setupAuthorizers: func() []AuthorizePublishFn {
				return []AuthorizePublishFn{
					func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
						return false, nil
					},
				}
			},
			setupRequest: func() *mockConnectRequest {
				return &mockConnectRequest{
					header:    http.Header{},
					peer:      connect.Peer{Addr: "192.168.1.1:1234"},
					procedure: "/xmtp.xmtpv4.payer_api.PayerApi/PublishClientEnvelopes",
					msg: &payer_api.PublishClientEnvelopesRequest{
						Envelopes: []*envelopes.ClientEnvelope{{}},
					},
				}
			},
			wantError:         true,
			wantAuthzCalled:   true,
			wantHandlerCalled: false,
			checkContext:      func(t *testing.T, ctx context.Context) {},
		},
		{
			name: "Authorization Not Called For Other Methods",
			setupIdentityFn: func() IdentityFn {
				return func(headers http.Header, peer string) (Identity, error) {
					return Identity{
						Kind:     identityKindIP,
						Identity: "192.168.1.1",
					}, nil
				}
			},
			setupAuthorizers: func() []AuthorizePublishFn {
				return []AuthorizePublishFn{
					func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
						t.Fatal("authorizer should not be called")
						return true, nil
					},
				}
			},
			setupRequest: func() *mockConnectRequest {
				return &mockConnectRequest{
					header:    http.Header{},
					peer:      connect.Peer{Addr: "192.168.1.1:1234"},
					procedure: "/xmtp.xmtpv4.payer_api.PayerApi/GetNodes",
					msg:       &payer_api.GetNodesRequest{},
				}
			},
			wantError:         false,
			wantAuthzCalled:   false,
			wantHandlerCalled: true,
			checkContext:      func(t *testing.T, ctx context.Context) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewGatewayInterceptor(
				logger,
				tt.setupIdentityFn(),
				tt.setupAuthorizers(),
			)

			var handlerCtx context.Context
			handlerCalled := false
			next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
				handlerCtx = ctx
				handlerCalled = true
				return nil, nil
			}

			wrappedUnary := interceptor.WrapUnary(next)
			req := tt.setupRequest()

			_, err := wrappedUnary(context.Background(), req)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantHandlerCalled, handlerCalled)

			if !tt.wantError && tt.checkContext != nil {
				tt.checkContext(t, handlerCtx)
			}
		})
	}
}

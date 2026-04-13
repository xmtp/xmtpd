package server

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"go.uber.org/zap/zaptest"
)

func newRequireNodeAuthInterceptor(t *testing.T) *RequireNodeAuthInterceptor {
	t.Helper()
	return NewRequireNodeAuthInterceptor(zaptest.NewLogger(t))
}

func ctxWithVerifiedNode(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.VerifiedNodeRequestCtxKey{}, true)
}

func TestRequireNodeAuth_Unary_Authenticated(t *testing.T) {
	interceptor := newRequireNodeAuthInterceptor(t)

	called := false
	wrapped := interceptor.WrapUnary(
		func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			called = true
			return nil, nil
		},
	)

	ctx := ctxWithVerifiedNode(t.Context())
	req := &mockConnectRequest{
		header: http.Header{},
		spec: connect.Spec{
			Procedure: "/xmtp.xmtpv4.message_api.ReplicationApi/SubscribeOriginators",
		},
	}

	_, err := wrapped(ctx, req)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestRequireNodeAuth_Unary_Unauthenticated(t *testing.T) {
	interceptor := newRequireNodeAuthInterceptor(t)

	called := false
	wrapped := interceptor.WrapUnary(
		func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			called = true
			return nil, nil
		},
	)

	req := &mockConnectRequest{
		header: http.Header{},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes"},
	}

	_, err := wrapped(t.Context(), req)
	require.Error(t, err)
	assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
	assert.False(t, called)
}

func TestRequireNodeAuth_Stream_Authenticated(t *testing.T) {
	interceptor := newRequireNodeAuthInterceptor(t)

	called := false
	wrapped := interceptor.WrapStreamingHandler(
		func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			called = true
			return nil
		},
	)

	ctx := ctxWithVerifiedNode(t.Context())
	conn := &mockStreamingConn{
		header: http.Header{},
		spec: connect.Spec{
			Procedure: "/xmtp.xmtpv4.message_api.ReplicationApi/SubscribeEnvelopes",
		},
	}

	err := wrapped(ctx, conn)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestRequireNodeAuth_Stream_Unauthenticated(t *testing.T) {
	interceptor := newRequireNodeAuthInterceptor(t)

	called := false
	wrapped := interceptor.WrapStreamingHandler(
		func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			called = true
			return nil
		},
	)

	conn := &mockStreamingConn{
		header: http.Header{},
		spec: connect.Spec{
			Procedure: "/xmtp.xmtpv4.message_api.ReplicationApi/SubscribeOriginators",
		},
	}

	err := wrapped(t.Context(), conn)
	require.Error(t, err)
	assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
	assert.False(t, called)
}

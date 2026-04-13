package server

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"go.uber.org/zap/zaptest"
)

// fakeStreamLimiter is a test double for ratelimiter.StreamLimiter.
type fakeStreamLimiter struct {
	mu           sync.Mutex
	acquireAllow bool
	acquireErr   error
	releaseErr   error
	refreshErr   error
	acquireCount int
	releaseCount int
}

func (f *fakeStreamLimiter) Acquire(_ context.Context, _ string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.acquireCount++
	return f.acquireAllow, f.acquireErr
}

func (f *fakeStreamLimiter) Release(_ context.Context, _ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.releaseCount++
	return f.releaseErr
}

func (f *fakeStreamLimiter) RefreshTTL(_ context.Context, _ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.refreshErr
}

func TestStreamLimit_Tier0Bypasses(t *testing.T) {
	limiter := &fakeStreamLimiter{acquireAllow: true}
	interceptor := NewStreamLimitInterceptor(
		zaptest.NewLogger(t), limiter, nil, 5*time.Minute,
	)

	ctx := context.WithValue(
		context.Background(),
		constants.VerifiedNodeRequestCtxKey{},
		true,
	)
	conn := &mockStreamingConn{
		spec: connect.Spec{
			Procedure: message_apiconnect.NotificationApiSubscribeAllEnvelopesProcedure,
		},
		peer:   connect.Peer{Addr: "10.0.0.1:1234"},
		header: http.Header{},
	}

	called := false
	handler := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		called = true
		return nil
	}

	err := interceptor.WrapStreamingHandler(handler)(ctx, conn)
	require.NoError(t, err)
	assert.True(t, called, "handler should be called")
	assert.Equal(t, 0, limiter.acquireCount, "should not call acquire for tier 0")
}

func TestStreamLimit_Tier2Allowed(t *testing.T) {
	limiter := &fakeStreamLimiter{acquireAllow: true}
	interceptor := NewStreamLimitInterceptor(
		zaptest.NewLogger(t), limiter, nil, 5*time.Minute,
	)

	conn := &mockStreamingConn{
		spec: connect.Spec{
			Procedure: message_apiconnect.NotificationApiSubscribeAllEnvelopesProcedure,
		},
		peer:   connect.Peer{Addr: "10.0.0.1:1234"},
		header: http.Header{},
	}

	called := false
	handler := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		called = true
		return nil
	}

	err := interceptor.WrapStreamingHandler(handler)(context.Background(), conn)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, 1, limiter.acquireCount)
	assert.Equal(t, 1, limiter.releaseCount, "release should fire after handler returns")
}

func TestStreamLimit_Tier2Denied(t *testing.T) {
	limiter := &fakeStreamLimiter{acquireAllow: false}
	interceptor := NewStreamLimitInterceptor(
		zaptest.NewLogger(t), limiter, nil, 5*time.Minute,
	)

	conn := &mockStreamingConn{
		spec: connect.Spec{
			Procedure: message_apiconnect.NotificationApiSubscribeAllEnvelopesProcedure,
		},
		peer:   connect.Peer{Addr: "10.0.0.1:1234"},
		header: http.Header{},
	}

	called := false
	handler := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		called = true
		return nil
	}

	err := interceptor.WrapStreamingHandler(handler)(context.Background(), conn)
	require.Error(t, err)
	assert.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err))
	assert.False(t, called, "handler should not be called when denied")
	assert.Equal(t, 0, limiter.releaseCount, "release should not fire when denied")
}

func TestStreamLimit_FailOpenOnRedisError(t *testing.T) {
	limiter := &fakeStreamLimiter{
		acquireAllow: false,
		acquireErr:   errors.New("redis connection refused"),
	}
	interceptor := NewStreamLimitInterceptor(
		zaptest.NewLogger(t), limiter, nil, 5*time.Minute,
	)

	conn := &mockStreamingConn{
		spec: connect.Spec{
			Procedure: message_apiconnect.NotificationApiSubscribeAllEnvelopesProcedure,
		},
		peer:   connect.Peer{Addr: "10.0.0.1:1234"},
		header: http.Header{},
	}

	called := false
	handler := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		called = true
		return nil
	}

	err := interceptor.WrapStreamingHandler(handler)(context.Background(), conn)
	require.NoError(t, err)
	assert.True(t, called, "should fail open on Redis error")
}

func TestStreamLimit_NonNotificationApiPassesThrough(t *testing.T) {
	limiter := &fakeStreamLimiter{acquireAllow: true}
	interceptor := NewStreamLimitInterceptor(
		zaptest.NewLogger(t), limiter, nil, 5*time.Minute,
	)

	conn := &mockStreamingConn{
		spec: connect.Spec{
			Procedure: message_apiconnect.QueryApiSubscribeTopicsProcedure,
		},
		peer:   connect.Peer{Addr: "10.0.0.1:1234"},
		header: http.Header{},
	}

	called := false
	handler := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		called = true
		return nil
	}

	err := interceptor.WrapStreamingHandler(handler)(context.Background(), conn)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, 0, limiter.acquireCount, "should not touch limiter for non-NotificationApi")
}

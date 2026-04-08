package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap/zaptest"
)

// These tests exercise the Connect interceptor against a real Redis instance
// (via pkg/testutils/redis.NewRedisForTest, which expects redis on
// localhost:6379 DB 15 — see dev/up). They cover the same paths as the
// unit tests but with the real RedisLimiter and BreakerLimiter wired in.

func newQueryRequest(peerAddr string) *mockConnectRequest {
	return &mockConnectRequest{
		header: http.Header{},
		peer:   connect.Peer{Addr: peerAddr},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds"},
		body:   &message_api.GetInboxIdsRequest{},
	}
}

func newSubscribeConn(peerAddr string) *mockStreamingConn {
	return &mockStreamingConn{
		header: http.Header{},
		peer:   connect.Peer{Addr: peerAddr},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics"},
	}
}

// Test 7.1: A Tier 2 client is allowed up to the per-minute capacity, then
// the very next call is rejected with ResourceExhausted.
func TestRateLimitIntegration_Tier2QueryDenyAfterBudgetExhausted(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	queryLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":q", []ratelimiter.Limit{
		{Capacity: 5, RefillEvery: time.Minute},
		{Capacity: 100, RefillEvery: time.Hour},
	})
	require.NoError(t, err)
	opensLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":o", []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	rl := NewRateLimitInterceptor(
		zaptest.NewLogger(t),
		queryLimiter,
		opensLimiter,
		nil,
		RateLimitInterceptorConfig{},
	)

	calls := 0
	next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		calls++
		return nil, nil
	}
	wrapped := rl.WrapUnary(next)

	for i := range 5 {
		_, err := wrapped(context.Background(), newQueryRequest("203.0.113.1:5001"))
		require.NoError(t, err, "call %d should succeed", i)
	}

	_, err = wrapped(context.Background(), newQueryRequest("203.0.113.1:5001"))
	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())
	assert.Equal(t, 5, calls, "inner handler should not be called on the 6th request")
}

// Test 7.2: ForceDebit can drive the bucket negative, after which a normal
// Allow call is rejected.
func TestRateLimitIntegration_RetrospectiveDrainGoesNegative(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	ctx := context.Background()

	res, err := limiter.Allow(ctx, "subj", 9)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	res, err = limiter.ForceDebit(ctx, "subj", 100)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.Less(t, res.Balances[0].Remaining, -80.0)

	res, err = limiter.Allow(ctx, "subj", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed, "next normal call should be rejected after retrospective drain")
}

// Test 7.3: A Tier 0 client (verified node) bypasses even an absurdly tight
// limit and is never charged.
func TestRateLimitIntegration_Tier0Bypass(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 1, RefillEvery: time.Hour}, // very tight
	})
	require.NoError(t, err)

	rl := NewRateLimitInterceptor(
		zaptest.NewLogger(t),
		limiter,
		limiter,
		nil,
		RateLimitInterceptorConfig{},
	)

	calls := 0
	next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		calls++
		return nil, nil
	}
	wrapped := rl.WrapUnary(next)

	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)
	for i := range 100 {
		_, err := wrapped(ctx, newQueryRequest("10.0.0.5:5001"))
		require.NoError(t, err, "call %d should succeed under Tier 0 bypass", i)
	}
	assert.Equal(t, 100, calls)
}

// Test 7.4: The subscribe-opens-per-minute sub-limit denies further opens
// once the bucket is exhausted.
func TestRateLimitIntegration_SubscribeOpensSubLimit(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	queryLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":q", []ratelimiter.Limit{
		{Capacity: 1000, RefillEvery: time.Minute},
	})
	require.NoError(t, err)
	opensLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":o", []ratelimiter.Limit{
		{Capacity: 3, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	rl := NewRateLimitInterceptor(
		zaptest.NewLogger(t),
		queryLimiter,
		opensLimiter,
		nil,
		RateLimitInterceptorConfig{},
	)

	calls := 0
	next := func(ctx context.Context, c connect.StreamingHandlerConn) error {
		calls++
		return nil
	}
	wrapped := rl.WrapStreamingHandler(next)

	for i := range 3 {
		require.NoError(t,
			wrapped(context.Background(), newSubscribeConn("203.0.113.1:5001")),
			"open %d should succeed", i)
	}

	err = wrapped(context.Background(), newSubscribeConn("203.0.113.1:5001"))
	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())
	assert.Equal(t, 3, calls)
}

// Test 7.5: When Redis goes away, the BreakerLimiter trips after the failure
// threshold and subsequent calls fail open without touching Redis.
func TestRateLimitIntegration_RedisDownTripsBreakerAndFailsOpen(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	inner, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 100, RefillEvery: time.Minute},
	})
	require.NoError(t, err)
	cb := ratelimiter.NewCircuitBreaker(2, 200*time.Millisecond)
	bl := ratelimiter.NewBreakerLimiter(inner, cb)

	// Sanity check: a normal call works while Redis is up.
	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// Kill the underlying connection: every subsequent inner call will error.
	require.NoError(t, client.Close())

	// First two failing calls trip the breaker; both fail open.
	for i := range 2 {
		res, err := bl.Allow(context.Background(), "subj", 1)
		require.NoError(t, err)
		require.True(t, res.Allowed, "call %d should fail open", i)
	}
	require.Equal(t, ratelimiter.BreakerOpen, cb.State())

	// Subsequent calls bypass Redis entirely and remain fail-open.
	for i := range 5 {
		res, err := bl.Allow(context.Background(), "subj", 1)
		require.NoError(t, err)
		require.True(t, res.Allowed, "post-trip call %d should fail open", i)
	}
}

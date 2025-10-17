package authorizers

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/gateway"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap"
)

func TestRateLimitAuthorizerBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *RateLimitBuilder
		wantErr string
	}{
		{
			name: "missing redis client",
			setup: func() *RateLimitBuilder {
				return NewRateLimitBuilder().
					WithLimits(RateLimit{Capacity: 10, RefillEvery: time.Minute})
			},
			wantErr: "redis client is not set",
		},
		{
			name: "empty limits",
			setup: func() *RateLimitBuilder {
				client, keyPrefix := redistestutils.NewRedisForTest(t)
				return NewRateLimitBuilder().
					WithKeyPrefix(keyPrefix).
					WithRedis(client)
			},
			wantErr: "no rate limits configured",
		},
		{
			name: "both missing",
			setup: func() *RateLimitBuilder {
				return NewRateLimitBuilder()
			},
			wantErr: "redis client is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.setup()
			_, err := builder.Build()
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestRateLimitAuthorizerBuilder_Success(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	logger := zap.NewNop()

	authFn, err := NewRateLimitBuilder().
		WithKeyPrefix(keyPrefix).
		WithRedis(client).
		WithLogger(logger).
		WithLimits(RateLimit{Capacity: 10, RefillEvery: time.Minute}).
		Build()

	require.NoError(t, err)
	require.NotNil(t, authFn)
}

func TestRateLimitAuthorizer_RespectsLimits(t *testing.T) {
	tests := []struct {
		name       string
		limits     []RateLimit
		requests   []int
		wantAllow  []bool
		wantErrors []bool
	}{
		{
			name:       "within single limit",
			limits:     []RateLimit{{Capacity: 10, RefillEvery: time.Minute}},
			requests:   []int{5, 5},
			wantAllow:  []bool{true, true},
			wantErrors: []bool{false, false},
		},
		{
			name:       "exceeds single limit",
			limits:     []RateLimit{{Capacity: 10, RefillEvery: time.Minute}},
			requests:   []int{10, 1},
			wantAllow:  []bool{true, false},
			wantErrors: []bool{false, true},
		},
		{
			name: "within multiple limits",
			limits: []RateLimit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			requests:   []int{5, 5},
			wantAllow:  []bool{true, true},
			wantErrors: []bool{false, false},
		},
		{
			name: "exceeds first limit",
			limits: []RateLimit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			requests:   []int{10, 1},
			wantAllow:  []bool{true, false},
			wantErrors: []bool{false, true},
		},
		{
			name: "exceeds second limit",
			limits: []RateLimit{
				{Capacity: 100, RefillEvery: time.Hour},
				{Capacity: 10, RefillEvery: time.Minute},
			},
			requests:   []int{10, 1},
			wantAllow:  []bool{true, false},
			wantErrors: []bool{false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)

			authFn, err := NewRateLimitBuilder().
				WithRedis(client).
				WithLimits(tt.limits...).
				WithKeyPrefix(keyPrefix).
				Build()
			require.NoError(t, err)

			identity := gateway.Identity{Identity: "test-user"}
			ctx := t.Context()

			for i, envelopes := range tt.requests {
				req := gateway.PublishRequestSummary{TotalEnvelopes: envelopes}
				allowed, err := authFn(ctx, identity, req)

				if tt.wantErrors[i] {
					require.Error(t, err, "request %d (envelopes=%d) should error", i, envelopes)
				} else {
					require.NoError(t, err, "request %d (envelopes=%d) should not error. Error: %v", i, envelopes, err)
				}

				require.Equal(t, tt.wantAllow[i], allowed, "request %d allowed mismatch", i)
			}
		})
	}
}

func TestRateLimitAuthorizer_ErrorType(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	authFn, err := NewRateLimitBuilder().
		WithKeyPrefix(keyPrefix).
		WithRedis(client).
		WithLimits(RateLimit{Capacity: 10, RefillEvery: time.Minute}).
		Build()
	require.NoError(t, err)

	identity := gateway.Identity{Identity: "test-user"}
	ctx := t.Context()

	// Consume all tokens
	req := gateway.PublishRequestSummary{TotalEnvelopes: 10}
	allowed, err := authFn(ctx, identity, req)
	require.NoError(t, err)
	require.True(t, allowed)

	// Next request should fail with rate limit error
	req = gateway.PublishRequestSummary{TotalEnvelopes: 1}
	allowed, err = authFn(ctx, identity, req)
	require.Error(t, err)
	require.False(t, allowed)

	// Verify error is GatewayServiceError
	var gwErr gateway.GatewayServiceError
	require.True(t, errors.As(err, &gwErr), "error should be GatewayServiceError")

	// Verify RetryAfter is set
	retryAfter := gwErr.RetryAfter()
	require.NotNil(t, retryAfter, "RetryAfter should not be nil")
	require.Greater(t, *retryAfter, time.Duration(0), "RetryAfter should be positive")
}

func TestRateLimitAuthorizer_RetryAfterAccuracy(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	authFn, err := NewRateLimitBuilder().
		WithRedis(client).
		WithKeyPrefix(keyPrefix).
		WithLimits(RateLimit{Capacity: 10, RefillEvery: time.Second}).
		Build()
	require.NoError(t, err)

	identity := gateway.Identity{Identity: "test-user"}
	ctx := t.Context()

	// Consume all tokens
	req := gateway.PublishRequestSummary{TotalEnvelopes: 10}
	allowed, err := authFn(ctx, identity, req)
	require.NoError(t, err)
	require.True(t, allowed)

	// Request more than capacity
	req = gateway.PublishRequestSummary{TotalEnvelopes: 5}
	allowed, err = authFn(ctx, identity, req)
	require.Error(t, err)
	require.False(t, allowed)

	var gwErr gateway.GatewayServiceError
	require.True(t, errors.As(err, &gwErr))

	retryAfter := gwErr.RetryAfter()
	require.NotNil(t, retryAfter)

	// RetryAfter should be approximately 500ms (need 5 tokens, refill rate is 1s/10 = 100ms per token)
	require.InDelta(t, 500*time.Millisecond, *retryAfter, float64(100*time.Millisecond))
}

func TestRateLimitAuthorizer_SubjectIsolation(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	authFn, err := NewRateLimitBuilder().
		WithRedis(client).
		WithKeyPrefix(keyPrefix).
		WithLimits(RateLimit{Capacity: 10, RefillEvery: time.Minute}).
		Build()
	require.NoError(t, err)

	ctx := t.Context()

	// User 1 consumes all tokens
	user1 := gateway.Identity{Identity: "user1"}
	req := gateway.PublishRequestSummary{TotalEnvelopes: 10}
	allowed, err := authFn(ctx, user1, req)
	require.NoError(t, err)
	require.True(t, allowed)

	// User 1 should be blocked
	req = gateway.PublishRequestSummary{TotalEnvelopes: 1}
	allowed, err = authFn(ctx, user1, req)
	require.Error(t, err)
	require.False(t, allowed)

	// User 2 should still have full capacity
	user2 := gateway.Identity{Identity: "user2"}
	req = gateway.PublishRequestSummary{TotalEnvelopes: 10}
	allowed, err = authFn(ctx, user2, req)
	require.NoError(t, err)
	require.True(t, allowed)
}

package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	testredis "github.com/xmtp/xmtpd/pkg/testutils/redis"
)

func TestRedisStreamLimiter_AcquireUnderLimit(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 2, 15*time.Minute)

	allowed, err := limiter.Acquire(context.Background(), "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestRedisStreamLimiter_AcquireAtLimit(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 2, 15*time.Minute)
	ctx := context.Background()

	allowed1, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed1)

	allowed2, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed2)

	allowed3, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.False(t, allowed3, "third stream should be denied")
}

func TestRedisStreamLimiter_ReleaseFreesSlot(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 2, 15*time.Minute)
	ctx := context.Background()

	_, _ = limiter.Acquire(ctx, "10.0.0.1")
	_, _ = limiter.Acquire(ctx, "10.0.0.1")

	err := limiter.Release(ctx, "10.0.0.1")
	require.NoError(t, err)

	allowed, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed, "slot should be freed after release")
}

func TestRedisStreamLimiter_DifferentIPsIndependent(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 1, 15*time.Minute)
	ctx := context.Background()

	allowed1, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed1)

	allowed2, err := limiter.Acquire(ctx, "10.0.0.2")
	require.NoError(t, err)
	assert.True(t, allowed2, "different IP should have independent limit")
}

func TestRedisStreamLimiter_RefreshTTL(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	ttl := 2 * time.Second
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 2, ttl)
	ctx := context.Background()

	_, _ = limiter.Acquire(ctx, "10.0.0.1")

	err := limiter.RefreshTTL(ctx, "10.0.0.1")
	require.NoError(t, err)

	remaining, err := client.TTL(ctx, prefix+"streams:10.0.0.1").Result()
	require.NoError(t, err)
	assert.Greater(t, remaining, time.Duration(0), "TTL should be positive after refresh")
	assert.LessOrEqual(t, remaining, ttl, "TTL should not exceed configured value")
}

func TestRedisStreamLimiter_TTLExpiryFreesSlot(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	shortTTL := 1 * time.Second
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 1, shortTTL)
	ctx := context.Background()

	allowed, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Wait for TTL to expire
	time.Sleep(shortTTL + 500*time.Millisecond)

	// Slot should be freed by TTL expiry
	allowed, err = limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed, "slot should be freed after TTL expiry")
}

func TestRedisStreamLimiter_ReleaseNeverGoesNegative(t *testing.T) {
	client, prefix := testredis.NewRedisForTest(t)
	limiter := ratelimiter.NewRedisStreamLimiter(client, prefix+"streams:", 2, 15*time.Minute)
	ctx := context.Background()

	// Release without acquire should not make count negative
	err := limiter.Release(ctx, "10.0.0.1")
	require.NoError(t, err)

	// Should still be able to acquire max streams
	allowed1, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed1)

	allowed2, err := limiter.Acquire(ctx, "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, allowed2)
}

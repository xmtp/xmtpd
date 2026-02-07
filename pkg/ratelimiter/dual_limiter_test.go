package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
)

func setupDualLimiterTest(t *testing.T) *ratelimiter.DualRedisLimiter {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	// Gateway: higher limits (100 requests per minute)
	gatewayLimits := []ratelimiter.Limit{
		{Capacity: 100, RefillEvery: time.Minute},
	}
	gatewayLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+"gateway", gatewayLimits)
	require.NoError(t, err)

	// User: lower limits (10 requests per minute)
	userLimits := []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	}
	userLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+"user", userLimits)
	require.NoError(t, err)

	dualLimiter, err := ratelimiter.NewDualRedisLimiter(gatewayLimiter, userLimiter)
	require.NoError(t, err)

	return dualLimiter
}

func TestDualLimiter_BothAllowed(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	result, err := limiter.AllowDual(ctx, "gateway1", "user1", 1)
	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, "", result.FailedSubject)
	assert.NotNil(t, result.GatewayResult)
	assert.NotNil(t, result.UserResult)
	assert.True(t, result.GatewayResult.Allowed)
	assert.True(t, result.UserResult.Allowed)
}

func TestDualLimiter_UserExceeded(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	// Exhaust user limit (10 requests)
	for i := 0; i < 10; i++ {
		result, err := limiter.AllowDual(ctx, "gateway1", "user1", 1)
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	// Next request should fail on user limit
	result, err := limiter.AllowDual(ctx, "gateway1", "user1", 1)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "user", result.FailedSubject)
	assert.NotNil(t, result.UserResult)
	assert.False(t, result.UserResult.Allowed)
}

func TestDualLimiter_GatewayExceeded(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	// Exhaust gateway limit by using different users (100 requests)
	for i := 0; i < 100; i++ {
		user := "user" + string(rune('a'+i%26))
		result, err := limiter.AllowDual(ctx, "gateway1", user, 1)
		require.NoError(t, err)
		assert.True(t, result.Allowed, "request %d should be allowed", i)
	}

	// Next request should fail on gateway limit
	result, err := limiter.AllowDual(ctx, "gateway1", "userNew", 1)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "gateway", result.FailedSubject)
	assert.NotNil(t, result.GatewayResult)
	assert.False(t, result.GatewayResult.Allowed)
}

func TestDualLimiter_NoUserSubject(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	// Empty user subject - only gateway limit checked
	result, err := limiter.AllowDual(ctx, "gateway1", "", 1)
	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.NotNil(t, result.GatewayResult)
	assert.Nil(t, result.UserResult)
}

func TestDualLimiter_DifferentUsersIndependent(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	// Exhaust limit for user1
	for i := 0; i < 10; i++ {
		result, err := limiter.AllowDual(ctx, "gateway1", "user1", 1)
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	// user1 is now rate-limited
	result, err := limiter.AllowDual(ctx, "gateway1", "user1", 1)
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// user2 should still work
	result, err = limiter.AllowDual(ctx, "gateway1", "user2", 1)
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}

func TestDualLimiter_AllowGatewayOnly(t *testing.T) {
	limiter := setupDualLimiterTest(t)

	ctx := context.Background()

	result, err := limiter.AllowGatewayOnly(ctx, "gateway1", 1)
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}

func TestNewDualRedisLimiter_NilGatewayLimiter(t *testing.T) {
	_, err := ratelimiter.NewDualRedisLimiter(nil, nil)
	assert.Error(t, err)
}

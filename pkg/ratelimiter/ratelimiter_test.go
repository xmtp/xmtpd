package ratelimiter

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
)

func TestNewRateLimiter(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	tests := []struct {
		name      string
		client    redis.UniversalClient
		limits    []Limit
		wantError bool
	}{
		{
			name:   "valid single limit",
			client: client,
			limits: []Limit{
				{Limit: 10, Window: time.Minute},
			},
			wantError: false,
		},
		{
			name:   "valid multiple limits",
			client: client,
			limits: []Limit{
				{Limit: 10, Window: time.Minute},
				{Limit: 100, Window: time.Hour},
			},
			wantError: false,
		},
		{
			name:      "nil client",
			client:    nil,
			limits:    []Limit{{Limit: 10, Window: time.Minute}},
			wantError: true,
		},
		{
			name:      "no limits",
			client:    client,
			limits:    []Limit{},
			wantError: true,
		},
		{
			name:   "zero limit",
			client: client,
			limits: []Limit{
				{Limit: 0, Window: time.Minute},
			},
			wantError: true,
		},
		{
			name:   "negative limit",
			client: client,
			limits: []Limit{
				{Limit: -1, Window: time.Minute},
			},
			wantError: true,
		},
		{
			name:   "zero window",
			client: client,
			limits: []Limit{
				{Limit: 10, Window: 0},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl, err := NewRateLimiterWithOptions(tt.client, tt.limits, Options{
				KeyPrefix: keyPrefix,
			})
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, rl)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rl)
			}
		})
	}
}

func TestRateLimiterSpend(t *testing.T) {
	ctx := t.Context()

	t.Run("single limit basic", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 5, Window: time.Second},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// First 5 requests should succeed
		for i := range 5 {
			allowed, err := rl.Spend(ctx, "user1", 1)
			require.NoError(t, err)
			require.True(t, allowed, "request %d should be allowed", i+1)
		}

		// 6th request should fail
		allowed, err := rl.Spend(ctx, "user1", 1)
		require.NoError(t, err)
		require.False(t, allowed, "6th request should be denied")
	})

	t.Run("multiple limits", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 5, Window: 100 * time.Millisecond},
			{Limit: 10, Window: 200 * time.Millisecond},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// First 5 requests should succeed
		for range 5 {
			allowed, err := rl.Spend(ctx, "user2", 1)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// 6th request should fail (hits first limit)
		allowed, err := rl.Spend(ctx, "user2", 1)
		require.NoError(t, err)
		require.False(t, allowed)

		// Wait for first window to expire
		time.Sleep(110 * time.Millisecond)

		// Next 5 requests should succeed (total 10 in 200ms)
		for range 5 {
			allowed, err := rl.Spend(ctx, "user2", 1)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// 11th request should fail (hits second limit)
		allowed, err = rl.Spend(ctx, "user2", 1)
		require.NoError(t, err)
		require.False(t, allowed)
	})

	t.Run("cost parameter", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 10, Window: time.Second},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// Single request with cost 5
		allowed, err := rl.Spend(ctx, "user3", 5)
		require.NoError(t, err)
		require.True(t, allowed)

		// Another request with cost 5
		allowed, err = rl.Spend(ctx, "user3", 5)
		require.NoError(t, err)
		require.True(t, allowed)

		// Request with cost 1 should fail (would exceed limit)
		allowed, err = rl.Spend(ctx, "user3", 1)
		require.NoError(t, err)
		require.False(t, allowed)
	})

	t.Run("different identifiers", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 2, Window: time.Second},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// Max out user1
		for range 2 {
			allowed, err := rl.Spend(ctx, "user1", 1)
			require.NoError(t, err)
			require.True(t, allowed)
		}
		allowed, err := rl.Spend(ctx, "user1", 1)
		require.NoError(t, err)
		require.False(t, allowed)

		// user2 should still be allowed
		allowed, err = rl.Spend(ctx, "user2", 1)
		require.NoError(t, err)
		require.True(t, allowed)
	})

	t.Run("leaky bucket refill", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 3, Window: 300 * time.Millisecond},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// Use 3 requests
		for range 3 {
			allowed, err := rl.Spend(ctx, "user4", 1)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// 4th should fail immediately
		allowed, err := rl.Spend(ctx, "user4", 1)
		require.NoError(t, err)
		require.False(t, allowed)

		// Wait for tokens to refill (1 token per 100ms)
		time.Sleep(110 * time.Millisecond)

		// Should be allowed again (1 token refilled)
		allowed, err = rl.Spend(ctx, "user4", 1)
		require.NoError(t, err)
		require.True(t, allowed)

		// But next should fail again
		allowed, err = rl.Spend(ctx, "user4", 1)
		require.NoError(t, err)
		require.False(t, allowed)
	})

	t.Run("invalid cost", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 10, Window: time.Second},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// Zero cost
		allowed, err := rl.Spend(ctx, "user5", 0)
		require.Error(t, err)
		require.False(t, allowed)

		// Negative cost
		allowed, err = rl.Spend(ctx, "user5", -1)
		require.Error(t, err)
		require.False(t, allowed)
	})

	t.Run("atomicity", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 10, Window: time.Second},
			{Limit: 15, Window: 2 * time.Second},
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		// Use 9 requests
		for range 9 {
			allowed, err := rl.Spend(ctx, "user6", 1)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// Request with cost 2 should fail (would exceed first limit)
		// But it should not consume from any limit
		allowed, err := rl.Spend(ctx, "user6", 2)
		require.NoError(t, err)
		require.False(t, allowed)

		// A request with cost 1 should still succeed
		allowed, err = rl.Spend(ctx, "user6", 1)
		require.NoError(t, err)
		require.True(t, allowed)
	})

	t.Run("concurrent requests", func(t *testing.T) {
		client, keyPrefix := redistestutils.NewRedisForTest(t)
		rl, err := NewRateLimiterWithOptions(client, []Limit{
			{Limit: 100, Window: 10 * time.Second}, // Longer window to minimize refill during test
		}, Options{
			KeyPrefix: keyPrefix,
		})
		require.NoError(t, err)

		var allowed atomic.Int32
		var wg sync.WaitGroup

		// Launch 150 concurrent requests
		for range 150 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ok, err := rl.Spend(ctx, "concurrent", 1)
				if err == nil && ok {
					allowed.Add(1)
				}
			}()
		}

		wg.Wait()

		// Should be close to 100 (may be slightly more due to refill during execution)
		allowedCount := allowed.Load()
		require.GreaterOrEqual(t, allowedCount, int32(100))
		require.LessOrEqual(t, allowedCount, int32(102)) // Allow small margin for refill
	})
}

func TestKeyGeneration(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	rl := &redisRateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		limits: []Limit{
			{Limit: 10, Window: time.Second},
			{Limit: 100, Window: time.Minute},
		},
	}

	t.Run("generateKey produces consistent keys", func(t *testing.T) {
		// Same inputs should produce same key
		key1 := rl.generateKey("user123", Limit{Limit: 10, Window: time.Second})
		key2 := rl.generateKey("user123", Limit{Limit: 10, Window: time.Second})
		require.Equal(t, key1, key2)

		// Different identifiers should produce different keys
		key3 := rl.generateKey("user456", Limit{Limit: 10, Window: time.Second})
		require.NotEqual(t, key1, key3)

		// Different limits should produce different keys
		key4 := rl.generateKey("user123", Limit{Limit: 20, Window: time.Second})
		require.NotEqual(t, key1, key4)

		// Different windows should produce different keys
		key5 := rl.generateKey("user123", Limit{Limit: 10, Window: time.Minute})
		require.NotEqual(t, key1, key5)
	})

	t.Run("generateKey uses prefix correctly", func(t *testing.T) {
		key := rl.generateKey("user123", Limit{Limit: 10, Window: time.Second})
		require.True(t, strings.HasPrefix(key, keyPrefix))
	})

	t.Run("generateKeys produces correct number of keys", func(t *testing.T) {
		keys := rl.generateKeys("user123")
		require.Len(t, keys, len(rl.limits))

		// All keys should be unique
		keySet := make(map[string]bool)
		for _, key := range keys {
			require.False(t, keySet[key], "duplicate key found")
			keySet[key] = true
		}
	})

	t.Run("generateKeys matches individual key generation", func(t *testing.T) {
		keys := rl.generateKeys("user123")
		for i, limit := range rl.limits {
			expectedKey := rl.generateKey("user123", limit)
			require.Equal(t, expectedKey, keys[i])
		}
	})
}

func BenchmarkRateLimiterSpend(b *testing.B) {
	ctx := context.Background()
	t := &testing.T{}
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	defer client.Close()

	rl, err := NewRateLimiterWithOptions(client, []Limit{
		{Limit: 1000000, Window: time.Hour}, // High limit to avoid hitting it
	}, Options{
		KeyPrefix: keyPrefix,
	})
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = rl.Spend(ctx, "bench", 1)
		}
	})
}

package ratelimiter_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
)

func TestRedisLimiter_BasicLimits(t *testing.T) {
	tests := []struct {
		name      string
		limits    []ratelimiter.Limit
		requests  []uint64
		wantAllow []bool
	}{
		{
			name:      "single limit allows within capacity",
			limits:    []ratelimiter.Limit{{Capacity: 10, RefillEvery: time.Minute}},
			requests:  []uint64{5, 5},
			wantAllow: []bool{true, true},
		},
		{
			name:      "single limit blocks over capacity",
			limits:    []ratelimiter.Limit{{Capacity: 10, RefillEvery: time.Minute}},
			requests:  []uint64{10, 1},
			wantAllow: []bool{true, false},
		},
		{
			name: "multiple limits all pass",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			requests:  []uint64{5, 5},
			wantAllow: []bool{true, true},
		},
		{
			name: "first limit blocks",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			requests:  []uint64{10, 1},
			wantAllow: []bool{true, false},
		},
		{
			name: "second limit blocks",
			limits: []ratelimiter.Limit{
				{Capacity: 100, RefillEvery: time.Hour},
				{Capacity: 10, RefillEvery: time.Minute},
			},
			requests:  []uint64{10, 1},
			wantAllow: []bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.NoError(t, err)

			for i, cost := range tt.requests {
				res, err := limiter.Allow(context.Background(), "test-subject", cost)
				require.NoError(t, err)
				require.Equal(t, tt.wantAllow[i], res.Allowed, "request %d", i)
			}
		})
	}
}

func TestRedisLimiter_Atomicity(t *testing.T) {
	tests := []struct {
		name               string
		limits             []ratelimiter.Limit
		cost               uint64
		wantAllow          bool
		wantRemaining      []float64
		wantFailedLimitIdx *int
	}{
		{
			name: "all limits pass - all decremented",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 20, RefillEvery: time.Minute},
			},
			cost:               5,
			wantAllow:          true,
			wantRemaining:      []float64{5, 15},
			wantFailedLimitIdx: nil,
		},
		{
			name: "first limit fails - none decremented",
			limits: []ratelimiter.Limit{
				{Capacity: 3, RefillEvery: time.Minute},
				{Capacity: 20, RefillEvery: time.Minute},
			},
			cost:               5,
			wantAllow:          false,
			wantRemaining:      []float64{3, 20},
			wantFailedLimitIdx: func() *int { i := 0; return &i }(),
		},
		{
			name: "second limit fails - none decremented",
			limits: []ratelimiter.Limit{
				{Capacity: 20, RefillEvery: time.Minute},
				{Capacity: 3, RefillEvery: time.Minute},
			},
			cost:               5,
			wantAllow:          false,
			wantRemaining:      []float64{20, 3},
			wantFailedLimitIdx: func() *int { i := 1; return &i }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.NoError(t, err)

			res, err := limiter.Allow(context.Background(), "test-subject", tt.cost)
			require.NoError(t, err)
			require.Equal(t, tt.wantAllow, res.Allowed)

			if tt.wantFailedLimitIdx == nil {
				require.Nil(t, res.FailedLimit)
				require.Nil(t, res.RetryAfter)
			} else {
				require.NotNil(t, res.FailedLimit)
				require.Equal(t, tt.limits[*tt.wantFailedLimitIdx], *res.FailedLimit)
				require.NotNil(t, res.RetryAfter)
				require.Greater(t, *res.RetryAfter, time.Duration(0))
			}

			require.Len(t, res.Balances, len(tt.limits))
			for i, want := range tt.wantRemaining {
				require.Equal(t, tt.limits[i], res.Balances[i].Limit, "balance[%d].Limit", i)
				require.InDelta(
					t,
					want,
					res.Balances[i].Remaining,
					0.01,
					"balance[%d].Remaining",
					i,
				)
			}

			// Second request with cost=1 to verify state wasn't modified on failure
			// After a failed request, the state should remain unchanged
			res2, err := limiter.Allow(context.Background(), "test-subject", 1)
			if tt.wantAllow {
				// First request succeeded, so second request should reflect deducted tokens
				require.NoError(t, err)
				for i, want := range tt.wantRemaining {
					require.InDelta(
						t,
						want-1,
						res2.Balances[i].Remaining,
						0.01,
						"balance[%d].Remaining after second check",
						i,
					)
				}
			} else {
				// First request failed, state unchanged, so check with another attempt
				// Result depends on whether limits still have enough tokens
				require.NoError(t, err)
				for i, want := range tt.wantRemaining {
					// If the original cost was more than capacity, this will still fail
					// Otherwise, verify remaining is unchanged
					if want >= 1 {
						require.InDelta(
							t,
							want-1,
							res2.Balances[i].Remaining,
							0.01,
							"balance[%d].Remaining after second check",
							i,
						)
					}
				}
			}
		})
	}
}

func TestRedisLimiter_RemainingAccuracy(t *testing.T) {
	tests := []struct {
		name          string
		limits        []ratelimiter.Limit
		requests      []uint64
		wantRemaining [][]float64
	}{
		{
			name:     "single limit tracks remaining correctly",
			limits:   []ratelimiter.Limit{{Capacity: 10, RefillEvery: time.Minute}},
			requests: []uint64{3, 2, 4},
			wantRemaining: [][]float64{
				{7},
				{5},
				{1},
			},
		},
		{
			name: "multiple limits track independently",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			requests: []uint64{5, 3},
			wantRemaining: [][]float64{
				{5, 95},
				{2, 92},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.NoError(t, err)

			for i, cost := range tt.requests {
				res, err := limiter.Allow(context.Background(), "test-subject", cost)
				require.NoError(t, err)
				require.True(t, res.Allowed, "request %d should be allowed", i)
				require.Len(t, res.Balances, len(tt.limits))
				for j, want := range tt.wantRemaining[i] {
					require.Equal(
						t,
						tt.limits[j],
						res.Balances[j].Limit,
						"request %d, balance[%d].Limit",
						i,
						j,
					)
					require.InDelta(
						t,
						want,
						res.Balances[j].Remaining,
						0.01,
						"request %d, balance[%d].Remaining",
						i,
						j,
					)
				}
			}
		})
	}
}

func TestRedisLimiter_TTL(t *testing.T) {
	tests := []struct {
		name          string
		limits        []ratelimiter.Limit
		cost          uint64
		wantTSTTL     time.Duration
		wantLimitTTLs []time.Duration // Expected TTL for each limit key
		wantTTLDelta  time.Duration
	}{
		{
			name:      "single limit sets TTL to refill time",
			limits:    []ratelimiter.Limit{{Capacity: 10, RefillEvery: 5 * time.Second}},
			cost:      1,
			wantTSTTL: 5 * time.Second,
			wantLimitTTLs: []time.Duration{
				500 * time.Millisecond,
			}, // (10-1)/10 * 5s = 4.5s remaining, but after deduction
			wantTTLDelta: 200 * time.Millisecond,
		},
		{
			name: "multiple limits use max refill time for timestamp",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: 2 * time.Second},
				{Capacity: 20, RefillEvery: 10 * time.Second},
			},
			cost:      1,
			wantTSTTL: 10 * time.Second,
			wantLimitTTLs: []time.Duration{
				200 * time.Millisecond, // (10-1)/10 * 2s = 1.8s
				500 * time.Millisecond, // (20-1)/20 * 10s = 9.5s
			},
			wantTTLDelta: 300 * time.Millisecond,
		},
		{
			name: "full bucket expires at refill period",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: 3 * time.Second},
			},
			cost:      10, // Consume all tokens
			wantTSTTL: 3 * time.Second,
			wantLimitTTLs: []time.Duration{
				3 * time.Second,
			}, // Bucket empty, will be full in refill_ms
			wantTTLDelta: 200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.NoError(t, err)

			// Make a request to trigger key creation
			_, err = limiter.Allow(context.Background(), "test-subject", tt.cost)
			require.NoError(t, err)

			// Check TTL on the timestamp key (which has the max refill time)
			tsKey := keyPrefix + ":test-subject:ts"
			tsTTL, err := client.PTTL(context.Background(), tsKey).Result()
			require.NoError(t, err)
			require.Greater(t, tsTTL, time.Duration(0), "timestamp key should have a TTL set")
			require.InDelta(
				t,
				float64(tt.wantTSTTL.Milliseconds()),
				float64(tsTTL.Milliseconds()),
				float64(tt.wantTTLDelta.Milliseconds()),
				"timestamp TTL should match max refill time",
			)

			// Check TTL on each individual limit key
			for i, wantLimitTTL := range tt.wantLimitTTLs {
				limitKey := keyPrefix + ":test-subject:" + strconv.Itoa(i+1)
				limitTTL, err := client.PTTL(context.Background(), limitKey).Result()
				require.NoError(t, err)
				require.Greater(
					t,
					limitTTL,
					time.Duration(0),
					"limit %d key should have a TTL set",
					i+1,
				)
				require.InDelta(
					t,
					float64(wantLimitTTL.Milliseconds()),
					float64(limitTTL.Milliseconds()),
					float64(tt.wantTTLDelta.Milliseconds()),
					"limit %d TTL should be based on time to refill", i+1,
				)
			}
		})
	}
}

func TestRedisLimiter_IndependentKeyExpiration(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	// Create limiter with two limits: short (1 second) and long (10 seconds)
	// After consuming 5 tokens:
	// - Limit 1: 5/10 tokens remaining, needs 500ms to refill
	// - Limit 2: 15/20 tokens remaining, needs 2500ms to refill
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix,
		[]ratelimiter.Limit{
			{Capacity: 10, RefillEvery: 1 * time.Second},
			{Capacity: 20, RefillEvery: 10 * time.Second},
		})
	require.NoError(t, err)

	// Make a request to consume 5 tokens from each bucket
	res, err := limiter.Allow(context.Background(), "test-subject", 5)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.InDelta(
		t,
		5.0,
		res.Balances[0].Remaining,
		0.01,
		"limit 1 should have 5 tokens remaining",
	)
	require.InDelta(
		t,
		15.0,
		res.Balances[1].Remaining,
		0.01,
		"limit 2 should have 15 tokens remaining",
	)

	// Check that both limit keys exist and have TTLs
	limit1Key := keyPrefix + ":test-subject:1"
	limit2Key := keyPrefix + ":test-subject:2"
	tsKey := keyPrefix + ":test-subject:ts"

	ttl1, err := client.PTTL(context.Background(), limit1Key).Result()
	require.NoError(t, err)
	require.Greater(t, ttl1, time.Duration(0), "limit 1 key should have TTL")

	ttl2, err := client.PTTL(context.Background(), limit2Key).Result()
	require.NoError(t, err)
	require.Greater(t, ttl2, time.Duration(0), "limit 2 key should have TTL")

	tsTTL, err := client.PTTL(context.Background(), tsKey).Result()
	require.NoError(t, err)
	require.Greater(t, tsTTL, time.Duration(0), "timestamp key should have TTL")

	// Verify specific TTL values
	// Limit 1: (10-5)/10 * 1000ms = 500ms to refill
	expectedTTL1 := 500 * time.Millisecond
	require.InDelta(t,
		float64(expectedTTL1.Milliseconds()),
		float64(ttl1.Milliseconds()),
		float64(100*time.Millisecond.Milliseconds()),
		"limit 1 should expire after time to refill 5 tokens (~500ms)")

	// Limit 2: (20-15)/20 * 10000ms = 2500ms to refill
	expectedTTL2 := 2500 * time.Millisecond
	require.InDelta(t,
		float64(expectedTTL2.Milliseconds()),
		float64(ttl2.Milliseconds()),
		float64(200*time.Millisecond.Milliseconds()),
		"limit 2 should expire after time to refill 5 tokens (~2500ms)")

	// Timestamp key should have the longest TTL (10 seconds = max refill time)
	expectedTSTTL := 10 * time.Second
	require.InDelta(t,
		float64(expectedTSTTL.Milliseconds()),
		float64(tsTTL.Milliseconds()),
		float64(200*time.Millisecond.Milliseconds()),
		"timestamp should expire at max refill time (10s)")

	// Verify that limit 1 expires much sooner than limit 2
	require.Less(t, ttl1, ttl2, "short refill limit should expire before long refill limit")

	// Verify both limits expire before timestamp
	require.Less(t, ttl1, tsTTL, "limit 1 should expire before timestamp")
	require.Less(t, ttl2, tsTTL, "limit 2 should expire before timestamp")
}

func TestRedisLimiter_Refill(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix,
		[]ratelimiter.Limit{{Capacity: 10, RefillEvery: 100 * time.Millisecond}})
	require.NoError(t, err)

	// Consume all tokens
	res, err := limiter.Allow(context.Background(), "test-subject", 10)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.InDelta(t, 0, res.Balances[0].Remaining, 0.01)

	// Should be blocked immediately
	res, err = limiter.Allow(context.Background(), "test-subject", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed)

	// Wait for partial refill
	time.Sleep(50 * time.Millisecond)
	res, err = limiter.Allow(context.Background(), "test-subject", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed, "should allow after partial refill")
	require.Greater(t, res.Balances[0].Remaining, 3.0, "should have refilled ~5 tokens")
}

func TestRedisLimiter_RetryAfter(t *testing.T) {
	tests := []struct {
		name              string
		limits            []ratelimiter.Limit
		cost              uint64
		wantRetryAfter    time.Duration
		wantRetryAfterMax time.Duration
	}{
		{
			name: "exact cost equals capacity",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Second},
			},
			cost: 11,
			// Need 1 token, refill rate is 1 second / 10 tokens = 100ms per token
			wantRetryAfter:    100 * time.Millisecond,
			wantRetryAfterMax: 110 * time.Millisecond,
		},
		{
			name: "cost is double capacity",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Second},
			},
			cost: 30,
			// Need 10 tokens, refill rate is 1 second / 10 tokens = 100ms per token
			wantRetryAfter:    time.Second,
			wantRetryAfterMax: 1100 * time.Millisecond,
		},
		{
			name: "partial tokens remaining",
			limits: []ratelimiter.Limit{
				{Capacity: 100, RefillEvery: 10 * time.Second},
			},
			cost: 50,
			// First consume 60 tokens, leaving 40, then try to consume 50
			// Need 10 more tokens, refill rate is 10s / 100 = 100ms per token
			wantRetryAfter:    time.Second,
			wantRetryAfterMax: 1100 * time.Millisecond,
		},
		{
			name: "multiple limits - first fails",
			limits: []ratelimiter.Limit{
				{Capacity: 5, RefillEvery: 500 * time.Millisecond},
				{Capacity: 100, RefillEvery: time.Hour},
			},
			cost: 10,
			// First limit needs 5 tokens, refill rate is 500ms / 5 = 100ms per token
			wantRetryAfter:    500 * time.Millisecond,
			wantRetryAfterMax: 550 * time.Millisecond,
		},
		{
			name: "cost exceeds capacity",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Second},
			},
			cost: 100,
			// Cost exceeds capacity, so need to wait full refill period
			wantRetryAfter:    time.Second,
			wantRetryAfterMax: 1100 * time.Millisecond,
		},
		{
			name: "cost far exceeds capacity",
			limits: []ratelimiter.Limit{
				{Capacity: 5, RefillEvery: 200 * time.Millisecond},
			},
			cost: 1000,
			// Cost far exceeds capacity, capped at refill period
			wantRetryAfter:    200 * time.Millisecond,
			wantRetryAfterMax: 220 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.NoError(t, err)

			// For "partial tokens remaining" test, first consume some tokens
			if tt.name == "partial tokens remaining" {
				_, err := limiter.Allow(context.Background(), "test-subject", 60)
				require.NoError(t, err)
			}

			// Make the request that should fail
			res, err := limiter.Allow(context.Background(), "test-subject", tt.cost)
			require.NoError(t, err)
			require.False(t, res.Allowed, "request should be blocked")
			require.NotNil(t, res.FailedLimit)
			require.NotNil(t, res.RetryAfter)

			// Verify RetryAfter is in expected range
			require.GreaterOrEqual(t, *res.RetryAfter, tt.wantRetryAfter,
				"RetryAfter should be at least %v, got %v", tt.wantRetryAfter, *res.RetryAfter)
			require.LessOrEqual(t, *res.RetryAfter, tt.wantRetryAfterMax,
				"RetryAfter should be at most %v, got %v", tt.wantRetryAfterMax, *res.RetryAfter)
		})
	}
}

func TestRedisLimiter_Errors(t *testing.T) {
	tests := []struct {
		name    string
		limits  []ratelimiter.Limit
		wantErr string
	}{
		{
			name:    "empty limits",
			limits:  []ratelimiter.Limit{},
			wantErr: "no limits provided",
		},
		{
			name:    "zero capacity",
			limits:  []ratelimiter.Limit{{Capacity: 0, RefillEvery: time.Minute}},
			wantErr: "invalid limit at index 0",
		},
		{
			name:    "negative capacity",
			limits:  []ratelimiter.Limit{{Capacity: -1, RefillEvery: time.Minute}},
			wantErr: "invalid limit at index 0",
		},
		{
			name:    "zero refill time",
			limits:  []ratelimiter.Limit{{Capacity: 10, RefillEvery: 0}},
			wantErr: "invalid limit at index 0",
		},
		{
			name:    "negative refill time",
			limits:  []ratelimiter.Limit{{Capacity: 10, RefillEvery: -1}},
			wantErr: "invalid limit at index 0",
		},
		{
			name: "invalid second limit",
			limits: []ratelimiter.Limit{
				{Capacity: 10, RefillEvery: time.Minute},
				{Capacity: 0, RefillEvery: time.Minute},
			},
			wantErr: "invalid limit at index 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, keyPrefix := redistestutils.NewRedisForTest(t)
			_, err := ratelimiter.NewRedisLimiter(client, keyPrefix, tt.limits)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestRedisLimiter_SubjectIsolation(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix,
		[]ratelimiter.Limit{{Capacity: 10, RefillEvery: time.Minute}})
	require.NoError(t, err)

	// Subject 1 consumes all tokens
	res, err := limiter.Allow(context.Background(), "subject1", 10)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// Subject 1 should be blocked
	res, err = limiter.Allow(context.Background(), "subject1", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed)

	// Subject 2 should still have full capacity
	res, err = limiter.Allow(context.Background(), "subject2", 10)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.InDelta(t, 0, res.Balances[0].Remaining, 0.01)
}

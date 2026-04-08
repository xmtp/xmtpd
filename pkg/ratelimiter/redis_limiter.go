package ratelimiter

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed script.lua
var luaScript string

type RedisLimiter struct {
	client    redis.UniversalClient
	script    *redis.Script
	keyPrefix string
	limits    []Limit
}

func NewRedisLimiter(
	client redis.UniversalClient,
	keyPrefix string,
	limits []Limit,
) (*RedisLimiter, error) {
	if err := validateLimits(limits); err != nil {
		return nil, err
	}

	return &RedisLimiter{
		client:    client,
		script:    redis.NewScript(luaScript),
		keyPrefix: keyPrefix,
		limits:    limits,
	}, nil
}

func (l *RedisLimiter) baseKey(subject string) string {
	return fmt.Sprintf("%s:%s", l.keyPrefix, subject)
}

func (l *RedisLimiter) buildKeys(subject string) []string {
	baseKey := l.baseKey(subject)
	// First key is timestamp, then one key per limit
	keys := make([]string, 1+len(l.limits))
	keys[0] = baseKey + ":ts"
	for i := range l.limits {
		keys[i+1] = fmt.Sprintf("%s:%d", baseKey, i+1)
	}
	return keys
}

func (l *RedisLimiter) buildArgs(requestTime time.Time, cost uint64) []any {
	args := make([]any, 0, 3+len(l.limits)*2)
	args = append(args, requestTime.UnixMilli(), len(l.limits), cost)
	for _, lim := range l.limits {
		args = append(args, lim.Capacity, lim.RefillEvery.Milliseconds())
	}

	return args
}

func (l *RedisLimiter) Allow(ctx context.Context, subject string, cost uint64) (*Result, error) {
	if cost == 0 {
		return nil, ErrCostMustBeGreaterThanZero
	}
	now := time.Now()
	keys := l.buildKeys(subject)
	args := l.buildArgs(now, cost)

	raw, err := l.script.Run(ctx, l.client, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	arr, ok := raw.([]any)
	if !ok || len(arr) < 1 {
		return nil, ErrUnexpectedScriptResponse
	}

	res, err := l.transformResult(arr, cost)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (l *RedisLimiter) transformResult(arr []any, cost uint64) (*Result, error) {
	var err error

	allowed := arr[0].(int64) == 1
	res := &Result{
		Allowed: allowed,
	}

	if allowed {
		// {1, rem1, rem2, ...}
		if res.Balances, err = l.parseBalances(arr[1:]); err != nil {
			return nil, err
		}
	} else {
		// {0, failed_index, rem1, rem2, ...}
		if len(arr) != 2+len(l.limits) {
			return nil, ErrUnexpectedScriptResponse
		}

		failedIndex := int(arr[1].(int64)) - 1 // Convert from 1-based to 0-based
		// Make sure the failed index is within bounds
		if failedIndex < 0 || failedIndex >= len(l.limits) {
			return nil, ErrInvalidFailedLimit
		}

		limit := l.limits[failedIndex]
		res.FailedLimit = &limit

		// Calculate how long until the limit refills enough for this request
		remaining := toFloat(arr[2+failedIndex])
		res.RetryAfter = calculateRetryAfter(limit, remaining, cost)

		if res.Balances, err = l.parseBalances(arr[2:]); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (l *RedisLimiter) parseBalances(arr []any) ([]LimitBalance, error) {
	if len(arr) != len(l.limits) {
		return nil, ErrUnexpectedScriptResponse
	}

	balances := make([]LimitBalance, len(l.limits))
	for i := range l.limits {
		balances[i] = LimitBalance{
			Limit:     l.limits[i],
			Remaining: toFloat(arr[i]),
		}
	}
	return balances, nil
}

func validateLimits(limits []Limit) error {
	if len(limits) == 0 {
		return ErrNoLimitsProvided
	}

	for i, lim := range limits {
		if lim.Capacity <= 0 || lim.RefillEvery <= 0 {
			return fmt.Errorf("invalid limit at index %d", i)
		}
	}

	return nil
}

func toFloat(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case int64:
		return float64(t)
	default:
		return 0
	}
}

// calculateRetryAfter computes the duration until a limit will have enough tokens
// for a request with the given cost to succeed.
// The returned duration is capped at RefillEvery since that's when the bucket is full.
func calculateRetryAfter(limit Limit, remaining float64, cost uint64) *time.Duration {
	tokensNeeded := float64(cost) - remaining
	if tokensNeeded <= 0 {
		return nil
	}

	// Cap tokens needed at capacity, since the bucket can't hold more
	if tokensNeeded > float64(limit.Capacity) {
		tokensNeeded = float64(limit.Capacity)
	}

	// Time to refill = (tokens_needed / capacity) * refill_every
	refillRate := float64(limit.RefillEvery.Nanoseconds()) / float64(limit.Capacity)
	retryAfterNanos := tokensNeeded * refillRate

	// RetryAfter cannot be longer than the time it takes to refill the bucket
	retryAfter := min(time.Duration(retryAfterNanos), limit.RefillEvery)

	return &retryAfter
}

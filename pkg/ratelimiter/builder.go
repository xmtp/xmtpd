package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils/clientip"
	"go.uber.org/zap"
)

// BuiltLimiter is the result of constructing the rate-limit subsystem at
// startup. It exposes a query limiter and an opens limiter (separate Redis
// key prefixes), the parsed trusted-proxy CIDRs, and is consumed by the server
// when constructing the rate-limit interceptor.
type BuiltLimiter struct {
	QueryLimiter RateLimiter // BreakerLimiter wrapping a RedisLimiter([per-minute, per-hour])
	OpensLimiter RateLimiter // BreakerLimiter wrapping a RedisLimiter([opens-per-minute])
	TrustedCIDRs []*net.IPNet
}

// Build constructs the rate-limit subsystem from server configuration. If
// rlOpts.Enable is false, returns (nil, nil) — the caller should treat the
// nil result as "rate limiting disabled."
//
// When enabled, Build pings Redis and returns an error if it is unreachable.
// This implements the spec's fail-fast-at-startup behavior.
func Build(
	ctx context.Context,
	logger *zap.Logger,
	redisOpts config.RedisOptions,
	rlOpts config.RateLimitOptions,
) (*BuiltLimiter, error) {
	if !rlOpts.Enable {
		return nil, nil
	}
	if redisOpts.RedisURL == "" {
		return nil, errors.New("rate limiting enabled but XMTPD_REDIS_URL is empty")
	}

	parsed, err := redis.ParseURL(redisOpts.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}
	client := redis.NewClient(parsed)

	pingCtx, cancel := context.WithTimeout(ctx, redisOpts.ConnectTimeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("rate limiting enabled but redis ping failed: %w", err)
	}

	queryInner, err := NewRedisLimiter(client, redisOpts.KeyPrefix+"rl:t2:q", []Limit{
		{Capacity: rlOpts.T2PerMinuteCapacity, RefillEvery: time.Minute},
		{Capacity: rlOpts.T2PerHourCapacity, RefillEvery: time.Hour},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to construct query limiter: %w", err)
	}
	opensInner, err := NewRedisLimiter(client, redisOpts.KeyPrefix+"rl:t2:o", []Limit{
		{Capacity: rlOpts.T2SubscribeOpensPerMinute, RefillEvery: time.Minute},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to construct opens limiter: %w", err)
	}

	queryWrapped := NewBreakerLimiter(
		queryInner,
		NewCircuitBreaker(rlOpts.BreakerFailureThreshold, rlOpts.BreakerCooldown),
	)
	opensWrapped := NewBreakerLimiter(
		opensInner,
		NewCircuitBreaker(rlOpts.BreakerFailureThreshold, rlOpts.BreakerCooldown),
	)

	cidrs, err := clientip.ParseTrustedProxyCIDRs(rlOpts.TrustedProxyCIDRs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trusted proxy CIDRs: %w", err)
	}

	logger.Info("rate limit interceptor enabled",
		zap.Int("t2_per_minute", rlOpts.T2PerMinuteCapacity),
		zap.Int("t2_per_hour", rlOpts.T2PerHourCapacity),
		zap.Int("t2_subscribe_opens_per_minute", rlOpts.T2SubscribeOpensPerMinute),
	)
	return &BuiltLimiter{
		QueryLimiter: queryWrapped,
		OpensLimiter: opensWrapped,
		TrustedCIDRs: cidrs,
	}, nil
}

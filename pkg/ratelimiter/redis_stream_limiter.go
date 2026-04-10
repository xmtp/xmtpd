package ratelimiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed stream_script.lua
var streamLuaScript string

// RedisStreamLimiter tracks concurrent stream counts using a Lua script
// for atomic acquire/release. Each subject (IP) gets a single Redis key
// whose value is the current stream count. A TTL on the key provides crash
// recovery: if a process dies without calling Release, the key expires and
// the slot is freed.
type RedisStreamLimiter struct {
	client    redis.UniversalClient
	script    *redis.Script
	keyPrefix string
	maxCount  int
	ttlMs     int64
}

// NewRedisStreamLimiter creates a RedisStreamLimiter.
//   - keyPrefix: prepended to the subject to form the Redis key (e.g. "xmtpd:rl:streams:").
//   - maxCount: maximum allowed concurrent streams per subject.
//   - ttl: Redis key TTL; acts as the crash self-heal window.
func NewRedisStreamLimiter(
	client redis.UniversalClient,
	keyPrefix string,
	maxCount int,
	ttl time.Duration,
) *RedisStreamLimiter {
	return &RedisStreamLimiter{
		client:    client,
		script:    redis.NewScript(streamLuaScript),
		keyPrefix: keyPrefix,
		maxCount:  maxCount,
		ttlMs:     ttl.Milliseconds(),
	}
}

func (l *RedisStreamLimiter) key(subject string) string {
	return l.keyPrefix + subject
}

// Acquire atomically checks the current count and increments only if below
// maxCount. Returns allowed=true if the stream was admitted.
func (l *RedisStreamLimiter) Acquire(ctx context.Context, subject string) (bool, error) {
	result, err := l.script.Run(
		ctx, l.client, []string{l.key(subject)},
		"acquire", l.maxCount, l.ttlMs,
	).Int64()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Release atomically decrements the stream count, clamped at zero.
func (l *RedisStreamLimiter) Release(ctx context.Context, subject string) error {
	_, err := l.script.Run(
		ctx, l.client, []string{l.key(subject)},
		"release", 0, l.ttlMs,
	).Result()
	return err
}

// RefreshTTL resets the key's TTL to keep it alive while a stream is open.
func (l *RedisStreamLimiter) RefreshTTL(ctx context.Context, subject string) error {
	return l.client.PExpire(ctx, l.key(subject), time.Duration(l.ttlMs)*time.Millisecond).Err()
}

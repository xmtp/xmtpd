package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStreamLimiter tracks concurrent stream counts using Redis INCR/DECR.
// Each subject (IP) gets a single Redis key whose value is the current stream
// count. A TTL on the key provides crash recovery: if a process dies without
// calling Release, the key expires and the slot is freed.
type RedisStreamLimiter struct {
	client    redis.UniversalClient
	keyPrefix string
	maxCount  int64
	ttl       time.Duration
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
		keyPrefix: keyPrefix,
		maxCount:  int64(maxCount),
		ttl:       ttl,
	}
}

func (l *RedisStreamLimiter) key(subject string) string {
	return l.keyPrefix + subject
}

// Acquire atomically increments the stream count for the subject. If the new
// count exceeds maxCount, it immediately decrements and returns allowed=false.
func (l *RedisStreamLimiter) Acquire(ctx context.Context, subject string) (bool, error) {
	k := l.key(subject)

	count, err := l.client.Incr(ctx, k).Result()
	if err != nil {
		return false, err
	}

	if count > l.maxCount {
		// Over limit — roll back the increment.
		l.client.Decr(ctx, k)
		return false, nil
	}

	// Set/refresh TTL on successful acquire.
	l.client.Expire(ctx, k, l.ttl)
	return true, nil
}

// Release decrements the stream count for the subject. The count is clamped
// at zero to prevent negative drift from orphan releases.
func (l *RedisStreamLimiter) Release(ctx context.Context, subject string) error {
	k := l.key(subject)

	val, err := l.client.Decr(ctx, k).Result()
	if err != nil {
		return err
	}

	// Clamp at zero: if DECR produced a negative value, reset to 0.
	if val < 0 {
		l.client.Set(ctx, k, 0, l.ttl)
	}
	return nil
}

// RefreshTTL resets the key's TTL to keep it alive while a stream is open.
func (l *RedisStreamLimiter) RefreshTTL(ctx context.Context, subject string) error {
	return l.client.Expire(ctx, l.key(subject), l.ttl).Err()
}

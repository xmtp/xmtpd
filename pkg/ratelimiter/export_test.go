package ratelimiter

import "time"

// SetNowForTest overrides the RedisLimiter clock. Test-only.
func (l *RedisLimiter) SetNowForTest(now func() time.Time) {
	l.now = now
}

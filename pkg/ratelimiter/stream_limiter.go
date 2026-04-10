package ratelimiter

import (
	"context"

	"go.uber.org/zap"
)

// StreamLimiter tracks concurrent stream counts per subject.
// Acquire increments the count and returns whether the stream is allowed.
// Release decrements the count when the stream closes.
// RefreshTTL extends the key's TTL while the stream is alive.
type StreamLimiter interface {
	Acquire(ctx context.Context, subject string) (allowed bool, err error)
	Release(ctx context.Context, subject string) error
	RefreshTTL(ctx context.Context, subject string) error
}

// BreakerStreamLimiter wraps a StreamLimiter with a CircuitBreaker.
// On Redis errors, Acquire fails open (allowed=true), Release and RefreshTTL
// errors are swallowed. When the breaker is OPEN, calls bypass Redis entirely.
type BreakerStreamLimiter struct {
	inner   StreamLimiter
	breaker *CircuitBreaker
	logger  *zap.Logger
}

// NewBreakerStreamLimiter wraps inner with the provided CircuitBreaker.
func NewBreakerStreamLimiter(
	inner StreamLimiter,
	breaker *CircuitBreaker,
	logger *zap.Logger,
) *BreakerStreamLimiter {
	return &BreakerStreamLimiter{inner: inner, breaker: breaker, logger: logger}
}

func (b *BreakerStreamLimiter) Acquire(ctx context.Context, subject string) (bool, error) {
	if !b.breaker.Allow() {
		return true, nil // fail open
	}
	allowed, err := b.inner.Acquire(ctx, subject)
	if err != nil {
		b.breaker.RecordFailure()
		return true, nil //nolint:nilerr // intentional fail-open
	}
	b.breaker.RecordSuccess()
	return allowed, nil
}

func (b *BreakerStreamLimiter) Release(ctx context.Context, subject string) error {
	if !b.breaker.Allow() {
		return nil // breaker open, skip
	}
	err := b.inner.Release(ctx, subject)
	if err != nil {
		b.breaker.RecordFailure()
		b.logger.Warn("stream limiter release failed (breaker)",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return nil //nolint:nilerr // best-effort release
	}
	b.breaker.RecordSuccess()
	return nil
}

func (b *BreakerStreamLimiter) RefreshTTL(ctx context.Context, subject string) error {
	if !b.breaker.Allow() {
		return nil // breaker open, skip
	}
	err := b.inner.RefreshTTL(ctx, subject)
	if err != nil {
		b.breaker.RecordFailure()
		return nil //nolint:nilerr // best-effort refresh
	}
	b.breaker.RecordSuccess()
	return nil
}

package ratelimiter

import (
	"context"
	"sync"
	"time"
)

// BreakerState represents the current state of a CircuitBreaker.
type BreakerState int

const (
	BreakerClosed   BreakerState = iota // normal operation
	BreakerOpen                         // short-circuiting all calls
	BreakerHalfOpen                     // probing with one call after cooldown
)

func (s BreakerState) String() string {
	switch s {
	case BreakerClosed:
		return "closed"
	case BreakerOpen:
		return "open"
	case BreakerHalfOpen:
		return "half_open"
	default:
		return "unknown"
	}
}

// CircuitBreaker is a simple consecutive-failure circuit breaker.
//
// Closed: every call passes. Failures increment a counter; success resets it.
// Open: every call is short-circuited (Allow returns false) until cooldown
// elapses, then transitions to HalfOpen.
// HalfOpen: the next call is allowed as a probe. Success → Closed.
// Failure → Open with the cooldown timer reset.
type CircuitBreaker struct {
	mu               sync.Mutex
	failureThreshold int
	cooldown         time.Duration

	state        BreakerState
	failureCount int
	openedAt     time.Time
}

// NewCircuitBreaker creates a CircuitBreaker that opens after failureThreshold
// consecutive failures and attempts to recover after cooldown has elapsed.
func NewCircuitBreaker(failureThreshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		cooldown:         cooldown,
		state:            BreakerClosed,
	}
}

// State returns the current breaker state.
func (cb *CircuitBreaker) State() BreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Allow reports whether the request should be forwarded to the inner resource.
// When the breaker is Open and the cooldown has elapsed, it transitions to
// HalfOpen and returns true for the probe call.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case BreakerClosed:
		return true
	case BreakerHalfOpen:
		return true
	case BreakerOpen:
		if time.Since(cb.openedAt) >= cb.cooldown {
			cb.state = BreakerHalfOpen
			BreakerStateGauge.Set(1)
			return true
		}
		return false
	}
	return true
}

// RecordSuccess resets the failure count and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount = 0
	if cb.state != BreakerClosed {
		cb.state = BreakerClosed
		BreakerStateGauge.Set(0)
	}
}

// RecordFailure increments the consecutive failure count. If the count reaches
// the threshold (or the breaker is in HalfOpen), the breaker opens and the
// cooldown timer resets. Failures recorded while the breaker is already Open
// are ignored — they would otherwise double-count BreakerTripsTotal and reset
// openedAt repeatedly, extending the cooldown for inflight callers.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == BreakerOpen {
		return
	}
	if cb.state == BreakerHalfOpen {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		BreakerStateGauge.Set(2)
		BreakerTripsTotal.Inc()
		return
	}
	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		BreakerStateGauge.Set(2)
		BreakerTripsTotal.Inc()
	}
}

// BreakerLimiter wraps a RateLimiter with a circuit breaker. On any error from
// the inner limiter, the breaker counts a failure and the request fails open
// (Allowed=true). When the breaker is OPEN, calls bypass the inner limiter
// entirely and fail open. Denials from the inner limiter are not failures.
type BreakerLimiter struct {
	inner   RateLimiter
	breaker *CircuitBreaker
}

// NewBreakerLimiter wraps inner with the provided CircuitBreaker.
func NewBreakerLimiter(inner RateLimiter, breaker *CircuitBreaker) *BreakerLimiter {
	return &BreakerLimiter{inner: inner, breaker: breaker}
}

// Allow implements RateLimiter. Errors from the inner limiter increment the
// breaker and return a fail-open result. Denials are passed through as-is.
func (b *BreakerLimiter) Allow(ctx context.Context, subject string, cost uint64) (*Result, error) {
	if !b.breaker.Allow() {
		return &Result{Allowed: true}, nil // fail open
	}
	res, err := b.inner.Allow(ctx, subject, cost)
	if err != nil {
		b.breaker.RecordFailure()
		// Fail open: swallow the error and admit the request. This is the
		// whole point of the breaker — Redis outages must not block traffic.
		return &Result{Allowed: true}, nil //nolint:nilerr // intentional fail-open
	}
	b.breaker.RecordSuccess()
	return res, nil
}

// ForceDebit implements RateLimiter with the same fail-open semantics as Allow.
func (b *BreakerLimiter) ForceDebit(
	ctx context.Context,
	subject string,
	cost uint64,
) (*Result, error) {
	if !b.breaker.Allow() {
		return &Result{Allowed: true}, nil
	}
	res, err := b.inner.ForceDebit(ctx, subject, cost)
	if err != nil {
		b.breaker.RecordFailure()
		return &Result{Allowed: true}, nil //nolint:nilerr // intentional fail-open
	}
	b.breaker.RecordSuccess()
	return res, nil
}

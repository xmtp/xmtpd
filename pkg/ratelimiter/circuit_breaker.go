package ratelimiter

import (
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
	cb.state = BreakerClosed
}

// RecordFailure increments the consecutive failure count. If the count reaches
// the threshold (or the breaker is in HalfOpen), the breaker opens and the
// cooldown timer resets.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == BreakerHalfOpen {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		return
	}
	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
	}
}

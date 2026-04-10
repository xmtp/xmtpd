package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var errSentinel = errors.New("sentinel")

func TestCircuitBreaker_OpensAfterThresholdFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	require.Equal(t, BreakerClosed, cb.State())

	for range 3 {
		require.True(t, cb.Allow())
		cb.RecordFailure()
	}

	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())
}

func TestCircuitBreaker_HalfOpenAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	fakeNow := time.Now()
	cb.now = func() time.Time { return fakeNow }
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())

	// Advance the fake clock past the cooldown without sleeping.
	fakeNow = fakeNow.Add(80 * time.Millisecond)

	require.True(t, cb.Allow())
	require.Equal(t, BreakerHalfOpen, cb.State())
}

// TestCircuitBreaker_RecordFailureWhileOpenIsIdempotent is a regression for
// PR #1938 macroscope Medium: when several inflight calls fail after the
// breaker has already opened, RecordFailure must not double-count the trip
// in BreakerTripsTotal nor reset openedAt (which would extend the cooldown).
func TestCircuitBreaker_RecordFailureWhileOpenIsIdempotent(t *testing.T) {
	cb := NewCircuitBreaker(1, time.Hour)
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	openedAt := cb.openedAt
	failureCount := cb.failureCount

	// Simulate further inflight calls failing while already open.
	for range 10 {
		cb.RecordFailure()
	}
	require.Equal(t, BreakerOpen, cb.State())
	require.Equal(
		t,
		openedAt,
		cb.openedAt,
		"openedAt must not be reset by failures while already open",
	)
	require.Equal(t, failureCount, cb.failureCount, "failureCount must not grow while already open")
}

func TestCircuitBreaker_HalfOpenSuccessClosesCircuit(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	fakeNow := time.Now()
	cb.now = func() time.Time { return fakeNow }
	cb.RecordFailure()
	fakeNow = fakeNow.Add(80 * time.Millisecond)

	require.True(t, cb.Allow())
	cb.RecordSuccess()
	require.Equal(t, BreakerClosed, cb.State())
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	fakeNow := time.Now()
	cb.now = func() time.Time { return fakeNow }
	cb.RecordFailure()
	fakeNow = fakeNow.Add(80 * time.Millisecond)

	require.True(t, cb.Allow())
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	cb.RecordFailure()
	require.Equal(t, BreakerClosed, cb.State())
}

// fakeLimiter is a test double for RateLimiter.
type fakeLimiter struct {
	allowResult *Result
	allowErr    error
}

func (f *fakeLimiter) Allow(_ context.Context, _ string, _ uint64) (*Result, error) {
	return f.allowResult, f.allowErr
}

func TestBreakerLimiter_FailOpenWhenBreakerOpen(t *testing.T) {
	inner := &fakeLimiter{allowErr: errSentinel}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(1, time.Hour))

	// First call: inner errors → breaker opens but call fails open
	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// Subsequent calls: breaker is open → bypass inner entirely, fail open
	inner.allowErr = nil
	inner.allowResult = &Result{Allowed: false}
	res, err = bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)
}

func TestBreakerLimiter_PassesThroughOnSuccess(t *testing.T) {
	inner := &fakeLimiter{allowResult: &Result{Allowed: true}}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(3, time.Hour))

	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)
}

func TestBreakerLimiter_PassesThroughOnDenial(t *testing.T) {
	inner := &fakeLimiter{allowResult: &Result{Allowed: false}}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(3, time.Hour))

	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed) // denial is not a failure
}

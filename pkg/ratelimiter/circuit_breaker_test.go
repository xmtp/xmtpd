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

	for i := 0; i < 3; i++ {
		require.True(t, cb.Allow())
		cb.RecordFailure()
	}

	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())
}

func TestCircuitBreaker_HalfOpenAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())

	time.Sleep(80 * time.Millisecond)

	require.True(t, cb.Allow())
	require.Equal(t, BreakerHalfOpen, cb.State())
}

func TestCircuitBreaker_HalfOpenSuccessClosesCircuit(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(80 * time.Millisecond)

	require.True(t, cb.Allow())
	cb.RecordSuccess()
	require.Equal(t, BreakerClosed, cb.State())
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(80 * time.Millisecond)

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
	debitResult *Result
	debitErr    error
}

func (f *fakeLimiter) Allow(ctx context.Context, subject string, cost uint64) (*Result, error) {
	return f.allowResult, f.allowErr
}
func (f *fakeLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*Result, error) {
	return f.debitResult, f.debitErr
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

func TestBreakerLimiter_ForceDebitFailOpen(t *testing.T) {
	inner := &fakeLimiter{debitErr: errSentinel}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(1, time.Hour))
	res, err := bl.ForceDebit(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)
}

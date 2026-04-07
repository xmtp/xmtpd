package ratelimiter

import (
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

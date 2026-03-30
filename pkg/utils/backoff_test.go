package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryWithBackoff_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	attempts, err := RetryWithBackoff(func() error {
		calls++
		return nil
	}, 3, 50*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, 1, attempts)
	assert.Equal(t, 1, calls)
}

func TestRetryWithBackoff_SucceedsOnRetry(t *testing.T) {
	calls := 0
	attempts, err := RetryWithBackoff(func() error {
		calls++
		if calls < 2 {
			return errors.New("transient error")
		}
		return nil
	}, 3, 50*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, 2, attempts)
	assert.Equal(t, 2, calls)
}

func TestRetryWithBackoff_ExhaustsRetries(t *testing.T) {
	calls := 0
	attempts, err := RetryWithBackoff(func() error {
		calls++
		return errors.New("persistent error")
	}, 3, 50*time.Millisecond)
	require.Error(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, 3, calls)
	assert.Contains(t, err.Error(), "persistent error")
}

func TestRetryWithBackoff_SucceedsOnLastAttempt(t *testing.T) {
	calls := 0
	attempts, err := RetryWithBackoff(func() error {
		calls++
		if calls < 3 {
			return errors.New("transient error")
		}
		return nil
	}, 3, 50*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, 3, calls)
}

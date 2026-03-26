package redis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryWithBackoff_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	attempts, err := retryWithBackoff(func() error {
		calls++
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, attempts)
	assert.Equal(t, 1, calls)
}

func TestRetryWithBackoff_SucceedsOnRetry(t *testing.T) {
	calls := 0
	attempts, err := retryWithBackoff(func() error {
		calls++
		if calls < 2 {
			return errors.New("transient error")
		}
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 2, attempts)
	assert.Equal(t, 2, calls)
}

func TestRetryWithBackoff_ExhaustsRetries(t *testing.T) {
	calls := 0
	attempts, err := retryWithBackoff(func() error {
		calls++
		return errors.New("persistent error")
	})
	require.Error(t, err)
	assert.Equal(t, maxRetries, attempts)
	assert.Equal(t, maxRetries, calls)
	assert.Contains(t, err.Error(), "persistent error")
}

func TestRetryWithBackoff_SucceedsOnLastAttempt(t *testing.T) {
	calls := 0
	attempts, err := retryWithBackoff(func() error {
		calls++
		if calls < maxRetries {
			return errors.New("transient error")
		}
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, maxRetries, attempts)
	assert.Equal(t, maxRetries, calls)
}

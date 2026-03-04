package utils

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// NewBackoff creates an ExponentialBackOff with the given timing parameters.
// Multiplier is fixed at 2.0 and randomization factor at 0.5.
// Set maxElapsedTime to 0 for no limit.
func NewBackoff(
	initialInterval, maxInterval, maxElapsedTime time.Duration,
) *backoff.ExponentialBackOff {
	return backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(initialInterval),
		backoff.WithMaxInterval(maxInterval),
		backoff.WithMultiplier(2.0),
		backoff.WithRandomizationFactor(0.5),
		backoff.WithMaxElapsedTime(maxElapsedTime),
	)
}

package ratelimiter

import "context"

// StreamLimiter tracks concurrent stream counts per subject.
// Acquire increments the count and returns whether the stream is allowed.
// Release decrements the count when the stream closes.
// RefreshTTL extends the key's TTL while the stream is alive.
type StreamLimiter interface {
	Acquire(ctx context.Context, subject string) (allowed bool, err error)
	Release(ctx context.Context, subject string) error
	RefreshTTL(ctx context.Context, subject string) error
}

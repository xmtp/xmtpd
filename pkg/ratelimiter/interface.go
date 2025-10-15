package ratelimiter

import (
	"context"
	"time"
)

type Limit struct {
	Capacity    int           // Maximum number of tokens in the bucket (e.g. 10)
	RefillEvery time.Duration // Time to fully refill the bucket from 0 to Capacity (e.g. 1 * time.Minute)
}

type LimitBalance struct {
	Limit     Limit
	Remaining float64 // remaining tokens
}

type Result struct {
	Allowed     bool
	FailedLimit *Limit         // The limit that was exceeded (nil if allowed)
	RetryAfter  *time.Duration // Time until the failed limit will allow the request (nil if allowed)
	Balances    []LimitBalance // remaining tokens after this check (post-decrement on success; current on reject)
}

type RateLimiter interface {
	Allow(ctx context.Context, subject string, cost uint64) (*Result, error)
}

// Package ratelimiter provides a rate limiter interface and implementation.
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

// RateLimiter is the basic rate limiter interface for single-subject limiting.
type RateLimiter interface {
	Allow(ctx context.Context, subject string, cost uint64) (*Result, error)
}

// DualResult represents the result of checking both gateway and user limits.
type DualResult struct {
	Allowed       bool
	GatewayResult *Result // Result from gateway limit check
	UserResult    *Result // Result from user limit check (nil if no user specified)
	FailedSubject string  // "gateway" or "user" if denied
}

// DualRateLimiter applies rate limits to both gateway and user in parallel.
// This is used for delegated signing where we want to limit both:
// - The gateway's overall throughput
// - Individual user's message rate
type DualRateLimiter interface {
	// AllowDual checks both gateway and user limits in parallel.
	// If userSubject is empty, only gateway limit is checked.
	// Returns denied if either limit is exceeded.
	AllowDual(
		ctx context.Context,
		gatewaySubject, userSubject string,
		cost uint64,
	) (*DualResult, error)
}

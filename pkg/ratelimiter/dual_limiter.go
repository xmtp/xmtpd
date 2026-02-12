package ratelimiter

import (
	"context"
	"fmt"
)

// DualRedisLimiter implements DualRateLimiter using two RedisLimiter instances.
// It applies separate limits to gateway and user subjects in parallel.
type DualRedisLimiter struct {
	gatewayLimiter *RedisLimiter
	userLimiter    *RedisLimiter
}

// NewDualRedisLimiter creates a new DualRedisLimiter with separate limits for gateway and user.
// gatewayLimits are applied to the gateway subject (e.g., gateway address)
// userLimits are applied to individual user subjects (e.g., user payer address)
func NewDualRedisLimiter(
	gatewayLimiter *RedisLimiter,
	userLimiter *RedisLimiter,
) (*DualRedisLimiter, error) {
	if gatewayLimiter == nil {
		return nil, fmt.Errorf("gateway limiter cannot be nil")
	}
	// userLimiter can be nil if we don't want per-user limits
	return &DualRedisLimiter{
		gatewayLimiter: gatewayLimiter,
		userLimiter:    userLimiter,
	}, nil
}

// AllowDual checks both gateway and user limits.
// If userSubject is empty, only gateway limit is checked.
// Both limits must pass for the request to be allowed.
func (d *DualRedisLimiter) AllowDual(
	ctx context.Context,
	gatewaySubject, userSubject string,
	cost uint64,
) (*DualResult, error) {
	result := &DualResult{
		Allowed: true,
	}

	// Check gateway limit first
	gatewayResult, err := d.gatewayLimiter.Allow(ctx, gatewaySubject, cost)
	if err != nil {
		return nil, fmt.Errorf("gateway rate limit check failed: %w", err)
	}
	result.GatewayResult = gatewayResult

	if !gatewayResult.Allowed {
		result.Allowed = false
		result.FailedSubject = "gateway"
		return result, nil
	}

	// If no user subject, we're done
	if userSubject == "" || d.userLimiter == nil {
		return result, nil
	}

	// Check user limit
	userResult, err := d.userLimiter.Allow(ctx, userSubject, cost)
	if err != nil {
		return nil, fmt.Errorf("user rate limit check failed: %w", err)
	}
	result.UserResult = userResult

	if !userResult.Allowed {
		result.Allowed = false
		result.FailedSubject = "user"
		// Note: We've already consumed gateway tokens. This is intentional -
		// if a user is rate-limited, the gateway still sees their request.
		// An alternative would be to check limits without consuming first.
		return result, nil
	}

	return result, nil
}

// AllowGatewayOnly is a convenience method for non-delegated requests.
// It only checks the gateway limit.
func (d *DualRedisLimiter) AllowGatewayOnly(
	ctx context.Context,
	gatewaySubject string,
	cost uint64,
) (*Result, error) {
	return d.gatewayLimiter.Allow(ctx, gatewaySubject, cost)
}

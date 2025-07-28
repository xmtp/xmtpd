package ratelimiter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limit defines a rate limit with a maximum count and time window
type Limit struct {
	Limit  int
	Window time.Duration
}

// RateLimiter interface defines the contract for rate limiting
type RateLimiter interface {
	Spend(ctx context.Context, identifier string, cost int) (bool, error)
}

// redisRateLimiter implements RateLimiter using Redis with a sliding window algorithm
type redisRateLimiter struct {
	client    redis.UniversalClient
	limits    []Limit
	keyPrefix string
	luaScript *redis.Script
}

// CreateRateLimiterFn is the function signature for creating a rate limiter
type CreateRateLimiterFn func(client redis.UniversalClient, limits []Limit) (*RateLimiter, error)

// Options contains configuration options for the rate limiter
type Options struct {
	// KeyPrefix is the prefix for all Redis keys used by the rate limiter
	// Defaults to "xmtpd:ratelimit:" if not specified
	KeyPrefix string
}

// NewRateLimiter creates a new Redis-backed rate limiter
func NewRateLimiter(client redis.UniversalClient, limits []Limit) (RateLimiter, error) {
	return NewRateLimiterWithOptions(client, limits, Options{})
}

// NewRateLimiterWithOptions creates a new Redis-backed rate limiter with custom options
func NewRateLimiterWithOptions(
	client redis.UniversalClient,
	limits []Limit,
	opts Options,
) (RateLimiter, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is required")
	}
	if len(limits) == 0 {
		return nil, fmt.Errorf("at least one limit must be specified")
	}

	// Validate limits
	for i, limit := range limits {
		if limit.Limit <= 0 {
			return nil, fmt.Errorf("limit %d: limit must be positive", i)
		}
		if limit.Window <= 0 {
			return nil, fmt.Errorf("limit %d: window must be positive", i)
		}
	}

	// Set default key prefix if not provided
	keyPrefix := opts.KeyPrefix
	if keyPrefix == "" {
		keyPrefix = "xmtpd:ratelimit:"
	}

	rl := &redisRateLimiter{
		client:    client,
		limits:    limits,
		keyPrefix: keyPrefix,
		luaScript: redis.NewScript(luaScript),
	}

	return rl, nil
}

// generateKey creates a unique Redis key for a specific identifier and limit combination.
// It uses a hash to keep keys short and avoid issues with special characters in identifiers.
func (rl *redisRateLimiter) generateKey(identifier string, limit Limit) string {
	// Create a unique key for each identifier+limit combination
	var keyBytes []byte
	keyBytes = fmt.Appendf(keyBytes, "%s:%d:%d", identifier, limit.Limit, limit.Window)
	hash := sha256.Sum256(keyBytes)
	return rl.keyPrefix + hex.EncodeToString(hash[:8])
}

// generateKeys creates Redis keys for all limits for a given identifier
func (rl *redisRateLimiter) generateKeys(identifier string) []string {
	keys := make([]string, len(rl.limits))
	for i, limit := range rl.limits {
		keys[i] = rl.generateKey(identifier, limit)
	}
	return keys
}

// Spend attempts to consume the specified cost from the rate limiter
func (rl *redisRateLimiter) Spend(ctx context.Context, identifier string, cost int) (bool, error) {
	if cost <= 0 {
		return false, fmt.Errorf("cost must be positive")
	}

	// Generate keys for each limit
	keys := rl.generateKeys(identifier)

	// Prepare arguments: currentTime, cost, then pairs of (limit, window)
	currentTime := time.Now().UnixMilli()
	args := make([]interface{}, 2+len(rl.limits)*2)
	args[0] = currentTime
	args[1] = cost
	for i, limit := range rl.limits {
		args[2+i*2] = limit.Limit
		args[2+i*2+1] = int64(limit.Window.Milliseconds())
	}

	// Execute the Lua script
	result, err := rl.luaScript.Run(ctx, rl.client, keys, args...).Result()
	if err != nil {
		return false, fmt.Errorf("failed to execute rate limit check: %w", err)
	}

	// Parse the result
	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) == 0 {
		return false, fmt.Errorf("unexpected result format from Lua script")
	}

	allowed, ok := resultSlice[0].(int64)
	if !ok {
		return false, fmt.Errorf("unexpected result type from Lua script")
	}

	return allowed == 1, nil
}

// Lua script for atomic multi-limit checking using leaky bucket algorithm
const luaScript = `
-- Keys: One key per rate limit configuration
-- Args: currentTime, cost, then pairs of (limit, window) for each rate limit

local currentTime = tonumber(ARGV[1])
local cost = tonumber(ARGV[2])
local numLimits = (#ARGV - 2) / 2

-- Helper function to calculate current tokens
local function getCurrentTokens(tokens, lastRefill, limit, window, currentTime)
    local elapsed = math.max(0, currentTime - lastRefill)
    local refillRate = limit / window
    local tokensToAdd = elapsed * refillRate
    return math.min(limit, tokens + tokensToAdd)
end

-- First, check all limits without modifying anything
for i = 1, numLimits do
    local key = KEYS[i]
    local limit = tonumber(ARGV[2 + (i-1)*2 + 1])
    local window = tonumber(ARGV[2 + (i-1)*2 + 2])
    
    -- Get current bucket state
    local bucket = redis.call('HMGET', key, 'tokens', 'lastRefill')
    local tokens = tonumber(bucket[1]) or limit
    local lastRefill = tonumber(bucket[2]) or currentTime
    
    -- Calculate current available tokens
    tokens = getCurrentTokens(tokens, lastRefill, limit, window, currentTime)
    
    -- Check if we have enough tokens
    if tokens < cost then
        return {0} -- Failed
    end
end

-- All limits pass, now deduct tokens from all buckets
for i = 1, numLimits do
    local key = KEYS[i]
    local limit = tonumber(ARGV[2 + (i-1)*2 + 1])
    local window = tonumber(ARGV[2 + (i-1)*2 + 2])
    
    -- Get current bucket state
    local bucket = redis.call('HMGET', key, 'tokens', 'lastRefill')
    local tokens = tonumber(bucket[1]) or limit
    local lastRefill = tonumber(bucket[2]) or currentTime
    
    -- Calculate current available tokens
    tokens = getCurrentTokens(tokens, lastRefill, limit, window, currentTime)
    
    -- Deduct the cost
    tokens = tokens - cost
    
    -- Update bucket state
    redis.call('HSET', key, 'tokens', tokens, 'lastRefill', currentTime)
    
    -- Set expiry to at least 2x the window to handle edge cases
    redis.call('EXPIRE', key, math.max(2, math.ceil(window / 500)))
end

return {1} -- Success
`

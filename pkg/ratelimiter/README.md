# Rate Limiter

A Redis-backed token bucket rate limiter that supports multiple concurrent limits with atomic operations.

## Features

- **Token Bucket Algorithm**: Continuously refilling token buckets with configurable capacity and refill rates
- **Multiple Concurrent Limits**: Check multiple rate limits atomically (e.g., per-minute and per-hour limits)
- **Atomic Operations**: All limits are checked and decremented atomically - if any limit fails, none are decremented
- **Variable Cost**: Support for requests with different costs (e.g., batch operations)
- **Retry-After Calculation**: Automatically calculates when a blocked request can be retried
- **Subject Isolation**: Each subject (user, IP, etc.) has independent rate limit tracking
- **TTL Management**: Keys automatically expire based on refill periods to minimize Redis memory usage

## Usage

### Basic Example

```go
import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/xmtp/xmtpd/pkg/ratelimiter"
    "go.uber.org/zap"
)

// Create a Redis client
client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Define rate limits
limits := []ratelimiter.Limit{
    {
        Capacity:    10,              // 10 requests
        RefillEvery: time.Minute,     // per minute
    },
    {
        Capacity:    100,             // 100 requests
        RefillEvery: time.Hour,       // per hour
    },
}

// Create the rate limiter
logger := zap.NewNop()
limiter, err := ratelimiter.NewRedisLimiter(
    logger,
    client,
    "myapp:ratelimit",  // Redis key prefix
    limits,
)
if err != nil {
    panic(err)
}

// Check if a request is allowed
ctx := context.Background()
result, err := limiter.Allow(ctx, "user:123", 1)
if err != nil {
    panic(err)
}

if result.Allowed {
    // Process the request
    fmt.Println("Request allowed")
    for i, balance := range result.Balances {
        fmt.Printf("Limit %d: %.2f tokens remaining\n", i, balance.Remaining)
    }
} else {
    // Reject the request
    fmt.Printf("Request blocked by limit: %+v\n", result.FailedLimit)
    if result.RetryAfter != nil {
        fmt.Printf("Retry after: %v\n", *result.RetryAfter)
    }
}
```

### Variable Cost Requests

For requests that consume multiple tokens (e.g., batch operations):

```go
// A batch request that costs 5 tokens
result, err := limiter.Allow(ctx, "user:123", 5)
```

### Handling Rate Limit Responses

The `Result` struct provides information about the rate limit check:

```go
type Result struct {
    Allowed     bool                // Whether the request is allowed
    FailedLimit *Limit              // The limit that was exceeded (nil if allowed)
    RetryAfter  *time.Duration      // Time until retry is possible (nil if allowed)
    Balances    []LimitBalance      // Current balance for each limit
}

type LimitBalance struct {
    Limit     Limit      // The limit configuration
    Remaining float64    // Remaining tokens after this check
}
```

## How It Works

### Token Bucket Algorithm

Each limit operates as a token bucket:

1. **Initial State**: Bucket starts full with `Capacity` tokens
2. **Request Processing**: Each request consumes `cost` tokens
3. **Refill**: Tokens continuously refill at a rate of `Capacity / RefillEvery`
4. **Blocking**: If insufficient tokens are available, the request is blocked

**Example:**
- Limit: 10 tokens per second
- Refill rate: 10 tokens / 1000ms = 0.01 tokens/ms = 1 token per 100ms
- If 7 tokens remain and a request costs 5 tokens:
  - After success: 2 tokens remain
  - After 800ms: 10 tokens available (bucket is full)

### Atomic Multi-Limit Checking

When multiple limits are configured, all limits are checked atomically:

1. Check all limits to ensure sufficient tokens
2. If **any** limit fails, **no** tokens are decremented
3. If **all** limits pass, **all** tokens are decremented

This ensures consistent state across all limits.

### Retry-After Calculation

When a request is blocked, `RetryAfter` tells you how long to wait:

```
RetryAfter = (tokens_needed - tokens_available) × (RefillEvery / Capacity)
```

The duration is capped at `RefillEvery` since the bucket can never exceed its capacity.

**Example:**
- Limit: 10 tokens per second
- Available: 3 tokens
- Cost: 5 tokens
- Calculation: (5 - 3) × (1000ms / 10) = 2 × 100ms = 200ms

### Redis Key Management

Keys are automatically managed with TTLs:

- **Timestamp key** (`prefix:subject:ts`): Expires after the longest `RefillEvery`
- **Limit keys** (`prefix:subject:1`, `prefix:subject:2`, etc.): Expire based on time to refill to capacity

This ensures Redis memory is automatically cleaned up for inactive subjects.

## Performance Considerations

### Redis Operations

Each `Allow()` call executes a single Lua script in Redis:
- **O(n)** complexity where n = number of limits
- Atomic execution (no race conditions)
- Typically completes in < 1ms

### Memory Usage

Redis memory usage depends on the number of keys and values stored. Each subject creates:
- 1 timestamp key
- N limit keys (where N = number of configured limits)

**Calculation for keyPrefix="rl", subject="192.168.1.1" (IPv4), 2 limits:**

Keys and values per subject:
1. `rl:192.168.1.1:ts` → `1735689600000` (13-byte timestamp)
2. `rl:192.168.1.1:1` → `7.5` (floating point, ~3-4 bytes)
3. `rl:192.168.1.1:2` → `95.0` (floating point, ~4 bytes)

**Memory overhead per key:**
Redis (v7+) with jemalloc allocator has significant per-key overhead:
- Redis object header: ~16 bytes
- String metadata (SDS): ~8 bytes
- Dict entry: ~24 bytes (can be 32-40 bytes with jemalloc alignment)
- **Base overhead per key: ~64 bytes** (before key/value data)

**Per-key calculation (approximate):**
- Timestamp key: 64 (overhead) + 18 (key) + 13 (value) = **~95 bytes**
- Limit key 1: 64 (overhead) + 16 (key) + 3 (value) = **~83 bytes**
- Limit key 2: 64 (overhead) + 16 (key) + 4 (value) = **~84 bytes**
- **Total per subject: ~262 bytes**

**For 10,000 active IPv4 subjects with 2 limits:**
- 10,000 × 262 bytes = **~2.62 MB**
```

**Notes:**
- Keys automatically expire based on inactivity, so memory usage reflects only active subjects

## Implementation Details

### Lua Script

The rate limiter uses a single Lua script (`script.lua`) for atomic operations:
1. Read current timestamp and token balances
2. Calculate refilled tokens based on elapsed time
3. Check if all limits have sufficient tokens
4. If yes: decrement all limits and return success
5. If no: return failure with the first failing limit

### Token Precision

Tokens are stored and calculated as floating-point numbers to support:
- Sub-second refill rates (e.g., 1 token per 100ms)
- Precise remaining token counts
- Accurate Retry-After calculations

## License

See the root LICENSE file for license information.

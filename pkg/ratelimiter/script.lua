-- ============================================================================
-- Multi-Limit Token Bucket Rate Limiter
-- ============================================================================
-- This Lua script implements a rate limiter that can enforce multiple limits
-- simultaneously using the token bucket algorithm with continuous refill.
--
-- INPUTS:
--   KEYS[1]           = Timestamp key for this subject (e.g., "rl:gateway:user123:ts")
--   KEYS[2..N+1]      = Token bucket keys for each limit (e.g., "rl:gateway:user123:1", ":2", ...)
--   ARGV[1]           = now_ms - Current timestamp in milliseconds
--   ARGV[2]           = num_limits (N) - Number of rate limits to enforce
--   ARGV[3]           = cost - Number of tokens to consume for this request
--   ARGV[4..4+2N-1]   = Flattened array of limit configurations:
--                       [capacity_1, refill_ms_1, capacity_2, refill_ms_2, ..., capacity_N, refill_ms_N]
--
-- REDIS STORAGE:
--   KEYS[1]     : Last update timestamp (integer ms) - expires at max(refill_ms)
--   KEYS[2..N+1]: Current token count for each limit (float) - expires at refill_ms[i] when full
--
-- RETURN VALUES:
--   On success: {1, remaining_1, remaining_2, ..., remaining_N}
--     - 1 indicates the request is allowed
--     - remaining_i is the token count AFTER deduction for each limit
--   On failure: {0, failed_index, remaining_1, remaining_2, ..., remaining_N}
--     - 0 indicates the request is denied
--     - failed_index is the 1-based index of the first limit that failed
--     - remaining_i is the CURRENT token count (no deduction) for each limit
--
-- BEHAVIOR:
--   - ALL limits must have sufficient tokens for the request to be allowed
--   - If ANY limit fails, NO tokens are deducted (atomic operation)
--   - Tokens refill continuously based on elapsed time since last update
--   - Each limit key expires independently when bucket is full for efficient cleanup
--   - Timestamp key expires at max(refill_ms) to track last activity
-- ============================================================================

-- Parse input arguments
local ts_key    = KEYS[1]           -- Timestamp key (shared across all limits)
local now_ms    = tonumber(ARGV[1]) -- Current time in milliseconds
local n         = tonumber(ARGV[2]) -- Number of limits to check
local cost      = tonumber(ARGV[3]) -- Tokens to consume

-- Initialize arrays to store limit configurations and current state
local caps      = {} -- Maximum capacity for each limit (e.g., 100 tokens)
local refill_ms = {} -- Time in ms to fully refill each limit (e.g., 60000 for 1 minute)
local tokens    = {} -- Current token count for each limit (after refill)

-- Parse limit configurations from ARGV
-- Arguments come in pairs: [capacity, refill_time, capacity, refill_time, ...]
local idx       = 4 -- Start after the first 3 arguments
for i = 1, n do
    caps[i] = tonumber(ARGV[idx])
    idx = idx + 1
    refill_ms[i] = tonumber(ARGV[idx])
    idx = idx + 1
end

-- ============================================================================
-- STEP 1: Fetch existing state from Redis
-- ============================================================================
-- Retrieve the shared timestamp and all token counts in a single pipeline
-- Using MGET: [ts, tokens_1, tokens_2, ..., tokens_N]
local keys_to_fetch = { ts_key }
for i = 1, n do
    table.insert(keys_to_fetch, KEYS[i + 1])
end

local raw = redis.call("MGET", unpack(keys_to_fetch))
local last_ts = raw[1] and tonumber(raw[1]) or now_ms -- Shared timestamp for all limits

-- ============================================================================
-- STEP 2: Calculate current token counts with refill
-- ============================================================================
-- For each limit, determine how many tokens are currently available by:
-- 1. Reading the stored token count (or using full capacity for new limits)
-- 2. Calculating how many tokens have refilled since the last update
-- 3. Capping the refilled amount at the maximum capacity
for i = 1, n do
    -- Extract token count from MGET response (offset by 1 for timestamp)
    local raw_tokens = raw[i + 1]

    -- Initialize with stored value, or full capacity if this is a new limit
    local t = raw_tokens and tonumber(raw_tokens) or caps[i]

    -- Calculate token refill based on elapsed time
    local delta = now_ms - last_ts -- Milliseconds since last update
    if delta < 0 then
        -- Clock skew protection: if timestamp is in the future, treat as no time passed
        delta = 0
    end

    if refill_ms[i] > 0 and delta > 0 then
        -- Calculate refill rate: tokens per millisecond
        -- Example: 100 tokens / 60000ms = 0.00167 tokens per millisecond
        local rate = caps[i] / refill_ms[i]

        -- Add refilled tokens based on elapsed time, but cap at max capacity
        -- Example: 0.00167 tokens/ms * 30000ms = 50 tokens refilled
        t = math.min(caps[i], t + rate * delta)
    end

    -- Store the calculated current state
    tokens[i] = t -- Current token count (may be fractional)
end

-- ============================================================================
-- STEP 3: Check if all limits can satisfy the request
-- ============================================================================
-- ALL limits must have enough tokens for the request to succeed.
-- If any limit fails, we track which one failed first (for debugging/reporting).
local failed_index = 0
for i = 1, n do
    if tokens[i] < cost then
        failed_index = i -- Store the 1-based index of the failing limit
        break
    end
end

-- ============================================================================
-- STEP 4: Apply the decision (all-or-nothing)
-- ============================================================================
if failed_index == 0 then
    -- ========================================================================
    -- SUCCESS PATH: All limits passed - deduct tokens and persist state
    -- ========================================================================

    -- Deduct tokens from each limit
    for i = 1, n do
        tokens[i] = tokens[i] - cost
    end

    -- Set timestamp with expiration in a single call
    -- Timestamp expires at the longest refill period across all limits
    local max_refill = math.max(unpack(refill_ms))
    if max_refill > 0 then
        redis.call("SET", ts_key, tostring(now_ms), "PX", max_refill)
    else
        redis.call("SET", ts_key, tostring(now_ms))
    end

    -- Set each limit key with its value and expiration in a single call
    -- When a bucket is full, it expires after its refill period
    -- When not full, set TTL based on time to refill to capacity
    for i = 1, n do
        local ttl
        if tokens[i] >= caps[i] then
            -- Bucket is full - expire after full refill period
            ttl = refill_ms[i]
        else
            -- Bucket not full - calculate time to refill to capacity
            -- TTL = time needed to refill from current level to full capacity
            local time_to_fill = (caps[i] - tokens[i]) * refill_ms[i] / caps[i]
            ttl = math.ceil(time_to_fill)
        end

        -- Use SET with PX to set value and expiration atomically
        redis.call("SET", KEYS[i + 1], tostring(tokens[i]), "PX", ttl)
    end

    -- Return success response: {1, remaining_tokens_1, remaining_tokens_2, ...}
    local resp = { 1 }                -- 1 = allowed
    for i = 1, n do
        table.insert(resp, tokens[i]) -- Remaining after deduction
    end
    return resp
else
    -- ========================================================================
    -- FAILURE PATH: At least one limit failed - reject without modifying state
    -- ========================================================================

    -- Return rejection response: {0, failed_index, current_tokens_1, current_tokens_2, ...}
    local resp = { 0, failed_index }  -- 0 = denied, index = which limit failed
    for i = 1, n do
        table.insert(resp, tokens[i]) -- Current tokens (no deduction)
    end
    return resp
end

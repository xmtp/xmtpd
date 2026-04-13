-- ============================================================================
-- Concurrent Stream Limiter
-- ============================================================================
-- Atomic acquire/release for a per-subject concurrent stream counter.
-- Single-key access per call (cluster-safe).
--
-- INPUTS:
--   KEYS[1] = Counter key (e.g., "xmtpd:rl:streams:10.0.0.1")
--   ARGV[1] = Operation: "acquire" or "release"
--   ARGV[2] = max_count (only used by acquire; ignored by release)
--   ARGV[3] = ttl_ms (TTL in milliseconds to set on the key)
--
-- RETURN VALUES:
--   acquire: 1 if allowed, 0 if denied
--   release: current count after release (clamped at 0)
-- ============================================================================

local key     = KEYS[1]
local op      = ARGV[1]
local max     = tonumber(ARGV[2])
local ttl_ms  = tonumber(ARGV[3])

if op == "acquire" then
    local count = tonumber(redis.call("GET", key) or "0")
    if count < max then
        count = redis.call("INCR", key)
        redis.call("PEXPIRE", key, ttl_ms)
        return 1
    end
    return 0
elseif op == "release" then
    local count = tonumber(redis.call("GET", key) or "0")
    if count > 0 then
        count = redis.call("DECR", key)
        redis.call("PEXPIRE", key, ttl_ms)
        return count
    end
    return 0
else
    return redis.error_reply("unknown operation: " .. tostring(op))
end

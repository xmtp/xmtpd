# Rate Limiting Spec — Node `QueryApi` (v1)

Implementation spec for node-side rate limiting on `QueryApi`. Derived from `tasks/rate-limiting-design.md`. Tracks issue [xmtp/xmtpd#366](https://github.com/xmtp/xmtpd/issues/366).

## Scope

**In scope**

- Rate limiting on the `QueryApi` gRPC service: `QueryEnvelopes`, `SubscribeTopics`, `GetInboxIds`, `GetNewestEnvelope`.
- Tier 0 / Tier 2 classification (Tier 1 deferred).
- Redis-backed token bucket with the existing `pkg/ratelimiter/redis_limiter.go`, extended for force-debit.
- Subscription retrospective billing via deferred drain.
- Connection-rate sub-limit on `SubscribeTopics` opens.
- Stream lifetime caps (idle timeout, max duration).
- Fail-open + circuit breaker on Redis runtime failure.

**Out of scope (v1)**

- `NotificationApi`, `PublishApi`, `ReplicationApi`, `GatewayApi` interceptors. (`GatewayApi` already has its own limiter.)
- Tier 1 (fan-out) classification — fan-out services live with Tier 2 limits on `QueryApi` until the `NotificationApi` design lands.
- In-memory limiter fallback. Redis is required when limits are enforced.
- Limit hot reload — config changes require restart.
- Per-method limits — single shared bucket per identity.
- Deprecated `ReplicationApi.SubscribeAllEnvelopes`, `PayerApi.*` paths. These are removed before launch.
- All gateway-side concerns. The gateway already has its own Redis-backed limiter.

## Architecture

### Package Layout

```
pkg/ratelimiter/
  redis_limiter.go         # existing — extended with force-debit support
  classifier.go            # NEW: tier classification from gRPC context
  cost.go                  # NEW: cost formulas (sqrt, drain)
  circuit_breaker.go       # NEW: failure-counting wrapper around the limiter
pkg/interceptors/server/
  rate_limit.go            # NEW: unary + stream interceptors for QueryApi
  rate_limit_test.go       # NEW
pkg/server/server.go       # MODIFIED: register interceptors on QueryApi handler
pkg/config/                # MODIFIED: add rate limit config struct
```

The interceptor is **service-scoped**: it is wrapped only around the `QueryApi` handler when the new d14n services are registered in `pkg/server/server.go`. Other services (`ReplicationApi`, `PublishApi`, `NotificationApi`) and deprecated paths receive no wrapping.

### Data Flow

```
gRPC request
  → auth interceptor (sets tier-related fields in context)
  → rate-limit interceptor (this spec)
      → classifier(ctx) → Tier 0 | Tier 2
      → if Tier 0: pass through
      → if Tier 2: cost(method, request) → debit via circuit-broken limiter
        → if denied: return ResourceExhausted
        → if allowed: call next handler
      → for SubscribeTopics: schedule deferred drain on stream end
  → handler
```

## Tier Classification

### Algorithm

Classify the incoming request based on context state populated by the auth interceptor. **Order matters**: Tier 0 is checked first because a valid node JWT is the only way to bypass limits.

```
classify(ctx) -> Tier
  if ctx has nodeJWT:
    if verifyNodeJWT(nodeJWT) == OK:
      return Tier0
    else:
      return Reject(Unauthenticated)   # all-or-nothing
  return Tier2
```

The classifier returns either a `Tier` value or a `Reject` sentinel. A `Reject` short-circuits the entire interceptor and returns `Unauthenticated` to the caller.

### Tier 2 Identity Key

For Tier 2, the bucket is keyed on **client IP**. The IP is extracted as follows:

1. Read the immediate gRPC peer address from `peer.FromContext(ctx)`.
2. Read the `X-Forwarded-For` metadata header from the incoming gRPC metadata, if present.
3. Walk `X-Forwarded-For` from right to left:
   - If the immediate peer is in the configured trusted-proxy CIDR list, peel one entry from the right of `X-Forwarded-For` and treat that as the new immediate peer.
   - Repeat until the immediate peer is no longer in any trusted CIDR.
4. The final immediate peer is the client IP.

If `X-Forwarded-For` is missing or the immediate peer is not in any trusted CIDR, use the gRPC peer address directly.

For IPv6, normalize the resulting IP to its `/64` prefix. The bucket key is `xmtpd:rl:t2:<normalized-ip>`.

### Trusted Proxy Defaults

The default trusted-proxy list ships with current AWS ALB CIDRs (sourced from AWS published ranges; reviewed when the IP ranges change). Operators may override or extend the list via configuration.

## Cost Model

```
cost_query(num_topics) = ceil(sqrt(max(num_topics, 1)))
cost_subscribe_open(num_topics) = ceil(sqrt(max(num_topics, 1)))
cost_lookup() = 1                              # GetInboxIds, GetNewestEnvelope
cost_subscribe_drain(elapsed_minutes) = ceil(elapsed_minutes / N) * D
```

Where:
- `N` = drain interval in minutes (config: `XMTPD_RATE_LIMIT_DRAIN_INTERVAL_MINUTES`, default `5`)
- `D` = drain amount per interval (config: `XMTPD_RATE_LIMIT_DRAIN_AMOUNT`, default `1`)

The `max(1, ...)` clamp on topic count ensures `cost_query(0) = 1` rather than `0`. A query with zero topics is malformed but should still be charged a baseline cost rather than rejected separately.

A stream that opens and closes within the first interval pays no drain cost (the `ceil` of `elapsed/N` is 0 when `elapsed = 0`). This is intentional: short-lived streams have already paid the `ceil(sqrt(num_topics))` admission cost, and "open and immediately close" abuse is bounded by the separate subscribe-opens-per-minute sub-limit. The drain is for **held** resources; if nothing was held, nothing is owed.

## Bucket Configuration

### Per-Tier Limits

Tier 2 buckets have two concurrent limits enforced in a single Redis call (using the existing multi-limit support in `redis_limiter.go`):

| Limit             | Default        | Config key                                               |
| ----------------- | -------------- | -------------------------------------------------------- |
| Per-minute tokens | 60             | `XMTPD_RATE_LIMIT_T2_PER_MINUTE_CAPACITY`                |
| Per-hour tokens   | 1200           | `XMTPD_RATE_LIMIT_T2_PER_HOUR_CAPACITY`                  |
| Subscribe opens/min | 10           | `XMTPD_RATE_LIMIT_T2_SUBSCRIBE_OPENS_PER_MINUTE`         |

The subscribe-open sub-limit is checked **only** for `SubscribeTopics` and is scoped to the same identity key but with a distinct logical limit name inside the Lua script.

Defaults are starting points; staging traffic will inform the final values before launch.

### Bucket Key Schema

```
xmtpd:rl:t2:<ip-or-prefix>           # main token bucket
```

The Lua script tracks the multiple limits as separate fields within a single Redis hash at this key, so all checks happen in one round trip and the script is single-key (cluster-safe without hash tag concerns).

## Redis Integration

### Required Lua Extension: Force Debit

The existing token-bucket Lua script clamps the result at zero. The retrospective subscription drain requires the bucket to go negative. Add a new operation mode to the script:

- **Mode `check_and_debit`** (existing): atomically check if `current >= cost`; if yes, debit and return allowed; if no, return denied without modification.
- **Mode `force_debit`** (NEW): unconditionally subtract `cost` from current. Allowed to go negative. Always returns allowed. Used for retrospective subscription drain only.

The mode is selected by the limiter call site, not by the client. Only the deferred drain in the subscribe handler issues `force_debit` calls.

### HA Support

The Redis client must support both Sentinel and Cluster modes. Configuration:

| Config key                                | Description                                          |
| ----------------------------------------- | ---------------------------------------------------- |
| `XMTPD_RATE_LIMIT_REDIS_MODE`             | `single` \| `sentinel` \| `cluster` (default: `single`) |
| `XMTPD_RATE_LIMIT_REDIS_ADDRS`            | Comma-separated address list                         |
| `XMTPD_RATE_LIMIT_REDIS_SENTINEL_MASTER`  | Master name (sentinel mode only)                     |
| `XMTPD_RATE_LIMIT_REDIS_PASSWORD`         | Password (optional)                                  |
| `XMTPD_RATE_LIMIT_REDIS_TIMEOUT_MS`       | Per-call timeout (default: `50`)                     |

The Lua script must be cluster-safe: each `EVAL` accesses exactly one key, removing any need for hash tags.

### Startup Behavior

If `XMTPD_RATE_LIMIT_ENABLE=true`, the node attempts a Redis `PING` at startup. On failure, the node logs a fatal error and exits with non-zero status. There is no in-memory fallback. Operators must either run Redis or set `XMTPD_RATE_LIMIT_ENABLE=false`.

### Runtime Failure: Circuit Breaker

The circuit breaker wraps every limiter call. State machine:

```
CLOSED  ──(N consecutive failures)──→  OPEN
OPEN    ──(after cooldown_seconds)──→  HALF_OPEN
HALF_OPEN ──(success)──→ CLOSED
HALF_OPEN ──(failure)──→ OPEN
```

Behavior:

- **CLOSED**: every request hits Redis. Failures increment a counter; success resets it.
- **OPEN**: every request bypasses Redis and is allowed. Counter is not touched. After `cooldown_seconds`, transition to `HALF_OPEN`.
- **HALF_OPEN**: the next request actually hits Redis. Success → `CLOSED`; failure → back to `OPEN` with the cooldown timer reset.

Configuration:

| Config key                                       | Default |
| ------------------------------------------------ | ------- |
| `XMTPD_RATE_LIMIT_BREAKER_FAILURE_THRESHOLD`     | `5`     |
| `XMTPD_RATE_LIMIT_BREAKER_COOLDOWN`              | `10s`   |

A "failure" is any non-nil error from the limiter call (including timeout). Denials are not failures.

## Subscription Lifetime

### Retrospective Drain

In the `SubscribeTopics` handler:

```go
func (s *Service) SubscribeTopics(req *...Request, stream ...) error {
    // 1. Classify and admission-check (interceptor handles this)
    // 2. Record start time
    startedAt := time.Now()
    // 3. Defer drain
    defer func() {
        elapsed := time.Since(startedAt)
        cost := costSubscribeDrain(elapsed)
        // force_debit, allow negative, ignore denial
        s.limiter.ForceDebit(ctx, identityKey, cost)
    }()
    // 4. Run the existing subscribe loop unchanged
    return s.runSubscribeLoop(req, stream)
}
```

The drain runs on every unwind path: graceful close, client disconnect, server-side cancel, panic, context cancellation. Process crashes are the only path that misses it; this is acceptable.

### Termination Layers

| Layer                       | Mechanism                                                   | Config key                                       |
| --------------------------- | ----------------------------------------------------------- | ------------------------------------------------ |
| Transport (gRPC keepalive)  | Existing gRPC server keepalive config                       | (existing — verify settings)                     |
| Application idle timeout    | Cancel stream if no activity for `idleTimeout`              | `XMTPD_RATE_LIMIT_STREAM_IDLE_TIMEOUT` (default: `15m`) |
| Hard maximum duration       | Cancel stream when `time.Since(start) > maxDuration`        | `XMTPD_RATE_LIMIT_STREAM_MAX_DURATION` (default: `60m`) |

The idle timeout and max duration are implemented as additional `select` cases in the existing subscribe loop's main goroutine — no new goroutines, no new registries.

## Errors

| Condition                                     | gRPC Code           | Notes                                                |
| --------------------------------------------- | ------------------- | ---------------------------------------------------- |
| Bucket exhausted                              | `ResourceExhausted` | Includes `retry-after-seconds` in trailer metadata   |
| Subscribe opens/min exceeded                  | `ResourceExhausted` | Same trailer                                         |
| Tier 0 JWT verification failed                | `Unauthenticated`   | No fall-through                                      |
| Stream cancelled by idle timeout              | `DeadlineExceeded`  | Server-initiated                                     |
| Stream cancelled by max-duration              | `DeadlineExceeded`  | Server-initiated; client should reconnect            |
| Redis circuit-breaker OPEN                    | (request allowed)   | Counter increments; alert fires                      |

The `retry-after-seconds` trailer is computed from the bucket's refill rate and current deficit. Conservative estimate is fine — clients should treat it as a lower bound.

## Observability

### Metrics (Prometheus)

| Metric                                               | Type     | Labels                                  |
| ---------------------------------------------------- | -------- | --------------------------------------- |
| `xmtpd_rate_limit_decisions_total`                   | Counter  | `service`, `method`, `tier`, `outcome`  |
| `xmtpd_rate_limit_redis_call_duration_seconds`       | Histogram | `outcome` (ok / error / timeout)        |
| `xmtpd_rate_limit_circuit_breaker_state`             | Gauge    | (0=closed, 1=half_open, 2=open)         |
| `xmtpd_rate_limit_circuit_breaker_trips_total`       | Counter  | -                                       |
| `xmtpd_rate_limit_subscribe_drain_total`             | Counter  | `outcome` (drained / breaker_skipped)   |
| `xmtpd_rate_limit_classifier_failures_total`         | Counter  | `reason` (jwt_expired / jwt_invalid / signer_unknown) |
| `xmtpd_rate_limit_stream_terminations_total`         | Counter  | `reason` (idle / max_duration / client_close / server_close) |

`outcome` for the decisions counter: `allowed` | `denied` | `bypassed` (Tier 0) | `failed_open` (circuit breaker).

### Logs

Lowercase, structured (zap), guarded for hot paths. Logger name: `"xmtpd.rate-limiter"`.

- **Info, once at startup**: `"rate limit interceptor enabled"` with config snapshot.
- **Warn, on circuit breaker open/close**: `"rate limit circuit breaker tripped"` / `"recovered"`.
- **Debug, on every decision (guarded)**: `"rate limit decision"` with `tier`, `method`, `cost`, `outcome`.
- **Warn, on Tier 0 verification failure**: `"tier 0 jwt verification failed"` with `reason`.

## Configuration Reference

All config is read from environment variables (and matching CLI flags via `go-flags` struct tags). Master switch:

| Key                              | Default  | Effect                                                 |
| -------------------------------- | -------- | ------------------------------------------------------ |
| `XMTPD_RATE_LIMIT_ENABLE`       | `false`  | When `false`, the interceptor is not registered at all |

When `false`, the entire rate-limiting subsystem is bypassed and Redis is not required. This is the dev/local default. Production deployments set this to `true`.

(Other config keys listed above in Bucket Configuration, Redis Integration, and Subscription Lifetime sections.)

## Test Plan

### Unit Tests

- **`classifier_test.go`**
  - Tier 0 with valid JWT → Tier 0
  - Tier 0 with expired JWT → reject
  - Tier 0 with invalid signature → reject
  - Tier 0 with valid signature, signer not in registry → reject
  - No JWT → Tier 2
  - Tier 2 IP extraction: direct peer
  - Tier 2 IP extraction: trusted proxy with `X-Forwarded-For`
  - Tier 2 IP extraction: untrusted proxy ignores `X-Forwarded-For`
  - Tier 2 IP extraction: IPv6 normalized to /64

- **`cost_test.go`**
  - `cost_query(0) = 1`, `cost_query(1) = 1`, `cost_query(100) = 10`, `cost_query(1000) = 32`
  - `cost_subscribe_drain(0 min)` → 1 interval charged
  - `cost_subscribe_drain(N min)` → 1 interval
  - `cost_subscribe_drain(N+1 min)` → 2 intervals

- **`circuit_breaker_test.go`**
  - Closed → open after threshold failures
  - Open → half_open after cooldown
  - Half_open success → closed
  - Half_open failure → open with reset cooldown
  - Open state allows requests through without calling Redis

- **`rate_limit_interceptor_test.go`** (with fake limiter)
  - Tier 0 bypasses limiter
  - Tier 2 calls limiter with correct cost
  - Limiter denial → `ResourceExhausted`
  - Limiter error → fail open
  - Subscribe registers deferred drain
  - Subscribe drain fires on context cancel
  - Subscribe drain fires on graceful close

### Integration Tests

Use a real Redis instance via testcontainers (`pkg/testutils/`).

- **End-to-end Tier 2 deny**: hit `QueryEnvelopes` from a fake client until `ResourceExhausted` is returned.
- **End-to-end retrospective drain**: open a subscribe, hold it for >N minutes, close it, verify the bucket reflects the drain.
- **End-to-end subscribe-opens/min**: rapidly open and close subscriptions until the connection-rate sub-limit fires.
- **End-to-end Tier 0 bypass**: hammer `QueryEnvelopes` with a node JWT and verify no rate limit fires.
- **Redis failover**: kill the Redis container mid-traffic; verify the circuit breaker opens, requests continue serving, and recovery works when Redis comes back.
- **Stream lifetime**: open a subscribe and verify it is cancelled at `idleTimeout` and at `maxDuration`.

### Verification Items (before merge)

1. Verify the Lua script's force-debit mode works correctly with `redis-cli EVAL` against a real Redis.
2. Verify the Lua script is single-key (cluster-safe) by running it against a `redis-cluster` testcontainer.
3. Verify the existing gRPC keepalive settings on the node server are reasonable for stream lifetime layer 1.
4. Verify the auth interceptor populates the context fields the classifier needs (or modify it to do so).
5. Verify the subscribe handler structure can accommodate the deferred drain and idle/max-duration `select` cases without restructuring.
6. Verify AWS ALB CIDR list is current at merge time.

## Open Questions

None for v1 — all design questions resolved during brainstorming. Defaults (limits, intervals, durations) will be tuned in staging before launch but starting values are specified above.

## Out of Scope (Restated)

- `NotificationApi` rate limiting — explicitly deferred to a future design.
- `PublishApi` rate limiting — handled by WAF.
- Tier 1 (fan-out) classification — deferred with `NotificationApi`.
- In-memory fallback limiter — by design, never implemented.
- Hot reload of limit values — restart required.
- Per-method bucket scoping — single bucket per identity.
- Deprecated d14n endpoints — being removed before launch.

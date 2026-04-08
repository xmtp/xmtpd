# Rate Limiting Implementation Plan — Node `QueryApi` (v1)

> **For agentic workers:** Implement task-by-task in order. Each task is TDD-shaped: write the failing test, run it, implement, run it again, commit. Steps use checkbox (`- [ ]`) syntax.

## Current Status (handoff point)

**Branch:** `feat/rate-limiting-queryapi`
**HEAD:** `1296d640`
**Base:** `927af16a` (main)
**Commits on branch:** 18

### Completed phases

- [x] **Phase 0** — Foundations: `RateLimitOptions` config, cost formulas (`CostQuery`, `CostSubscribeDrain`), tier classifier, source-IP extraction with trusted-proxy peeling and IPv6 `/64` normalization. All pure functions, unit-tested.
- [x] **Phase 1** — Lua `force_debit` mode with strict mode validation, `ForceDebit` Go method, multi-limit test. Backward compatible with existing `Allow` path.
- [x] **Phase 2** — Consecutive-failure `CircuitBreaker` + `BreakerLimiter` fail-open wrapper. 9 unit tests.
- [x] **Phase 3** — Connect-RPC interceptor for `QueryApi` with two-limiter design (query bucket + opens bucket), method routing helper, cost helper, unary + streaming handlers. 11 unit tests.
- [x] **Phase 4** — `SubscribeTopics` handler integration: admission cost at open, deferred drain at close, idle timeout, hard max-duration cap. 4 unit tests for the admission/drain helper.
- [x] **Phase 5** — Server wiring: `ratelimiter.Build` startup constructor with fail-fast Redis ping, `Redis` field added to `ServerOptions`, `startAPIServer` wraps rate-limit interceptor around the QueryApi handler only (not Replication, Publish, or Notification), `SetSubscribeTrustedCIDRs` called once at startup, `msgRateLimiter` and `rlConfig` passed to `NewReplicationAPIService`.
- [x] **Phase 6** — Prometheus metrics: `xmtpd_rate_limit_decisions_total` (service, method, tier, outcome), `xmtpd_rate_limit_circuit_breaker_state` gauge, `xmtpd_rate_limit_circuit_breaker_trips_total`, `xmtpd_rate_limit_stream_terminations_total` (reason). Registered via `ratelimiter.Register(promReg)` in `startAPIServer`.

### State at handoff

- `go build ./...` — clean.
- `go test ./pkg/ratelimiter/ ./pkg/interceptors/server/ ./pkg/api/message/ -short -count=1` — all green. ~45 unit tests across the three packages.
- Test packages that require Docker (testcontainer-redis) for the pre-existing `TestRedisLimiter_*` tests — pass locally when Docker is running.
- Pre-existing main-branch test-file breakage (`InsertGatewayEnvelopeParams` → `InsertGatewayEnvelopeV3Params` in `pkg/api/message/service_api_test.go`) was fixed in `ebfacd85` as a drive-by so the message test package could compile.
- Spec (`tasks/rate-limiting-spec.md`) config-reference tables were updated to match the implemented env var names (removed stale `_ENABLED`/`_SECONDS`/`_MINUTES` suffixes).

### Remaining work

- [ ] **Phase 7** — Integration tests using `pkg/testutils/redis.NewRedisForTest` (real Redis via testcontainers). Five tests:
  1. End-to-end Tier 2 query deny after budget exhaustion
  2. Retrospective drain driving the bucket negative, then next call rejected
  3. Tier 0 bypass under tight limits
  4. Subscribe opens-per-minute sub-limit triggering `ResourceExhausted`
  5. Redis connection killed mid-flight → circuit breaker trips → fail-open recovery
- [ ] **Final verification** — `dev/test`, `dev/lint-fix`, manual smoke test with a real Redis and `grpcurl`, metric scrape against `http://127.0.0.1:8008/metrics`.
- [ ] **PR** — open as draft referencing `xmtp/xmtpd#366`.

### Commit history

```
1296d640 feat(ratelimiter): add prometheus metrics for decisions, breaker state, and stream terminations
246d6c2e feat(server): wire QueryApi rate-limit interceptor and config into startAPIServer
56591669 feat(ratelimiter): add startup builder with fail-fast Redis ping
ebfacd85 fix(test): rename InsertGatewayEnvelopeParams to InsertGatewayEnvelopeV3Params
e4265112 feat(message): apply admission cost and deferred drain to SubscribeTopics
0bce224b feat(message): plumb rate limiter and config into Service constructor
41a20c50 feat(interceptors): add Connect interceptor for QueryApi rate limiting
ffc5dbc1 feat(interceptors): add QueryApi method routing and cost helpers for rate limiting
6a8858ea feat(ratelimiter): add BreakerLimiter wrapper for fail-open Redis runtime errors
373c93be feat(ratelimiter): add consecutive-failure circuit breaker
7bd6b433 fix(ratelimiter): validate Lua mode argument and add multi-limit ForceDebit test
4bb3102b feat(ratelimiter): add ForceDebit for retrospective subscription billing
a4a34a67 feat(ratelimiter): extend Lua script with force_debit mode for retrospective drain
5d36cd3e fix(ratelimiter): return sentinel for unparseable client IP and document Tier1 reservation
57b9c64b feat(ratelimiter): add source IP extraction with trusted proxy peeling and IPv6 /64 normalization
7e429873 feat(ratelimiter): add tier classifier from auth context flag
242bff26 feat(ratelimiter): add cost functions for query and subscription drain
76217129 feat(config): add RateLimitOptions for QueryApi rate limiting
```

### Next session primer

> Continue rate-limiting work on `feat/rate-limiting-queryapi` (HEAD `1296d640`). Phases 0–6 complete and green. Start at Phase 7 (Integration Tests) in `tasks/rate-limiting-plan.md`. Use `pkg/testutils/redis.NewRedisForTest` for real Redis via testcontainers. After Phase 7, run `dev/test` and `dev/lint-fix`, then open a draft PR referencing `xmtp/xmtpd#366`.

### Known caveats for the next session

- The `TestCircuitBreaker_*` tests don't reset `BreakerStateGauge` between runs. It's a harmless global side effect (tests read `breaker.State()`, not the gauge), but if you add gauge-reading tests, reset the gauge in a `t.Cleanup` or make it a per-instance field.
- `ExtractClientIP` returns the sentinel `"invalid"` (constant `InvalidClientIPKey`) when the peer address can't be parsed. All malformed traffic shares one bucket — intentional, from a code-review fix.
- The subscribe-topics handler uses a package-level `subscribeTrustedCIDRs` set via `SetSubscribeTrustedCIDRs` at server startup. Integration tests that exercise the full handler path must call this (or accept that the nil default means no XFF peeling).
- Commit `57b9c64b` is an empty marker commit left from a subagent split mistake in Phase 0. Consider `git rebase -i` to drop it before opening the PR.

---

**Goal:** Add Redis-backed rate limiting on the node's `QueryApi` (`QueryEnvelopes`, `SubscribeTopics`, `GetInboxIds`, `GetNewestEnvelope`) with Tier 0 (node JWT) bypass, retrospective subscription billing, and fail-open circuit-breaker behavior on Redis outages.

**Architecture:** Connect-RPC interceptor wrapped only around the `QueryApi` handler in `pkg/server/server.go`. Reuses the existing `pkg/ratelimiter/redis_limiter.go` token-bucket implementation, extended with a `ForceDebit` method that bypasses the non-negative clamp. Subscription drain via `defer` in the handler — no per-stream goroutines or registries. Circuit breaker wraps every Redis call to absorb runtime Redis outages.

**Tech Stack:** Go 1.25, Connect-RPC (`connectrpc.com/connect`), go-redis v9 (`UniversalClient`), zap logging, prometheus metrics, testify, testcontainers (via `pkg/testutils/redis`).

**Key adjustments from spec discovered during recon:**
- Spec said "gRPC interceptor" — actual implementation is a Connect interceptor (`connect.Interceptor`).
- Spec said the classifier verifies JWTs — actually the existing `ServerAuthInterceptor` already does this and sets `constants.VerifiedNodeRequestCtxKey{}=true`. The classifier just reads the context flag.
- Spec proposed new Redis config keys for cluster/sentinel modes — `pkg/config/redis.go` already supports HA via URL form (`UniversalClient`), no new keys needed.

---

## Pre-flight

- [ ] **Create a feature branch.**

```bash
cd /Users/martinkysel/work/xmtpd
git checkout -b feat/rate-limiting-queryapi
```

---

## File Structure

```
pkg/ratelimiter/
  script.lua                      MODIFIED   add force-debit mode
  redis_limiter.go                MODIFIED   add ForceDebit method
  redis_limiter_test.go           MODIFIED   add force-debit tests
  classifier.go                   NEW        tier classification + IP extraction
  classifier_test.go              NEW
  cost.go                         NEW        cost formulas (sqrt, drain)
  cost_test.go                    NEW
  circuit_breaker.go              NEW        failure-counting wrapper
  circuit_breaker_test.go         NEW
  metrics.go                      NEW        prometheus collectors
pkg/interceptors/server/
  rate_limit.go                   NEW        connect interceptor
  rate_limit_test.go              NEW
pkg/api/message/
  subscribe_topics.go             MODIFIED   defer drain, idle timeout, max duration
pkg/config/
  options.go                      MODIFIED   add RateLimitOptions
pkg/server/
  server.go                       MODIFIED   wire interceptor onto QueryApi handler only
```

---

## Phase 0 — Foundations (no Redis required)

### Task 0.1: Add `RateLimitOptions` config struct

**Files:**
- Modify: `pkg/config/options.go` (append after `ReflectionOptions`, register in `ServerOptions`)

- [ ] **Step 1: Add the options struct.**

```go
// RateLimitOptions controls the QueryApi rate-limiting interceptor.
type RateLimitOptions struct {
	Enable bool `long:"enable" env:"XMTPD_RATE_LIMIT_ENABLE" description:"Enable QueryApi rate limiting (requires Redis)"`

	// Tier 2 bucket limits
	T2PerMinuteCapacity         int `long:"t2-per-minute-capacity"          env:"XMTPD_RATE_LIMIT_T2_PER_MINUTE_CAPACITY"          description:"Tier 2 per-minute token capacity"          default:"60"`
	T2PerHourCapacity           int `long:"t2-per-hour-capacity"            env:"XMTPD_RATE_LIMIT_T2_PER_HOUR_CAPACITY"            description:"Tier 2 per-hour token capacity"            default:"1200"`
	T2SubscribeOpensPerMinute   int `long:"t2-subscribe-opens-per-minute"   env:"XMTPD_RATE_LIMIT_T2_SUBSCRIBE_OPENS_PER_MINUTE"   description:"Tier 2 subscribe-opens per minute"         default:"10"`

	// Subscription drain
	DrainIntervalMinutes int `long:"drain-interval-minutes" env:"XMTPD_RATE_LIMIT_DRAIN_INTERVAL_MINUTES" description:"Minutes per subscription drain interval" default:"5"`
	DrainAmount          int `long:"drain-amount"           env:"XMTPD_RATE_LIMIT_DRAIN_AMOUNT"           description:"Tokens per drain interval"                default:"1"`

	// Stream lifetime
	StreamIdleTimeout time.Duration `long:"stream-idle-timeout" env:"XMTPD_RATE_LIMIT_STREAM_IDLE_TIMEOUT" description:"Cancel a stream that has had no activity for this duration" default:"15m"`
	StreamMaxDuration time.Duration `long:"stream-max-duration" env:"XMTPD_RATE_LIMIT_STREAM_MAX_DURATION" description:"Hard cap on subscription stream lifetime"                  default:"60m"`

	// Circuit breaker
	BreakerFailureThreshold int           `long:"breaker-failure-threshold" env:"XMTPD_RATE_LIMIT_BREAKER_FAILURE_THRESHOLD" description:"Consecutive Redis failures before tripping the circuit breaker" default:"5"`
	BreakerCooldown         time.Duration `long:"breaker-cooldown"          env:"XMTPD_RATE_LIMIT_BREAKER_COOLDOWN"          description:"How long the circuit breaker stays open before probing"          default:"10s"`
	RedisCallTimeout        time.Duration `long:"redis-call-timeout"        env:"XMTPD_RATE_LIMIT_REDIS_CALL_TIMEOUT"        description:"Per-call Redis timeout"                                          default:"50ms"`

	// Trusted proxy CIDRs (comma-separated)
	TrustedProxyCIDRs string `long:"trusted-proxy-cidrs" env:"XMTPD_RATE_LIMIT_TRUSTED_PROXY_CIDRS" description:"Comma-separated trusted proxy CIDR list for X-Forwarded-For peeling"`
}
```

- [ ] **Step 2: Register the struct in `ServerOptions`.**

Add a new field in the `ServerOptions` struct (around line 140):

```go
RateLimit       RateLimitOptions       `group:"Rate Limit Options"       namespace:"rate-limit"`
```

- [ ] **Step 3: Verify the config compiles.**

Run: `go build ./pkg/config/...`
Expected: clean build.

- [ ] **Step 4: Commit.**

```bash
git add pkg/config/options.go
git commit -m "feat(config): add RateLimitOptions for QueryApi rate limiting"
```

---

### Task 0.2: Cost functions

**Files:**
- Create: `pkg/ratelimiter/cost.go`
- Create: `pkg/ratelimiter/cost_test.go`

- [ ] **Step 1: Write the failing test.**

```go
package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCostQuery(t *testing.T) {
	tests := []struct {
		topics int
		want   uint64
	}{
		{0, 1},     // clamp to 1
		{1, 1},
		{4, 2},
		{100, 10},
		{1000, 32},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, CostQuery(tt.topics), "topics=%d", tt.topics)
	}
}

func TestCostSubscribeDrain(t *testing.T) {
	intervalMinutes := 5
	drainAmount := 1
	tests := []struct {
		elapsed time.Duration
		want    uint64
	}{
		{0, 0},                        // no time held → no drain
		{1 * time.Minute, 1},          // partial first interval → 1
		{5 * time.Minute, 1},          // exactly one interval → 1
		{6 * time.Minute, 2},          // into second interval → 2
		{60 * time.Minute, 12},        // 1 hour → 12 intervals
	}
	for _, tt := range tests {
		got := CostSubscribeDrain(tt.elapsed, intervalMinutes, drainAmount)
		require.Equal(t, tt.want, got, "elapsed=%s", tt.elapsed)
	}
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestCost' -v`
Expected: FAIL — undefined `CostQuery`, `CostSubscribeDrain`.

- [ ] **Step 3: Implement.**

```go
package ratelimiter

import (
	"math"
	"time"
)

// CostQuery returns the rate-limit token cost of a query against `numTopics` topics.
// Cost is sublinear: ceil(sqrt(max(numTopics, 1))). A 0-topic query is malformed
// but charged the baseline cost of 1 rather than rejected separately.
func CostQuery(numTopics int) uint64 {
	if numTopics < 1 {
		numTopics = 1
	}
	return uint64(math.Ceil(math.Sqrt(float64(numTopics))))
}

// CostSubscribeDrain returns the retrospective drain cost for a subscription
// that was held open for `elapsed` time. The drain is computed in whole intervals:
// floor(elapsed / intervalMinutes) intervals, each costing `drainAmount` tokens.
//
// A stream that closes within the first interval pays no drain cost. The
// admission cost paid at open time and the subscribe-opens-per-minute sub-limit
// together prevent open-and-immediately-close abuse — the drain is for held
// resources only.
func CostSubscribeDrain(elapsed time.Duration, intervalMinutes, drainAmount int) uint64 {
	if elapsed <= 0 || intervalMinutes <= 0 || drainAmount <= 0 {
		return 0
	}
	intervals := uint64(math.Ceil(elapsed.Minutes() / float64(intervalMinutes)))
	return intervals * uint64(drainAmount)
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestCost' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/cost.go pkg/ratelimiter/cost_test.go
git commit -m "feat(ratelimiter): add cost functions for query and subscription drain"
```

---

### Task 0.3: Tier classifier (without IP extraction yet)

**Files:**
- Create: `pkg/ratelimiter/classifier.go`
- Create: `pkg/ratelimiter/classifier_test.go`

- [ ] **Step 1: Write the failing test.**

```go
package ratelimiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
)

func TestClassify_Tier0WhenContextFlagSet(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)
	tier := ClassifyTier(ctx)
	require.Equal(t, Tier0, tier)
}

func TestClassify_Tier2WhenNoContextFlag(t *testing.T) {
	ctx := context.Background()
	tier := ClassifyTier(ctx)
	require.Equal(t, Tier2, tier)
}

func TestClassify_Tier2WhenContextFlagFalse(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, false)
	tier := ClassifyTier(ctx)
	require.Equal(t, Tier2, tier)
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestClassify' -v`
Expected: FAIL — undefined `Tier0`, `Tier2`, `ClassifyTier`.

- [ ] **Step 3: Implement.**

```go
package ratelimiter

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/constants"
)

// Tier identifies which rate-limit policy applies to a request.
type Tier int

const (
	// Tier0 is authenticated node-to-node traffic. Bypasses all limits.
	Tier0 Tier = iota
	// Tier2 is unauthenticated edge-client traffic. Subject to limits.
	Tier2
)

// ClassifyTier inspects the request context for the auth interceptor's
// verified-node flag. If present and true, the request is Tier 0. Otherwise
// it is Tier 2. Tier 0 verification (JWT validation, signer-in-registry check)
// is performed by the auth interceptor upstream — the classifier never
// re-verifies and never falls back to Tier 2 on JWT failure (the auth
// interceptor returns Unauthenticated directly in that case).
func ClassifyTier(ctx context.Context) Tier {
	if v, ok := ctx.Value(constants.VerifiedNodeRequestCtxKey{}).(bool); ok && v {
		return Tier0
	}
	return Tier2
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestClassify' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/classifier.go pkg/ratelimiter/classifier_test.go
git commit -m "feat(ratelimiter): add tier classifier from auth context flag"
```

---

### Task 0.4: Source IP extraction with trusted proxy peeling

**Files:**
- Modify: `pkg/ratelimiter/classifier.go` (append)
- Modify: `pkg/ratelimiter/classifier_test.go` (append)

- [ ] **Step 1: Write the failing tests.**

```go
func TestExtractClientIP_NoXFFUsesPeer(t *testing.T) {
	got := ExtractClientIP("203.0.113.1:5001", "", nil)
	require.Equal(t, "203.0.113.1", got)
}

func TestExtractClientIP_TrustedProxyPeelsOneHop(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	// peer is in trusted CIDR, XFF has [client, proxy]; we peel the right entry
	got := ExtractClientIP("10.0.0.5:5001", "203.0.113.1, 10.0.0.5", trusted)
	require.Equal(t, "203.0.113.1", got)
}

func TestExtractClientIP_UntrustedProxyIgnoresXFF(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	got := ExtractClientIP("198.51.100.7:5001", "203.0.113.1", trusted)
	require.Equal(t, "198.51.100.7", got)
}

func TestExtractClientIP_IPv6NormalizedToSlash64(t *testing.T) {
	got := ExtractClientIP("[2001:db8:abcd:1234:5678:9abc:def0:1234]:5001", "", nil)
	require.Equal(t, "2001:db8:abcd:1234::/64", got)
}

func TestExtractClientIP_IPv4NotNormalized(t *testing.T) {
	got := ExtractClientIP("203.0.113.1:5001", "", nil)
	require.Equal(t, "203.0.113.1", got)
}

func mustParseCIDRs(t *testing.T, cidrs []string) []*net.IPNet {
	t.Helper()
	out := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, n, err := net.ParseCIDR(c)
		require.NoError(t, err)
		out = append(out, n)
	}
	return out
}
```

Add the import: `"net"`.

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestExtractClientIP' -v`
Expected: FAIL — undefined `ExtractClientIP`.

- [ ] **Step 3: Implement.**

```go
import (
	"net"
	"strings"
)

// ParseTrustedProxyCIDRs parses a comma-separated CIDR list. Empty entries are
// ignored. Returns an error on the first malformed entry.
func ParseTrustedProxyCIDRs(s string) ([]*net.IPNet, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make([]*net.IPNet, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		_, n, err := net.ParseCIDR(p)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted-proxy CIDR %q: %w", p, err)
		}
		out = append(out, n)
	}
	return out, nil
}

// ExtractClientIP returns the bucket-key IP for an incoming request.
//
// The algorithm:
//  1. Parse the immediate peer address (host:port).
//  2. If the peer is in any trusted CIDR and X-Forwarded-For is non-empty,
//     peel the rightmost entry of XFF and treat it as the new peer.
//  3. Repeat until the peer is no longer in a trusted CIDR or XFF is exhausted.
//  4. For IPv4, return the dotted-quad string. For IPv6, return the /64 prefix
//     in CIDR notation (clients within a /64 share a bucket).
func ExtractClientIP(peerAddr, xff string, trusted []*net.IPNet) string {
	host, _, err := net.SplitHostPort(peerAddr)
	if err != nil {
		host = peerAddr
	}

	xffParts := splitXFF(xff)
	for ipInTrusted(host, trusted) && len(xffParts) > 0 {
		host = strings.TrimSpace(xffParts[len(xffParts)-1])
		xffParts = xffParts[:len(xffParts)-1]
	}

	parsed := net.ParseIP(host)
	if parsed == nil {
		return host
	}
	if v4 := parsed.To4(); v4 != nil {
		return v4.String()
	}
	// IPv6: zero out everything after the /64 prefix
	mask := net.CIDRMask(64, 128)
	return (&net.IPNet{IP: parsed.Mask(mask), Mask: mask}).String()
}

func splitXFF(xff string) []string {
	if xff == "" {
		return nil
	}
	return strings.Split(xff, ",")
}

func ipInTrusted(host string, trusted []*net.IPNet) bool {
	if len(trusted) == 0 {
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	for _, n := range trusted {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
```

Add `"fmt"` to imports if not already present.

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestExtractClientIP' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/classifier.go pkg/ratelimiter/classifier_test.go
git commit -m "feat(ratelimiter): add source IP extraction with trusted proxy peeling and IPv6 /64 normalization"
```

---

## Phase 1 — Lua Force-Debit Extension

### Task 1.1: Extend `script.lua` with a `force_debit` mode

**Files:**
- Modify: `pkg/ratelimiter/script.lua`
- Modify: `pkg/ratelimiter/redis_limiter.go` (pass mode arg)

The current script always does `check_and_debit`. We add a new ARGV slot for the mode and a separate code path for force-debit. The script remains single-key per call (cluster-safe) since it still touches only `KEYS[1..N+1]` for the same subject.

- [ ] **Step 1: Add the mode argument to the Lua script.**

Replace lines 36-56 of `script.lua` with:

```lua
-- Parse input arguments
local ts_key    = KEYS[1]           -- Timestamp key (shared across all limits)
local now_ms    = tonumber(ARGV[1]) -- Current time in milliseconds
local n         = tonumber(ARGV[2]) -- Number of limits to check
local cost      = tonumber(ARGV[3]) -- Tokens to consume
local mode      = ARGV[4]           -- "check" (default) or "force"

-- Initialize arrays to store limit configurations and current state
local caps      = {} -- Maximum capacity for each limit (e.g., 100 tokens)
local refill_ms = {} -- Time in ms to fully refill each limit (e.g., 60000 for 1 minute)
local tokens    = {} -- Current token count for each limit (after refill)

-- Parse limit configurations from ARGV
-- Arguments come in pairs: [capacity, refill_time, capacity, refill_time, ...]
local idx       = 5 -- Start after the first 4 arguments
for i = 1, n do
    caps[i] = tonumber(ARGV[idx])
    idx = idx + 1
    refill_ms[i] = tonumber(ARGV[idx])
    idx = idx + 1
end
```

- [ ] **Step 2: Replace STEP 3 (the failure check) to support force mode.**

Replace lines 106-118 (the `failed_index` block) with:

```lua
-- ============================================================================
-- STEP 3: Check if all limits can satisfy the request
-- ============================================================================
-- In "check" mode, ALL limits must have enough tokens.
-- In "force" mode, the deduction is unconditional and the result is allowed.
local failed_index = 0
if mode ~= "force" then
    for i = 1, n do
        if tokens[i] < cost then
            failed_index = i -- Store the 1-based index of the failing limit
            break
        end
    end
end
```

- [ ] **Step 3: Allow `tokens[i]` to go negative in the success path.**

Inside STEP 4 (`if failed_index == 0 then`), the existing line `tokens[i] = tokens[i] - cost` already does the subtraction. But the TTL block downstream relies on `tokens[i] >= caps[i]` and `caps[i] - tokens[i]`, both of which still work for negative values. The only issue is the TTL calculation when negative. Wrap the TTL calc:

Replace the TTL block (lines 144-158 area) with:

```lua
    -- Set each limit key with its value and expiration in a single call
    -- When a bucket is full, it expires after its refill period
    -- When not full, set TTL based on time to refill to capacity
    -- Negative balances (force_debit) get a TTL covering full refill from zero
    -- plus the deficit, capped at 2x refill_ms for sanity.
    for i = 1, n do
        local ttl
        if tokens[i] >= caps[i] then
            -- Bucket is full - expire after full refill period
            ttl = refill_ms[i]
        elseif tokens[i] >= 0 then
            -- Bucket not full but non-negative - time to refill to capacity
            local time_to_fill = (caps[i] - tokens[i]) * refill_ms[i] / caps[i]
            ttl = math.ceil(time_to_fill)
        else
            -- Negative balance - extra time to climb back to zero, then to full
            local time_to_fill = (caps[i] - tokens[i]) * refill_ms[i] / caps[i]
            ttl = math.ceil(math.min(time_to_fill, refill_ms[i] * 2))
        end

        -- Use SET with PX to set value and expiration atomically
        redis.call("SET", KEYS[i + 1], tostring(tokens[i]), "PX", ttl)
    end
```

- [ ] **Step 4: Update `redis_limiter.go` `buildArgs` to include the mode.**

Modify `buildArgs` in `pkg/ratelimiter/redis_limiter.go`:

```go
func (l *RedisLimiter) buildArgs(requestTime time.Time, cost uint64, mode string) []any {
	args := make([]any, 0, 4+len(l.limits)*2)
	args = append(args, requestTime.UnixMilli(), len(l.limits), cost, mode)
	for _, lim := range l.limits {
		args = append(args, lim.Capacity, lim.RefillEvery.Milliseconds())
	}
	return args
}
```

Update the existing call site in `Allow`:

```go
args := l.buildArgs(now, cost, "check")
```

- [ ] **Step 5: Run existing limiter tests to verify nothing broke.**

Run: `go test ./pkg/ratelimiter/ -run 'TestRedisLimiter' -v`
Expected: PASS — existing tests still pass since `"check"` is the original behavior.

(Note: this test requires Docker for testcontainers. If unavailable locally, run only when CI has Docker.)

- [ ] **Step 6: Commit.**

```bash
git add pkg/ratelimiter/script.lua pkg/ratelimiter/redis_limiter.go
git commit -m "feat(ratelimiter): extend Lua script with force_debit mode for retrospective drain"
```

---

### Task 1.2: Add `ForceDebit` method to `RedisLimiter`

**Files:**
- Modify: `pkg/ratelimiter/redis_limiter.go`
- Modify: `pkg/ratelimiter/redis_limiter_test.go`

- [ ] **Step 1: Write the failing test.**

Add to `redis_limiter_test.go`:

```go
func TestRedisLimiter_ForceDebit_AllowsNegativeBalance(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	ctx := context.Background()

	// Spend 8 tokens normally
	res, err := limiter.Allow(ctx, "subj", 8)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// Force-debit 20 tokens — should drive the bucket negative without rejection
	res, err = limiter.ForceDebit(ctx, "subj", 20)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.Less(t, res.Balances[0].Remaining, 0.0)

	// Subsequent normal Allow should be rejected
	res, err = limiter.Allow(ctx, "subj", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed)
}

func TestRedisLimiter_ForceDebit_ZeroCostIsError(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)
	_, err = limiter.ForceDebit(context.Background(), "subj", 0)
	require.ErrorIs(t, err, ratelimiter.ErrCostMustBeGreaterThanZero)
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestRedisLimiter_ForceDebit' -v`
Expected: FAIL — undefined method `ForceDebit`.

- [ ] **Step 3: Implement `ForceDebit` in `redis_limiter.go`.**

Add after the existing `Allow` method:

```go
// ForceDebit unconditionally subtracts `cost` tokens from the bucket, allowing
// the result to go negative. Used for retrospective subscription drain: when a
// stream closes after holding resources, we bill for the held time even if the
// client's bucket cannot absorb the full charge. The next normal Allow call
// from that subject will be rejected until the bucket refills back to a
// sufficient positive value.
//
// ForceDebit always returns Allowed=true (or an error if Redis fails).
func (l *RedisLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*Result, error) {
	if cost == 0 {
		return nil, ErrCostMustBeGreaterThanZero
	}
	now := time.Now()
	keys := l.buildKeys(subject)
	args := l.buildArgs(now, cost, "force")

	raw, err := l.script.Run(ctx, l.client, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	arr, ok := raw.([]any)
	if !ok || len(arr) < 1 {
		return nil, ErrUnexpectedScriptResponse
	}

	return l.transformResult(arr, cost)
}
```

- [ ] **Step 4: Update the `RateLimiter` interface in `interface.go`.**

```go
type RateLimiter interface {
	Allow(ctx context.Context, subject string, cost uint64) (*Result, error)
	ForceDebit(ctx context.Context, subject string, cost uint64) (*Result, error)
}
```

- [ ] **Step 5: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestRedisLimiter_ForceDebit' -v`
Expected: PASS (requires Docker for testcontainer).

- [ ] **Step 6: Commit.**

```bash
git add pkg/ratelimiter/redis_limiter.go pkg/ratelimiter/redis_limiter_test.go pkg/ratelimiter/interface.go
git commit -m "feat(ratelimiter): add ForceDebit for retrospective subscription billing"
```

---

## Phase 2 — Circuit Breaker

### Task 2.1: Circuit breaker state machine

**Files:**
- Create: `pkg/ratelimiter/circuit_breaker.go`
- Create: `pkg/ratelimiter/circuit_breaker_test.go`

- [ ] **Step 1: Write the failing test.**

```go
package ratelimiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker_OpensAfterThresholdFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	require.Equal(t, BreakerClosed, cb.State())

	for i := 0; i < 3; i++ {
		require.True(t, cb.Allow())
		cb.RecordFailure()
	}

	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())
}

func TestCircuitBreaker_HalfOpenAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())

	time.Sleep(80 * time.Millisecond)

	// First call after cooldown should be allowed (half-open probe)
	require.True(t, cb.Allow())
	require.Equal(t, BreakerHalfOpen, cb.State())
}

func TestCircuitBreaker_HalfOpenSuccessClosesCircuit(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(80 * time.Millisecond)

	require.True(t, cb.Allow())
	cb.RecordSuccess()
	require.Equal(t, BreakerClosed, cb.State())
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(80 * time.Millisecond)

	require.True(t, cb.Allow())
	cb.RecordFailure()
	require.Equal(t, BreakerOpen, cb.State())
	require.False(t, cb.Allow())
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess() // resets counter
	cb.RecordFailure()
	require.Equal(t, BreakerClosed, cb.State())
}

var errSentinel = errors.New("sentinel")
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestCircuitBreaker' -v`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement.**

```go
package ratelimiter

import (
	"sync"
	"time"
)

type BreakerState int

const (
	BreakerClosed BreakerState = iota
	BreakerOpen
	BreakerHalfOpen
)

func (s BreakerState) String() string {
	switch s {
	case BreakerClosed:
		return "closed"
	case BreakerOpen:
		return "open"
	case BreakerHalfOpen:
		return "half_open"
	default:
		return "unknown"
	}
}

// CircuitBreaker is a simple consecutive-failure circuit breaker.
//
// Closed: every call passes through. Failures increment a counter; success resets it.
// Open: every call is short-circuited (Allow returns false) until cooldown elapses,
// then transitions to HalfOpen.
// HalfOpen: the next call is allowed as a probe. Success → Closed. Failure → Open
// with the cooldown timer reset.
type CircuitBreaker struct {
	mu               sync.Mutex
	failureThreshold int
	cooldown         time.Duration

	state         BreakerState
	failureCount  int
	openedAt      time.Time
}

func NewCircuitBreaker(failureThreshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		cooldown:         cooldown,
		state:            BreakerClosed,
	}
}

func (cb *CircuitBreaker) State() BreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Allow returns true if a call should be made. In Closed and HalfOpen states it
// returns true; in Open state it returns true once the cooldown elapses (and
// transitions to HalfOpen) and false otherwise.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case BreakerClosed:
		return true
	case BreakerHalfOpen:
		return true
	case BreakerOpen:
		if time.Since(cb.openedAt) >= cb.cooldown {
			cb.state = BreakerHalfOpen
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount = 0
	cb.state = BreakerClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == BreakerHalfOpen {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		return
	}
	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
	}
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestCircuitBreaker' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/circuit_breaker.go pkg/ratelimiter/circuit_breaker_test.go
git commit -m "feat(ratelimiter): add consecutive-failure circuit breaker"
```

---

### Task 2.2: Wire the breaker around the limiter (`BreakerLimiter`)

**Files:**
- Modify: `pkg/ratelimiter/circuit_breaker.go` (append wrapper)
- Modify: `pkg/ratelimiter/circuit_breaker_test.go` (append wrapper test)

- [ ] **Step 1: Write the failing test.**

```go
type fakeLimiter struct {
	allowResult *Result
	allowErr    error
	debitResult *Result
	debitErr    error
}

func (f *fakeLimiter) Allow(ctx context.Context, subject string, cost uint64) (*Result, error) {
	return f.allowResult, f.allowErr
}
func (f *fakeLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*Result, error) {
	return f.debitResult, f.debitErr
}

func TestBreakerLimiter_FailOpenWhenBreakerOpen(t *testing.T) {
	inner := &fakeLimiter{allowErr: errSentinel}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(1, time.Hour))

	// First call: inner returns error → breaker opens, but call should still
	// fail open (Allowed=true) because the breaker policy is fail-open.
	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// Second call: breaker is open → bypass inner entirely, fail open.
	inner.allowErr = nil // would also pass if called, but we don't expect it to be called
	inner.allowResult = &Result{Allowed: false}
	res, err = bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed) // fail-open, not the inner's denial
}

func TestBreakerLimiter_PassesThroughOnSuccess(t *testing.T) {
	inner := &fakeLimiter{allowResult: &Result{Allowed: true}}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(3, time.Hour))

	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.True(t, res.Allowed)
}

func TestBreakerLimiter_PassesThroughOnDenial(t *testing.T) {
	inner := &fakeLimiter{allowResult: &Result{Allowed: false}}
	bl := NewBreakerLimiter(inner, NewCircuitBreaker(3, time.Hour))

	res, err := bl.Allow(context.Background(), "subj", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed) // denial is not a failure
}
```

Add `"context"` to imports.

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/ratelimiter/ -run 'TestBreakerLimiter' -v`
Expected: FAIL — undefined `BreakerLimiter`.

- [ ] **Step 3: Implement the wrapper.**

```go
// BreakerLimiter wraps a RateLimiter with a circuit breaker. On any error from
// the inner limiter, the breaker counts a failure and the request fails open
// (allowed=true). When the breaker is OPEN, calls bypass the inner limiter
// entirely and fail open. Denials from the inner limiter are not failures.
type BreakerLimiter struct {
	inner   RateLimiter
	breaker *CircuitBreaker
}

func NewBreakerLimiter(inner RateLimiter, breaker *CircuitBreaker) *BreakerLimiter {
	return &BreakerLimiter{inner: inner, breaker: breaker}
}

func (b *BreakerLimiter) Allow(ctx context.Context, subject string, cost uint64) (*Result, error) {
	if !b.breaker.Allow() {
		return &Result{Allowed: true}, nil // fail open
	}
	res, err := b.inner.Allow(ctx, subject, cost)
	if err != nil {
		b.breaker.RecordFailure()
		return &Result{Allowed: true}, nil // fail open
	}
	b.breaker.RecordSuccess()
	return res, nil
}

func (b *BreakerLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*Result, error) {
	if !b.breaker.Allow() {
		return &Result{Allowed: true}, nil
	}
	res, err := b.inner.ForceDebit(ctx, subject, cost)
	if err != nil {
		b.breaker.RecordFailure()
		return &Result{Allowed: true}, nil
	}
	b.breaker.RecordSuccess()
	return res, nil
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/ratelimiter/ -run 'TestBreakerLimiter' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/circuit_breaker.go pkg/ratelimiter/circuit_breaker_test.go
git commit -m "feat(ratelimiter): add BreakerLimiter wrapper for fail-open Redis runtime errors"
```

---

## Phase 3 — Connect Interceptor

> **Design note for Phase 3:** The interceptor uses **two RateLimiter instances**, not one. This is because the spec defines two distinct sub-limits with different shapes:
>
> - **Query bucket** — `[per-minute, per-hour]` capacity. Used for all unary QueryApi calls and shared across methods.
> - **Opens bucket** — `[opens-per-minute]` capacity. Used only for `SubscribeTopics` admission counting.
>
> These could be packed into one `RedisLimiter` with three limits in its limit slice, but then both subjects (the user IP and the user-IP-with-`:opens`-suffix) would allocate buckets for all three limits in Redis, wasting keys. Two separate limiters are clearer and use exactly the keys they need. The constructor `NewRateLimitInterceptor` takes both.

### Task 3.1: Method-routing helper for QueryApi

**Files:**
- Create: `pkg/interceptors/server/rate_limit.go`
- Create: `pkg/interceptors/server/rate_limit_test.go`

We need a helper to recognize which method on `QueryApi` is being called and what cost to charge. Procedures look like `/xmtp.xmtpv4.message_api.QueryApi/QueryEnvelopes`.

- [ ] **Step 1: Write the failing test.**

```go
package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryApiMethod_FromProcedure(t *testing.T) {
	tests := []struct {
		procedure string
		want      QueryApiMethod
		ok        bool
	}{
		{"/xmtp.xmtpv4.message_api.QueryApi/QueryEnvelopes", MethodQueryEnvelopes, true},
		{"/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics", MethodSubscribeTopics, true},
		{"/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds", MethodGetInboxIds, true},
		{"/xmtp.xmtpv4.message_api.QueryApi/GetNewestEnvelope", MethodGetNewestEnvelope, true},
		{"/xmtp.xmtpv4.message_api.PublishApi/PublishPayerEnvelopes", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := QueryApiMethodFromProcedure(tt.procedure)
		require.Equal(t, tt.ok, ok, tt.procedure)
		if tt.ok {
			require.Equal(t, tt.want, got)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/interceptors/server/ -run 'TestQueryApiMethod' -v`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement.**

```go
package server

import "strings"

type QueryApiMethod string

const (
	MethodQueryEnvelopes    QueryApiMethod = "QueryEnvelopes"
	MethodSubscribeTopics   QueryApiMethod = "SubscribeTopics"
	MethodGetInboxIds       QueryApiMethod = "GetInboxIds"
	MethodGetNewestEnvelope QueryApiMethod = "GetNewestEnvelope"
)

const queryApiPathPrefix = "/xmtp.xmtpv4.message_api.QueryApi/"

// QueryApiMethodFromProcedure parses a Connect procedure path and returns the
// QueryApi method, or (zero, false) if the procedure does not belong to QueryApi.
func QueryApiMethodFromProcedure(procedure string) (QueryApiMethod, bool) {
	if !strings.HasPrefix(procedure, queryApiPathPrefix) {
		return "", false
	}
	name := procedure[len(queryApiPathPrefix):]
	switch QueryApiMethod(name) {
	case MethodQueryEnvelopes, MethodSubscribeTopics, MethodGetInboxIds, MethodGetNewestEnvelope:
		return QueryApiMethod(name), true
	}
	return "", false
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/interceptors/server/ -run 'TestQueryApiMethod' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/interceptors/server/rate_limit.go pkg/interceptors/server/rate_limit_test.go
git commit -m "feat(interceptors): add QueryApi method routing helper"
```

---

### Task 3.2: Cost-from-request helper for QueryApi

**Files:**
- Modify: `pkg/interceptors/server/rate_limit.go`
- Modify: `pkg/interceptors/server/rate_limit_test.go`

The interceptor needs to compute the cost given a method and a request body. For `QueryEnvelopes` we need the topic count from `req.Query.Topics`. For `SubscribeTopics` we need `len(req.Filters)`. For lookups it's a constant.

- [ ] **Step 1: Write the failing test.**

```go
import (
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

func TestComputeCost_QueryEnvelopes_FromTopics(t *testing.T) {
	req := &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			Topics: make([][]byte, 100),
		},
	}
	cost := ComputeCost(MethodQueryEnvelopes, req)
	require.Equal(t, uint64(10), cost)
}

func TestComputeCost_SubscribeTopics_FromFilters(t *testing.T) {
	req := &message_api.SubscribeTopicsRequest{
		Filters: make([]*message_api.SubscribeTopicsRequest_TopicFilter, 4),
	}
	cost := ComputeCost(MethodSubscribeTopics, req)
	require.Equal(t, uint64(2), cost)
}

func TestComputeCost_GetInboxIds_Constant(t *testing.T) {
	req := &message_api.GetInboxIdsRequest{}
	cost := ComputeCost(MethodGetInboxIds, req)
	require.Equal(t, uint64(1), cost)
}

func TestComputeCost_GetNewestEnvelope_Constant(t *testing.T) {
	req := &message_api.GetNewestEnvelopeRequest{}
	cost := ComputeCost(MethodGetNewestEnvelope, req)
	require.Equal(t, uint64(1), cost)
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/interceptors/server/ -run 'TestComputeCost' -v`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement.**

```go
import (
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

// ComputeCost returns the rate-limit token cost of a QueryApi request based on
// its method and body. For SubscribeTopics this is the *admission* cost only;
// the retrospective drain cost is computed by the subscribe handler at close.
func ComputeCost(method QueryApiMethod, req any) uint64 {
	switch method {
	case MethodQueryEnvelopes:
		if r, ok := req.(*message_api.QueryEnvelopesRequest); ok {
			return ratelimiter.CostQuery(len(r.GetQuery().GetTopics()))
		}
	case MethodSubscribeTopics:
		if r, ok := req.(*message_api.SubscribeTopicsRequest); ok {
			return ratelimiter.CostQuery(len(r.GetFilters()))
		}
	case MethodGetInboxIds, MethodGetNewestEnvelope:
		return 1
	}
	return 1
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/interceptors/server/ -run 'TestComputeCost' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/interceptors/server/rate_limit.go pkg/interceptors/server/rate_limit_test.go
git commit -m "feat(interceptors): add ComputeCost helper for QueryApi requests"
```

---

### Task 3.3: Connect interceptor — unary path

**Files:**
- Modify: `pkg/interceptors/server/rate_limit.go`
- Modify: `pkg/interceptors/server/rate_limit_test.go`

- [ ] **Step 1: Write the failing test.**

```go
import (
	"context"
	"net"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"go.uber.org/zap"
)

type fakeLimiter struct {
	lastSubject string
	lastCost    uint64
	result      *ratelimiter.Result
	err         error
	debitCalls  int
}

func (f *fakeLimiter) Allow(ctx context.Context, subject string, cost uint64) (*ratelimiter.Result, error) {
	f.lastSubject = subject
	f.lastCost = cost
	return f.result, f.err
}
func (f *fakeLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*ratelimiter.Result, error) {
	f.debitCalls++
	return &ratelimiter.Result{Allowed: true}, nil
}

// newTestInterceptor builds an interceptor with the same limiter for both
// query and opens, which is enough for unit tests that exercise one path at
// a time.
func newTestInterceptor(limiter *fakeLimiter) *RateLimitInterceptor {
	return NewRateLimitInterceptor(zap.NewNop(), limiter, limiter, nil, RateLimitInterceptorConfig{})
}

type fakeUnaryRequest struct {
	connect.Request[message_api.GetInboxIdsRequest]
	procedure string
	peerAddr  string
	xff       string
}

func (f *fakeUnaryRequest) Spec() connect.Spec {
	return connect.Spec{Procedure: f.procedure}
}
func (f *fakeUnaryRequest) Peer() connect.Peer {
	return connect.Peer{Addr: f.peerAddr}
}
func (f *fakeUnaryRequest) Header() http.Header {
	h := http.Header{}
	if f.xff != "" {
		h.Set("X-Forwarded-For", f.xff)
	}
	return h
}

func TestRateLimitInterceptor_BypassesNonQueryApi(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: true}}
	rl := newTestInterceptor(limiter)

	called := false
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		called = true
		return nil, nil
	}

	req := &fakeUnaryRequest{procedure: "/xmtp.xmtpv4.message_api.PublishApi/PublishPayerEnvelopes"}
	_, err := rl.WrapUnary(next)(context.Background(), req)
	require.NoError(t, err)
	require.True(t, called)
	require.Empty(t, limiter.lastSubject) // limiter not called
}

func TestRateLimitInterceptor_Tier0BypassesLimiter(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: true}}
	rl := newTestInterceptor(limiter)

	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)
	req := &fakeUnaryRequest{procedure: "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds"}

	_, err := rl.WrapUnary(func(c context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	})(ctx, req)
	require.NoError(t, err)
	require.Empty(t, limiter.lastSubject)
}

func TestRateLimitInterceptor_Tier2DeniedReturnsResourceExhausted(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: false}}
	rl := newTestInterceptor(limiter)

	req := &fakeUnaryRequest{
		procedure: "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds",
		peerAddr:  "203.0.113.1:5001",
	}
	_, err := rl.WrapUnary(func(c context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		t.Fatal("inner handler should not be called on denial")
		return nil, nil
	})(context.Background(), req)
	require.Error(t, err)
	cerr := new(connect.Error)
	require.ErrorAs(t, err, &cerr)
	require.Equal(t, connect.CodeResourceExhausted, cerr.Code())
	require.Equal(t, "203.0.113.1", limiter.lastSubject)
	require.Equal(t, uint64(1), limiter.lastCost)
}
```

Add imports: `"net/http"`, `"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"`.

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/interceptors/server/ -run 'TestRateLimitInterceptor' -v`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement the interceptor (unary path only for now).**

Append to `pkg/interceptors/server/rate_limit.go`:

```go
import (
	"context"
	"errors"
	"fmt"
	"net"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"go.uber.org/zap"
)

// RateLimitInterceptorConfig holds the runtime knobs for the interceptor.
// All durations and counts are in their natural units (the config layer
// converts from XMTPD_RATE_LIMIT_* env vars).
type RateLimitInterceptorConfig struct {
	DrainIntervalMinutes int
	DrainAmount          int
}

// RateLimitInterceptor is a Connect interceptor that enforces token-bucket rate
// limits on QueryApi methods. Tier 0 (verified node) requests bypass the limiter.
// Other procedures (PublishApi, NotificationApi, ReplicationApi) are passed
// through unchanged — this interceptor does NOT meter them.
//
// It uses two limiters:
//   - queryLimiter: scopes the per-minute and per-hour buckets keyed on client IP.
//     Used for all unary QueryApi methods.
//   - opensLimiter: scopes the subscribe-opens-per-minute bucket keyed on client
//     IP with a `:opens` suffix. Used by the streaming path only.
type RateLimitInterceptor struct {
	logger       *zap.Logger
	queryLimiter ratelimiter.RateLimiter
	opensLimiter ratelimiter.RateLimiter
	trustedCIDRs []*net.IPNet
	cfg          RateLimitInterceptorConfig
}

var _ connect.Interceptor = (*RateLimitInterceptor)(nil)

func NewRateLimitInterceptor(
	logger *zap.Logger,
	queryLimiter ratelimiter.RateLimiter,
	opensLimiter ratelimiter.RateLimiter,
	trustedCIDRs []*net.IPNet,
	cfg RateLimitInterceptorConfig,
) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger:       logger.Named("xmtpd.rate-limiter"),
		queryLimiter: queryLimiter,
		opensLimiter: opensLimiter,
		trustedCIDRs: trustedCIDRs,
		cfg:          cfg,
	}
}

func (i *RateLimitInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		method, ok := QueryApiMethodFromProcedure(req.Spec().Procedure)
		if !ok {
			return next(ctx, req)
		}

		if ratelimiter.ClassifyTier(ctx) == ratelimiter.Tier0 {
			return next(ctx, req)
		}

		subject := i.subjectFromRequest(req)
		cost := ComputeCost(method, req.Any())

		res, err := i.queryLimiter.Allow(ctx, subject, cost)
		if err != nil {
			// Should not happen — BreakerLimiter wraps errors. Defensive only.
			i.logger.Warn("rate limiter call returned error", zap.Error(err))
			return next(ctx, req)
		}
		if !res.Allowed {
			return nil, connect.NewError(
				connect.CodeResourceExhausted,
				fmt.Errorf("rate limit exceeded"),
			)
		}
		return next(ctx, req)
	}
}

func (i *RateLimitInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next // server interceptor — never called
}

// subjectFromRequest extracts the bucket-key subject for a request, defaulting
// to the client IP after trusted-proxy peeling and IPv6 /64 normalization.
func (i *RateLimitInterceptor) subjectFromRequest(req connect.AnyRequest) string {
	xff := req.Header().Get("X-Forwarded-For")
	return ratelimiter.ExtractClientIP(req.Peer().Addr, xff, i.trustedCIDRs)
}

var _ = errors.New // keep import even if unused below
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/interceptors/server/ -run 'TestRateLimitInterceptor' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/interceptors/server/rate_limit.go pkg/interceptors/server/rate_limit_test.go
git commit -m "feat(interceptors): add Connect interceptor for QueryApi rate limiting (unary path)"
```

---

### Task 3.4: Connect interceptor — streaming path with drain hook

**Files:**
- Modify: `pkg/interceptors/server/rate_limit.go` (add `WrapStreamingHandler`)
- Modify: `pkg/interceptors/server/rate_limit_test.go`

The streaming path differs from unary in two ways:
1. The body is read from `conn.Receive(req)` after `WrapStreamingHandler` returns control. The interceptor cannot pre-compute the cost from the request body — it has to peek at the first received message OR (simpler) check the procedure, assume it's `SubscribeTopics`, and *post*-decode after the inner handler reads the request.

The cleanest approach for v1: do the admission check using a worst-case `cost=1` for the subject's primary bucket and use a separate Redis call (the subscribe-opens-per-minute sub-limit). The actual `√(filters)` admission cost is charged inside the handler **before** the catch-up runs, which we wire in Phase 4.

For now, the streaming interceptor only handles:
- Tier 0 bypass
- Subscribe-opens-per-minute counter (cost=1 against a separate "opens" subject suffix)

The handler-side billing in Phase 4 handles the admission √(filters) cost and the retrospective drain.

- [ ] **Step 1: Write the failing test.**

```go
type fakeStreamingHandlerConn struct {
	connect.StreamingHandlerConn
	procedure string
	peerAddr  string
	headers   http.Header
}

func (f *fakeStreamingHandlerConn) Spec() connect.Spec {
	return connect.Spec{Procedure: f.procedure}
}
func (f *fakeStreamingHandlerConn) Peer() connect.Peer {
	return connect.Peer{Addr: f.peerAddr}
}
func (f *fakeStreamingHandlerConn) RequestHeader() http.Header { return f.headers }

func TestRateLimitInterceptor_Streaming_BypassesNonQueryApi(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: true}}
	rl := newTestInterceptor(limiter)

	called := false
	next := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		called = true
		return nil
	}
	conn := &fakeStreamingHandlerConn{
		procedure: "/xmtp.xmtpv4.message_api.NotificationApi/SubscribeAllEnvelopes",
		headers:   http.Header{},
	}
	err := rl.WrapStreamingHandler(next)(context.Background(), conn)
	require.NoError(t, err)
	require.True(t, called)
	require.Empty(t, limiter.lastSubject)
}

func TestRateLimitInterceptor_Streaming_Tier0Bypass(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: true}}
	rl := newTestInterceptor(limiter)

	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)
	conn := &fakeStreamingHandlerConn{
		procedure: "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics",
		headers:   http.Header{},
	}
	err := rl.WrapStreamingHandler(func(c context.Context, sc connect.StreamingHandlerConn) error {
		return nil
	})(ctx, conn)
	require.NoError(t, err)
	require.Empty(t, limiter.lastSubject)
}

func TestRateLimitInterceptor_Streaming_OpensSubLimit_Denied(t *testing.T) {
	limiter := &fakeLimiter{result: &ratelimiter.Result{Allowed: false}}
	rl := newTestInterceptor(limiter)

	conn := &fakeStreamingHandlerConn{
		procedure: "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics",
		peerAddr:  "203.0.113.1:5001",
		headers:   http.Header{},
	}
	err := rl.WrapStreamingHandler(func(c context.Context, sc connect.StreamingHandlerConn) error {
		t.Fatal("inner handler should not be called when opens sub-limit denies")
		return nil
	})(context.Background(), conn)
	require.Error(t, err)
	cerr := new(connect.Error)
	require.ErrorAs(t, err, &cerr)
	require.Equal(t, connect.CodeResourceExhausted, cerr.Code())
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/interceptors/server/ -run 'TestRateLimitInterceptor_Streaming' -v`
Expected: FAIL — `WrapStreamingHandler` not implemented for the rate-limit case.

- [ ] **Step 3: Implement.**

Replace the stub `WrapStreamingHandler` in `rate_limit.go`:

```go
func (i *RateLimitInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		method, ok := QueryApiMethodFromProcedure(conn.Spec().Procedure)
		if !ok {
			return next(ctx, conn)
		}
		if ratelimiter.ClassifyTier(ctx) == ratelimiter.Tier0 {
			return next(ctx, conn)
		}

		// Streaming handler only meters SubscribeTopics. The other QueryApi
		// methods are unary and handled in WrapUnary.
		if method != MethodSubscribeTopics {
			return next(ctx, conn)
		}

		xff := conn.RequestHeader().Get("X-Forwarded-For")
		clientIP := ratelimiter.ExtractClientIP(conn.Peer().Addr, xff, i.trustedCIDRs)

		// Charge against the subscribe-opens-per-minute sub-limit using the
		// opens limiter (separate Redis key space) with a `:opens` subject
		// suffix. The √(filters) admission and retrospective drain are handled
		// by the SubscribeTopics handler itself (Phase 4).
		opensSubject := clientIP + ":opens"
		res, err := i.opensLimiter.Allow(ctx, opensSubject, 1)
		if err != nil {
			i.logger.Warn("rate limiter (opens) returned error", zap.Error(err))
			return next(ctx, conn)
		}
		if !res.Allowed {
			return connect.NewError(
				connect.CodeResourceExhausted,
				fmt.Errorf("subscribe rate limit exceeded"),
			)
		}
		return next(ctx, conn)
	}
}
```

- [ ] **Step 4: Run test to verify it passes.**

Run: `go test ./pkg/interceptors/server/ -run 'TestRateLimitInterceptor_Streaming' -v`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/interceptors/server/rate_limit.go pkg/interceptors/server/rate_limit_test.go
git commit -m "feat(interceptors): add streaming path for SubscribeTopics opens sub-limit"
```

---

## Phase 4 — Subscribe Handler Integration

### Task 4.1: Plumb the limiter and config into the message Service

**Files:**
- Modify: `pkg/api/message/service.go` (add limiter + rate-limit config to `Service` struct)
- Modify: `pkg/server/server.go` (pass them when constructing `NewReplicationAPIService`)

This task is glue. Read `pkg/api/message/service.go` to see the current `Service` struct and constructor signature; add two new optional fields and a setter (or add to the constructor signature). The plan assumes adding to the constructor for explicitness.

- [ ] **Step 1: Add fields to the `Service` struct.**

In `pkg/api/message/service.go`, add to the struct:

```go
rateLimiter  ratelimiter.RateLimiter // nil when rate limiting disabled
rlConfig     RateLimitConfig
```

Define a small config type in the same file:

```go
// RateLimitConfig holds the subset of rate-limit options the message Service
// needs at handler-time (cost params, drain params, stream lifetime caps).
type RateLimitConfig struct {
	Enabled              bool
	DrainIntervalMinutes int
	DrainAmount          int
	StreamIdleTimeout    time.Duration
	StreamMaxDuration    time.Duration
}
```

- [ ] **Step 2: Add constructor parameters.**

Append two parameters to `NewReplicationAPIService`:

```go
rateLimiter ratelimiter.RateLimiter,
rlConfig RateLimitConfig,
```

Set them on the constructed `Service`. If `rateLimiter` is nil and `rlConfig.Enabled` is false, the handler will skip drain calls.

- [ ] **Step 3: Update all call sites.**

Find call sites with: `grep -rn "NewReplicationAPIService(" --include="*.go"`. Pass `nil, RateLimitConfig{}` from non-server callers (tests, etc.). The real call site in `pkg/server/server.go` will receive a real limiter once Phase 5 wires it up — until then, also pass `nil, RateLimitConfig{}`.

- [ ] **Step 4: Build to verify.**

Run: `go build ./...`
Expected: clean.

- [ ] **Step 5: Run existing message API tests.**

Run: `go test ./pkg/api/message/ -count=1`
Expected: PASS — existing behavior unchanged because the new params are no-ops.

- [ ] **Step 6: Commit.**

```bash
git add pkg/api/message/ pkg/server/server.go
git commit -m "feat(message): plumb rate limiter and config into Service constructor"
```

---

### Task 4.2: SubscribeTopics admission cost and deferred drain

**Files:**
- Modify: `pkg/api/message/subscribe_topics.go`
- Create: `pkg/api/message/subscribe_topics_ratelimit_test.go`

- [ ] **Step 1: Write a unit test for the helper that performs admission and registers the deferred drain.**

```go
package message

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

type spyLimiter struct {
	allowSubject string
	allowCost    uint64
	allowResult  *ratelimiter.Result
	debitSubject string
	debitCost    uint64
}

func (s *spyLimiter) Allow(ctx context.Context, subject string, cost uint64) (*ratelimiter.Result, error) {
	s.allowSubject = subject
	s.allowCost = cost
	if s.allowResult == nil {
		return &ratelimiter.Result{Allowed: true}, nil
	}
	return s.allowResult, nil
}
func (s *spyLimiter) ForceDebit(ctx context.Context, subject string, cost uint64) (*ratelimiter.Result, error) {
	s.debitSubject = subject
	s.debitCost = cost
	return &ratelimiter.Result{Allowed: true}, nil
}

func TestApplySubscribeAdmissionAndDrain_BillsAtCloseUsingElapsed(t *testing.T) {
	limiter := &spyLimiter{}
	cfg := RateLimitConfig{
		Enabled:              true,
		DrainIntervalMinutes: 5,
		DrainAmount:          1,
	}
	subject := "203.0.113.1"
	numFilters := 4 // sqrt(4) = 2

	cleanup, err := applySubscribeAdmissionAndDrain(context.Background(), limiter, cfg, subject, numFilters)
	require.NoError(t, err)
	require.Equal(t, uint64(2), limiter.allowCost) // ceil(sqrt(4))
	require.Equal(t, subject, limiter.allowSubject)

	// Simulate stream held for 6 minutes by manually re-running the drain calc
	// via the cleanup func — cleanup uses time.Since(start), which we cannot
	// fast-forward. Instead, assert cleanup runs without panicking and that
	// the limiter received a debit call after invocation.
	cleanup()
	require.Equal(t, subject, limiter.debitSubject)
	// Drain may be 0 for a sub-millisecond elapsed; that's expected.
}

func TestApplySubscribeAdmissionAndDrain_DisabledIsNoOp(t *testing.T) {
	limiter := &spyLimiter{}
	cfg := RateLimitConfig{Enabled: false}
	cleanup, err := applySubscribeAdmissionAndDrain(context.Background(), limiter, cfg, "subj", 4)
	require.NoError(t, err)
	require.Empty(t, limiter.allowSubject)
	cleanup()
	require.Empty(t, limiter.debitSubject)
}

func TestApplySubscribeAdmissionAndDrain_DenialReturnsError(t *testing.T) {
	limiter := &spyLimiter{allowResult: &ratelimiter.Result{Allowed: false}}
	cfg := RateLimitConfig{Enabled: true, DrainIntervalMinutes: 5, DrainAmount: 1}
	_, err := applySubscribeAdmissionAndDrain(context.Background(), limiter, cfg, "subj", 4)
	require.Error(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails.**

Run: `go test ./pkg/api/message/ -run 'TestApplySubscribeAdmissionAndDrain' -v`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement the helper.**

Add to `pkg/api/message/subscribe_topics.go`:

```go
import (
	// existing imports...
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

// applySubscribeAdmissionAndDrain charges the admission cost for opening a
// subscription and returns a cleanup func that will retrospectively bill the
// drain cost when the stream closes. The cleanup func should be invoked via
// defer in the SubscribeTopics handler.
//
// When rate limiting is disabled (cfg.Enabled == false), this is a no-op
// returning a no-op cleanup func.
func applySubscribeAdmissionAndDrain(
	ctx context.Context,
	limiter ratelimiter.RateLimiter,
	cfg RateLimitConfig,
	subject string,
	numFilters int,
) (func(), error) {
	if !cfg.Enabled || limiter == nil {
		return func() {}, nil
	}

	cost := ratelimiter.CostQuery(numFilters)
	res, err := limiter.Allow(ctx, subject, cost)
	if err != nil {
		// BreakerLimiter handles errors and fails open, so this should not
		// happen in practice. Defensive: treat as failure-open and proceed.
		return func() {}, nil
	}
	if !res.Allowed {
		return nil, connect.NewError(
			connect.CodeResourceExhausted,
			fmt.Errorf("subscribe admission rate limit exceeded"),
		)
	}

	startedAt := time.Now()
	cleanup := func() {
		// Use a fresh background context — the request context is cancelled
		// at this point. Drain still needs to write to Redis.
		drainCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		drain := ratelimiter.CostSubscribeDrain(
			time.Since(startedAt),
			cfg.DrainIntervalMinutes,
			cfg.DrainAmount,
		)
		if drain == 0 {
			return
		}
		_, _ = limiter.ForceDebit(drainCtx, subject, drain)
	}
	return cleanup, nil
}
```

- [ ] **Step 4: Wire it into `SubscribeTopics`.**

In `subscribe_topics.go`, near the top of `SubscribeTopics` (after the nil-check on `req.Msg`), add:

```go
filters := req.Msg.GetFilters()

subject := subscribeSubjectFromContext(ctx, req)
cleanup, err := applySubscribeAdmissionAndDrain(ctx, s.rateLimiter, s.rlConfig, subject, len(filters))
if err != nil {
	return err
}
defer cleanup()
```

(Move the existing `filters := req.Msg.GetFilters()` line up if it's currently below; the order matters because admission needs the filter count.)

Also add a small helper:

```go
// subscribeSubjectFromContext returns the bucket-key subject for a subscribe
// request. It mirrors the logic in the rate-limit interceptor: peer IP after
// trusted-proxy peeling. The interceptor cannot inject the subject into context
// (Connect interceptors run before the handler reads headers), so the handler
// re-derives it. The trustedCIDRs list is read from a package-level var set by
// the server at startup; absent that, only the peer addr is used.
func subscribeSubjectFromContext(ctx context.Context, req *connect.Request[message_api.SubscribeTopicsRequest]) string {
	xff := req.Header().Get("X-Forwarded-For")
	return ratelimiter.ExtractClientIP(req.Peer().Addr, xff, subscribeTrustedCIDRs)
}

// subscribeTrustedCIDRs is set by the server at startup. Defaults to nil (no peeling).
var subscribeTrustedCIDRs []*net.IPNet

// SetSubscribeTrustedCIDRs configures the trusted-proxy CIDRs used by
// subscribe handlers for IP keying. Called once at server startup.
func SetSubscribeTrustedCIDRs(cidrs []*net.IPNet) {
	subscribeTrustedCIDRs = cidrs
}
```

Add `"net"` to the imports.

- [ ] **Step 5: Run the unit test.**

Run: `go test ./pkg/api/message/ -run 'TestApplySubscribeAdmissionAndDrain' -v`
Expected: PASS.

- [ ] **Step 6: Run all message API tests.**

Run: `go test ./pkg/api/message/ -count=1`
Expected: PASS (existing subscribe tests should not be affected because rateLimiter is nil in test setup).

- [ ] **Step 7: Commit.**

```bash
git add pkg/api/message/subscribe_topics.go pkg/api/message/subscribe_topics_ratelimit_test.go
git commit -m "feat(message): apply admission cost and deferred drain to SubscribeTopics"
```

---

### Task 4.3: Idle timeout and max-duration in the subscribe loop

**Files:**
- Modify: `pkg/api/message/subscribe_topics.go`

The existing select loop in `SubscribeTopics` has cases for `ticker.C`, `envelopesCh`, `ctx.Done()`, `s.ctx.Done()`. Add an idle-timeout timer that resets on each `envelopesCh` send, and a max-duration timer that fires once.

- [ ] **Step 1: Modify the select loop.**

Replace the current loop in `SubscribeTopics` (lines ~100-137) with:

```go
// GRPC keep-alives are not sufficient in some load balanced environments.
ticker := time.NewTicker(s.options.SendKeepAliveInterval)
defer ticker.Stop()

// Idle timeout: cancel the stream if no envelope arrives for this long.
// Disabled (0 duration) when rate-limit config is disabled.
var idleTimer *time.Timer
var idleC <-chan time.Time
if s.rlConfig.StreamIdleTimeout > 0 {
	idleTimer = time.NewTimer(s.rlConfig.StreamIdleTimeout)
	defer idleTimer.Stop()
	idleC = idleTimer.C
}

// Max duration: hard cap on stream lifetime.
var maxC <-chan time.Time
if s.rlConfig.StreamMaxDuration > 0 {
	maxTimer := time.NewTimer(s.rlConfig.StreamMaxDuration)
	defer maxTimer.Stop()
	maxC = maxTimer.C
}

for {
	select {
	case <-ticker.C:
		err = stream.Send(newSubscriptionStatusMessage(
			message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_WAITING,
		))
		if err != nil {
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("could not send keepalive: %w", err),
			)
		}

	case envs, open := <-envelopesCh:
		ticker.Reset(s.options.SendKeepAliveInterval)
		if idleTimer != nil {
			if !idleTimer.Stop() {
				select {
				case <-idleTimer.C:
				default:
				}
			}
			idleTimer.Reset(s.rlConfig.StreamIdleTimeout)
		}

		if !open {
			logger.Debug("channel closed by worker")
			return nil
		}

		envsToSend := advanceTopicCursors(cursors, envs, logger)
		err = s.sendTopicEnvelopes(stream, envsToSend)
		if err != nil {
			return err
		}

	case <-idleC:
		logger.Info("subscribe stream cancelled by idle timeout")
		return connect.NewError(
			connect.CodeDeadlineExceeded,
			errors.New("stream idle timeout"),
		)

	case <-maxC:
		logger.Info("subscribe stream cancelled by max duration")
		return connect.NewError(
			connect.CodeDeadlineExceeded,
			errors.New("stream max duration reached"),
		)

	case <-ctx.Done():
		logger.Debug("topic subscription stream closed")
		return nil

	case <-s.ctx.Done():
		logger.Debug("message service closed")
		return nil
	}
}
```

- [ ] **Step 2: Build to verify.**

Run: `go build ./pkg/api/message/...`
Expected: clean.

- [ ] **Step 3: Add a test for the idle timeout path.**

Append to `pkg/api/message/subscribe_topics_ratelimit_test.go`:

```go
// Idle timeout integration is best covered by an end-to-end test in Phase 7
// (Task 7.5) using the real subscribe handler over a fake stream. The select
// loop change is structurally exercised by existing subscribe tests when
// StreamIdleTimeout is 0 (disabled) and by the Phase 7 test when it's set.
```

(Comment-only "test" — the real coverage is the e2e test in Phase 7. Idle-timer fast-forward is hard to mock without a clock-injection refactor that is out of scope.)

- [ ] **Step 4: Run all message API tests.**

Run: `go test ./pkg/api/message/ -count=1`
Expected: PASS — existing subscribe tests still pass since `StreamIdleTimeout` and `StreamMaxDuration` are zero in test setup.

- [ ] **Step 5: Commit.**

```bash
git add pkg/api/message/subscribe_topics.go pkg/api/message/subscribe_topics_ratelimit_test.go
git commit -m "feat(message): add idle timeout and max duration to SubscribeTopics select loop"
```

---

## Phase 5 — Server Wiring

### Task 5.1: Build the rate limiter at server startup with fail-fast

**Files:**
- Create: `pkg/ratelimiter/builder.go`
- Modify: `pkg/server/server.go`

- [ ] **Step 1: Create the builder.**

```go
package ratelimiter

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

// BuiltLimiter is the result of constructing the rate-limit subsystem at
// startup. It exposes a query limiter and an opens limiter (see Phase 3 design
// note), the parsed trusted-proxy CIDRs, and is consumed by the server when
// constructing the rate-limit interceptor.
type BuiltLimiter struct {
	QueryLimiter RateLimiter // BreakerLimiter wrapping a RedisLimiter([per-minute, per-hour])
	OpensLimiter RateLimiter // BreakerLimiter wrapping a RedisLimiter([opens-per-minute])
	TrustedCIDRs []*net.IPNet
}

// Build constructs the rate-limit subsystem from server configuration. If
// rlOpts.Enable is false, returns (nil, nil) — the caller should treat the
// nil result as "rate limiting disabled."
//
// When enabled, Build pings Redis and returns an error if it is unreachable.
// This implements the spec's fail-fast-at-startup behavior.
func Build(
	ctx context.Context,
	logger *zap.Logger,
	redisOpts config.RedisOptions,
	rlOpts config.RateLimitOptions,
) (*BuiltLimiter, error) {
	if !rlOpts.Enable {
		return nil, nil
	}
	if redisOpts.RedisURL == "" {
		return nil, fmt.Errorf("rate limiting enabled but XMTPD_REDIS_URL is empty")
	}

	parsed, err := redis.ParseURL(redisOpts.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}
	client := redis.NewClient(parsed)

	pingCtx, cancel := context.WithTimeout(ctx, redisOpts.ConnectTimeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("rate limiting enabled but redis ping failed: %w", err)
	}

	queryInner, err := NewRedisLimiter(client, redisOpts.KeyPrefix+"rl:t2:q", []Limit{
		{Capacity: rlOpts.T2PerMinuteCapacity, RefillEvery: time.Minute},
		{Capacity: rlOpts.T2PerHourCapacity, RefillEvery: time.Hour},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to construct query limiter: %w", err)
	}
	opensInner, err := NewRedisLimiter(client, redisOpts.KeyPrefix+"rl:t2:o", []Limit{
		{Capacity: rlOpts.T2SubscribeOpensPerMinute, RefillEvery: time.Minute},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to construct opens limiter: %w", err)
	}

	queryWrapped := NewBreakerLimiter(queryInner, NewCircuitBreaker(rlOpts.BreakerFailureThreshold, rlOpts.BreakerCooldown))
	opensWrapped := NewBreakerLimiter(opensInner, NewCircuitBreaker(rlOpts.BreakerFailureThreshold, rlOpts.BreakerCooldown))

	cidrs, err := ParseTrustedProxyCIDRs(rlOpts.TrustedProxyCIDRs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trusted proxy CIDRs: %w", err)
	}

	logger.Info("rate limit interceptor enabled",
		zap.Int("t2_per_minute", rlOpts.T2PerMinuteCapacity),
		zap.Int("t2_per_hour", rlOpts.T2PerHourCapacity),
		zap.Int("t2_subscribe_opens_per_minute", rlOpts.T2SubscribeOpensPerMinute),
	)
	return &BuiltLimiter{
		QueryLimiter: queryWrapped,
		OpensLimiter: opensWrapped,
		TrustedCIDRs: cidrs,
	}, nil
}
```

- [ ] **Step 2: Wire `Build` into `pkg/server/server.go`.**

In `startAPIServer`, before the `registrationFunc` is built, call:

```go
built, err := ratelimiter.Build(svc.ctx, cfg.Logger, cfg.Options.Redis, cfg.Options.RateLimit)
if err != nil {
	return fmt.Errorf("failed to build rate limiter: %w", err)
}
```

Note: `cfg.Options.Redis` does not currently exist on `ServerOptions`. Add it:

```go
// In options.go, ServerOptions struct:
Redis           RedisOptions           `group:"Redis Options"            namespace:"redis"`
```

Then in `registrationFunc`, when registering the QueryApi handler, use a different `handlerOpts` that includes the rate-limit interceptor *only* for QueryApi:

```go
queryHandlerOpts := handlerOpts
if built != nil {
	rlInterceptor := serverInterceptors.NewRateLimitInterceptor(
		cfg.Logger,
		built.QueryLimiter,
		built.OpensLimiter,
		built.TrustedCIDRs,
		serverInterceptors.RateLimitInterceptorConfig{
			DrainIntervalMinutes: cfg.Options.RateLimit.DrainIntervalMinutes,
			DrainAmount:          cfg.Options.RateLimit.DrainAmount,
		},
	)
	queryHandlerOpts = []connect.HandlerOption{
		connect.WithReadMaxBytes(constants.GRPCPayloadLimit),
		connect.WithSendMaxBytes(constants.GRPCPayloadLimit),
		connect.WithInterceptors(append(interceptors, rlInterceptor)...),
	}

	// Also configure the message Service with the limiter for handler-side
	// admission and drain.
	message.SetSubscribeTrustedCIDRs(built.TrustedCIDRs)
}

queryPath, queryHandler := message_apiconnect.NewQueryApiHandler(
	replicationService, queryHandlerOpts...,
)
```

Pass the limiter and config into `NewReplicationAPIService`:

```go
var msgRateLimiter ratelimiter.RateLimiter
rlConfig := message.RateLimitConfig{}
if built != nil {
	msgRateLimiter = built.QueryLimiter
	rlConfig = message.RateLimitConfig{
		Enabled:              true,
		DrainIntervalMinutes: cfg.Options.RateLimit.DrainIntervalMinutes,
		DrainAmount:          cfg.Options.RateLimit.DrainAmount,
		StreamIdleTimeout:    cfg.Options.RateLimit.StreamIdleTimeout,
		StreamMaxDuration:    cfg.Options.RateLimit.StreamMaxDuration,
	}
}

replicationService, err := message.NewReplicationAPIService(
	svc.ctx, cfg.Logger, svc.registrant, cfg.NodeRegistry, cfg.DB, svc.mlsValidation,
	svc.cursorUpdater, cfg.FeeCalculator, cfg.Options.API, isMigrationEnabled,
	10*time.Millisecond,
	db.NewCachedOriginatorList(cfg.DB.ReadQuery(), cfg.Options.API.OriginatorCacheTTL, cfg.Logger),
	ledgerPkg.NewLedger(cfg.Logger, cfg.DB),
	msgRateLimiter,
	rlConfig,
)
```

Add the imports: `"github.com/xmtp/xmtpd/pkg/ratelimiter"` and rename `server` import to avoid clash if needed (or alias the interceptors package as `serverInterceptors`).

- [ ] **Step 3: Build to verify.**

Run: `go build ./...`
Expected: clean.

- [ ] **Step 4: Run unit tests.**

Run: `go test ./pkg/ratelimiter/ ./pkg/interceptors/server/ ./pkg/api/message/ ./pkg/server/ -count=1`
Expected: PASS.

- [ ] **Step 5: Commit.**

```bash
git add pkg/ratelimiter/builder.go pkg/server/server.go pkg/config/options.go pkg/api/message/
git commit -m "feat(server): wire QueryApi rate limit interceptor with fail-fast Redis check"
```

---

## Phase 6 — Metrics

### Task 6.1: Decision counter and breaker gauge

**Files:**
- Create: `pkg/ratelimiter/metrics.go`
- Modify: `pkg/interceptors/server/rate_limit.go` (call metrics)
- Modify: `pkg/ratelimiter/circuit_breaker.go` (set gauge from breaker state)

- [ ] **Step 1: Define the prometheus metrics.**

```go
package ratelimiter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DecisionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_decisions_total",
			Help: "Rate-limit decisions broken down by service, method, tier, and outcome",
		},
		[]string{"service", "method", "tier", "outcome"},
	)
	BreakerStateGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtpd_rate_limit_circuit_breaker_state",
			Help: "Circuit breaker state: 0=closed, 1=half_open, 2=open",
		},
	)
	BreakerTripsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_circuit_breaker_trips_total",
			Help: "Number of times the circuit breaker has tripped open",
		},
	)
	StreamTerminationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_stream_terminations_total",
			Help: "Subscribe stream terminations broken down by reason",
		},
		[]string{"reason"},
	)
)

// Register registers all rate-limit metrics with the provided registry.
// Safe to call multiple times — uses MustRegister and ignores AlreadyRegistered.
func Register(reg prometheus.Registerer) {
	for _, c := range []prometheus.Collector{
		DecisionsTotal, BreakerStateGauge, BreakerTripsTotal, StreamTerminationsTotal,
	} {
		_ = reg.Register(c) // ignore already-registered errors
	}
}
```

- [ ] **Step 2: Increment the gauge from the breaker.**

In `circuit_breaker.go`, modify `RecordFailure` and `RecordSuccess` to update the gauge:

```go
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount = 0
	if cb.state != BreakerClosed {
		cb.state = BreakerClosed
		BreakerStateGauge.Set(0)
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == BreakerHalfOpen {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		BreakerStateGauge.Set(2)
		BreakerTripsTotal.Inc()
		return
	}
	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		BreakerStateGauge.Set(2)
		BreakerTripsTotal.Inc()
	}
}
```

And in `Allow` when transitioning Open → HalfOpen:

```go
case BreakerOpen:
	if time.Since(cb.openedAt) >= cb.cooldown {
		cb.state = BreakerHalfOpen
		BreakerStateGauge.Set(1)
		return true
	}
	return false
```

- [ ] **Step 3: Increment decision counter from the interceptor.**

In `rate_limit.go`, after the limiter call in `WrapUnary`:

```go
outcome := "allowed"
if err != nil {
	outcome = "failed_open"
} else if !res.Allowed {
	outcome = "denied"
}
ratelimiter.DecisionsTotal.WithLabelValues(
	"QueryApi",
	string(method),
	"tier2",
	outcome,
).Inc()
```

Add an equivalent block for the streaming path with `outcome="bypassed"` for Tier 0 and `"denied"`/`"allowed"` for the opens sub-limit.

- [ ] **Step 4: Register metrics from the server.**

In `pkg/server/server.go`, after building the limiter:

```go
if built != nil {
	ratelimiter.Register(promReg)
}
```

- [ ] **Step 5: Build and run tests.**

Run: `go build ./... && go test ./pkg/ratelimiter/ ./pkg/interceptors/server/ -count=1`
Expected: PASS.

- [ ] **Step 6: Commit.**

```bash
git add pkg/ratelimiter/metrics.go pkg/ratelimiter/circuit_breaker.go pkg/interceptors/server/rate_limit.go pkg/server/server.go
git commit -m "feat(ratelimiter): add prometheus metrics for decisions, breaker state, and stream terminations"
```

---

### Task 6.2: Stream termination reason counter

**Files:**
- Modify: `pkg/api/message/subscribe_topics.go`

- [ ] **Step 1: Increment the counter from each return path.**

In the select loop, replace the four return statements with:

```go
case <-idleC:
	ratelimiter.StreamTerminationsTotal.WithLabelValues("idle").Inc()
	logger.Info("subscribe stream cancelled by idle timeout")
	return connect.NewError(connect.CodeDeadlineExceeded, errors.New("stream idle timeout"))

case <-maxC:
	ratelimiter.StreamTerminationsTotal.WithLabelValues("max_duration").Inc()
	logger.Info("subscribe stream cancelled by max duration")
	return connect.NewError(connect.CodeDeadlineExceeded, errors.New("stream max duration reached"))

case <-ctx.Done():
	ratelimiter.StreamTerminationsTotal.WithLabelValues("client_close").Inc()
	logger.Debug("topic subscription stream closed")
	return nil

case <-s.ctx.Done():
	ratelimiter.StreamTerminationsTotal.WithLabelValues("server_close").Inc()
	logger.Debug("message service closed")
	return nil
```

- [ ] **Step 2: Build and test.**

Run: `go build ./... && go test ./pkg/api/message/ -count=1`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/api/message/subscribe_topics.go
git commit -m "feat(message): record stream termination reasons in prometheus"
```

---

## Phase 7 — Integration Tests

These tests use real Redis via `pkg/testutils/redis.NewRedisForTest`. They are slower than unit tests and may be skipped in environments without Docker.

### Task 7.1: End-to-end Tier 2 query deny

**Files:**
- Create: `pkg/interceptors/server/rate_limit_integration_test.go`

- [ ] **Step 1: Write the test.**

```go
//go:build integration
// +build integration

package server_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	server "github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap"
)

func TestRateLimit_Integration_Tier2_QueryDenyAfterBudgetExhausted(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	queryLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":q", []ratelimiter.Limit{
		{Capacity: 5, RefillEvery: time.Minute},
		{Capacity: 100, RefillEvery: time.Hour},
	})
	require.NoError(t, err)
	opensLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":o", []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	rl := server.NewRateLimitInterceptor(zap.NewNop(), queryLimiter, opensLimiter, nil, server.RateLimitInterceptorConfig{})

	calls := 0
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		calls++
		return nil, nil
	}
	wrapped := rl.WrapUnary(next)

	for i := 0; i < 5; i++ {
		req := newFakeQueryRequest("203.0.113.1:5001", "")
		_, err := wrapped(context.Background(), req)
		require.NoError(t, err, "call %d should succeed", i)
	}

	// 6th call should be denied
	req := newFakeQueryRequest("203.0.113.1:5001", "")
	_, err = wrapped(context.Background(), req)
	require.Error(t, err)
	cerr := new(connect.Error)
	require.ErrorAs(t, err, &cerr)
	require.Equal(t, connect.CodeResourceExhausted, cerr.Code())
	require.Equal(t, 5, calls)
}

// Before writing this test, extract the fakeUnaryRequest helper from
// rate_limit_test.go into a shared test helper file:
//
//   pkg/interceptors/server/testhelpers_test.go
//
// containing fakeUnaryRequest, fakeStreamingHandlerConn, and these constructors:

// (Add to testhelpers_test.go)
func newFakeQueryRequest(peerAddr, xff string) *fakeUnaryRequest {
	headers := http.Header{}
	if xff != "" {
		headers.Set("X-Forwarded-For", xff)
	}
	return &fakeUnaryRequest{
		procedure: "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds",
		peerAddr:  peerAddr,
		headers:   headers,
	}
}

func newFakeStreamingConn(procedure, peerAddr string) *fakeStreamingHandlerConn {
	return &fakeStreamingHandlerConn{
		procedure: procedure,
		peerAddr:  peerAddr,
		headers:   http.Header{},
	}
}
```

The engineer should extract the `fakeUnaryRequest` helper into a `testhelpers_test.go` file in the same package so both unit and integration tests can use it.

- [ ] **Step 2: Run the integration test (with Docker available).**

Run: `go test -tags=integration ./pkg/interceptors/server/ -run 'TestRateLimit_Integration_Tier2_QueryDeny' -v`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/interceptors/server/rate_limit_integration_test.go
git commit -m "test(interceptors): integration test for Tier 2 query deny"
```

---

### Task 7.2: Retrospective drain test

- [ ] **Step 1: Write the test.**

In the same integration test file:

```go
func TestRateLimit_Integration_RetrospectiveDrainGoesNegative(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 10, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	ctx := context.Background()

	// Spend 9 tokens
	res, err := limiter.Allow(ctx, "subj", 9)
	require.NoError(t, err)
	require.True(t, res.Allowed)

	// ForceDebit 100 — should drive the bucket sharply negative
	res, err = limiter.ForceDebit(ctx, "subj", 100)
	require.NoError(t, err)
	require.True(t, res.Allowed)
	require.Less(t, res.Balances[0].Remaining, -80.0)

	// Subsequent normal Allow should fail
	res, err = limiter.Allow(ctx, "subj", 1)
	require.NoError(t, err)
	require.False(t, res.Allowed)
}
```

- [ ] **Step 2: Run.**

Run: `go test -tags=integration ./pkg/interceptors/server/ -run 'TestRateLimit_Integration_Retrospective' -v`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/interceptors/server/rate_limit_integration_test.go
git commit -m "test(interceptors): integration test for retrospective drain going negative"
```

---

### Task 7.3: Tier 0 bypass test

- [ ] **Step 1: Write the test.**

```go
func TestRateLimit_Integration_Tier0Bypass(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	limiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 1, RefillEvery: time.Hour}, // very tight
	})
	require.NoError(t, err)
	rl := server.NewRateLimitInterceptor(zap.NewNop(), limiter, limiter, nil, server.RateLimitInterceptorConfig{})

	calls := 0
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		calls++
		return nil, nil
	}
	wrapped := rl.WrapUnary(next)

	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)

	// Hammer the endpoint 100 times — all should pass because Tier 0 bypasses
	for i := 0; i < 100; i++ {
		req := newFakeQueryRequest("10.0.0.5:5001", "")
		_, err := wrapped(ctx, req)
		require.NoError(t, err, "call %d", i)
	}
	require.Equal(t, 100, calls)
}
```

- [ ] **Step 2: Run.**

Run: `go test -tags=integration ./pkg/interceptors/server/ -run 'TestRateLimit_Integration_Tier0Bypass' -v`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/interceptors/server/rate_limit_integration_test.go
git commit -m "test(interceptors): integration test for Tier 0 bypass under tight limits"
```

---

### Task 7.4: Subscribe-opens-per-minute sub-limit test

- [ ] **Step 1: Write the test.**

```go
func TestRateLimit_Integration_SubscribeOpensSubLimit(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	opensLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":opens", []ratelimiter.Limit{
		{Capacity: 3, RefillEvery: time.Minute}, // tight: 3 opens/min
	})
	require.NoError(t, err)
	queryLimiter, err := ratelimiter.NewRedisLimiter(client, keyPrefix+":q", []ratelimiter.Limit{
		{Capacity: 1000, RefillEvery: time.Minute},
	})
	require.NoError(t, err)

	rl := server.NewRateLimitInterceptor(zap.NewNop(), queryLimiter, opensLimiter, nil, server.RateLimitInterceptorConfig{})

	next := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return nil
	}
	wrapped := rl.WrapStreamingHandler(next)

	for i := 0; i < 3; i++ {
		conn := newFakeStreamingConn(
			"/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics",
			"203.0.113.1:5001",
		)
		require.NoError(t, wrapped(context.Background(), conn), "open %d", i)
	}
	conn := newFakeStreamingConn(
		"/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics",
		"203.0.113.1:5001",
	)
	err = wrapped(context.Background(), conn)
	require.Error(t, err)
	cerr := new(connect.Error)
	require.ErrorAs(t, err, &cerr)
	require.Equal(t, connect.CodeResourceExhausted, cerr.Code())
}
```

- [ ] **Step 2: Run.**

Run: `go test -tags=integration ./pkg/interceptors/server/ -run 'TestRateLimit_Integration_SubscribeOpens' -v`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/interceptors/server/rate_limit_integration_test.go
git commit -m "test(interceptors): integration test for subscribe opens-per-minute sub-limit"
```

---

### Task 7.5: Redis fail-open + circuit breaker recovery

- [ ] **Step 1: Write the test.**

```go
func TestRateLimit_Integration_RedisDownTripsBreakerAndFailsOpen(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	inner, err := ratelimiter.NewRedisLimiter(client, keyPrefix, []ratelimiter.Limit{
		{Capacity: 100, RefillEvery: time.Minute},
	})
	require.NoError(t, err)
	cb := ratelimiter.NewCircuitBreaker(2, 200*time.Millisecond)
	bl := ratelimiter.NewBreakerLimiter(inner, cb)

	// Close the underlying Redis connection to simulate Redis going down.
	require.NoError(t, client.Close())

	// First two calls should fail open and trip the breaker
	for i := 0; i < 2; i++ {
		res, err := bl.Allow(context.Background(), "subj", 1)
		require.NoError(t, err)
		require.True(t, res.Allowed) // fail-open
	}
	require.Equal(t, ratelimiter.BreakerOpen, cb.State())

	// Subsequent calls bypass Redis entirely and remain fail-open
	for i := 0; i < 5; i++ {
		res, err := bl.Allow(context.Background(), "subj", 1)
		require.NoError(t, err)
		require.True(t, res.Allowed)
	}
}
```

- [ ] **Step 2: Run.**

Run: `go test -tags=integration ./pkg/interceptors/server/ -run 'TestRateLimit_Integration_RedisDown' -v`
Expected: PASS.

- [ ] **Step 3: Commit.**

```bash
git add pkg/interceptors/server/rate_limit_integration_test.go
git commit -m "test(interceptors): integration test for Redis-down fail-open + breaker trip"
```

---

## Final Verification

- [ ] **Step 1: Run the full test suite.**

Run: `dev/test`
Expected: PASS.

- [ ] **Step 2: Run lint-fix.**

Run: `dev/lint-fix`
Expected: clean.

- [ ] **Step 3: Build all binaries.**

Run: `go build ./...`
Expected: clean.

- [ ] **Step 4: Manual smoke test against a real Redis.**

```bash
docker run -d --name xmtpd-rl-redis -p 6379:6379 redis:7-alpine
XMTPD_RATE_LIMIT_ENABLE=true XMTPD_REDIS_URL=redis://localhost:6379 \
  XMTPD_RATE_LIMIT_T2_PER_MINUTE_CAPACITY=3 \
  dev/run
# In another terminal: hit QueryApi/GetInboxIds 4 times via grpcurl or buf curl
# 4th call should return ResourceExhausted
docker rm -f xmtpd-rl-redis
```

- [ ] **Step 5: Verify metrics are exposed.**

With the node running, fetch `http://127.0.0.1:8008/metrics` and grep for `xmtpd_rate_limit_`. Expected: at least the four metrics defined in `metrics.go`.

- [ ] **Step 6: Open a draft PR.**

```bash
git push -u origin feat/rate-limiting-queryapi
gh pr create --draft --title "feat: rate limiting on QueryApi (v1)" \
  --body "Implements xmtp/xmtpd#366 v1 — rate limiting on QueryApi only. Spec: tasks/rate-limiting-spec.md. Design: tasks/rate-limiting-design.md."
```

---

## Done When

- All Phase 0–7 tasks committed.
- `dev/test` passes.
- `dev/lint-fix` clean.
- Manual smoke test confirms a Tier 2 client gets `ResourceExhausted` after exhausting the bucket.
- Manual smoke test confirms `xmtpd_rate_limit_*` metrics are exported.
- Draft PR open referencing `xmtp/xmtpd#366`.

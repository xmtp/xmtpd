# Rate Limiting Design

## Background

xmtpd nodes are publicly reachable gRPC endpoints. Without rate limiting, a small number of clients can generate disproportionate load — either through high connection counts, wide subscriptions spanning thousands of topics, or abusive write patterns that incur real on-chain settlement costs. This document describes the rate limiting strategy for both the node and the gateway.

This design is intended to align with the service restructuring proposed in [xmtp/proto#327](https://github.com/xmtp/proto/issues/327), which cleanly separates consumers into distinct gRPC services.

---

## Client Tiers

Three distinct client types have different trust levels and usage patterns.

**Tier 0 — Authenticated Nodes**
Node-to-node replication traffic. Nodes identify themselves via a signed JWT verified against the on-chain NodeRegistry contract. These clients must not be rate limited; full sync requires unrestricted throughput. Authentication is mandatory. Tier 0 is **all-or-nothing**: any JWT verification failure (expired, invalid signature, signer not in NodeRegistry) is rejected with `Unauthenticated`. There is no fall-through to Tier 2 — operators are expected to monitor for verification failures and rotate credentials promptly.

**Tier 1 — Fan-out Services**
High-throughput services such as push notification servers. These clients consume the full envelope stream or fan out across large numbers of topics on behalf of many users. Explicit registration is required before elevated limits apply. High limits, but not exempt. **Deferred from v1**: Tier 1 lives with Tier 2 limits on `QueryApi` until the `NotificationApi` design lands.

**Tier 2 — Edge Clients**
Mobile apps, browsers, bots, and other end-user clients. Unauthenticated, identified by source IP. Modeled for single-user usage patterns: a small number of topics, infrequent queries, one or two open subscriptions. Low limits.

### Source IP Identification

Tier 2 keying depends on extracting the real client IP, not the immediate gRPC peer (which is the load balancer or ingress). Source IP is determined by walking `X-Forwarded-For` from the right, peeling exactly one hop per trusted proxy. A proxy is "trusted" only if the immediate peer's IP falls within an operator-configured list of trusted CIDR blocks. The default ships with AWS ALB CIDR ranges; operators running their own ingress configure their own list.

For IPv6, clients are keyed at the **/64 prefix** rather than the full /128 address. This matches typical IPv6 allocation boundaries and prevents trivial bucket evasion by hopping addresses within an allocation. The trade-off — that a household or allocation shares one bucket — is the same compromise IPv4 NAT already imposes, and Tier 2 limits are sized accordingly.

---

## Enforcement Points

Rate limiting applies at two places.

**The Node** enforces limits on all read traffic (queries and subscriptions). Writes to the node are not rate limited — a signed envelope from a gateway is proof of payment and the node must accept it unconditionally. If a node cannot keep up with write volume, that is a performance problem to be solved, not a signal to rate limit. DDoS protection at the load balancer layer applies, but that is distinct from rate limiting.

**The Gateway** enforces limits on all write traffic before messages are signed and submitted. This is the critical enforcement point for payer safety: the gateway incurs real settlement costs on behalf of its clients, so it must reject abusive or runaway write patterns before they generate fees.

---

## Architecture

### Iteration 1 — In-Process

Rate limiting will be implemented directly inside the node process as interceptors in the existing gRPC middleware chain. This is the lowest-infrastructure path: no new binaries, no new deployment topology, and the auth context needed for tier classification is already available in-process.

AWS WAFv2 at the load balancer layer remains in place as a blunt DDoS shield. The in-process interceptors handle all tier-aware, semantically-informed rate limiting.

### Future — Standalone Proxy

The longer-term target architecture is a standalone proxy that sits in front of the node, with the node itself moved to an internal-only endpoint:

```
[public traffic] → [rate-limit proxy] → [node, internal only]
```

This separation is desirable because it allows the proxy to be deployed and scaled independently, simplifies the node process, and makes it possible for node operators to run the proxy in front of an HA pool of nodes. The proxy would share a codebase with the node and reuse the same auth and rate limiting packages, so the migration from in-process to proxy requires moving where the interceptors run, not rewriting them.

### Backend: Redis

Rate limit state lives in Redis. When rate limiting is enforced, **Redis is a hard dependency** — there is no in-memory fallback for single-node deployments. The node must reach Redis at startup or fail fast. This is intentional: rate limiting must be consistent across replicas of the same node, and an in-memory fallback would create silent inconsistency under failover.

Operators must run Redis in an HA configuration. The client supports both Redis Sentinel and Redis Cluster. The token-bucket Lua script is constrained to single-key access per call (with hash-tagged keys where multiple keys are needed in one script) so it is safe in cluster mode.

### Failure Mode: Redis Unreachable at Runtime

Redis going down at runtime is handled with **fail-open + circuit breaker**. The reasoning: rate limiting at the node is defense-in-depth — WAFv2 already sits in front of the node and absorbs volumetric attacks. Taking the read path down because the rate-limit backend is unreachable trades a brief unprotected window for a hard outage, which is the wrong trade.

Behavior:
- Each Redis call has a tight timeout (target: 50ms).
- On timeout or error, the request is allowed through. A metric is incremented.
- After N consecutive failures, the limiter trips a circuit breaker and stops calling Redis entirely for a cooldown period (target: 10s), then probes. This avoids piling synchronous Redis-timeout latency on every request when Redis is hard-down.
- An alert fires on circuit-breaker-open.

Startup behavior is the opposite: if Redis is unreachable at process start, the node refuses to start. This catches misconfiguration loudly rather than running unprotected by accident.

---

## Service-to-Tier Mapping

After the proto restructuring in xmtp/proto#327, the gRPC service path directly identifies the consumer tier. Rate limiting policy is applied based on service, with no request body inspection required for routing decisions.

| Service          | Consumer                  | Policy                                                                                       |
| ---------------- | ------------------------- | -------------------------------------------------------------------------------------------- |
| `ReplicationApi` | Tier 0 — nodes            | Require JWT auth token. No rate limit.                                                       |
| `NotificationApi`| Tier 1 — fan-out services | Deferred for v1. Eventually: require registration, high connection limit.                   |
| `QueryApi`       | Tier 2 — edge clients     | IP-based identity. Topic-weighted rate limits.                                               |
| `PublishApi`     | Gateway only              | No rate limit. Signed envelope = proof of payment. DDoS protection at load balancer only.    |
| `GatewayApi`     | Edge clients → gateway    | Per-JWT (operator gateway) or per-IP (public gateway).                                       |

---

## Fan-out Service Registration

Fan-out services (notification servers) need a registration signal that the node can verify. Registration uses the existing JWT mechanism: the node operator issues a signed JWT to the fan-out service using their node signing key. The node verifies the JWT against the operator's known key — no new token format, no new infrastructure. This scopes the fan-out service to a specific operator's node, which reflects the actual deployment relationship.

On-chain registration is the natural long-term evolution of this: a fan-out service registers a public key and deposit on-chain, making registration permissionless and Sybil-resistant without operator involvement.

This design explicitly underspecifies this problem space.

---

## Rate Limiting Model

### Queries

Queries are request/response. The cost of a query scales sublinearly with the number of topics in the request: a 1000-topic query is more expensive than a 1-topic query, but not 1000× more.

**Model: leaky bucket, cost = `ceil(sqrt(num_topics))`.**

Each client has a token budget that refills over time. Issuing a query deducts `ceil(sqrt(num_topics))` tokens — 1 topic costs 1, 100 topics costs 10, 1000 topics costs 32. This reflects the sublinear relationship between topic count and actual DB work, dominated by partitions hit. The square-root shape also creates the right incentive: a client batching 1000 topics into one query (cost 32) pays far less than the same client issuing 1000 single-topic queries (cost 1000).

The request `limit` parameter is **not** factored into cost. While `limit` bounds the worst-case bandwidth a query can produce, the server CPU and DB work — partition access, index seeks, blob fetches — is dominated by topic count, not row count. Bandwidth-driven attacks are addressed at the WAFv2 / load balancer layer, not by application-level rate limiting.

The budget is shared across all `QueryApi` methods from the same identity: a client cannot evade the bucket by parallelizing across `QueryEnvelopes`, `SubscribeTopics`, `GetInboxIds`, and `GetNewestEnvelope`.

### Subscriptions

Subscriptions are long-lived streams, not request/response. Per-envelope rate limiting is the wrong model here: by the time an envelope reaches a delivery point, the server has already done the work to fetch and decode it. Dropping an envelope at that stage wastes work without preventing cost, and creates cursor consistency problems on reconnect.

The cost of a subscription is not in the ongoing stream — per-envelope dispatch is constant-time regardless of topic count. The cost is at establishment (catch-up queries) and in holding server resources for the lifetime of the connection.

**Model: charge at open, retrospective drain at close.**

Opening a subscription deducts `ceil(sqrt(num_topics))` tokens from the client's shared bucket — the same cost as an equivalent query, since the catch-up phase is essentially a query. While the subscription is open, no further metering happens: zero ongoing Redis traffic, no per-stream goroutine for billing.

When the stream closes — for any reason: graceful client close, client disconnect, server-side cancellation, or process unwind — a deferred drain fires and debits the bucket by `elapsed_minutes / N * drain_amount`, reflecting held server resources for the duration the stream existed. **The bucket is allowed to go negative.** A client who held a stream long enough to owe more than they had earns no immediate punishment, but their next request is rejected until the bucket refills back to positive. This is fair retrospective billing rather than predictive metering.

Two operational notes:

1. The deferred drain runs via `defer` in the subscribe handler, so any unwind path fires it. The only path that misses it is a process crash, which is acceptable since rate limiting is defensive, not accounting.
2. The token-bucket Lua script must be extended with a "force debit" path that bypasses the standard non-negative clamp. This is a targeted change to `pkg/ratelimiter/redis_limiter.go`.

This avoids semaphore state tracking per subscription, eliminates per-stream tickers, and reduces total Redis ops per stream to **two**: one at open, one at close.

Subscription open events are also subject to a separate **connection-rate sub-limit** (e.g. 10 opens per minute per identity), enforced in the same Redis call as the main bucket via the existing multi-limit support in `redis_limiter.go`. This prevents rapid-reconnect patterns from generating repeated catch-up load while bypassing the cheap-streaming phase that retrospective billing assumes.

### Stream Lifetime Management

Retrospective billing only works if streams actually close in bounded time. Three layers of termination ensure this:

1. **gRPC keepalive (transport layer).** The gRPC server is configured with keepalive pings; peers that fail to pong within the deadline are disconnected. This catches dead TCP connections (NAT timeouts, network drops, killed clients).
2. **Application-level idle timeout.** If no envelope has been sent and no client activity has been observed within `idleTimeout`, the server cancels the stream. This catches alive-but-stalled clients (slow consumers with full send buffers, half-open NAT entries that pings still traverse).
3. **Hard maximum stream duration.** A stream older than `maxStreamDuration` (target: 1 hour) is unconditionally cancelled by the server, regardless of activity. This forces periodic re-subscription, bounds worst-case unbilled hold time, and is good operational hygiene independent of rate limiting.

All three paths cancel the stream's context, which fires the deferred drain and debits the bucket for the time the stream was actually held.

### Write Traffic (Gateway)

There are two distinct gateway deployment models with different rate limiting needs.

**App Operator Gateway (authenticated)**

Most deployed gateways issue JWT tokens to their clients. The gateway operator subsidizes all writes on behalf of those clients and controls who receives tokens. Rate limiting is per-JWT: each token represents a known client identity, and limits are enforced at that granularity. The operator's decision to issue a token is an implicit grant of access; the rate limit bounds how much any single client can consume of the operator's subsidy.

The payer is the gateway itself — one payer address covers all clients of that gateway. Per-JWT rate limiting prevents a single misbehaving client from exhausting the operator's on-chain balance.

**Public Gateway (unauthenticated)**

A public gateway accepts writes from anyone without authentication. This is a free tier: the gateway operator absorbs costs for unknown clients. Rate limiting falls back to per-source-IP. Initial limits will be generous to allow legitimate use, but this is the highest-risk surface for abuse and limits will be tightened over time as usage patterns become clear.

The positive-balance check remains in place as a hard gate for both gateway types, independent of rate limiting. A gateway with an empty payer balance rejects all writes regardless of rate limit state.

---

## Handling the Notification Server Case

Push notification servers subscribe to `NotificationApi/SubscribeAllEnvelopes` — the full envelope stream with no topic filter. Topic-count rate limiting does not apply to this service by definition. The rate control for Tier 1 fan-out services on this path is:

- Two concurrent `SubscribeAllEnvelopes` subscriptions per registered identity, to support HA push server deployments. Duplicate delivery across the two connections is handled by idempotency-key deduplication in the push server.
- A connection-rate limit to prevent rapid reconnect.

The registration requirement itself is the primary control. An unregistered client cannot access `NotificationApi` at Tier 1 limits.

---

## What This Does Not Cover

- **Message content validation** — MLS validation happens inside the node and is unchanged.
- **Payer balance enforcement** — The positive-balance check is a separate gate, not a rate limit.
- **Geographic blocking** — Remains at the WAFv2 layer.
- **DDoS absorption** — The load balancer and WAFv2 handle volumetric attacks before traffic reaches the node.

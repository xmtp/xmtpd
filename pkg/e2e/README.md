# xmtpd E2E Test Framework

End-to-end testing framework for xmtpd that orchestrates multi-node clusters using Docker containers, with built-in support for chaos injection and database-level verification.

## Quick Start

```sh
dev/up
dev/e2e
```

Or run directly:

```sh
go run ./cmd/xmtpd-e2e --xmtpd-image ghcr.io/xmtp/xmtpd:latest
go run ./cmd/xmtpd-e2e --test smoke
go run ./cmd/xmtpd-e2e --output json
```

## Architecture

```ascii
┌─────────────────────────────────────────────────────────┐
│  Runner                                                 │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Environment (per-test, isolated Docker network) │   │
│  │                                                  │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐           │   │
│  │  │ node-100│  │ node-200│  │ node-300│  ...      │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘           │   │
│  │       │            │            │                │   │
│  │       └────────────┼────────────┘                │   │
│  │                    │                             │   │
│  │              ┌─────┴─────┐                       │   │
│  │              │ toxiproxy │  (chaos injection)    │   │
│  │              └───────────┘                       │   │
│  │                                                  │   │
│  │  ┌─────────┐  ┌───────┐  ┌───────┐               │   │
│  │  │  anvil  │  │ redis │  │gateway│  ...          │   │
│  │  └─────────┘  └───────┘  └───────┘               │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  Host: Postgres (localhost:8765), MLS validation        │
└─────────────────────────────────────────────────────────┘
```

Each test gets its own Docker network. The runner creates/destroys the full environment per test.

## Package Structure

| Package    | Purpose                                                                        |
| ---------- | ------------------------------------------------------------------------------ |
| `types/`   | Core data structures: `Environment`, `NodeHandle`, `GatewayHandle`, `TestingT` |
| `runner/`  | Test discovery, execution loop, environment lifecycle                          |
| `node/`    | xmtpd node container management                                                |
| `gateway/` | xmtpd gateway container management                                             |
| `chain/`   | Anvil blockchain container (pre-deployed XMTP contracts)                       |
| `chaos/`   | Toxiproxy-based network fault injection                                        |
| `keys/`    | Anvil private key allocation and funding                                       |
| `client/`  | Traffic generation via gRPC envelope publishing                                |
| `observe/` | Database observation (direct Postgres queries for verification)                |
| `cmd/`     | Cobra CLI entry point                                                          |

## API

### Test Interface

Tests implement this interface and are registered in `runner/registry.go`:

```go
type Test interface {
    Name() string
    Description() string
    Run(ctx context.Context, env *Environment) error
}
```

### Environment

The `Environment` is the central orchestrator passed to every test:

```go
// Nodes
env.AddNode(ctx)                        // register on-chain + start container
env.AddNode(ctx, node.WithAPI(true))    // with options
env.Node(100)                           // access by on-chain ID
env.Nodes()                             // all nodes

// Gateways
env.AddGateway(ctx)
env.Gateway(0)                          // access by creation index
env.Gateways()

// Traffic clients
env.NewClient(100)                      // create client targeting node-100
env.Client(100).PublishEnvelopes(ctx, 10)

// Assertions (testify-compatible)
require := require.New(env.T())
require.NoError(err)
```

### NodeHandle

Fluent API for node lifecycle, chaos, and observation:

```go
node := env.Node(100)

// Lifecycle
node.Stop(ctx)
node.Start(ctx)

// On-chain operations
node.AddToCanonicalNetwork(ctx)
node.RemoveFromCanonicalNetwork(ctx)

// Chaos injection (requires toxiproxy)
node.AddLatency(ctx, 500)              // 500ms latency
node.AddBandwidthLimit(ctx, 1024)
node.AddConnectionReset(ctx)
node.AddTimeout(ctx, 5000)
node.RemoveAllToxics(ctx)

// Database observation
count, err := node.GetEnvelopeCount()
node.WaitForEnvelopes(ctx, 10)         // block until count >= 10
node.GetPayerReportCount()
node.GetPayerReportStatusCounts()
```

### Traffic Generation

```go
client := env.Client(100)

// Synchronous: publish N envelopes and return
client.PublishEnvelopes(ctx, 50)

// Asynchronous: generate traffic in background
gen := client.GenerateTraffic(ctx, TrafficOptions{
    Duration:  5 * time.Minute,
    BatchSize: 10,
    Interval:  time.Second,
})
// ... do other things ...
client.Stop()
```

## Design Decisions

### Node IDs Increment by 100

Node IDs are assigned as 100, 200, 300, etc. This avoids collisions with Anvil account indices and provides clear, readable identifiers in logs and test assertions.

### Chaos via Toxiproxy Proxies

When chaos is enabled, each node registers a toxiproxy proxy. The node's **on-chain HTTP address** points to the proxy (`toxiproxy:<port>`) rather than the node directly (`node-100:5050`). This means all inter-node replication traffic flows through toxiproxy, enabling selective fault injection without modifying node code. Toxics are applied per-node, so you can degrade one node's connectivity while leaving others healthy.

### Database-Level Verification

Tests verify correctness by querying each node's Postgres database directly (via `observe/`), not by polling application APIs or parsing logs. This provides deterministic, reliable assertions: "node-200 has exactly 10 envelopes" rather than "the API eventually returned a success response." Polling uses 2-second intervals with context-based timeouts.

### TestingT Adapter

The framework runs outside `go test`, so it provides a `TestingT` type that implements `testify.TestingT`. Tests use standard `require.New(env.T())` / `assert.New(env.T())` patterns. When `require.FailNow()` is called, it panics with a `TestFailedError` that the runner recovers and reports as a test failure.

### Key Management

Anvil ships with 10 pre-funded accounts. The key manager partitions them by role:

- Account 0: Admin (registration, canonical network operations)
- Account 1: Client (payer envelope signing)
- Accounts 2-4: Gateways (3 slots)
- Accounts 5-9: Nodes (5 slots)

When a pool is exhausted, new ECDSA keys are generated and funded via Anvil RPC with 1000 ETH each.

### Container Networking

- **Internal:** Containers communicate via Docker network aliases (`node-100:5050`, `ws://anvil:8545`, `redis://redis:6379/0`)
- **External:** Test clients connect via host-mapped ports (`localhost:<random>`) allocated by testcontainers
- **Host services:** Containers reach host Postgres and MLS validation via `host.docker.internal`

### Node Registration Flow

`AddNode()` performs these steps automatically:

1. Allocate a signer key from the key manager
2. Run `xmtpd-cli:main` container to register the node on-chain via NodeRegistry
3. Create a toxiproxy proxy (if chaos enabled) and set the on-chain address to the proxy
4. Create a per-node Postgres database (`e2e_node_100`)
5. Start the xmtpd container with `XMTPD_CONTRACTS_ENVIRONMENT=anvil` to load embedded contract addresses

### Environment Isolation

Each test gets a fresh environment with its own Docker network, chain, chaos controller, and containers. Cleanup runs in reverse order (gateways → nodes → Redis → chain → chaos → network), capturing the first error while continuing through remaining teardown steps.

## Existing Tests

| Test              | Description                                                                                                 |
| ----------------- | ----------------------------------------------------------------------------------------------------------- |
| `smoke`           | Start 3 nodes + 1 gateway, publish 10 envelopes, verify replication across all nodes                        |
| `chaos-node-down` | Generate traffic, stop a node, restart it, verify it catches up                                             |
| `chaos-latency`   | Inject 500ms latency on a node, verify system recovers after removal                                        |
| `gateway-scale`   | Add/remove gateways dynamically while publishing traffic                                                    |
| `payer-lifecycle` | Long-running test verifying payer reports progress through creation → attestation → submission → settlement |

## Writing a New Test

1. Create a file in `tests/` implementing the `Test` interface
2. Register it in `runner/registry.go`

```go
package tests

type MyTest struct{}

func (t *MyTest) Name() string        { return "my-test" }
func (t *MyTest) Description() string { return "Verifies something important" }

func (t *MyTest) Run(ctx context.Context, env *types.Environment) error {
    require := require.New(env.T())

    // Add infrastructure
    require.NoError(env.AddNode(ctx))
    require.NoError(env.AddNode(ctx))
    require.NoError(env.AddGateway(ctx))

    // Create a client and publish traffic
    require.NoError(env.NewClient(100))
    require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

    // Verify replication
    for _, n := range env.Nodes() {
        require.NoError(n.WaitForEnvelopes(ctx, 10))
    }

    return nil
}
```

## Docker Images

| Image                               | Purpose                                            |
| ----------------------------------- | -------------------------------------------------- |
| `ghcr.io/xmtp/xmtpd:latest`         | Node binary (overridable via `--xmtpd-image`)      |
| `ghcr.io/xmtp/xmtpd-gateway:latest` | Gateway binary (overridable via `--gateway-image`) |
| `ghcr.io/xmtp/contracts:latest`     | Anvil with pre-deployed XMTP contracts             |
| `ghcr.io/xmtp/xmtpd-cli:main`       | CLI for on-chain node registration                 |
| `redis:7-alpine`                    | Redis for gateway caching                          |
| `ghcr.io/shopify/toxiproxy:2.12.0`  | Network fault injection proxy                      |

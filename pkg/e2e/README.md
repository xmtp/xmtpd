# xmtpd E2E Test Framework

- [xmtpd E2E Test Framework](#xmtpd-e2e-test-framework)
  - [Quick Start](#quick-start)
  - [Architecture](#architecture)
  - [Package Structure](#package-structure)
  - [API](#api)
    - [Test Interface](#test-interface)
    - [Environment](#environment)
    - [NodeHandle](#nodehandle)
    - [GatewayHandle](#gatewayhandle)
    - [ClientHandle](#clienthandle)
  - [Design Decisions](#design-decisions)
    - [Node IDs Increment by 100](#node-ids-increment-by-100)
    - [Chaos via Toxiproxy Proxies](#chaos-via-toxiproxy-proxies)
    - [Database-Level Verification](#database-level-verification)
    - [TestingT Adapter](#testingt-adapter)
    - [Key Management](#key-management)
    - [Container Networking](#container-networking)
    - [Node Registration Flow](#node-registration-flow)
    - [Environment Isolation](#environment-isolation)
  - [TODO](#todo)
    - [Client is not a pure client](#client-is-not-a-pure-client)
  - [Existing Tests](#existing-tests)
  - [Writing a New Test](#writing-a-new-test)
  - [Docker Images](#docker-images)

End-to-end testing framework for xmtpd that orchestrates multi-node clusters using Docker containers, with built-in support for chaos injection and database-level verification.

## Quick Start

```sh
dev/up                                                              # start local deps (postgres, etc.)
go run ./cmd/xmtpd-e2e run                                         # run all tests
go run ./cmd/xmtpd-e2e run --xmtpd-image ghcr.io/xmtp/xmtpd:latest
go run ./cmd/xmtpd-e2e run --test smoke
go run ./cmd/xmtpd-e2e run --output-format json
go run ./cmd/xmtpd-e2e list                                        # list available tests
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

| Package    | Purpose                                                                                        |
| ---------- | ---------------------------------------------------------------------------------------------- |
| `types/`   | Core data structures: `Environment`, `NodeHandle`, `GatewayHandle`, `ClientHandle`, `TestingT` |
| `runner/`  | Test discovery, execution loop, environment lifecycle                                          |
| `tests/`   | Test implementations (one file per test)                                                       |
| `node/`    | xmtpd node container management                                                                |
| `gateway/` | xmtpd gateway container management                                                             |
| `chain/`   | Anvil blockchain container (pre-deployed XMTP contracts)                                       |
| `chaos/`   | Toxiproxy-based network fault injection                                                        |
| `keys/`    | Anvil private key allocation and funding                                                       |
| `client/`  | Traffic generation via gRPC envelope publishing                                                |
| `observe/` | Database observation (direct Postgres queries for verification)                                |
| `cmd/`     | Cobra CLI entry point (`run` and `list` subcommands)                                           |

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
// --- Node management ---
env.AddNode(ctx)                              // register on-chain + start container
env.AddNode(ctx, WithNodeImage("img:v1"))     // custom Docker image
env.AddNode(ctx, WithNodeEnvVars(map[string]string{"KEY": "VAL"}))
env.Node(100)                                 // access by on-chain ID (panics if missing)
env.Nodes()                                   // all nodes in creation order

// --- Gateway management ---
env.AddGateway(ctx)                           // start gateway container
env.AddGateway(ctx, WithGatewayAlias("gw"))   // custom alias (default: "gateway-{index}")
env.AddGateway(ctx, WithGatewayImage("img"))  // custom Docker image
env.AddGateway(ctx, WithGatewayEnvVars(map[string]string{"KEY": "VAL"}))
env.Gateway(0)                                // access by creation index (panics if missing)
env.Gateways()                                // all gateways in creation order

// --- Client management ---
env.NewClient(100)                            // create client targeting node-100
env.NewClient(200, WithPayerKey(customKey))   // custom payer identity
env.Client(100)                               // access by nodeID (panics if missing)
env.Clients()                                 // all clients

// --- On-chain network operations ---
env.AddNodeToCanonicalNetwork(ctx, nodeID)
env.RemoveNodeFromCanonicalNetwork(ctx, nodeID)

// --- Rate operations ---
env.UpdateRates(ctx, RateOptions{
    MessageFee:    1000,
    StorageFee:    500,
    CongestionFee: 0,
    TargetRate:    100,
    StartTime:     0,                         // 0 = default (2h from now)
})

// --- Payer / token operations ---
env.FundPayer(ctx, address, amount)           // mint + wrap + approve + deposit (full flow)
env.MintFeeToken(ctx, amount)                 // mint underlying → wrap → approve PayerRegistry
env.DepositPayer(ctx, address, amount)        // deposit into PayerRegistry
env.GetPayerBalance(ctx, address)             // query PayerRegistry balance
env.GetFeeTokenBalance(ctx, address)          // query xUSD balance
env.GetGasBalance(ctx, address)               // query native ETH balance
env.RequestPayerWithdrawal(ctx, privKey, amt) // request withdrawal (payer's key required)
env.CancelPayerWithdrawal(ctx, privKey)       // cancel pending withdrawal
env.FinalizePayerWithdrawal(ctx, privKey, recipient)

// --- Settlement / distribution operations ---
env.SendExcessToFeeDistributor(ctx)           // move excess from PayerRegistry → DistributionManager
env.GetPayerRegistryExcess(ctx)               // query excess balance
env.GetDistributionManagerOwedFees(ctx, nodeID)
env.ClaimFromDistributionManager(ctx, nodeOwnerKey, nodeID, originatorNodeIDs, reportIndices)
env.WithdrawFromDistributionManager(ctx, nodeOwnerKey, nodeID)

// --- Assertions (testify-compatible) ---
require := require.New(env.T())
require.NoError(err)
```

### NodeHandle

Fluent API for node lifecycle, chaos, observation, and on-chain queries:

```go
node := env.Node(100)

// Identity
node.ID()                                     // on-chain nodeID (uint32)
node.Alias()                                  // container alias (e.g. "node-100")
node.Endpoint()                               // host-accessible gRPC address
node.Address()                                // Ethereum address (common.Address)
node.SignerKey()                               // private key hex string
node.DBConnectionString()                     // Postgres connection string

// Lifecycle
node.Stop(ctx)
node.Start(ctx)                               // restart a stopped node

// On-chain operations
node.AddToCanonicalNetwork(ctx)
node.RemoveFromCanonicalNetwork(ctx)

// Chaos injection (requires toxiproxy)
node.AddLatency(ctx, 500)                     // 500ms latency
node.AddBandwidthLimit(ctx, 1024)             // 1024 KB/s
node.AddConnectionReset(ctx, 5000)            // RST after 5s
node.AddTimeout(ctx, 5000)                    // black hole after 5s (0 = indefinite)
node.RemoveAllToxics(ctx)                     // restore normal networking

// Balance queries
node.GetFeeTokenBalance(ctx)                  // xUSD balance for node owner
node.GetGasBalance(ctx)                       // native ETH balance for node owner

// Database observation
node.GetEnvelopeCount(ctx)                    // total envelopes in this node's DB
node.GetVectorClock(ctx)                      // vector clock entries (per originator)
node.GetStagedEnvelopeCount(ctx)              // staged envelopes awaiting processing
node.GetPayerReportCount(ctx)                 // total payer reports
node.GetPayerReportStatusCounts(ctx)          // breakdown by attestation/submission status
node.GetUnsettledUsage(ctx)                   // per-payer spending stats
node.GetSettledPayerReports(ctx)              // settled reports with originator + report index
node.GetNodeInfo(ctx)                         // node_id from the database
node.WaitForEnvelopes(ctx, 10)                // poll until >= 10 envelopes (2s interval)
node.WaitForPayerReports(ctx, checkFn, desc)  // poll until checkFn returns true
```

### GatewayHandle

Fluent API for gateway lifecycle, chaos, and payer operations:

```go
gw := env.Gateway(0)

// Identity
gw.Index()                                    // creation index (0, 1, 2, ...)
gw.Alias()                                    // container alias (e.g. "gateway-0")
gw.Endpoint()                                 // internal network address
gw.Address()                                  // Ethereum payer address (common.Address)

// Lifecycle
gw.Stop(ctx)

// Chaos injection
gw.AddLatency(ctx, 500)
gw.AddBandwidthLimit(ctx, 1024)
gw.AddConnectionReset(ctx, 5000)
gw.AddTimeout(ctx, 5000)
gw.RemoveAllToxics(ctx)

// Balance queries
gw.GetPayerBalance(ctx)                       // PayerRegistry balance
gw.GetFeeTokenBalance(ctx)                    // xUSD balance
gw.GetGasBalance(ctx)                         // native ETH balance

// Payer operations
gw.Deposit(ctx, amount)                       // full flow: mint → wrap → approve → deposit
gw.RequestWithdrawal(ctx, amount)
gw.CancelWithdrawal(ctx)
gw.FinalizeWithdrawal(ctx, recipientAddress)
```

### ClientHandle

Traffic generation client with payer convenience methods:

```go
client := env.Client(100)

// Traffic
client.PublishEnvelopes(ctx, 50)              // synchronous: publish N envelopes
gen := client.GenerateTraffic(ctx, client.TrafficOptions{
    Duration:  5 * time.Minute,
    BatchSize: 10,
})
client.Stop()                                 // stop background traffic
client.NodeID()                               // target node's on-chain ID

// Identity
client.Address()                              // Ethereum payer address

// Balance queries
client.GetPayerBalance(ctx)                   // PayerRegistry balance
client.GetFeeTokenBalance(ctx)                // xUSD balance
client.GetGasBalance(ctx)                     // native ETH balance

// Payer operations (direct-publish mode only; in production the gateway is the payer)
client.Deposit(ctx, amount)                   // full flow: mint → wrap → approve → deposit
client.RequestWithdrawal(ctx, amount)
client.CancelWithdrawal(ctx)
client.FinalizeWithdrawal(ctx, recipientAddress)
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
2. Run `xmtpd-cli` container to register the node on-chain via NodeRegistry
3. Create a toxiproxy proxy (if chaos enabled) and set the on-chain address to the proxy
4. Reset the per-node Postgres database (`e2e_node_100`) — drop if exists to ensure clean state
5. Start the xmtpd container with `XMTPD_CONTRACTS_ENVIRONMENT=anvil` to load embedded contract addresses (the node creates its DB via migrations on startup)

### Environment Isolation

Each test gets a fresh environment with its own Docker network, chain, chaos controller, and containers. Cleanup runs in reverse order (gateways → nodes → Redis → chain → chaos → network), capturing the first error while continuing through remaining teardown steps.

## TODO

### Client is not a pure client

The current e2e client publishes directly to nodes and signs its own payer
envelopes — making it act as both client AND payer. In production, clients
send messages through a gateway, and the gateway is the actual payer (see
`XMTPD_PAYER_PRIVATE_KEY` in the gateway container). We need a client mode
that publishes through a gateway rather than directly to a node, so we can
test the real client → gateway → node flow.

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
| `ghcr.io/xmtp/xmtpd-cli:latest`     | CLI for on-chain node registration                 |
| `redis:7-alpine`                    | Redis for gateway caching                          |
| `ghcr.io/shopify/toxiproxy:2.12.0`  | Network fault injection proxy                      |

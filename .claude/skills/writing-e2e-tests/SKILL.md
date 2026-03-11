---
name: writing-e2e-tests
description: >-
  Use when creating, modifying, or reviewing E2E tests in pkg/e2e/tests/.
  Triggers on "e2e test", "add e2e test", "new e2e test", or when working
  with files under pkg/e2e/tests/.
---

# Writing E2E Tests

## Overview

E2E tests live in `pkg/e2e/tests/` and implement the `types.Test` interface.
Each test gets a fully isolated `Environment` with its own Docker network,
Anvil chain, Redis, toxiproxy, and containers. The runner handles all
setup/teardown.

## Creating a New Test

### Step 1: Create the test file

Create `pkg/e2e/tests/<snake_case_name>.go` implementing `types.Test`:

```go
package tests

import (
    "context"
    "time"

    "github.com/stretchr/testify/require"
    "github.com/xmtp/xmtpd/pkg/e2e/types"
)

type MyTest struct{}

func NewMyTest() *MyTest { return &MyTest{} }

func (t *MyTest) Name() string        { return "my-test" }
func (t *MyTest) Description() string { return "One-line description of what this verifies" }

func (t *MyTest) Run(ctx context.Context, env *types.Environment) error {
    require := require.New(env.T())

    // 1. Set up infrastructure
    require.NoError(env.AddNode(ctx))
    require.NoError(env.AddNode(ctx))
    require.NoError(env.AddGateway(ctx))

    // 2. Create clients and generate traffic
    require.NoError(env.NewClient(100))
    require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

    // 3. Assert with timeouts
    checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    for _, n := range env.Nodes() {
        require.NoError(n.WaitForEnvelopes(checkCtx, 1))
    }

    return nil
}

var _ types.Test = (*MyTest)(nil)
```

### Step 2: Register the test

Add the constructor to `pkg/e2e/runner/registry.go`:

```go
func AllTests() []types.Test {
    return []types.Test{
        // ... existing tests ...
        tests.NewMyTest(),
    }
}
```

### Step 3: Lint and verify

Run `dev/lint-fix` to ensure the code passes all linters.

## Conventions

### File and naming

- File name: `snake_case.go` matching the test name (e.g. `chaos_node_down.go` for "chaos-node-down")
- Struct name: `PascalCase` + `Test` suffix (e.g. `ChaosNodeDownTest`)
- Constructor: `NewXxxTest()` returning a pointer
- `Name()` returns a kebab-case identifier (e.g. `"chaos-node-down"`)
- `Description()` returns a short, human-readable sentence (no period)
- Add compile-time interface check: `var _ types.Test = (*MyTest)(nil)`

### Test structure pattern

Every test follows this structure:

1. **Infrastructure setup** -- `env.AddNode(ctx)`, `env.AddGateway(ctx)`
2. **Client creation** -- `env.NewClient(nodeID)`
3. **Action** -- publish traffic, inject chaos, stop/start nodes, etc.
4. **Assertion with timeout** -- always wrap waits in `context.WithTimeout`
5. **Cleanup** (optional) -- stop traffic generators, remove toxics

### Assertions

- Use `require := require.New(env.T())` for fatal assertions (stops test on failure)
- Use `assert := assert.New(env.T())` for non-fatal checks (continues after failure)
- Prefer `require` for setup steps (AddNode, AddGateway, NewClient)
- Include descriptive messages: `require.NoError(err, "failed to do X for node %d", nodeID)`

### Timeouts

- Always wrap polling waits in `context.WithTimeout` -- never wait indefinitely
- Short operations (envelope replication): 30-60 seconds
- Payer report creation (depends on worker scheduling): up to 65 minutes
- Post-generator operations (attestation, submission): 10-15 minutes
- Use `ctx` (the test's parent context) as the parent for all timeouts

### Logging

- Use `env.Logger` for structured logging (zap)
- Log phase transitions: `env.Logger.Info("phase N: description")`
- Log important state: balances, counts, node IDs
- Use zap fields, not string formatting: `zap.Uint32("node_id", id)`

## Available API

See [pkg/e2e/README.md](../../../pkg/e2e/README.md) for the full API reference.

### Quick reference

```go
// Node management
env.AddNode(ctx)                              // register on-chain + start container
env.AddNode(ctx, types.WithAlias("name"))     // custom alias
env.AddNode(ctx, types.WithNodeImage("img"))  // custom image
env.AddNode(ctx, types.WithNodeEnvVars(m))    // extra env vars
env.Node(100)                                 // access by on-chain ID
env.Nodes()                                   // all nodes

// Gateway management
env.AddGateway(ctx)
env.AddGateway(ctx, types.WithGatewayAlias("name"))
env.Gateway(0)                                // access by creation index
env.Gateways()

// Client management
env.NewClient(100)                            // create client for node 100
env.NewClient(200, types.WithPayerKey(key))   // custom payer key
env.Client(100)                               // access by node ID
env.Clients()

// NodeHandle -- lifecycle
node.Stop(ctx)
node.Start(ctx)

// NodeHandle -- on-chain
node.AddToCanonicalNetwork(ctx)
node.RemoveFromCanonicalNetwork(ctx)

// NodeHandle -- chaos (requires toxiproxy)
node.AddLatency(ctx, ms)
node.AddBandwidthLimit(ctx, kbps)
node.AddConnectionReset(ctx, timeoutMs)
node.AddTimeout(ctx, timeoutMs)               // 0 = black hole
node.RemoveAllToxics(ctx)

// NodeHandle -- observation
node.GetEnvelopeCount(ctx)
node.GetVectorClock(ctx)
node.GetStagedEnvelopeCount(ctx)
node.GetPayerReportCount(ctx)
node.GetPayerReportStatusCounts(ctx)
node.GetUnsettledUsage(ctx)
node.GetSettledPayerReports(ctx)
node.WaitForEnvelopes(ctx, minCount)
node.WaitForPayerReports(ctx, checkFn, description)

// NodeHandle -- balances
node.GetFeeTokenBalance(ctx)
node.GetGasBalance(ctx)

// GatewayHandle
gw.Stop(ctx)
gw.Deposit(ctx, amount)
gw.RequestWithdrawal(ctx, amount)
gw.GetPayerBalance(ctx)

// ClientHandle -- traffic
client.PublishEnvelopes(ctx, count)
client.GenerateTraffic(ctx, client.TrafficOptions{BatchSize: 10, Duration: 5*time.Minute})
client.Stop()

// ClientHandle -- payer ops
client.Deposit(ctx, amount)
client.GetPayerBalance(ctx)
client.RequestWithdrawal(ctx, amount)

// Environment -- on-chain operations
env.UpdateRates(ctx, types.RateOptions{...})
env.FundPayer(ctx, address, amount)
env.MintFeeToken(ctx, amount)
env.DepositPayer(ctx, address, amount)
env.GetPayerBalance(ctx, address)
env.GetFeeTokenBalance(ctx, address)
env.GetGasBalance(ctx, address)
env.SendExcessToFeeDistributor(ctx)
env.GetPayerRegistryExcess(ctx)
env.ClaimFromDistributionManager(ctx, ownerKey, nodeID, originatorIDs, indices)
env.WithdrawFromDistributionManager(ctx, ownerKey, nodeID)
```

## Common Patterns

### Verify replication across all nodes

```go
checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
defer cancel()
for _, n := range env.Nodes() {
    require.NoError(n.WaitForEnvelopes(checkCtx, expectedCount))
}
```

### Chaos: stop and restart a node

```go
node := env.Node(200)
require.NoError(node.Stop(ctx))
// ... do something while node is down ...
require.NoError(node.Start(ctx))
// ... verify recovery ...
```

### Chaos: inject and remove latency

```go
require.NoError(node.AddLatency(ctx, 500))
// ... traffic while latency is active ...
require.NoError(node.RemoveAllToxics(ctx))
// ... verify recovery ...
```

### Background traffic with cleanup

```go
env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
    BatchSize: 10,
    Duration:  5 * time.Minute,
})
defer env.Client(100).Stop()
```

### Fund a payer and verify deposit

```go
amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1e18
require.NoError(env.Client(100).Deposit(ctx, amount))
balance, err := env.Client(100).GetPayerBalance(ctx)
require.NoError(err)
require.Positive(balance.Sign())
```

### Wait for payer report status

```go
require.NoError(node.WaitForPayerReports(
    ctx,
    func(c *observe.PayerReportStatusCounts) bool {
        return c.SubmissionSettled > 0
    },
    "at least 1 settled payer report",
))
```

## Node IDs

Nodes are assigned IDs 100, 200, 300, etc. in creation order. After calling
`env.AddNode(ctx)` three times, you have nodes 100, 200, and 300. Use these
IDs consistently when creating clients: `env.NewClient(100)` creates a client
that publishes to node 100.

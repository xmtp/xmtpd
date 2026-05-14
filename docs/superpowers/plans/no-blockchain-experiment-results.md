# XMTPD Without Chain — Experiment Results

## Summary

Successfully implemented and tested a `--no-blockchain` mode for XMTPD that routes
commit messages and identity updates to dedicated nodes instead of through the
Arbitrum L2 blockchain. The experiment demonstrates that the blockchain dependency
can be removed from the critical message path with minimal code changes.

## Architecture

### Production (with chain)
```
Client → Gateway/Payer → shouldSendToBlockchain?
  YES → PublishGroupMessage/PublishIdentityUpdate → Arbitrum L2 → indexer → DB
  NO  → NodeSelector → Node → DB
```

### Experiment (no chain)
```
Client → Gateway/Payer → shouldSendToBlockchain?
  YES (but NoBlockchain=true) → dedicated node ID → Node → DB
  NO  → NodeSelector → Node → DB (unchanged)
```

## Code Changes

| File | Change |
|------|--------|
| `pkg/config/options.go` | Added `NoBlockchain`, `CommitNodeID`, `IdentityNodeID` to `PayerOptions` and `NoBlockchain` to `APIOptions` |
| `pkg/api/payer/config.go` | Added `NoBlockchain`, `CommitNodeID`, `IdentityNodeID` to payer `Config` + `WithNoBlockchain` option |
| `pkg/api/payer/service.go` | Modified `groupEnvelopes()` to route blockchain-bound messages to nodes; fixed `determineRetentionPolicy` to handle identity updates and commits |
| `pkg/api/message/service.go` | Made 3 rejection points conditional: identity update rejection, commit/proposal rejection, DependsOn nodeID≥100 check |
| `pkg/gateway/builder.go` | Wire `WithNoBlockchain` option from config to payer service |
| `pkg/testutils/api/api.go` | Added `WithTestNoBlockchain()` option |

**Total: ~100 lines of production code changed, ~200 lines of tests added.**

## Performance Results

| Path | Avg Latency | Notes |
|------|------------|-------|
| No-blockchain commit (via node) | **~13ms** | Same as regular message |
| Regular message (baseline) | **~13ms** | Unchanged |
| Blockchain commit (production) | **2,000-15,000ms** | Arbitrum L2 finality |

**Improvement: 150-1000x faster for commit/proposal messages.**

## Failure Mode Analysis

### 1. Commit Node Failure
- **Scenario**: Node designated for commits (e.g., node 100) goes down
- **Impact**: All commits fail until node recovers
- **Mitigation**: Existing retry + banlist mechanism attempts failover. Could add
  fallback commit node ID config.
- **Comparison**: With blockchain, Arbitrum L2 availability is the dependency instead.
  Node failures are local and recoverable; L2 outages are external.

### 2. Identity Node Failure
- **Scenario**: Node designated for identity updates (e.g., node 200) goes down
- **Impact**: Identity updates fail, new users cannot register
- **Mitigation**: Same retry/banlist mechanism. Could configure multiple identity nodes.
- **Risk Level**: HIGH — identity updates are critical for onboarding.

### 3. Split Brain / Ordering
- **Scenario**: Without blockchain ordering guarantee, two commits for same group
  could arrive at different nodes simultaneously
- **Impact**: In production, blockchain provides total ordering for commits.
  Without it, ordering depends on the node's local sequencing.
- **Mitigation**: Stable hashing ensures same group always routes to same commit node.
  Single-node-per-group provides sequencing. BUT if commit node changes (failover),
  ordering gap is possible.
- **Risk Level**: MEDIUM — acceptable for experiment, needs resolution for production.

### 4. DependsOn Validation Relaxed
- **Scenario**: `nodeID >= 100` check is relaxed in no-blockchain mode
- **Impact**: Messages can now declare dependencies on node-originated commits
  (previously only blockchain-originated commits with ID 0 were allowed)
- **Mitigation**: Logical — if commits come from nodes, DependsOn must accept node IDs
- **Risk Level**: LOW — correct behavior for no-blockchain mode

### 5. Replay/Dedup
- **Scenario**: Without blockchain's inherent dedup, same commit could be published twice
- **Impact**: Duplicate commits stored in node DB
- **Mitigation**: Node's existing dedup logic (topic+sequence_id uniqueness) handles this
- **Risk Level**: LOW

### 6. Settlement/Audit Trail
- **Scenario**: No blockchain record of commits means no on-chain audit trail
- **Impact**: Cannot verify message ordering disputes on-chain
- **Mitigation**: Node DB provides audit trail. For experiment purposes, acceptable.
- **Risk Level**: LOW for experiment, HIGH for production

## Configuration

```bash
# Payer (gateway) flags
--payer.no-blockchain          # Enable no-blockchain mode
--payer.commit-node-id=100     # Route commits to node 100
--payer.identity-node-id=200   # Route identity updates to node 200

# Node (replication) flags
--api.no-blockchain            # Accept commits and identity updates from payers

# Environment variables
XMTPD_PAYER_NO_BLOCKCHAIN=true
XMTPD_PAYER_COMMIT_NODE_ID=100
XMTPD_PAYER_IDENTITY_NODE_ID=200
XMTPD_API_NO_BLOCKCHAIN=true
```

## Dev Environment

```bash
# Start 3-node cluster with triple profile
dev/docker/up triple

# Run all 3 nodes with no-blockchain
dev/run-no-chain 1  # Node 100 on :5050
dev/run-no-chain 2  # Node 200 on :5051
dev/run-no-chain 3  # Node 300 on :5052
```

## Test Coverage

- `TestNoBlockchain_IdentityUpdateRoutedToNode` — identity updates bypass blockchain
- `TestNoBlockchain_CommitRoutedToNode` — commits route to dedicated node, stored successfully
- `TestNoBlockchain_RegularMessageStillUsesNodeSelector` — regular messages unaffected
- `TestNoBlockchain_LatencyComparison` — 50-iteration latency benchmark

## Next Steps

1. **Gateway integration test**: Full E2E with gateway → payer → node flow
2. **xdbg validation**: Run xdbg debug tool against no-blockchain cluster
3. **Multi-group ordering test**: Verify ordering across multiple groups
4. **Failover testing**: Kill commit/identity nodes, verify recovery
5. **Production readiness**: If experiment succeeds, design permanent architecture

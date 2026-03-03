# Batched Publish Worker Design

## Problem

The publish worker processes envelopes sequentially, issuing ~4 DB round trips per envelope. For a batch of 100 envelopes, this means ~400 DB round trips. The batch insert infrastructure (`insert_gateway_envelope_batch`) already exists but isn't used by the publish worker.

## Design

Restructure the publish worker into two phases: CPU prep (per-envelope) then batch DB ops (3 round trips total).

### Phase 1: CPU Prep (per-envelope, sequential)

For each staged envelope in the batch:

1. Parse payer envelope and topic, determine `isReserved`
2. Calculate base fee (CPU-only, no DB)
3. Calculate congestion fee using DB state + in-memory running count
4. Sign envelope with calculated fees
5. Recover payer address from signed envelope

**Congestion fee accuracy:** Query `GetRecentOriginatorCongestion` once at the start of the batch. Maintain an in-memory running count of messages processed so far (per minute bucket). Each envelope's congestion fee uses `DB congestion + running count`, producing identical fees to sequential processing regardless of batch boundaries.

**Error model:** If any envelope fails CPU prep, return false and retry the entire batch (all-or-nothing). This preserves the current ordering guarantee.

### Phase 2: Batch DB Ops (3 round trips)

1. **`BulkFindOrCreatePayers`** — new sqlc query using `unnest` + `ON CONFLICT` to upsert all unique payer addresses in one call. Returns `(address, id)` pairs.
2. **`InsertGatewayEnvelopeBatchV2`** — extended SQL function that inserts gateway envelopes, increments unsettled usage, and increments originator congestion. Accepts a `p_is_reserved` boolean array: reserved topics skip both usage and congestion increments.
3. **`BulkDeleteStagedOriginatorEnvelopes`** — new sqlc query: `DELETE WHERE id = ANY(@ids)`. Replaces the old single-row `DeleteStagedOriginatorEnvelope`.

If the batch insert fails, return false and retry the entire batch. CPU prep results are reusable on retry.

### Reserved Topics (Unified Path)

All envelopes go through the same batch pipeline. Reserved topics are marked with `isReserved=true` per row. The SQL function uses this flag to:

- Skip `unsettled_usage` increment (no zero-value rows created)
- Skip `originator_congestion` increment

This matches today's behavior exactly while eliminating the branching code path.

## New SQL Artifacts

### Migration: `insert_gateway_envelope_batch_v2`

New SQL function (does not modify existing `insert_gateway_envelope_batch`). Adds:

- `p_is_reserved boolean[]` parameter
- `c_prep` / `c` CTEs for congestion increment (filtered by `WHERE NOT is_reserved`)
- Usage CTE also filtered by `WHERE NOT is_reserved` (instead of relying on `payer_id IS NOT NULL`)

### sqlc Query: `BulkFindOrCreatePayers`

```sql
-- name: BulkFindOrCreatePayers :many
WITH input AS (
    SELECT address FROM unnest(@addresses::TEXT[]) AS t(address)
),
ins AS (
    INSERT INTO payers(address)
    SELECT address FROM input
    ON CONFLICT (address) DO NOTHING
    RETURNING id, address
)
SELECT address, id
FROM ins
UNION ALL
SELECT i.address, p.id
FROM input i
JOIN payers p ON p.address = i.address
WHERE i.address NOT IN (SELECT address FROM ins);
```

### sqlc Query: `BulkDeleteStagedOriginatorEnvelopes`

```sql
-- name: BulkDeleteStagedOriginatorEnvelopes :execrows
DELETE FROM staged_originator_envelopes WHERE id = ANY(@ids::BIGINT[]);
```

### sqlc Query: `InsertGatewayEnvelopeBatchV2`

Points to `insert_gateway_envelope_batch_v2` function with the added `p_is_reserved` parameter.

## Go Changes

| File | Change |
|------|--------|
| `pkg/db/types/gateway_envelope_batch.go` | Add `IsReserved` field to `GatewayEnvelopeRow`, update `ToParams` for V2 |
| `pkg/db/gateway_envelope_batch.go` | Add V2 wrapper calling new SQL function |
| `pkg/api/message/publish_worker.go` | Replace per-envelope `publishStagedEnvelope` loop with `publishBatch` two-phase method |
| `pkg/fees/calculator.go` | Add `additionalMessages` parameter to `CalculateCongestionFee` |
| `pkg/db/bench/hot_path_bench_test.go` | Update to use bulk delete |

## Removed

- `DeleteStagedOriginatorEnvelope` sqlc query — replaced by `BulkDeleteStagedOriginatorEnvelopes`

## Concurrency Safety

The system uses optimistic concurrency with idempotency guards, not pessimistic locking. This extends cleanly to batch operations:

- `InsertGatewayEnvelope` uses `ON CONFLICT DO NOTHING`. The batch SQL function's CTE chain flows from actually-inserted rows (`m` CTE), so a losing worker increments zero usage and zero congestion.
- Two workers computing fees for the same envelopes produce identical results (same base congestion, same in-memory tracking logic). The first to insert wins; the other's work is discarded.
- `BulkDeleteStagedOriginatorEnvelopes` naturally handles partial deletes — already-deleted rows are simply not matched.
- No additional locking needed.

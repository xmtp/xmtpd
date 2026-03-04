# CLAUDE.md — pkg/db

## Overview

Database layer using **sqlc** for query generation, **pgx v5** driver, and **golang-migrate** for schema migrations. Supports read/write replica routing via `Handler`.

## File Layout

```
pkg/db/
  sqlc/           # source .sql query files (input to sqlc)
  queries/        # generated Go code (output of sqlc) — do not edit
  types/          # custom Go types (e.g. GatewayEnvelopeBatch)
  db.go           # Handler (read/write routing)
  pgx.go          # connection management, NewNamespacedDB
  tx.go           # transaction helpers
pkg/db/migrations/   # SQL migration files (up/down)
sqlc.yaml         # sqlc config (project root)
```

## Schema Summary

| Table                         | Purpose                                                | Key Columns                                                                                                   |
| ----------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------- |
| `gateway_envelopes_meta`      | Envelope metadata (hot path, partitioned)              | `originator_node_id`, `originator_sequence_id`, `topic`, `payer_id`, `expiry`, `gateway_time`                 |
| `gateway_envelope_blobs`      | Envelope payloads (cold path, partitioned)             | `originator_node_id`, `originator_sequence_id`, `originator_envelope`                                         |
| `gateway_envelopes_view`      | View joining meta + blobs                              | —                                                                                                             |
| `gateway_envelopes_latest`    | Latest sequence ID per originator (trigger-maintained) | `originator_node_id`, `originator_sequence_id`, `gateway_time`                                                |
| `staged_originator_envelopes` | Publish queue before ordering                          | `id` (serial), `topic`, `payer_envelope`                                                                      |
| `node_info`                   | Local node identity (singleton)                        | `node_id`, `public_key`                                                                                       |
| `address_log`                 | Address → inbox_id mapping                             | `address`, `inbox_id`, `association_sequence_id`, `revocation_sequence_id`                                    |
| `payers`                      | Payer addresses                                        | `id`, `address`                                                                                               |
| `unsettled_usage`             | Per-payer per-originator usage tracking                | `payer_id`, `originator_id`, `minutes_since_epoch`, `spend_picodollars`                                       |
| `payer_reports`               | Cross-node payer settlement reports                    | `id`, `originator_node_id`, `start_sequence_id`, `end_sequence_id`, `submission_status`, `attestation_status` |
| `payer_report_attestations`   | Attestation signatures on reports                      | `payer_report_id`, `node_id`, `signature`                                                                     |
| `payer_ledger_events`         | Deposit/withdrawal/settlement events                   | `event_id`, `payer_id`, `amount_picodollars`, `event_type`                                                    |
| `blockchain_messages`         | Links envelopes to blockchain blocks                   | `block_number`, `block_hash`, `originator_node_id`, `originator_sequence_id`, `is_canonical`                  |
| `latest_block`                | Last indexed block per contract                        | `contract_address`, `block_number`                                                                            |
| `originator_congestion`       | Per-originator message rate tracking                   | `originator_id`, `minutes_since_epoch`, `num_messages`                                                        |
| `nonce_table`                 | Transaction nonce management                           | `nonce`                                                                                                       |
| `migration_tracker`           | Data migration progress tracking                       | `source_table`, `last_migrated_id`                                                                            |
| `migration_dead_letter_box`   | Failed migration records for retry                     | `source_table`, `sequence_id`, `payload`, `reason`, `retryable`                                               |

## Partitioning

`gateway_envelopes_meta` and `gateway_envelope_blobs` use two-level partitioning:

1. **Level 1 — LIST** by `originator_node_id`: one child table per originator (e.g. `gateway_envelopes_meta_o100`)
2. **Level 2 — RANGE** by `originator_sequence_id`: 1M-row bands (e.g. `gateway_envelopes_meta_o100_s0_1000000`)

**Dynamic creation:** `ensure_gateway_parts_v2(originator_node_id, sequence_id, band_width DEFAULT 1000000)` creates both levels if missing. Called via savepoint pattern on "no partition" errors:

```sql
-- See pkg/db/sqlc/partitions.sql
SAVEPOINT sp_part;
-- attempt insert, on partition error:
ROLLBACK TO SAVEPOINT sp_part;
-- call EnsureGatewayParts, then retry
RELEASE SAVEPOINT sp_part;
```

## Hot/Cold Path

- **Meta table** (`gateway_envelopes_meta`): small rows, heavily indexed, used in all query filters
- **Blob table** (`gateway_envelope_blobs`): large payloads, joined only when envelope content is needed
- **View** (`gateway_envelopes_view`): convenience join of meta + blobs, FK with ON DELETE CASCADE

## Key Query Patterns

**Advisory locks** (`sqlc/advisory_locks.sql`): `pg_advisory_xact_lock` / `pg_try_advisory_xact_lock` for serializing operations (e.g. staged envelope insertion, dead letter box)

**Upserts** (`ON CONFLICT`): used throughout — `InsertNodeInfo`, `IncrementUnsettledUsage`, `FindOrCreatePayer`, `InsertOrIgnorePayerReport`

**Batch operations** (`unnest()`): `InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage` accepts parallel arrays for bulk inserts

**Time bucketing**: `minutes_since_epoch` column for usage aggregation; payer reports use minute-level granularity

**LATERAL per (topic, originator)** (`sqlc/envelopes_v2.sql`): `SelectGatewayEnvelopesByTopics` uses `CROSS JOIN LATERAL` for per-(topic, originator) index probes with per-originator blob join for cache locality. Callers must include all originators in cursor arrays (use `FillMissingOriginators` on the `VectorClock` before `SetVectorClockByTopics`).

**LATERAL per originator** (`sqlc/envelopes_v2.sql`): `SelectGatewayEnvelopesByOriginators` uses `CROSS JOIN LATERAL` for per-originator cursor-based pagination

**Trigger-maintained latest**: `gateway_envelopes_latest` auto-updated via `AFTER INSERT` trigger on `gateway_envelopes_meta`

## Indexes

On `gateway_envelopes_meta`:

- `gem_topic_orig_seq_idx` — `(topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time)` — covering index for V3b LATERAL
- `gem_topic_time_desc_idx` — `(topic, gateway_time DESC) INCLUDE (originator_node_id, originator_sequence_id)`
- `gem_expiry_idx` — `(expiry) INCLUDE (...) WHERE expiry IS NOT NULL`

Other notable indexes:

- `blockchain_messages`: `(block_number, is_canonical)`
- `unsettled_usage`: `(originator_id, minutes_since_epoch DESC)`
- `payer_reports`: `(submission_status, created_at)`, `(attestation_status, created_at)`
- `payer_ledger_events`: `(payer_id)`

## Go Utilities

**Handler** (`db.go`): Routes queries to write or read replica. `NewDBHandler(db, WithReadReplica(readDB))`. Methods: `Write()`, `Read()`, `WriteQuery()`, `ReadQuery()`.

**Transaction helpers** (`tx.go`):

- `RunInTx(ctx, db, opts, func(ctx, *queries.Queries) error)` — run in transaction with auto-rollback
- `RunInTxWithResult[T](ctx, db, opts, func(ctx, *queries.Queries) (T, error))` — same with return value
- `RunInTxRaw(ctx, db, opts, func(ctx, *sql.Tx) error)` — raw transaction access

**NewNamespacedDB** (`pgx.go`): Creates database if not exists, runs migrations, returns `*sql.DB`. Used for test isolation.

**ConnectToDB** (`pgx.go`): Connects to existing database without creation or migration.

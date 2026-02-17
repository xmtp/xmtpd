---
name: writing-queries
description: >-
  Writes and modifies sqlc queries for the xmtpd PostgreSQL database.
  Use when creating, editing, or reviewing .sql query files in pkg/db/sqlc/,
  when the user mentions sqlc, database queries, or adding new database operations.
---

# Writing sqlc Queries

## Project Setup

- **sqlc config:** `sqlc.yaml` at repo root
- **Query files:** `pkg/db/sqlc/*.sql`
- **Generated Go:** `pkg/db/queries/`
- **Regenerate:** `dev/gen/all`

## File Organization

Queries are grouped by domain into separate `.sql` files (e.g., `envelopes.sql`, `payer_reports.sql`, `congestion.sql`). New queries go in the file matching their domain; create a new file only if no existing file fits.

## Query Annotations

Four types used in this project:

- `:one` — single row (InsertGatewayEnvelope, GetLatestSequenceId)
- `:many` — multiple rows (SelectVectorClock, FetchPayerReports)
- `:exec` — no return (InsertPayerLedgerEvent, SetReportAttestationStatus)
- `:execrows` — returns affected count (DeleteStagedOriginatorEnvelope, InsertAddressLog)

## Query Naming

Format: `-- name: VerbNoun :type` in PascalCase.

Verbs observed: Select, Get, Insert, Delete, Update, Set, Find, Fetch, Build, Count, Clear, Increment, Revoke.

Examples: `FindOrCreatePayer`, `DeleteExpiredEnvelopesBatch`, `IncrementOriginatorCongestion`.

## Parameter Styles

Three styles, in order of preference:

1. **`@named_param`** — preferred for most queries (vast majority of codebase)
2. **`sqlc.narg('name')`** — for nullable/optional filter parameters (generates `sql.NullXXX` in Go)
3. **`sqlc.arg(name)`** — when used alongside positional `$1` in same query, or when explicit naming is needed with type casts

Positional `$1` — only in `configuration.sql` and mixed with `sqlc.arg()` in `congestion.sql`.

## SQL Formatting Conventions

- **Tab indentation** throughout
- **UPPERCASE:** keywords (SELECT, FROM, WHERE, JOIN, INSERT, ON CONFLICT, DO UPDATE, RETURNING, ORDER BY, LIMIT, WITH, VALUES, SET, DELETE, UPDATE, USING, UNION ALL, GROUP BY, HAVING), functions (COALESCE, COUNT, MAX, SUM, EXTRACT, GREATEST, ANY, EXISTS, NULLIF, NOW), types in casts (::BIGINT, ::INT, ::BYTEA[], ::TEXT, ::SMALLINT, ::TIMESTAMP)
- **lowercase:** table names, column names, aliases, parameter names
- **One column per line** in SELECT, INSERT column lists, and VALUES
- **Blank line** between queries
- **Semicolon** terminates each query
- **Inline comments** with `--` for non-obvious logic

## Common Patterns

### CTE multi-table insert

InsertGatewayEnvelope pattern:

```sql
WITH m AS (INSERT ... RETURNING 1),
b AS (INSERT ... RETURNING 1)
SELECT ...
```

### Find or create (upsert + fallback)

FindOrCreatePayer pattern:

```sql
WITH ins AS (
	INSERT ... ON CONFLICT DO NOTHING RETURNING id
), u AS (
	SELECT id FROM ins
	UNION ALL
	SELECT id FROM table WHERE ...
)
SELECT id FROM u LIMIT 1
```

### ON CONFLICT DO NOTHING

Idempotent inserts: InsertNodeInfo, InsertOrIgnorePayerReport.

### ON CONFLICT DO UPDATE

Increment counters (IncrementOriginatorCongestion, IncrementUnsettledUsage), conditional updates with WHERE clause (InsertAddressLog).

### RETURNING

Return inserted/affected data. `RETURNING *` (InsertStagedOriginatorEnvelope), `RETURNING id` (FindOrCreatePayer).

### Array parameters

`@param::TYPE[]` with ANY or unnest:

- `ANY(@topics::BYTEA[])` (SelectNewestFromTopics)
- `unnest(@addresses::TEXT[])` for batch operations

### Nullable optional filters

```sql
sqlc.narg('name')::TYPE IS NULL OR condition
```

Used in FetchPayerReports, FetchAttestations.

### Zero-means-unset optional filters

```sql
@param::BIGINT = 0 OR condition
```

Used in SumOriginatorCongestion, GetPayerUnsettledUsage.

### NULLIF for optional limits

```sql
LIMIT NULLIF(@row_limit::INT, 0)
```

0 = no limit.

### LATERAL joins

**Per (topic, originator):** `SelectGatewayEnvelopesByTopics` uses `CROSS JOIN LATERAL` for per-(topic, originator) index probes on `gem_topic_orig_seq_idx`, with a second per-originator LATERAL for blob join cache locality. Callers must include ALL originators in cursor arrays — use `FillMissingOriginators(vc, allOriginators)` on the `VectorClock` before `SetVectorClockByTopics`.

**Per originator:** `SelectGatewayEnvelopesByOriginators` uses `CROSS JOIN LATERAL` for per-originator cursor-based pagination.

### FOR UPDATE SKIP LOCKED

Concurrent-safe batch processing: GetNextAvailableNonce, DeleteExpiredEnvelopesBatch.

### CTE batch delete

```sql
WITH to_delete AS (
	SELECT ... LIMIT @batch_size FOR UPDATE SKIP LOCKED
)
DELETE ... USING to_delete ... RETURNING ...
```

### Window functions

`COUNT(...) OVER (PARTITION BY ...)` in FetchPayerReports attestations_count.

### DISTINCT ON

Newest per group: SelectNewestFromTopics.

### Advisory locks

`pg_advisory_xact_lock()`, `pg_try_advisory_xact_lock()`.

### Stored function calls

`SELECT * FROM function_name(params)` or `SELECT function_name(params)`.

### COALESCE for defaults

```sql
COALESCE(SUM(amount), 0)::BIGINT
```

Used in GetPayerBalance, SumOriginatorCongestion.

### Encode/decode hex

`encode(inbox_id, 'hex')`, `decode(@inbox_id, 'hex')` in identity_updates.sql.

### SAVEPOINT management

InsertSavePoint, InsertSavePointRelease, InsertSavePointRollback in partitions.sql.

## After Writing Queries

1. Run `dev/gen/all` to regenerate Go code
2. Verify the generated Go in `pkg/db/queries/` compiles
3. Run `dev/test` for integration tests

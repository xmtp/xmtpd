---
name: writing-migrations
description: >-
  Creates and modifies PostgreSQL database migrations for xmtpd using golang-migrate.
  Use when adding or altering tables, columns, indexes, functions, triggers,
  constraints, or partitions, or when the user mentions migrations or schema changes.
---

# Writing Database Migrations

## Creating a Migration

Run `dev/gen/migration {name}` to produce:

- `pkg/migrations/NNNNN_name.up.sql`
- `pkg/migrations/NNNNN_name.down.sql`

5-digit sequential numbering. Hyphen-separated lowercase names (e.g., `add-latest-block`, `payer-nonces`, `add-dead-letter-box`).

## How Migrations Run

- Embedded via `//go:embed *.sql` in `pkg/migrations/migrations.go`
- Uses golang-migrate with PostgreSQL driver
- Each migration runs in a transaction
- Tracked in `schema_migrations` table

## SQL Formatting

Same as queries: tab indentation, UPPERCASE keywords/functions/types, lowercase identifiers.

Additional conventions:

- Comment at top explains purpose when non-obvious
- `SET statement_timeout = 0` for long-running DDL (00006, 00018)
- `IF NOT EXISTS` / `IF EXISTS` for idempotency

## Up Migration Patterns

### CREATE TABLE

Project data type conventions:

| Type | Usage |
|------|-------|
| `BYTEA` | Binary data (keys, hashes, envelopes, topics) |
| `TEXT` | String identifiers (addresses) |
| `BIGINT` | Sequence IDs, unix timestamps, picodollar amounts |
| `INTEGER` | Node IDs, payer IDs |
| `SMALLINT` | Enum/status codes |
| `SERIAL` / `BIGSERIAL` | Auto-increment PKs |
| `TIMESTAMP` | With `DEFAULT NOW()` or `DEFAULT CURRENT_TIMESTAMP` |
| `BOOLEAN` | Flags |
| `INT[]` | Integer arrays |

### Primary keys

Composite PKs are common:

```sql
PRIMARY KEY (originator_node_id, originator_sequence_id)
```

### CHECK constraints

Singleton pattern:

```sql
CONSTRAINT is_singleton CHECK (singleton_id = 1)
```

### CREATE INDEX

Naming: `{table_prefix}_{columns}_idx` (e.g., `gem_time_node_seq_idx`, `gem_expiry_idx`).

Use `INCLUDE` for covering indexes, filtered indexes with `WHERE`.

### PL/pgSQL functions

```sql
CREATE OR REPLACE FUNCTION name(params)
RETURNS type AS $$
BEGIN
	...
END;
$$ LANGUAGE plpgsql;
```

### SQL functions

```sql
CREATE OR REPLACE FUNCTION name(params)
RETURNS TABLE(...)
LANGUAGE SQL AS $$
	...
$$;
```

### Triggers

Row-level:

```sql
FOR EACH ROW EXECUTE FUNCTION func()
```

Statement-level (preferred for bulk):

```sql
REFERENCING NEW TABLE AS new FOR EACH STATEMENT EXECUTE FUNCTION func()
```

### Views

```sql
CREATE OR REPLACE VIEW name AS SELECT ...
```

### Partitioning

Two-level: LIST by `originator_node_id`, then RANGE by `originator_sequence_id` bands.

### ALTER TABLE

For adding constraints, FKs:

```sql
ADD CONSTRAINT fk_name FOREIGN KEY (col) REFERENCES table(col)
```

### Data seeding

```sql
INSERT INTO table VALUES (...), (...);
```

### Versioning

Functions versioned with `_v2` suffix rather than dropped/recreated (e.g., `ensure_gateway_parts_v2`, `update_latest_envelope_v2`).

## Down Migration Patterns

Complete reversal in reverse dependency order:

1. `DROP TRIGGER IF EXISTS name ON table` (triggers first)
2. `DROP FUNCTION IF EXISTS name` (then functions)
3. `DROP VIEW IF EXISTS name`
4. `DROP INDEX IF EXISTS name`
5. `ALTER TABLE ... DROP CONSTRAINT IF EXISTS name`
6. `DROP TABLE IF EXISTS name CASCADE` (parent tables last)

Always use `IF EXISTS` for safety.

## Updating Migration Tests

In `pkg/migrations/migrations_test.go`:

1. Increment `const currentMigration` to match new migration number
2. Add a `checkXxx` function asserting new schema objects exist using helpers: `tableExists()`, `indexExists()`, `functionExists()`, `triggerExists()`, `viewExists()`, `constraintExists()`
3. Add corresponding `t.Run("NNNNN_name", ...)` call in `TestMigrations`
4. Run `dev/test ./pkg/migrations/...`

## After Writing Migrations

1. Run `dev/gen/all` to regenerate sqlc (migrations are the schema source)
2. Run `dev/test` to verify

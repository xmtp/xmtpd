# Per-Topic Cursor Query Strategy: Summary of Findings

## Problem Statement

The current production query (`SelectGatewayEnvelopesByTopics`, called "V0") uses a **single vector clock** (one cursor per originator) shared across all subscribed topics. This means a client subscribed to 500 topics must re-fetch data it has already seen on fast-moving topics whenever a slow-moving topic needs to catch up. Per-topic cursors — where each topic maintains its own per-originator cursor position — eliminate this redundant re-fetching.

The question: can per-topic cursors be queried efficiently at production scale?

## Production Data Profile (Feb 2026)

| Originator | Max Seq | Est. Rows | RANGE Partitions |
| ---------- | ------- | --------- | ---------------- |
| o0         | 5.6M    | ~5.6M     | ~6               |
| o1         | 3.5M    | ~3.5M     | ~4               |
| o10        | 22.7M   | ~22.7M    | ~23              |
| o11        | 111.5M  | ~111.5M   | ~112             |
| o13        | 3.1M    | ~3.1M     | ~4               |
| o100       | 40K     | ~40K      | 1                |
| o200       | 31K     | ~31K      | 1                |

- **Total**: ~146M rows, ~53K distinct topics, 7 originators, ~151 leaf partitions
- **Skew**: o11 has 76% of all data
- **Typical subscription**: 50–500 topics

## Is it performant to introduce per-topic cursors?

**Yes**, with the right query and index. At 10M rows with production-realistic skew (76%/15%/7.6% across 3 originators, 50K topics, 500 subscribed):

| Cursor Position                  | V0 (status quo) | V3b + Index (per-topic) |
| -------------------------------- | --------------- | ----------------------- |
| 80% (nearly caught up)           | 1,476ms         | **23ms**                |
| 20% (far behind)                 | 1,749ms         | **23ms**                |
| Mixed (half at 80%, half at 20%) | 1,880ms         | **23ms**                |

V3b with `gem_topic_orig_seq_idx` is **57-82x faster** than V0 across all cursor distributions, while also providing per-topic cursor precision that eliminates redundant data transfer.

## Recommended Query and Index Strategy

### Index

```sql
CREATE INDEX gem_topic_orig_seq_idx ON gateway_envelopes_meta
  (topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time);
```

The `INCLUDE (gateway_time)` makes it a **covering index**, enabling index-only scans with zero heap fetches. This index also replaces three existing indexes that can be dropped (`gem_topic_time_idx`, `gem_time_node_seq_idx`, `gem_originator_node_id`), resulting in a net reduction from 5 indexes to 3.

### Query (V3b)

The query uses three phases:

1. **LATERAL per (topic, originator) on meta only** — For each (topic, originator) cursor entry, do a targeted index probe. Most probes return 0 rows instantly (topic doesn't exist for that originator past the cursor). Only matching rows are collected.
2. **CTE with sort + LIMIT** — Materializes at most `row_limit` rows (e.g., 500).
3. **Per-originator blob join** — Groups blob lookups by originator for cache locality, reducing per-lookup time from ~1ms to ~0.2ms.

The caller pre-populates all known originators per topic with `seq_id=0` for unseen originators, eliminating the need for a UNION ALL fallback branch.

### Caller Requirements

V3b requires that **all known originators** appear in the cursor arrays for every query. Without this, the LATERAL subquery skips originators not in the cursor, silently missing data.

The API layer handles this via two components:

1. **`CachedOriginatorList`** (`pkg/db/originator_list.go`) — Queries `gateway_envelopes_latest` for all known originator node IDs, cached with a configurable TTL (5 minutes in production). Thread-safe with double-checked locking.

2. **`FillMissingOriginators`** (`pkg/db/types.go`) — Takes the caller's `SelectGatewayEnvelopesByTopicsParams` and the full originator list, appending any missing originators with `seq_id=0`.

The call chain in `pkg/api/message/service.go` is:

```
allOriginators := s.originatorList.GetOriginatorNodeIDs(ctx)
db.FillMissingOriginators(&params, allOriginators)
```

For specialized callers (e.g., `IdentityUpdateStorer`) that query a single known originator, the originator can be hardcoded directly in the params without using `FillMissingOriginators`.

## Comparison to V0 (Status Quo)

| Dimension                          | V0                                                        | V3b + Index                                                           |
| ---------------------------------- | --------------------------------------------------------- | --------------------------------------------------------------------- |
| **Cursor model**                   | Shared (one per originator)                               | Per-topic (one per topic x originator)                                |
| **Meta access**                    | O(partitions) — one scan per partition with `ANY(topics)` | O(topics x originators) — one index probe per pair                    |
| **Planning time**                  | 9-35ms (grows with topic count due to `ANY(large_array)`) | 0.4-0.5ms (stable)                                                    |
| **Execution at 10M rows**          | 1,024-1,845ms                                             | 22-23ms                                                               |
| **Sensitivity to cursor position** | Moderate (44ms-416ms meta phase depending on cursor)      | Minimal (22-23ms regardless)                                          |
| **Redundant re-fetching**          | Yes — shared cursor forces re-scanning already-seen data  | No — each topic tracks its own position                               |
| **Index requirement**              | Uses existing `gem_topic_time_idx`                        | Requires `gem_topic_orig_seq_idx` (but allows dropping 3 old indexes) |
| **Write overhead**                 | N/A                                                       | Marginal — net fewer indexes after cleanup                            |

**Key trade-off**: V3b makes O(topics x originators) individual index probes instead of V0's O(partitions) batch scans. Without the covering index, V3b's per-probe cost is high enough (heap fetches, missing indexes on some partitions) that V0's batch approach wins. With the index, each probe is a sub-microsecond index-only seek, and the 1500 probes (500 topics x 3 originators) complete in ~22ms total.

## Why V3b Without the Index Is Not Suitable

V3b's LATERAL-per-(topic, originator) pattern makes O(topics x originators) individual probes against the meta table. Each probe needs to seek by `(topic, originator_node_id, originator_sequence_id)`. Without an index on exactly those columns, three compounding problems emerge:

### 1. Hot-spot amplification on the dominant originator

The table is LIST-partitioned by `originator_node_id`, then RANGE-partitioned by `originator_sequence_id` in 1M bands. The dominant originator (76% of data) has ~8 RANGE subpartitions. Without the covering index, V3b must scan the dominant originator's partitions **once per subscribed topic** — 500 topics means 500 scans of the same partitions.

This was observed directly at ~9M rows in production-scale testing. The `o11_s109M` partition (a single RANGE leaf of the dominant originator) lacked the composite index and fell back to a bitmap heap scan on the single-column `originator_node_id` index:

- **V0** scanned this partition **once** (loops=1): 70 heap blocks, 5,227 rows, 0.47ms.
- **V3b** scanned it **1,000 times** (loops=1000): 70,000 heap blocks, 5.2M cumulative row reads, 224ms total.

This single partition accounted for **~70% of V3b's entire meta phase** at low cursor positions. V0 handles the same partition gracefully because its `ANY(topics)` predicate batches all topic filtering into one pass.

### 2. Heap fetches instead of index-only scans

V3b's LATERAL subquery selects `gateway_time` alongside the key columns. Without INCLUDE (gateway_time) in the index, every probe that finds a matching row must do a heap fetch to retrieve `gateway_time`. At 10M rows with the covering index, V3b achieves **zero heap fetches** — the entire query is satisfied from the index alone. Without it, each of the ~500 result rows (plus any intermediate rows examined and discarded) requires a random I/O to the heap, adding milliseconds per probe across thousands of probes.

### 3. Cursor-position sensitivity

Without the index, V3b's performance degrades sharply as cursors move backward:

| Cursor | V3b (no index) | V3b (with index) |
| ------ | -------------- | ---------------- |
| 80%    | 321ms          | 23ms             |
| 20%    | 603ms          | 23ms             |
| Mixed  | 412ms          | 22ms             |

At low cursor positions, more rows exist past each cursor, so each LATERAL probe scans more data before hitting the LIMIT. With the covering index, the probe is an index-only range scan that starts exactly at the cursor position and reads forward — the cost is proportional to the number of rows _returned_, not the number of rows _after the cursor_. Without the index, the probe uses bitmap heap scans that degrade with the volume of post-cursor data.

### 4. The numbers in context

At ~9M rows with 1000 topics and 7 originators (no covering index, warm cache), V3b's meta phase was **6-39x slower** than V0's depending on cursor position. V0's meta phase was stable at ~7ms regardless of cursor position, while V3b ranged from 46ms to 322ms. Without the index, V3b's lower per-probe cost (targeted seek vs batch scan) is overwhelmed by the 300x more operations it must perform.

The covering index flips this equation: each probe becomes a sub-microsecond index-only seek, making the aggregate cost of thousands of probes (~22ms) dramatically less than V0's fewer but heavier partition scans (~1,000ms+).

## Methodology

### Test Infrastructure

- **10M rows** seeded in a Dockerized PostgreSQL instance with production-realistic originator skew
- **3 originators** (7.6%, 15%, 76% of data) producing 14 RANGE partitions per table
- **50K topics**, 500 subscribed, with power-law frequency distribution
- **JIT disabled** for consistent timing

### Warm-Cache Protocol

Early testing produced a false result: V3b appeared 5.3x faster than V0 at ~9M rows. Multi-distribution testing revealed this was a **cache warming artifact** — V0 ran first on cold cache (372ms blob phase), V3b ran second on warm cache (2.5ms blob phase). The 200x blob timing difference was entirely cache state, not query efficiency.

To eliminate this bias, the 10M-row test runs each query **3 times** for warm-up before the measured `EXPLAIN (ANALYZE, BUFFERS)` pass. All variants see equally warm caches.

### Cursor Distributions

Three distributions tested per variant: 80% (nearly caught up), 20% (far behind), and mixed (per-topic variation). The mixed distribution is the most realistic and the scenario where per-topic cursors provide the most value over shared cursors.

### Correctness Verification

Before performance testing: V0c confirmed identical to V0, V3b confirmed as superset of V0, V4 capture rate verified at 100% for uniform cursors, V5 confirmed identical to V3b.

## Alternatives Tested and Why They Failed

### V0 variants (shared cursor optimizations)

**V0c** (CTE topic list, no UNION ALL): Fixed V0's O(N) planning time (1ms vs 35ms) but execution remained ~1,000ms. The UNION ALL removal helps but doesn't change the fundamental scan pattern.

**V0d / V6** (LATERAL per originator): Only 3-7 iterations instead of thousands, but each iteration scans all post-cursor rows for one originator in sequence order, filtering by topic via join. With sparse topic selectivity (~500 topics out of 50K distinct), the query scans hundreds of thousands of rows per originator before finding enough matches. **Timed out at >120 seconds** on ~9M rows. Works at 10K rows but catastrophically fails at scale.

### Per-topic cursor alternatives

**V4** (hybrid V0 scan + per-topic post-filter): Uses floor cursors (min per originator across topics) for a V0-style coarse scan, then post-filters with per-topic cursors. The floor cursor is too low when topics diverge, causing the 3x over-fetch limit to scan excessive data. Execution: 970-1,335ms — no better than V0c.

**V5** (temp table cursors): Inserts per-topic cursors into a temp table, lets the planner choose join strategy. Viable at 80% cursors (352ms) but degrades at 20% (1,104ms). The temp table overhead (CREATE + INSERT + ANALYZE) doesn't justify itself when V3b+index does 22ms.

### LATERAL variants on blobs

**V2b** (blob join inside LATERAL): Excellent at 1 partition per originator but catastrophic with multiple RANGE subpartitions — the planner cannot prune by `originator_sequence_id` from a join, scanning all subpartitions per iteration. 499ms at 10 partitions, 1009ms at 50.

**V1/V1b** (LATERAL per topic with scalar subquery cursor): Correlated subquery scans full cursor CTE for each meta row. O(N x M x cursor_size) complexity. Disqualified at >400ms on 10K rows.

## Source Documents

- `doc/plans/2026-02-13-per-topic-cursor-query-design.md` — Initial investigation at 10K and ~9M rows
- `doc/plans/2026-02-16-scaled-query-investigation-design.md` — 10M-row test design and results
- `doc/plans/2026-02-16-scaled-query-investigation-plan.md` — Implementation plan for the test harness

---

## Appendix A: Full Query Definitions

All query variants tested during this investigation, organized by cursor model.

### Shared-Cursor Variants

These variants use a single vector clock (one cursor per originator) shared across all subscribed topics.

#### V0: Original Baseline

Production query at the start of the investigation. Uses `ANY(topics)` with UNION ALL to handle known and unknown originators separately.

**Parameters**: `$1` = cursor_node_ids `INT[]`, `$2` = cursor_seq_ids `BIGINT[]`, `$3` = row_limit `INT`, `$4` = topics `BYTEA[]`

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest($1::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest($2::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq
    FROM unnest($2::BIGINT[]) AS t(seq_id)
),
filtered AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN cursors AS c
         ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.topic = ANY ($4::BYTEA[])
      AND m.originator_node_id = ANY($1::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)

    UNION ALL

    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    WHERE m.topic = ANY ($4::BYTEA[])
      AND m.originator_sequence_id > 0
      AND NOT EXISTS (
          SELECT 1 FROM cursors AS c
          WHERE c.cursor_node_id = m.originator_node_id
      )

    ORDER BY originator_node_id, originator_sequence_id
    LIMIT NULLIF($3::INT, 0)
)
SELECT f.originator_node_id,
       f.originator_sequence_id,
       f.gateway_time,
       f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b
     ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

#### V0b: CTE Topic List with UNION ALL

Replaces `ANY($topics)` with a CTE `topic_list` joined via `unnest()`. Keeps UNION ALL. Eliminates O(N) planning regression at 1000 topics.

**Parameters**: Same as V0.

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord) USING (ord)
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq
    FROM unnest(@cursor_seq_ids::BIGINT[]) AS t(seq_id)
),
topic_list AS (
    SELECT t.topic FROM unnest(@topics::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT m.originator_node_id, m.originator_sequence_id, m.gateway_time, m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    JOIN cursors AS c ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.originator_node_id = ANY(@cursor_node_ids::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    UNION ALL
    SELECT m.originator_node_id, m.originator_sequence_id, m.gateway_time, m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    WHERE m.originator_sequence_id > 0
      AND NOT EXISTS (SELECT 1 FROM cursors AS c WHERE c.cursor_node_id = m.originator_node_id)
    ORDER BY originator_node_id, originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
)
SELECT f.originator_node_id, f.originator_sequence_id, f.gateway_time, f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

#### V0c: CTE Topic List, No UNION ALL

Same as V0b but removes the UNION ALL branch. Requires the caller to pre-populate all originators in the cursor with `seq_id=0` for unknown originators.

**Parameters**: Same as V0.

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord) USING (ord)
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq
    FROM unnest(@cursor_seq_ids::BIGINT[]) AS t(seq_id)
),
topic_list AS (
    SELECT t.topic FROM unnest(@topics::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT m.originator_node_id, m.originator_sequence_id, m.gateway_time, m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    JOIN cursors AS c ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.originator_node_id = ANY(@cursor_node_ids::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
)
SELECT f.originator_node_id, f.originator_sequence_id, f.gateway_time, f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

#### V0d / V6: LATERAL Per Originator with Topic Join

Only 3–7 LATERAL iterations (one per originator) instead of O(N) topic expansion. Each iteration scans one originator's partition, filtering by topics via join. Works at 10K rows but **catastrophically fails at production scale** with sparse topic selectivity — scans millions of rows per originator before finding enough matches.

**Parameters**: Same as V0.

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord) USING (ord)
),
topic_list AS (
    SELECT t.topic FROM unnest(@topics::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT sub.originator_node_id, sub.originator_sequence_id, sub.gateway_time, sub.topic
    FROM cursors AS c
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id, m.originator_sequence_id, m.gateway_time, m.topic
        FROM gateway_envelopes_meta AS m
        JOIN topic_list AS tl ON m.topic = tl.topic
        WHERE m.originator_node_id = c.cursor_node_id
          AND m.originator_sequence_id > c.cursor_sequence_id
        ORDER BY m.originator_sequence_id
        LIMIT @row_limit
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
)
SELECT f.originator_node_id, f.originator_sequence_id, f.gateway_time, f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

### Per-Topic Cursor Variants

These variants use per-topic cursors — each topic maintains its own per-originator cursor position.

#### V3b (Investigation Version): LATERAL Per (Topic, Originator) with Per-Originator Blob Join

The recommended query. Uses three phases: LATERAL per (topic, originator) on meta, CTE with sort + LIMIT, per-originator blob join for cache locality. The caller pre-populates all known originators per topic with `seq_id=0` for unseen originators.

**Parameters**: `$1` = cursor_topics `BYTEA[]`, `$2` = cursor_node_ids `INT[]`, `$3` = cursor_seq_ids `BIGINT[]`, `$4` = rows_per_entry `INT`, `$5` = row_limit `INT`

```sql
WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
filtered AS (
    SELECT sub.originator_node_id,
           sub.originator_sequence_id,
           sub.gateway_time,
           sub.topic
    FROM cursor_entries AS ce
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id,
               m.originator_sequence_id,
               m.gateway_time,
               m.topic
        FROM gateway_envelopes_meta AS m
        WHERE m.topic = ce.topic
          AND m.originator_node_id = ce.node_id
          AND m.originator_sequence_id > ce.seq_id
        ORDER BY m.originator_sequence_id
        LIMIT $4
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT $5
),
originator_ids AS (
    SELECT DISTINCT originator_node_id FROM filtered
)
SELECT bl.originator_node_id,
       bl.originator_sequence_id,
       bl.gateway_time,
       bl.topic,
       bl.originator_envelope
FROM originator_ids AS oi
CROSS JOIN LATERAL (
    SELECT f.originator_node_id,
           f.originator_sequence_id,
           f.gateway_time,
           f.topic,
           b.originator_envelope
    FROM filtered AS f
    JOIN gateway_envelope_blobs AS b
        ON b.originator_node_id = oi.originator_node_id
       AND b.originator_sequence_id = f.originator_sequence_id
    WHERE f.originator_node_id = oi.originator_node_id
) AS bl
ORDER BY bl.originator_node_id, bl.originator_sequence_id
```

#### V3b (Production Version): sqlc-Deployed Query

The deployed version in `pkg/db/sqlc/envelopes_v2.sql`. Differs from the investigation version: uses a CROSS JOIN to expand topics × cursors (instead of pre-flattened arrays), and dynamically calculates `rows_per_entry` from `row_limit / num_originators` with a floor of 10.

**Parameters**: `@cursor_node_ids` = `INT[]`, `@cursor_sequence_ids` = `BIGINT[]`, `@topics` = `BYTEA[]`, `@row_limit` = `INT`

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
cursor_entries AS (
    SELECT t.topic, c.cursor_node_id AS node_id, c.cursor_sequence_id AS seq_id
    FROM unnest(@topics::BYTEA[]) AS t(topic)
    CROSS JOIN cursors AS c
),
rows_per_entry AS (
    SELECT GREATEST(
        NULLIF(@row_limit::INT, 0) / GREATEST(array_length(@cursor_node_ids::INT[], 1), 1),
        10
    ) AS val
),
filtered AS (
    SELECT sub.originator_node_id,
           sub.originator_sequence_id,
           sub.gateway_time,
           sub.topic
    FROM cursor_entries AS ce
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id,
               m.originator_sequence_id,
               m.gateway_time,
               m.topic
        FROM gateway_envelopes_meta AS m
        WHERE m.topic = ce.topic
          AND m.originator_node_id = ce.node_id
          AND m.originator_sequence_id > ce.seq_id
        ORDER BY m.originator_sequence_id
        LIMIT (SELECT val FROM rows_per_entry)
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
),
originator_ids AS (
    SELECT DISTINCT originator_node_id FROM filtered
)
SELECT bl.originator_node_id,
       bl.originator_sequence_id,
       bl.gateway_time,
       bl.topic,
       bl.originator_envelope
FROM originator_ids AS oi
CROSS JOIN LATERAL (
    SELECT f.originator_node_id,
           f.originator_sequence_id,
           f.gateway_time,
           f.topic,
           b.originator_envelope
    FROM filtered AS f
    JOIN gateway_envelope_blobs AS b
        ON b.originator_node_id = oi.originator_node_id
       AND b.originator_sequence_id = f.originator_sequence_id
    WHERE f.originator_node_id = oi.originator_node_id
) AS bl
ORDER BY bl.originator_node_id, bl.originator_sequence_id
```

#### V4: Hybrid V0 + Per-Topic Post-Filter

Combines V0's efficient single-pass meta scan with per-topic cursor precision via post-filtering. Uses floor cursors (min per originator across all topics) for a coarse scan, then joins with per-topic cursor entries to filter precisely. The 3× over-fetch limit handles moderate cursor divergence but fails when topics have widely divergent cursor positions.

**Parameters**: `$1` = cursor_topics `BYTEA[]`, `$2` = cursor_node_ids `INT[]`, `$3` = cursor_seq_ids `BIGINT[]`, `$4` = rows_per_entry `INT` (unused), `$5` = row_limit `INT`

```sql
WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
distinct_topics AS (
    SELECT DISTINCT topic FROM cursor_entries
),
floor_cursors AS (
    SELECT node_id AS cursor_node_id,
           MIN(seq_id) AS cursor_sequence_id
    FROM cursor_entries
    GROUP BY node_id
),
distinct_node_ids AS (
    SELECT DISTINCT node_id FROM cursor_entries
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(cursor_sequence_id), 0) AS min_seq
    FROM floor_cursors
),
coarse AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN floor_cursors AS c
         ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    JOIN distinct_topics AS dt ON m.topic = dt.topic
    WHERE m.originator_node_id = ANY(
        (SELECT array_agg(node_id) FROM distinct_node_ids)
    )
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF($5::INT, 0) * 3
),
filtered AS (
    SELECT co.originator_node_id,
           co.originator_sequence_id,
           co.gateway_time,
           co.topic
    FROM coarse AS co
    JOIN cursor_entries AS ce
         ON ce.topic = co.topic
         AND ce.node_id = co.originator_node_id
    WHERE co.originator_sequence_id > ce.seq_id
    ORDER BY co.originator_node_id, co.originator_sequence_id
    LIMIT NULLIF($5::INT, 0)
)
SELECT f.originator_node_id,
       f.originator_sequence_id,
       f.gateway_time,
       f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b
     ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

#### V5: Temp Table Cursors

Inserts per-topic cursors into an unlogged temp table, then hash-joins with meta. Planner sees real statistics on the temp table. Two-step execution is more complex to implement.

**Step 1 — Setup** (params: `$1` = cursor_topics `BYTEA[]`, `$2` = cursor_node_ids `INT[]`, `$3` = cursor_seq_ids `BIGINT[]`):

```sql
CREATE TEMP TABLE IF NOT EXISTS _cursor_entries (
    topic BYTEA NOT NULL,
    node_id INT NOT NULL,
    seq_id BIGINT NOT NULL
);
TRUNCATE _cursor_entries;
INSERT INTO _cursor_entries (topic, node_id, seq_id)
SELECT t.topic, n.node_id, s.seq_id
FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord);
ANALYZE _cursor_entries
```

**Step 2 — Main query** (params: `$1` = distinct_node_ids `INT[]`, `$2` = row_limit `INT`):

```sql
WITH min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq FROM _cursor_entries
),
filtered AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN _cursor_entries AS ce
         ON m.topic = ce.topic
         AND m.originator_node_id = ce.node_id
         AND m.originator_sequence_id > ce.seq_id
    WHERE m.originator_node_id = ANY($1::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF($2::INT, 0)
)
SELECT f.originator_node_id,
       f.originator_sequence_id,
       f.gateway_time,
       f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b
     ON b.originator_node_id = f.originator_node_id
     AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id
```

## Appendix B: Index Definitions

### Recommended New Index

```sql
CREATE INDEX gem_topic_orig_seq_idx ON gateway_envelopes_meta
  (topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time);
```

### Indexes to Drop (Superseded)

| Index                    | Definition                                                          | Why Redundant                                                                         |
| ------------------------ | ------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| `gem_topic_time_idx`     | `(topic, gateway_time, originator_node_id, originator_sequence_id)` | New index is strictly better for all consumers; `gateway_time` never appears in WHERE |
| `gem_time_node_seq_idx`  | `(gateway_time, originator_node_id, originator_sequence_id)`        | No query uses `gateway_time` as leading filter; legacy index                          |
| `gem_originator_node_id` | `(originator_node_id)`                                              | Redundant with LIST partitioning + PRIMARY KEY                                        |

### Indexes to Keep

| Index                     | Definition                                                                        | Reason                               |
| ------------------------- | --------------------------------------------------------------------------------- | ------------------------------------ |
| `gem_topic_time_desc_idx` | `(topic, gateway_time DESC) INCLUDE (originator_node_id, originator_sequence_id)` | Required by `SelectNewestFromTopics` |
| `gem_expiry_idx`          | `(expiry) INCLUDE (...) WHERE expiry IS NOT NULL`                                 | Required by all pruning queries      |

**Net result**: Adding 1 index and dropping 3 reduces total count from 5 to 3, lowering write amplification while dramatically improving read performance.

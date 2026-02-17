# Scaled Query Investigation Design

## Goal

Expand the per-topic cursor query investigation from 10K rows to 10M rows with production-realistic data distribution, eliminate cache bias with warm-up passes, and explore new query shapes and indexing patterns that may outperform both V0 and V3b.

## Background

The [initial investigation](2026-02-13-per-topic-cursor-query-design.md) found:

- **V3b** (LATERAL per topic×originator, per-originator blob join) appeared 5.3× faster than V0, but multi-distribution testing revealed this was a **cache warming artifact**.
- **V0's meta phase is 6–39× faster** than V3b because it scans O(partitions) instead of O(topics × originators).
- **Per-topic cursors remain valuable** as a data model (eliminates redundant re-fetching) but need a query strategy that avoids V3b's O(N×M) LATERAL overhead.

## Production Data Profile

From production queries (Feb 2026):

| Originator | Max Seq | Est. Rows | RANGE Partitions |
|---|---|---|---|
| o0 | 5.6M | ~5.6M | ~6 |
| o1 | 3.5M | ~3.5M | ~4 |
| o10 | 22.7M | ~22.7M | ~23 |
| o11 | 111.5M | ~111.5M | ~112 |
| o13 | 3.1M | ~3.1M | ~4 |
| o100 | 40K | ~40K | 1 |
| o200 | 31K | ~31K | 1 |

- **Total**: ~146M rows, ~53K distinct topics, 7 originators, ~151 leaf partitions
- **Skew**: o11 has 76% of all data
- **Typical subscription**: 50–500 topics

## Test Infrastructure

### File

`pkg/db/query_scalability_test.go` — new file, separate from existing explain/correctness tests.

### Shared Database via TestMain

TestMain seeds a single 10M row database shared across all subtests:

- **3 originators** modeling production skew:
  - o1: ~760K rows (7.6%), max_seq ~760K, 1 RANGE partition
  - o2: ~1.5M rows (15%), max_seq ~1.5M, 2 RANGE partitions
  - o3: ~7.6M rows (76%), max_seq ~7.6M, 8 RANGE partitions
- **~50K distinct topics** with power-law frequency distribution (most topics have 1–10 rows, few topics have 1000+)
- **500 "subscribed" topics** selected from the pool for query parameters
- `ANALYZE` both tables after seeding
- `SET jit = off` for consistent timing

### Warm-Up Protocol

Before each variant's measured EXPLAIN ANALYZE:

1. Run the query **3 times** discarding results (warms buffer cache + plan cache)
2. Then run `EXPLAIN (ANALYZE, BUFFERS)` for the measured execution

This eliminates the cold-cache bias that invalidated the original V3b comparison.

### Cursor Distributions

Each variant runs with 3 cursor positions:

| Distribution | Description |
|---|---|
| **80%** | Cursors at 80% of each originator's max_seq (nearly caught up) |
| **20%** | Cursors at 20% of each originator's max_seq (far behind) |
| **Mixed** | Per-topic variation: half of subscribed topics at 80%, half at 20%. For V0 (shared cursor), uses min() across topic cursors per originator. |

## Query Variants

### Baseline (Carried Forward)

**V0** — Original `ANY(topics)` + UNION ALL + shared cursor. Production query.

**V0c** — CTE topic list, no UNION ALL, pre-populated originators. Best V0 optimization from prior work.

**V3b** — LATERAL per (topic, originator) on meta, per-originator blob join. Best per-topic cursor variant from prior work.

### New Variants

#### V4: Hybrid V0 + Per-Topic Post-Filter

Combines V0's efficient single-pass meta scan with per-topic cursor precision via post-filtering.

```sql
WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest(@cursor_topics::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
-- Compute floor cursor per originator (minimum across all topics)
floor_cursors AS (
    SELECT node_id AS cursor_node_id,
           MIN(seq_id) AS cursor_sequence_id
    FROM cursor_entries
    GROUP BY node_id
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(cursor_sequence_id), 0) AS min_seq
    FROM floor_cursors
),
-- V0-style single pass using floor cursors
coarse AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN floor_cursors AS c
         ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.topic = ANY(@topics::BYTEA[])
      AND m.originator_node_id = ANY(@cursor_node_ids::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    -- Over-fetch to account for post-filter reduction
    LIMIT NULLIF(@row_limit::INT, 0) * 3
),
-- Post-filter: apply per-topic cursors
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
    LIMIT NULLIF(@row_limit::INT, 0)
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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

**Trade-off**: Gets V0's O(partitions) meta speed. Over-fetches rows between floor cursor and per-topic cursor — the 3× over-fetch limit handles moderate cursor divergence but may miss rows if divergence is extreme. The `ANY(@topics)` still has O(N) planning, but V0c's CTE topic approach could be used instead.

**V4b variant**: Same approach but using CTE `topic_list` join instead of `ANY(@topics)` to avoid O(N) planning overhead.

#### V5: Temp Table Cursors

Insert per-topic cursors into an unlogged temp table, then hash-join with meta.

```sql
-- Step 1 (executed separately before the main query):
CREATE TEMP TABLE IF NOT EXISTS _cursor_entries (
    topic BYTEA NOT NULL,
    node_id INT NOT NULL,
    seq_id BIGINT NOT NULL
) ON COMMIT DROP;

INSERT INTO _cursor_entries (topic, node_id, seq_id)
SELECT t.topic, n.node_id, s.seq_id
FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord);

-- Optional: ANALYZE _cursor_entries;

-- Step 2: Main query
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
    WHERE m.originator_node_id = ANY($2::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF($4::INT, 0)
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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

**Trade-off**: Planner sees real statistics on the temp table and can choose hash/merge join strategies. But temp table creation + insert + analyze has per-query overhead. Two-step execution is more complex to implement.

#### V6: LATERAL Per Originator + Topic Join (with Index)

Only N_originator iterations (3–7) instead of N_topics × N_originators. Each iteration scans one originator's partition, filtering by topics via join. Requires `gem_topic_orig_seq_idx` to avoid the sequential scan problem that killed V0d.

```sql
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
topic_list AS (
    SELECT t.topic FROM unnest(@topics::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT sub.originator_node_id,
           sub.originator_sequence_id,
           sub.gateway_time,
           sub.topic
    FROM cursors AS c
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id,
               m.originator_sequence_id,
               m.gateway_time,
               m.topic
        FROM topic_list AS tl
        JOIN gateway_envelopes_meta AS m
            ON m.topic = tl.topic
           AND m.originator_node_id = c.cursor_node_id
           AND m.originator_sequence_id > c.cursor_sequence_id
        ORDER BY m.originator_sequence_id
        LIMIT @row_limit
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

**Trade-off**: Only 3–7 LATERAL iterations (vs V3b's thousands). With the right index, each iteration does index probes per topic. But this uses a shared cursor (not per-topic), so it doesn't solve the per-topic cursor problem on its own. Could be combined with V4's post-filter approach.

**Note**: This is a shared-cursor query (same params as V0). It's included to test whether the LATERAL-per-originator approach works at scale with the right index, since V0d failed without it.

## Index Experiments

| Index | Definition | Tested With |
|---|---|---|
| **gem_topic_orig_seq_idx** | `(topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time)` | V3b, V6 |
| **BRIN on dominant partition** | `USING BRIN (originator_sequence_id)` on o3's partition only | V0, V3b |

Each index experiment runs the full variant suite with and without the index, measuring the delta.

## Output Format

Each variant prints EXPLAIN (ANALYZE, BUFFERS) output. A summary table is printed at the end:

```
Cursor   | Variant          | Plan(ms) | Exec(ms) | Total(ms) | Rows | Buffers
80%      | V0_baseline      |      ... |      ... |       ... |  ... | ...
80%      | V0c_no_union     |      ... |      ... |       ... |  ... | ...
80%      | V3b_lateral      |      ... |      ... |       ... |  ... | ...
80%      | V4_hybrid        |      ... |      ... |       ... |  ... | ...
80%      | V5_temp_table    |      ... |      ... |       ... |  ... | ...
80%      | V6_lateral_orig  |      ... |      ... |       ... |  ... | ...
20%      | ...              |          |          |           |      |
mixed    | ...              |          |          |           |      |
```

## Correctness Validation

Before performance testing, verify result equivalence:

1. V0c returns same rows as V0 (when all originators pre-populated)
2. V4 returns same rows as V3b (both use per-topic cursors)
3. V5 returns same rows as V3b
4. V3b returns a superset of V0's results (per-topic cursors are more precise, never less)

## Success Criteria

1. Identify the fastest query variant at 10M rows with production-like skew
2. Determine whether any new variant can match V0's ~7ms meta speed while providing per-topic cursor precision
3. Quantify the over-fetch cost of the hybrid approach (V4)
4. Determine whether `gem_topic_orig_seq_idx` eliminates V3b's hot-spot problem on skewed partitions
5. Produce updated performance data for the design doc

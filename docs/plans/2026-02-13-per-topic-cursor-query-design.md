# Per-Topic Cursor Query Design

## Problem

`SelectGatewayEnvelopesByTopics` uses a single vector clock (per-originator cursor) shared across all topics. This means:
- All topics share the same cursor position per originator
- Planning time grows linearly with topic count (6.8ms at 1000 topics)
- `ANY(large_array)` degrades to sequential scans

## Solution: V3b — LATERAL per (topic, originator) on meta, per-originator blob join

Each topic gets its own vector clock (`originator_id => sequence_id`). The query uses CROSS JOIN LATERAL over flattened cursor entries on the meta table only, collects and limits results in a CTE, then joins with blobs via a per-originator LATERAL for optimal cache locality. Production validated at 5.3× faster than V0 on 9M rows.

### SQL

```sql
WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest(@cursor_topics::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
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
        LIMIT @rows_per_entry
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT @row_limit
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

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `cursor_topics` | `BYTEA[]` | One entry per (topic, originator) pair |
| `cursor_node_ids` | `INT[]` | Originator for each entry |
| `cursor_seq_ids` | `BIGINT[]` | Last seen sequence for each entry |
| `rows_per_entry` | `INT` | Max rows per (topic, originator) pair |
| `row_limit` | `INT` | Total max rows |

For N topics and M originators, the arrays have N*M entries. The node pre-populates all known originators per topic, using `seq_id=0` for originators not yet seen by the client.

## Design Decisions

### LATERAL on meta only, blob join deferred (V3 over V2b)

Initial testing with 1 partition per originator showed V2b (blob join inside LATERAL) as the winner. However, multi-partition testing revealed V2b degrades catastrophically with many RANGE subpartitions — each LATERAL iteration scans ALL blob subpartitions via sequential scan because `originator_sequence_id` comes from a join and can't be used for RANGE partition pruning.

V3 avoids this by:
1. Running the LATERAL on meta only (fast index seeks, no blob scanning)
2. Collecting results in a CTE with sort + LIMIT (materializes at most `row_limit` rows)
3. Joining with blobs once for just the limited result set (500 rows, not 3000+)

At 10 partitions per originator, V3 is 2x faster than V0 and 25x faster than V2b.

### Caller pre-populates all originators

The node knows all active originators and fills in `seq_id=0` for any originator missing from the client's cursor. This eliminates the need for a UNION ALL branch to handle unknown originators.

### No topic ordering in output

Results are ordered by `(originator_node_id, originator_sequence_id)` only. Topic ordering is not required.

## Performance Summary

All timings: 10,000 rows, 3 originators, 1000 topics, JIT disabled.

### Single partition per originator (6 total partitions)

| Variant | Plan | Exec | Total | Bottleneck |
|---------|------|------|-------|------------|
| V0 original | 6.60ms | 3.95ms | 10.5ms | O(N) planning time |
| V1 blob outside | 0.22ms | 419.6ms | 419.8ms | CTE scan O(N^2) + blob cross-join |
| V1b blob inside | 0.28ms | 2294.8ms | 2295.0ms | 10M blob rows scanned |
| V2 blob outside | 0.30ms | 186.2ms | 186.5ms | Blob cross-join (5.3M filtered) |
| V2b blob inside | 0.24ms | 8.35ms | 8.6ms | Near-optimal at 1 partition |
| **V3 meta+deferred blob** | **0.22ms** | **5.50ms** | **5.7ms** | **Near-optimal** |

### Multi-partition performance (1000 topics, scaled cursors)

| Parts/originator | V0 Total | V2b Total | V3 Total |
|---|---|---|---|
| 1 (6 total) | 10.5ms | 8.6ms | 5.7ms |
| 10 (34 total) | 47.1ms | 499ms | 19.7ms |
| 50 (154 total) | 80.6ms | 1009ms | 322ms |

### Multi-partition with gem_topic_seq_idx (1000 topics)

| Parts/originator | V0 Total | V2b Total | V3 Total |
|---|---|---|---|
| 10 (34 total) | 51.8ms | 615ms | **25.2ms** |
| 50 (154 total) | 109ms | 990ms | 347ms |

### Blob size impact (1 partition, 1000 topics)

| Blob Size | V0 Total | V2b Total | V3 Total |
|---|---|---|---|
| 256B | 11.7ms | 8.9ms | 6.0ms |
| 10KB | 10.6ms | 9.8ms | 11.6ms |
| 100KB | 10.5ms | 9.4ms | 6.6ms |

Blob payload size has minimal impact — all variants use PK lookups or hash joins on blob keys, so TOAST overhead is negligible at these result set sizes.

## Detailed Findings Per Query Variant

All timings from 10,000 rows, 3 originators, 1000 topics, JIT disabled unless noted.

### V0: Original query (baseline)

**Execution: 3.49ms | Planning: 6.12ms | Total: 9.61ms** (1 partition)

```sql
-- Pattern: CTE cursors + ANY($topics) + UNION ALL (known/unknown originators) + blob hash join
WHERE m.topic = ANY ($4::BYTEA[])
  AND m.originator_node_id = ANY($1::INT[])
```

The planner generates a bitmap index scan with all topic conditions OR'd together, then hash-joins with blobs (3097kB hash table). Execution is fast because a single scan covers all topics. However, planning time grows linearly with topic count — the planner must evaluate each topic literal. At 1000 topics this accounts for 64% of total time. At higher topic counts or larger tables, `ANY(large_array)` degrades to sequential scans.

Planning time also grows with partition count: 6.6ms at 1 part, 42ms at 10 parts, 79ms at 50 parts.

**Strengths**: Fast execution, efficient hash join for blobs, scales well with partitions.
**Weaknesses**: O(N) planning time, shared cursor prevents per-topic positioning, seq scans at scale.

### V1: LATERAL per topic, scalar subquery cursor, blob join outside

**Execution: 419.6ms | Planning: 0.22ms | Total: 419.8ms**

Two compounding problems:

1. **Correlated CTE scan**: The scalar subquery scans the full cursor_entries CTE (3000 rows) for EACH meta row. With 3 partitions per topic and 1000 topics = 3017 invocations, each scanning 3000 rows. Total: ~9M CTE rows scanned.

2. **Blob cross-join**: The blob join sits outside the LATERAL. Postgres materializes all 2467 LATERAL result rows, then for each of 2157 blob rows, checks against all materialized rows. Result: 5,319,089 rows removed by join filter.

**Verdict**: Disqualified. Correlated subquery makes execution O(N * M * cursor_size).

### V1b: LATERAL per topic, scalar subquery cursor, blob join inside

**Execution: 2294.8ms | Planning: 0.28ms | Total: 2295.0ms**

Moving the blob join inside the LATERAL made things dramatically worse. Without fixed `originator_node_id` in the LATERAL scope, all blob partitions are scanned every iteration. Total: 10M blob rows scanned.

**Verdict**: Worst performer. Disqualified.

### V2: LATERAL per (topic, originator) pair, blob join outside

**Execution: 186.2ms | Planning: 0.30ms | Total: 186.5ms**

The inner LATERAL is near-optimal — all three filter values are constants per iteration, enabling direct index seeks. However, the outer blob cross-join destroys performance: 5,319,245 rows removed by join filter.

**Verdict**: Proves the LATERAL structure is correct, but blob must not be joined outside.

### V2b: LATERAL per (topic, originator) pair, blob join inside

**At 1 partition: 9.75ms exec | 0.31ms plan | 10.0ms total** — Near-optimal.
**At 10 partitions: 499ms** — Catastrophic regression.
**At 50 partitions: 1009ms** — Unusable.

At 1 partition, `b.originator_node_id = ce.node_id` enables partition pruning on blobs, making each iteration a single PK lookup. But with multiple RANGE subpartitions, the planner cannot prune by `originator_sequence_id` (it comes from a join, not a constant), so it scans ALL subpartitions per iteration.

At 50 partitions: 3000 iterations * 50 subpartitions * ~67 rows/partition = ~10M seq scan rows.

**Verdict**: Excellent at 1 partition, disqualified at production partition counts.

### V3: LATERAL on meta only, deferred blob join (RECOMMENDED)

**At 1 partition: 5.50ms exec | 0.22ms plan | 5.7ms total**
**At 10 partitions: 18.7ms exec | 1.0ms plan | 19.7ms total**
**At 50 partitions: 316ms exec | 5.6ms plan | 322ms total**

V3 separates the meta LATERAL from the blob join:
1. LATERAL per (topic, originator) on meta only — fast index seeks
2. CTE with sort + LIMIT — materializes at most `row_limit` rows
3. Single blob join for limited results — 500 PK lookups via Memoize

At 1 partition, V3 is faster than V2b because it avoids per-iteration blob overhead entirely. At 10 partitions, V3 is 25x faster than V2b. At 50 partitions, V3 still degrades (meta LATERAL must seq-scan all subpartitions after cursor), but the blob join is negligible (2ms of 316ms).

**Verdict**: Best overall. O(1) planning, efficient at realistic partition counts (10/originator). The 50-partition scenario (150 total) is unrealistic for production — 100M rows with 1M bands = ~10 partitions per originator.

## Partition Sensitivity Analysis

The dominant factor is the number of RANGE subpartitions per originator. With two-level partitioning (LIST by originator, RANGE by sequence_id in 1M bands):

- **1 partition/originator**: All LATERAL variants perform well. Planner uses index scans.
- **10 partitions/originator**: V3 excels (19.7ms), V0 degrades on planning (47ms), V2b collapses (499ms).
- **50 partitions/originator**: Only V0 remains competitive on execution, but planning time is 79ms. V3 degrades to 322ms because each LATERAL iteration seq-scans all remaining subpartitions.

In production (100M rows, 10 originators, 1M bands): each originator has ~10M rows = ~10 RANGE subpartitions. This is the sweet spot for V3, where it's 2x faster than V0.

## Scaling Projections

For production (100M rows, 10 originators, ~10 partitions/originator):
- 1000 topics * 10 originators = 10,000 LATERAL iterations
- Each iteration: index seek on meta partition (~0.002ms)
- CTE construction: ~2ms
- Sort + limit: ~0.5ms
- Blob join (500 rows): ~1ms
- Estimated total: ~25ms
- V0 at same scale: ~50ms+ (planning time dominant)

---

## V0 Variant Exploration (Single-Cursor Optimization)

Separately from the per-topic cursor work (V1–V3), we explored optimizations to V0 while keeping the **single cursor shared across all topics** constraint. The goal: eliminate the O(N) planning time caused by `ANY(large_array)` without changing the cursor model.

### Variants Tested

**V0** (baseline): `ANY($topics)` + UNION ALL for unknown originators.

**V0b**: Replace `ANY($topics)` with a CTE `topic_list` joined via `unnest()`. Keeps UNION ALL.

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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

**V0c**: CTE `topic_list` + no UNION ALL. Caller pre-populates all originators in the cursor (seq_id=0 for unknown). Eliminates the NOT EXISTS anti-join branch.

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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

**V0d**: LATERAL per originator with CTE `topic_list`. Only 3–10 iterations (one per originator) instead of O(N) topic expansion. Each LATERAL subquery filters by originator + sequence cursor, joined with the topic CTE inside.

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
ORDER BY f.originator_node_id, f.originator_sequence_id;
```

### V0 Variant Performance (1000 topics, scaled cursors)

| Parts/originator | V0 Total | V0b Total | V0c Total | V0d Total |
|---|---|---|---|---|
| 1 (6 total) | 9.0ms | 4.4ms | 3.4ms | **3.2ms** |
| 10 (34 total) | 45.3ms | 6.3ms | 5.4ms | **3.2ms** |
| 50 (154 total) | 96.4ms | 18.7ms | 17.0ms | **8.7ms** |

### V0 Variant Analysis

**V0b** (CTE topics, UNION ALL): Replacing `ANY(topics)` with a CTE join eliminates the O(N) planning regression at 1000 topics. Planning drops from 40ms to 1.9ms at 10 partitions. However, the UNION ALL branch adds unnecessary overhead, and at 50 partitions the CTE join itself grows in planning cost (12ms).

**V0c** (CTE topics, no UNION ALL): Removing the UNION ALL branch gives a modest improvement over V0b. Requires the caller to pre-populate all originators in the cursor. Simpler plan, but still uses hash/merge joins that scan more data than necessary.

**V0d** (LATERAL per originator): Winner on small datasets (10K rows), but **catastrophically fails at production scale** — see production validation below. Each LATERAL iteration must scan all post-cursor rows for one originator in sequence order, filtering by topic via join. With millions of rows per originator and sparse topic matches, this degrades to full sequential scans.

At 10K rows V0d achieves O(1) planning with respect to topics and enables LIST partition pruning. But these advantages are irrelevant when the inner scan touches millions of rows.

---

## Production Validation (~9M rows, 7 originators, 22 leaf partitions)

Testing against a production-scale database reveals that local 10K-row benchmarks were misleading for V0d. The database has 7 originators (0, 1, 10, 11, 13, 100, 200) with 30K–3.5M rows each, 22 leaf RANGE partitions, and cursors positioned at ~80% of each originator's max sequence.

### Results (1000 topics, JIT disabled)

| Variant | Planning | Execution | Total | vs V0 |
|---|---|---|---|---|
| V0 (original) | 9.3ms | 990ms | **999ms** | baseline |
| V3 (per-topic cursor) | 6.8ms | 443ms | **450ms** | 2.2× faster |
| **V3b (per-originator blob)** | **7.2ms** | **180ms** | **187ms** | **5.3× faster** |
| V0d (LATERAL/originator) | — | >120,000ms | **timeout** | Disqualified |

### V0 Breakdown (999ms)

The meta filtering (416ms) breaks into two parts:
- **First UNION ALL branch** (275ms): Scans all 22 partitions with `ANY(1000 topics)`. Index-only scans on `gem_topic_time_idx` find 8466 matching rows, then hash-joins with cursors → 1452 rows after cursor filter (7014 removed by join filter).
- **Second UNION ALL branch** (141ms): Rescans all 22 partitions for the same 1000 topics, then anti-joins against cursors → **0 rows**. All 7 originators are known, making this branch pure waste.

The **blob join** (572ms) dominates: 500 PK lookups at ~1.14ms each via nested loop with runtime partition pruning. Only partitions matching each row's originator are actually scanned.

### V3 Breakdown (450ms)

The **meta LATERAL** (49ms) iterates over 7000 (topic, originator) pairs. Each iteration performs index-only scans using `(topic, originator_node_id, originator_sequence_id)` index conditions. Most iterations return 0 rows instantly (topic doesn't exist for that originator after cursor). Only 833 total rows found across all iterations. Runtime partition pruning skips irrelevant partitions — most show "never executed."

The **blob join** (393ms) fetches 500 rows at ~0.78ms each, similar to V0 but slightly faster due to different partition distribution of results.

### V3b Breakdown (187ms)

V3b wraps the blob join in a per-originator LATERAL, providing `originator_node_id` as a constant for partition pruning. The planner chose a Merge Join on `originator_node_id` between `filtered` and `originator_ids`, then fed the sorted result into the blob nested loop.

The **meta LATERAL** (48ms) is unchanged from V3.

The **blob join** (130ms) is 3× faster than V3. Per-lookup times dropped from 0.68–1.4ms to 0.19–0.63ms. The Merge Join sorts result rows by originator, so all blob lookups for originator 0 happen consecutively, then all for originator 10, etc. This maximizes cache locality on the partitioned blob table — each originator's pages stay hot for all its lookups instead of being evicted between interleaved lookups for other originators.

```sql
WITH cursor_entries AS (
    SELECT t.topic, cd.node_id, cd.seq_id
    FROM topics AS t
    CROSS JOIN (
        SELECT x.node_id, y.seq_id
        FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
        JOIN unnest(@cursor_seq_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord) USING (ord)
    ) AS cd
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
        LIMIT @rows_per_entry
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT @row_limit
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
ORDER BY bl.originator_node_id, bl.originator_sequence_id;
```

### Why V0d Fails at Scale

V0d's access pattern is fundamentally wrong for large datasets. Each of 7 LATERAL iterations must:

1. Scan all rows for one originator after the cursor position (100K–700K rows per originator at 80% cursor)
2. Order by `originator_sequence_id` (requires full sort or index scan)
3. Filter each row against 1000 topics via hash join with `topic_list`
4. Stop after finding 500 matches

With sparse topic selectivity (~1000 topics out of potentially millions of distinct topics), the query scans through hundreds of thousands of rows per originator before finding enough matches. Total rows scanned: millions.

V3 avoids this by inverting the access pattern: it starts from (topic, originator) pairs and seeks directly via index. Each of 7000 iterations does a single index probe, most returning 0 rows in microseconds.

**Lesson**: Small-dataset benchmarks (10K rows) cannot detect sequential scan problems. At 10K rows, scanning everything is fast. At 9M rows, the same scan pattern is 10,000× slower. Always validate against production-scale data.

# Per-Topic Cursor Query Strategy: Summary of Findings

## Problem Statement

The current production query (`SelectGatewayEnvelopesByTopics`, called "V0") uses a **single vector clock** (one cursor per originator) shared across all subscribed topics. This means a client subscribed to 500 topics must re-fetch data it has already seen on fast-moving topics whenever a slow-moving topic needs to catch up. Per-topic cursors — where each topic maintains its own per-originator cursor position — eliminate this redundant re-fetching.

The question: can per-topic cursors be queried efficiently at production scale?

## Is it performant to introduce per-topic cursors?

**Yes**, with the right query and index. At 10M rows with production-realistic skew (76%/15%/7.6% across 3 originators, 50K topics, 500 subscribed):

| Cursor Position | V0 (status quo) | V3b + Index (per-topic) |
|---|---|---|
| 80% (nearly caught up) | 1,476ms | **23ms** |
| 20% (far behind) | 1,749ms | **23ms** |
| Mixed (half at 80%, half at 20%) | 1,880ms | **23ms** |

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

## Comparison to V0 (Status Quo)

| Dimension | V0 | V3b + Index |
|---|---|---|
| **Cursor model** | Shared (one per originator) | Per-topic (one per topic x originator) |
| **Meta access** | O(partitions) — one scan per partition with `ANY(topics)` | O(topics x originators) — one index probe per pair |
| **Planning time** | 9-35ms (grows with topic count due to `ANY(large_array)`) | 0.4-0.5ms (stable) |
| **Execution at 10M rows** | 1,024-1,845ms | 22-23ms |
| **Sensitivity to cursor position** | Moderate (44ms-416ms meta phase depending on cursor) | Minimal (22-23ms regardless) |
| **Redundant re-fetching** | Yes — shared cursor forces re-scanning already-seen data | No — each topic tracks its own position |
| **Index requirement** | Uses existing `gem_topic_time_idx` | Requires `gem_topic_orig_seq_idx` (but allows dropping 3 old indexes) |
| **Write overhead** | N/A | Marginal — net fewer indexes after cleanup |

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
|---|---|---|
| 80% | 321ms | 23ms |
| 20% | 603ms | 23ms |
| Mixed | 412ms | 22ms |

At low cursor positions, more rows exist past each cursor, so each LATERAL probe scans more data before hitting the LIMIT. With the covering index, the probe is an index-only range scan that starts exactly at the cursor position and reads forward — the cost is proportional to the number of rows *returned*, not the number of rows *after the cursor*. Without the index, the probe uses bitmap heap scans that degrade with the volume of post-cursor data.

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

- `docs/plans/2026-02-13-per-topic-cursor-query-design.md` — Initial investigation at 10K and ~9M rows
- `docs/plans/2026-02-16-scaled-query-investigation-design.md` — 10M-row test design and results
- `docs/plans/2026-02-16-scaled-query-investigation-plan.md` — Implementation plan for the test harness

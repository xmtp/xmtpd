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

## 10M Row Results

### Test Configuration

- **Database**: PostgreSQL in Docker (local), 10M rows, 3 originators (76%/15%/7.6% skew)
- **50K topics**, 500 subscribed, power-law frequency distribution
- **14 RANGE partitions** per table (meta and blobs)
- **Warm cache**: 3 warm-up passes before each measured EXPLAIN ANALYZE
- **JIT disabled** for consistent timing

### Without Extra Index

| Cursor | Variant | Plan (ms) | Exec (ms) | Notes |
|--------|---------|-----------|-----------|-------|
| **80%** | V0_baseline | 9.5 | 1,466 | Production query |
| | V0c_no_union | 1.0 | 1,024 | CTE topics helps planning, not execution |
| | **V3b_lateral** | **0.4** | **321** | **Best per-topic cursor variant** |
| | V4_hybrid | 0.7 | 1,030 | Coarse scan too expensive |
| | V5_temp_table | 0.6 | 352 | Close to V3b |
| | V6_lateral_orig | 0.4 | 223 | Best shared-cursor, but no per-topic |
| **20%** | V0_baseline | 9.3 | 1,739 | More data to scan |
| | V0c_no_union | 1.0 | 1,162 | |
| | **V3b_lateral** | **0.4** | **603** | Still fastest per-topic |
| | V4_hybrid | 0.7 | 1,335 | Floor cursor too low, scans too much |
| | V5_temp_table | 0.6 | 1,104 | Worse without index |
| | V6_lateral_orig | 0.5 | 748 | |
| **Mixed** | V0_baseline | 34.5 | 1,845 | Planning time spike |
| | V0c_no_union | 1.0 | 1,191 | |
| | **V3b_lateral** | **0.5** | **412** | Per-topic cursors shine here |
| | V4_hybrid | 0.8 | 1,328 | |
| | V5_temp_table | 0.6 | 718 | |
| | V6_lateral_orig | 0.5 | 759 | |

### With `gem_topic_orig_seq_idx` Index

Index: `(topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time)`

| Cursor | Variant | Plan (ms) | Exec (ms) | vs. No Index |
|--------|---------|-----------|-----------|--------------|
| **80%** | V0_with_idx | 9.9 | 1,321 | 1.1× faster |
| | **V3b_with_idx** | **0.5** | **22.9** | **14× faster** |
| | V4_with_idx | 0.9 | 971 | 1.1× faster |
| | V6_with_idx | 0.5 | 264 | ~same |
| **20%** | V0_with_idx | 9.7 | 1,680 | ~same |
| | **V3b_with_idx** | **0.5** | **22.5** | **27× faster** |
| | V4_with_idx | 0.9 | 1,417 | ~same |
| | V6_with_idx | 0.5 | 810 | ~same |
| **Mixed** | V0_with_idx | 10.1 | 1,628 | ~same |
| | **V3b_with_idx** | **0.5** | **22.4** | **18× faster** |
| | V4_with_idx | 0.9 | 1,279 | ~same |
| | V6_with_idx | 0.5 | 700 | ~same |

### Correctness Verification

| Check | Result |
|-------|--------|
| V0c matches V0 | ✅ Identical rows |
| V3b ≥ V0 rows | ✅ Superset confirmed |
| V4 capture rate (80% uniform) | ✅ 100% (500/500) |
| V5 matches V3b | ✅ Identical rows |

### Analysis

#### 1. V3b + Index is the clear winner

V3b with `gem_topic_orig_seq_idx` achieves **22-23ms execution** regardless of cursor position — a **14-27× improvement** over V3b without the index, and **57-75× faster** than V0 baseline.

The index enables **index-only scans** on each partition's LATERAL subquery, with zero heap fetches. The covering `INCLUDE (gateway_time)` is critical — without it, each per-(topic, originator) probe would need a heap fetch.

#### 2. V4 (hybrid) does not achieve its goal

V4 was designed to combine V0's efficient single-pass scan with per-topic cursor precision. In practice, the coarse scan (using floor cursors) scans too many rows — the floor cursor is the minimum across all topics, so the 3× over-fetch limit is insufficient to cover the gap. V4's execution time (~1,000-1,300ms) is comparable to V0c and doesn't benefit from the per-topic precision.

#### 3. V5 (temp table) is viable but V3b+index is better

V5 achieves good results at 80% cursors (352ms, close to V3b's 321ms) but degrades at 20% cursors (1,104ms). The temp table overhead (CREATE + INSERT + ANALYZE) doesn't pay for itself when V3b+index does 22ms directly.

#### 4. V6 (LATERAL per originator) works but is shared-cursor only

V6 (reusing queryV0dSQL) performs well (223-759ms) but only supports shared cursors, not per-topic. It's competitive for the shared-cursor use case but doesn't solve the per-topic cursor problem.

#### 5. V0 baseline planning time degrades with mixed cursors

V0's planning time spiked to 34.5ms with mixed cursors (vs 9.5ms at 80%). V0c eliminates this (~1ms consistently) but execution remains slow.

### Recommendation

**Deploy V3b + `gem_topic_orig_seq_idx`** for the per-topic cursor query path:

- **22ms execution** is fast enough for real-time subscription queries
- **0.5ms planning** is negligible
- The index adds write overhead but the covering columns make it index-only (no heap fetches)
- Per-topic cursors eliminate redundant re-fetching, reducing bandwidth
- Works consistently across all cursor positions (80%, 20%, mixed)

**Action items:**
1. Add migration to create `gem_topic_orig_seq_idx` on `gateway_envelopes_meta`
2. Implement V3b as the production per-topic cursor query in `envelopes_v2.sql`
3. Monitor index size and write amplification in production
4. Drop redundant indexes (see analysis below)

### Index Drop Analysis

When adding `gem_topic_orig_seq_idx`, three existing indexes on `gateway_envelopes_meta` were evaluated for removal. Every `.sql` file in `pkg/db/sqlc/` and all Go files with raw SQL were audited.

#### Queries Touching `gateway_envelopes_meta`

| # | Query | File | WHERE Filters | ORDER BY |
|---|-------|------|---------------|----------|
| 1 | `SelectGatewayEnvelopesByTopics` (Branch A) | envelopes_v2.sql:108 | `topic = ANY(...)`, `originator_node_id = ANY(...)`, `originator_sequence_id > cursor` | `originator_node_id, originator_sequence_id` |
| 2 | `SelectGatewayEnvelopesByTopics` (Branch B) | envelopes_v2.sql:138 | `topic = ANY(...)`, `originator_sequence_id > 0` | `originator_node_id, originator_sequence_id` |
| 3 | `SelectNewestFromTopics` | envelopes_v2.sql:36 | `topic = ANY(...)` | `topic, gateway_time DESC` |
| 4 | `SelectGatewayEnvelopesBySingleOriginator` | envelopes_v2.sql:55 | `originator_node_id = @val`, `originator_sequence_id > cursor` | `originator_sequence_id` |
| 5 | `SelectGatewayEnvelopesByOriginators` (LATERAL) | envelopes_v2.sql:72 | `originator_node_id = o.node_id`, `originator_sequence_id > cursor` | `originator_sequence_id` |
| 6 | `SelectGatewayEnvelopesUnfiltered` | envelopes_v2.sql:166 | `originator_sequence_id > cursor` (via view) | `originator_node_id, originator_sequence_id` |
| 7 | `CountExpiredEnvelopes` | prune.sql:1 | `expiry IS NOT NULL`, `expiry < now()`, `originator_sequence_id <= val` | — |
| 8 | `DeleteExpiredEnvelopesBatch` | prune.sql:16 | `expiry IS NOT NULL`, `expiry < now()`, `originator_sequence_id <= val` | `expiry, originator_node_id, originator_sequence_id` |
| 9 | `CountExpiredMigratedEnvelopes` | prune.sql:40 | `expiry IS NOT NULL`, `expiry < now()`, `originator_node_id BETWEEN 10 AND 14` | — |
| 10 | `DeleteExpiredMigratedEnvelopesBatch` | prune.sql:47 | same as #9 | `expiry, originator_node_id, originator_sequence_id` |
| 11 | `InsertGatewayEnvelope` | envelopes_v2.sql:1 | — (INSERT) | — |
| 12 | `InsertGatewayEnvelopeBatch...` | envelopes_v2.sql:184 | — (INSERT via function) | — |

No other `.sql` files or Go files contain queries against `gateway_envelopes_meta`.

#### 1. `gem_topic_time_idx` — CAN DROP

**Definition:** `(topic, gateway_time, originator_node_id, originator_sequence_id)`

**Only consumer:** `SelectGatewayEnvelopesByTopics` (queries #1, #2) — the migration comment explicitly says "required for SelectGatewayEnvelopesByTopics".

**Why the new index is strictly better for this query:**
- Branch A filters on `topic`, `originator_node_id`, `originator_sequence_id` — all three are key columns in `gem_topic_orig_seq_idx`, in the correct order for range scan. The old index has `gateway_time` as the 2nd key column despite `gateway_time` never appearing in any WHERE clause, forcing the planner to skip over it.
- Branch B filters on `topic` and `originator_sequence_id > 0` — both indexes lead with `topic`. The new index reaches `originator_sequence_id` as the 3rd column vs 4th in the old index.
- ORDER BY is `(originator_node_id, originator_sequence_id)` which matches the new index's key order, not the old index's.
- `gateway_time` is only needed in SELECT, not WHERE/ORDER — the new index provides it via INCLUDE.

**Other queries:** `SelectNewestFromTopics` (query #3) uses `gem_topic_time_desc_idx` (topic, gateway_time DESC), not `gem_topic_time_idx`.

**Verdict:** Drop. The new `gem_topic_orig_seq_idx` is a direct upgrade for every query that used this index.

#### 2. `gem_time_node_seq_idx` — CAN DROP

**Definition:** `(gateway_time, originator_node_id, originator_sequence_id)`

**Migration comment:** "required for most gateway sorted selects"

**Audit finding: no current query uses `gateway_time` as a leading filter.** Every query on `gateway_envelopes_meta` filters first by `topic`, `originator_node_id`, or `expiry` — never by `gateway_time` alone. Specifically:
- Queries #1-3 lead with `topic`
- Queries #4-5 lead with `originator_node_id`
- Query #6 leads with `originator_sequence_id`
- Queries #7-10 lead with `expiry` (served by `gem_expiry_idx`)
- `SelectNewestFromTopics` uses `ORDER BY gateway_time DESC` but with a `topic` filter, so it uses `gem_topic_time_desc_idx`

The migration comment appears to be legacy — the query patterns evolved to use topic-leading or originator-leading access, making this index dead weight.

**Verdict:** Drop. No query in the codebase benefits from a `gateway_time`-leading index. Confirm no external consumers before deploying.

#### 3. `gem_originator_node_id` — CAN DROP

**Definition:** `(originator_node_id)` — single-column index

**Why redundant:**
- The table is **LIST-partitioned** by `originator_node_id`, so any query specifying `originator_node_id` in WHERE automatically prunes to just that partition — no index needed.
- Within a partition, the **PRIMARY KEY** `(originator_node_id, originator_sequence_id)` provides the index for range scans on `originator_sequence_id`.
- Queries #4-5 filter by `originator_node_id` equality + `originator_sequence_id` range — partition pruning + PK covers this exactly.
- No query does `SELECT DISTINCT originator_node_id FROM gateway_envelopes_meta` — the `gateway_envelopes_latest` table serves that purpose.
- Queries #9-10 use `originator_node_id BETWEEN 10 AND 14` but `gem_expiry_idx` is the primary access path (filters by `expiry` first).

**Verdict:** Drop. Partition pruning + PK fully subsumes this index for every existing query.

#### Summary

| Index | Verdict | Reason |
|-------|---------|--------|
| `gem_topic_time_idx` | **DROP** | New `gem_topic_orig_seq_idx` is strictly better for all consumers |
| `gem_time_node_seq_idx` | **DROP** | No query uses `gateway_time` as leading filter; legacy comment |
| `gem_originator_node_id` | **DROP** | Redundant with LIST partitioning + PRIMARY KEY |
| `gem_topic_time_desc_idx` | **KEEP** | Required by `SelectNewestFromTopics` (ORDER BY topic, gateway_time DESC) |
| `gem_expiry_idx` | **KEEP** | Required by all pruning queries; unrelated to subscription path |

**Net result:** Adding 1 index (`gem_topic_orig_seq_idx`) and dropping 3 indexes reduces total index count from 5 to 3, lowering write amplification while dramatically improving read performance.

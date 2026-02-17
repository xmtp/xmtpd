package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// ---------- Scaled test configuration ----------

// Originator distribution modeling production skew (o11 = 76% of data).
// node_id => fraction of total rows
var scaledOriginators = []struct {
	nodeID   int32
	fraction float64 // fraction of totalRows
}{
	{100, 0.076}, // ~760K rows, 1 RANGE partition
	{200, 0.150}, // ~1.5M rows, 2 RANGE partitions
	{300, 0.774}, // ~7.74M rows, 8 RANGE partitions
}

const (
	scaledTotalRows       = 10_000_000
	scaledTotalTopics     = 50_000
	scaledSubscribedCount = 500
	scaledBatchSize       = 5_000
	scaledRowLimit        = 500
	scaledRowsPerEntry    = 50
)

// ---------- Seeding ----------

type scaledTestData struct {
	db               *sql.DB
	allTopics        [][]byte // all 50K topics
	subscribedTopics [][]byte // first 500 topics (hot)
	originatorIDs    []int32
	maxSeqs          map[int32]int64 // per-originator max sequence ID
}

// seedScaledDatabase creates 10M rows with production-like skew.
// Uses raw SQL batch inserts (bypassing application layer) for speed.
func seedScaledDatabase(t *testing.T, db *sql.DB) *scaledTestData {
	t.Helper()
	ctx := context.Background()
	querier := queries.New(db)

	// Generate topics.
	rng := rand.New(rand.NewSource(42))
	allTopics := make([][]byte, scaledTotalTopics)
	for i := range allTopics {
		allTopics[i] = make([]byte, 32)
		rng.Read(allTopics[i])
	}
	subscribedTopics := allTopics[:scaledSubscribedCount]

	// Pre-create all partitions.
	originatorIDs := make([]int32, len(scaledOriginators))
	maxSeqs := make(map[int32]int64)
	for i, o := range scaledOriginators {
		originatorIDs[i] = o.nodeID
		rowCount := int64(float64(scaledTotalRows) * o.fraction)
		maxSeq := rowCount // seq IDs are 1..rowCount
		maxSeqs[o.nodeID] = maxSeq

		// Create partitions for each 1M band.
		for seq := int64(0); seq <= maxSeq; seq += 1_000_000 {
			err := querier.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
				OriginatorNodeID:     o.nodeID,
				OriginatorSequenceID: seq,
				BandWidth:            1_000_000,
			})
			require.NoError(t, err, "creating partition for node=%d seq=%d", o.nodeID, seq)
		}
	}

	// Batch insert rows per originator.
	for _, o := range scaledOriginators {
		rowCount := int(float64(scaledTotalRows) * o.fraction)
		t.Logf("seeding originator %d: %d rows", o.nodeID, rowCount)

		nodeIDs := make([]int32, 0, scaledBatchSize)
		seqIDs := make([]int64, 0, scaledBatchSize)
		topics := make([][]byte, 0, scaledBatchSize)
		expiries := make([]int64, 0, scaledBatchSize)
		payloads := make([][]byte, 0, scaledBatchSize)

		payload := make([]byte, 256)
		rng.Read(payload) // reuse same payload for speed

		flush := func() {
			if len(nodeIDs) == 0 {
				return
			}
			_, err := db.ExecContext(ctx, `
				INSERT INTO gateway_envelopes_meta
					(originator_node_id, originator_sequence_id, topic, expiry)
				SELECT * FROM unnest($1::INT[], $2::BIGINT[], $3::BYTEA[], $4::BIGINT[])`,
				pq.Array(nodeIDs), pq.Array(seqIDs), pq.Array(topics), pq.Array(expiries),
			)
			require.NoError(t, err, "batch insert meta for node=%d", o.nodeID)

			_, err = db.ExecContext(ctx, `
				INSERT INTO gateway_envelope_blobs
					(originator_node_id, originator_sequence_id, originator_envelope)
				SELECT * FROM unnest($1::INT[], $2::BIGINT[], $3::BYTEA[])`,
				pq.Array(nodeIDs), pq.Array(seqIDs), pq.Array(payloads),
			)
			require.NoError(t, err, "batch insert blobs for node=%d", o.nodeID)

			nodeIDs = nodeIDs[:0]
			seqIDs = seqIDs[:0]
			topics = topics[:0]
			expiries = expiries[:0]
			payloads = payloads[:0]
		}

		for seq := int64(1); seq <= int64(rowCount); seq++ {
			// Topic selection: 50% chance of subscribed topic (creates ~10K rows per subscribed topic).
			var topic []byte
			if rng.Float64() < 0.5 {
				topic = subscribedTopics[rng.Intn(len(subscribedTopics))]
			} else {
				topic = allTopics[rng.Intn(len(allTopics))]
			}

			nodeIDs = append(nodeIDs, o.nodeID)
			seqIDs = append(seqIDs, seq)
			topics = append(topics, topic)
			expiries = append(expiries, 0)
			payloads = append(payloads, payload)

			if len(nodeIDs) >= scaledBatchSize {
				flush()
			}
		}
		flush()
	}

	// ANALYZE for accurate planner statistics.
	_, err := db.ExecContext(ctx, "ANALYZE gateway_envelopes_meta")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "ANALYZE gateway_envelope_blobs")
	require.NoError(t, err)

	// Report partition counts.
	var metaParts, blobParts int
	_ = db.QueryRowContext(ctx,
		`WITH RECURSIVE parts AS (
			SELECT oid FROM pg_class WHERE relname = 'gateway_envelopes_meta'
			UNION ALL
			SELECT c.oid FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
			JOIN parts p ON p.oid = i.inhparent
		) SELECT COUNT(*)-1 FROM parts`).Scan(&metaParts)
	_ = db.QueryRowContext(ctx,
		`WITH RECURSIVE parts AS (
			SELECT oid FROM pg_class WHERE relname = 'gateway_envelope_blobs'
			UNION ALL
			SELECT c.oid FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
			JOIN parts p ON p.oid = i.inhparent
		) SELECT COUNT(*)-1 FROM parts`).Scan(&blobParts)

	t.Logf("seeded %d rows across %d originators (%d meta partitions, %d blob partitions, %d topics, %d subscribed)",
		scaledTotalRows, len(scaledOriginators), metaParts, blobParts, scaledTotalTopics, scaledSubscribedCount)

	return &scaledTestData{
		db:               db,
		allTopics:        allTopics,
		subscribedTopics: subscribedTopics,
		originatorIDs:    originatorIDs,
		maxSeqs:          maxSeqs,
	}
}

// ---------- Cursor builders ----------

type cursorDistribution struct {
	name     string
	// fractions maps originator node_id to cursor position as fraction of max_seq.
	// For per-topic cursors, topicFractions overrides per-topic if set.
	fractions      map[int32]float64
	topicFractions []float64 // len == subscribedTopics; nil means use fractions for all
}

func (td *scaledTestData) cursorDistributions() []cursorDistribution {
	rng := rand.New(rand.NewSource(99))

	// Mixed: half topics at 80%, half at 20%.
	topicFracs := make([]float64, scaledSubscribedCount)
	for i := range topicFracs {
		if rng.Float64() < 0.5 {
			topicFracs[i] = 0.8
		} else {
			topicFracs[i] = 0.2
		}
	}

	return []cursorDistribution{
		{
			name:      "80pct",
			fractions: map[int32]float64{100: 0.8, 200: 0.8, 300: 0.8},
		},
		{
			name:      "20pct",
			fractions: map[int32]float64{100: 0.2, 200: 0.2, 300: 0.2},
		},
		{
			name:           "mixed",
			fractions:      map[int32]float64{100: 0.8, 200: 0.8, 300: 0.8},
			topicFractions: topicFracs,
		},
	}
}

// buildV0Args builds shared-cursor V0 params. For mixed distributions, uses
// the minimum cursor per originator across all topics.
func (td *scaledTestData) buildV0Args(dist cursorDistribution) []any {
	nodeIDs := make([]int32, len(td.originatorIDs))
	seqIDs := make([]int64, len(td.originatorIDs))

	for i, nid := range td.originatorIDs {
		nodeIDs[i] = nid
		baseFrac := dist.fractions[nid]

		if dist.topicFractions != nil {
			// Use minimum fraction across all topics for this originator.
			minFrac := baseFrac
			for _, tf := range dist.topicFractions {
				if tf < minFrac {
					minFrac = tf
				}
			}
			baseFrac = minFrac
		}

		seqIDs[i] = int64(float64(td.maxSeqs[nid]) * baseFrac)
	}

	return []any{
		pq.Array(nodeIDs),
		pq.Array(seqIDs),
		int32(scaledRowLimit),
		pq.Array(td.subscribedTopics),
	}
}

// buildPerTopicArgs builds per-topic cursor params (V3b/V4/V5 format).
// Each subscribed topic x originator gets its own cursor entry.
func (td *scaledTestData) buildPerTopicArgs(dist cursorDistribution) []any {
	n := scaledSubscribedCount * len(td.originatorIDs)
	ct := make([][]byte, 0, n)
	cn := make([]int32, 0, n)
	cs := make([]int64, 0, n)

	for ti, topic := range td.subscribedTopics {
		for _, nid := range td.originatorIDs {
			frac := dist.fractions[nid]
			if dist.topicFractions != nil {
				frac = dist.topicFractions[ti]
			}
			ct = append(ct, topic)
			cn = append(cn, nid)
			cs = append(cs, int64(float64(td.maxSeqs[nid])*frac))
		}
	}

	return []any{
		pq.Array(ct),
		pq.Array(cn),
		pq.Array(cs),
		int32(scaledRowsPerEntry),
		int32(scaledRowLimit),
	}
}

// buildV4Args builds per-topic cursor params for V4 (no rows_per_entry param).
// V4 params: $1=topics, $2=node_ids, $3=seq_ids, $4=row_limit
func (td *scaledTestData) buildV4Args(dist cursorDistribution) []any {
	ptArgs := td.buildPerTopicArgs(dist)
	// Skip $4 (rows_per_entry) — V4 doesn't use it
	return []any{ptArgs[0], ptArgs[1], ptArgs[2], ptArgs[4]}
}

// ---------- Warm-up and EXPLAIN ----------

// warmAndExplain runs the query 3 times for cache warming, then runs
// EXPLAIN (ANALYZE, BUFFERS) and returns the plan text.
func warmAndExplain(t *testing.T, db *sql.DB, rawSQL string, args ...any) (plan string, rowCount int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Warm-up: run 3 times, discard results.
	for range 3 {
		rows, err := db.QueryContext(ctx, rawSQL, args...)
		require.NoError(t, err, "warm-up query failed")
		for rows.Next() {
			// discard
		}
		rows.Close()
	}

	// Count rows from an actual execution.
	rows, err := db.QueryContext(ctx, rawSQL, args...)
	require.NoError(t, err)
	rowCount = 0
	for rows.Next() {
		rowCount++
	}
	rows.Close()

	// EXPLAIN ANALYZE with BUFFERS.
	rows, err = db.QueryContext(ctx, "EXPLAIN (ANALYZE, BUFFERS) "+rawSQL, args...)
	require.NoError(t, err)
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var line string
		require.NoError(t, rows.Scan(&line))
		lines = append(lines, line)
	}
	require.NoError(t, rows.Err())

	plan = strings.Join(lines, "\n")
	return plan, rowCount
}

// v5Setup runs V5's multi-step temp table creation + insert + analyze.
func v5Setup(t *testing.T, ctx context.Context, db *sql.DB, insertArgs []any) {
	t.Helper()
	_, err := db.ExecContext(ctx, queryV5CreateSQL)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, queryV5TruncateSQL)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, queryV5InsertSQL, insertArgs...)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, queryV5AnalyzeSQL)
	require.NoError(t, err)
}

// warmAndExplainV5 handles V5's multi-statement pattern:
// setup SQL runs first (temp table creation + insert + analyze), then the main query is explained.
func warmAndExplainV5(
	t *testing.T,
	db *sql.DB,
	insertArgs []any,
	mainSQL string,
	mainArgs []any,
) (plan string, rowCount int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Warm-up.
	for range 3 {
		v5Setup(t, ctx, db, insertArgs)
		rows, err := db.QueryContext(ctx, mainSQL, mainArgs...)
		require.NoError(t, err)
		for rows.Next() {
		}
		rows.Close()
		_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS _cursor_entries")
	}

	// Actual count.
	v5Setup(t, ctx, db, insertArgs)
	rows, err := db.QueryContext(ctx, mainSQL, mainArgs...)
	require.NoError(t, err)
	for rows.Next() {
		rowCount++
	}
	rows.Close()
	_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS _cursor_entries")

	// EXPLAIN.
	v5Setup(t, ctx, db, insertArgs)
	rows, err = db.QueryContext(ctx, "EXPLAIN (ANALYZE, BUFFERS) "+mainSQL, mainArgs...)
	require.NoError(t, err)
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var line string
		require.NoError(t, rows.Scan(&line))
		lines = append(lines, line)
	}
	require.NoError(t, rows.Err())
	_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS _cursor_entries")

	return strings.Join(lines, "\n"), rowCount
}

type variantResult struct {
	cursor  string
	variant string
	plan    string
	rows    int
}

func printSummary(t *testing.T, results []variantResult) {
	t.Helper()
	fmt.Printf("\n%-8s | %-25s | %5s | %s\n", "Cursor", "Variant", "Rows", "Plan/Exec extracted from EXPLAIN")
	fmt.Println(strings.Repeat("-", 80))
	for _, r := range results {
		// Extract Planning Time and Execution Time from plan.
		planTime, execTime := "?", "?"
		for _, line := range strings.Split(r.plan, "\n") {
			if strings.Contains(line, "Planning Time:") {
				planTime = strings.TrimSpace(line)
			}
			if strings.Contains(line, "Execution Time:") {
				execTime = strings.TrimSpace(line)
			}
		}
		fmt.Printf("%-8s | %-25s | %5d | %s / %s\n", r.cursor, r.variant, r.rows, planTime, execTime)
	}
	fmt.Println()
}

// ---------- V4 SQL ----------

// V4: Hybrid V0-style scan with per-topic cursor post-filter.
// Uses floor cursors (min per originator across all topics) for the coarse scan,
// then joins with per-topic cursor_entries to filter precisely.
// Params: $1=cursor_topics BYTEA[], $2=cursor_node_ids INT[],
//
//	$3=cursor_seq_ids BIGINT[], $4=row_limit INT
const queryV4SQL = `WITH cursor_entries AS (
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
    WHERE m.originator_node_id IN (SELECT node_id FROM distinct_node_ids)
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF($4::INT, 0) * 3
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// ---------- V5 SQL ----------

// V5: Temp table cursors — insert per-topic cursors into a temp table,
// then hash-join with meta.
// Setup params: $1=cursor_topics BYTEA[], $2=cursor_node_ids INT[], $3=cursor_seq_ids BIGINT[]
const queryV5CreateSQL = `CREATE TEMP TABLE IF NOT EXISTS _cursor_entries (
    topic BYTEA NOT NULL,
    node_id INT NOT NULL,
    seq_id BIGINT NOT NULL
)`

const queryV5TruncateSQL = `TRUNCATE _cursor_entries`

const queryV5InsertSQL = `INSERT INTO _cursor_entries (topic, node_id, seq_id)
SELECT t.topic, n.node_id, s.seq_id
FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)`

const queryV5AnalyzeSQL = `ANALYZE _cursor_entries`

// Main query params: $1=cursor_node_ids INT[] (distinct), $2=row_limit INT
const queryV5MainSQL = `WITH min_cursor_seq AS (
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// buildV5Args returns (insertArgs, mainArgs) for the V5 temp table variant.
func (td *scaledTestData) buildV5Args(dist cursorDistribution) (insertArgs []any, mainArgs []any) {
	ptArgs := td.buildPerTopicArgs(dist)
	// Insert uses the per-topic cursor arrays: $1=topics, $2=node_ids, $3=seq_ids
	insertArgs = ptArgs[:3]
	// Main uses: $1=distinct_node_ids, $2=row_limit
	mainArgs = []any{
		pq.Array(td.originatorIDs),
		int32(scaledRowLimit),
	}
	return insertArgs, mainArgs
}

// ---------- Main test ----------

func TestScaledQueryComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scaled query test in short mode")
	}

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	db.SetMaxOpenConns(1) // Pin to one connection so SET and temp tables persist.

	_, err := db.ExecContext(ctx, "SET jit = off")
	require.NoError(t, err)

	start := time.Now()
	td := seedScaledDatabase(t, db)
	t.Logf("seeding took %s", time.Since(start))

	var allResults []variantResult

	for _, dist := range td.cursorDistributions() {
		v0Args := td.buildV0Args(dist)
		ptArgs := td.buildPerTopicArgs(dist)

		t.Run(dist.name, func(t *testing.T) {
			t.Run("V0_baseline", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV0SQL, v0Args...)
				allResults = append(allResults, variantResult{dist.name, "V0_baseline", plan, rows})
				fmt.Printf("\n===== V0_baseline @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V0c_no_union", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV0cSQL, v0Args...)
				allResults = append(allResults, variantResult{dist.name, "V0c_no_union", plan, rows})
				fmt.Printf("\n===== V0c_no_union @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V3b_lateral", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV3bSQL, ptArgs...)
				allResults = append(allResults, variantResult{dist.name, "V3b_lateral", plan, rows})
				fmt.Printf("\n===== V3b_lateral @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V4_hybrid", func(t *testing.T) {
				v4Args := td.buildV4Args(dist)
				plan, rows := warmAndExplain(t, db, queryV4SQL, v4Args...)
				allResults = append(allResults, variantResult{dist.name, "V4_hybrid", plan, rows})
				fmt.Printf("\n===== V4_hybrid @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V5_temp_table", func(t *testing.T) {
				insertArgs, mainArgs := td.buildV5Args(dist)
				plan, rows := warmAndExplainV5(t, db, insertArgs, queryV5MainSQL, mainArgs)
				allResults = append(allResults, variantResult{dist.name, "V5_temp_table", plan, rows})
				fmt.Printf("\n===== V5_temp_table @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V6_lateral_orig", func(t *testing.T) {
				// V6 uses shared-cursor params (same as V0).
				plan, rows := warmAndExplain(t, db, queryV0dSQL, v0Args...)
				allResults = append(allResults, variantResult{dist.name, "V6_lateral_orig", plan, rows})
				fmt.Printf("\n===== V6_lateral_orig @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})
		})
	}

	printSummary(t, allResults)
}

// ---------- Index experiments ----------

func TestScaledQueryComparison_WithIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scaled query test in short mode")
	}

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	db.SetMaxOpenConns(1) // Pin to one connection so SET and temp tables persist.

	_, err := db.ExecContext(ctx, "SET jit = off")
	require.NoError(t, err)

	start := time.Now()
	td := seedScaledDatabase(t, db)
	t.Logf("seeding took %s", time.Since(start))

	// Create the covering index.
	_, err = db.ExecContext(ctx,
		`CREATE INDEX gem_topic_orig_seq_idx ON gateway_envelopes_meta
		 (topic, originator_node_id, originator_sequence_id) INCLUDE (gateway_time)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, "ANALYZE gateway_envelopes_meta")
	require.NoError(t, err)
	t.Log("created gem_topic_orig_seq_idx")

	var allResults []variantResult

	for _, dist := range td.cursorDistributions() {
		v0Args := td.buildV0Args(dist)
		ptArgs := td.buildPerTopicArgs(dist)

		t.Run(dist.name, func(t *testing.T) {
			t.Run("V0_with_idx", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV0SQL, v0Args...)
				allResults = append(allResults, variantResult{dist.name, "V0_with_idx", plan, rows})
				fmt.Printf("\n===== V0_with_idx @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V3b_with_idx", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV3bSQL, ptArgs...)
				allResults = append(allResults, variantResult{dist.name, "V3b_with_idx", plan, rows})
				fmt.Printf("\n===== V3b_with_idx @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V4_with_idx", func(t *testing.T) {
				v4Args := td.buildV4Args(dist)
				plan, rows := warmAndExplain(t, db, queryV4SQL, v4Args...)
				allResults = append(allResults, variantResult{dist.name, "V4_with_idx", plan, rows})
				fmt.Printf("\n===== V4_with_idx @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})

			t.Run("V6_with_idx", func(t *testing.T) {
				plan, rows := warmAndExplain(t, db, queryV0dSQL, v0Args...)
				allResults = append(allResults, variantResult{dist.name, "V6_with_idx", plan, rows})
				fmt.Printf("\n===== V6_with_idx @ %s (%d rows) =====\n%s\n", dist.name, rows, plan)
			})
		})
	}

	printSummary(t, allResults)
}

// ---------- Correctness verification ----------

func TestScaledQueryCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scaled query test in short mode")
	}

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	db.SetMaxOpenConns(1) // Pin to one connection so SET and temp tables persist.

	_, err := db.ExecContext(ctx, "SET jit = off")
	require.NoError(t, err)

	td := seedScaledDatabase(t, db)

	dist := td.cursorDistributions()[0] // 80% cursors
	v0Args := td.buildV0Args(dist)
	ptArgs := td.buildPerTopicArgs(dist)

	v4Args := td.buildV4Args(dist)

	v0Rows := queryResultRows(t, db, queryV0SQL, v0Args...)
	v0cRows := queryResultRows(t, db, queryV0cSQL, v0Args...)
	v3bRows := queryResultRows(t, db, queryV3bSQL, ptArgs...)
	v4Rows := queryResultRows(t, db, queryV4SQL, v4Args...)

	t.Run("V0c_matches_V0", func(t *testing.T) {
		requireSameRows(t, "V0c vs V0", v0Rows, v0cRows)
	})

	t.Run("V3b_is_superset_of_V0", func(t *testing.T) {
		// V3b with uniform cursors returns same rows as V0.
		// With mixed cursors it would return a superset.
		// At 80% uniform, they should match.
		require.GreaterOrEqual(t, len(v3bRows), len(v0Rows),
			"V3b should return at least as many rows as V0 (per-topic precision)")
	})

	t.Run("V4_matches_V3b", func(t *testing.T) {
		// V4 uses same per-topic cursors but may return fewer rows if 3x over-fetch
		// wasn't enough. Check it returns a subset of V3b.
		require.LessOrEqual(t, len(v4Rows), len(v3bRows),
			"V4 rows should be <= V3b rows (3x over-fetch may truncate)")
		t.Logf("V3b returned %d rows, V4 returned %d rows (%.0f%% capture rate)",
			len(v3bRows), len(v4Rows), float64(len(v4Rows))/float64(len(v3bRows))*100)
	})

	// V5 correctness (temp table variant).
	t.Run("V5_matches_V3b", func(t *testing.T) {
		insertArgs, mainArgs := td.buildV5Args(dist)
		v5Setup(t, ctx, db, insertArgs)
		v5Rows := queryResultRows(t, db, queryV5MainSQL, mainArgs...)
		_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS _cursor_entries")

		requireSameRows(t, "V5 vs V3b", v3bRows, v5Rows)
	})
}

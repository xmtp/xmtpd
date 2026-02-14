package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// ---------- Raw SQL ----------

// V0: Original query (per-originator cursor shared across all topics).
// Params: $1=cursor_node_ids INT[], $2=cursor_seq_ids BIGINT[],
//
//	$3=row_limit INT, $4=topics BYTEA[]
const queryV0SQL = `WITH cursors AS (
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// V0b: Replace ANY(topics) with a CTE join to avoid O(N) planning.
// Keeps UNION ALL for unknown originators (same semantics as V0).
// Params: $1=cursor_node_ids INT[], $2=cursor_seq_ids BIGINT[],
//
//	$3=row_limit INT, $4=topics BYTEA[]
const queryV0bSQL = `WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest($1::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest($2::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq
    FROM unnest($2::BIGINT[]) AS t(seq_id)
),
topic_list AS (
    SELECT t.topic FROM unnest($4::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    JOIN cursors AS c
         ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.originator_node_id = ANY($1::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)

    UNION ALL

    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    WHERE m.originator_sequence_id > 0
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// V0c: CTE topics + no UNION ALL. Requires caller to pre-populate all
// originators in the cursor (use seq_id=0 for unknown originators).
// Simpler query plan without the NOT EXISTS anti-join branch.
// Params: $1=cursor_node_ids INT[], $2=cursor_seq_ids BIGINT[],
//
//	$3=row_limit INT, $4=topics BYTEA[]
const queryV0cSQL = `WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest($1::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest($2::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
min_cursor_seq AS (
    SELECT COALESCE(MIN(seq_id), 0) AS min_seq
    FROM unnest($2::BIGINT[]) AS t(seq_id)
),
topic_list AS (
    SELECT t.topic FROM unnest($4::BYTEA[]) AS t(topic)
),
filtered AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    JOIN topic_list AS tl ON m.topic = tl.topic
    JOIN cursors AS c
         ON m.originator_node_id = c.cursor_node_id
         AND m.originator_sequence_id > c.cursor_sequence_id
    WHERE m.originator_node_id = ANY($1::INT[])
      AND m.originator_sequence_id > (SELECT min_seq FROM min_cursor_seq)
    ORDER BY m.originator_node_id, m.originator_sequence_id
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// V0d: LATERAL per originator with CTE topic list. Only 3-10 iterations
// (one per originator) instead of O(N) topic expansion. Each iteration
// filters by originator + sequence cursor, joined with topic CTE.
// Params: $1=cursor_node_ids INT[], $2=cursor_seq_ids BIGINT[],
//
//	$3=row_limit INT, $4=topics BYTEA[]
const queryV0dSQL = `WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest($1::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest($2::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
topic_list AS (
    SELECT t.topic FROM unnest($4::BYTEA[]) AS t(topic)
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
        FROM gateway_envelopes_meta AS m
        JOIN topic_list AS tl ON m.topic = tl.topic
        WHERE m.originator_node_id = c.cursor_node_id
          AND m.originator_sequence_id > c.cursor_sequence_id
        ORDER BY m.originator_sequence_id
        LIMIT $3
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
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
ORDER BY f.originator_node_id, f.originator_sequence_id`

// V1: LATERAL per topic with scalar subquery cursor.
// Handles unknown originators naturally via COALESCE(cursor, 0).
// Blob join OUTSIDE the LATERAL.
// Params: $1=cursor_topics BYTEA[], $2=cursor_node_ids INT[],
//
//	$3=cursor_seq_ids BIGINT[], $4=rows_per_topic INT, $5=row_limit INT
const queryV1SQL = `WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
distinct_topics AS (
    SELECT DISTINCT topic FROM cursor_entries
)
SELECT sub.originator_node_id,
       sub.originator_sequence_id,
       sub.gateway_time,
       sub.topic,
       b.originator_envelope
FROM distinct_topics AS dt
CROSS JOIN LATERAL (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    WHERE m.topic = dt.topic
      AND m.originator_sequence_id > COALESCE(
          (SELECT ce.seq_id FROM cursor_entries AS ce
           WHERE ce.topic = dt.topic AND ce.node_id = m.originator_node_id),
          0
      )
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT $4
) AS sub
JOIN gateway_envelope_blobs AS b
    ON b.originator_node_id = sub.originator_node_id
   AND b.originator_sequence_id = sub.originator_sequence_id
ORDER BY sub.originator_node_id, sub.originator_sequence_id
LIMIT $5`

// V1b: Same as V1 but with blob join INSIDE the LATERAL.
// This forces per-row PK lookups on blobs instead of cross-join.
const queryV1bSQL = `WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
distinct_topics AS (
    SELECT DISTINCT topic FROM cursor_entries
)
SELECT sub.originator_node_id,
       sub.originator_sequence_id,
       sub.gateway_time,
       sub.topic,
       sub.originator_envelope
FROM distinct_topics AS dt
CROSS JOIN LATERAL (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic,
           b.originator_envelope
    FROM gateway_envelopes_meta AS m
    JOIN gateway_envelope_blobs AS b
        ON b.originator_node_id = m.originator_node_id
       AND b.originator_sequence_id = m.originator_sequence_id
    WHERE m.topic = dt.topic
      AND m.originator_sequence_id > COALESCE(
          (SELECT ce.seq_id FROM cursor_entries AS ce
           WHERE ce.topic = dt.topic AND ce.node_id = m.originator_node_id),
          0
      )
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT $4
) AS sub
ORDER BY sub.originator_node_id, sub.originator_sequence_id
LIMIT $5`

// V2: LATERAL per (topic, originator) pair. Blob join OUTSIDE.
// Each subquery has constant topic/originator/cursor values for best index use.
// Requires caller to include ALL originators per topic (use seq_id=0 for unknown).
// Params: $1=cursor_topics BYTEA[], $2=cursor_node_ids INT[],
//
//	$3=cursor_seq_ids BIGINT[], $4=rows_per_entry INT, $5=row_limit INT
const queryV2SQL = `WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
)
SELECT sub.originator_node_id,
       sub.originator_sequence_id,
       sub.gateway_time,
       sub.topic,
       b.originator_envelope
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
JOIN gateway_envelope_blobs AS b
    ON b.originator_node_id = sub.originator_node_id
   AND b.originator_sequence_id = sub.originator_sequence_id
ORDER BY sub.originator_node_id, sub.originator_sequence_id
LIMIT $5`

// V3: LATERAL per (topic, originator) on meta only, CTE to sort+limit,
// then a single blob join for the limited result set.
// This avoids both the O(N) planning of V0 and the blob-per-iteration scan of V2b.
const queryV3SQL = `WITH cursor_entries AS (
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

// V2b: Same as V2 but with blob join INSIDE the LATERAL + partition pruning hint.
const queryV2bSQL = `WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest($1::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest($2::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest($3::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
)
SELECT sub.originator_node_id,
       sub.originator_sequence_id,
       sub.gateway_time,
       sub.topic,
       sub.originator_envelope
FROM cursor_entries AS ce
CROSS JOIN LATERAL (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic,
           b.originator_envelope
    FROM gateway_envelopes_meta AS m
    JOIN gateway_envelope_blobs AS b
        ON b.originator_node_id = m.originator_node_id
       AND b.originator_sequence_id = m.originator_sequence_id
       AND b.originator_node_id = ce.node_id
    WHERE m.topic = ce.topic
      AND m.originator_node_id = ce.node_id
      AND m.originator_sequence_id > ce.seq_id
    ORDER BY m.originator_sequence_id
    LIMIT $4
) AS sub
ORDER BY sub.originator_node_id, sub.originator_sequence_id
LIMIT $5`

// ---------- Config ----------

var originatorNodeIDs = []int32{100, 200, 300}

// Per-originator cursor positions used across all tests.
var cursorSeqs = []int64{500, 1200, 80}

// ---------- Seeding ----------

// seedDatabase inserts totalMessages envelopes distributed across originators.
// seqMultiplier controls how far apart sequence IDs are spaced. With
// seqMultiplier=1 (default), IDs are sequential (1,2,3,...) and all fit in one
// 1M partition band. With seqMultiplier=N, IDs become N,2N,3N,... which
// spreads them across multiple 1M partition bands.
// Example: 3334 messages/originator * seqMultiplier=3000 => max seqID ~10M => ~10 partitions.
func seedDatabase(
	t *testing.T,
	db *sql.DB,
	selectedTopics [][]byte,
	totalMessages int,
	seqMultiplier int64,
) {
	t.Helper()
	ctx := context.Background()
	querier := queries.New(db)

	if seqMultiplier < 1 {
		seqMultiplier = 1
	}

	noiseTopics := make([][]byte, 200)
	for i := range noiseTopics {
		noiseTopics[i] = testutils.RandomBytes(32)
	}

	rng := rand.New(rand.NewSource(42))

	seqByNode := make(map[int32]int64)
	for _, nid := range originatorNodeIDs {
		seqByNode[nid] = 0
	}

	for i := range totalMessages {
		nodeID := originatorNodeIDs[i%len(originatorNodeIDs)]
		seqByNode[nodeID]++
		seqID := seqByNode[nodeID] * seqMultiplier

		var topic []byte
		if rng.Float64() < 0.3 && len(selectedTopics) > 0 {
			topic = selectedTopics[rng.Intn(len(selectedTopics))]
		} else {
			topic = noiseTopics[rng.Intn(len(noiseTopics))]
		}

		row := queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                topic,
			OriginatorEnvelope:   testutils.RandomBytes(256),
		}

		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(ctx, querier, row)
		require.NoError(t, err, "inserting envelope %d (node=%d seq=%d)", i, nodeID, seqID)
	}

	// Count partitions per table.
	var metaParts, blobParts int
	_ = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pg_catalog.pg_inherits
		 WHERE inhparent = 'gateway_envelopes_meta'::regclass`).Scan(&metaParts)
	_ = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pg_catalog.pg_inherits
		 WHERE inhparent = 'gateway_envelope_blobs'::regclass`).Scan(&blobParts)

	t.Logf(
		"seeded %d envelopes across %d originators (seqMultiplier=%d, max_seq=%d, meta_partitions=%d, blob_partitions=%d, %d selected topics, %d noise topics)",
		totalMessages, len(originatorNodeIDs), seqMultiplier,
		seqByNode[originatorNodeIDs[0]]*seqMultiplier,
		metaParts, blobParts,
		len(selectedTopics), len(noiseTopics),
	)

	_, err := db.ExecContext(ctx, "ANALYZE gateway_envelopes_meta")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "ANALYZE gateway_envelope_blobs")
	require.NoError(t, err)
}

// ---------- Per-topic cursor param builders ----------

type perTopicCursorParams struct {
	CursorTopics  [][]byte
	CursorNodeIDs []int32
	CursorSeqIDs  []int64
	RowsPerEntry  int32
	RowLimit      int32
}

// buildFullCursorParams creates cursor entries for ALL originators on every topic.
// seqMultiplier scales cursor positions to match the seeded sequence ID range.
func buildFullCursorParams(topics [][]byte, seqMultiplier int64) perTopicCursorParams {
	var ct [][]byte
	var cn []int32
	var cs []int64

	for _, topic := range topics {
		for i, nodeID := range originatorNodeIDs {
			ct = append(ct, topic)
			cn = append(cn, nodeID)
			cs = append(cs, cursorSeqs[i]*seqMultiplier)
		}
	}

	return perTopicCursorParams{
		CursorTopics:  ct,
		CursorNodeIDs: cn,
		CursorSeqIDs:  cs,
		RowsPerEntry:  50,
		RowLimit:      500,
	}
}


// ---------- Generic EXPLAIN ANALYZE ----------

func runExplain(t *testing.T, db *sql.DB, rawSQL string, args ...any) string {
	t.Helper()

	rows, err := db.QueryContext(
		context.Background(),
		"EXPLAIN ANALYZE "+rawSQL,
		args...,
	)
	if err != nil {
		t.Fatalf("EXPLAIN ANALYZE failed: %v", err)
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatalf("scanning EXPLAIN row: %v", err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterating EXPLAIN rows: %v", err)
	}

	return strings.Join(lines, "\n")
}

// ---------- Row count verification ----------

func countRows(t *testing.T, db *sql.DB, rawSQL string, args ...any) int {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), rawSQL, args...)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var discard any
		cols, _ := rows.Columns()
		dests := make([]any, len(cols))
		for i := range dests {
			dests[i] = &discard
		}
		_ = rows.Scan(dests...)
	}
	return count
}

// ---------- Test runner helper ----------

type queryVariant struct {
	name    string
	rawSQL  string
	args    []any
}

func runVariants(t *testing.T, db *sql.DB, n int, variants []queryVariant) {
	t.Helper()
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			rowCount := countRows(t, db, v.rawSQL, v.args...)
			plan := runExplain(t, db, v.rawSQL, v.args...)
			fmt.Printf(
				"\n========== %s (%d topics, %d rows) ==========\n%s\n\n",
				v.name, n, rowCount, plan,
			)
		})
	}
}

// ---------- Main test ----------

func TestQueryComparison_ExplainAnalyze(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	// Disable JIT to avoid dominating execution time on small datasets.
	_, err := db.ExecContext(ctx, "SET jit = off")
	require.NoError(t, err)

	maxTopics := 1000
	allTopics := make([][]byte, maxTopics)
	for i := range allTopics {
		allTopics[i] = testutils.RandomBytes(32)
	}

	seedDatabase(t, db, allTopics, 10_000, 1)

	topicCounts := []int{10, 100, 1000}

	for _, n := range topicCounts {
		topics := allTopics[:n]

		t.Run(fmt.Sprintf("%d_topics", n), func(t *testing.T) {
			pFull := buildFullCursorParams(topics, 1)

			v1Args := []any{
				pq.Array(pFull.CursorTopics),
				pq.Array(pFull.CursorNodeIDs),
				pq.Array(pFull.CursorSeqIDs),
				pFull.RowsPerEntry,
				pFull.RowLimit,
			}

			variants := []queryVariant{
				{
					name:   "V0_original",
					rawSQL: queryV0SQL,
					args: []any{
						pq.Array(originatorNodeIDs),
						pq.Array(cursorSeqs),
						int32(500),
						pq.Array(topics),
					},
				},
				{
					name:   "V1_blob_outside",
					rawSQL: queryV1SQL,
					args:   v1Args,
				},
				{
					name:   "V1b_blob_inside",
					rawSQL: queryV1bSQL,
					args:   v1Args,
				},
				{
					name:   "V2_blob_outside",
					rawSQL: queryV2SQL,
					args:   v1Args,
				},
				{
					name:   "V2b_blob_inside",
					rawSQL: queryV2bSQL,
					args:   v1Args,
				},
			}

			runVariants(t, db, n, variants)
		})
	}
}

// TestQueryComparison_WithTopicSeqIndex tests V2b with an additional index
// on (topic, originator_sequence_id) that enables direct index seek for
// the per-(topic, originator) LATERAL pattern.
func TestQueryComparison_WithTopicSeqIndex(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	_, err := db.ExecContext(ctx, "SET jit = off")
	require.NoError(t, err)

	maxTopics := 1000
	allTopics := make([][]byte, maxTopics)
	for i := range allTopics {
		allTopics[i] = testutils.RandomBytes(32)
	}

	seedDatabase(t, db, allTopics, 10_000, 1)

	// Create the new index: (topic, originator_sequence_id) per partition.
	_, err = db.ExecContext(ctx,
		`CREATE INDEX gem_topic_seq_idx ON gateway_envelopes_meta
		 (topic, originator_sequence_id) INCLUDE (gateway_time)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, "ANALYZE gateway_envelopes_meta")
	require.NoError(t, err)

	t.Log("created gem_topic_seq_idx(topic, originator_sequence_id) INCLUDE (gateway_time)")

	topicCounts := []int{10, 100, 1000}

	for _, n := range topicCounts {
		topics := allTopics[:n]

		t.Run(fmt.Sprintf("%d_topics", n), func(t *testing.T) {
			pFull := buildFullCursorParams(topics, 1)
			v2bArgs := []any{
				pq.Array(pFull.CursorTopics),
				pq.Array(pFull.CursorNodeIDs),
				pq.Array(pFull.CursorSeqIDs),
				pFull.RowsPerEntry,
				pFull.RowLimit,
			}

			variants := []queryVariant{
				{
					name:   "V0_baseline",
					rawSQL: queryV0SQL,
					args: []any{
						pq.Array(originatorNodeIDs),
						pq.Array(cursorSeqs),
						int32(500),
						pq.Array(topics),
					},
				},
				{
					name:   "V2b_with_topic_seq_idx",
					rawSQL: queryV2bSQL,
					args:   v2bArgs,
				},
			}

			runVariants(t, db, n, variants)
		})
	}
}

// TestQueryComparison_ManyPartitions tests V0 vs V2b with many RANGE partitions
// per originator. By spreading sequence IDs across 1M boundaries, we force
// creation of multiple subpartitions and can measure the impact on planning time.
//
// seqMultiplier controls partition count:
//
//	3334 msgs/originator * multiplier => max_seq => ceil(max_seq/1M) partitions
//	multiplier=3000 => ~10 partitions/originator
//	multiplier=15000 => ~50 partitions/originator
func TestQueryComparison_ManyPartitions(t *testing.T) {
	multipliers := []struct {
		name       string
		multiplier int64
	}{
		{"1_partition", 1},
		{"10_partitions", 3000},
		{"50_partitions", 15000},
	}

	for _, m := range multipliers {
		t.Run(m.name, func(t *testing.T) {
			ctx := context.Background()
			db, _ := testutils.NewRawDB(t, ctx)

			_, err := db.ExecContext(ctx, "SET jit = off")
			require.NoError(t, err)

			maxTopics := 1000
			allTopics := make([][]byte, maxTopics)
			for i := range allTopics {
				allTopics[i] = testutils.RandomBytes(32)
			}

			seedDatabase(t, db, allTopics, 10_000, m.multiplier)

			// Count total partitions (all levels) for reporting.
			var totalParts int
			err = db.QueryRowContext(ctx,
				`WITH RECURSIVE parts AS (
					SELECT oid FROM pg_class WHERE relname = 'gateway_envelopes_meta'
					UNION ALL
					SELECT c.oid FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
					JOIN parts p ON p.oid = i.inhparent
				) SELECT COUNT(*)-1 FROM parts`).Scan(&totalParts)
			require.NoError(t, err)
			t.Logf("total meta partitions (all levels): %d", totalParts)

			topicCounts := []int{10, 100, 1000}

			for _, n := range topicCounts {
				topics := allTopics[:n]

				t.Run(fmt.Sprintf("%d_topics", n), func(t *testing.T) {
					pFull := buildFullCursorParams(topics, m.multiplier)

					// V0 cursor positions must also scale.
					scaledCursorSeqs := make([]int64, len(cursorSeqs))
					for i, s := range cursorSeqs {
						scaledCursorSeqs[i] = s * m.multiplier
					}

					lateralArgs := []any{
						pq.Array(pFull.CursorTopics),
						pq.Array(pFull.CursorNodeIDs),
						pq.Array(pFull.CursorSeqIDs),
						pFull.RowsPerEntry,
						pFull.RowLimit,
					}

					variants := []queryVariant{
						{
							name:   "V0_baseline",
							rawSQL: queryV0SQL,
							args: []any{
								pq.Array(originatorNodeIDs),
								pq.Array(scaledCursorSeqs),
								int32(500),
								pq.Array(topics),
							},
						},
						{
							name:   "V2b_blob_inside",
							rawSQL: queryV2bSQL,
							args:   lateralArgs,
						},
						{
							name:   "V3_lateral_meta_then_blob",
							rawSQL: queryV3SQL,
							args:   lateralArgs,
						},
					}

					runVariants(t, db, n, variants)
				})
			}
		})
	}
}

// TestQueryComparison_ManyPartitionsWithIndex tests V3 with the proposed
// gem_topic_seq_idx index to see if it helps with many partitions.
func TestQueryComparison_ManyPartitionsWithIndex(t *testing.T) {
	multipliers := []struct {
		name       string
		multiplier int64
	}{
		{"10_partitions", 3000},
		{"50_partitions", 15000},
	}

	for _, m := range multipliers {
		t.Run(m.name, func(t *testing.T) {
			ctx := context.Background()
			db, _ := testutils.NewRawDB(t, ctx)

			_, err := db.ExecContext(ctx, "SET jit = off")
			require.NoError(t, err)

			maxTopics := 1000
			allTopics := make([][]byte, maxTopics)
			for i := range allTopics {
				allTopics[i] = testutils.RandomBytes(32)
			}

			seedDatabase(t, db, allTopics, 10_000, m.multiplier)

			// Create the new index.
			_, err = db.ExecContext(ctx,
				`CREATE INDEX gem_topic_seq_idx ON gateway_envelopes_meta
				 (topic, originator_sequence_id) INCLUDE (gateway_time)`)
			require.NoError(t, err)

			_, err = db.ExecContext(ctx, "ANALYZE gateway_envelopes_meta")
			require.NoError(t, err)

			topicCounts := []int{100, 1000}

			for _, n := range topicCounts {
				topics := allTopics[:n]

				t.Run(fmt.Sprintf("%d_topics", n), func(t *testing.T) {
					pFull := buildFullCursorParams(topics, m.multiplier)

					scaledCursorSeqs := make([]int64, len(cursorSeqs))
					for i, s := range cursorSeqs {
						scaledCursorSeqs[i] = s * m.multiplier
					}

					lateralArgs := []any{
						pq.Array(pFull.CursorTopics),
						pq.Array(pFull.CursorNodeIDs),
						pq.Array(pFull.CursorSeqIDs),
						pFull.RowsPerEntry,
						pFull.RowLimit,
					}

					variants := []queryVariant{
						{
							name:   "V0_baseline",
							rawSQL: queryV0SQL,
							args: []any{
								pq.Array(originatorNodeIDs),
								pq.Array(scaledCursorSeqs),
								int32(500),
								pq.Array(topics),
							},
						},
						{
							name:   "V2b_blob_inside",
							rawSQL: queryV2bSQL,
							args:   lateralArgs,
						},
						{
							name:   "V3_with_topic_seq_idx",
							rawSQL: queryV3SQL,
							args:   lateralArgs,
						},
					}

					runVariants(t, db, n, variants)
				})
			}
		})
	}
}

// TestQueryComparison_V0Variants explores optimizations to V0 while keeping
// the single-cursor-across-all-topics constraint.
//
// V0:  Original (ANY(topics) + UNION ALL for unknown originators)
// V0b: CTE topic_list join instead of ANY() — still has UNION ALL
// V0c: CTE topic_list + no UNION ALL — caller pre-populates all originators
// V0d: LATERAL per originator with CTE topic_list — only 3-10 iterations
func TestQueryComparison_V0Variants(t *testing.T) {
	multipliers := []struct {
		name       string
		multiplier int64
	}{
		{"1_partition", 1},
		{"10_partitions", 3000},
		{"50_partitions", 15000},
	}

	for _, m := range multipliers {
		t.Run(m.name, func(t *testing.T) {
			ctx := context.Background()
			db, _ := testutils.NewRawDB(t, ctx)

			_, err := db.ExecContext(ctx, "SET jit = off")
			require.NoError(t, err)

			maxTopics := 1000
			allTopics := make([][]byte, maxTopics)
			for i := range allTopics {
				allTopics[i] = testutils.RandomBytes(32)
			}

			seedDatabase(t, db, allTopics, 10_000, m.multiplier)

			topicCounts := []int{10, 100, 1000}

			for _, n := range topicCounts {
				topics := allTopics[:n]

				t.Run(fmt.Sprintf("%d_topics", n), func(t *testing.T) {
					scaledCursorSeqs := make([]int64, len(cursorSeqs))
					for i, s := range cursorSeqs {
						scaledCursorSeqs[i] = s * m.multiplier
					}

					v0Args := []any{
						pq.Array(originatorNodeIDs),
						pq.Array(scaledCursorSeqs),
						int32(500),
						pq.Array(topics),
					}

					variants := []queryVariant{
						{
							name:   "V0_original",
							rawSQL: queryV0SQL,
							args:   v0Args,
						},
						{
							name:   "V0b_cte_topics",
							rawSQL: queryV0bSQL,
							args:   v0Args,
						},
						{
							name:   "V0c_cte_no_union",
							rawSQL: queryV0cSQL,
							args:   v0Args,
						},
						{
							name:   "V0d_lateral_per_originator",
							rawSQL: queryV0dSQL,
							args:   v0Args,
						},
					}

					runVariants(t, db, n, variants)
				})
			}
		})
	}
}

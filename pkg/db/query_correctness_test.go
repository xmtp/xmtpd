package db_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// ---------- Correctness Test Helpers ----------

type resultRow struct {
	nodeID int32
	seqID  int64
	topic  []byte
}

type testEnvelope struct {
	nodeID int32
	seqID  int64
	topic  []byte
}

func queryResultRows(t *testing.T, db *sql.DB, rawSQL string, args ...any) []resultRow {
	t.Helper()
	rows, err := db.QueryContext(context.Background(), rawSQL, args...)
	require.NoError(t, err)
	defer rows.Close()

	var results []resultRow
	for rows.Next() {
		var r resultRow
		var skip1, skip2 any
		err := rows.Scan(&r.nodeID, &r.seqID, &skip1, &r.topic, &skip2)
		require.NoError(t, err)
		results = append(results, r)
	}
	require.NoError(t, rows.Err())
	return results
}

func insertTestEnvelopes(t *testing.T, db *sql.DB, envelopes []testEnvelope) {
	t.Helper()
	ctx := context.Background()
	q := queries.New(db)
	for _, e := range envelopes {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(ctx, q, queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     e.nodeID,
			OriginatorSequenceID: e.seqID,
			Topic:                e.topic,
			OriginatorEnvelope:   []byte("payload"),
		})
		require.NoError(t, err)
	}
}

// v3bArgsFromV0 converts V0-style shared cursor to V3b per-topic cursor arrays.
// Every topic gets the same cursor for each originator.
func v3bArgsFromV0(
	topics [][]byte,
	nodeIDs []int32,
	seqIDs []int64,
	rowsPerEntry, rowLimit int32,
) []any {
	var ct [][]byte
	var cn []int32
	var cs []int64
	for _, topic := range topics {
		for i, nodeID := range nodeIDs {
			ct = append(ct, topic)
			cn = append(cn, nodeID)
			cs = append(cs, seqIDs[i])
		}
	}
	return []any{
		pq.Array(ct),
		pq.Array(cn),
		pq.Array(cs),
		rowsPerEntry,
		rowLimit,
	}
}

func requireSameRows(t *testing.T, label string, expected, actual []resultRow) {
	t.Helper()
	require.Len(t, actual, len(expected), "%s: row count mismatch", label)
	for i := range expected {
		require.Equalf(t, expected[i].nodeID, actual[i].nodeID, "%s: row %d nodeID", label, i)
		require.Equalf(t, expected[i].seqID, actual[i].seqID, "%s: row %d seqID", label, i)
		require.Equalf(t, expected[i].topic, actual[i].topic, "%s: row %d topic", label, i)
	}
}

// ---------- Correctness Tests ----------
//
// These tests verify functional equivalence and behavioral differences between
// V0 (shared cursor) and V3/V3b (per-topic cursor) query variants.

// TestQueryCorrectness_V3_V3b_Identical verifies that V3b is a pure performance
// optimization over V3 — both must return identical results.
func TestQueryCorrectness_V3_V3b_Identical(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01
	topicB := make([]byte, 32)
	topicB[0] = 0x02

	insertTestEnvelopes(t, db, []testEnvelope{
		{100, 1, topicA}, {100, 2, topicA}, {100, 3, topicB},
		{200, 1, topicA}, {200, 2, topicB},
		{300, 1, topicB},
	})

	args := v3bArgsFromV0(
		[][]byte{topicA, topicB},
		[]int32{100, 200, 300},
		[]int64{0, 0, 0},
		int32(500), int32(500),
	)

	v3Rows := queryResultRows(t, db, queryV3SQL, args...)
	v3bRows := queryResultRows(t, db, queryV3bSQL, args...)

	requireSameRows(t, "V3 vs V3b", v3Rows, v3bRows)
	require.Len(t, v3Rows, 6, "expected all 6 envelopes")
}

// TestQueryCorrectness_V0_V3b_EquivalentWithFullCursors verifies V0 and V3b
// return the same rows when given equivalent parameters: same cursor for all
// topics, all originators pre-populated, and rows_per_entry >= row_limit.
func TestQueryCorrectness_V0_V3b_EquivalentWithFullCursors(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01
	topicB := make([]byte, 32)
	topicB[0] = 0x02

	insertTestEnvelopes(t, db, []testEnvelope{
		{100, 1, topicA}, {100, 2, topicA}, {100, 3, topicA},
		{100, 4, topicB}, {100, 5, topicB},
		{200, 1, topicB}, {200, 2, topicA}, {200, 3, topicB},
		{300, 1, topicA}, {300, 2, topicA}, {300, 3, topicB},
	})

	topics := [][]byte{topicA, topicB}
	nodeIDs := []int32{100, 200, 300}
	seqIDs := []int64{1, 0, 1}

	v0Args := []any{
		pq.Array(nodeIDs),
		pq.Array(seqIDs),
		int32(500),
		pq.Array(topics),
	}
	v3bArgs := v3bArgsFromV0(topics, nodeIDs, seqIDs, int32(500), int32(500))

	v0Rows := queryResultRows(t, db, queryV0SQL, v0Args...)
	v3bRows := queryResultRows(t, db, queryV3bSQL, v3bArgs...)

	requireSameRows(t, "V0 vs V3b (aligned)", v0Rows, v3bRows)
	// node 100 seq>1: 2,3,4,5 = 4 rows
	// node 200 seq>0: 1,2,3 = 3 rows
	// node 300 seq>1: 2,3 = 2 rows
	require.Len(t, v0Rows, 9)
}

// TestQueryCorrectness_RowsPerEntryTruncation demonstrates that V3b's
// rows_per_entry parameter can cause fewer rows than V0 when a single
// (topic, originator) pair has many rows.
func TestQueryCorrectness_RowsPerEntryTruncation(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01

	var envs []testEnvelope
	for i := int64(1); i <= 20; i++ {
		envs = append(envs, testEnvelope{100, i, topicA})
	}
	insertTestEnvelopes(t, db, envs)

	topics := [][]byte{topicA}
	nodeIDs := []int32{100}
	seqIDs := []int64{0}

	// V0: returns all 20 rows.
	v0Rows := queryResultRows(t, db, queryV0SQL, []any{
		pq.Array(nodeIDs), pq.Array(seqIDs), int32(500), pq.Array(topics),
	}...)
	require.Len(t, v0Rows, 20)

	// V3b with rows_per_entry=5: capped to 5 rows per (topic, originator).
	v3bRows := queryResultRows(t, db, queryV3bSQL,
		v3bArgsFromV0(topics, nodeIDs, seqIDs, int32(5), int32(500))...)
	require.Len(t, v3bRows, 5, "V3b should return only 5 rows due to rows_per_entry")
	for i, r := range v3bRows {
		require.Equal(t, int64(i+1), r.seqID, "should be lowest sequence IDs")
	}

	// V3b with rows_per_entry=row_limit: matches V0.
	v3bRowsFull := queryResultRows(t, db, queryV3bSQL,
		v3bArgsFromV0(topics, nodeIDs, seqIDs, int32(500), int32(500))...)
	requireSameRows(t, "V3b(rows_per_entry=500) vs V0", v0Rows, v3bRowsFull)
}

// TestQueryCorrectness_EmptyCursors demonstrates the behavioral difference
// with empty cursors: V0 returns all matching rows (via UNION ALL branch B),
// while V3b returns nothing (no cursor entries to iterate).
//
// The fix is for the caller to always pre-populate all known originator IDs
// with seq_id=0 before calling V3b.
func TestQueryCorrectness_EmptyCursors(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01

	insertTestEnvelopes(t, db, []testEnvelope{
		{100, 1, topicA}, {100, 2, topicA},
		{200, 1, topicA},
	})

	topics := [][]byte{topicA}

	// V0 with empty cursor arrays: returns all rows via UNION ALL branch B.
	v0Rows := queryResultRows(t, db, queryV0SQL, []any{
		pq.Array([]int32{}), pq.Array([]int64{}), int32(500), pq.Array(topics),
	}...)
	require.Len(t, v0Rows, 3, "V0 with empty cursors returns all matching rows")

	// V3b with empty cursor entries: returns nothing.
	v3bRows := queryResultRows(t, db, queryV3bSQL, []any{
		pq.Array([][]byte{}), pq.Array([]int32{}), pq.Array([]int64{}),
		int32(500), int32(500),
	}...)
	require.Len(t, v3bRows, 0, "V3b with empty cursors returns no rows")

	// Fix: pre-populate all originators with seq_id=0.
	allNodeIDs := []int32{100, 200}
	allSeqIDs := []int64{0, 0}
	v3bRowsFixed := queryResultRows(t, db, queryV3bSQL,
		v3bArgsFromV0(topics, allNodeIDs, allSeqIDs, int32(500), int32(500))...)
	requireSameRows(t, "V3b(pre-populated) vs V0", v0Rows, v3bRowsFixed)
}

// TestQueryCorrectness_UnknownOriginators demonstrates that V0's UNION ALL
// branch B automatically catches rows from originators not listed in the
// cursor, while V3b only returns rows for originators explicitly included.
// The fix is the same: pre-populate all known originator IDs.
func TestQueryCorrectness_UnknownOriginators(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01

	insertTestEnvelopes(t, db, []testEnvelope{
		{100, 1, topicA}, {100, 2, topicA},
		{200, 1, topicA}, {200, 2, topicA}, {200, 3, topicA},
	})

	topics := [][]byte{topicA}
	// Cursor only includes originator 100. Originator 200 is "unknown".
	knownNodeIDs := []int32{100}
	knownSeqIDs := []int64{0}

	// V0: returns 100's rows + ALL of 200's rows (via UNION ALL branch B).
	v0Rows := queryResultRows(t, db, queryV0SQL, []any{
		pq.Array(knownNodeIDs), pq.Array(knownSeqIDs), int32(500), pq.Array(topics),
	}...)
	require.Len(t, v0Rows, 5, "V0 returns rows from both originators")

	// V3b: only returns originator 100's rows. 200 is missing from cursor entries.
	v3bRows := queryResultRows(t, db, queryV3bSQL,
		v3bArgsFromV0(topics, knownNodeIDs, knownSeqIDs, int32(500), int32(500))...)
	require.Len(t, v3bRows, 2, "V3b only returns rows for originators in cursor entries")

	// Fix: pre-populate ALL originators with seq_id=0.
	allNodeIDs := []int32{100, 200}
	allSeqIDs := []int64{0, 0}
	v3bRowsFull := queryResultRows(t, db, queryV3bSQL,
		v3bArgsFromV0(topics, allNodeIDs, allSeqIDs, int32(500), int32(500))...)
	requireSameRows(t, "V3b(all originators) vs V0", v0Rows, v3bRowsFull)
}

// TestQueryCorrectness_PerTopicCursorIndependence demonstrates V3b's key
// advantage: independent cursor positions per topic. V0's shared cursor
// must use the minimum position, re-fetching already-seen data for topics
// that are further ahead.
func TestQueryCorrectness_PerTopicCursorIndependence(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	topicA := make([]byte, 32)
	topicA[0] = 0x01
	topicB := make([]byte, 32)
	topicB[0] = 0x02

	// Originator 100: topicA at seq 1-10, topicB at seq 11-20.
	var envs []testEnvelope
	for i := int64(1); i <= 10; i++ {
		envs = append(envs, testEnvelope{100, i, topicA})
	}
	for i := int64(11); i <= 20; i++ {
		envs = append(envs, testEnvelope{100, i, topicB})
	}
	insertTestEnvelopes(t, db, envs)

	// Scenario: client has seen topicB up to seq 15 but is new to topicA.
	// Desired: all of topicA (seq 1-10) + topicB after 15 (seq 16-20) = 15 rows.

	// V3b: independent cursors per topic achieve exactly this.
	v3bRows := queryResultRows(t, db, queryV3bSQL, []any{
		pq.Array([][]byte{topicA, topicB}),
		pq.Array([]int32{100, 100}),
		pq.Array([]int64{0, 15}),
		int32(500), int32(500),
	}...)
	require.Len(t, v3bRows, 15, "V3b: 10 from topicA + 5 from topicB")

	topicACount := 0
	for _, r := range v3bRows {
		if r.topic[0] == 0x01 {
			topicACount++
			require.Greater(t, r.seqID, int64(0), "topicA rows after cursor 0")
		} else {
			require.Greater(t, r.seqID, int64(15), "topicB rows after cursor 15")
		}
	}
	require.Equal(t, 10, topicACount)

	// V0 must use cursor=0 (minimum) to not miss topicA rows.
	// This re-fetches topicB seq 11-15 that the client already has.
	v0Rows := queryResultRows(t, db, queryV0SQL, []any{
		pq.Array([]int32{100}),
		pq.Array([]int64{0}),
		int32(500),
		pq.Array([][]byte{topicA, topicB}),
	}...)
	require.Len(t, v0Rows, 20, "V0 returns all 20 rows including 5 already seen")
}

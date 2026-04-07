package migrations_test

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/migrations"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
)

// TestMigration23_RenamesExistingPartitionsSafely verifies that migration 00023
// (rename gateway_envelope_blobs -> gateway_envelopes_blob) can be applied
// safely to a database that already contains data and child partitions.
//
// The test stands up a database at schema version 22 (one step before the
// rename), populates it via the *v2* sqlc query path — which goes through the
// pre-rename insert_gateway_envelope_batch_v2 stored function and creates L1
// and L2 partitions under the old gateway_envelope_blobs_* naming — and then
// migrates forward to HEAD and asserts zero data loss.
//
// The pre-rename queries (InsertGatewayEnvelopeBatchV2, EnsureGatewayParts)
// remain in sqlc precisely so this test can populate at v22 without falling
// back to raw SQL. They wrap stored functions, so sqlc validates them on
// signature alone — the function bodies are opaque and may reference the
// soon-to-be-renamed table.
func TestMigration23_RenamesExistingPartitionsSafely(t *testing.T) {
	ctx := t.Context()
	database, _ := testutils.NewRawDBAtVersion(t, ctx, 22)

	// Sanity: at version 22 the table still has the old name.
	tableExists(t, database, "gateway_envelope_blobs")
	assertTableAbsent(t, database, "gateway_envelopes_blob")

	// Populate using the v2 sqlc query path against the pre-rename schema.
	populateDatabaseAtV22(t, database)

	// Snapshot data + partition counts under the old name.
	blobsBefore := snapshotEnvelopeKeys(t, database, "gateway_envelope_blobs")
	viewBefore := snapshotEnvelopeKeys(t, database, "gateway_envelopes_view")
	require.NotEmpty(t, blobsBefore, "expected blob rows after populate")
	require.NotEmpty(t, viewBefore, "expected view rows after populate")

	oldNameChildCount := countTablesLike(t, database, "gateway_envelope_blobs\\_%")
	require.Positive(t, oldNameChildCount, "expected old-name child partitions at v22")

	// Apply migration 23.
	require.NoError(
		t,
		migrations.Migrate(ctx, database),
		"migration 22 -> 23 failed against populated DB",
	)

	// Parent table is under the new name; old name is gone.
	tableExists(t, database, "gateway_envelopes_blob")
	assertTableAbsent(t, database, "gateway_envelope_blobs")

	// Every child partition was renamed 1:1.
	newNameChildCount := countTablesLike(t, database, "gateway_envelopes_blob\\_%")
	assert.Equal(
		t,
		oldNameChildCount,
		newNameChildCount,
		"migration should rename every child partition 1:1",
	)
	assert.Zero(
		t,
		countTablesLike(t, database, "gateway_envelope_blobs\\_%"),
		"no old-name child partitions should remain after migration",
	)

	// Verify each expected L1 + L2 partition exists explicitly.
	for _, oid := range originatorIDs {
		l1 := "gateway_envelopes_blob_o" + strconv.Itoa(int(oid))
		tableExists(t, database, l1)
		for band := range 5 {
			start := int64(band) * 1_000_000
			end := start + 1_000_000
			l2 := l1 + "_s" + strconv.FormatInt(start, 10) + "_" + strconv.FormatInt(end, 10)
			tableExists(t, database, l2)
		}
	}

	// Data is fully intact via both the renamed blob table and the recreated view.
	blobsAfter := snapshotEnvelopeKeys(t, database, "gateway_envelopes_blob")
	viewAfter := snapshotEnvelopeKeys(t, database, "gateway_envelopes_view")
	assert.Equal(t, blobsBefore, blobsAfter, "blob rows lost across migration")
	assert.Equal(t, viewBefore, viewAfter, "view rows lost across migration")

	// View exists and the new v3 partition + batch-insert functions are in place.
	viewExists(t, database, "gateway_envelopes_view")
	functionExists(t, database, "make_blob_originator_part_v3")
	functionExists(t, database, "make_blob_seq_subpart_v3")
	functionExists(t, database, "ensure_gateway_parts_v3")
	functionExists(t, database, "insert_gateway_envelope_batch_v3")

	// Migration is append-only: pre-existing v1/v2 functions are still in
	// pg_proc (their bodies reference the old table name and are no longer
	// callable, but the catalog entries remain so this assertion holds).
	functionExists(t, database, "make_blob_originator_part")
	functionExists(t, database, "make_blob_seq_subpart")
	functionExists(t, database, "make_blob_originator_part_v2")
	functionExists(t, database, "make_blob_seq_subpart_v2")
	functionExists(t, database, "insert_gateway_envelope_batch")
	functionExists(t, database, "insert_gateway_envelope_batch_v2")
	functionExists(t, database, "ensure_gateway_parts_v2")
}

type envelopeKey struct {
	originatorNodeID     int32
	originatorSequenceID int64
}

func snapshotEnvelopeKeys(
	t *testing.T,
	database *sql.DB,
	table string,
) map[envelopeKey]struct{} {
	t.Helper()
	rows, err := database.QueryContext(
		t.Context(),
		"SELECT originator_node_id, originator_sequence_id FROM "+table,
	)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	out := map[envelopeKey]struct{}{}
	for rows.Next() {
		var k envelopeKey
		require.NoError(t, rows.Scan(&k.originatorNodeID, &k.originatorSequenceID))
		out[k] = struct{}{}
	}
	require.NoError(t, rows.Err())
	return out
}

func countTablesLike(t *testing.T, database *sql.DB, pattern string) int {
	t.Helper()
	var n int
	err := database.QueryRowContext(
		t.Context(),
		`SELECT COUNT(*) FROM pg_class c
		 JOIN pg_namespace n ON n.oid = c.relnamespace
		 WHERE n.nspname = 'public'
		   AND c.relkind IN ('r', 'p')
		   AND c.relname LIKE $1`,
		pattern,
	).Scan(&n)
	require.NoError(t, err)
	return n
}

// populateDatabaseAtV22 populates a database at schema version 22 (pre-rename)
// using the surviving v2 sqlc query path. It calls EnsureGatewayParts (v2
// stored function) to create partitions under the old gateway_envelope_blobs_*
// naming, then calls InsertGatewayEnvelopeBatchV2 to insert envelopes via the
// pre-rename batch insert function. Both stored functions reference the
// old-named blob table, which is exactly the schema present at v22.
//
// Mirrors populateDatabase: 3 originators × 5 bands × 3 topics = 45 envelopes,
// 3 L1 partitions, 15 L2 subpartitions.
func populateDatabaseAtV22(t *testing.T, database *sql.DB) {
	ctx := t.Context()
	q := queries.New(database)

	for _, nodeID := range originatorIDs {
		privKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		_, err = q.InsertNodeInfo(ctx, queries.InsertNodeInfoParams{
			NodeID:    nodeID,
			PublicKey: crypto.FromECDSAPub(&privKey.PublicKey),
		})
		require.NoError(t, err)
	}

	payerIDs := make([]int32, 3)
	for i := range payerIDs {
		payerIDs[i] = testutils.CreatePayer(t, database, testutils.RandomAddress().Hex())
		require.NotZero(t, payerIDs[i])
	}

	topics := [][]byte{topicA, topicB, topicC}

	// Pre-create every (originator, band) partition pair via the v2 ensure
	// function — at v22 this targets gateway_envelope_blobs_*.
	for _, originatorID := range originatorIDs {
		for band := range 5 {
			baseSeqID := int64(band) * db.GatewayEnvelopeBandWidth
			require.NoError(t, q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
				OriginatorNodeID:     originatorID,
				OriginatorSequenceID: baseSeqID,
				BandWidth:            db.GatewayEnvelopeBandWidth,
			}))
		}
	}

	// Build a single batch and insert via the v2 stored function.
	var (
		nodeIDs   []int32
		seqIDs    []int64
		topicArr  [][]byte
		payerArr  []int32
		times     []time.Time
		expiries  []int64
		envelopes [][]byte
		spends    []int64
		usage     []bool
		congest   []bool
	)
	now := time.Now().UTC()

	for _, originatorID := range originatorIDs {
		for band := range 5 {
			baseSeqID := int64(band) * db.GatewayEnvelopeBandWidth
			for k, topic := range topics {
				seqID := baseSeqID + int64(k)
				envelope := testutils.Marshal(
					t,
					envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
						t,
						uint32(originatorID),
						uint64(seqID),
						topic,
					),
				)
				nodeIDs = append(nodeIDs, originatorID)
				seqIDs = append(seqIDs, seqID)
				topicArr = append(topicArr, topic)
				payerArr = append(payerArr, payerIDs[k])
				times = append(times, now)
				expiries = append(expiries, 0)
				envelopes = append(envelopes, envelope)
				spends = append(spends, 0)
				usage = append(usage, false)
				congest = append(congest, false)
			}
		}
	}

	result, err := q.InsertGatewayEnvelopeBatchV2(ctx, queries.InsertGatewayEnvelopeBatchV2Params{
		OriginatorNodeIds:     nodeIDs,
		OriginatorSequenceIds: seqIDs,
		Topics:                topicArr,
		PayerIds:              payerArr,
		GatewayTimes:          times,
		Expiries:              expiries,
		OriginatorEnvelopes:   envelopes,
		SpendPicodollars:      spends,
		CountUsage:            usage,
		CountCongestion:       congest,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(nodeIDs)), result.InsertedMetaRows)
	require.Equal(t, int64(len(nodeIDs)), result.InsertedBlobRows)
}

func assertTableAbsent(t *testing.T, database *sql.DB, tableName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)`,
		tableName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists, "table %s should not exist", tableName)
}

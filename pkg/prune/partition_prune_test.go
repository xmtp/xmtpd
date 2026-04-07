package prune_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func tableExists(t *testing.T, ctx context.Context, db *sql.DB, tableName string) bool {
	t.Helper()

	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM pg_class c
			JOIN pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relname = $1
		)
	`, tableName).Scan(&exists)
	require.NoError(t, err)

	return exists
}

func ensureBandExists(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorID int32,
	seqID int64,
) {
	t.Helper()

	_, err := db.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		originatorID,
		seqID,
		int64(1000000),
	)
	require.NoError(t, err)
}

func insertEnvelope(
	t *testing.T,
	db *sql.DB,
	originatorID int32,
	seqID int64,
) {
	t.Helper()

	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV3Params{{
		OriginatorNodeID:     originatorID,
		OriginatorSequenceID: seqID,
		Topic:                []byte("topic"),
		OriginatorEnvelope:   []byte("payload"),
		GatewayTime:          time.Now(),
		Expiry:               time.Now().Add(1 * time.Hour).Unix(),
	}})
}

func TestExecutor_DropPrunablePartitions_DropsEmptyLowerMetaAndBlob(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const oid int32 = 100

	// Create two partitions:
	//   [0, 1000000)       -> empty, should be prunable
	//   [1000000, 2000000) -> ceiling, should be preserved
	ensureBandExists(t, ctx, db, oid, 1)
	ensureBandExists(t, ctx, db, oid, 1000001)

	upperMeta := "gateway_envelopes_meta_o100_s1000000_2000000"
	lowerMeta := "gateway_envelopes_meta_o100_s0_1000000"
	upperBlob := "gateway_envelopes_blob_o100_s1000000_2000000"
	lowerBlob := "gateway_envelopes_blob_o100_s0_1000000"

	require.True(t, tableExists(t, ctx, db, lowerMeta))
	require.True(t, tableExists(t, ctx, db, upperMeta))
	require.True(t, tableExists(t, ctx, db, lowerBlob))
	require.True(t, tableExists(t, ctx, db, upperBlob))

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})

	err := exec.DropPrunablePartitions()
	require.NoError(t, err)

	assert.False(t, tableExists(t, ctx, db, lowerMeta))
	assert.False(t, tableExists(t, ctx, db, lowerBlob))

	assert.True(t, tableExists(t, ctx, db, upperMeta))
	assert.True(t, tableExists(t, ctx, db, upperBlob))
}

func TestExecutor_DropPrunablePartitions_DoesNotDropCeilingEvenIfEmpty(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const oid int32 = 100

	// Both partitions empty; only the lower one should be dropped.
	ensureBandExists(t, ctx, db, oid, 1)
	ensureBandExists(t, ctx, db, oid, 1000001)

	lowerMeta := "gateway_envelopes_meta_o100_s0_1000000"
	upperMeta := "gateway_envelopes_meta_o100_s1000000_2000000"
	lowerBlob := "gateway_envelopes_blob_o100_s0_1000000"
	upperBlob := "gateway_envelopes_blob_o100_s1000000_2000000"

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})

	err := exec.DropPrunablePartitions()
	require.NoError(t, err)

	assert.False(t, tableExists(t, ctx, db, lowerMeta))
	assert.False(t, tableExists(t, ctx, db, lowerBlob))

	assert.True(t, tableExists(t, ctx, db, upperMeta))
	assert.True(t, tableExists(t, ctx, db, upperBlob))
}

func TestExecutor_DropPrunablePartitions_DryRun(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const oid int32 = 100

	ensureBandExists(t, ctx, db, oid, 1)
	ensureBandExists(t, ctx, db, oid, 1000001)

	lowerMeta := "gateway_envelopes_meta_o100_s0_1000000"
	upperMeta := "gateway_envelopes_meta_o100_s1000000_2000000"
	lowerBlob := "gateway_envelopes_blob_o100_s0_1000000"
	upperBlob := "gateway_envelopes_blob_o100_s1000000_2000000"

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    true,
		MaxCycles: 1,
	})

	err := exec.DropPrunablePartitions()
	require.NoError(t, err)

	assert.True(t, tableExists(t, ctx, db, lowerMeta))
	assert.True(t, tableExists(t, ctx, db, upperMeta))
	assert.True(t, tableExists(t, ctx, db, lowerBlob))
	assert.True(t, tableExists(t, ctx, db, upperBlob))
}

func TestExecutor_DropPrunablePartitions_DoesNotDropNonEmptyLowerPartition(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const oid int32 = 100

	ensureBandExists(t, ctx, db, oid, 1)
	ensureBandExists(t, ctx, db, oid, 1000001)

	lowerMeta := "gateway_envelopes_meta_o100_s0_1000000"
	upperMeta := "gateway_envelopes_meta_o100_s1000000_2000000"
	lowerBlob := "gateway_envelopes_blob_o100_s0_1000000"
	upperBlob := "gateway_envelopes_blob_o100_s1000000_2000000"

	// Make lower partition non-empty.
	insertEnvelope(t, db, oid, 1)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})

	err := exec.DropPrunablePartitions()
	require.NoError(t, err)

	assert.True(t, tableExists(t, ctx, db, lowerMeta))
	assert.True(t, tableExists(t, ctx, db, lowerBlob))

	assert.True(t, tableExists(t, ctx, db, upperMeta))
	assert.True(t, tableExists(t, ctx, db, upperBlob))
}

func TestExecutor_DropPrunablePartitions_PerOriginatorCeiling(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const oid1 int32 = 100
	const oid2 int32 = 101

	for _, oid := range []int32{oid1, oid2} {
		ensureBandExists(t, ctx, db, oid, 1)
		ensureBandExists(t, ctx, db, oid, 1000001)
	}

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})

	err := exec.DropPrunablePartitions()
	require.NoError(t, err)

	assert.False(t, tableExists(t, ctx, db, "gateway_envelopes_meta_o100_s0_1000000"))
	assert.False(t, tableExists(t, ctx, db, "gateway_envelopes_blob_o100_s0_1000000"))
	assert.True(t, tableExists(t, ctx, db, "gateway_envelopes_meta_o100_s1000000_2000000"))
	assert.True(t, tableExists(t, ctx, db, "gateway_envelopes_blob_o100_s1000000_2000000"))

	assert.False(t, tableExists(t, ctx, db, "gateway_envelopes_meta_o101_s0_1000000"))
	assert.False(t, tableExists(t, ctx, db, "gateway_envelopes_blob_o101_s0_1000000"))
	assert.True(t, tableExists(t, ctx, db, "gateway_envelopes_meta_o101_s1000000_2000000"))
	assert.True(t, tableExists(t, ctx, db, "gateway_envelopes_blob_o101_s1000000_2000000"))
}

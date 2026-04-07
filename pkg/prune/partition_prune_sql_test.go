package prune_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func mustExec(t *testing.T, ctx context.Context, db *sql.DB, sql string, args ...any) {
	t.Helper()
	_, err := db.ExecContext(ctx, sql, args...)
	require.NoError(t, err)
}

func TestGetPrunableMetaPartitions_ReturnsOnlyTrulyEmptyNonCeiling(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	const oid int32 = 100

	// Create two bands for the same originator.
	// Lower band should become prunable if empty.
	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1),
		int64(1000000),
	)
	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1000001),
		int64(1000000),
	)

	// Put one row in the upper band so it is definitely not empty.
	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV3Params{{
		OriginatorNodeID:     oid,
		OriginatorSequenceID: 1000001,
		Topic:                []byte("topic"),
		OriginatorEnvelope:   []byte("payload"),
		GatewayTime:          time.Now(),
		Expiry:               time.Now().Add(1 * time.Hour).Unix(),
	}})

	parts, err := q.GetPrunableMetaPartitions(ctx)
	require.NoError(t, err)

	require.Len(t, parts, 1)
	assert.Equal(t, oid, parts[0].OriginatorNodeID)
	assert.Equal(t, "gateway_envelopes_meta_o100_s0_1000000", parts[0].Tablename)
	assert.EqualValues(t, 0, parts[0].BandStart)
	assert.EqualValues(t, 1000000, parts[0].BandEnd)
}

func TestGetPrunableMetaPartitions_DoesNotReturnCeilingEvenIfEmpty(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	const oid int32 = 100

	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1),
		int64(1000000),
	)
	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1000001),
		int64(1000000),
	)

	// Both partitions remain empty.

	parts, err := q.GetPrunableMetaPartitions(ctx)
	require.NoError(t, err)

	require.Len(t, parts, 1)
	assert.Equal(t, "gateway_envelopes_meta_o100_s0_1000000", parts[0].Tablename)

	for _, p := range parts {
		assert.NotEqual(t, "gateway_envelopes_meta_o100_s1000000_2000000", p.Tablename)
	}
}

func TestGetPrunableMetaPartitions_DoesNotReturnNonEmptyLowerPartition(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	const oid int32 = 100

	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1),
		int64(1000000),
	)
	mustExec(
		t,
		ctx,
		db,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid,
		int64(1000001),
		int64(1000000),
	)

	// Insert into the lower band, making it non-empty.
	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV3Params{{
		OriginatorNodeID:     oid,
		OriginatorSequenceID: 1,
		Topic:                []byte("topic"),
		OriginatorEnvelope:   []byte("payload"),
		GatewayTime:          time.Now(),
		Expiry:               time.Now().Add(1 * time.Hour).Unix(),
	}})

	parts, err := q.GetPrunableMetaPartitions(ctx)
	require.NoError(t, err)

	assert.Empty(t, parts)
}

func TestGetPrunableMetaPartitions_PerOriginatorCeiling(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	const oid1 int32 = 100
	const oid2 int32 = 101

	for _, oid := range []int32{oid1, oid2} {
		mustExec(
			t,
			ctx,
			db,
			`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
			oid,
			int64(1),
			int64(1000000),
		)
		mustExec(
			t,
			ctx,
			db,
			`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
			oid,
			int64(1000001),
			int64(1000000),
		)
	}

	// Put a row only in the top partition for each originator.
	for _, oid := range []int32{oid1, oid2} {
		testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV3Params{{
			OriginatorNodeID:     oid,
			OriginatorSequenceID: 1000001,
			Topic:                []byte("topic"),
			OriginatorEnvelope:   []byte("payload"),
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(1 * time.Hour).Unix(),
		}})
	}

	parts, err := q.GetPrunableMetaPartitions(ctx)
	require.NoError(t, err)

	require.Len(t, parts, 2)

	got := map[int32]string{}
	for _, p := range parts {
		got[p.OriginatorNodeID] = p.Tablename
	}

	assert.Equal(t, "gateway_envelopes_meta_o100_s0_1000000", got[oid1])
	assert.Equal(t, "gateway_envelopes_meta_o101_s0_1000000", got[oid2])
}

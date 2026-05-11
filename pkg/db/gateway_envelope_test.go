package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildParams(
	payerID int32,
	originatorID int32,
	sequenceID int64,
	spendPicodollars int64,
) (queries.InsertGatewayEnvelopeV3Params, queries.IncrementUnsettledUsageParams) {
	insertParams := queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     originatorID,
		OriginatorSequenceID: sequenceID,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(100),
		PayerID:              xmtpd_db.NullInt32(payerID),
	}

	incrementParams := queries.IncrementUnsettledUsageParams{
		PayerID:           payerID,
		OriginatorID:      originatorID,
		MinutesSinceEpoch: 1,
		SpendPicodollars:  spendPicodollars,
	}

	return insertParams, incrementParams
}

func TestInsertAndIncrement(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)
	// Create a payer
	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := int32(100)
	sequenceID := int64(10)

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	numInserted, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
		true,
	)
	require.NoError(t, err)
	require.Equal(t, int64(1), numInserted)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(100), payerSpend.TotalSpendPicodollars)
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)

	originatorCongestion, err := querier.SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{OriginatorID: originatorID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(1), originatorCongestion)
}

func TestPayerMustExist(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	payerID := testutils.RandomInt32()
	originatorID := int32(100)
	sequenceID := int64(10)

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	_, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
		true,
	)
	require.Error(t, err)
}

func TestInsertAndIncrementParallel(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)
	// Create a payer
	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := int32(100)
	sequenceID := int64(10)
	numberOfInserts := 20

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	var wg sync.WaitGroup

	totalInserted := int64(0)

	attemptInsert := func() {
		defer wg.Done()
		numInserted, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
			ctx,
			db,
			insertParams,
			incrementParams,
			true,
		)
		require.NoError(t, err)
		atomic.AddInt64(&totalInserted, numInserted)
	}

	// insert one to avoid DDL creation deadlocks
	wg.Add(1)
	attemptInsert()

	for range numberOfInserts {
		wg.Add(1)
		go attemptInsert()
	}

	wg.Wait()

	require.Equal(t, int64(1), totalInserted)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(100), payerSpend.TotalSpendPicodollars)
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)

	originatorCongestion, err := querier.SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{OriginatorID: originatorID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(1), originatorCongestion)
}

func TestInsertAndIncrementWithOutOfOrderSequenceID(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)

	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := int32(100)
	sequenceID := int64(10)

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	_, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
		true,
	)
	require.NoError(t, err)

	lowerSequenceID := int64(5)

	insertParams, incrementParams = buildParams(payerID, originatorID, lowerSequenceID, 100)

	_, err = xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
		true,
	)
	require.NoError(t, err)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)
}

func TestInsertGatewayEnvelopeWithChecksStandalone(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     int32(i),
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(100),
			},
		)
		require.NoError(t, err)

	}
}

func TestInsertGatewayEnvelopeWithChecksStandalone_FailsInTransaction(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	querier := queries.New(db).WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     int32(i),
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(100),
			},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "current transaction is aborted")
	}
}

func TestInsertGatewayEnvelopeWithChecksStandalone_AutoCreateAndRetry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	q := queries.New(db)

	const nodeID int32 = 42
	const seqID int64 = 1

	_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(128),
		},
	)
	require.NoError(t, err)

	// A second insert into the SAME band should succeed without needing to create parts again.
	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID + 123, // still within [0..1_000_000) default band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(128),
		},
	)
	require.NoError(t, err)
}

// TestEnsureGatewayPartsV4_ConcurrentCreate verifies that many concurrent
// callers racing to create partitions for the same (originator, band) are
// all serialized via the advisory locks added in migration 00024, and that
// the end state is exactly one L1 + one L2 pair on both the meta and blob
// sides.
//
// Before migration 00024, ensure_gateway_parts_v3 did not serialize callers
// and relied on a regex match of SQLERRM to swallow "already a partition"
// errors — a race could leave the child unattached while the caller saw
// success, which later surfaced as SQLSTATE 23514 ("no partition of
// relation...") on the insert path. The v4 helpers hold a per-(originator,
// band) `pg_advisory_xact_lock` around the CREATE/ATTACH window and short-
// circuit via `pg_inherits` when the partition is already attached, so
// concurrent callers converge on the same committed state with no errors.
//
// Note: this test races the PARTITION CREATION path specifically (via
// EnsureGatewayPartsV4). It does not overlap the creation with concurrent
// inserts, because PostgreSQL's intrinsic lock ordering between INSERT
// (RowExclusive on parent) and ATTACH PARTITION (ShareRowExclusive on
// parent) can deadlock independently of this code — see the existing
// `TestInsertAndIncrementParallel` which works around the same limitation
// with an "insert one to avoid DDL creation deadlocks" warmup.
func TestEnsureGatewayPartsV4_ConcurrentCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	q := queries.New(db)

	const (
		nodeID       int32 = 4242
		seqID        int64 = 1
		bandWidth          = xmtpd_db.GatewayEnvelopeBandWidth
		numGoroutine       = 32
	)

	var (
		wg      sync.WaitGroup
		errCh   = make(chan error, numGoroutine)
		startCh = make(chan struct{})
	)

	for range numGoroutine {
		wg.Go(func() {
			<-startCh
			errCh <- q.EnsureGatewayPartsV4(ctx, queries.EnsureGatewayPartsV4Params{
				OriginatorNodeID:     nodeID,
				OriginatorSequenceID: seqID,
				BandWidth:            bandWidth,
			})
		})
	}

	// Release all goroutines at once to maximise the chance of racing on
	// partition creation.
	close(startCh)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "concurrent ensure_gateway_parts_v4 must not fail")
	}

	// End state: exactly one L1 (and one L2) child under each of the meta
	// and blob parents for this originator.
	for _, parent := range []string{"gateway_envelopes_meta", "gateway_envelopes_blob"} {
		l1Name := fmt.Sprintf("%s_o%d", parent, nodeID)

		var l1Count int
		require.NoError(
			t,
			db.QueryRowContext(
				ctx,
				`SELECT COUNT(*)
				 FROM pg_inherits i
				 JOIN pg_class c ON c.oid = i.inhrelid
				 JOIN pg_class p ON p.oid = i.inhparent
				 WHERE p.relname = $1 AND c.relname = $2`,
				parent, l1Name,
			).Scan(&l1Count),
		)
		require.Equal(t, 1, l1Count, "expected exactly one L1 child under %s", parent)

		var l2Count int
		require.NoError(
			t,
			db.QueryRowContext(
				ctx,
				`SELECT COUNT(*)
				 FROM pg_inherits i
				 JOIN pg_class c ON c.oid = i.inhrelid
				 JOIN pg_class p ON p.oid = i.inhparent
				 WHERE p.relname = $1`,
				l1Name,
			).Scan(&l2Count),
		)
		require.Equal(t, 1, l2Count, "expected exactly one L2 child under %s", l1Name)
	}
}

// TestInsertGatewayEnvelopeWithChecksStandalone_ConcurrentWithWarmup verifies
// that, once partitions exist, concurrent standalone inserts all land
// successfully without deadlocks or "no partition of relation" errors.
//
// This mirrors the pattern in `TestInsertAndIncrementParallel`: a single
// initial insert warms up the partitions so subsequent concurrent inserts
// hit the fast path. It exercises the EnsureGatewayPartsV4 short-circuit
// (pg_inherits says "already attached") when many goroutines retry after
// racing on a subsequent band boundary, without triggering the PG-intrinsic
// INSERT-vs-ATTACH deadlock.
func TestInsertGatewayEnvelopeWithChecksStandalone_ConcurrentWithWarmup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	q := queries.New(db)

	const (
		nodeID       int32 = 4343
		numGoroutine int32 = 16
	)

	// Warmup insert to create partitions and avoid the PG-intrinsic
	// INSERT-vs-ATTACH deadlock in the concurrent phase.
	_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	var (
		wg      sync.WaitGroup
		errCh   = make(chan error, numGoroutine)
		startCh = make(chan struct{})
	)

	for i := range numGoroutine {
		wg.Add(1)
		go func(seq int64) {
			defer wg.Done()
			<-startCh
			_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
				ctx,
				q,
				queries.InsertGatewayEnvelopeV3Params{
					OriginatorNodeID:     nodeID,
					OriginatorSequenceID: seq,
					Topic:                testutils.RandomBytes(32),
					OriginatorEnvelope:   testutils.RandomBytes(64),
				},
			)
			errCh <- err
		}(int64(i + 2))
	}

	close(startCh)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "concurrent standalone insert must not fail after warmup")
	}

	var count int64
	require.NoError(
		t,
		db.QueryRowContext(
			ctx,
			`SELECT COUNT(*) FROM gateway_envelopes_meta WHERE originator_node_id = $1`,
			nodeID,
		).Scan(&count),
	)
	require.Equal(t, int64(numGoroutine+1), count, "expected warmup + concurrent rows")
}

func TestInsertGatewayEnvelopeWithChecksStandalone_PreexistingPartitions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	q := queries.New(db)

	const nodeID int32 = 7
	const seqID int64 = 10

	// Explicitly ensure parts up-front to simulate "partitions already exist" path.
	err := q.EnsureGatewayPartsV3(ctx, queries.EnsureGatewayPartsV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		BandWidth:            1_000_000,
	})
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksStandalone_BandBoundaries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	q := queries.New(db)

	const nodeID int32 = 99

	// Sequence values straddling a band boundary:
	seqLeft := xmtpd_db.GatewayEnvelopeBandWidth - 1 // falls into band [0, bw)
	seqRight := xmtpd_db.GatewayEnvelopeBandWidth    // falls into band [bw, 2*bw)

	_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqRight,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft + 123, // still within first band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqRight + 456, // still within second band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTransactional_FailsWithoutTx(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     int32(i),
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(100),
			},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "SAVEPOINT can only be used in transaction blocks")
	}
}

func TestInsertGatewayEnvelopeWithCheckTransactional(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	querier := queries.New(db).WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     int32(i),
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(100),
			},
		)
		require.NoError(t, err)
	}
}

func TestInsertGatewayEnvelopeWithChecksTx_AutoCreateAndRetry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	const nodeID int32 = 42
	const seqID int64 = 1

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(128),
		},
	)
	require.NoError(t, err)

	// A second insert into the SAME band should succeed without needing to create parts again.
	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID + 123, // still within [0..1_000_000) default band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(128),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTx_PreexistingPartitions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	const nodeID int32 = 7
	const seqID int64 = 10

	// Explicitly ensure parts up-front to simulate "partitions already exist" path.
	err = q.EnsureGatewayPartsV3(ctx, queries.EnsureGatewayPartsV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		BandWidth:            xmtpd_db.GatewayEnvelopeBandWidth,
	})
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTxn_BandBoundaries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	const nodeID int32 = 99

	// Sequence values straddling a band boundary:
	seqLeft := xmtpd_db.GatewayEnvelopeBandWidth - 1 // falls into band [0, bw)
	seqRight := xmtpd_db.GatewayEnvelopeBandWidth    // falls into band [bw, 2*bw)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqRight,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft + 123, // still within first band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqRight + 456, // still within second band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTx_RollbackDDL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)

	const nodeID int32 = 7
	const seqID int64 = 10

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_ = tx.Rollback()

	tx2, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q2 := queries.New(db).WithTx(tx2)
	defer func(tx2 *sql.Tx) {
		_ = tx2.Rollback()
	}(tx2)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q2,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTx_Commit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)

	const nodeID int32 = 7
	const seqID int64 = 10

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqID,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_ = tx.Commit()
	t.Logf("tx committed")
}

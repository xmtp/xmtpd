package db_test

import (
	"context"
	"database/sql"
	"errors"
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

// db_runInsertWithEnsure mirrors the production caller contract: run the insert in a
// transaction, and on ErrGatewayPartitionMissing create the partition out-of-band (its own
// transaction under the exclusive lock) and retry once.
func db_runInsertWithEnsure(
	ctx context.Context,
	t *testing.T,
	db *sql.DB,
	params queries.InsertGatewayEnvelopeV3Params,
) error {
	t.Helper()
	insertTx := func(ctx context.Context, q *queries.Queries) error {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(ctx, q, params)
		return err
	}
	err := xmtpd_db.RunInTx(ctx, db, &sql.TxOptions{}, insertTx)
	if errors.Is(err, xmtpd_db.ErrGatewayPartitionMissing) {
		require.NoError(t, xmtpd_db.EnsureGatewayPartitions(
			ctx, db, params.OriginatorNodeID, params.OriginatorSequenceID,
		))
		err = xmtpd_db.RunInTx(ctx, db, &sql.TxOptions{}, insertTx)
	}
	return err
}

func TestInsertGatewayEnvelopeWithCheckTransactional(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	for i := 1; i < 10; i++ {
		// Partitions don't exist yet, so the in-transaction insert reports the missing
		// partition; the caller creates it out-of-band and retries.
		require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
			OriginatorNodeID:     int32(i),
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(100),
		}))
	}
}

func TestInsertGatewayEnvelopeWithChecksTx_AutoCreateAndRetry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	const nodeID int32 = 42
	const seqID int64 = 1

	// First insert creates the partition out-of-band and retries.
	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(128),
	}))

	// A second insert into the SAME band should succeed without needing to create parts again.
	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID + 123, // still within [0..1_000_000) default band
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(128),
	}))
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

	const nodeID int32 = 99

	// Sequence values straddling a band boundary:
	seqLeft := xmtpd_db.GatewayEnvelopeBandWidth - 1 // falls into band [0, bw)
	seqRight := xmtpd_db.GatewayEnvelopeBandWidth    // falls into band [bw, 2*bw)

	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqLeft,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}))

	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqRight,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}))

	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqLeft + 123, // still within first band
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}))

	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqRight + 456, // still within second band
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}))
}

// TestInsertGatewayEnvelopeWithChecksTx_RollbackInsert verifies that rolling back the insert
// transaction discards the inserted row while the out-of-band-created partition persists, so a
// fresh insert into the same partition succeeds.
func TestInsertGatewayEnvelopeWithChecksTx_RollbackInsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	const nodeID int32 = 7
	const seqID int64 = 10

	params := queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}

	// Create the partition out-of-band, then insert in a transaction that we roll back.
	require.NoError(t, xmtpd_db.EnsureGatewayPartitions(ctx, db, nodeID, seqID))

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q := queries.New(db).WithTx(tx)
	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(ctx, q, params)
	require.NoError(t, err)
	_ = tx.Rollback()

	// The partition still exists, so re-inserting the same row in a new transaction succeeds.
	tx2, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	q2 := queries.New(db).WithTx(tx2)
	defer func(tx2 *sql.Tx) {
		_ = tx2.Rollback()
	}(tx2)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(ctx, q2, params)
	require.NoError(t, err)
}

func TestInsertGatewayEnvelopeWithChecksTx_Commit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	const nodeID int32 = 7
	const seqID int64 = 10

	require.NoError(t, db_runInsertWithEnsure(ctx, t, db, queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		Topic:                testutils.RandomBytes(32),
		OriginatorEnvelope:   testutils.RandomBytes(64),
	}))
}

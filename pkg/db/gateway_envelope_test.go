package db_test

import (
	"context"
	"database/sql"
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
) (queries.InsertGatewayEnvelopeParams, queries.IncrementUnsettledUsageParams) {
	insertParams := queries.InsertGatewayEnvelopeParams{
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
	db, _ := testutils.NewDB(t, ctx)

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
	)
	require.NoError(t, err)
	require.Equal(t, numInserted, int64(1))

	payerSpend, err := db.Query().GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend.TotalSpendPicodollars, int64(100))
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)

	originatorCongestion, err := db.Query().SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{OriginatorID: originatorID},
	)
	require.NoError(t, err)
	require.Equal(t, originatorCongestion, int64(1))
}

func TestPayerMustExist(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	payerID := testutils.RandomInt32()
	originatorID := int32(100)
	sequenceID := int64(10)

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	_, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
	)
	require.Error(t, err)
}

func TestInsertAndIncrementParallel(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

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

	require.Equal(t, totalInserted, int64(1))

	payerSpend, err := db.Query().GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend.TotalSpendPicodollars, int64(100))
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)

	originatorCongestion, err := db.Query().SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{OriginatorID: originatorID},
	)
	require.NoError(t, err)
	require.Equal(t, originatorCongestion, int64(1))
}

func TestInsertAndIncrementWithOutOfOrderSequenceID(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := int32(100)
	sequenceID := int64(10)

	insertParams, incrementParams := buildParams(payerID, originatorID, sequenceID, 100)

	_, err := xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
	)
	require.NoError(t, err)

	lowerSequenceID := int64(5)

	insertParams, incrementParams = buildParams(payerID, originatorID, lowerSequenceID, 100)

	_, err = xmtpd_db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		db,
		insertParams,
		incrementParams,
	)
	require.NoError(t, err)

	payerSpend, err := db.Query().GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend.LastSequenceID, sequenceID)
}

func TestInsertGatewayEnvelopeWithChecksStandalone(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			db,
			queries.InsertGatewayEnvelopeParams{
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
	db, _ := testutils.NewDB(t, ctx)
	tx, err := db.DB().BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	querier := db.Query().WithTx(tx)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	for i := 1; i < 10; i++ {
		_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandaloneWithQuerier(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeParams{
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
	db, _ := testutils.NewDB(t, ctx)

	const nodeID int32 = 42
	const seqID int64 = 1

	_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
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
		db,
		queries.InsertGatewayEnvelopeParams{
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
	db, _ := testutils.NewDB(t, ctx)

	const nodeID int32 = 7
	const seqID int64 = 10

	// Explicitly ensure parts up-front to simulate "partitions already exist" path.
	err := db.Query().EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		BandWidth:            1_000_000,
	})
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
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
	db, _ := testutils.NewDB(t, ctx)
	const nodeID int32 = 99

	// Sequence values straddling a band boundary:
	seqLeft := xmtpd_db.GatewayEnvelopeBandWidth - 1 // falls into band [0, bw)
	seqRight := xmtpd_db.GatewayEnvelopeBandWidth    // falls into band [bw, 2*bw)

	_, err := xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqRight,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     nodeID,
			OriginatorSequenceID: seqLeft + 123, // still within first band
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(64),
		},
	)
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		db,
		queries.InsertGatewayEnvelopeParams{
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
			queries.InsertGatewayEnvelopeParams{
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
			queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
	err = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		BandWidth:            xmtpd_db.GatewayEnvelopeBandWidth,
	})
	require.NoError(t, err)

	_, err = xmtpd_db.InsertGatewayEnvelopeWithChecksTransactional(
		ctx,
		q,
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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
		queries.InsertGatewayEnvelopeParams{
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

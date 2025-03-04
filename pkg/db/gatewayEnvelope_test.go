package db_test

import (
	"context"
	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"
	"sync"
	"sync/atomic"
	"testing"

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
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)
	// Create a payer
	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := testutils.RandomInt32()
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

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend, int64(100))
}

func TestPayerMustExist(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	payerID := testutils.RandomInt32()
	originatorID := testutils.RandomInt32()
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
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)
	// Create a payer
	payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
	originatorID := testutils.RandomInt32()
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

	for range numberOfInserts {
		wg.Add(1)
		go attemptInsert()
	}

	wg.Wait()

	require.Equal(t, totalInserted, int64(1))

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, payerSpend, int64(100))
}

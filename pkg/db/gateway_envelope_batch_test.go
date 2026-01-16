package db_test

import (
	"context"
	"database/sql"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// buildBatchInput builds a batch of gateway envelope parameters.
// From startSequenceID to startSequenceID + count - 1.
func buildBatchInput(
	payerID int32,
	originatorID int32,
	startSequenceID int64,
	count int,
	spendPerMessage int64,
) *types.GatewayEnvelopeBatch {
	batch := types.NewGatewayEnvelopeBatch()

	now := time.Now()
	for i := 0; i < count; i++ {
		batch.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: startSequenceID + int64(i),
			Topic:                testutils.RandomBytes(32),
			PayerID:              payerID,
			GatewayTime:          now,
			Expiry:               now.Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   testutils.RandomBytes(100),
			SpendPicodollars:     spendPerMessage,
		})
	}

	return batch
}

func TestBatchInsert_Basic(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		batchSize       = rand.Intn(10) + 1
		input           = buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
	)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize)*spendPerMessage, payerSpend.TotalSpendPicodollars)
	require.Equal(t, int64(batchSize), payerSpend.LastSequenceID)
}

func TestBatchInsert_OnlyEnvelopesBatch(t *testing.T) {
	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		querier      = queries.New(db)
		originatorID = testutils.RandomInt32()
	)

	batch := types.NewGatewayEnvelopeBatch()
	for i := 0; i < 3; i++ {
		batch.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: int64(i + 1),
			Topic:                testutils.RandomBytes(32),
			PayerID:              0,
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   testutils.RandomBytes(100),
			SpendPicodollars:     100,
		})
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		batch,
	)
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: 0},
	)
	require.NoError(t, err)

	// Unsettled usage should not be incremented, as all the payers are null.
	require.Equal(t, int64(0), payerSpend.TotalSpendPicodollars)
	require.Equal(t, int64(0), payerSpend.LastSequenceID)
}

func TestBatchInsert_EmptyInput(t *testing.T) {
	var (
		ctx   = context.Background()
		db, _ = testutils.NewRawDB(t, ctx)
		input = types.NewGatewayEnvelopeBatch()
	)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty input")
	require.Equal(t, int64(0), result)
}

func TestBatchInsert_DuplicatesIgnored(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		batchSize       = rand.Intn(10) + 1
		input           = buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
	)

	// Insert first batch.
	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	// Insert same batch again (duplicates).
	result, err = xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), result)

	// Verify usage was NOT double-counted.
	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize)*spendPerMessage, payerSpend.TotalSpendPicodollars)
}

func TestBatchInsert_PartialDuplicates(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
	)

	// Insert first batch (seq 1-5).
	input1 := buildBatchInput(payerID, originatorID, 1, 5, spendPerMessage)
	result1, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input1,
	)
	require.NoError(t, err)
	require.Equal(t, int64(5), result1)

	// Insert overlapping batch (seq 3-7). 3,4,5 are duplicates, 6,7 are new.
	input2 := buildBatchInput(payerID, originatorID, 3, 5, spendPerMessage)
	result2, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input2,
	)
	require.NoError(t, err)
	require.Equal(t, int64(2), result2)

	// Verify usage is 5 + 2 = 7 messages.
	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(700), payerSpend.TotalSpendPicodollars)
}

func TestBatchInsert_MultipleOriginators(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID1   = int32(100)
		originatorID2   = int32(200)
		spendPerMessage = int64(100)
		batch1          = buildBatchInput(payerID, originatorID1, 1, 4, spendPerMessage)
		batch2          = buildBatchInput(payerID, originatorID2, 1, 4, spendPerMessage)
	)

	input := types.NewGatewayEnvelopeBatch()
	for _, envelope := range batch1.Envelopes {
		input.Add(envelope)
	}

	for _, envelope := range batch2.Envelopes {
		input.Add(envelope)
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(8), result)
}

func TestBatchInsert_InvalidSequenceOrder(t *testing.T) {
	// Function has to handle ordering of sequence IDs for each originator node.
	var (
		ctx           = context.Background()
		db, _         = testutils.NewRawDB(t, ctx)
		originatorID1 = int32(100)
		originatorID2 = int32(200)
	)

	testCases := []struct {
		name                  string
		originatorNodeIds     []int32
		originatorSequenceIds []int64
	}{
		{
			name: "out of order for same originator",
			originatorNodeIds: []int32{
				originatorID1,
				originatorID1,
				originatorID2,
				originatorID2,
			},
			originatorSequenceIds: []int64{1, 2, 2, 1},
		},
		{
			name: "duplicate sequence ID for same originator",
			originatorNodeIds: []int32{
				originatorID1,
				originatorID1,
				originatorID2,
				originatorID2,
			},
			originatorSequenceIds: []int64{1, 2, 2, 2},
		},
		{
			name:                  "single originator out of order",
			originatorNodeIds:     []int32{originatorID1, originatorID1, originatorID1},
			originatorSequenceIds: []int64{5, 3, 10},
		},
		{
			name:                  "single originator duplicate",
			originatorNodeIds:     []int32{originatorID1, originatorID1},
			originatorSequenceIds: []int64{5, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payerID := testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
			count := len(tc.originatorNodeIds)

			input := types.NewGatewayEnvelopeBatch()

			now := time.Now()
			for i := 0; i < count; i++ {
				input.Add(types.GatewayEnvelopeRow{
					OriginatorNodeID:     tc.originatorNodeIds[i],
					OriginatorSequenceID: tc.originatorSequenceIds[i],
					Topic:                testutils.RandomBytes(32),
					PayerID:              payerID,
					GatewayTime:          now,
					Expiry:               now.Add(24 * time.Hour).Unix(),
					OriginatorEnvelope:   testutils.RandomBytes(100),
					SpendPicodollars:     100,
				})
			}

			_, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				input,
			)
			require.NoError(t, err)
		})
	}
}

func TestBatchInsert_NullPayerID(t *testing.T) {
	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		originatorID = int32(100)

		// 0 is a null payer ID.
		input = buildBatchInput(0, originatorID, 1, 3, 100)
	)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	// Verify payer_id is NULL in the database.
	var payerID sql.NullInt32
	err = db.QueryRowContext(
		ctx,
		`SELECT payer_id FROM gateway_envelopes_meta 
		 WHERE originator_node_id = $1 AND originator_sequence_id = $2`,
		originatorID,
		int64(1),
	).Scan(&payerID)
	require.NoError(t, err)
	require.False(t, payerID.Valid, "payer_id should be NULL")
}

func TestBatchInsert_PayerMustExist(t *testing.T) {
	var (
		ctx                = context.Background()
		db, _              = testutils.NewRawDB(t, ctx)
		nonExistentPayerID = testutils.RandomInt32()
		originatorID       = int32(100)
		input              = buildBatchInput(nonExistentPayerID, originatorID, 1, 3, 100)
	)

	_, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)

	// Note, the test does not create the payer, so it should fail with a FK constraint violation.
	require.Error(t, err)
}

func TestBatchInsert_BandBoundaries(t *testing.T) {
	t.Parallel()

	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		payerID      = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID = int32(99)
		seqLeft      = xmtpd_db.GatewayEnvelopeBandWidth - 1 // falls into band [0, bw)
		seqRight     = xmtpd_db.GatewayEnvelopeBandWidth + 1 // falls into band [bw, 2*bw)
	)

	message1 := buildBatchInput(payerID, originatorID, seqLeft, 1, 100)
	message2 := buildBatchInput(payerID, originatorID, seqRight, 1, 100)

	input := types.NewGatewayEnvelopeBatch()
	for _, envelope := range message1.Envelopes {
		input.Add(envelope)
	}
	for _, envelope := range message2.Envelopes {
		input.Add(envelope)
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(2), result)
}

func TestBatchInsert_Parallel(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		numGoroutines   = 10
		batchSize       = 5
		totalInserted   = int64(0)
		input           = buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
	)

	// First insert to create partitions (avoid DDL deadlocks).
	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)

	atomic.AddInt64(&totalInserted, result)

	var wg sync.WaitGroup

	// Parallel inserts with different sequence ranges.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startSeq int64) {
			defer wg.Done()
			p := buildBatchInput(payerID, originatorID, startSeq, batchSize, spendPerMessage)
			n, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				p,
			)
			require.NoError(t, err)
			atomic.AddInt64(&totalInserted, n)
		}(int64(1 + (i+1)*batchSize)) // Non-overlapping ranges.
	}

	wg.Wait()

	expectedTotal := int64((numGoroutines + 1) * batchSize)
	require.Equal(t, expectedTotal, totalInserted)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, expectedTotal*spendPerMessage, payerSpend.TotalSpendPicodollars)
}

func TestBatchInsert_ParallelDuplicates(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		numGoroutines   = 20
		batchSize       = 5
		totalInserted   = int64(0)
		input           = buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
	)

	// First insert to create partitions.
	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)

	atomic.AddInt64(&totalInserted, result)

	var wg sync.WaitGroup

	// All goroutines try to insert the SAME batch (duplicates).
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
			n, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				p,
			)
			require.NoError(t, err)
			atomic.AddInt64(&totalInserted, n)
		}()
	}

	wg.Wait()

	// Only the first batch should have been inserted.
	require.Equal(t, int64(batchSize), totalInserted)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize)*spendPerMessage, payerSpend.TotalSpendPicodollars)
}

func TestBatchInsert_LargeBatch(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		batchSize       = 1000
		spendPerMessage = int64(50)
		input           = buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
	)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize)*spendPerMessage, payerSpend.TotalSpendPicodollars)
}

func TestBatchInsert_PreexistingPartitions(t *testing.T) {
	t.Parallel()

	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		querier      = queries.New(db)
		payerID      = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID = int32(7)
		input        = buildBatchInput(payerID, originatorID, 1, 5, 100)
	)

	// Pre-create partitions.
	err := querier.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     originatorID,
		OriginatorSequenceID: 1,
		BandWidth:            xmtpd_db.GatewayEnvelopeBandWidth,
	})
	require.NoError(t, err)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(5), result)
}

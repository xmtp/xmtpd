package db_test

import (
	"context"
	"database/sql"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	for i := range count {
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
	for i := range 3 {
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
			for i := range count {
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
	for i := range numGoroutines {
		wg.Add(1)
		go func(startSeq int64) {
			defer wg.Done()
			p := buildBatchInput(payerID, originatorID, startSeq, batchSize, spendPerMessage)
			n, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				p,
			)
			assert.NoError(t, err)
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
	for range numGoroutines {
		wg.Go(func() {
			p := buildBatchInput(payerID, originatorID, 1, batchSize, spendPerMessage)
			n, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				p,
			)
			require.NoError(t, err)
			atomic.AddInt64(&totalInserted, n)
		})
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

func TestBatchInsertV2_Basic(t *testing.T) {
	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		payerID      = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID = int32(100)
		batchSize    = 6
	)

	batch := types.NewGatewayEnvelopeBatch()
	now := time.Now()
	for i := range batchSize {
		batch.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: int64(i + 1),
			Topic:                testutils.RandomBytes(32),
			PayerID:              payerID,
			GatewayTime:          now,
			Expiry:               now.Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   testutils.RandomBytes(100),
			SpendPicodollars:     100,
			IsReserved:           i%2 == 0, // alternate: true, false, true, false, ...
		})
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		ctx,
		db,
		batch,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	// Verify all envelopes were inserted (meta + blob rows).
	var metaCount, blobCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gateway_envelopes_meta WHERE originator_node_id = $1`,
		originatorID,
	).Scan(&metaCount)
	require.NoError(t, err)
	require.Equal(t, batchSize, metaCount)

	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gateway_envelope_blobs WHERE originator_node_id = $1`,
		originatorID,
	).Scan(&blobCount)
	require.NoError(t, err)
	require.Equal(t, batchSize, blobCount)
}

func TestBatchInsertV2_ReservedTopicsNoUsageNoCongestion(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(500)
	)

	batch := types.NewGatewayEnvelopeBatch()
	now := time.Now()
	// 3 reserved + 2 non-reserved = 5 total
	for i := range 5 {
		batch.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: int64(i + 1),
			Topic:                testutils.RandomBytes(32),
			PayerID:              payerID,
			GatewayTime:          now,
			Expiry:               now.Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   testutils.RandomBytes(100),
			SpendPicodollars:     spendPerMessage,
			IsReserved:           i < 3, // first 3 are reserved
		})
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		ctx,
		db,
		batch,
	)
	require.NoError(t, err)
	require.Equal(t, int64(5), result)

	// Unsettled usage should only reflect the 2 non-reserved messages.
	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(2)*spendPerMessage, payerSpend.TotalSpendPicodollars,
		"spend_picodollars should be sum of non-reserved only")

	// Verify message_count in unsettled_usage is 2 (not 5).
	var messageCount int32
	err = db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(message_count), 0)::INT FROM unsettled_usage
		 WHERE payer_id = $1`,
		payerID,
	).Scan(&messageCount)
	require.NoError(t, err)
	require.Equal(t, int32(2), messageCount,
		"unsettled_usage message_count should be 2 (non-reserved only)")

	// Originator congestion should also be 2 (not 5).
	// Pass 0 for both boundaries to fetch all rows (the query treats 0 as "no filter").
	congestion, err := querier.SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{
			OriginatorID:        originatorID,
			MinutesSinceEpochGt: 0,
			MinutesSinceEpochLt: 0,
		},
	)
	require.NoError(t, err)
	require.Equal(t, int64(2), congestion,
		"originator_congestion num_messages should be 2 (non-reserved only)")
}

func TestBatchInsertV2_AllReservedNoSideEffects(t *testing.T) {
	var (
		ctx          = context.Background()
		db, _        = testutils.NewRawDB(t, ctx)
		querier      = queries.New(db)
		payerID      = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID = int32(200)
		batchSize    = 5
	)

	batch := types.NewGatewayEnvelopeBatch()
	now := time.Now()
	for i := range batchSize {
		batch.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: int64(i + 1),
			Topic:                testutils.RandomBytes(32),
			PayerID:              payerID,
			GatewayTime:          now,
			Expiry:               now.Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   testutils.RandomBytes(100),
			SpendPicodollars:     100,
			IsReserved:           true, // ALL reserved
		})
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		ctx,
		db,
		batch,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	// Zero unsettled_usage rows for this payer.
	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), payerSpend.TotalSpendPicodollars,
		"all-reserved batch should produce zero unsettled_usage spend")
	require.Equal(t, int64(0), payerSpend.LastSequenceID,
		"all-reserved batch should produce zero unsettled_usage last_sequence_id")

	// Zero congestion rows for this originator.
	congestion, err := querier.SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{
			OriginatorID:        originatorID,
			MinutesSinceEpochGt: 0,
			MinutesSinceEpochLt: 0,
		},
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), congestion,
		"all-reserved batch should produce zero originator_congestion")

	// All 5 meta + blob rows should still be inserted.
	var metaCount, blobCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gateway_envelopes_meta WHERE originator_node_id = $1`,
		originatorID,
	).Scan(&metaCount)
	require.NoError(t, err)
	require.Equal(t, batchSize, metaCount, "all meta rows should be inserted")

	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gateway_envelope_blobs WHERE originator_node_id = $1`,
		originatorID,
	).Scan(&blobCount)
	require.NoError(t, err)
	require.Equal(t, batchSize, blobCount, "all blob rows should be inserted")
}

func TestBulkFindOrCreatePayers_DuplicateAddresses(t *testing.T) {
	var (
		ctx   = context.Background()
		db, _ = testutils.NewRawDB(t, ctx)
		q     = queries.New(db)
		addr1 = testutils.RandomAddress().Hex()
		addr2 = testutils.RandomAddress().Hex()
	)

	// Pass duplicate addresses: [addr1, addr2, addr1, addr1]
	rows, err := q.BulkFindOrCreatePayers(ctx, []string{addr1, addr2, addr1, addr1})
	require.NoError(t, err)

	// Should return exactly 2 unique addresses.
	addressSet := make(map[string]int32)
	for _, row := range rows {
		addressSet[row.Address] = row.ID
	}

	require.Len(t, addressSet, 2, "should return exactly 2 unique addresses")
	require.Contains(t, addressSet, addr1)
	require.Contains(t, addressSet, addr2)

	// Each address should have a valid non-zero ID.
	require.NotZero(t, addressSet[addr1], "addr1 should have a non-zero ID")
	require.NotZero(t, addressSet[addr2], "addr2 should have a non-zero ID")
}

func TestBatchInsertV2_ConcurrentIdempotency(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(300)
		spendPerMessage = int64(100)
		batchSize       = 5
	)

	// Build a single batch that both goroutines will try to insert.
	makeBatch := func() *types.GatewayEnvelopeBatch {
		batch := types.NewGatewayEnvelopeBatch()
		now := time.Now()
		for i := range batchSize {
			batch.Add(types.GatewayEnvelopeRow{
				OriginatorNodeID:     originatorID,
				OriginatorSequenceID: int64(i + 1),
				Topic:                []byte("fixed-topic-for-idempotency"),
				PayerID:              payerID,
				GatewayTime:          now,
				Expiry:               now.Add(24 * time.Hour).Unix(),
				OriginatorEnvelope:   testutils.RandomBytes(100),
				SpendPicodollars:     spendPerMessage,
				IsReserved:           false,
			})
		}
		return batch
	}

	// First insert to create partitions (avoids DDL deadlocks in concurrent test).
	firstBatch := makeBatch()
	result, err := xmtpd_db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		ctx,
		db,
		firstBatch,
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), result)

	// Now run concurrent duplicate inserts.
	var (
		wg            sync.WaitGroup
		totalInserted atomic.Int64
	)

	numGoroutines := 10
	for range numGoroutines {
		wg.Go(func() {
			batch := makeBatch()
			n, insertErr := xmtpd_db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
				ctx,
				db,
				batch,
			)
			require.NoError(t, insertErr)
			totalInserted.Add(n)
		})
	}

	wg.Wait()

	// All concurrent duplicates should have inserted 0 rows.
	require.Equal(t, int64(0), totalInserted.Load(),
		"concurrent duplicate inserts should insert 0 new rows")

	// Total meta rows should be exactly batchSize.
	var metaCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gateway_envelopes_meta WHERE originator_node_id = $1`,
		originatorID,
	).Scan(&metaCount)
	require.NoError(t, err)
	require.Equal(t, batchSize, metaCount,
		"total meta rows should be exactly batchSize")

	// Unsettled usage should reflect only the first insert (incremented once).
	payerSpend, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{PayerID: payerID},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize)*spendPerMessage, payerSpend.TotalSpendPicodollars,
		"unsettled_usage should be incremented exactly once")

	// Congestion should also reflect only the first insert.
	congestion, err := querier.SumOriginatorCongestion(
		ctx,
		queries.SumOriginatorCongestionParams{
			OriginatorID:        originatorID,
			MinutesSinceEpochGt: 0,
			MinutesSinceEpochLt: 0,
		},
	)
	require.NoError(t, err)
	require.Equal(t, int64(batchSize), congestion,
		"originator_congestion should be incremented exactly once")
}

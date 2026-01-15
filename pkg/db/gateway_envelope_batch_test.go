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
) queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams {
	input := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     make([]int32, count),
		OriginatorSequenceIds: make([]int64, count),
		Topics:                make([][]byte, count),
		PayerIds:              make([]int32, count),
		GatewayTimes:          make([]time.Time, count),
		Expiries:              make([]int64, count),
		OriginatorEnvelopes:   make([][]byte, count),
		SpendPicodollars:      make([]int64, count),
	}

	now := time.Now()
	for i := 0; i < count; i++ {
		input.OriginatorNodeIds[i] = originatorID
		input.OriginatorSequenceIds[i] = startSequenceID + int64(i)
		input.Topics[i] = testutils.RandomBytes(32)
		input.PayerIds[i] = payerID
		input.GatewayTimes[i] = now
		input.Expiries[i] = now.Add(24 * time.Hour).Unix()
		input.OriginatorEnvelopes[i] = testutils.RandomBytes(100)
		input.SpendPicodollars[i] = spendPerMessage
	}

	return input
}

func TestBatchInsert_Basic(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		batchSize       = rand.Intn(100)
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

	input := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     []int32{originatorID, originatorID, originatorID},
		OriginatorSequenceIds: []int64{1, 2, 3},
		Topics: [][]byte{
			testutils.RandomBytes(32),
			testutils.RandomBytes(32),
			testutils.RandomBytes(32),
		},
		PayerIds:     []int32{0, 0, 0},
		GatewayTimes: []time.Time{time.Now(), time.Now(), time.Now()},
		Expiries: []int64{
			time.Now().Add(24 * time.Hour).Unix(),
			time.Now().Add(24 * time.Hour).Unix(),
			time.Now().Add(24 * time.Hour).Unix(),
		},
		OriginatorEnvelopes: [][]byte{
			testutils.RandomBytes(100),
			testutils.RandomBytes(100),
			testutils.RandomBytes(100),
		},
		// It doesn't matter what the spend picodollars are, as all the payers are null.
		// SQL will not insert any rows into unsettled_usage.
		SpendPicodollars: []int64{1, 2, 3},
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
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
		input = queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{}
	)

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(0), result)
}

func TestBatchInsert_ArrayLengthMismatch(t *testing.T) {
	var (
		ctx   = context.Background()
		db, _ = testutils.NewRawDB(t, ctx)
	)

	testCases := []struct {
		name   string
		modify func(*queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams)
	}{
		{
			name: "OriginatorSequenceIds shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.OriginatorSequenceIds = p.OriginatorSequenceIds[:len(p.OriginatorSequenceIds)-1]
			},
		},
		{
			name: "Topics shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.Topics = p.Topics[:len(p.Topics)-1]
			},
		},
		{
			name: "PayerIds shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.PayerIds = p.PayerIds[:len(p.PayerIds)-1]
			},
		},
		{
			name: "GatewayTimes shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.GatewayTimes = p.GatewayTimes[:len(p.GatewayTimes)-1]
			},
		},
		{
			name: "Expiries shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.Expiries = p.Expiries[:len(p.Expiries)-1]
			},
		},
		{
			name: "OriginatorEnvelopes shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.OriginatorEnvelopes = p.OriginatorEnvelopes[:len(p.OriginatorEnvelopes)-1]
			},
		},
		{
			name: "SpendPicodollars shorter",
			modify: func(p *queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams) {
				p.SpendPicodollars = p.SpendPicodollars[:len(p.SpendPicodollars)-1]
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				payerID = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
				input   = buildBatchInput(payerID, int32(100), 1, 5, 100)
			)

			tc.modify(&input)

			_, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				input,
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "input array length mismatch")
		})
	}
}

func TestBatchInsert_DuplicatesIgnored(t *testing.T) {
	var (
		ctx             = context.Background()
		db, _           = testutils.NewRawDB(t, ctx)
		querier         = queries.New(db)
		payerID         = testutils.CreatePayer(t, db, testutils.RandomAddress().Hex())
		originatorID    = int32(100)
		spendPerMessage = int64(100)
		batchSize       = rand.Intn(100)
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
	)

	input := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     []int32{originatorID1, originatorID1, originatorID2, originatorID2},
		OriginatorSequenceIds: []int64{1, 2, 1, 2},
		Topics:                make([][]byte, 4),
		PayerIds:              []int32{payerID, payerID, payerID, payerID},
		GatewayTimes:          make([]time.Time, 4),
		Expiries:              make([]int64, 4),
		OriginatorEnvelopes:   make([][]byte, 4),
		SpendPicodollars: []int64{
			spendPerMessage,
			spendPerMessage,
			spendPerMessage,
			spendPerMessage,
		},
	}

	now := time.Now()
	for i := 0; i < 4; i++ {
		input.Topics[i] = testutils.RandomBytes(32)
		input.GatewayTimes[i] = now
		input.Expiries[i] = now.Add(24 * time.Hour).Unix()
		input.OriginatorEnvelopes[i] = testutils.RandomBytes(100)
	}

	result, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
		ctx,
		db,
		input,
	)
	require.NoError(t, err)
	require.Equal(t, int64(4), result)
}

func TestBatchInsert_InvalidSequenceOrder(t *testing.T) {
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

			input := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
				OriginatorNodeIds:     tc.originatorNodeIds,
				OriginatorSequenceIds: tc.originatorSequenceIds,
				Topics:                make([][]byte, count),
				PayerIds:              make([]int32, count),
				GatewayTimes:          make([]time.Time, count),
				Expiries:              make([]int64, count),
				OriginatorEnvelopes:   make([][]byte, count),
				SpendPicodollars:      make([]int64, count),
			}

			now := time.Now()
			for i := 0; i < count; i++ {
				input.Topics[i] = testutils.RandomBytes(32)
				input.PayerIds[i] = payerID
				input.GatewayTimes[i] = now
				input.Expiries[i] = now.Add(24 * time.Hour).Unix()
				input.OriginatorEnvelopes[i] = testutils.RandomBytes(100)
				input.SpendPicodollars[i] = 100
			}

			_, err := xmtpd_db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				db,
				input,
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "sequence IDs must be strictly ascending")
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

	input := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     []int32{originatorID, originatorID},
		OriginatorSequenceIds: []int64{seqLeft, seqRight},
		Topics:                [][]byte{testutils.RandomBytes(32), testutils.RandomBytes(32)},
		PayerIds:              []int32{payerID, payerID},
		GatewayTimes:          []time.Time{time.Now(), time.Now()},
		Expiries: []int64{
			time.Now().Add(24 * time.Hour).Unix(),
			time.Now().Add(24 * time.Hour).Unix(),
		},
		OriginatorEnvelopes: [][]byte{testutils.RandomBytes(100), testutils.RandomBytes(100)},
		SpendPicodollars:    []int64{100, 100},
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

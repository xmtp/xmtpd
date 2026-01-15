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
) []types.GatewayEnvelopeRow {
	input := make([]types.GatewayEnvelopeRow, count)

	now := time.Now()
	for i := 0; i < count; i++ {
		input[i].OriginatorNodeID = originatorID
		input[i].OriginatorSequenceID = startSequenceID + int64(i)
		input[i].Topic = testutils.RandomBytes(32)
		input[i].PayerID = payerID
		input[i].GatewayTime = now
		input[i].Expiry = now.Add(24 * time.Hour).Unix()
		input[i].OriginatorEnvelope = testutils.RandomBytes(100)
		input[i].SpendPicodollars = spendPerMessage
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
		batchSize       = rand.Intn(10)
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

	input := make([]types.GatewayEnvelopeRow, 3)
	for i := 0; i < 3; i++ {
		input[i].OriginatorNodeID = originatorID
		input[i].OriginatorSequenceID = int64(i + 1)
		input[i].Topic = testutils.RandomBytes(32)
		input[i].PayerID = 0
		input[i].GatewayTime = time.Now()
		input[i].Expiry = time.Now().Add(24 * time.Hour).Unix()
		input[i].OriginatorEnvelope = testutils.RandomBytes(100)
		input[i].SpendPicodollars = 100
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
		input = []types.GatewayEnvelopeRow{}
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
		batchSize       = rand.Intn(10)
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
		messages1       = buildBatchInput(payerID, originatorID1, 1, 4, spendPerMessage)
		messages2       = buildBatchInput(payerID, originatorID2, 1, 4, spendPerMessage)
	)

	input := append(messages1, messages2...)

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

			input := make([]types.GatewayEnvelopeRow, count)

			now := time.Now()
			for i := 0; i < count; i++ {
				input[i].Topic = testutils.RandomBytes(32)
				input[i].PayerID = payerID
				input[i].GatewayTime = now
				input[i].Expiry = now.Add(24 * time.Hour).Unix()
				input[i].OriginatorEnvelope = testutils.RandomBytes(100)
				input[i].SpendPicodollars = 100
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
	input := append(message1, message2...)

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

// TestGatewayEnvelopeRow_ValueParseable verifies that the Value() method
// produces a string that PostgreSQL can parse as a gateway_envelope_row composite type.
func TestGatewayEnvelopeRow_ValueParseable(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	// Verify the type exists after migrations.
	var typeExists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'gateway_envelope_row')`,
	).Scan(&typeExists)
	require.NoError(t, err)
	require.True(t, typeExists, "gateway_envelope_row type should exist after migrations")

	// Create a row with representative data including bytea fields.
	row := types.GatewayEnvelopeRow{
		OriginatorNodeID:     42,
		OriginatorSequenceID: 123456,
		Topic:                []byte{0xDE, 0xAD, 0xBE, 0xEF},
		PayerID:              7,
		GatewayTime:          time.Date(2025, 6, 15, 10, 30, 0, 123456000, time.UTC),
		Expiry:               1750000000,
		OriginatorEnvelope:   []byte{0xCA, 0xFE, 0xBA, 0xBE},
		SpendPicodollars:     999,
	}

	// Get the serialized value.
	val, err := row.Value()
	require.NoError(t, err)

	literal, ok := val.(string)
	require.True(t, ok, "Value() should return a string")
	t.Logf("Serialized literal: %s", literal)

	// Test 1: Parse as a single element.
	var parsedNodeID int32
	err = db.QueryRowContext(ctx,
		`SELECT (v).originator_node_id FROM (SELECT $1::gateway_envelope_row AS v) sub`,
		literal,
	).Scan(&parsedNodeID)
	require.NoError(t, err, "PostgreSQL should be able to parse the single element literal")
	require.Equal(t, row.OriginatorNodeID, parsedNodeID)

	// Test 2: Parse as an array element (this is how pq.Array uses it).
	var parsedSeqID int64
	err = db.QueryRowContext(ctx,
		`SELECT (arr[1]).originator_sequence_id 
		 FROM (SELECT ARRAY[$1::gateway_envelope_row] AS arr) sub`,
		literal,
	).Scan(&parsedSeqID)
	require.NoError(t, err, "PostgreSQL should be able to parse element inside array")
	require.Equal(t, row.OriginatorSequenceID, parsedSeqID)

	// Test 3: Parse bytea field to verify hex encoding.
	var parsedTopic []byte
	err = db.QueryRowContext(ctx,
		`SELECT (v).topic FROM (SELECT $1::gateway_envelope_row AS v) sub`,
		literal,
	).Scan(&parsedTopic)
	require.NoError(t, err, "PostgreSQL should be able to parse bytea field")
	require.Equal(t, row.Topic, parsedTopic)

	// Test 4: Parse timestamp field.
	var parsedTime time.Time
	err = db.QueryRowContext(ctx,
		`SELECT (v).gateway_time FROM (SELECT $1::gateway_envelope_row AS v) sub`,
		literal,
	).Scan(&parsedTime)
	require.NoError(t, err, "PostgreSQL should be able to parse timestamp field")
	// Compare truncated to microseconds (PostgreSQL precision).
	require.Equal(t, row.GatewayTime.Truncate(time.Microsecond), parsedTime.UTC())
}

// TestGatewayEnvelopeRow_ArrayValuer verifies that multiple rows can be
// serialized and parsed as a PostgreSQL array of composite types.
func TestGatewayEnvelopeRow_ArrayValuer(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	rows := []types.GatewayEnvelopeRow{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 100,
			Topic:                []byte{0x01, 0x02},
			PayerID:              10,
			GatewayTime:          time.Now().UTC().Truncate(time.Microsecond),
			Expiry:               1750000000,
			OriginatorEnvelope:   []byte{0xAA, 0xBB},
			SpendPicodollars:     50,
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 200,
			Topic:                []byte{0x03, 0x04},
			PayerID:              20,
			GatewayTime:          time.Now().UTC().Truncate(time.Microsecond),
			Expiry:               1750000001,
			OriginatorEnvelope:   []byte{0xCC, 0xDD},
			SpendPicodollars:     75,
		},
	}

	// Build array literal manually: ARRAY[elem1, elem2]
	var literals []string
	for _, r := range rows {
		val, err := r.Value()
		require.NoError(t, err)
		literals = append(literals, val.(string))
	}

	// Query: count elements and sum a field.
	var count int
	var totalSpend int64
	query := `
		SELECT 
			array_length(arr, 1),
			(SELECT SUM((e).spend_picodollars) FROM unnest(arr) AS e)
		FROM (
			SELECT ARRAY[$1::gateway_envelope_row, $2::gateway_envelope_row] AS arr
		) sub
	`
	err := db.QueryRowContext(ctx, query, literals[0], literals[1]).Scan(&count, &totalSpend)
	require.NoError(t, err, "PostgreSQL should parse array of composite types")
	require.Equal(t, 2, count)
	require.Equal(t, int64(125), totalSpend)
}

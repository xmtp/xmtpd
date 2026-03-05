package fees

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	rateMessageFee      = 100
	rateStorageFee      = 50
	rateCongestionFee   = 200
	targetRatePerMinute = 4
)

func setupCalculator() *FeeCalculator {
	rates := &Rates{
		MessageFee:          rateMessageFee,
		StorageFee:          rateStorageFee,
		CongestionFee:       rateCongestionFee,
		TargetRatePerMinute: targetRatePerMinute,
	}
	ratesFetcher := NewFixedRatesFetcher(rates)
	return NewFeeCalculator(ratesFetcher)
}

func addCongestion(
	t *testing.T,
	querier *queries.Queries,
	originatorID uint32,
	minutesSinceEpoch int32,
	congestion int,
) {
	for range congestion {
		err := querier.IncrementOriginatorCongestion(
			context.Background(),
			queries.IncrementOriginatorCongestionParams{
				OriginatorID:      int32(originatorID),
				MinutesSinceEpoch: minutesSinceEpoch,
			},
		)

		require.NoError(t, err)
	}
}

func TestCalculateBaseFee(t *testing.T) {
	calculator := setupCalculator()

	messageTime := time.Now()
	messageSize := int64(100)
	storageDurationDays := int64(1)

	baseFee, err := calculator.CalculateBaseFee(
		messageTime,
		messageSize,
		uint32(storageDurationDays),
	)
	require.NoError(t, err)

	expectedFee := rateMessageFee + (rateStorageFee * messageSize * storageDurationDays)
	require.Equal(t, currency.PicoDollar(expectedFee), baseFee)
}

func TestCalculateCongestionFee(t *testing.T) {
	calculator := setupCalculator()
	db, _ := testutils.NewRawDB(t, context.Background())

	ctx := context.Background()
	querier := queries.New(db)
	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)

	// Should return 0 if no congestion
	congestionFee, err := calculator.CalculateCongestionFee(
		ctx,
		querier,
		messageTime,
		originatorID,
	)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(0), congestionFee)

	// Congestion rate is 100 because this is double the max
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 8)
	congestionFee, err = calculator.CalculateCongestionFee(
		ctx,
		querier,
		messageTime,
		originatorID,
	)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(20000), congestionFee)
}

func TestCongestionFeeParity_BatchVsSequential(t *testing.T) {
	calculator := setupCalculator()
	ctx := context.Background()

	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)
	seedCongestion := 3

	// --- Sequential path: start from seeded state, compute 10 fees, incrementing DB after each ---
	seqDB, _ := testutils.NewRawDB(t, ctx)
	seqQuerier := queries.New(seqDB)
	seqOriginatorID := uint32(testutils.RandomInt32())

	addCongestion(t, seqQuerier, seqOriginatorID, minutesSinceEpoch, seedCongestion)

	sequentialFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := calculator.CalculateCongestionFee(
			ctx, seqQuerier, messageTime, seqOriginatorID,
		)
		require.NoError(t, err)
		sequentialFees[i] = fee

		// Increment congestion in DB to simulate the message being committed.
		addCongestion(t, seqQuerier, seqOriginatorID, minutesSinceEpoch, 1)
	}

	// --- Batched path: same seed, compute 10 fees using BatchFeeCalculator ---
	batchDB, _ := testutils.NewRawDB(t, ctx)
	batchQuerier := queries.New(batchDB)
	batchOriginatorID := uint32(testutils.RandomInt32())

	addCongestion(t, batchQuerier, batchOriginatorID, minutesSinceEpoch, seedCongestion)

	batchCalc := calculator.NewBatchFeeCalculator(ctx, batchQuerier, batchOriginatorID)
	batchedFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := batchCalc.CalculateCongestionFee(messageTime)
		require.NoError(t, err)
		batchedFees[i] = fee
	}

	// Assert every fee matches between sequential and batched.
	for i := range 10 {
		require.Equal(t, sequentialFees[i], batchedFees[i],
			"fee mismatch at message %d: sequential=%d batched=%d",
			i, sequentialFees[i], batchedFees[i])
	}
}

func TestBatchFeeCalculator_SameMinute(t *testing.T) {
	calculator := setupCalculator()
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(db)
	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)

	// Seed some congestion
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 3)

	// Use BatchFeeCalculator
	batchCalc := calculator.NewBatchFeeCalculator(ctx, querier, originatorID)
	batchFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := batchCalc.CalculateCongestionFee(messageTime)
		require.NoError(t, err)
		batchFees[i] = fee
	}

	// Compare against sequential CalculateCongestionFee calls
	seqDB, _ := testutils.NewRawDB(t, ctx)
	seqQuerier := queries.New(seqDB)
	seqOriginatorID := uint32(testutils.RandomInt32())
	addCongestion(t, seqQuerier, seqOriginatorID, minutesSinceEpoch, 3)

	for i := range 10 {
		seqFee, err := calculator.CalculateCongestionFee(
			ctx, seqQuerier, messageTime, seqOriginatorID,
		)
		require.NoError(t, err)
		require.Equal(t, seqFee, batchFees[i],
			"fee mismatch at message %d", i)
		addCongestion(t, seqQuerier, seqOriginatorID, minutesSinceEpoch, 1)
	}
}

func TestBatchFeeCalculator_CrossMinuteBoundary(t *testing.T) {
	calculator := setupCalculator()
	ctx := context.Background()

	now := time.Now()
	minute := utils.MinutesSinceEpoch(now)
	timeMinute0 := time.Unix(int64(minute)*60, 0)
	timeMinute1 := time.Unix(int64(minute+1)*60, 0)

	// Message schedule: 3 in minute0, then 5 in minute1.
	// Minute1's sliding window includes minute0, so cross-minute tracking matters.
	messageTimes := []time.Time{
		timeMinute0, timeMinute0, timeMinute0,
		timeMinute1, timeMinute1, timeMinute1, timeMinute1, timeMinute1,
	}

	seedCongestion := 2

	// --- Batch path ---
	batchDB, _ := testutils.NewRawDB(t, ctx)
	batchQuerier := queries.New(batchDB)
	batchOriginatorID := uint32(testutils.RandomInt32())
	addCongestion(t, batchQuerier, batchOriginatorID, minute, seedCongestion)

	batchCalc := calculator.NewBatchFeeCalculator(ctx, batchQuerier, batchOriginatorID)
	batchFees := make([]currency.PicoDollar, len(messageTimes))
	for i, mt := range messageTimes {
		fee, err := batchCalc.CalculateCongestionFee(mt)
		require.NoError(t, err)
		batchFees[i] = fee
	}

	// --- Sequential path: commit each message to DB before computing the next fee ---
	seqDB, _ := testutils.NewRawDB(t, ctx)
	seqQuerier := queries.New(seqDB)
	seqOriginatorID := uint32(testutils.RandomInt32())
	addCongestion(t, seqQuerier, seqOriginatorID, minute, seedCongestion)

	for i, mt := range messageTimes {
		seqMinute := utils.MinutesSinceEpoch(mt)
		seqFee, err := calculator.CalculateCongestionFee(
			ctx, seqQuerier, mt, seqOriginatorID,
		)
		require.NoError(t, err)
		require.Equal(t, seqFee, batchFees[i],
			"fee mismatch at message %d (minute %d)", i, seqMinute)

		// Commit this message to the DB so the next sequential call sees it
		addCongestion(t, seqQuerier, seqOriginatorID, seqMinute, 1)
	}
}

func TestOverflow(t *testing.T) {
	calculator := setupCalculator()

	messageTime := time.Now()

	// Test overflow in CalculateBaseFee
	messageSize := math.MaxInt64
	storageDurationDays := math.MaxInt64

	_, err := calculator.CalculateBaseFee(
		messageTime,
		int64(messageSize),
		uint32(storageDurationDays),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "storage fee calculation overflow")
}

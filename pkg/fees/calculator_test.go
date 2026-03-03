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
		0,
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
		0,
	)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(20000), congestionFee)
}

func TestCalculateCongestionFeeWithAdditionalMessages(t *testing.T) {
	calculator := setupCalculator()
	db, _ := testutils.NewRawDB(t, context.Background())

	ctx := context.Background()
	querier := queries.New(db)
	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)

	// Add congestion just at threshold (targetRatePerMinute=4)
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 4)

	// With 0 additional messages, fee should be 0 (at target, not above)
	fee0, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 0)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(0), fee0)

	// With 1 additional message pushing above target, fee should be > 0
	// currentMinute=5, ratio=5/4=1.25, in exponential range
	fee1, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 1)
	require.NoError(t, err)
	require.Greater(t, fee1, currency.PicoDollar(0))

	// With 2 additional messages, ratio=6/4=1.5, hitting max congestion
	// Fee should be higher than with 1 additional message
	fee2, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 2)
	require.NoError(t, err)
	require.Greater(t, fee2, fee1)
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
			ctx, seqQuerier, messageTime, seqOriginatorID, 0,
		)
		require.NoError(t, err)
		sequentialFees[i] = fee

		// Increment congestion in DB to simulate the message being committed.
		addCongestion(t, seqQuerier, seqOriginatorID, minutesSinceEpoch, 1)
	}

	// --- Batched path: same seed, compute 10 fees using additionalMessages=0..9 ---
	batchDB, _ := testutils.NewRawDB(t, ctx)
	batchQuerier := queries.New(batchDB)
	batchOriginatorID := uint32(testutils.RandomInt32())

	addCongestion(t, batchQuerier, batchOriginatorID, minutesSinceEpoch, seedCongestion)

	batchedFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := calculator.CalculateCongestionFee(
			ctx, batchQuerier, messageTime, batchOriginatorID, int32(i),
		)
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

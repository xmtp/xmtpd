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

// TestCongestionFeeParity_StaleReadReplicaCausesDivergence demonstrates the bug described
// in https://github.com/xmtp/xmtpd/issues/1818.
//
// When the remote node (EnvelopeSink) reads congestion from a stale read replica, it sees
// lower (or zero) message counts and therefore computes lower congestion fees than the
// originating node's BatchFeeCalculator.  The fix is to read congestion from the primary
// (write) DB so both paths always see the same committed state.
//
// The test uses two separate DB instances to simulate the primary and a frozen replica:
//   - writeDB: receives all IncrementOriginatorCongestion writes (primary)
//   - staleDB: seeded once and never updated (frozen replica)
//
// With targetRatePerMinute=4 and 3 messages seeded, messages 3-5 cross the congestion
// threshold and accrue non-zero fees on the batch/primary path but zero fees on the stale path.
func TestCongestionFeeParity_StaleReadReplicaCausesDivergence(t *testing.T) {
	calculator := setupCalculator()
	ctx := context.Background()

	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)
	seedCongestion := 3
	numMessages := 5

	// writeDB acts as the primary. All congestion increments land here.
	writeDB, _ := testutils.NewRawDB(t, ctx)
	writeQuerier := queries.New(writeDB)
	originatorID := uint32(testutils.RandomInt32())
	addCongestion(t, writeQuerier, originatorID, minutesSinceEpoch, seedCongestion)

	// staleDB simulates a lagging read replica frozen at 0 messages (no seed, no updates).
	staleDB, _ := testutils.NewRawDB(t, ctx)
	staleQuerier := queries.New(staleDB)

	// --- Local node: BatchFeeCalculator reads from writeDB ---
	batchCalc := calculator.NewBatchFeeCalculator(ctx, writeQuerier, originatorID)
	batchTotal := currency.PicoDollar(0)
	for range numMessages {
		fee, err := batchCalc.CalculateCongestionFee(messageTime)
		require.NoError(t, err)
		batchTotal += fee
	}

	// --- Remote node (BUG): CalculateCongestionFee reads from stale replica ---
	staleOriginatorID := uint32(testutils.RandomInt32())
	staleTotal := currency.PicoDollar(0)
	for range numMessages {
		fee, err := calculator.CalculateCongestionFee(
			ctx, staleQuerier, messageTime, staleOriginatorID,
		)
		require.NoError(t, err)
		staleTotal += fee
		// Increment only goes to writeDB; staleDB never sees it.
		addCongestion(t, writeQuerier, staleOriginatorID, minutesSinceEpoch, 1)
	}

	// Stale path under-counts: all fees are zero because the replica has no congestion data.
	require.Equal(
		t, currency.PicoDollar(0), staleTotal,
		"stale read replica should always return zero congestion (frozen at 0)",
	)
	// The batch path sees congestion and charges non-zero fees: proves the divergence.
	require.Greater(t, int64(batchTotal), int64(staleTotal),
		"local batch total (%d) should exceed stale remote total (%d): divergence confirmed",
		batchTotal, staleTotal,
	)

	// --- Remote node (FIX): CalculateCongestionFee reads from writeDB ---
	// Re-seed a fresh originator on writeDB so both paths start from the same state.
	fixedOriginatorID := uint32(testutils.RandomInt32())
	addCongestion(t, writeQuerier, fixedOriginatorID, minutesSinceEpoch, seedCongestion)

	batchCalcFixed := calculator.NewBatchFeeCalculator(ctx, writeQuerier, fixedOriginatorID)
	batchFixedTotal := currency.PicoDollar(0)
	for range numMessages {
		fee, err := batchCalcFixed.CalculateCongestionFee(messageTime)
		require.NoError(t, err)
		batchFixedTotal += fee
	}

	seqOriginatorID := uint32(testutils.RandomInt32())
	addCongestion(t, writeQuerier, seqOriginatorID, minutesSinceEpoch, seedCongestion)

	fixedTotal := currency.PicoDollar(0)
	for range numMessages {
		fee, err := calculator.CalculateCongestionFee(
			ctx, writeQuerier, messageTime, seqOriginatorID,
		)
		require.NoError(t, err)
		fixedTotal += fee
		// Increment goes to the same writeDB that the fee read came from: parity is preserved.
		addCongestion(t, writeQuerier, seqOriginatorID, minutesSinceEpoch, 1)
	}

	require.Equal(
		t, batchFixedTotal, fixedTotal,
		"remote (write DB) total should match local batch total: fix confirmed",
	)
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

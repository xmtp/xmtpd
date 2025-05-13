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
	RATE_MESSAGE_FEE       = 100
	RATE_STORAGE_FEE       = 50
	RATE_CONGESTION_FEE    = 200
	TARGET_RATE_PER_MINUTE = 4
)

func setupCalculator() *FeeCalculator {
	rates := &Rates{
		MessageFee:          RATE_MESSAGE_FEE,
		StorageFee:          RATE_STORAGE_FEE,
		CongestionFee:       RATE_CONGESTION_FEE,
		TargetRatePerMinute: TARGET_RATE_PER_MINUTE,
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

	expectedFee := RATE_MESSAGE_FEE + (RATE_STORAGE_FEE * messageSize * storageDurationDays)
	require.Equal(t, currency.PicoDollar(expectedFee), baseFee)
}

func TestCalculateCongestionFee(t *testing.T) {
	calculator := setupCalculator()
	db, _ := testutils.NewDB(t, context.Background())

	ctx := context.Background()
	querier := queries.New(db)
	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)

	// Should return 0 if no congestion
	congestionFee, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(0), congestionFee)

	// Congestion rate is 100 because this is double the max
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 8)
	congestionFee, err = calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(20000), congestionFee)
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

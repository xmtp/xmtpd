package fees

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
)

const (
	RATE_MESSAGE_FEE    = 100
	RATE_STORAGE_FEE    = 50
	RATE_CONGESTION_FEE = 200
)

func setupCalculator() *FeeCalculator {
	rates := &Rates{
		MessageFee:    RATE_MESSAGE_FEE,
		StorageFee:    RATE_STORAGE_FEE,
		CongestionFee: RATE_CONGESTION_FEE,
	}
	ratesFetcher := NewFixedRatesFetcher(rates)
	return NewFeeCalculator(ratesFetcher)
}

func TestCalculateBaseFee(t *testing.T) {
	calculator := setupCalculator()

	messageTime := time.Now()
	messageSize := int64(100)
	storageDurationDays := int64(1)

	baseFee, err := calculator.CalculateBaseFee(messageTime, messageSize, storageDurationDays)
	require.NoError(t, err)

	expectedFee := RATE_MESSAGE_FEE + (RATE_STORAGE_FEE * messageSize * storageDurationDays)
	require.Equal(t, currency.PicoDollar(expectedFee), baseFee)
}

func TestCalculateCongestionFee(t *testing.T) {
	calculator := setupCalculator()

	messageTime := time.Now()
	congestionPercent := int64(50)

	congestionFee, err := calculator.CalculateCongestionFee(messageTime, congestionPercent)
	require.NoError(t, err)

	expectedFee := RATE_CONGESTION_FEE * congestionPercent
	require.Equal(t, currency.PicoDollar(expectedFee), congestionFee)
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
		int64(storageDurationDays),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "storage fee calculation overflow")
}

func TestInvalidCongestionPercent(t *testing.T) {
	calculator := setupCalculator()

	messageTime := time.Now()

	_, err := calculator.CalculateCongestionFee(messageTime, int64(101))
	require.Error(t, err)
	require.Contains(t, err.Error(), "congestionPercent must be between 0 and 100")

	_, err = calculator.CalculateCongestionFee(messageTime, int64(-1))
	require.Error(t, err)
	require.Contains(t, err.Error(), "congestionPercent must be between 0 and 100")
}

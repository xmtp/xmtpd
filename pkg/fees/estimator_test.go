package fees_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/fees"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/fees"
)

func TestNewFeeEstimator(t *testing.T) {
	mockCalculator := mocks.NewMockIFeeCalculator(t)
	estimator := fees.NewFeeEstimator(mockCalculator)

	assert.NotNil(t, estimator)
	assert.NotNil(t, estimator)
}

func TestEstimateFees(t *testing.T) {
	t.Run("returns base fee when no congestion cached", func(t *testing.T) {
		mockCalculator := mocks.NewMockIFeeCalculator(t)
		estimator := fees.NewFeeEstimator(mockCalculator)

		baseFee := currency.PicoDollar(100)
		mockCalculator.EXPECT().
			CalculateBaseFee(mock.AnythingOfType("time.Time"), int64(1000), uint32(7)).
			Return(baseFee, nil)

		fee, err := estimator.EstimateFees(1, 1000, 7)
		require.NoError(t, err)
		assert.Equal(t, baseFee, fee)
	})

	t.Run("adds cached congestion to base fee", func(t *testing.T) {
		mockCalculator := mocks.NewMockIFeeCalculator(t)
		estimator := fees.NewFeeEstimator(mockCalculator)

		// Pre-populate cache with congestion fee
		ctx := context.Background()
		congestionFee := currency.PicoDollar(50)
		mockCalculator.EXPECT().
			CalculateCongestionFee(ctx, mock.Anything, mock.AnythingOfType("time.Time"), uint32(2)).
			Return(congestionFee, nil)

		_, _ = estimator.CalculateCongestionFee(ctx, &queries.Queries{}, time.Now(), 2)

		baseFee := currency.PicoDollar(200)
		mockCalculator.EXPECT().
			CalculateBaseFee(mock.AnythingOfType("time.Time"), int64(2000), uint32(14)).
			Return(baseFee, nil)

		fee, err := estimator.EstimateFees(2, 2000, 14)
		require.NoError(t, err)
		assert.Equal(t, baseFee+congestionFee, fee)
	})

	t.Run("propagates base fee errors", func(t *testing.T) {
		mockCalculator := mocks.NewMockIFeeCalculator(t)
		estimator := fees.NewFeeEstimator(mockCalculator)

		expectedErr := errors.New("base fee error")
		mockCalculator.EXPECT().
			CalculateBaseFee(mock.AnythingOfType("time.Time"), int64(3000), uint32(21)).
			Return(currency.PicoDollar(0), expectedErr)

		_, err := estimator.EstimateFees(3, 3000, 21)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestCalculateCongestionFeeUpdatesCache(t *testing.T) {
	mockCalculator := mocks.NewMockIFeeCalculator(t)
	estimator := fees.NewFeeEstimator(mockCalculator)

	ctx := context.Background()
	querier := &queries.Queries{}
	originatorID := uint32(1)
	messageTime := time.Now()
	congestionFee := currency.PicoDollar(25)

	mockCalculator.EXPECT().
		CalculateCongestionFee(ctx, querier, messageTime, originatorID).
		Return(congestionFee, nil)

	fee, err := estimator.CalculateCongestionFee(ctx, querier, messageTime, originatorID)
	require.NoError(t, err)
	assert.Equal(t, congestionFee, fee)

	// Verify cache was updated by calling EstimateFees
	baseFee := currency.PicoDollar(100)
	mockCalculator.EXPECT().
		CalculateBaseFee(mock.AnythingOfType("time.Time"), int64(1000), uint32(7)).
		Return(baseFee, nil)

	totalFee, err := estimator.EstimateFees(originatorID, 1000, 7)
	require.NoError(t, err)
	assert.Equal(t, baseFee+congestionFee, totalFee)
}

func TestConcurrentAccess(t *testing.T) {
	mockCalculator := mocks.NewMockIFeeCalculator(t)
	estimator := fees.NewFeeEstimator(mockCalculator)

	mockCalculator.EXPECT().
		CalculateBaseFee(mock.Anything, mock.Anything, mock.Anything).
		Return(currency.PicoDollar(100), nil).
		Maybe()

	mockCalculator.EXPECT().
		CalculateCongestionFee(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(currency.PicoDollar(50), nil).
		Maybe()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			originatorID := uint32(id % 10)

			switch id % 3 {
			case 0:
				_, _ = estimator.EstimateFees(originatorID, 1000, 7)
			case 1:
				_, _ = estimator.CalculateCongestionFee(
					context.Background(),
					&queries.Queries{},
					time.Now(),
					originatorID,
				)
			default:
				_, _ = estimator.CalculateBaseFee(time.Now(), 1000, 7)
			}
		}(i)
	}

	wg.Wait()
}

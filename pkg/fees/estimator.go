package fees

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type FeeEstimator struct {
	calculator       IFeeCalculator
	recentCongestion map[uint32]currency.PicoDollar
	mutex            sync.RWMutex
}

func NewFeeEstimator(calculator IFeeCalculator) *FeeEstimator {
	return &FeeEstimator{
		calculator:       calculator,
		recentCongestion: make(map[uint32]currency.PicoDollar),
		mutex:            sync.RWMutex{},
	}
}

func (e *FeeEstimator) EstimateFees(
	originatorID uint32,
	payerEnvelopeLength int64,
	retentionDays uint32,
) (currency.PicoDollar, error) {
	baseFee, err := e.CalculateBaseFee(time.Now(), payerEnvelopeLength, retentionDays)
	if err != nil {
		return 0, err
	}

	e.mutex.RLock()
	defer e.mutex.RUnlock()
	congestion := e.recentCongestion[originatorID]

	return baseFee + congestion, nil
}

func (e *FeeEstimator) updateCongestionEstimates(
	originatorID uint32,
	congestion currency.PicoDollar,
) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.recentCongestion[originatorID] = congestion
}

func (e *FeeEstimator) CalculateBaseFee(
	messageTime time.Time,
	messageSize int64,
	storageDurationDays uint32,
) (currency.PicoDollar, error) {
	return e.calculator.CalculateBaseFee(messageTime, messageSize, storageDurationDays)
}

func (e *FeeEstimator) CalculateCongestionFee(
	ctx context.Context,
	querier *queries.Queries,
	messageTime time.Time,
	originatorID uint32,
) (currency.PicoDollar, error) {
	congestion, err := e.calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID)
	if err != nil {
		return 0, err
	}

	e.updateCongestionEstimates(originatorID, congestion)

	return congestion, nil
}

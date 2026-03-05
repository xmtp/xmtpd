package fees

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// BatchFeeCalculator is a stateful calculator for batch congestion fee processing.
// It lazily caches per-minute congestion snapshots from the DB and tracks
// in-memory per-minute message counts across successive calls.
type BatchFeeCalculator struct {
	calculator   *FeeCalculator
	ctx          context.Context
	querier      *queries.Queries
	originatorID int32
	// DB snapshot cache: minute -> [5]int32 congestion window
	snapshots map[int32][5]int32
	// Batch-processed message counts per minute
	batchCounts map[int32]int32
}

func (c *FeeCalculator) NewBatchFeeCalculator(
	ctx context.Context,
	querier *queries.Queries,
	originatorID uint32,
) *BatchFeeCalculator {
	return &BatchFeeCalculator{
		calculator:   c,
		ctx:          ctx,
		querier:      querier,
		originatorID: int32(originatorID),
		snapshots:    make(map[int32][5]int32),
		batchCounts:  make(map[int32]int32),
	}
}

// CalculateCongestionFee computes the congestion fee for a message at the given time.
// It lazily fetches the DB snapshot for this minute (cached after first fetch),
// adjusts the 5-minute window for batch-processed messages, computes the fee,
// then increments the batch count for this minute.
func (b *BatchFeeCalculator) CalculateCongestionFee(
	messageTime time.Time,
) (currency.PicoDollar, error) {
	minute := int32(utils.MinutesSinceEpoch(messageTime))

	// Lazily fetch and cache the DB snapshot for this minute
	baseSnapshot, ok := b.snapshots[minute]
	if !ok {
		var err error
		baseSnapshot, err = db.Get5MinutesOfCongestion(
			b.ctx, b.querier, b.originatorID, minute,
		)
		if err != nil {
			return 0, err
		}
		b.snapshots[minute] = baseSnapshot
	}

	// Copy the snapshot (safe: [5]int32 is a value type, not a reference)
	// and adjust all 5 indices for batch-processed messages in the window
	adjusted := baseSnapshot
	for i := range int32(5) {
		adjusted[i] += b.batchCounts[minute-i]
	}

	rates, err := b.calculator.ratesFetcher.GetRates(messageTime)
	if err != nil {
		return 0, err
	}

	congestionUnits := CalculateCongestion(adjusted, int32(rates.TargetRatePerMinute))

	if congestionUnits < 0 || congestionUnits > 100 {
		return 0, fmt.Errorf(
			"congestionUnits must be between 0 and 100, got %d",
			congestionUnits,
		)
	}

	// Increment count AFTER computing fee (this message hasn't been "committed" yet)
	b.batchCounts[minute]++

	if congestionUnits == 0 {
		return 0, nil
	}

	result := rates.CongestionFee * currency.PicoDollar(congestionUnits)
	if result/currency.PicoDollar(congestionUnits) != rates.CongestionFee {
		return 0, errors.New("congestion fee calculation overflow")
	}
	return result, nil
}

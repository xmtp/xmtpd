package fees

import (
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/currency"
)

type FeeCalculator struct {
	ratesFetcher IRatesFetcher
}

func NewFeeCalculator(ratesFetcher IRatesFetcher) *FeeCalculator {
	return &FeeCalculator{ratesFetcher: ratesFetcher}
}

func (c *FeeCalculator) CalculateBaseFee(
	messageTime time.Time,
	messageSize int64,
	storageDurationDays int64,
) (currency.PicoDollar, error) {
	if messageSize <= 0 {
		return 0, fmt.Errorf("messageSize must be greater than 0, got %d", messageSize)
	}
	if storageDurationDays <= 0 {
		return 0, fmt.Errorf(
			"storageDurationDays must be greater than 0, got %d",
			storageDurationDays,
		)
	}

	rates, err := c.ratesFetcher.GetRates(messageTime)
	if err != nil {
		return 0, err
	}

	// Calculate storage fee components separately to check for overflow
	storageFeePerByte := rates.StorageFee * currency.PicoDollar(messageSize)
	if storageFeePerByte/currency.PicoDollar(messageSize) != rates.StorageFee {
		return 0, fmt.Errorf("storage fee calculation overflow")
	}

	totalStorageFee := storageFeePerByte * currency.PicoDollar(storageDurationDays)
	if totalStorageFee/currency.PicoDollar(storageDurationDays) != storageFeePerByte {
		return 0, fmt.Errorf("storage fee calculation overflow")
	}

	return rates.MessageFee + totalStorageFee, nil
}

func (c *FeeCalculator) CalculateCongestionFee(
	messageTime time.Time,
	congestionUnits int64,
) (currency.PicoDollar, error) {
	if congestionUnits < 0 || congestionUnits > 100 {
		return 0, fmt.Errorf(
			"congestionPercent must be between 0 and 100, got %d",
			congestionUnits,
		)
	}

	if congestionUnits == 0 {
		return 0, nil
	}

	rates, err := c.ratesFetcher.GetRates(messageTime)
	if err != nil {
		return 0, err
	}

	result := rates.CongestionFee * currency.PicoDollar(congestionUnits)
	if result/currency.PicoDollar(congestionUnits) != rates.CongestionFee {
		return 0, fmt.Errorf("congestion fee calculation overflow")
	}
	return result, nil
}

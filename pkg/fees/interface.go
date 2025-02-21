package fees

import (
	"time"

	"github.com/xmtp/xmtpd/pkg/currency"
)

// Rates containt the cost for each fee component at a given message time.
// Values in the rates struct are denominated in USD PicoDollars
type Rates struct {
	MessageFee    currency.PicoDollar // The flat per-message fee
	StorageFee    currency.PicoDollar // The fee per byte-day of storage
	CongestionFee currency.PicoDollar // The fee per unit of congestion
}

// The RatesFetcher is responsible for loading the rates for a given message time.
// This allows us to roll out new rates over time, and apply them to messages consistently.
type IRatesFetcher interface {
	GetRates(messageTime time.Time) (*Rates, error)
}

type IFeeCalculator interface {
	CalculateBaseFee(
		messageTime time.Time,
		messageSize int64,
		storageDurationDays int64,
	) (currency.PicoDollar, error)
	CalculateCongestionFee(
		messageTime time.Time,
		congestionUnits int64,
	) (currency.PicoDollar, error)
}

// Package fees implements the fees test utils.
package fees

import (
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
)

var testRates = &fees.Rates{
	MessageFee:          currency.PicoDollar(100),
	StorageFee:          currency.PicoDollar(100),
	CongestionFee:       currency.PicoDollar(100),
	TargetRatePerMinute: 100,
}

func NewTestRatesFetcher() *fees.FixedRatesFetcher {
	return fees.NewFixedRatesFetcher(testRates)
}

func NewTestFeeCalculator() *fees.FeeCalculator {
	return fees.NewFeeCalculator(NewTestRatesFetcher())
}

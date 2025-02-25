package fees

import (
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
)

var TEST_RATES = &fees.Rates{
	MessageFee:    currency.PicoDollar(100),
	StorageFee:    currency.PicoDollar(100),
	CongestionFee: currency.PicoDollar(100),
}

func NewTestRatesFetcher() *fees.FixedRatesFetcher {
	return fees.NewFixedRatesFetcher(TEST_RATES)
}

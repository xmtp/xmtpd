package fees

import "time"

// A fixed fee schedule that doesn't rely on the blockchain.
// Used primarily for testing, and until we have fees onchain.
type FixedRatesFetcher struct {
	rates *Rates
}

func NewFixedRatesFetcher(rates *Rates) *FixedRatesFetcher {
	return &FixedRatesFetcher{rates: rates}
}

func (f *FixedRatesFetcher) GetRates(_messageTime time.Time) (*Rates, error) {
	return f.rates, nil
}

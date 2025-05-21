package fees

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const MAX_REFRESH_INTERVAL = 60 * time.Minute

// Dumbed down version of the RatesManager contract interface
type RatesContract interface {
	GetRates(
		opts *bind.CallOpts,
		fromIndex *big.Int,
		count *big.Int,
	) ([]rateregistry.IRateRegistryRates, error)
}

type indexedRates struct {
	startTime time.Time
	rates     *Rates
}

// ContractRatesFetcher pulls all the rates from the RatesManager contract
// and stores them in a sorted set to find the appropriate rate for a given timestamp
type ContractRatesFetcher struct {
	ctx             context.Context
	wg              sync.WaitGroup
	logger          *zap.Logger
	contract        RatesContract
	rates           []*indexedRates
	refreshInterval time.Duration
	lastRefresh     time.Time
}

// NewContractRatesFetcher creates a new ContractRatesFetcher using the provided eth client
func NewContractRatesFetcher(
	ctx context.Context,
	ethclient bind.ContractCaller,
	logger *zap.Logger,
	options config.ContractsOptions,
) (*ContractRatesFetcher, error) {
	contract, err := rateregistry.NewRateRegistryCaller(
		common.HexToAddress(options.SettlementChain.RateRegistryAddress),
		ethclient,
	)
	if err != nil {
		return nil, err
	}

	return &ContractRatesFetcher{
		logger:          logger.Named("contractRatesFetcher"),
		contract:        contract,
		ctx:             ctx,
		refreshInterval: options.SettlementChain.RateRegistryRefreshInterval,
	}, nil
}

// Start the ContractRatesFetcher and begin fetching rates from the smart contract
// periodically.
func (c *ContractRatesFetcher) Start() error {
	// If we can't load the data at least once, fail to start the service
	if err := c.refreshData(); err != nil {
		c.logger.Error("Failed to refresh data", zap.Error(err))
		return err
	}

	tracing.GoPanicWrap(
		c.ctx,
		&c.wg,
		"rates-fetcher",
		func(ctx context.Context) { c.refreshLoop() },
	)

	return nil
}

// refreshData fetches all rates from the smart contract and validates them
func (c *ContractRatesFetcher) refreshData() error {
	var err error

	fromIndex := big.NewInt(0)
	newRates := make([]*indexedRates, 0)
	// for {
	c.logger.Info("getting page", zap.Int64("fromIndex", fromIndex.Int64()))
	resp, err := c.contract.GetRates(&bind.CallOpts{Context: c.ctx}, fromIndex, big.NewInt(1))
	if err != nil {
		c.logger.Error(
			"error calling contract",
			zap.Error(err),
			zap.Int64("fromIndex", fromIndex.Int64()),
		)
		return err
	}

	newRates = append(newRates, transformRates(resp)...)
	// fromIndex = fromIndex.Add(fromIndex, big.NewInt(int64(len(resp))))

	c.logger.Info("getting next page")

	// TODO mkysel fix paging
	//	break

	//}

	if err = validateRates(newRates); err != nil {
		c.logger.Error("failed to validate rates", zap.Error(err))
		return err
	}

	c.rates = newRates
	c.lastRefresh = time.Now()
	c.logger.Debug("refreshed rates", zap.Int("numRates", len(newRates)))

	return err
}

func (c *ContractRatesFetcher) GetRates(timestamp time.Time) (*Rates, error) {
	if time.Since(c.lastRefresh) > MAX_REFRESH_INTERVAL {
		c.logger.Warn(
			"last rates refresh was too long ago for accurate rates",
			zap.Duration("duration", time.Since(c.lastRefresh)),
		)
		return nil, errors.New("last rates refresh was too long ago")
	}

	if len(c.rates) == 0 {
		return nil, errors.New("no rates found")
	}

	// If the timestamp is before the oldest rate, return an error
	if timestamp.Before(c.rates[0].startTime) {
		return nil, errors.New("timestamp is before the oldest rate")
	}

	// Most messages should using the current rate, so check that before doing a binary search
	newestRate := c.rates[len(c.rates)-1]
	if timestamp.After(newestRate.startTime) {
		return newestRate.rates, nil
	}

	return c.findMatchingRate(timestamp), nil
}

func (c *ContractRatesFetcher) findMatchingRate(timestamp time.Time) *Rates {
	// Binary search to find the rate with the closest startTime that is before or equal to the provided timestamp
	left, right := 0, len(c.rates)-1

	for left <= right {
		mid := left + (right-left)/2

		if c.rates[mid].startTime.Equal(timestamp) {
			return c.rates[mid].rates
		}

		if c.rates[mid].startTime.Before(timestamp) {
			// Check if this is the closest rate before the timestamp
			if mid == len(c.rates)-1 || c.rates[mid+1].startTime.After(timestamp) {
				return c.rates[mid].rates
			}

			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// Fallback to the first rate if no exact or closest match is found
	return c.rates[0].rates
}

func (c *ContractRatesFetcher) refreshLoop() {
	ticker := time.NewTicker(c.refreshInterval)
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.refreshData(); err != nil {
				c.logger.Error("Failed to refresh data", zap.Error(err))
			}
		}
	}
}

func transformRates(rates []rateregistry.IRateRegistryRates) []*indexedRates {
	newIndexedRates := make([]*indexedRates, len(rates))
	for i, rate := range rates {
		newIndexedRates[i] = &indexedRates{
			startTime: time.Unix(int64(rate.StartTime), 0),
			rates: &Rates{
				MessageFee:          currency.PicoDollar(rate.MessageFee),
				StorageFee:          currency.PicoDollar(rate.StorageFee),
				CongestionFee:       currency.PicoDollar(rate.CongestionFee),
				TargetRatePerMinute: rate.TargetRatePerMinute,
			},
		}
	}

	return newIndexedRates
}

func validateRates(rates []*indexedRates) error {
	if len(rates) == 0 {
		return errors.New("no rates found")
	}
	earliestStart := rates[0].startTime
	for _, rate := range rates[1:] {
		if rate.startTime.Equal(earliestStart) {
			return errors.New("duplicate rate start time")
		}

		if rate.startTime.Before(earliestStart) {
			return errors.New("rates are not sorted")
		}
		earliestStart = rate.startTime
	}
	return nil
}

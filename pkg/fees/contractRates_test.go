package fees

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/ratesmanager"
	"github.com/xmtp/xmtpd/pkg/currency"
	feesMock "github.com/xmtp/xmtpd/pkg/mocks/fees"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildFetcher(t *testing.T) (*ContractRatesFetcher, *feesMock.MockRatesContract, func()) {
	mockContract := feesMock.NewMockRatesContract(t)
	ctx, cancel := context.WithCancel(context.Background())

	fetcher := &ContractRatesFetcher{
		logger:          testutils.NewLog(t),
		ctx:             ctx,
		contract:        mockContract,
		refreshInterval: 100 * time.Millisecond,
	}

	return fetcher, mockContract, func() {
		cancel()
	}
}

func buildRates(fees uint64, startTime uint64) ratesmanager.RatesManagerRates {
	return ratesmanager.RatesManagerRates{
		MessageFee:    fees,
		StorageFee:    fees,
		CongestionFee: fees,
		StartTime:     startTime,
	}
}

func TestLoadGetRates(t *testing.T) {
	fetcher, mockContract, cancel := buildFetcher(t)
	defer cancel()

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0)).Return(ratesResponse{
		Rates:   []ratesmanager.RatesManagerRates{buildRates(100, 1), buildRates(200, 2)},
		HasMore: false,
	}, nil)

	require.NoError(t, fetcher.Start())

	require.Len(t, fetcher.rates, 2)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
}

func TestCanPaginate(t *testing.T) {
	fetcher, mockContract, cancel := buildFetcher(t)
	defer cancel()

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0)).Return(ratesResponse{
		Rates:   []ratesmanager.RatesManagerRates{buildRates(100, 1), buildRates(200, 2)},
		HasMore: true,
	}, nil).Times(1)

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(2)).Return(ratesResponse{
		Rates:   []ratesmanager.RatesManagerRates{buildRates(300, 3)},
		HasMore: false,
	}, nil).Times(1)

	require.NoError(t, fetcher.Start())

	require.Len(t, fetcher.rates, 3)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
	require.Equal(t, fetcher.rates[2].rates.MessageFee, currency.PicoDollar(300))
}

func TestGetRates(t *testing.T) {
	fetcher, mockContract, cancel := buildFetcher(t)
	defer cancel()

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0)).Return(ratesResponse{
		Rates: []ratesmanager.RatesManagerRates{
			buildRates(100, 100),
			buildRates(200, 200),
			buildRates(300, 300),
		},
		HasMore: false,
	}, nil)

	require.NoError(t, fetcher.Start())

	// Exactly equals the first rate
	rates, err := fetcher.GetRates(time.Unix(100, 0))
	require.NoError(t, err)
	require.Equal(t, rates.MessageFee, currency.PicoDollar(100))

	// Between the first and second rate
	rates, err = fetcher.GetRates(time.Unix(101, 0))
	require.NoError(t, err)
	require.Equal(t, rates.MessageFee, currency.PicoDollar(100))

	// After the second rate
	rates, err = fetcher.GetRates(time.Unix(202, 0))
	require.NoError(t, err)
	require.Equal(t, rates.MessageFee, currency.PicoDollar(200))

	// After the third rate
	rates, err = fetcher.GetRates(time.Unix(303, 0))
	require.NoError(t, err)
	require.Equal(t, rates.MessageFee, currency.PicoDollar(300))
}

func TestFailIfNoRates(t *testing.T) {
	fetcher, mockContract, cancel := buildFetcher(t)
	defer cancel()

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0)).Return(ratesResponse{
		Rates:   []ratesmanager.RatesManagerRates{},
		HasMore: false,
	}, nil)

	require.Error(t, fetcher.Start())
}

func TestGetRatesBeforeFirstRate(t *testing.T) {
	fetcher, mockContract, cancel := buildFetcher(t)
	defer cancel()

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0)).Return(ratesResponse{
		Rates: []ratesmanager.RatesManagerRates{
			buildRates(100, 100),
			buildRates(200, 200),
			buildRates(300, 300),
		},
		HasMore: false,
	}, nil)

	require.NoError(t, fetcher.Start())

	rates, err := fetcher.GetRates(time.Unix(50, 0))
	require.Error(t, errors.New("timestamp is before the oldest rate"), err)
	require.Nil(t, rates)
}

func TestGetRatesUninitialized(t *testing.T) {
	fetcher, _, cancel := buildFetcher(t)
	defer cancel()

	rates, err := fetcher.GetRates(time.Unix(100, 0))
	require.Error(t, errors.New("no rates found"), err)
	require.Nil(t, rates)
}

package fees

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/currency"
	feesMock "github.com/xmtp/xmtpd/pkg/mocks/fees"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildFetcher(t *testing.T) (*ContractRatesFetcher, *feesMock.MockRatesContract) {
	mockContract := feesMock.NewMockRatesContract(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	fetcher := &ContractRatesFetcher{
		logger:          testutils.NewLog(t),
		ctx:             ctx,
		contract:        mockContract,
		refreshInterval: 100 * time.Millisecond,
	}

	return fetcher, mockContract
}

func buildRates(fees uint64, startTime uint64) rateregistry.IRateRegistryRates {
	return rateregistry.IRateRegistryRates{
		MessageFee:          fees,
		StorageFee:          fees,
		CongestionFee:       fees,
		StartTime:           startTime,
		TargetRatePerMinute: 100 * 60,
	}
}

func TestLoadGetRates(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(100, 1), buildRates(200, 2)}, nil)

	require.NoError(t, fetcher.Start())

	require.Len(t, fetcher.rates, 2)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
	require.Equal(t, fetcher.rates[0].rates.TargetRatePerMinute, uint64(100*60))
	require.Equal(t, fetcher.rates[1].rates.TargetRatePerMinute, uint64(100*60))
}

func TestCanPaginate(t *testing.T) {
	t.Skip()
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(100, 1), buildRates(200, 2)}, nil).
		Times(1)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(2), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(300, 3)}, nil).
		Times(1)

	require.NoError(t, fetcher.Start())

	require.Len(t, fetcher.rates, 3)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
	require.Equal(t, fetcher.rates[2].rates.MessageFee, currency.PicoDollar(300))
}

func TestGetRates(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0), mock.Anything).Return([]rateregistry.IRateRegistryRates{
		buildRates(100, 100),
		buildRates(200, 200),
		buildRates(300, 300),
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
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0), mock.Anything).Return([]rateregistry.IRateRegistryRates{},
		nil)

	require.Error(t, fetcher.Start())
}

func TestGetRatesBeforeFirstRate(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().GetRates(mock.Anything, big.NewInt(0), mock.Anything).Return([]rateregistry.IRateRegistryRates{
		buildRates(100, 100),
		buildRates(200, 200),
		buildRates(300, 300),
	}, nil)

	require.NoError(t, fetcher.Start())

	rates, err := fetcher.GetRates(time.Unix(50, 0))
	require.ErrorContains(t, err, "timestamp is before the oldest rate")
	require.Nil(t, rates)
}

func TestGetRatesUninitialized(t *testing.T) {
	fetcher, _ := buildFetcher(t)

	rates, err := fetcher.GetRates(time.Unix(100, 0))
	require.ErrorContains(t, err, "last rates refresh was too long ago")
	require.Nil(t, rates)
}

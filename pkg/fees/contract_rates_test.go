package fees

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/testutils"
	feesMock "github.com/xmtp/xmtpd/pkg/testutils/mocks/fees"
)

const testPageSize = 5

func buildFetcher(t *testing.T) (*ContractRatesFetcher, *feesMock.MockRatesContract) {
	mockContract := feesMock.NewMockRatesContract(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	fetcher := &ContractRatesFetcher{
		logger:          testutils.NewLog(t),
		ctx:             ctx,
		contract:        mockContract,
		refreshInterval: 100 * time.Millisecond,
		pageSize:        testPageSize,
		currentIndex:    big.NewInt(0),
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
		GetRatesCount(mock.Anything).
		Return(big.NewInt(1), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(100, 1), buildRates(200, 2)}, nil)

	require.NoError(t, fetcher.refreshData())

	require.Len(t, fetcher.rates, 2)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
	require.Equal(t, fetcher.rates[0].rates.TargetRatePerMinute, uint64(100*60))
	require.Equal(t, fetcher.rates[1].rates.TargetRatePerMinute, uint64(100*60))
}

func TestCanPaginate(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(6), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(100, 1), buildRates(200, 2), buildRates(300, 3), buildRates(400, 4), buildRates(500, 5)}, nil).
		Times(1)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(5), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(600, 6)}, nil).
		Times(1)

	require.NoError(t, fetcher.refreshData())

	require.Len(t, fetcher.rates, 6)
	require.Equal(t, fetcher.rates[0].rates.MessageFee, currency.PicoDollar(100))
	require.Equal(t, fetcher.rates[1].rates.MessageFee, currency.PicoDollar(200))
	require.Equal(t, fetcher.rates[2].rates.MessageFee, currency.PicoDollar(300))
	require.Equal(t, fetcher.rates[3].rates.MessageFee, currency.PicoDollar(400))
	require.Equal(t, fetcher.rates[4].rates.MessageFee, currency.PicoDollar(500))
	require.Equal(t, fetcher.rates[5].rates.MessageFee, currency.PicoDollar(600))
}

func TestGetRates(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(3), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{
			buildRates(100, 100),
			buildRates(200, 200),
			buildRates(300, 300),
		}, nil)

	require.NoError(t, fetcher.refreshData())

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

	mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(1), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{},
			nil)

	require.Error(t, fetcher.refreshData())
}

func TestGetRatesBeforeFirstRate(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(3), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{
			buildRates(100, 100),
			buildRates(200, 200),
			buildRates(300, 300),
		}, nil)

	require.NoError(t, fetcher.refreshData())

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

func TestCanContinue(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	counts := mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(5), nil)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(100, 1), buildRates(200, 2), buildRates(300, 3), buildRates(400, 4), buildRates(500, 5)}, nil).
		Times(1)

	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(5), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{buildRates(600, 6)}, nil).
		Times(1)

	require.NoError(t, fetcher.refreshData())
	require.Len(t, fetcher.rates, 5)

	counts.Return(big.NewInt(6), nil)

	require.NoError(t, fetcher.refreshData())
	require.Len(t, fetcher.rates, 6)
}

// TestPartialRefreshFailureDoesNotSkipRates guards the incremental-refresh
// cursor against a mid-pagination failure.
//
// refreshData accumulates fetched pages into a local slice and only commits it
// to c.rates at the very end, but it advances c.currentIndex inside the loop as
// each page is fetched. If a later page fails, the already-fetched pages are
// discarded (the local slice is dropped) while the cursor stays advanced. The
// next successful refresh then resumes past the discarded rates, leaving a
// permanent gap in the rate set.
//
// Concretely: 10 rates, page size 5. First refresh fetches [0,5) then fails on
// [5,10). Recovery refresh must re-fetch [0,5) — not resume at [5,10) and drop
// the first five rates.
func TestPartialRefreshFailureDoesNotSkipRates(t *testing.T) {
	fetcher, mockContract := buildFetcher(t)

	mockContract.EXPECT().
		GetRatesCount(mock.Anything).
		Return(big.NewInt(10), nil)

	// Page [0,5) always succeeds.
	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(0), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{
			buildRates(100, 1), buildRates(200, 2), buildRates(300, 3),
			buildRates(400, 4), buildRates(500, 5),
		}, nil)

	// Page [5,10): fails the first time it is requested, succeeds the second.
	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(5), mock.Anything).
		Return(nil, errors.New("transient rpc error")).
		Once()
	mockContract.EXPECT().
		GetRates(mock.Anything, big.NewInt(5), mock.Anything).
		Return([]rateregistry.IRateRegistryRates{
			buildRates(600, 6), buildRates(700, 7), buildRates(800, 8),
			buildRates(900, 9), buildRates(1000, 10),
		}, nil).
		Once()

	// First refresh fails on the second page and must not commit anything.
	require.Error(t, fetcher.refreshData())
	require.Empty(t, fetcher.rates, "a failed refresh must not commit partial rates")

	// Recovery refresh: with a correct cursor this re-fetches [0,5) then [5,10)
	// for all 10 rates; a leaked cursor resumes at [5,10) and drops the first 5.
	require.NoError(t, fetcher.refreshData())

	require.Len(t, fetcher.rates, 10,
		"recovery refresh skipped the rates dropped by the failed one (cursor advanced without committing)")
	require.Equal(t, currency.PicoDollar(100), fetcher.rates[0].rates.MessageFee)
	require.Equal(t, currency.PicoDollar(1000), fetcher.rates[9].rates.MessageFee)
}

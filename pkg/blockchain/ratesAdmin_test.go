package blockchain_test

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildRatesAdmin(t *testing.T) *blockchain.RatesAdmin {
	ctx := context.Background()
	logger := testutils.NewLog(t)
	rpcUrl := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.AppChain.RpcURL)
	require.NoError(t, err)

	ratesAdmin, err := blockchain.NewRatesAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	return ratesAdmin
}

func TestAddRates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	rates := fees.Rates{
		MessageFee:          100,
		StorageFee:          200,
		CongestionFee:       300,
		TargetRatePerMinute: 100 * 60,
	}

	err := ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)
}

func TestAddNegativeRates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee: -100,
	})
	require.ErrorContains(t, err, "must be positive")

	err = ratesAdmin.AddRates(context.Background(), fees.Rates{
		StorageFee: -100,
	})
	require.ErrorContains(t, err, "must be positive")

	err = ratesAdmin.AddRates(context.Background(), fees.Rates{
		CongestionFee: -100,
	})
	require.ErrorContains(t, err, "must be positive")
}

func TestAdd0Rates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee:          0,
		StorageFee:          0,
		CongestionFee:       0,
		TargetRatePerMinute: 0,
	})
	require.NoError(t, err)
}

func TestAddLargeRates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee:          math.MaxInt64,
		StorageFee:          math.MaxInt64,
		CongestionFee:       math.MaxInt64,
		TargetRatePerMinute: math.MaxUint64,
	})
	require.NoError(t, err)
}

func TestAddRatesAgain(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	rates := fees.Rates{
		MessageFee:          5,
		StorageFee:          16,
		CongestionFee:       700,
		TargetRatePerMinute: 1000,
	}

	err := ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)

	err = ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)
}

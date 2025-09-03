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

func buildRatesAdmin(t *testing.T) (*blockchain.RatesAdmin, *blockchain.ParameterAdmin) {
	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	// TODO(mkysel) aren't rates on the settlement chain?
	client, err := blockchain.NewRPCClient(
		ctx,
		contractsOptions.AppChain.RPCURL,
	)
	require.NoError(t, err)

	paramAdmin, err := blockchain.NewParameterAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	ratesAdmin, err := blockchain.NewRatesAdmin(logger, paramAdmin, client, contractsOptions)
	require.NoError(t, err)

	return ratesAdmin, paramAdmin
}

func TestAddRates(t *testing.T) {
	ratesAdmin, _ := buildRatesAdmin(t)

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
	ratesAdmin, _ := buildRatesAdmin(t)

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
	ratesAdmin, _ := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee:          0,
		StorageFee:          0,
		CongestionFee:       0,
		TargetRatePerMinute: 0,
	})
	require.NoError(t, err)
}

func TestAddLargeRates(t *testing.T) {
	ratesAdmin, _ := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee:          math.MaxInt64,
		StorageFee:          math.MaxInt64,
		CongestionFee:       math.MaxInt64,
		TargetRatePerMinute: math.MaxUint64,
	})
	require.NoError(t, err)
}

func TestAddRatesAgain(t *testing.T) {
	ratesAdmin, _ := buildRatesAdmin(t)

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

func TestRates_ReadDefaults(t *testing.T) {
	_, paramAdmin := buildRatesAdmin(t)
	ctx := context.Background()

	// 1) Defaults should be zero (unset) for all four rate params.
	msg, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_MESSAGE_FEE_KEY)
	require.NoError(t, err)
	store, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_STORAGE_FEE_KEY)
	require.NoError(t, err)
	cong, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_CONGESTION_FEE_KEY)
	require.NoError(t, err)
	target, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY,
	)
	require.NoError(t, err)

	require.EqualValues(t, 0, msg, "default messageFee should be 0")
	require.EqualValues(t, 0, store, "default storageFee should be 0")
	require.EqualValues(t, 0, cong, "default congestionFee should be 0")
	require.EqualValues(t, 0, target, "default targetRatePerMinute should be 0")
}

func TestRates_WriteThenRead(t *testing.T) {
	ratesAdmin, paramAdmin := buildRatesAdmin(t)
	ctx := context.Background()

	want := fees.Rates{
		MessageFee:          123,
		StorageFee:          456,
		CongestionFee:       789,
		TargetRatePerMinute: 60,
	}
	require.NoError(t, ratesAdmin.AddRates(ctx, want))

	msg2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_MESSAGE_FEE_KEY)
	require.NoError(t, err)
	store2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_STORAGE_FEE_KEY)
	require.NoError(t, err)
	cong2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_CONGESTION_FEE_KEY)
	require.NoError(t, err)
	target2, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY,
	)
	require.NoError(t, err)

	require.EqualValues(t, uint64(want.MessageFee), msg2)
	require.EqualValues(t, uint64(want.StorageFee), store2)
	require.EqualValues(t, uint64(want.CongestionFee), cong2)
	require.EqualValues(t, want.TargetRatePerMinute, target2)
}

func TestRates_WriteZeroes_ReadZeroes(t *testing.T) {
	ratesAdmin, paramAdmin := buildRatesAdmin(t)
	ctx := context.Background()

	zero := fees.Rates{
		MessageFee:          0,
		StorageFee:          0,
		CongestionFee:       0,
		TargetRatePerMinute: 0,
	}
	require.NoError(t, ratesAdmin.AddRates(ctx, zero))

	msg, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_MESSAGE_FEE_KEY)
	require.NoError(t, err)
	store, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_STORAGE_FEE_KEY)
	require.NoError(t, err)
	cong, err := paramAdmin.GetParameterUint64(ctx, blockchain.RATE_REGISTRY_CONGESTION_FEE_KEY)
	require.NoError(t, err)
	target, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY,
	)
	require.NoError(t, err)

	require.EqualValues(t, 0, msg)
	require.EqualValues(t, 0, store)
	require.EqualValues(t, 0, cong)
	require.EqualValues(t, 0, target)
}

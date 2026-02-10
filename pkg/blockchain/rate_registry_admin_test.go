package blockchain_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildRatesAdmin(t *testing.T) (blockchain.IRatesAdmin, blockchain.IParameterAdmin) {
	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.TestPrivateKey,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(
		ctx,
		contractsOptions.SettlementChain.RPCURL,
	)
	require.NoError(t, err)

	paramAdmin, err := blockchain.NewSettlementParameterAdmin(
		logger,
		client,
		signer,
		contractsOptions,
	)
	require.NoError(t, err)

	ratesAdmin, err := blockchain.NewRatesAdmin(
		logger,
		client,
		signer,
		paramAdmin,
		contractsOptions,
	)
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
		StartTime:           uint64(time.Now().Add(2 * time.Hour).Unix()),
	}

	err := ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)
}

func TestAddNegativeRates(t *testing.T) {
	ratesAdmin, _ := buildRatesAdmin(t)

	err := ratesAdmin.AddRates(context.Background(), fees.Rates{
		MessageFee: -100,
		StartTime:  uint64(time.Now().Add(2 * time.Hour).Unix()),
	})
	require.ErrorContains(t, err, "must be positive")

	err = ratesAdmin.AddRates(context.Background(), fees.Rates{
		StorageFee: -100,
		StartTime:  uint64(time.Now().Add(2 * time.Hour).Unix()),
	})
	require.ErrorContains(t, err, "must be positive")

	err = ratesAdmin.AddRates(context.Background(), fees.Rates{
		CongestionFee: -100,
		StartTime:     uint64(time.Now().Add(2 * time.Hour).Unix()),
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
		StartTime:           uint64(time.Now().Add(2 * time.Hour).Unix()),
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
		StartTime:           uint64(time.Now().Add(2 * time.Hour).Unix()),
	})
	require.NoError(t, err)
}

func TestAddRatesAgain(t *testing.T) {
	ratesAdmin, _ := buildRatesAdmin(t)

	startTime := uint64(time.Now().Add(2 * time.Hour).Unix())

	rates := fees.Rates{
		MessageFee:          5,
		StorageFee:          16,
		CongestionFee:       700,
		TargetRatePerMinute: 1000,
		StartTime:           startTime,
	}

	err := ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)

	rates.StartTime = startTime + 3600
	err = ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)
}

func TestRates_ReadDefaults(t *testing.T) {
	t.Skip("Some defaults seem to be update - https://github.com/xmtp/smart-contracts/issues/126")
	_, paramAdmin := buildRatesAdmin(t)
	ctx := context.Background()

	// 1) Defaults should be zero (unset) for all four rate params.
	msg, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryMessageFeeKey)
	require.NoError(t, err)
	store, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryStorageFeeKey)
	require.NoError(t, err)
	cong, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryCongestionFeeKey)
	require.NoError(t, err)
	target, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RateRegistryTargetRatePerMinuteKey,
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
		StartTime:           uint64(time.Now().Add(2 * time.Hour).Unix()),
	}
	require.NoError(t, ratesAdmin.AddRates(ctx, want))

	msg2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryMessageFeeKey)
	require.NoError(t, err)
	store2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryStorageFeeKey)
	require.NoError(t, err)
	cong2, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryCongestionFeeKey)
	require.NoError(t, err)
	target2, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RateRegistryTargetRatePerMinuteKey,
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
		StartTime:           uint64(time.Now().Add(2 * time.Hour).Unix()),
	}
	require.NoError(t, ratesAdmin.AddRates(ctx, zero))

	msg, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryMessageFeeKey)
	require.NoError(t, err)
	store, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryStorageFeeKey)
	require.NoError(t, err)
	cong, err := paramAdmin.GetParameterUint64(ctx, blockchain.RateRegistryCongestionFeeKey)
	require.NoError(t, err)
	target, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.RateRegistryTargetRatePerMinuteKey,
	)
	require.NoError(t, err)

	require.EqualValues(t, 0, msg)
	require.EqualValues(t, 0, store)
	require.EqualValues(t, 0, cong)
	require.EqualValues(t, 0, target)
}

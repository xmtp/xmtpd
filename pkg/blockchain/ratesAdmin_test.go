package blockchain_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildRatesAdmin(t *testing.T) *blockchain.RatesAdmin {
	ctx := context.Background()
	logger := testutils.NewLog(t)
	rpcUrl, cleanup := anvil.StartAnvil(t, false)
	t.Cleanup(cleanup)
	contractsOptions := testutils.NewContractsOptions(rpcUrl)

	// Set the nodes contract address to a random smart contract instead of the fixed deployment
	contractsOptions.RatesManagerContractAddress = testutils.DeployRatesManagerContract(t, rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	ratesAdmin, err := blockchain.NewRatesAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	return ratesAdmin
}

func TestAddRates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	rates := rateregistry.RateRegistryRates{
		MessageFee:          100,
		StorageFee:          200,
		CongestionFee:       300,
		TargetRatePerMinute: 100 * 60,
		StartTime:           1000,
	}

	var err error

	err = ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)

	var returnedRates struct {
		Rates   []rateregistry.RateRegistryRates
		HasMore bool
	}

	returnedRates, err = ratesAdmin.Contract().GetRates(&bind.CallOpts{}, big.NewInt(0))

	require.NoError(t, err)
	require.Len(t, returnedRates.Rates, 1)
	require.False(t, returnedRates.HasMore)
	require.Equal(t, returnedRates.Rates[0], rates)
}

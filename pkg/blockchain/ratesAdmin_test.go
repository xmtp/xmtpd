package blockchain

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/ratesmanager"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildRatesAdmin(t *testing.T) *RatesAdmin {
	ctx := context.Background()
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)
	// Set the nodes contract address to a random smart contract instead of the fixed deployment
	contractsOptions.RatesManagerContractAddress = testutils.DeployRatesManagerContract(t)

	signer, err := NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	ratesAdmin, err := NewRatesAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	return ratesAdmin
}

func TestAddRates(t *testing.T) {
	ratesAdmin := buildRatesAdmin(t)

	rates := ratesmanager.RatesManagerRates{
		MessageFee:    100,
		StorageFee:    200,
		CongestionFee: 300,
		StartTime:     1000,
	}

	err := ratesAdmin.AddRates(context.Background(), rates)
	require.NoError(t, err)

	returnedRates, err := ratesAdmin.contract.GetRates(&bind.CallOpts{}, big.NewInt(0))
	require.NoError(t, err)
	require.Len(t, returnedRates.Rates, 1)
	require.False(t, returnedRates.HasMore)
	require.Equal(t, returnedRates.Rates[0], rates)
}

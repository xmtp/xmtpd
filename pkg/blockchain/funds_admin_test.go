package blockchain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildFundsAdmin(
	t *testing.T,
) (blockchain.IFundsAdmin, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(
		ctx,
		contractsOptions.SettlementChain.RPCURL,
	)
	require.NoError(t, err)

	admin, err := blockchain.NewFundsAdmin(
		logger,
		client,
		signer,
		contractsOptions,
	)
	require.NoError(t, err)

	return admin, ctx
}

func TestFundsAdmin(t *testing.T) {
	_, _ = buildFundsAdmin(t)
}

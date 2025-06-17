package blockchain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildReportsManager(t *testing.T) blockchain.PayerReportsManager {
	logger := testutils.NewLog(t)
	rpcUrl := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(t.Context(), contractsOptions.SettlementChain.WssURL)
	require.NoError(t, err)

	reportsManager, err := blockchain.NewReportsManager(
		logger,
		client,
		signer,
		contractsOptions.SettlementChain,
	)
	require.NoError(t, err)

	return reportsManager
}

func TestDomainSeparator(t *testing.T) {
	reportsManager := buildReportsManager(t)

	domainSeparator, err := reportsManager.GetDomainSeparator(t.Context())
	require.NoError(t, err)

	require.Len(t, domainSeparator, 32)

	domainSeparatorCached, err := reportsManager.GetDomainSeparator(t.Context())
	require.NoError(t, err)

	require.Equal(t, domainSeparator, domainSeparatorCached)
}

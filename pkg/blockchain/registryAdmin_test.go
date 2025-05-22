package blockchain_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildRegistry(
	t *testing.T,
) (blockchain.INodeRegistryAdmin, blockchain.INodeRegistryCaller, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := testutils.NewLog(t)
	rpcUrl := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.SettlementChain.RpcURL)
	require.NoError(t, err)

	registry, err := blockchain.NewNodeRegistryAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	caller, err := blockchain.NewNodeRegistryCaller(logger, client, contractsOptions)
	require.NoError(t, err)

	return registry, caller, ctx
}

func addRandomNode(
	t *testing.T,
	registry blockchain.INodeRegistryAdmin,
	ctx context.Context,
) {
	privateKey := testutils.RandomPrivateKey(t)
	owner := testutils.RandomAddress()
	httpAddress := testutils.RandomString(32)

	require.Eventually(t, func() bool {
		err := registry.AddNode(
			ctx,
			owner.String(),
			&privateKey.PublicKey,
			httpAddress,
		)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
}

func TestAddNode(t *testing.T) {
	registry, _, ctx := buildRegistry(t)

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	require.Eventually(t, func() bool {
		err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
}

func TestAddNodeBadOwner(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomString(10) // This is an invalid hex address

	registry, _, ctx := buildRegistry(t)
	err := registry.AddNode(ctx, owner, &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "invalid owner address provided")
}

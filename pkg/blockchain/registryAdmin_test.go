package blockchain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildRegistry(t *testing.T) (*NodeRegistryAdmin, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	_, baseChainOptions := testutils.GetContractsOptions(t)
	// Set the nodes contract address to a random smart contract instead of the fixed deployment
	baseChainOptions.NodesContractAddress = testutils.DeployNodesContract(t)

	signer, err := NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		baseChainOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := NewClient(ctx, baseChainOptions.RpcUrl)
	require.NoError(t, err)

	registry, err := NewNodeRegistryAdmin(logger, client, signer, baseChainOptions)
	require.NoError(t, err)

	return registry, ctx, func() {
		cancel()
	}
}

func TestAddNode(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
	require.NoError(t, err)
}

func TestAddNodeBadOwner(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	// This is an invalid hex address
	owner := testutils.RandomString(10)

	err := registry.AddNode(ctx, owner, &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "invalid owner address provided")
}

func TestAddNodeUnauthorized(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	// Create a signer that won't work
	_, baseChainOptions := testutils.GetContractsOptions(t)
	signer, err := NewPrivateKeySigner(
		utils.EcdsaPrivateKeyToString(testutils.RandomPrivateKey(t)),
		baseChainOptions.ChainID,
	)
	require.NoError(t, err)
	registry.signer = signer

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	err = registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "Out of gas")
}

package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildRegistry(t *testing.T) (*NodeRegistryAdmin, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)
	// Set the nodes contract address to a random smart contract instead of the fixed deployment
	contractsOptions.NodesContractAddress = testutils.DeployNodesContract(t)

	signer, err := NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	registry, err := NewNodeRegistryAdmin(logger, client, signer, contractsOptions)
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

	require.Eventually(t, func() bool {
		err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
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
	contractsOptions := testutils.GetContractsOptions(t)
	signer, err := NewPrivateKeySigner(
		utils.EcdsaPrivateKeyToString(testutils.RandomPrivateKey(t)),
		contractsOptions.ChainID,
	)
	require.NoError(t, err)
	registry.signer = signer

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	err = registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "Out of gas")
}

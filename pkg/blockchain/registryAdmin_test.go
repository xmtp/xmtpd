package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildRegistry(t *testing.T) (*NodeRegistryAdmin, *PrivateKeySigner, context.Context, func()) {
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

	return registry, signer, ctx, func() {
		cancel()
	}
}

func addRandomNode(
	t *testing.T,
	ctx context.Context,
	registry NodeRegistry,
) (uint32, *ecdsa.PrivateKey, error) {
	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	nodeId, err := registry.AddNode(
		ctx,
		owner.String(),
		&privateKey.PublicKey,
		httpAddress,
		big.NewInt(0),
	)
	require.NoError(t, err)
	require.NotZero(t, nodeId)

	return nodeId, privateKey, nil
}

func TestAddNode(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()
	_, _, err := addRandomNode(t, ctx, registry)
	require.NoError(t, err)
}

func TestAddNodeBadOwner(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	// This is an invalid hex address
	owner := testutils.RandomString(10)

	_, err := registry.AddNode(ctx, owner, &privateKey.PublicKey, httpAddress, big.NewInt(0))
	require.ErrorContains(t, err, "invalid owner address provided")
}

func TestAddNodeBadMinMonthlyFee(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	_, err := registry.AddNode(
		ctx,
		owner.String(),
		&privateKey.PublicKey,
		httpAddress,
		big.NewInt(-1),
	)
	require.ErrorContains(t, err, "invalid min monthly fee provided")
}

func TestAddNodeUnauthorized(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
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

	_, err = registry.AddNode(
		ctx,
		owner.String(),
		&privateKey.PublicKey,
		httpAddress,
		big.NewInt(0),
	)
	require.ErrorContains(t, err, "Out of gas")
}

func TestUpdateIsApiEnabled(t *testing.T) {
	// UpdateIsApiEnabled is callable only by the node owner.
	// For this test, registry signer is the node owner.
	registry, signer, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	signingKey := testutils.RandomPrivateKey(t)

	nodeId, err := registry.AddNode(
		ctx,
		signer.FromAddress().String(),
		&signingKey.PublicKey,
		testutils.RandomString(32),
		big.NewInt(0),
	)
	require.NoError(t, err)

	err = registry.UpdateIsApiEnabled(ctx, nodeId)
	require.NoError(t, err)
}

func TestUpdateIsApiEnabledUnauthorized(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	nodeId, _, err := addRandomNode(t, ctx, registry)
	require.NoError(t, err)

	err = registry.UpdateIsApiEnabled(ctx, nodeId)

	// 0x82b42900 is the signature of Unauthorized() error
	require.Equal(t, err.Error(), "execution reverted: custom error 0x82b42900")
}

func TestUpdateIsApiEnabledBadNodeId(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	err := registry.UpdateIsApiEnabled(ctx, 1)
	require.Equal(t, err.Error(), "execution reverted: custom error 0x82b42900")
}

func TestUpdateIsReplicationEnabled(t *testing.T) {
	registry, _, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	nodeId, _, err := addRandomNode(t, ctx, registry)
	require.NoError(t, err)

	err = registry.UpdateIsReplicationEnabled(ctx, nodeId, true)
	require.NoError(t, err)
}

func TestUpdateActive(t *testing.T) {
	registry, signer, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	signingKey := testutils.RandomPrivateKey(t)

	nodeId, err := registry.AddNode(
		ctx,
		signer.FromAddress().String(),
		&signingKey.PublicKey,
		testutils.RandomString(32),
		big.NewInt(0),
	)
	require.NoError(t, err)

	err = registry.UpdateIsApiEnabled(ctx, nodeId)
	require.NoError(t, err)

	err = registry.UpdateIsReplicationEnabled(ctx, nodeId, true)
	require.NoError(t, err)

	err = registry.UpdateActive(ctx, nodeId, true)
	require.NoError(t, err)
}

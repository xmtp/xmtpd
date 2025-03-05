package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildRegistry(t *testing.T) (INodeRegistryAdmin, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)

	// Deploy the contract always, so the tests are deterministic.
	contractsOptions.NodesContractAddress = testutils.DeployNodesV2Contract(t)

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

func addRandomNode(
	t *testing.T,
	registry INodeRegistryAdmin,
	ctx context.Context,
) {
	privateKey := testutils.RandomPrivateKey(t)
	owner := testutils.RandomAddress()
	httpAddress := testutils.RandomString(32)
	minMonthlyFee := int64(1000)

	require.Eventually(t, func() bool {
		err := registry.AddNode(
			ctx,
			owner.String(),
			&privateKey.PublicKey,
			httpAddress,
			minMonthlyFee,
		)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
}

func TestAddNode(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	require.Eventually(t, func() bool {
		err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress, 1000)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
}

func TestAddNodeBadOwner(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomString(10) // This is an invalid hex address

	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()
	err := registry.AddNode(ctx, owner, &privateKey.PublicKey, httpAddress, 1000)
	require.ErrorContains(t, err, "invalid owner address provided")
}

func TestAddNodeBadMinMonthlyFee(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()
	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress, -1)
	require.ErrorContains(t, err, "invalid min monthly fee provided")
}

func TestDisableNode(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	addRandomNode(t, registry, ctx)

	err := registry.DisableNode(ctx, 100)
	require.NoError(t, err)
}

func TestEnableNode(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	addRandomNode(t, registry, ctx)

	err := registry.EnableNode(ctx, 100)
	require.NoError(t, err)
}

func TestRemoveFromApiNodes(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	addRandomNode(t, registry, ctx)

	err := registry.RemoveFromApiNodes(ctx, 100)
	require.NoError(t, err)
}

func TestRemoveFromReplicationNodes(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	addRandomNode(t, registry, ctx)

	err := registry.RemoveFromReplicationNodes(ctx, 100)
	require.NoError(t, err)
}

func TestSetMaxActiveNodes(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	err := registry.SetMaxActiveNodes(ctx, 100)
	require.NoError(t, err)
}

func TestSetNodeOperatorCommissionPercent(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	err := registry.SetNodeOperatorCommissionPercent(ctx, 100)
	require.NoError(t, err)
}

func TestSetNodeOperatorCommissionPercentInvalid(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	err := registry.SetNodeOperatorCommissionPercent(ctx, -1)
	require.ErrorContains(t, err, "invalid commission percent provided")
}

func TestSetBaseURI(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	err := registry.SetBaseURI(ctx, "https://example.com/")
	require.NoError(t, err)
}

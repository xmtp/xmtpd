package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildRegistry(
	t *testing.T,
	version RegistryAdminVersion,
) (INodeRegistryAdmin, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)

	if version == RegistryAdminV1 {
		contractsOptions.NodesContractAddress = testutils.DeployNodesContract(t)
	}

	if version == RegistryAdminV2 {
		contractsOptions.NodesContractAddress = testutils.DeployNodesV2Contract(t)
	}

	signer, err := NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	registry, err := NewNodeRegistryAdmin(logger, client, signer, contractsOptions, version)
	require.NoError(t, err)

	return registry, ctx, func() {
		cancel()
	}
}

func TestAddNode(t *testing.T) {
	registry, ctx, cleanup := buildRegistry(t, RegistryAdminV1)
	defer cleanup()

	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomAddress()

	require.Eventually(t, func() bool {
		err := registry.AddNode(ctx, owner.String(), &privateKey.PublicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)

	registryV2, ctxV2, cleanupV2 := buildRegistry(t, RegistryAdminV2)
	defer cleanupV2()

	privateKeyV2 := testutils.RandomPrivateKey(t)
	httpAddressV2 := testutils.RandomString(32)
	ownerV2 := testutils.RandomAddress()

	require.Eventually(t, func() bool {
		err := registryV2.AddNode(ctxV2, ownerV2.String(), &privateKeyV2.PublicKey, httpAddressV2)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)
}

func TestAddNodeBadOwner(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	httpAddress := testutils.RandomString(32)
	owner := testutils.RandomString(10) // This is an invalid hex address

	registry, ctx, cleanup := buildRegistry(t, RegistryAdminV1)
	defer cleanup()
	err := registry.AddNode(ctx, owner, &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "invalid owner address provided")

	registryV2, ctxV2, cleanupV2 := buildRegistry(t, RegistryAdminV2)
	defer cleanupV2()
	err = registryV2.AddNode(ctxV2, owner, &privateKey.PublicKey, httpAddress)
	require.ErrorContains(t, err, "invalid owner address provided")
}

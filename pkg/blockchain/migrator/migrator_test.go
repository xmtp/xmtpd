package migrator

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func setupRegistry(
	t *testing.T,
	adminVersion blockchain.RegistryAdminVersion,
	callerVersion blockchain.RegistryCallerVersion,
) (blockchain.INodeRegistryAdmin, blockchain.INodeRegistryCaller, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)

	if adminVersion == blockchain.RegistryAdminV1 {
		contractsOptions.NodesContractAddress = testutils.DeployNodesContract(t)
	}

	if adminVersion == blockchain.RegistryAdminV2 {
		contractsOptions.NodesContractAddress = testutils.DeployNodesV2Contract(t)
	}

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		client,
		signer,
		contractsOptions,
		adminVersion,
	)
	require.NoError(t, err)

	registryCaller, err := blockchain.NewNodeRegistryCaller(
		logger,
		client,
		contractsOptions,
		callerVersion,
	)
	require.NoError(t, err)

	return registryAdmin, registryCaller, func() {
		cancel()
		client.Close()
	}
}

func registerRandomNode(
	t *testing.T,
	registryAdmin blockchain.INodeRegistryAdmin,
) SerializableNodeV1 {
	ownerAddress := testutils.RandomAddress().Hex()
	httpAddress := testutils.RandomString(30)
	publicKey := testutils.RandomPrivateKey(t).PublicKey
	require.Eventually(t, func() bool {
		err := registryAdmin.AddNode(context.Background(), ownerAddress, &publicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)

	return SerializableNodeV1{
		OwnerAddress:  ownerAddress,
		HttpAddress:   httpAddress,
		SigningKeyPub: utils.EcdsaPublicKeyToString(&publicKey),
		IsHealthy:     true,
		// We don't get the node ID here. Can live without for the tests
	}
}

func TestRegistryRead(t *testing.T) {
	registryAdminV1, registryCaller, cleanup := setupRegistry(
		t,
		blockchain.RegistryAdminV1,
		blockchain.RegistryCallerV1,
	)
	defer cleanup()

	node1 := registerRandomNode(t, registryAdminV1)
	node2 := registerRandomNode(t, registryAdminV1)

	nodes, err := ReadFromRegistry[SerializableNodeV1](registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	require.Equal(t, node1.OwnerAddress, nodes[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, nodes[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, nodes[0].SigningKeyPub)
	require.Equal(t, node1.IsHealthy, nodes[0].IsHealthy)

	require.Equal(t, node2.OwnerAddress, nodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, nodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, nodes[1].SigningKeyPub)
	require.Equal(t, node2.IsHealthy, nodes[1].IsHealthy)
}

func TestFileDump(t *testing.T) {
	registryAdminV1, registryCaller, cleanup := setupRegistry(
		t,
		blockchain.RegistryAdminV1,
		blockchain.RegistryCallerV1,
	)
	defer cleanup()

	_ = registerRandomNode(t, registryAdminV1)

	nodes, err := ReadFromRegistry[SerializableNodeV1](registryCaller)
	require.NoError(t, err)
	require.Len(t, nodes, 1)

	// Create a temporary file path for testing
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.json", testutils.RandomString(12)))

	err = DumpNodesToFile(nodes, tempFilePath)
	require.NoError(t, err)

	restoredNodes, err := ImportNodesFromFile(tempFilePath)
	require.NoError(t, err)

	require.Equal(t, nodes, restoredNodes)
}

func TestRegistryWrite(t *testing.T) {
	registryAdminV1, registryCallerV1, cleanupV1 := setupRegistry(
		t,
		blockchain.RegistryAdminV1,
		blockchain.RegistryCallerV1,
	)
	defer cleanupV1()

	node1 := registerRandomNode(t, registryAdminV1)
	node2 := registerRandomNode(t, registryAdminV1)

	nodes, err := ReadFromRegistry[SerializableNodeV1](registryCallerV1)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	// Build a NodesV2 registry
	registryAdminV2, registryCallerV2, cleanupV2 := setupRegistry(
		t,
		blockchain.RegistryAdminV2,
		blockchain.RegistryCallerV2,
	)
	defer cleanupV2()

	err = WriteToRegistryV2(testutils.NewLog(t), nodes, registryAdminV2)
	require.NoError(t, err)

	restoredNodes, err := ReadFromRegistry[SerializableNodeV2](registryCallerV2)
	require.NoError(t, err)
	require.Equal(t, 2, len(restoredNodes))

	// Old parameters should be the same.
	require.Equal(t, node1.OwnerAddress, restoredNodes[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, restoredNodes[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, restoredNodes[0].SigningKeyPub)
	require.Equal(t, node2.OwnerAddress, restoredNodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, restoredNodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, restoredNodes[1].SigningKeyPub)

	// New parameters should be the default values.
	require.Equal(t, "0", restoredNodes[0].MinMonthlyFee)
	require.Equal(t, false, restoredNodes[0].IsReplicationEnabled)
	require.Equal(t, false, restoredNodes[0].IsApiEnabled)
	require.Equal(t, "0", restoredNodes[1].MinMonthlyFee)
	require.Equal(t, false, restoredNodes[1].IsReplicationEnabled)
	require.Equal(t, false, restoredNodes[1].IsApiEnabled)
}

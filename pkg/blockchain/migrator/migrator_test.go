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

func setupTest(
	t *testing.T,
) (*blockchain.NodeRegistryAdmin, *blockchain.NodeRegistryCaller, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)
	// Set the nodes contract address to a random smart contract instead of the fixed deployment
	contractsOptions.NodesContractAddress = testutils.DeployNodesContract(t)

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
	)
	require.NoError(t, err)

	registryCaller, err := blockchain.NewNodeRegistryCaller(
		logger,
		client,
		contractsOptions,
	)
	require.NoError(t, err)

	return registryAdmin, registryCaller, func() {
		cancel()
		client.Close()
	}
}

func registerRandomNode(
	t *testing.T,
	registryAdmin *blockchain.NodeRegistryAdmin,
) SerializableNode {
	ownerAddress := testutils.RandomAddress().Hex()
	httpAddress := testutils.RandomString(30)
	publicKey := testutils.RandomPrivateKey(t).PublicKey
	require.Eventually(t, func() bool {
		err := registryAdmin.AddNode(context.Background(), ownerAddress, &publicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)

	return SerializableNode{
		OwnerAddress:  ownerAddress,
		HttpAddress:   httpAddress,
		SigningKeyPub: utils.EcdsaPublicKeyToString(&publicKey),
		IsHealthy:     true,
		// We don't get the node ID here. Can live without for the tests
	}
}

func TestRegistryRead(t *testing.T) {
	registryAdmin, registryCaller, cleanup := setupTest(t)
	defer cleanup()

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
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
	registryAdmin, registryCaller, cleanup := setupTest(t)
	defer cleanup()

	_ = registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
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
	registryAdmin, registryCaller, cleanup := setupTest(t)
	defer cleanup()

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	// Build a new registry admin with a different nodes contract
	newRegistryAdmin, newRegistryCaller, newCleanup := setupTest(t)
	defer newCleanup()
	err = WriteToRegistry(testutils.NewLog(t), nodes, newRegistryAdmin)
	require.NoError(t, err)

	restoredNodes, err := ReadFromRegistry(newRegistryCaller)
	require.NoError(t, err)
	require.Equal(t, nodes, restoredNodes)

	require.Equal(t, node1.OwnerAddress, restoredNodes[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, restoredNodes[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, restoredNodes[0].SigningKeyPub)
	require.Equal(t, node1.IsHealthy, restoredNodes[0].IsHealthy)

	require.Equal(t, node2.OwnerAddress, restoredNodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, restoredNodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, restoredNodes[1].SigningKeyPub)
	require.Equal(t, node2.IsHealthy, restoredNodes[1].IsHealthy)
}

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
) (blockchain.INodeRegistryAdmin, blockchain.INodeRegistryCaller, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)
	contractsOptions.NodesContractAddress = testutils.DeployNodesV2Contract(t)

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
	registryAdmin blockchain.INodeRegistryAdmin,
) SerializableNode {
	ownerAddress := testutils.RandomAddress().Hex()
	httpAddress := testutils.RandomString(30)
	publicKey := testutils.RandomPrivateKey(t).PublicKey
	require.Eventually(t, func() bool {
		err := registryAdmin.AddNode(context.Background(), ownerAddress, &publicKey, httpAddress, 0)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)

	return SerializableNode{
		OwnerAddress:         ownerAddress,
		HttpAddress:          httpAddress,
		SigningKeyPub:        utils.EcdsaPublicKeyToString(&publicKey),
		MinMonthlyFee:        0,
		IsReplicationEnabled: false,
		IsApiEnabled:         false,
	}
}

func TestRegistryRead(t *testing.T) {
	registryAdmin, registryCaller, cleanup := setupRegistry(t)
	defer cleanup()

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	require.Equal(t, node1.OwnerAddress, nodes[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, nodes[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, nodes[0].SigningKeyPub)
	require.Equal(t, node1.MinMonthlyFee, nodes[0].MinMonthlyFee)
	require.Equal(t, node1.IsReplicationEnabled, nodes[0].IsReplicationEnabled)
	require.Equal(t, node1.IsApiEnabled, nodes[0].IsApiEnabled)

	require.Equal(t, node2.OwnerAddress, nodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, nodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, nodes[1].SigningKeyPub)
	require.Equal(t, node2.MinMonthlyFee, nodes[1].MinMonthlyFee)
	require.Equal(t, node2.IsReplicationEnabled, nodes[1].IsReplicationEnabled)
	require.Equal(t, node2.IsApiEnabled, nodes[1].IsApiEnabled)
}

func TestFileDump(t *testing.T) {
	registryAdmin, registryCaller, cleanup := setupRegistry(t)
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
	registryAdmin, registryCaller, cleanup := setupRegistry(t)
	defer cleanup()

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	registryAdmin2, registryCaller2, cleanup2 := setupRegistry(t)
	defer cleanup2()
	err = WriteToRegistry(testutils.NewLog(t), nodes, registryAdmin2)
	require.NoError(t, err)

	restoredNodes, err := ReadFromRegistry(registryCaller2)
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
	require.Equal(t, int64(0), restoredNodes[0].MinMonthlyFee)
	require.Equal(t, false, restoredNodes[0].IsReplicationEnabled)
	require.Equal(t, false, restoredNodes[0].IsApiEnabled)
	require.Equal(t, int64(0), restoredNodes[1].MinMonthlyFee)
	require.Equal(t, false, restoredNodes[1].IsReplicationEnabled)
	require.Equal(t, false, restoredNodes[1].IsApiEnabled)
}

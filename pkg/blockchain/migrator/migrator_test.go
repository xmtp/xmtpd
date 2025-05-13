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
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func setupRegistry(
	t *testing.T,
) (blockchain.INodeRegistryAdmin, blockchain.INodeRegistryCaller) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})
	logger := testutils.NewLog(t)
	rpcUrl := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(rpcUrl)
	contractsOptions.SettlementChain.NodeRegistryAddress = testutils.DeployNodesContract(t, rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.SettlementChain.RpcURL)
	t.Cleanup(func() {
		client.Close()
	})
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

	return registryAdmin, registryCaller
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
		OwnerAddress:              ownerAddress,
		HttpAddress:               httpAddress,
		SigningKeyPub:             utils.EcdsaPublicKeyToString(&publicKey),
		MinMonthlyFeeMicroDollars: 0,
	}
}

func TestRegistryRead(t *testing.T) {
	registryAdmin, registryCaller := setupRegistry(t)

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	require.Equal(t, node1.OwnerAddress, nodes[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, nodes[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, nodes[0].SigningKeyPub)
	require.Equal(t, node1.MinMonthlyFeeMicroDollars, nodes[0].MinMonthlyFeeMicroDollars)
	require.Equal(t, node1.InCanonicalNetwork, nodes[0].InCanonicalNetwork)

	require.Equal(t, node2.OwnerAddress, nodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, nodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, nodes[1].SigningKeyPub)
	require.Equal(t, node2.MinMonthlyFeeMicroDollars, nodes[1].MinMonthlyFeeMicroDollars)
	require.Equal(t, node2.InCanonicalNetwork, nodes[1].InCanonicalNetwork)
}

func TestFileDump(t *testing.T) {
	registryAdmin, registryCaller := setupRegistry(t)

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
	registryAdmin, registryCaller := setupRegistry(t)

	node1 := registerRandomNode(t, registryAdmin)
	node2 := registerRandomNode(t, registryAdmin)

	nodes, err := ReadFromRegistry(registryCaller)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodes))

	registryAdmin2, registryCaller2 := setupRegistry(t)
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
	require.Equal(t, int64(0), restoredNodes[0].MinMonthlyFeeMicroDollars)
	require.Equal(t, false, restoredNodes[0].InCanonicalNetwork)
	require.Equal(t, int64(0), restoredNodes[1].MinMonthlyFeeMicroDollars)
	require.Equal(t, false, restoredNodes[1].InCanonicalNetwork)
}

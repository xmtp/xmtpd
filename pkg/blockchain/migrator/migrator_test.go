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
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(
		ctx,
		contractsOptions.SettlementChain.RPCURL,
	)
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
		_, err := registryAdmin.AddNode(context.Background(), ownerAddress, &publicKey, httpAddress)
		return err == nil
	}, 1*time.Second, 50*time.Millisecond)

	return SerializableNode{
		OwnerAddress:  ownerAddress,
		HttpAddress:   httpAddress,
		SigningKeyPub: utils.EcdsaPublicKeyToString(&publicKey),
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
	require.Equal(t, node1.InCanonicalNetwork, nodes[0].InCanonicalNetwork)

	require.Equal(t, node2.OwnerAddress, nodes[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, nodes[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, nodes[1].SigningKeyPub)
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

	oldNodes := make([]SerializableNode, 0)

	registryAdmin2, registryCaller2 := setupRegistry(t)
	err = WriteToRegistry(t.Context(), nodes, oldNodes, registryAdmin2)
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
	require.Equal(t, false, restoredNodes[0].InCanonicalNetwork)
	require.Equal(t, false, restoredNodes[1].InCanonicalNetwork)
}

// 1) Empty -> New registry: adds nodes, keeps defaults (not in canonical network)
func TestRegistryWrite_AddsNewNodesAndKeepsDefaults(t *testing.T) {
	registryAdminSrc, registryCallerSrc := setupRegistry(t)

	// Source: register two random nodes (not in canonical network by default).
	node1 := registerRandomNode(t, registryAdminSrc)
	node2 := registerRandomNode(t, registryAdminSrc)

	newNodes, err := ReadFromRegistry(registryCallerSrc)
	require.NoError(t, err)
	require.Equal(t, 2, len(newNodes))

	// Destination: fresh/empty registry.
	registryAdminDst, registryCallerDst := setupRegistry(t)

	oldNodesDst, err := ReadFromRegistry(registryCallerDst)
	require.NoError(t, err)
	require.Len(t, oldNodesDst, 0)

	err = WriteToRegistry(t.Context(), newNodes, oldNodesDst, registryAdminDst)
	require.NoError(t, err)

	restored, err := ReadFromRegistry(registryCallerDst)
	require.NoError(t, err)
	require.Equal(t, 2, len(restored))

	require.Equal(t, node1.OwnerAddress, restored[0].OwnerAddress)
	require.Equal(t, node1.HttpAddress, restored[0].HttpAddress)
	require.Equal(t, node1.SigningKeyPub, restored[0].SigningKeyPub)
	require.Equal(t, node2.OwnerAddress, restored[1].OwnerAddress)
	require.Equal(t, node2.HttpAddress, restored[1].HttpAddress)
	require.Equal(t, node2.SigningKeyPub, restored[1].SigningKeyPub)

	require.False(t, restored[0].InCanonicalNetwork)
	require.False(t, restored[1].InCanonicalNetwork)
}

// 2) Nodes already registered: do NOT re-add; only AddToNetwork when requested
func TestRegistryWrite_AddsToNetworkForExistingNodes(t *testing.T) {
	admin, caller := setupRegistry(t)

	// Seed registry with two nodes (not in canonical by default).
	registerRandomNode(t, admin)
	registerRandomNode(t, admin)

	existing, err := ReadFromRegistry(caller)
	require.NoError(t, err)
	require.Equal(t, 2, len(existing))

	desired := make([]SerializableNode, len(existing))
	copy(desired, existing)
	for i := range desired {
		desired[i].InCanonicalNetwork = true
	}

	err = WriteToRegistry(t.Context(), desired, existing, admin)
	require.NoError(t, err)

	after, err := ReadFromRegistry(caller)
	require.NoError(t, err)
	require.Equal(t, 2, len(after))
	require.True(t, after[0].InCanonicalNetwork)
	require.True(t, after[1].InCanonicalNetwork)

	for i := range after {
		require.Equal(t, existing[i].OwnerAddress, after[i].OwnerAddress)
		require.Equal(t, existing[i].HttpAddress, after[i].HttpAddress)
		require.Equal(t, existing[i].SigningKeyPub, after[i].SigningKeyPub)
	}
}

// 3) Mixed case: one node already exists; one is brand new.
// Existing should be added to network; new should be added then added to network.
func TestRegistryWrite_MixedExistingAndNew(t *testing.T) {
	// Target registry where migration will write.
	adminDst, callerDst := setupRegistry(t)

	// Pre-seed destination with ONE node (exists, not canonical).
	existingNode := registerRandomNode(t, adminDst)

	// Snapshot "oldNodes" from destination.
	oldDst, err := ReadFromRegistry(callerDst)
	require.NoError(t, err)
	require.Equal(t, 1, len(oldDst))
	require.Equal(t, existingNode.SigningKeyPub, oldDst[0].SigningKeyPub)

	// Create a NEW node we want to add (use another registry just to fabricate a valid node struct).
	adminTmp, callerTmp := setupRegistry(t)
	newNode := registerRandomNode(t, adminTmp)
	tmpNodes, err := ReadFromRegistry(callerTmp)
	require.NoError(t, err)
	require.Equal(t, 1, len(tmpNodes))

	desired := []SerializableNode{
		{
			NodeID:             oldDst[0].NodeID,
			OwnerAddress:       oldDst[0].OwnerAddress,
			HttpAddress:        oldDst[0].HttpAddress,
			SigningKeyPub:      oldDst[0].SigningKeyPub,
			InCanonicalNetwork: true,
		},
		{
			OwnerAddress:       tmpNodes[0].OwnerAddress,
			HttpAddress:        tmpNodes[0].HttpAddress,
			SigningKeyPub:      tmpNodes[0].SigningKeyPub,
			InCanonicalNetwork: true,
		},
	}

	err = WriteToRegistry(context.Background(), desired, oldDst, adminDst)
	require.NoError(t, err)

	after, err := ReadFromRegistry(callerDst)
	require.NoError(t, err)
	require.Equal(t, 2, len(after))

	byPub := map[string]SerializableNode{
		after[0].SigningKeyPub: after[0],
		after[1].SigningKeyPub: after[1],
	}

	gotExisting, ok1 := byPub[existingNode.SigningKeyPub]
	require.True(t, ok1)
	require.True(t, gotExisting.InCanonicalNetwork)
	require.Equal(t, existingNode.OwnerAddress, gotExisting.OwnerAddress)
	require.Equal(t, existingNode.HttpAddress, gotExisting.HttpAddress)

	gotNew, ok2 := byPub[newNode.SigningKeyPub]
	require.True(t, ok2)
	require.True(t, gotNew.InCanonicalNetwork)
	require.Equal(t, tmpNodes[0].OwnerAddress, gotNew.OwnerAddress)
	require.Equal(t, tmpNodes[0].HttpAddress, gotNew.HttpAddress)
}

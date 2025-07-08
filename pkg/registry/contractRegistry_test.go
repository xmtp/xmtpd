package registry_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	r "github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const TEST_PUBKEY = "04760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1b01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d"

func requireNodeEquals(t *testing.T, a, b r.Node) {
	require.Condition(t, func() bool {
		return a.NodeID == b.NodeID && a.HttpAddress == b.HttpAddress
	})
}

func requireAllNodesEqual(t *testing.T, a, b []r.Node) {
	require.Equal(t, len(a), len(b))
	for i, node := range a {
		requireNodeEquals(t, node, b[i])
	}
}

func TestContractRegistryNewNodes(t *testing.T) {
	registry, err := r.NewSmartContractRegistry(context.Background(),
		nil,
		testutils.NewLog(t),
		config.ContractsOptions{
			SettlementChain: config.SettlementChainOptions{
				NodeRegistryRefreshInterval: 100 * time.Millisecond,
			},
		},
	)
	require.NoError(t, err)

	enc, err := hex.DecodeString(TEST_PUBKEY)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodeRegistryContract(t)
	mockContract.EXPECT().
		GetAllNodes(mock.Anything).
		Return([]noderegistry.INodeRegistryNodeWithId{
			{
				NodeId: 1,
				Node: noderegistry.INodeRegistryNode{
					HttpAddress:      "http://foo.com",
					SigningPublicKey: enc,
					IsCanonical:      true,
				},
			},
			{
				NodeId: 2,
				Node: noderegistry.INodeRegistryNode{
					HttpAddress:      "https://bar.com",
					SigningPublicKey: enc,
					IsCanonical:      true,
				},
			},
		}, nil)

	registry.SetContractForTest(mockContract)

	sub := registry.OnNewNodes()
	require.NoError(t, registry.Start())
	defer registry.Stop()
	newNodes := <-sub
	requireAllNodesEqual(t, []r.Node{
		{NodeID: 1, HttpAddress: "http://foo.com"},
		{NodeID: 2, HttpAddress: "https://bar.com"},
	},
		newNodes)
}

func TestContractRegistryChangedNodes(t *testing.T) {
	registry, err := r.NewSmartContractRegistry(context.Background(),
		nil,
		testutils.NewLog(t),
		config.ContractsOptions{
			SettlementChain: config.SettlementChainOptions{
				NodeRegistryRefreshInterval: 10 * time.Millisecond,
			},
		},
	)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodeRegistryContract(t)

	enc, err := hex.DecodeString(TEST_PUBKEY)
	require.NoError(t, err)

	hasSentInitialValues := false
	// The first call, we'll set the address to foo.com.
	// Subsequent calls will set the address to bar.com
	mockContract.EXPECT().
		GetAllNodes(mock.Anything).
		RunAndReturn(func(*bind.CallOpts) ([]noderegistry.INodeRegistryNodeWithId, error) {
			httpAddress := "http://foo.com"
			if !hasSentInitialValues {
				hasSentInitialValues = true
			} else {
				httpAddress = "http://bar.com"
			}
			return []noderegistry.INodeRegistryNodeWithId{
				{
					NodeId: 1,
					Node: noderegistry.INodeRegistryNode{
						HttpAddress:      httpAddress,
						SigningPublicKey: enc,
						IsCanonical:      true,
					},
				},
			}, nil
		})

	// Override the contract in the registry with a mock before calling Start
	registry.SetContractForTest(mockContract)

	sub := registry.OnChangedNode(1)

	getCurrentCount := r.CountChannel(sub, func(node r.Node) {
		require.Equal(t, "http://bar.com", node.HttpAddress)
	})

	require.NoError(t, registry.Start())
	defer registry.Stop()

	// ensure there is at least notification
	require.Eventually(t, func() bool {
		return getCurrentCount() == 1
	}, 1*time.Second, 10*time.Millisecond)
}

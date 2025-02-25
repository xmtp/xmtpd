package registry_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
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
		config.ContractsOptions{RefreshInterval: 100 * time.Millisecond},
	)
	require.NoError(t, err)

	enc, err := hex.DecodeString(TEST_PUBKEY)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)
	mockContract.EXPECT().
		GetActiveNodes(mock.Anything).
		Return([]nodes.INodesNodeWithId{
			{
				NodeId: big.NewInt(1),
				Node: nodes.INodesNode{
					HttpAddress:   "http://foo.com",
					SigningKeyPub: enc,
					IsActive:      true,
				},
			},
			{
				NodeId: big.NewInt(2),
				Node: nodes.INodesNode{
					HttpAddress:   "https://bar.com",
					SigningKeyPub: enc,
					IsActive:      true,
				},
			},
			{
				NodeId: big.NewInt(3),
				Node: nodes.INodesNode{
					HttpAddress:   "https://baz.com",
					SigningKeyPub: enc,
					IsActive:      false,
				},
			},
		}, nil)

	registry.SetContractForTest(mockContract)

	sub, cancelSub := registry.OnNewNodes()
	defer cancelSub()
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
		config.ContractsOptions{RefreshInterval: 10 * time.Millisecond},
	)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)

	enc, err := hex.DecodeString(TEST_PUBKEY)
	require.NoError(t, err)

	hasSentInitialValues := false
	// The first call, we'll set the address to foo.com.
	// Subsequent calls will set the address to bar.com
	mockContract.EXPECT().
		GetActiveNodes(mock.Anything).
		RunAndReturn(func(*bind.CallOpts) ([]nodes.INodesNodeWithId, error) {
			httpAddress := "http://foo.com"
			if !hasSentInitialValues {
				hasSentInitialValues = true
			} else {
				httpAddress = "http://bar.com"
			}
			return []nodes.INodesNodeWithId{
				{
					NodeId: big.NewInt(1),
					Node: nodes.INodesNode{
						HttpAddress:   httpAddress,
						SigningKeyPub: enc,
						IsActive:      true,
					},
				},
			}, nil
		})

	// Override the contract in the registry with a mock before calling Start
	registry.SetContractForTest(mockContract)

	// Subscribe to new nodes
	newSub, cancelNewSub := registry.OnNewNodes()
	defer cancelNewSub()

	// Subscribe to changed nodes
	sub, cancelSub := registry.OnChangedNode(1)
	defer cancelSub()

	counterSub, cancelCounter := registry.OnChangedNode(1)
	getCurrentCount := r.CountChannel(counterSub)
	defer cancelCounter()

	// Verify initial state
	go func() {
		nodes := <-newSub
		require.Equal(t, nodes[0].NodeID, uint32(1))
		require.Equal(t, nodes[0].HttpAddress, "http://foo.com")
	}()

	// Verify changed nodes
	go func() {
		for node := range sub {
			require.Equal(t, node.HttpAddress, "http://bar.com")
		}
	}()

	require.NoError(t, registry.Start())
	defer registry.Stop()
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, getCurrentCount(), 1)
}

func TestStopOnContextCancel(t *testing.T) {
	registry, err := r.NewSmartContractRegistry(context.Background(),
		nil,
		testutils.NewLog(t),
		config.ContractsOptions{RefreshInterval: 10 * time.Millisecond},
	)
	require.NoError(t, err)

	enc, err := hex.DecodeString(TEST_PUBKEY)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)
	mockContract.EXPECT().
		GetActiveNodes(mock.Anything).
		RunAndReturn(func(*bind.CallOpts) ([]nodes.INodesNodeWithId, error) {
			return []nodes.INodesNodeWithId{
				{
					NodeId: big.NewInt(int64(rand.Intn(1000))),
					Node: nodes.INodesNode{
						HttpAddress:   "http://foo.com",
						SigningKeyPub: enc,
						IsActive:      true,
					},
				},
			}, nil
		})

	registry.SetContractForTest(mockContract)

	sub, cancelSub := registry.OnNewNodes()
	defer cancelSub()
	getCurrentCount := r.CountChannel(sub)

	require.NoError(t, registry.Start())

	// Accumulate some new nodes
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, getCurrentCount(), 0)

	// Stop the registry
	registry.Stop()

	// Wait for a little bit to give the cancellation time to take effect
	time.Sleep(10 * time.Millisecond)
	currentNodeCount := getCurrentCount()
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, currentNodeCount, getCurrentCount())
}

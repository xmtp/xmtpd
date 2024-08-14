package registry

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/mocks"
	testUtils "github.com/xmtp/xmtpd/pkg/testing"
)

func TestContractRegistryNewNodes(t *testing.T) {
	registry, err := NewSmartContractRegistry(
		nil,
		testUtils.NewLog(t),
		config.ContractsOptions{RefreshInterval: 100 * time.Millisecond},
	)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)
	mockContract.EXPECT().
		AllNodes(mock.Anything).
		Return([]abis.NodesNodeWithId{
			{NodeId: 1, Node: abis.NodesNode{HttpAddress: "http://foo.com"}},
			{NodeId: 2, Node: abis.NodesNode{HttpAddress: "https://bar.com"}},
		}, nil)

	registry.contract = mockContract

	sub, cancelSub := registry.OnNewNodes()
	defer cancelSub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	require.NoError(t, registry.Start(ctx))
	newNodes := <-sub
	require.Equal(
		t,
		[]Node{
			{NodeID: 1, HttpAddress: "http://foo.com"},
			{NodeID: 2, HttpAddress: "https://bar.com"},
		},
		newNodes,
	)
}

func TestContractRegistryChangedNodes(t *testing.T) {
	registry, err := NewSmartContractRegistry(
		nil,
		testUtils.NewLog(t),
		config.ContractsOptions{RefreshInterval: 10 * time.Millisecond},
	)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)

	hasSentInitialValues := false
	// The first call, we'll set the address to foo.com.
	// Subsequent calls will set the address to bar.com
	mockContract.EXPECT().
		AllNodes(mock.Anything).RunAndReturn(func(*bind.CallOpts) ([]abis.NodesNodeWithId, error) {
		httpAddress := "http://foo.com"
		if !hasSentInitialValues {
			hasSentInitialValues = true
		} else {
			httpAddress = "http://bar.com"
		}
		return []abis.NodesNodeWithId{
			{NodeId: 1, Node: abis.NodesNode{HttpAddress: httpAddress}},
		}, nil
	})

	// Override the contract in the registry with a mock before calling Start
	registry.contract = mockContract

	sub, cancelSub := registry.OnChangedNode(1)
	defer cancelSub()
	counterSub, cancelCounter := registry.OnChangedNode(1)
	getCurrentCount := countChannel(counterSub)
	defer cancelCounter()
	go func() {
		for node := range sub {
			require.Equal(t, node.HttpAddress, "http://bar.com")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	require.NoError(t, registry.Start(ctx))
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, getCurrentCount(), 1)
}

func TestStopOnContextCancel(t *testing.T) {
	registry, err := NewSmartContractRegistry(
		nil,
		testUtils.NewLog(t),
		config.ContractsOptions{RefreshInterval: 10 * time.Millisecond},
	)
	require.NoError(t, err)

	mockContract := mocks.NewMockNodesContract(t)
	mockContract.EXPECT().
		AllNodes(mock.Anything).
		RunAndReturn(func(*bind.CallOpts) ([]abis.NodesNodeWithId, error) {
			return []abis.NodesNodeWithId{
				{
					NodeId: uint16(rand.Intn(1000)),
					Node:   abis.NodesNode{HttpAddress: "http://foo.com"},
				},
			}, nil
		})

	registry.contract = mockContract

	sub, cancelSub := registry.OnNewNodes()
	defer cancelSub()
	getCurrentCount := countChannel(sub)

	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, registry.Start(ctx))
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, getCurrentCount(), 0)
	// Cancel the context
	cancel()
	// Wait for a little bit to give the cancellation time to take effect
	time.Sleep(10 * time.Millisecond)
	currentNodeCount := getCurrentCount()
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, currentNodeCount, getCurrentCount())
}

package payer_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/registry"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func TestGetNode_StableAssignment(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(700),
		nodeRegistry.GetHealthyNode(1200),
	}, nil)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("deadbeef"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)
	// Get assigned node multiple times
	node1, err := selector.GetNode(tpc)
	require.NoError(t, err)
	node2, err := selector.GetNode(tpc)
	require.NoError(t, err)

	if node1 != node2 {
		t.Errorf("Stable hashing failed, expected the same node but got %d and %d", node1, node2)
	}
}

func TestGetNode_EmptyNodes(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{}, nil)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))
	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	_, err := selector.GetNode(tpc)
	require.Error(t, err, "Expected an error for empty node list, but got none")
}

func TestGetNode_NoAvailableNodesError(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nil, errors.New("node fetch error"))

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))
	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	_, err := selector.GetNode(tpc)
	require.Error(t, err)
	require.Equal(t, "node fetch error", err.Error(), "Expected registry error message")
}

func TestGetNode_CorrectAssignment(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)
	tpc1 := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))
	tpc2 := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key_topic2"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)
	node2, err := selector.GetNode(tpc2)
	require.NoError(t, err)

	require.NotEqual(t, node1, node2, "Different topics should be assigned to different nodes")
}

func TestGetNode_NodeReassignment(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(700),
		nodeRegistry.GetHealthyNode(1200),
	}, nil)
	tpc1 := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	// pretend to remove the hashed node
	newMockRegistry := mocks.NewMockNodeRegistry(t)
	newMockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil)

	selector = payer.NewStableHashingNodeSelectorAlgorithm(newMockRegistry)

	newNode, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	require.NotEqual(t, node1, newNode, "Reassignment should be assigned to different node")
}

func TestGetNode_NodeReassignmentStability(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(700),
		nodeRegistry.GetHealthyNode(1200),
	}, nil)
	tpc1 := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	// pretend to remove the hashed node
	newMockRegistry := mocks.NewMockNodeRegistry(t)
	newMockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(700),
		nodeRegistry.GetHealthyNode(800),
		nodeRegistry.GetHealthyNode(1200),
	}, nil)

	selector = payer.NewStableHashingNodeSelectorAlgorithm(newMockRegistry)

	newNode, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	require.Equal(t, node1, newNode, "Reassignment be stable")
}

// Example usage
func TestGetNode_FindTopics(t *testing.T) {
	t.Skipf("This test helps with generation of payloads. No need to run it continuously")
	nodes := []registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(700),
		nodeRegistry.GetHealthyNode(1200),
	}

	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nodes, nil)

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	targetBuckets := make(map[uint32]topic.Topic)
	// Brute-force search for topics that hash into each bucket
	for i := 0; i < 1000000; i++ {
		tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte(strconv.Itoa(i)))
		node, err := selector.GetNode(tpc)
		require.NoError(t, err)
		if _, exists := targetBuckets[node]; !exists {
			targetBuckets[node] = tpc
		}

		// Stop early if all buckets are filled
		if len(targetBuckets) == len(nodes) {
			break
		}
	}

	fmt.Println("Generated Topics for Each Bucket:")
	for nodeID, tpc := range targetBuckets {
		fmt.Printf("Bucket %d -> Topic: %s\n", nodeID, tpc.String())
	}
}

func TestGetNode_ConfirmTopicBalance(t *testing.T) {
	nodes := []registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
		nodeRegistry.GetHealthyNode(400),
		nodeRegistry.GetHealthyNode(500),
		nodeRegistry.GetHealthyNode(1200),
		nodeRegistry.GetHealthyNode(8000),
	}

	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nodes, nil)

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	const totalRequests = 10000

	targetBuckets := make(map[uint32]uint32)
	// Brute-force search for topics that hash into each bucket
	for i := 0; i < totalRequests; i++ {
		tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte(strconv.Itoa(i)))
		node, err := selector.GetNode(tpc)
		require.NoError(t, err)
		if _, exists := targetBuckets[node]; !exists {
			targetBuckets[node] = 1
		} else {
			targetBuckets[node] += 1
		}
	}
	t.Logf("Target Buckets: %v", targetBuckets)

	// Compute expected balance
	expectedPerNode := totalRequests / uint32(len(nodes))
	tolerance := float64(expectedPerNode) * 0.05 // 10% tolerance
	t.Logf("Target Tolerance: %v", tolerance)

	// Verify that each bucket is within ±10% of expected distribution
	for nodeID, count := range targetBuckets {
		require.InDeltaf(t, expectedPerNode, count, tolerance,
			"Node %d has an imbalance: expected ~%d but got %d", nodeID, expectedPerNode, count)
	}
}

func TestGetNode_NodeGetNextIfBanned(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)
	tpc1 := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("stable_key"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)
	require.EqualValues(t, 200, node1)

	banlist := []uint32{node1}

	reselectedNode, err := selector.GetNode(tpc1, banlist)
	require.NoError(t, err)
	require.NotEqualValues(t, node1, reselectedNode)
	require.EqualValues(t, 300, reselectedNode)

	banlist = append(banlist, reselectedNode)

	reselectedNode, err = selector.GetNode(tpc1, banlist)
	require.NoError(t, err)
	require.NotEqualValues(t, node1, reselectedNode)
	require.EqualValues(t, 100, reselectedNode)

	banlist = append(banlist, reselectedNode)

	// now we are out of nodes to try
	reselectedNode, err = selector.GetNode(tpc1, banlist)
	require.Error(t, err)
	require.EqualValues(t, 0, reselectedNode)
}

func TestManualNodeSelector_SingleNode(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewManualNodeSelectorAlgorithm(mockRegistry, []uint32{200})
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.EqualValues(t, 200, node)
}

func TestManualNodeSelector_MultipleNodes(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewManualNodeSelectorAlgorithm(mockRegistry, []uint32{200, 300, 100})
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.EqualValues(t, 200, node)
}

func TestManualNodeSelector_WithBanlist(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewManualNodeSelectorAlgorithm(mockRegistry, []uint32{200, 300, 100})
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc, []uint32{200})
	require.NoError(t, err)
	require.EqualValues(t, 300, node)
}

func TestManualNodeSelector_NoNodesConfigured(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)

	selector := payer.NewManualNodeSelectorAlgorithm(mockRegistry, []uint32{})
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	_, err := selector.GetNode(tpc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no manual nodes configured")
}

func TestManualNodeSelector_NodeNotInRegistry(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil)

	selector := payer.NewManualNodeSelectorAlgorithm(mockRegistry, []uint32{200, 300})
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	_, err := selector.GetNode(tpc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no available manual nodes")
}

func TestOrderedPreferenceNodeSelector_FirstPreferred(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewOrderedPreferenceNodeSelectorAlgorithm(
		mockRegistry,
		[]uint32{300, 200, 100},
	)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.EqualValues(t, 300, node)
}

func TestOrderedPreferenceNodeSelector_FallbackToSecond(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewOrderedPreferenceNodeSelectorAlgorithm(
		mockRegistry,
		[]uint32{400, 300, 200},
	)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.EqualValues(t, 300, node)
}

func TestOrderedPreferenceNodeSelector_WithBanlist(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewOrderedPreferenceNodeSelectorAlgorithm(
		mockRegistry,
		[]uint32{300, 200, 100},
	)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc, []uint32{300, 200})
	require.NoError(t, err)
	require.EqualValues(t, 100, node)
}

func TestOrderedPreferenceNodeSelector_FallbackToAny(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewOrderedPreferenceNodeSelectorAlgorithm(
		mockRegistry,
		[]uint32{400, 500},
	)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.Contains(t, []uint32{100, 200, 300}, node)
}

func TestRandomNodeSelector_ReturnsValidNode(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewRandomNodeSelectorAlgorithm(mockRegistry)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	require.NoError(t, err)
	require.Contains(t, []uint32{100, 200, 300}, node)
}

func TestRandomNodeSelector_Distribution(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewRandomNodeSelectorAlgorithm(mockRegistry)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	counts := make(map[uint32]int)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		node, err := selector.GetNode(tpc)
		require.NoError(t, err)
		counts[node]++
	}

	require.Len(t, counts, 3)
	for nodeID, count := range counts {
		require.Greater(t, count, 200, "Node %d should be selected at least 200 times", nodeID)
	}
}

func TestRandomNodeSelector_WithBanlist(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	selector := payer.NewRandomNodeSelectorAlgorithm(mockRegistry)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc, []uint32{100, 200})
	require.NoError(t, err)
	require.EqualValues(t, 300, node)
}

func TestRandomNodeSelector_EmptyNodes(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{}, nil)

	selector := payer.NewRandomNodeSelectorAlgorithm(mockRegistry)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	_, err := selector.GetNode(tpc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no available nodes")
}

func TestClosestNodeSelector_ReturnsNode(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil)

	selector := payer.NewClosestNodeSelectorAlgorithm(mockRegistry, 0, 0)
	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	node, err := selector.GetNode(tpc)
	if err == nil {
		require.EqualValues(t, 100, node)
	} else {
		require.Contains(t, err.Error(), "no available nodes with latency measurements")
	}
} 
 
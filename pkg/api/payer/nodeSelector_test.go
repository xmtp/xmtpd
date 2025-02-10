package payer_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"strconv"
	"testing"
)

func TestGetNode_StableAssignment(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 1200},
	}, nil)
	tpc := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("deadbeef"))

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
	tpc := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key"))
	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	_, err := selector.GetNode(tpc)
	require.Error(t, err, "Expected an error for empty node list, but got none")
}

func TestGetNode_NoAvailableNodesError(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nil, errors.New("node fetch error"))

	tpc := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key"))
	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	_, err := selector.GetNode(tpc)
	require.Error(t, err)
	require.Equal(t, "node fetch error", err.Error(), "Expected registry error message")
}

func TestGetNode_CorrectAssignment(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 1200},
	}, nil)
	tpc1 := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key"))
	tpc2 := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key_topic2"))

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
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 1200},
	}, nil)
	tpc1 := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	// pretend to remove the hashed node
	newMockRegistry := mocks.NewMockNodeRegistry(t)
	newMockRegistry.On("GetNodes").Return([]registry.Node{
		{NodeID: 100},
	}, nil)

	selector = payer.NewStableHashingNodeSelectorAlgorithm(newMockRegistry)

	newNode, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	require.NotEqual(t, node1, newNode, "Reassignment should be assigned to different node")
}

func TestGetNode_NodeReassignmentStability(t *testing.T) {
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return([]registry.Node{
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 1200},
	}, nil)
	tpc1 := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte("stable_key"))

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	node1, err := selector.GetNode(tpc1)
	require.NoError(t, err)

	// pretend to remove the hashed node
	newMockRegistry := mocks.NewMockNodeRegistry(t)
	newMockRegistry.On("GetNodes").Return([]registry.Node{
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 800},
		{NodeID: 1200},
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
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 700},
		{NodeID: 1200},
	}

	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nodes, nil)

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	targetBuckets := make(map[uint32]topic.Topic)
	// Brute-force search for topics that hash into each bucket
	for i := 0; i < 1000000; i++ {
		tpc := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte(strconv.Itoa(i)))
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
		{NodeID: 100},
		{NodeID: 200},
		{NodeID: 300},
		{NodeID: 400},
		{NodeID: 500},
		{NodeID: 1200},
		{NodeID: 8000},
	}

	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Return(nodes, nil)

	selector := payer.NewStableHashingNodeSelectorAlgorithm(mockRegistry)

	const totalRequests = 10000

	targetBuckets := make(map[uint32]uint32)
	// Brute-force search for topics that hash into each bucket
	for i := 0; i < totalRequests; i++ {
		tpc := *topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, []byte(strconv.Itoa(i)))
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

package payer

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sort"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

type NodeSelectorAlgorithm interface {
	GetNode(topic topic.Topic, banlist ...[]uint32) (uint32, error)
}

type StableHashingNodeSelectorAlgorithm struct {
	reg registry.NodeRegistry
}

func NewStableHashingNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
) *StableHashingNodeSelectorAlgorithm {
	return &StableHashingNodeSelectorAlgorithm{reg: reg}
}

// hashKey hashes the topic to a stable uint16 hash
func HashKey(topic topic.Topic) uint32 {
	hash := sha256.Sum256(topic.Bytes())
	return binary.BigEndian.Uint32(hash[:4])
}

// GetNode selects a node for a given topic using stable hashing
func (s *StableHashingNodeSelectorAlgorithm) GetNode(
	topic topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	nodes, err := s.reg.GetNodes()
	if err != nil {
		return 0, err
	}

	// Flatten banlist
	banned := make(map[uint32]struct{})
	for _, list := range banlist {
		for _, id := range list {
			banned[id] = struct{}{}
		}
	}

	// Filter out banned, disabled, or non-API nodes
	var availableNodes []registry.Node
	for _, node := range nodes {
		if node.IsDisabled || !node.IsApiEnabled {
			continue
		}
		if _, exists := banned[node.NodeID]; !exists {
			availableNodes = append(availableNodes, node)
		}
	}

	if len(availableNodes) == 0 {
		if len(nodes) == 0 {
			return 0, errors.New("no available nodes")
		}
		return 0, errors.New("no available nodes after filtering")
	}

	// Sort availableNodes to ensure stability
	sort.Slice(availableNodes, func(i, j int) bool {
		return availableNodes[i].NodeID < availableNodes[j].NodeID
	})

	topicHash := HashKey(topic)

	numNodes := uint32(len(availableNodes))
	maxHashSpace := ^uint32(0)
	spacing := maxHashSpace / numNodes

	// Compute virtual positions for each available node
	nodeLocations := make([]uint32, numNodes)
	for i := range availableNodes {
		nodeLocations[i] = uint32(i) * spacing
	}

	// Binary search to find the first node with a virtual position >= topicHash
	idx := sort.Search(len(nodeLocations), func(i int) bool {
		return topicHash < nodeLocations[i]
	})

	// Select the appropriate node from availableNodes
	candidateIdx := idx % len(nodeLocations)
	return availableNodes[candidateIdx].NodeID, nil
}

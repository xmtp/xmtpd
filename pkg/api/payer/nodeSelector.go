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

	if len(nodes) == 0 {
		return 0, errors.New("no available nodes")
	}

	// Sort nodes to ensure stability
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].NodeID < nodes[j].NodeID })

	topicHash := HashKey(topic)

	numNodes := uint32(len(nodes))
	maxHashSpace := ^uint32(0)
	spacing := maxHashSpace / numNodes

	// Compute virtual positions for each node
	// Skip nodes that are disabled or do not have an API enabled
	nodeLocations := make([]uint32, numNodes)
	gotNodes := false
	for i, node := range nodes {
		if node.IsDisabled || !node.IsApiEnabled {
			continue
		}
		nodeLocations[i] = uint32(i) * spacing
		gotNodes = true
	}

	if !gotNodes {
		return 0, errors.New("no available active nodes")
	}

	// Binary search to find the first node with a virtual position >= topicHash
	idx := sort.Search(len(nodeLocations), func(i int) bool {
		return topicHash < nodeLocations[i]
	})

	// Flatten banlist
	banned := make(map[uint32]struct{})
	for _, list := range banlist {
		for _, id := range list {
			banned[id] = struct{}{}
		}
	}

	// Find the next available node
	for i := 0; i < len(nodes); i++ {
		candidateIdx := (idx + i) % len(nodeLocations)
		candidateNodeID := nodes[candidateIdx].NodeID

		if _, exists := banned[candidateNodeID]; !exists {
			return candidateNodeID, nil
		}
	}

	return 0, errors.New("no available nodes after considering banlist")
}

package payer

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"sort"
)

type NodeSelectorAlgorithm interface {
	GetNode(topic topic.Topic) (uint32, error)
}

type StableHashingNodeSelectorAlgorithm struct {
	reg registry.NodeRegistry
}

func NewStableHashingNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
) *StableHashingNodeSelectorAlgorithm {
	return &StableHashingNodeSelectorAlgorithm{reg: reg}
}

// hashKey hashes the topic to a stable uint32 hash
func HashKey(topic topic.Topic) uint32 {
	hash := sha256.Sum256(topic.Bytes())
	return binary.BigEndian.Uint32(hash[:4])
}

// GetNode selects a node for a given topic using stable hashing
func (s *StableHashingNodeSelectorAlgorithm) GetNode(topic topic.Topic) (uint32, error) {
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

	// we do not want to hash the entire available hash space, just consider up to the largest NodeID issued
	topicHash = topicHash % nodes[len(nodes)-1].NodeID

	// Use binary search to find the closest node
	// search returns the smallest element where condition is TRUE
	idx := sort.Search(
		len(nodes),
		func(i int) bool { return topicHash < nodes[i].NodeID },
	)

	//println(topicHash, idx, nodes[idx].NodeID, topic.String())

	return nodes[idx].NodeID, nil
}

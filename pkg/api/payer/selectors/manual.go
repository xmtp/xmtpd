package selectors

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// ManualNodeSelectorAlgorithm always selects the first available node from an explicit list of node IDs.
// The topic is ignored. If a node in the manual list is not present in the registry or is banned,
// it is skipped. If no node from the list is usable, selection fails.
//
// This strategy is useful for pinning traffic to specific nodes (e.g., during testing or partial rollouts).
type ManualNodeSelectorAlgorithm struct {
	reg     registry.NodeRegistry
	nodeIDs []uint32
}

var _ NodeSelectorAlgorithm = (*ManualNodeSelectorAlgorithm)(nil)

func NewManualNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
	nodeIDs []uint32,
) *ManualNodeSelectorAlgorithm {
	return &ManualNodeSelectorAlgorithm{
		reg:     reg,
		nodeIDs: nodeIDs,
	}
}

func (m *ManualNodeSelectorAlgorithm) GetNode(
	_ topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	if len(m.nodeIDs) == 0 {
		return 0, errors.New("no manual nodes configured")
	}

	nodes, err := m.reg.GetNodes()
	if err != nil {
		return 0, err
	}

	nodeMap := make(map[uint32]struct{})
	for _, node := range nodes {
		nodeMap[node.NodeID] = struct{}{}
	}

	banned := make(map[uint32]struct{})
	for _, list := range banlist {
		for _, id := range list {
			banned[id] = struct{}{}
		}
	}

	for _, nodeID := range m.nodeIDs {
		if _, exists := nodeMap[nodeID]; !exists {
			continue
		}
		if _, isBanned := banned[nodeID]; !isBanned {
			return nodeID, nil
		}
	}

	return 0, errors.New("no available manual nodes after considering banlist")
}

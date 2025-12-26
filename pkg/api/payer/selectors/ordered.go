package selectors

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

type OrderedPreferenceNodeSelectorAlgorithm struct {
	reg              registry.NodeRegistry
	preferredNodeIDs []uint32
}

func NewOrderedPreferenceNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
	preferredNodeIDs []uint32,
) *OrderedPreferenceNodeSelectorAlgorithm {
	return &OrderedPreferenceNodeSelectorAlgorithm{
		reg:              reg,
		preferredNodeIDs: preferredNodeIDs,
	}
}

func (o *OrderedPreferenceNodeSelectorAlgorithm) GetNode(
	_ topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	nodes, err := o.reg.GetNodes()
	if err != nil {
		return 0, err
	}

	if len(nodes) == 0 {
		return 0, errors.New("no available nodes")
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

	for _, nodeID := range o.preferredNodeIDs {
		if _, exists := nodeMap[nodeID]; !exists {
			continue
		}
		if _, isBanned := banned[nodeID]; !isBanned {
			return nodeID, nil
		}
	}

	for _, node := range nodes {
		if _, isBanned := banned[node.NodeID]; !isBanned {
			return node.NodeID, nil
		}
	}

	return 0, errors.New("no available nodes after considering banlist")
}

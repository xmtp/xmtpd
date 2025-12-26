package selectors

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

type RandomNodeSelectorAlgorithm struct {
	reg registry.NodeRegistry
}

func NewRandomNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
) *RandomNodeSelectorAlgorithm {
	return &RandomNodeSelectorAlgorithm{reg: reg}
}

func (r *RandomNodeSelectorAlgorithm) GetNode(
	_ topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	nodes, err := r.reg.GetNodes()
	if err != nil {
		return 0, err
	}

	if len(nodes) == 0 {
		return 0, errors.New("no available nodes")
	}

	banned := make(map[uint32]struct{})
	for _, list := range banlist {
		for _, id := range list {
			banned[id] = struct{}{}
		}
	}

	availableNodes := make([]uint32, 0, len(nodes))
	for _, node := range nodes {
		if _, isBanned := banned[node.NodeID]; !isBanned {
			availableNodes = append(availableNodes, node.NodeID)
		}
	}

	if len(availableNodes) == 0 {
		return 0, errors.New("no available nodes after considering banlist")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(availableNodes))))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}

	return availableNodes[n.Int64()], nil
}

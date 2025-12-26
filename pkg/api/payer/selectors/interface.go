package selectors

import (
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

type NodeSelectorAlgorithm interface {
	GetNode(topic topic.Topic, banlist ...[]uint32) (uint32, error)
}

type NodeSelectorStrategy string

const (
	NodeSelectorStrategyStable  NodeSelectorStrategy = "stable"
	NodeSelectorStrategyManual  NodeSelectorStrategy = "manual"
	NodeSelectorStrategyOrdered NodeSelectorStrategy = "ordered"
	NodeSelectorStrategyRandom  NodeSelectorStrategy = "random"
	NodeSelectorStrategyClosest NodeSelectorStrategy = "closest"
)

type NodeSelectorConfig struct {
	Strategy       NodeSelectorStrategy
	PreferredNodes []uint32
	CacheExpiry    time.Duration
	ConnectTimeout time.Duration
}

func NewNodeSelector(
	reg registry.NodeRegistry,
	config NodeSelectorConfig,
) (NodeSelectorAlgorithm, error) {
	switch config.Strategy {
	case NodeSelectorStrategyStable, "":
		return NewStableHashingNodeSelectorAlgorithm(reg), nil
	case NodeSelectorStrategyManual:
		if len(config.PreferredNodes) == 0 {
			return nil, errors.New("manual strategy requires at least one preferred node")
		}
		return NewManualNodeSelectorAlgorithm(reg, config.PreferredNodes), nil
	case NodeSelectorStrategyOrdered:
		if len(config.PreferredNodes) == 0 {
			return nil, errors.New("ordered strategy requires at least one preferred node")
		}
		return NewOrderedPreferenceNodeSelectorAlgorithm(reg, config.PreferredNodes), nil
	case NodeSelectorStrategyRandom:
		return NewRandomNodeSelectorAlgorithm(reg), nil
	case NodeSelectorStrategyClosest:
		return NewClosestNodeSelectorAlgorithm(
			reg,
			config.CacheExpiry,
			config.ConnectTimeout,
			config.PreferredNodes,
		), nil
	default:
		return nil, fmt.Errorf("unknown node selector strategy: %s", config.Strategy)
	}
}

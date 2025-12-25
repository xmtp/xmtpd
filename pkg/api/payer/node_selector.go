package payer

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"sort"
	"sync"
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

type StableHashingNodeSelectorAlgorithm struct {
	reg registry.NodeRegistry
}

func NewStableHashingNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
) *StableHashingNodeSelectorAlgorithm {
	return &StableHashingNodeSelectorAlgorithm{reg: reg}
}

// HashKey hashes the topic to a stable uint32 hash using SHA-256.
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
	nodeLocations := make([]uint32, numNodes)
	for i := range nodes {
		nodeLocations[i] = uint32(i) * spacing
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

type ManualNodeSelectorAlgorithm struct {
	reg     registry.NodeRegistry
	nodeIDs []uint32
}

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

type ClosestNodeSelectorAlgorithm struct {
	reg            registry.NodeRegistry
	preferredNodes []uint32
	latencyCache   map[uint32]time.Duration
	cacheMutex     sync.RWMutex
	cacheExpiry    time.Duration
	lastUpdate     time.Time
	connectTimeout time.Duration
}

func NewClosestNodeSelectorAlgorithm(
	reg registry.NodeRegistry,
	cacheExpiry time.Duration,
	connectTimeout time.Duration,
	preferredNodes ...[]uint32,
) *ClosestNodeSelectorAlgorithm {
	if cacheExpiry == 0 {
		cacheExpiry = 5 * time.Minute
	}
	if connectTimeout == 0 {
		connectTimeout = 2 * time.Second
	}
	
	var nodes []uint32
	if len(preferredNodes) > 0 && len(preferredNodes[0]) > 0 {
		nodes = preferredNodes[0]
	}
	
	return &ClosestNodeSelectorAlgorithm{
		reg:            reg,
		preferredNodes: nodes,
		latencyCache:   make(map[uint32]time.Duration),
		cacheExpiry:    cacheExpiry,
		connectTimeout: connectTimeout,
	}
}

func (c *ClosestNodeSelectorAlgorithm) GetNode(
	_ topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	nodes, err := c.reg.GetNodes()
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

	// Filter nodes to preferred list if specified
	nodesToConsider := nodes
	if len(c.preferredNodes) > 0 {
		preferredSet := make(map[uint32]struct{})
		for _, nodeID := range c.preferredNodes {
			preferredSet[nodeID] = struct{}{}
		}

		filtered := make([]registry.Node, 0, len(nodes))
		for _, node := range nodes {
			if _, isPreferred := preferredSet[node.NodeID]; isPreferred {
				filtered = append(filtered, node)
			}
		}

		// If we have preferred nodes available, use only those
		// Otherwise fall back to all nodes
		if len(filtered) > 0 {
			nodesToConsider = filtered
		}
	}

	c.cacheMutex.RLock()
	cacheExpired := time.Since(c.lastUpdate) > c.cacheExpiry
	c.cacheMutex.RUnlock()

	if cacheExpired {
		c.updateLatencyCache(nodesToConsider)
	}

	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	var closestNodeID uint32
	minLatency := time.Duration(1<<63 - 1)

	for _, node := range nodesToConsider {
		if _, isBanned := banned[node.NodeID]; isBanned {
			continue
		}

		latency, ok := c.latencyCache[node.NodeID]
		if !ok {
			continue
		}

		if latency < minLatency {
			minLatency = latency
			closestNodeID = node.NodeID
		}
	}

	if closestNodeID == 0 {
		return 0, errors.New("no available nodes with latency measurements")
	}

	return closestNodeID, nil
}

func (c *ClosestNodeSelectorAlgorithm) updateLatencyCache(nodes []registry.Node) {
	newCache := make(map[uint32]time.Duration)

	for _, node := range nodes {
		latency := c.measureLatency(node.HTTPAddress)
		if latency > 0 {
			newCache[node.NodeID] = latency
		}
	}

	// Only update cache if at least one latency measurement was successful
	// This prevents wiping out the previous cache when all probes fail
	if len(newCache) > 0 {
		c.cacheMutex.Lock()
		c.latencyCache = newCache
		c.lastUpdate = time.Now()
		c.cacheMutex.Unlock()
	}
}

func (c *ClosestNodeSelectorAlgorithm) measureLatency(httpAddress string) time.Duration {
	parsedURL, err := url.Parse(httpAddress)
	if err != nil {
		return -1
	}

	host := parsedURL.Hostname()
	if host == "" {
		return -1
	}

	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	address := net.JoinHostPort(host, port)

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, c.connectTimeout)
	latency := time.Since(start)

	if err != nil {
		return -1
	}
	_ = conn.Close()

	return latency
}

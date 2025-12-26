package selectors

import (
	"errors"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// ClosestNodeSelectorAlgorithm selects the node with the lowest measured TCP connect latency,
// excluding any nodes present in the banlist. Latency measurements are cached for a configured duration
// and refreshed periodically to reduce overhead.
type ClosestNodeSelectorAlgorithm struct {
	reg            registry.NodeRegistry
	preferredNodes []uint32
	latencyCache   map[uint32]time.Duration
	cacheMutex     sync.RWMutex
	cacheExpiry    time.Duration
	lastUpdate     time.Time
	connectTimeout time.Duration
}

var _ NodeSelectorAlgorithm = (*ClosestNodeSelectorAlgorithm)(nil)

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
		latency, err := c.measureLatency(node.HTTPAddress)
		if err == nil && latency > 0 {
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

func (c *ClosestNodeSelectorAlgorithm) measureLatency(httpAddress string) (time.Duration, error) {
	parsedURL, err := url.Parse(httpAddress)
	if err != nil {
		return 0, errors.New("invalid HTTP address")
	}

	host := parsedURL.Hostname()
	if host == "" {
		return 0, errors.New("invalid HTTP address")
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
	if err != nil {
		return 0, err
	}
	_ = conn.Close()

	latency := time.Since(start)

	return latency, nil
}

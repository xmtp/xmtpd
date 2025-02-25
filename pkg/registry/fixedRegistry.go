package registry

import (
	"errors"
	"sync"
)

// TODO: Delete this or move to a test file
type FixedNodeRegistry struct {
	nodes                     []Node
	newNodeNotifier           *notifier[[]Node]
	changedNodeNotifiers      map[uint32]*notifier[Node]
	changedNodeNotifiersMutex sync.Mutex
}

var _ NodeRegistry = &FixedNodeRegistry{}

func NewFixedNodeRegistry(nodes []Node) *FixedNodeRegistry {
	return &FixedNodeRegistry{nodes: nodes}
}

func (r *FixedNodeRegistry) GetNodes() ([]Node, error) {
	return r.nodes, nil
}

func (r *FixedNodeRegistry) GetNode(nodeId uint32) (*Node, error) {
	for _, node := range r.nodes {
		if node.NodeID == nodeId {
			return &node, nil
		}
	}
	return nil, errors.New("node not found")
}

func (f *FixedNodeRegistry) AddNode(node Node) {
	f.nodes = append(f.nodes, node)
}

func (f *FixedNodeRegistry) OnNewNodes() (<-chan []Node, CancelSubscription) {
	return f.newNodeNotifier.register()
}

func (f *FixedNodeRegistry) OnChangedNode(
	nodeId uint32,
) (<-chan Node, CancelSubscription) {
	f.changedNodeNotifiersMutex.Lock()
	defer f.changedNodeNotifiersMutex.Unlock()

	registry, ok := f.changedNodeNotifiers[nodeId]
	if !ok {
		registry = newNotifier[Node]()
		f.changedNodeNotifiers[nodeId] = registry
	}
	return registry.register()
}

func (f *FixedNodeRegistry) Stop() {
}

package registry

import "sync"

// TODO: Delete this or move to a test file
type FixedNodeRegistry struct {
	nodes                     []Node
	newNodeNotifier           *notifier[[]Node]
	changedNodeNotifiers      map[uint16]*notifier[Node]
	changedNodeNotifiersMutex sync.Mutex
}

func NewFixedNodeRegistry(nodes []Node) *FixedNodeRegistry {
	return &FixedNodeRegistry{nodes: nodes}
}

func (r *FixedNodeRegistry) GetNodes() ([]Node, error) {
	return r.nodes, nil
}

func (f *FixedNodeRegistry) AddNode(node Node) {
	f.nodes = append(f.nodes, node)
}

func (f *FixedNodeRegistry) OnNewNodes() (<-chan []Node, CancelSubscription) {
	return f.newNodeNotifier.register()
}

func (f *FixedNodeRegistry) OnChangedNode(
	nodeId uint16,
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

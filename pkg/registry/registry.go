package registry

type Node struct {
	NodeId      int
	SigningKey  []byte
	HttpAddress string
	MtlsCert    []byte
}

type NodeRegistry interface {
	GetNodes() ([]Node, error)
	// OnChange()
}

// TODO: Delete this or move to a test file

type FixedNodeRegistry struct {
	nodes []Node
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

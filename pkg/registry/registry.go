package registry

type Record struct {
	ID            int32
	PublicKey     []byte
	GrpcAddress   string
	DisabledBlock *uint64
	// Maybe add mTLS cert here
}

type NodeRegistry interface {
	GetNodes() ([]Record, error)
	// OnChange()
}

// TODO: Delete this or move to a test file

type FixedNodeRegistry struct {
	nodes []Record
}

func NewFixedNodeRegistry(nodes []Record) *FixedNodeRegistry {
	return &FixedNodeRegistry{nodes: nodes}
}

func (r *FixedNodeRegistry) GetNodes() ([]Record, error) {
	return r.nodes, nil
}

func (f *FixedNodeRegistry) AddNode(node Record) {
	f.nodes = append(f.nodes, node)
}

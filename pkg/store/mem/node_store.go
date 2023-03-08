package memstore

import (
	"sync"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
)

type NodeStore struct {
	ctx        context.Context
	topicsLock sync.RWMutex
	topics     map[string]*memstore.MemoryStore
}

func NewNodeStore(ctx context.Context) *NodeStore {
	return &NodeStore{
		ctx:    ctx,
		topics: make(map[string]*memstore.MemoryStore)}
}

func (n *NodeStore) NewTopic(topic string) (crdt.Store, error) {
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	if t, ok := n.topics[topic]; ok {
		return t, nil
	}
	t := memstore.New(n.ctx)
	n.topics[topic] = t
	return t, nil
}

func (n *NodeStore) Topics() (topics []string, err error) {
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	for name := range n.topics {
		topics = append(topics, name)
	}
	return topics, nil
}

func (n *NodeStore) Close() error {
	return nil
}

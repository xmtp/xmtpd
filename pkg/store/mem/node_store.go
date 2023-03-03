package memstore

import (
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
)

type NodeStore struct {
	ctx context.Context
}

func NewNodeStore(ctx context.Context) *NodeStore {
	return &NodeStore{ctx}
}

func (n *NodeStore) NewTopic(topic string) (crdt.Store, error) {
	return memstore.New(n.ctx), nil
}

func (n *NodeStore) Close() error {
	return nil
}

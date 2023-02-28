package memstore

import (
	"github.com/xmtp/xmtpd/pkg/crdt"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type NodeStore struct {
	log *zap.Logger
}

func NewNodeStore(log *zap.Logger) *NodeStore {
	return &NodeStore{log}
}

func (n *NodeStore) NewTopic(topic string) (crdt.Store, error) {
	return memstore.New(n.log), nil
}

func (n *NodeStore) Close() error {
	return nil
}

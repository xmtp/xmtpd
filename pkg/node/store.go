package node

import "github.com/xmtp/xmtpd/pkg/crdt"

type NodeStore interface {
	NewTopic(topic string) (crdt.Store, error)
	Close() error
}

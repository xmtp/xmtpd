package node

import "github.com/xmtp/xmtpd/pkg/crdt"

type NodeStore interface {
	// Open or create a topic store in the node store
	NewTopic(topic string) (crdt.Store, error)
	// Return list of names for all topic stores in the node store
	Topics() ([]string, error)
	// Close the node store
	Close() error
	// Remove topic from the store.
	// Assumes the replica for the topic has been closed.
	DeleteTopic(topic string) error
}

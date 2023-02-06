package crdt

import (
	mh "github.com/multiformats/go-multihash"
)

// NodeSyncer manages the syncing capability for a Node
type NodeSyncer interface {
	// AddNode registers the node with the broadcaster.
	AddNode(*Node)
	// NewTopic creates a TopicSyncer for given topic and node.
	NewTopic(name string, node *Node) TopicSyncer
}

// TopicSyncer provides syncing capability to a specific topic.
type TopicSyncer interface {
	// Fetch retrieves a set of Events from the network based on the provided CIDs.
	// It is a single attempt that can fail completely or return only some
	// of the requested events. If there is no error, the resulting slice is always
	// the same size as the CID slice, but there can be some nils instead of Events in it.
	Fetch([]mh.Multihash) ([]*Event, error)
	// FetchAll() ([]*Event, error) // used to seed new nodes
}

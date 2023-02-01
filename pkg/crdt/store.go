package crdt

import (
	"context"
	"errors"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
)

var InvalidCursor = errors.New("Invalid cursor")

// NodeStore manages the storage capacity for a Node.
type NodeStore interface {
	// NewTopic creates a TopicStore for the specified topic.
	NewTopic(name string, node *Node) TopicStore

	// Following methods are needed for bootstrapping topics
	// from a pre-existing store.

	// Topics returns the list of all the stored topic names
	Topics() (topics []string, err error)
}

// TopicStore represents the storage capacity for a specific topic.
type TopicStore interface {
	// NewEvent creates and stores a new Event,
	// making the current heads its links and
	// replacing the heads with the new Event.
	// Returns the new Event.
	NewEvent(*messagev1.Envelope) (*Event, error)
	// AddEvent stores the Event if it isn't know yet,
	// Returns whether it was actually added.
	AddEvent(ev *Event) (added bool, err error)
	// AddHead stores the Event if it isn't know yet,
	// and add it to the heads
	// Returns whether it was actually added.
	AddHead(ev *Event) (added bool, err error)
	// RemoveHead checks if we already have the event,
	// and also removes it from heads if it's there.
	// Returns whether we already have the event or not.
	RemoveHead(cid mh.Multihash) (haveAlready bool, err error)

	// Following methods are needed for bootstrapping a topic
	// from a pre-existing store.

	// GetPendingLinks scans the whole topic for links that
	// are not present in the topic.
	// Returns the list of all missing links.
	FindMissingLinks() (links []mh.Multihash, err error)

	// Following methods are needed for querying
	Query(ctx context.Context, req *messagev1.QueryRequest) ([]*messagev1.Envelope, *messagev1.PagingInfo, error)

	// Following methods are just for testing,
	// not needed for the protocol implementation

	// Get returns the Event based on its CID, nil if absent.
	Get(cid mh.Multihash) (*Event, error)
	// Count returns count of all stored events
	Count() (int, error)
}

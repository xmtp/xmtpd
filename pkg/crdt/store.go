package crdt

import (
	"errors"

	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

// The cursor event couldn't be find in the event range resulting from the parameters of the query.
// The store may have changed since the query that yielded that cursor.
var ErrStoreInvalidCursor = errors.New("cursor event not found")

// Store represents the storage capacity for a specific CRDT.
type Store interface {
	// AppendEvent creates and stores a new Event,
	// making the current heads its links and
	// replacing the heads with the new Event.
	// Returns the new Event.
	AppendEvent([]byte) (*types.Event, error)
	// AddEvent stores the Event if it isn't know yet,
	// Returns whether it was actually added.
	AddEvent(ev *types.Event) (added bool, err error)
	// AddHead stores the Event if it isn't know yet,
	// and add it to the heads
	// Returns whether it was actually added.
	AddHead(ev *types.Event) (added bool, err error)
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
}

package crdt

import (
	"context"
	"errors"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
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
	AppendEvent(ctx context.Context, env *messagev1.Envelope) (*types.Event, error)

	// AddEvent stores the Event if it isn't know yet,
	// Returns whether it was actually added.
	InsertEvent(ctx context.Context, ev *types.Event) (bool, error)

	// AddHead stores the Event if it isn't know yet,
	// and add it to the heads
	// Returns whether it was actually added.
	InsertHead(ctx context.Context, ev *types.Event) (bool, error)

	// RemoveHead checks if we already have the event,
	// and also removes it from heads if it's there.
	// Returns whether we already have the event or not.
	RemoveHead(ctx context.Context, cid multihash.Multihash) (bool, error)

	// GetEvents returns the set of events matching the given set of CIDs.
	GetEvents(ctx context.Context, links []multihash.Multihash) ([]*types.Event, error)

	// Following methods are needed for bootstrapping a topic
	// from a pre-existing store.

	// GetPendingLinks scans the whole topic for links that
	// are not present in the topic.
	// Returns the list of all missing links.
	FindMissingLinks(ctx context.Context) ([]multihash.Multihash, error)

	// Query returns a set of envelopes matching the query request criteria.
	Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error)

	// Close gracefully closes the store.
	Close() error
}

package crdt

import (
	"context"

	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

// Syncer provides syncing capability to a specific replica.
type Syncer interface {
	// Fetch retrieves a set of Events from the network DHT.
	// It is a single attempt that can fail completely or return only some
	// of the requested events. If there is no error, the resulting slice is always
	// the same size as the CID slice, but there can be some nils instead of Events in it.
	Fetch(context.Context, []multihash.Multihash) ([]*types.Event, error)

	// Close gracefully closes the syncer.
	Close() error
}

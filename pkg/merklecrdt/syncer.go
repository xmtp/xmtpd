package merklecrdt

import (
	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
)

// Syncer provides syncing capability to a specific CRDT.
type Syncer interface {
	// Fetch retrieves a set of Events from the network based on the provided CIDs.
	// It is a single attempt that can fail completely or return only some
	// of the requested events. If there is no error, the resulting slice is always
	// the same size as the CID slice, but there can be some nils instead of Events in it.
	Fetch([]mh.Multihash) ([]*types.Event, error)
	// FetchAll() ([]*Event, error) // used to seed new nodes
}

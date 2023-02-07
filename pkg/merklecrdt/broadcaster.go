package merklecrdt

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
)

// Broadcaster manages broadcasts for a CRDT instance.
type Broadcaster interface {
	// Broadcast sends an Event out to the network
	Broadcast(*types.Event) error

	// Obtain the next event received from the network.
	Next(ctx context.Context) (*types.Event, error)
}

package crdt

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

// Broadcaster manages broadcasts for a replica.
type Broadcaster interface {
	// Broadcast sends an Event out to the network
	Broadcast(*types.Event) error

	// Next obtains the next event received from the network.
	Next(ctx context.Context) (*types.Event, error)

	// Close gracefully closes the broadcaster.
	Close() error
}

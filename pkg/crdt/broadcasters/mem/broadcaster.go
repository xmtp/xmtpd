package membroadcaster

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemoryBroadcaster struct {
	log   *zap.Logger
	ch    chan *types.Event
	peers []*MemoryBroadcaster
}

func New(log *zap.Logger) *MemoryBroadcaster {
	return &MemoryBroadcaster{
		log: log,
		ch:  make(chan *types.Event, 100),
	}
}

func (b *MemoryBroadcaster) Broadcast(ev *types.Event) error {
	b.log.Debug("broadcast event", zap.Cid("event", ev.Cid))
	b.ch <- ev
	for _, peer := range b.peers {
		peer.ch <- ev
	}
	return nil
}

func (b *MemoryBroadcaster) Next(ctx context.Context) (*types.Event, error) {
	return <-b.ch, nil
}

func (b *MemoryBroadcaster) AddPeer(peer *MemoryBroadcaster) {
	b.peers = append(b.peers, peer)
}

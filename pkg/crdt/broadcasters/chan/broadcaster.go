package chanbroadcaster

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type ChannelBroadcaster struct {
	log   *zap.Logger
	ch    chan *types.Event
	peers []*ChannelBroadcaster
}

func New(log *zap.Logger) *ChannelBroadcaster {
	return &ChannelBroadcaster{
		log: log,
		ch:  make(chan *types.Event, 100),
	}
}

func (b *ChannelBroadcaster) Broadcast(ev *types.Event) error {
	b.log.Debug("broadcast event", zap.Cid("event", ev.Cid))
	b.ch <- ev
	for _, peer := range b.peers {
		peer.ch <- ev
	}
	return nil
}

func (b *ChannelBroadcaster) Next(ctx context.Context) (*types.Event, error) {
	return <-b.ch, nil
}

func (b *ChannelBroadcaster) AddPeer(peer *ChannelBroadcaster) {
	b.peers = append(b.peers, peer)
}

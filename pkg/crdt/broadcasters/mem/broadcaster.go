package membroadcaster

import (
	"reflect"
	"testing"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemoryBroadcaster struct {
	log   *zap.Logger
	ch    chan *types.Event
	peers []*MemoryBroadcaster
}

func New(ctx context.Context) *MemoryBroadcaster {
	return &MemoryBroadcaster{
		log: ctx.Logger(),
		ch:  make(chan *types.Event, 100),
	}
}

func (b *MemoryBroadcaster) Close() error {
	close(b.ch)
	return nil
}

func (b *MemoryBroadcaster) Broadcast(ctx context.Context, ev *types.Event) error {
	b.log.Debug("broadcast event", zap.Cid("event", ev.Cid))
	b.ch <- ev
	for _, peer := range b.peers {
		peer.ch <- ev
	}
	return nil
}

func (b *MemoryBroadcaster) Next(ctx context.Context) (*types.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ev, ok := <-b.ch:
		if !ok {
			return nil, nil
		}
		b.log.Debug("received broadcasted event", zap.Cid("event", ev.Cid))
		return ev, nil
	}
}

// AddPeer adds a memory broadcaster peer that new events will be shared with.
func (b *MemoryBroadcaster) AddPeer(t *testing.T, peer interface{}) {
	switch peer := peer.(type) {
	case *MemoryBroadcaster:
		b.peers = append(b.peers, peer)
	default:
		b.log.Warn("unknown broadcaster peer type", zap.String("type", reflect.TypeOf(peer).String()))
	}
}

package memsyncer

import (
	"math/rand"
	"reflect"

	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemorySyncer struct {
	log   *zap.Logger
	store crdt.Store
	peers []*MemorySyncer
}

func New(ctx context.Context, store crdt.Store) *MemorySyncer {
	return &MemorySyncer{
		log:   ctx.Logger(),
		store: store,
	}
}

func (s *MemorySyncer) Close() error {
	return nil
}

func (s *MemorySyncer) Fetch(ctx context.Context, cids []multihash.Multihash) ([]*types.Event, error) {
	peer := s.randomPeer()
	if peer == nil {
		return []*types.Event{}, nil
	}
	events, err := peer.GetStoreEvents(ctx, cids)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *MemorySyncer) AddStoreEvents(ctx context.Context, events []*types.Event) error {
	for _, ev := range events {
		_, err := s.store.InsertEvent(ctx, ev)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MemorySyncer) AddPeer(peer interface{}) {
	switch peer := peer.(type) {
	case *MemorySyncer:
		s.peers = append(s.peers, peer)
	default:
		s.log.Warn("unknown syncer peer type", zap.String("type", reflect.TypeOf(peer).String()))
	}
}

func (s *MemorySyncer) GetStoreEvents(ctx context.Context, cids []multihash.Multihash) ([]*types.Event, error) {
	return s.store.GetEvents(ctx, cids...)
}

func (s *MemorySyncer) randomPeer() *MemorySyncer {
	if len(s.peers) == 0 {
		return nil
	}
	i := rand.Intn(len(s.peers))
	return s.peers[i]
}

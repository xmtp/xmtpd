package memsyncer

import (
	"context"
	"math/rand"
	"reflect"

	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemorySyncer struct {
	log   *zap.Logger
	store crdt.Store
	peers []*MemorySyncer
}

func New(log *zap.Logger, store crdt.Store) *MemorySyncer {
	return &MemorySyncer{
		log:   log,
		store: store,
	}
}

func (s *MemorySyncer) Close() error {
	return nil
}

func (s *MemorySyncer) Fetch(ctx context.Context, cids []multihash.Multihash) ([]*types.Event, error) {
	localEvents, err := s.store.GetEvents(ctx, cids)
	if err != nil {
		return nil, err
	}
	localCids := map[string]struct{}{}
	for _, ev := range localEvents {
		localCids[ev.Cid.String()] = struct{}{}
	}

	missingCids := make([]multihash.Multihash, 0, len(cids))
	for _, cid := range cids {
		if _, ok := localCids[cid.String()]; ok {
			continue
		}
		missingCids = append(missingCids, cid)
	}
	if len(missingCids) == 0 {
		return localEvents, nil
	}

	peer := s.randomPeer()
	if peer == nil {
		return localEvents, nil
	}
	peerEvents, err := peer.Fetch(ctx, missingCids)
	if err != nil {
		return nil, err
	}
	events := append(localEvents, peerEvents...)

	return events, nil
}

func (s *MemorySyncer) AddStoreEvents(ctx context.Context, events []*types.Event) error {
	for _, ev := range events {
		_, err := s.store.AddEvent(ev)
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

func (s *MemorySyncer) randomPeer() *MemorySyncer {
	if len(s.peers) == 0 {
		return nil
	}
	i := rand.Intn(len(s.peers))
	return s.peers[i]
}

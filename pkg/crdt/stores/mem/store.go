package memstore

import (
	"bytes"
	"context"
	"sort"
	"sync"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// MemoryStore is an in-memory store used for testing.
type MemoryStore struct {
	sync.RWMutex

	log *zap.Logger

	heads        map[string]bool         // CIDs of current head events
	events       map[string]*types.Event // maps CIDs to all known Events
	eventsByTime []*types.Event
}

func New(log *zap.Logger) *MemoryStore {
	return &MemoryStore{
		log:    log,
		heads:  make(map[string]bool),
		events: make(map[string]*types.Event),
	}
}

func (s *MemoryStore) Close() error {
	return nil
}

func (s *MemoryStore) InsertEvent(ctx context.Context, ev *types.Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.Cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.log.Debug("inserting event", zap.Cid("event", ev.Cid))
	s.addEvent(key, ev)
	return true, nil
}

func (s *MemoryStore) AppendEvent(ctx context.Context, env *messagev1.Envelope) (*types.Event, error) {
	s.Lock()
	defer s.Unlock()
	heads, err := s.Heads(ctx)
	if err != nil {
		return nil, err
	}
	ev, err := types.NewEvent(env, heads)
	if err != nil {
		return nil, err
	}
	key := ev.Cid.String()
	s.log.Debug("appending event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))
	s.addEvent(key, ev)
	s.heads = map[string]bool{key: true}
	return ev, err
}

func (s *MemoryStore) InsertHead(ctx context.Context, ev *types.Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.Cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.addEvent(key, ev)
	s.heads[key] = true
	s.log.Debug("inserting head", zap.Cid("event", ev.Cid), zap.Int("heads", len(s.heads)))
	return true, nil
}

func (s *MemoryStore) RemoveHead(ctx context.Context, cid multihash.Multihash) (have bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := cid.String()
	if s.events[key] == nil {
		return false, nil
	}
	if s.heads[key] {
		s.log.Debug("removing head", zap.Cid("event", cid), zap.Int("heads", len(s.heads)-1))
	}
	delete(s.heads, key)
	return true, nil
}

func (s *MemoryStore) FindMissingLinks(ctx context.Context) (links []multihash.Multihash, err error) {
	s.RLock()
	defer s.RUnlock()
	for _, ev := range s.events {
		for _, cid := range ev.Links {
			if s.events[cid.String()] == nil {
				links = append(links, cid)
			}
		}
	}
	return links, nil
}

func (s *MemoryStore) GetEvents(ctx context.Context, cids ...multihash.Multihash) ([]*types.Event, error) {
	s.RLock()
	defer s.RUnlock()
	events := make([]*types.Event, 0, len(cids))
	for _, cid := range cids {
		ev, ok := s.events[cid.String()]
		if !ok {
			continue
		}
		events = append(events, ev)
	}
	return events, nil
}

func (s *MemoryStore) Events(ctx context.Context) ([]*types.Event, error) {
	s.RLock()
	defer s.RUnlock()
	events := make([]*types.Event, len(s.events))
	i := 0
	for _, ev := range s.events {
		events[i] = ev
		i++
	}
	return events, nil
}

func (s *MemoryStore) Heads(ctx context.Context) ([]multihash.Multihash, error) {
	cids := []multihash.Multihash{}
	for key := range s.heads {
		cids = append(cids, s.events[key].Cid)
	}
	return cids, nil
}

// private functions

// key MUST be equal to ev.Cid.String()
func (s *MemoryStore) addEvent(key string, ev *types.Event) {
	s.events[key] = ev

	// Add to index sorted by timestamp.
	i, _ := sort.Find(len(s.eventsByTime), func(i int) int {
		res := ev.TimestampNs - s.eventsByTime[i].TimestampNs
		if res != 0 {
			return int(res)
		}
		return bytes.Compare(ev.Cid, s.eventsByTime[i].Cid)
	})
	if i == len(s.eventsByTime) {
		s.eventsByTime = append(s.eventsByTime, ev)
	} else {
		s.eventsByTime = makeRoomAt(s.eventsByTime, i)
	}
	s.eventsByTime[i] = ev
}

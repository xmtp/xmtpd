package memstore

import (
	"sync"

	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// MemoryStore is an in-memory store used for testing.
type MemoryStore struct {
	sync.RWMutex

	log *zap.Logger

	heads  map[string]bool         // CIDs of current head events
	events map[string]*types.Event // maps CIDs to all known Events
}

func New(log *zap.Logger) *MemoryStore {
	return &MemoryStore{
		log:    log,
		heads:  make(map[string]bool),
		events: make(map[string]*types.Event),
	}
}

func (s *MemoryStore) AddEvent(ev *types.Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.Cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.log.Debug("adding event", zap.Cid("event", ev.Cid))
	s.addEvent(key, ev)
	return true, nil
}

func (s *MemoryStore) AddHead(ev *types.Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.Cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.addEvent(key, ev)
	s.heads[key] = true
	s.log.Debug("adding head", zap.Cid("event", ev.Cid), zap.Int("heads", len(s.heads)))
	return true, nil
}

func (s *MemoryStore) RemoveHead(cid multihash.Multihash) (have bool, err error) {
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

func (s *MemoryStore) NewEvent(payload []byte) (*types.Event, error) {
	s.Lock()
	defer s.Unlock()
	ev, err := types.NewEvent(payload, s.allHeads())
	if err != nil {
		return nil, err
	}
	key := ev.Cid.String()
	s.log.Debug("creating event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))
	s.addEvent(key, ev)
	s.heads = map[string]bool{key: true}
	return ev, err
}

func (s *MemoryStore) FindMissingLinks() (links []multihash.Multihash, err error) {
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

func (s *MemoryStore) Events() map[string]*types.Event {
	return s.events
}

// private functions

func (s *MemoryStore) allHeads() (cids []multihash.Multihash) {
	for key := range s.heads {
		cids = append(cids, s.events[key].Cid)
	}
	return cids
}

// key MUST be equal to ev.Cid.String()
func (s *MemoryStore) addEvent(key string, ev *types.Event) {
	s.events[key] = ev
}

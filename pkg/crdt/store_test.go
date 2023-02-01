package crdt

import (
	"context"
	"sort"
	"sync"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"go.uber.org/zap"
)

func Test_Query(t *testing.T) {
	// create a topic with some pre-existing traffic
	net := randomMsgTest(t, 1, 1, 20)
	defer net.Close()

	res, _, err := net.Query(0, t0, timeRange(5, 13))
	assert.NoError(t, err)
	net.assertQueryResult(res, 5, 6, 7, 8, 9, 10, 11, 12, 13)

	res, _, err = net.Query(0, t0, timeRange(5, 9), descending())
	assert.NoError(t, err)
	net.assertQueryResult(res, 9, 8, 7, 6, 5)

	res, _, err = net.Query(0, t0, timeRange(5, 15), limit(4))
	assert.NoError(t, err)
	net.assertQueryResult(res, 5, 6, 7, 8)

	res, _, err = net.Query(0, t0, timeRange(5, 15), limit(4), descending())
	assert.NoError(t, err)
	net.assertQueryResult(res, 15, 14, 13, 12)
}

// In-memory store using maps to store Events
type mapStore struct {
	sync.RWMutex
	topics map[string]*mapTopicStore
}

func newMapStore() *mapStore {
	return &mapStore{
		topics: make(map[string]*mapTopicStore),
	}
}

// NewTopic returns a store for a pre-existing topic or creates a new one.
func (s *mapStore) NewTopic(name string, n *Node) TopicStore {
	s.Lock()
	defer s.Unlock()
	ts := s.topics[name]
	if ts == nil {
		ts = &mapTopicStore{
			node:   n,
			log:    n.log.Named(name),
			heads:  make(map[string]bool),
			events: make(map[string]*Event),
		}
		s.topics[name] = ts
	}
	return ts
}

// Topics lists all topics in the store.
func (s *mapStore) Topics() (topics []string, err error) {
	s.RLock()
	defer s.RUnlock()
	for k := range s.topics {
		topics = append(topics, k)
	}
	return topics, nil
}

// In-memory TopicStore
type mapTopicStore struct {
	sync.RWMutex
	node   *Node
	heads  map[string]bool   // CIDs of current head events
	events map[string]*Event // maps CIDs to all known Events
	byTime []*Event          // events sorted by event.timestampNs
	log    *zap.Logger
}

var _ TopicStore = (*mapTopicStore)(nil)

func (s *mapTopicStore) AddEvent(ev *Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.log.Debug("adding event", zapCid("event", ev.cid))
	s.addEvent(key, ev)
	return true, nil
}

func (s *mapTopicStore) AddHead(ev *Event) (added bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := ev.cid.String()
	if s.events[key] != nil {
		return false, nil
	}
	s.addEvent(key, ev)
	s.heads[key] = true
	s.log.Debug("adding head", zapCid("event", ev.cid), zap.Int("heads", len(s.heads)))
	return true, nil
}

func (s *mapTopicStore) RemoveHead(cid mh.Multihash) (have bool, err error) {
	s.Lock()
	defer s.Unlock()
	key := cid.String()
	if s.events[key] == nil {
		return false, nil
	}
	if s.heads[key] {
		s.log.Debug("removing head", zapCid("event", cid), zap.Int("heads", len(s.heads)-1))
	}
	delete(s.heads, key)
	return true, nil
}

func (s *mapTopicStore) NewEvent(env *messagev1.Envelope) (*Event, error) {
	s.Lock()
	defer s.Unlock()
	ev, err := NewEvent(env, s.allHeads())
	if err != nil {
		return nil, err
	}
	key := ev.cid.String()
	s.log.Debug("creating event", zapCid("event", ev.cid), zap.Int("links", len(ev.links)))
	s.addEvent(key, ev)
	s.heads = map[string]bool{key: true}
	return ev, err
}

func (s *mapTopicStore) FindMissingLinks() (links []mh.Multihash, err error) {
	s.RLock()
	defer s.RUnlock()
	for _, ev := range s.events {
		for _, cid := range ev.links {
			if s.events[cid.String()] == nil {
				links = append(links, cid)
			}
		}
	}
	return links, nil
}

func (s *mapTopicStore) Query(ctx context.Context, req *messagev1.QueryRequest) ([]*messagev1.Envelope, *messagev1.PagingInfo, error) {
	s.RLock()
	defer s.RUnlock()
	from, _ := sort.Find(len(s.byTime), func(i int) int {
		return int(req.StartTimeNs - s.byTime[i].TimestampNs)
	})
	if from == len(s.byTime) {
		// everything is earlier than StartTimeNs
		return nil, nil, nil
	}
	upTo := req.EndTimeNs + 1
	end, _ := sort.Find(len(s.byTime), func(i int) int {
		return int(upTo - s.byTime[i].TimestampNs)
	})
	result := s.byTime[from:end]
	if req.PagingInfo == nil {
		return toEnvelopes(result, false), nil, nil
	}
	if cursor := req.PagingInfo.Cursor.GetIndex(); cursor != nil {
		return nil, nil, TODO
	}
	if req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING {
		if limit := req.PagingInfo.Limit; limit != 0 {
			result = result[len(result)-int(limit):]
		}
		return toEnvelopes(result, true), nil, nil

	}
	if limit := req.PagingInfo.Limit; limit != 0 {
		result = result[:limit]
	}
	return toEnvelopes(result, false), nil, nil
}

func (s *mapTopicStore) Get(cid mh.Multihash) (*Event, error) {
	s.RLock()
	defer s.RUnlock()
	return s.events[cid.String()], nil
}

func (s *mapTopicStore) Count() (int, error) {
	return len(s.events), nil
}

// private functions

func (s *mapTopicStore) allHeads() (cids []mh.Multihash) {
	for key := range s.heads {
		cids = append(cids, s.events[key].cid)
	}
	return cids
}

// key MUST be equal to ev.cid.String()
func (s *mapTopicStore) addEvent(key string, ev *Event) {
	i, _ := sort.Find(len(s.byTime), func(i int) int {
		return ev.Compare(s.byTime[i])
	})
	if i == len(s.byTime) {
		s.byTime = append(s.byTime, ev)
	} else {
		s.byTime = makeRoomAt(s.byTime, i)
	}
	s.byTime[i] = ev
	s.events[key] = ev
}

// shift events from index i to the right
// to create room at the index.
func makeRoomAt(events []*Event, i int) []*Event {
	// if there's enough capacity in the slice, just shift the tail
	if len(events) < cap(events) {
		events = events[:len(events)+1]
		copy(events[i+1:], events[i:])
		return events
	}
	// figure out desired capacity of a new slice
	var newCap int
	// don't need to worry about len(events) == 0
	// because of the !found append in addEvent
	if len(events) < 1024 {
		newCap = 2 * len(events)
	} else {
		newCap = len(events) + 1024
	}
	// copy events into a new slice, leaving a gap at index i
	newEvents := make([]*Event, len(events)+1, newCap)
	copy(newEvents, events[:i])
	copy(newEvents[i+1:], events[i:])
	return newEvents
}

func toEnvelopes(events []*Event, reversed bool) []*messagev1.Envelope {
	envs := make([]*messagev1.Envelope, len(events))
	if reversed {
		for i, j := 0, len(events)-1; i < len(envs); i, j = i+1, j-1 {
			envs[i] = events[j].Envelope
		}
	} else {
		for i := range envs {
			envs[i] = events[i].Envelope
		}
	}
	return envs
}

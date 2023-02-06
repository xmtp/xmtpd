package tests

import (
	"context"
	"sort"
	"sync"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// In-memory store using maps to store Events
type mapStore struct {
	sync.RWMutex
	topics map[string]*mapTopicStore
}

func NewMapStore(l *zap.Logger) crdt.NodeStore {
	return &mapStore{
		topics: make(map[string]*mapTopicStore),
	}
}

// NewTopic returns a store for a pre-existing topic or creates a new one.
func (s *mapStore) NewTopic(name string, n *crdt.Node) crdt.TopicStore {
	s.Lock()
	defer s.Unlock()
	ts := s.topics[name]
	if ts == nil {
		ts = &mapTopicStore{
			node:   n,
			log:    n.LogNamed(name),
			heads:  make(map[string]bool),
			events: make(map[string]*crdt.Event),
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
	node   *crdt.Node
	heads  map[string]bool        // CIDs of current head events
	events map[string]*crdt.Event // maps CIDs to all known Events
	byTime []*crdt.Event          // events sorted by event.timestampNs
	log    *zap.Logger
}

func (s *mapTopicStore) AddEvent(ev *crdt.Event) (added bool, err error) {
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

func (s *mapTopicStore) AddHead(ev *crdt.Event) (added bool, err error) {
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

func (s *mapTopicStore) RemoveHead(cid mh.Multihash) (have bool, err error) {
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

func (s *mapTopicStore) NewEvent(env *messagev1.Envelope) (*crdt.Event, error) {
	s.Lock()
	defer s.Unlock()
	ev, err := crdt.NewEvent(env, s.allHeads())
	if err != nil {
		return nil, err
	}
	key := ev.Cid.String()
	s.log.Debug("creating event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))
	s.addEvent(key, ev)
	s.heads = map[string]bool{key: true}
	return ev, err
}

func (s *mapTopicStore) FindMissingLinks() (links []mh.Multihash, err error) {
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

func (s *mapTopicStore) Query(ctx context.Context, req *messagev1.QueryRequest) ([]*messagev1.Envelope, *messagev1.PagingInfo, error) {
	s.RLock()
	defer s.RUnlock()
	var start int
	if req.StartTimeNs > 0 {
		start, _ = sort.Find(len(s.byTime), func(i int) int {
			return int(req.StartTimeNs - s.byTime[i].TimestampNs)
		})
	}
	if start == len(s.byTime) {
		// everything is earlier than StartTimeNs
		return nil, nil, nil
	}
	end := len(s.byTime)
	if req.EndTimeNs > 0 {
		upTo := req.EndTimeNs + 1
		end, _ = sort.Find(len(s.byTime), func(i int) int {
			return int(upTo - s.byTime[i].TimestampNs)
		})
	}
	result := s.byTime[start:end]
	if req.PagingInfo == nil {
		// if there's no paging info we're done
		return toEnvelopes(result, false), nil, nil
	}
	reversed := req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING
	cursor := req.PagingInfo.Cursor.GetIndex()
	if cursor != nil {
		// find the cursor event in the result
		cEvt := &crdt.Event{
			Cid:      cursor.Digest,
			Envelope: &messagev1.Envelope{TimestampNs: cursor.SenderTimeNs},
		}
		cIdx, found := sort.Find(len(result), func(i int) int {
			return cEvt.Compare(result[i])
		})
		if !found {
			return nil, nil, crdt.ErrInvalidCursor
		}
		// reslice the result from the cursor event to the end
		if reversed {
			result = result[:cIdx]
		} else {
			result = result[cIdx+1:]
		}
	}
	if reversed {
		if limit := req.PagingInfo.Limit; limit != 0 && int(limit) < len(result) {
			result = result[len(result)-int(limit):]
		}
		var newCursorEvent *crdt.Event
		if len(result) > 0 {
			newCursorEvent = result[0]
		}
		return toEnvelopes(result, reversed), updatedPagingInfo(req.PagingInfo, newCursorEvent), nil

	}
	if limit := req.PagingInfo.Limit; limit != 0 && int(limit) < len(result) {
		result = result[:limit]
	}
	var newCursorEvent *crdt.Event
	if len(result) > 0 {
		newCursorEvent = result[len(result)-1]
	}
	return toEnvelopes(result, reversed), updatedPagingInfo(req.PagingInfo, newCursorEvent), nil
}

func (s *mapTopicStore) Get(cid mh.Multihash) (*crdt.Event, error) {
	s.RLock()
	defer s.RUnlock()
	return s.events[cid.String()], nil
}

func (s *mapTopicStore) Count() (int, error) {
	s.RLock()
	defer s.RUnlock()
	return len(s.events), nil
}

// private functions

func (s *mapTopicStore) allHeads() (cids []mh.Multihash) {
	for key := range s.heads {
		cids = append(cids, s.events[key].Cid)
	}
	return cids
}

// key MUST be equal to ev.Cid.String()
func (s *mapTopicStore) addEvent(key string, ev *crdt.Event) {
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
func makeRoomAt(events []*crdt.Event, i int) []*crdt.Event {
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
	newEvents := make([]*crdt.Event, len(events)+1, newCap)
	copy(newEvents, events[:i])
	copy(newEvents[i+1:], events[i:])
	return newEvents
}

func toEnvelopes(events []*crdt.Event, reversed bool) []*messagev1.Envelope {
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

// updates paging info with a cursor for given event (or nil)
func updatedPagingInfo(pi *messagev1.PagingInfo, cursorEvent *crdt.Event) *messagev1.PagingInfo {
	var cursor *messagev1.Cursor
	if cursorEvent != nil {
		cursor = &messagev1.Cursor{
			Cursor: &messagev1.Cursor_Index{
				Index: &messagev1.IndexCursor{
					SenderTimeNs: cursorEvent.TimestampNs,
					Digest:       cursorEvent.Cid,
				},
			},
		}
	}
	// Note that we're modifying the original query's paging info here.
	pi.Cursor = cursor
	return pi
}

// generate a list of ints in from start to end inclusive.
// if end < start generate it in reverse.
func intRange(start, end int) (list []int) {
	if start < end {
		list = make([]int, end-start+1)
		for k := range list {
			list[k] = start + k
		}
		return list
	}
	list = make([]int, start-end+1)
	for k := range list {
		list[k] = start - k
	}
	return list
}

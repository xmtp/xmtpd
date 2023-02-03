package crdt

import (
	"context"
	"sort"
	"sync"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func Test_Query(t *testing.T) {
	// create a topic with 20 messages
	net := randomMsgTest(t, 1, 1, 20)
	defer net.Close()

	t.Run("all", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0)
		require.NoError(t, err)
		assert.Nil(t, pi, "paging info")
		net.AssertQueryResult(t, res, intRange(1, 20)...)
	})
	t.Run("descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, descending())
		require.NoError(t, err)
		assert.NotNil(t, pi, "paging info")
		net.AssertQueryResult(t, res, intRange(20, 1)...)
	})
	t.Run("limit", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, limit(5))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 5, pi.Cursor)
		net.AssertQueryResult(t, res, intRange(1, 5)...)
	})
	t.Run("limit descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, limit(5), descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 16, pi.Cursor)
		net.AssertQueryResult(t, res, intRange(20, 16)...)
	})
	t.Run("range", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(5, 13))
		require.NoError(t, err)
		assert.Nil(t, pi, "paging info")
		net.AssertQueryResult(t, res, 5, 6, 7, 8, 9, 10, 11, 12, 13)

	})
	t.Run("range descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(5, 9), descending())
		require.NoError(t, err)
		assert.NotNil(t, pi, "paging info")
		net.AssertQueryResult(t, res, 9, 8, 7, 6, 5)

	})
	t.Run("range limit", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(5, 15), limit(4))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 8, pi.Cursor)
		net.AssertQueryResult(t, res, 5, 6, 7, 8)

	})
	t.Run("range limit descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(5, 15), limit(4), descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 12, pi.Cursor)
		net.AssertQueryResult(t, res, 15, 14, 13, 12)

	})
	t.Run("cursor", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(5, 13), limit(5))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 9, pi.Cursor)
		net.AssertQueryResult(t, res, 5, 6, 7, 8, 9)

		res, pi, err = net.Query(t, 0, t0, timeRange(5, 13), limit(5), cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 13, pi.Cursor)
		net.AssertQueryResult(t, res, 10, 11, 12, 13)

		res, pi, err = net.Query(t, 0, t0, timeRange(5, 13), limit(5), cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.Nil(t, pi.Cursor)
		net.AssertQueryResult(t, res)

	})
	t.Run("cursor descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, timeRange(7, 15), limit(5), descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.NotNil(t, pi.Cursor)
		net.AssertQueryResult(t, res, 15, 14, 13, 12, 11)

		res, pi, err = net.Query(t, 0, t0, timeRange(7, 15), limit(5), descending(), cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.NotNil(t, pi.Cursor)
		net.AssertQueryResult(t, res, 10, 9, 8, 7)

		res, pi, err = net.Query(t, 0, t0, timeRange(7, 15), limit(5), descending(), cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.Nil(t, pi.Cursor)
		net.AssertQueryResult(t, res)
	})
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
	s.log.Debug("adding event", zap.Cid("event", ev.cid))
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
	s.log.Debug("adding head", zap.Cid("event", ev.cid), zap.Int("heads", len(s.heads)))
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

func (s *mapTopicStore) NewEvent(env *messagev1.Envelope) (*Event, error) {
	s.Lock()
	defer s.Unlock()
	ev, err := NewEvent(env, s.allHeads())
	if err != nil {
		return nil, err
	}
	key := ev.cid.String()
	s.log.Debug("creating event", zap.Cid("event", ev.cid), zap.Int("links", len(ev.links)))
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
		cEvt := &Event{
			cid:      cursor.Digest,
			Envelope: &messagev1.Envelope{TimestampNs: cursor.SenderTimeNs},
		}
		cIdx, found := sort.Find(len(result), func(i int) int {
			return cEvt.Compare(result[i])
		})
		if !found {
			return nil, nil, ErrInvalidCursor
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
		var newCursorEvent *Event
		if len(result) > 0 {
			newCursorEvent = result[0]
		}
		return toEnvelopes(result, reversed), updatedPagingInfo(req.PagingInfo, newCursorEvent), nil

	}
	if limit := req.PagingInfo.Limit; limit != 0 && int(limit) < len(result) {
		result = result[:limit]
	}
	var newCursorEvent *Event
	if len(result) > 0 {
		newCursorEvent = result[len(result)-1]
	}
	return toEnvelopes(result, reversed), updatedPagingInfo(req.PagingInfo, newCursorEvent), nil
}

func (s *mapTopicStore) Get(cid mh.Multihash) (*Event, error) {
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

// updates paging info with a cursor for given event (or nil)
func updatedPagingInfo(pi *messagev1.PagingInfo, cursorEvent *Event) *messagev1.PagingInfo {
	var cursor *messagev1.Cursor
	if cursorEvent != nil {
		cursor = &messagev1.Cursor{
			Cursor: &messagev1.Cursor_Index{
				Index: &messagev1.IndexCursor{
					SenderTimeNs: cursorEvent.TimestampNs,
					Digest:       cursorEvent.cid,
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

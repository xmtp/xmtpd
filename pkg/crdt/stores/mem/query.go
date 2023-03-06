package memstore

import (
	"bytes"
	"errors"
	"sort"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

var (
	ErrCursorNotFound = errors.New("cursor not found")
)

func (s *MemoryStore) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	s.RLock()
	defer s.RUnlock()

	var start int
	if req.StartTimeNs > 0 {
		start, _ = sort.Find(len(s.eventsByTime), func(i int) int {
			return int(req.StartTimeNs - s.eventsByTime[i].TimestampNs)
		})
	}

	if start == len(s.eventsByTime) {
		// everything is earlier than StartTimeNs
		return &messagev1.QueryResponse{}, nil
	}

	end := len(s.eventsByTime)
	if req.EndTimeNs > 0 {
		upTo := req.EndTimeNs + 1
		end, _ = sort.Find(len(s.eventsByTime), func(i int) int {
			return int(upTo - s.eventsByTime[i].TimestampNs)
		})
	}

	result := s.eventsByTime[start:end]
	if req.PagingInfo == nil {
		return &messagev1.QueryResponse{
			Envelopes: toEnvelopes(result),
		}, nil
	}

	reversed := req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING
	cursor := req.PagingInfo.Cursor.GetIndex()
	if cursor != nil {
		// find the cursor event in the result
		cIdx, found := sort.Find(len(result), func(i int) int {
			res := cursor.SenderTimeNs - result[i].TimestampNs
			if res != 0 {
				return int(res)
			}
			return bytes.Compare(cursor.Digest, result[i].Cid)
		})
		if !found {
			return nil, ErrCursorNotFound
		}
		// reslice the result from the cursor event to the end
		if reversed {
			result = result[:cIdx]
		} else {
			result = result[cIdx+1:]
		}
	}

	limit := int(req.PagingInfo.Limit)
	if reversed {
		if limit != 0 && limit < len(result) {
			result = result[len(result)-int(limit):]
		}
		result = reverseEvents(result)
	} else {
		if limit != 0 && limit < len(result) {
			result = result[:limit]
		}
	}
	resp := &messagev1.QueryResponse{
		Envelopes: toEnvelopes(result),
	}
	if limit > 0 && len(result) == limit {
		lastEvent := result[len(result)-1]
		resp.PagingInfo = &messagev1.PagingInfo{
			Limit:     req.PagingInfo.Limit,
			Direction: req.PagingInfo.Direction,
			Cursor: &messagev1.Cursor{
				Cursor: &messagev1.Cursor_Index{
					Index: &messagev1.IndexCursor{
						SenderTimeNs: lastEvent.TimestampNs,
						Digest:       lastEvent.Cid,
					},
				},
			},
		}
	}
	return resp, nil
}

// shift events from index i to the right
// to create room at the index.
func makeRoomAt(events []*types.Event, i int) []*types.Event {
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
	newEvents := make([]*types.Event, len(events)+1, newCap)
	copy(newEvents, events[:i])
	copy(newEvents[i+1:], events[i:])
	return newEvents
}

func toEnvelopes(events []*types.Event) []*messagev1.Envelope {
	envs := make([]*messagev1.Envelope, len(events))
	for i, ev := range events {
		envs[i] = ev.Envelope
	}
	return envs
}

package memstore

import (
	"bytes"
	"context"
	"errors"
	"sort"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	crdtmemstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	ErrCursorNotFound = errors.New("cursor not found")
)

type MemoryStore struct {
	crdtmemstore.MemoryStore

	log *zap.Logger

	envsByTime []*messagev1.Envelope
}

func New(log *zap.Logger) *MemoryStore {
	return &MemoryStore{
		log: log,
	}
}

func (s *MemoryStore) Close() error {
	return nil
}

func (s *MemoryStore) InsertEnvelope(ctx context.Context, env *messagev1.Envelope) error {
	i, _ := sort.Find(len(s.envsByTime), func(i int) int {
		res := env.TimestampNs - s.envsByTime[i].TimestampNs
		if res != 0 {
			return int(res)
		}
		envCid, _ := utils.BuildEnvelopeCid(env)
		otherEnvCid, _ := utils.BuildEnvelopeCid(s.envsByTime[i])
		return bytes.Compare(envCid, otherEnvCid)
	})
	if i == len(s.envsByTime) {
		s.envsByTime = append(s.envsByTime, env)
	} else {
		s.envsByTime = makeRoomAt(s.envsByTime, i)
	}
	s.envsByTime[i] = env
	return nil
}

func (s *MemoryStore) QueryEnvelopes(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	s.RLock()
	defer s.RUnlock()

	var start int
	if req.StartTimeNs > 0 {
		start, _ = sort.Find(len(s.envsByTime), func(i int) int {
			return int(req.StartTimeNs - s.envsByTime[i].TimestampNs)
		})
	}

	if start == len(s.envsByTime) {
		// everything is earlier than StartTimeNs
		return &messagev1.QueryResponse{}, nil
	}

	end := len(s.envsByTime)
	if req.EndTimeNs > 0 {
		upTo := req.EndTimeNs + 1
		end, _ = sort.Find(len(s.envsByTime), func(i int) int {
			return int(upTo - s.envsByTime[i].TimestampNs)
		})
	}

	result := s.envsByTime[start:end]
	if req.PagingInfo == nil {
		return &messagev1.QueryResponse{
			Envelopes: result,
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
			cid, _ := utils.BuildEnvelopeCid(result[i])
			return bytes.Compare(cursor.Digest, cid)
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

	if reversed {
		if limit := req.PagingInfo.Limit; limit != 0 && int(limit) < len(result) {
			result = result[len(result)-int(limit):]
		}
		var newCursorEnv *messagev1.Envelope
		if len(result) > 0 {
			newCursorEnv = result[0]
		}
		utils.Reverse(result)
		pi, err := updatedPagingInfo(req.PagingInfo, newCursorEnv)
		if err != nil {
			return nil, err
		}
		return &messagev1.QueryResponse{
			Envelopes:  result,
			PagingInfo: pi,
		}, nil
	}

	if limit := req.PagingInfo.Limit; limit != 0 && int(limit) < len(result) {
		result = result[:limit]
	}

	var newCursorEnv *messagev1.Envelope
	if len(result) > 0 {
		newCursorEnv = result[len(result)-1]
	}

	pi, err := updatedPagingInfo(req.PagingInfo, newCursorEnv)
	if err != nil {
		return nil, err
	}

	return &messagev1.QueryResponse{
		Envelopes:  result,
		PagingInfo: pi,
	}, nil
}

// shift events from index i to the right
// to create room at the index.
func makeRoomAt(envs []*messagev1.Envelope, i int) []*messagev1.Envelope {
	// if there's enough capacity in the slice, just shift the tail
	if len(envs) < cap(envs) {
		envs = envs[:len(envs)+1]
		copy(envs[i+1:], envs[i:])
		return envs
	}
	// figure out desired capacity of a new slice
	var newCap int
	// don't need to worry about len(events) == 0
	// because of the !found append in addEvent
	if len(envs) < 1024 {
		newCap = 2 * len(envs)
	} else {
		newCap = len(envs) + 1024
	}
	// copy events into a new slice, leaving a gap at index i
	newEnvs := make([]*messagev1.Envelope, len(envs)+1, newCap)
	copy(newEnvs, envs[:i])
	copy(newEnvs[i+1:], envs[i:])
	return newEnvs
}

// updates paging info with a cursor for given event (or nil)
func updatedPagingInfo(pi *messagev1.PagingInfo, cursorEnv *messagev1.Envelope) (*messagev1.PagingInfo, error) {
	var cursor *messagev1.Cursor
	if cursorEnv != nil {
		cid, err := utils.BuildEnvelopeCid(cursorEnv)
		if err != nil {
			return nil, err
		}
		cursor = &messagev1.Cursor{
			Cursor: &messagev1.Cursor_Index{
				Index: &messagev1.IndexCursor{
					SenderTimeNs: cursorEnv.TimestampNs,
					Digest:       cid,
				},
			},
		}
	}
	// Note that we're modifying the original query's paging info here.
	pi.Cursor = cursor
	return pi, nil
}

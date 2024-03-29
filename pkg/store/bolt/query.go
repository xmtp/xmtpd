package bolt

import (
	"bytes"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	bolt "go.etcd.io/bbolt"
)

// All CIDs should be less or equal than maxCID in terms of bytes.Compare().
// Given that all Multihash values start with the code and size prefix 0xFFFF should do it.
var maxCID = []byte{0xFF, 0xFF}

func (s *Store) Query(ctx context.Context, req *messagev1.QueryRequest) (resp *messagev1.QueryResponse, err error) {
	var lastKey []byte // key of the last envelope
	var envs []*messagev1.Envelope
	// Figure out the iteration parameters:
	// reversed - whether we iterate in reverse or not
	// start, stop - the key to start from (nil => starting from the beginning or end if reversed)
	// stop - the key to stop at (nil => run to the end or beginning if reversed)
	// count - max number of iteration steps (0 => no limit)
	start, stop, reversed, count := computeCursorLoopParameters(req)
	limit := int(count)
	withinLimits := buildCursorLoopCondition(&count, stop, reversed)
	hadCursor := req.PagingInfo != nil && req.PagingInfo.Cursor.GetIndex() != nil
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		c := events.Cursor()
		// advanceFn captures the iteration direction
		advanceFn := c.Next
		if reversed {
			advanceFn = c.Prev
		}
		// Cursor loop
		for k, v := positionCursor(c, start, reversed, hadCursor); k != nil && withinLimits(k); k, v = advanceFn() {
			env, err := types.EnvelopeFromBytes(v)
			if err != nil {
				return err
			}
			envs = append(envs, env)
			lastKey = k
			count--
		}
		// Have to make a copy of the last key before leaving the transaction scope.
		if lastKey != nil {
			lkc := make([]byte, len(lastKey))
			copy(lkc, lastKey)
			lastKey = lkc
		}
		return nil
	})
	resp = &messagev1.QueryResponse{
		Envelopes: envs,
	}
	// Set PagingInfo to pass the cursor back if there could be more pages to follow.
	if limit > 0 && len(envs) == limit {
		timestampNs, cid := types.FromByTimeKey(lastKey)
		resp.PagingInfo = &messagev1.PagingInfo{
			Limit:     req.PagingInfo.Limit,
			Direction: req.PagingInfo.Direction,
			Cursor: &messagev1.Cursor{
				Cursor: &messagev1.Cursor_Index{
					Index: &messagev1.IndexCursor{
						SenderTimeNs: timestampNs,
						Digest:       cid,
					},
				},
			},
		}
	}
	return resp, err
}

// loopCondition says whether k is still in the range of the iteration parameters.
type loopCondition func(k []byte) bool

// return the appropriate loop condition function based on the provided loop parameters.
func buildCursorLoopCondition(limit *uint32, stop []byte, reversed bool) loopCondition {
	var stopCondition loopCondition
	if stop != nil {
		// pick the right stopCondition based on whether we iterate in reverse or not
		if reversed {
			stopCondition = func(k []byte) bool { return bytes.Compare(k, stop) >= 0 }
		} else {
			stopCondition = func(k []byte) bool { return bytes.Compare(k, stop) <= 0 }
		}
	}
	if *limit > 0 {
		if stopCondition == nil {
			// limit only
			return func(k []byte) bool { return *limit > 0 }
		} else {
			// combined stop and limit
			return func(k []byte) bool { return *limit > 0 && stopCondition(k) }
		}
	}
	if stopCondition != nil {
		// stop condition only
		return stopCondition
	}
	// no conditions
	return func(k []byte) bool { return true }
}

// extract the iteration parameters from the query.
func computeCursorLoopParameters(req *messagev1.QueryRequest) (start, stop []byte, reversed bool, limit uint32) {
	if req.PagingInfo == nil {
		start = toKey(req.StartTimeNs, nil)
		stop = toKey(req.EndTimeNs, maxCID)
		return start, stop, reversed, limit
	}
	reversed = req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING
	limit = req.PagingInfo.Limit
	cursor := req.PagingInfo.Cursor.GetIndex()
	if cursor == nil {
		if reversed {
			start = toKey(req.EndTimeNs, maxCID)
		} else {
			start = toKey(req.StartTimeNs, nil)
		}
	} else {

		start = toKey(cursor.SenderTimeNs, cursor.Digest)
	}
	if reversed {
		stop = toKey(req.StartTimeNs, nil)
	} else {
		stop = toKey(req.EndTimeNs, maxCID)
	}
	return start, stop, reversed, limit
}

func toKey(timestamp uint64, cid []byte) []byte {
	if timestamp == 0 {
		return nil
	}
	return types.ToByTimeKey(timestamp, cid)
}

// Position the cursor for iteration start based on the iteration parameters.
// Return the key/value at that position.
func positionCursor(c *bolt.Cursor, start []byte, reversed, hadCursor bool) (k, v []byte) {
	if start == nil {
		// if no start we start at either end.
		if reversed {
			k, v = c.Last()
		} else {
			k, v = c.First()
		}
	} else {
		// otherwise seek to the start key
		// but if k != start and we are reversed, we need to back up 1 step
		// because Seek() will put the cursor on the NEXT higher key.
		k, v = c.Seek(start)
		if reversed && bytes.Compare(start, k) < 0 {
			k, v = c.Prev()
		}
	}
	if hadCursor {
		// if start comes from a cursor from the query we need to advance one step
		// because the value at the cursor was already sent with the previous page.
		if reversed {
			k, v = c.Prev()
		} else {
			k, v = c.Next()
		}
	}
	return k, v
}

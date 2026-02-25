package db

import (
	"database/sql"
	"math"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type (
	VectorClock = map[uint32]uint64
	Topic       = []byte
)

const GatewayEnvelopeBandWidth int64 = 1_000_000

func NullInt32(v int32) sql.NullInt32 {
	return sql.NullInt32{Int32: v, Valid: true}
}

func NullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

func SetVectorClockByTopics(
	q *queries.SelectGatewayEnvelopesByTopicsParams,
	vc VectorClock,
) *queries.SelectGatewayEnvelopesByTopicsParams {
	q.CursorNodeIds = make([]int32, 0, len(vc))
	q.CursorSequenceIds = make([]int64, 0, len(vc))
	for nodeID, sequenceID := range vc {
		if nodeID > math.MaxInt32 {
			continue
		}
		q.CursorNodeIds = append(q.CursorNodeIds, int32(nodeID))
		q.CursorSequenceIds = append(q.CursorSequenceIds, int64(sequenceID))
	}
	return q
}

func SetVectorClockByOriginators(
	q *queries.SelectGatewayEnvelopesByOriginatorsParams,
	vc VectorClock,
) *queries.SelectGatewayEnvelopesByOriginatorsParams {
	q.CursorNodeIds = make([]int32, 0, len(vc))
	q.CursorSequenceIds = make([]int64, 0, len(vc))
	for nodeID, sequenceID := range vc {
		if nodeID > math.MaxInt32 {
			continue
		}
		q.CursorNodeIds = append(q.CursorNodeIds, int32(nodeID))
		q.CursorSequenceIds = append(q.CursorSequenceIds, int64(sequenceID))
	}
	return q
}

func SetVectorClockUnfiltered(
	q *queries.SelectGatewayEnvelopesUnfilteredParams,
	vc VectorClock,
) *queries.SelectGatewayEnvelopesUnfilteredParams {
	q.CursorNodeIds = make([]int32, 0, len(vc))
	q.CursorSequenceIds = make([]int64, 0, len(vc))
	for nodeID, sequenceID := range vc {
		if nodeID > math.MaxInt32 {
			continue
		}
		q.CursorNodeIds = append(q.CursorNodeIds, int32(nodeID))
		q.CursorSequenceIds = append(q.CursorSequenceIds, int64(sequenceID))
	}
	return q
}

func ToVectorClock(rows []queries.GatewayEnvelopesLatest) VectorClock {
	vc := make(VectorClock)
	for _, row := range rows {
		vc[uint32(row.OriginatorNodeID)] = uint64(row.OriginatorSequenceID)
	}
	return vc
}

func TransformRowsByTopic(
	rows []queries.SelectGatewayEnvelopesByTopicsRow,
) []queries.GatewayEnvelopesView {
	result := make([]queries.GatewayEnvelopesView, len(rows))
	for i, row := range rows {
		result[i] = queries.GatewayEnvelopesView(row)
	}
	return result
}

func TransformRowsByOriginator(
	rows []queries.SelectGatewayEnvelopesByOriginatorsRow,
) []queries.GatewayEnvelopesView {
	result := make([]queries.GatewayEnvelopesView, len(rows))
	for i, row := range rows {
		result[i] = queries.GatewayEnvelopesView(row)
	}
	return result
}

// FillMissingOriginators ensures that every originator from allOriginators
// is present in the vector clock. Missing originators are added with
// a sequence ID of 0, meaning "start from the beginning".
func FillMissingOriginators(vc VectorClock, allOriginators []uint32) {
	for _, id := range allOriginators {
		if _, ok := vc[id]; !ok {
			vc[id] = 0
		}
	}
}

// TopicCursors maps raw topic bytes (as string key) to a per-topic VectorClock.
type TopicCursors map[string]VectorClock

// SetPerTopicCursors flattens TopicCursors into the parallel arrays required by
// SelectGatewayEnvelopesByPerTopicCursors. Each (topic, nodeID, seqID) triple
// produces one entry in the three arrays.
func SetPerTopicCursors(
	q *queries.SelectGatewayEnvelopesByPerTopicCursorsParams,
	tc TopicCursors,
) {
	// Count total entries for pre-allocation.
	total := 0
	for _, vc := range tc {
		total += len(vc)
	}

	q.CursorTopics = make([][]byte, 0, total)
	q.CursorNodeIds = make([]int32, 0, total)
	q.CursorSequenceIds = make([]int64, 0, total)

	for topicKey, vc := range tc {
		topicBytes := []byte(topicKey)
		for nodeID, seqID := range vc {
			if nodeID > math.MaxInt32 || seqID > uint64(math.MaxInt64) {
				continue
			}
			q.CursorTopics = append(q.CursorTopics, topicBytes)
			q.CursorNodeIds = append(q.CursorNodeIds, int32(nodeID))
			q.CursorSequenceIds = append(q.CursorSequenceIds, int64(seqID))
		}
	}
}

// TransformRowsByPerTopicCursors converts per-topic cursor rows to the common
// GatewayEnvelopesView type.
func TransformRowsByPerTopicCursors(
	rows []queries.SelectGatewayEnvelopesByPerTopicCursorsRow,
) []queries.GatewayEnvelopesView {
	result := make([]queries.GatewayEnvelopesView, len(rows))
	for i, row := range rows {
		result[i] = queries.GatewayEnvelopesView(row)
	}
	return result
}

// CalculateRowsPerEntry computes the per-(topic, originator) sub-limit
// for the per-topic cursor query. Returns at least 10 to avoid starving
// low-volume originators.
func CalculateRowsPerEntry(numEntries int, rowLimit int32) int32 {
	if numEntries == 0 {
		return rowLimit
	}
	if numEntries > math.MaxInt32 {
		numEntries = math.MaxInt32
	}
	rpe := max(rowLimit/int32(numEntries), 10)
	return rpe
}

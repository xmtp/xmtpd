package db

import (
	"database/sql"

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

package db

import (
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type VectorClock = map[uint32]uint64

func NullInt32(v int32) sql.NullInt32 {
	return sql.NullInt32{Int32: v, Valid: true}
}

func NullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

func SetVectorClock(
	q *queries.SelectGatewayEnvelopesParams,
	vc VectorClock,
) *queries.SelectGatewayEnvelopesParams {
	q.CursorNodeIds = make([]int32, 0, len(vc))
	q.CursorSequenceIds = make([]int64, 0, len(vc))
	for nodeID, sequenceID := range vc {
		q.CursorNodeIds = append(q.CursorNodeIds, int32(nodeID))
		q.CursorSequenceIds = append(q.CursorSequenceIds, int64(sequenceID))
	}
	return q
}

func ToVectorClock(rows []queries.SelectVectorClockRow) VectorClock {
	vc := make(VectorClock)
	for _, row := range rows {
		vc[uint32(row.OriginatorNodeID)] = uint64(row.OriginatorSequenceID)
	}
	return vc
}

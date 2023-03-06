package postgresstore

import (
	"fmt"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func (s *Store) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	if len(req.ContentTopics) == 0 {
		req.ContentTopics = []string{s.topic}
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}
	topic := req.ContentTopics[0]

	if topic != s.topic {
		return nil, ErrTopicMismatch
	}

	baseSQL := "SELECT cid, topic, timestamp_ns, message FROM events WHERE topic = $1"
	args := []any{topic}

	timeFilterSQL := ""
	if req.StartTimeNs > 0 {
		timeFilterSQL = fmt.Sprintf(" AND timestamp_ns >= $%d", len(args)+1)
		args = append(args, req.StartTimeNs)
	}
	if req.EndTimeNs > 0 {
		timeFilterSQL += fmt.Sprintf(" AND timestamp_ns <= $%d", len(args)+1)
		args = append(args, req.EndTimeNs)
	}

	sortSQL := " ORDER BY timestamp_ns, cid"
	limitSQL := ""
	cursorFilterSQL := ""
	var limit int
	if req.PagingInfo != nil {
		if req.PagingInfo.Cursor != nil && req.PagingInfo.Cursor.GetIndex() != nil {
			cursor := req.PagingInfo.Cursor.GetIndex()
			if req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING {
				cursorFilterSQL = fmt.Sprintf(" AND timestamp_ns < $%d OR (timestamp_ns = $%d AND cid < $%d)", len(args)+1, len(args)+2, len(args)+3)
				args = append(args, cursor.SenderTimeNs, cursor.SenderTimeNs, string(cursor.Digest))
			} else {
				cursorFilterSQL = fmt.Sprintf(" AND timestamp_ns > $%d OR (timestamp_ns = $%d AND cid > $%d)", len(args)+1, len(args)+2, len(args)+3)
				args = append(args, cursor.SenderTimeNs, cursor.SenderTimeNs, string(cursor.Digest))
			}
		}

		if req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING {
			sortSQL = " ORDER BY timestamp_ns DESC, cid"
		}

		if limit = int(req.PagingInfo.Limit); limit > 0 {
			limitSQL = fmt.Sprintf(" LIMIT $%d", len(args)+1)
			args = append(args, req.PagingInfo.Limit)
		}
	}

	sql := baseSQL + timeFilterSQL + cursorFilterSQL + sortSQL + limitSQL

	s.log.Debug("querying", zap.String("sql", sql))
	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		cidHex      string
		evTopic     string
		timestampNS uint64
		message     []byte
		envs        []*messagev1.Envelope
	)
	for rows.Next() {
		err := rows.Scan(&cidHex, &evTopic, &timestampNS, &message)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &messagev1.Envelope{
			ContentTopic: evTopic,
			TimestampNs:  timestampNS,
			Message:      message,
		})
	}

	resp := &messagev1.QueryResponse{
		Envelopes: envs,
	}
	if limit > 0 && len(envs) == limit {
		resp.PagingInfo = &messagev1.PagingInfo{
			Limit:     req.PagingInfo.Limit,
			Direction: req.PagingInfo.Direction,
			Cursor: &messagev1.Cursor{
				Cursor: &messagev1.Cursor_Index{
					Index: &messagev1.IndexCursor{
						SenderTimeNs: timestampNS,
						Digest:       []byte(cidHex),
					},
				},
			},
		}
	}

	return resp, nil
}

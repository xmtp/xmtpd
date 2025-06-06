// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: envelopes.sql

package queries

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const deleteStagedOriginatorEnvelope = `-- name: DeleteStagedOriginatorEnvelope :execrows
DELETE FROM staged_originator_envelopes
WHERE id = $1
`

func (q *Queries) DeleteStagedOriginatorEnvelope(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteStagedOriginatorEnvelope, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getLatestCursor = `-- name: GetLatestCursor :many
SELECT originator_node_id,
	MAX(originator_sequence_id)::BIGINT AS max_sequence_id
FROM gateway_envelopes
GROUP BY originator_node_id
`

type GetLatestCursorRow struct {
	OriginatorNodeID int32
	MaxSequenceID    int64
}

func (q *Queries) GetLatestCursor(ctx context.Context) ([]GetLatestCursorRow, error) {
	rows, err := q.db.QueryContext(ctx, getLatestCursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLatestCursorRow
	for rows.Next() {
		var i GetLatestCursorRow
		if err := rows.Scan(&i.OriginatorNodeID, &i.MaxSequenceID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLatestSequenceId = `-- name: GetLatestSequenceId :one
SELECT COALESCE(max(originator_sequence_id), 0)::BIGINT AS originator_sequence_id
FROM gateway_envelopes
WHERE originator_node_id = $1
`

func (q *Queries) GetLatestSequenceId(ctx context.Context, originatorNodeID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLatestSequenceId, originatorNodeID)
	var originator_sequence_id int64
	err := row.Scan(&originator_sequence_id)
	return originator_sequence_id, err
}

const insertGatewayEnvelope = `-- name: InsertGatewayEnvelope :execrows
INSERT INTO gateway_envelopes(
		originator_node_id,
		originator_sequence_id,
		topic,
		originator_envelope,
		payer_id,
		gateway_time,
        expiry
	)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		COALESCE($6, NOW()),
        $7
	) ON CONFLICT DO NOTHING
`

type InsertGatewayEnvelopeParams struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	OriginatorEnvelope   []byte
	PayerID              sql.NullInt32
	GatewayTime          interface{}
	Expiry               sql.NullInt64
}

func (q *Queries) InsertGatewayEnvelope(ctx context.Context, arg InsertGatewayEnvelopeParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertGatewayEnvelope,
		arg.OriginatorNodeID,
		arg.OriginatorSequenceID,
		arg.Topic,
		arg.OriginatorEnvelope,
		arg.PayerID,
		arg.GatewayTime,
		arg.Expiry,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const insertStagedOriginatorEnvelope = `-- name: InsertStagedOriginatorEnvelope :one
SELECT id, originator_time, topic, payer_envelope
FROM insert_staged_originator_envelope($1, $2)
`

type InsertStagedOriginatorEnvelopeParams struct {
	Topic         []byte
	PayerEnvelope []byte
}

func (q *Queries) InsertStagedOriginatorEnvelope(ctx context.Context, arg InsertStagedOriginatorEnvelopeParams) (StagedOriginatorEnvelope, error) {
	row := q.db.QueryRowContext(ctx, insertStagedOriginatorEnvelope, arg.Topic, arg.PayerEnvelope)
	var i StagedOriginatorEnvelope
	err := row.Scan(
		&i.ID,
		&i.OriginatorTime,
		&i.Topic,
		&i.PayerEnvelope,
	)
	return i, err
}

const selectGatewayEnvelopes = `-- name: SelectGatewayEnvelopes :many
SELECT gateway_time, originator_node_id, originator_sequence_id, topic, originator_envelope, payer_id, expiry
FROM select_gateway_envelopes(
		$1::INT [],
		$2::BIGINT [],
		$3::BYTEA [],
		$4::INT [],
		$5::INT
	)
`

type SelectGatewayEnvelopesParams struct {
	CursorNodeIds     []int32
	CursorSequenceIds []int64
	Topics            [][]byte
	OriginatorNodeIds []int32
	RowLimit          int32
}

func (q *Queries) SelectGatewayEnvelopes(ctx context.Context, arg SelectGatewayEnvelopesParams) ([]GatewayEnvelope, error) {
	rows, err := q.db.QueryContext(ctx, selectGatewayEnvelopes,
		pq.Array(arg.CursorNodeIds),
		pq.Array(arg.CursorSequenceIds),
		pq.Array(arg.Topics),
		pq.Array(arg.OriginatorNodeIds),
		arg.RowLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GatewayEnvelope
	for rows.Next() {
		var i GatewayEnvelope
		if err := rows.Scan(
			&i.GatewayTime,
			&i.OriginatorNodeID,
			&i.OriginatorSequenceID,
			&i.Topic,
			&i.OriginatorEnvelope,
			&i.PayerID,
			&i.Expiry,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectNewestFromTopics = `-- name: SelectNewestFromTopics :many
SELECT e.gateway_time, e.originator_node_id, e.originator_sequence_id, e.topic, e.originator_envelope, e.payer_id, e.expiry
FROM gateway_envelopes e
	INNER JOIN (
		SELECT topic,
			MAX(gateway_time) AS max_time
		FROM gateway_envelopes
		WHERE topic = ANY($1::BYTEA [])
		GROUP BY topic
	) t ON e.topic = t.topic
	AND e.gateway_time = t.max_time
ORDER BY e.topic
`

func (q *Queries) SelectNewestFromTopics(ctx context.Context, topics [][]byte) ([]GatewayEnvelope, error) {
	rows, err := q.db.QueryContext(ctx, selectNewestFromTopics, pq.Array(topics))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GatewayEnvelope
	for rows.Next() {
		var i GatewayEnvelope
		if err := rows.Scan(
			&i.GatewayTime,
			&i.OriginatorNodeID,
			&i.OriginatorSequenceID,
			&i.Topic,
			&i.OriginatorEnvelope,
			&i.PayerID,
			&i.Expiry,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectStagedOriginatorEnvelopes = `-- name: SelectStagedOriginatorEnvelopes :many
SELECT id, originator_time, topic, payer_envelope
FROM staged_originator_envelopes
WHERE id > $1
ORDER BY id ASC
LIMIT $2
`

type SelectStagedOriginatorEnvelopesParams struct {
	LastSeenID int64
	NumRows    int32
}

func (q *Queries) SelectStagedOriginatorEnvelopes(ctx context.Context, arg SelectStagedOriginatorEnvelopesParams) ([]StagedOriginatorEnvelope, error) {
	rows, err := q.db.QueryContext(ctx, selectStagedOriginatorEnvelopes, arg.LastSeenID, arg.NumRows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StagedOriginatorEnvelope
	for rows.Next() {
		var i StagedOriginatorEnvelope
		if err := rows.Scan(
			&i.ID,
			&i.OriginatorTime,
			&i.Topic,
			&i.PayerEnvelope,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectVectorClock = `-- name: SelectVectorClock :many
SELECT DISTINCT ON (originator_node_id) originator_node_id,
	originator_sequence_id,
	originator_envelope
FROM gateway_envelopes
ORDER BY originator_node_id,
	originator_sequence_id DESC
`

type SelectVectorClockRow struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	OriginatorEnvelope   []byte
}

func (q *Queries) SelectVectorClock(ctx context.Context) ([]SelectVectorClockRow, error) {
	rows, err := q.db.QueryContext(ctx, selectVectorClock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectVectorClockRow
	for rows.Next() {
		var i SelectVectorClockRow
		if err := rows.Scan(&i.OriginatorNodeID, &i.OriginatorSequenceID, &i.OriginatorEnvelope); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

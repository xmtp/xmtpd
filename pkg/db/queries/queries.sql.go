// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package queries

import (
	"context"
	"database/sql"
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

const insertGatewayEnvelope = `-- name: InsertGatewayEnvelope :execrows
SELECT
	insert_gateway_envelope($1, $2, $3, $4)
`

type InsertGatewayEnvelopeParams struct {
	OriginatorID         int32
	OriginatorSequenceID int64
	Topic                []byte
	OriginatorEnvelope   []byte
}

func (q *Queries) InsertGatewayEnvelope(ctx context.Context, arg InsertGatewayEnvelopeParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertGatewayEnvelope,
		arg.OriginatorID,
		arg.OriginatorSequenceID,
		arg.Topic,
		arg.OriginatorEnvelope,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const insertNodeInfo = `-- name: InsertNodeInfo :execrows
INSERT INTO node_info(node_id, public_key)
	VALUES ($1, $2)
ON CONFLICT
	DO NOTHING
`

type InsertNodeInfoParams struct {
	NodeID    int32
	PublicKey []byte
}

func (q *Queries) InsertNodeInfo(ctx context.Context, arg InsertNodeInfoParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertNodeInfo, arg.NodeID, arg.PublicKey)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const insertStagedOriginatorEnvelope = `-- name: InsertStagedOriginatorEnvelope :one
SELECT
	id, originator_time, topic, payer_envelope
FROM
	insert_staged_originator_envelope($1, $2)
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
SELECT
	id, originator_node_id, originator_sequence_id, topic, originator_envelope
FROM
	gateway_envelopes
WHERE ($1::BYTEA IS NULL
	OR topic = $1)
AND ($2::INT IS NULL
	OR originator_node_id = $2)
AND ($3::BIGINT IS NULL
	OR originator_sequence_id > $3)
AND ($4::BIGINT IS NULL
	OR id > $4)
LIMIT $5::INT
`

type SelectGatewayEnvelopesParams struct {
	Topic                []byte
	OriginatorNodeID     sql.NullInt32
	OriginatorSequenceID sql.NullInt64
	GatewaySequenceID    sql.NullInt64
	RowLimit             sql.NullInt32
}

func (q *Queries) SelectGatewayEnvelopes(ctx context.Context, arg SelectGatewayEnvelopesParams) ([]GatewayEnvelope, error) {
	rows, err := q.db.QueryContext(ctx, selectGatewayEnvelopes,
		arg.Topic,
		arg.OriginatorNodeID,
		arg.OriginatorSequenceID,
		arg.GatewaySequenceID,
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
			&i.ID,
			&i.OriginatorNodeID,
			&i.OriginatorSequenceID,
			&i.Topic,
			&i.OriginatorEnvelope,
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

const selectNodeInfo = `-- name: SelectNodeInfo :one
SELECT
	node_id, public_key, singleton_id
FROM
	node_info
WHERE
	singleton_id = 1
`

func (q *Queries) SelectNodeInfo(ctx context.Context) (NodeInfo, error) {
	row := q.db.QueryRowContext(ctx, selectNodeInfo)
	var i NodeInfo
	err := row.Scan(&i.NodeID, &i.PublicKey, &i.SingletonID)
	return i, err
}

const selectStagedOriginatorEnvelopes = `-- name: SelectStagedOriginatorEnvelopes :many
SELECT
	id, originator_time, topic, payer_envelope
FROM
	staged_originator_envelopes
WHERE
	id > $1
ORDER BY
	id ASC
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

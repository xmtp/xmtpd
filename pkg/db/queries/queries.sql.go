// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

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

const findOrCreatePayer = `-- name: FindOrCreatePayer :one
INSERT INTO payers(address)
	VALUES ($1)
ON CONFLICT (address)
	DO UPDATE SET
		address = $1
	RETURNING
		id
`

func (q *Queries) FindOrCreatePayer(ctx context.Context, address string) (int32, error) {
	row := q.db.QueryRowContext(ctx, findOrCreatePayer, address)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getAddressLogs = `-- name: GetAddressLogs :many
SELECT
	a.address,
	encode(a.inbox_id, 'hex') AS inbox_id,
	a.association_sequence_id
FROM
	address_log a
	INNER JOIN (
		SELECT
			address,
			MAX(association_sequence_id) AS max_association_sequence_id
		FROM
			address_log
		WHERE
			address = ANY ($1::TEXT[])
			AND revocation_sequence_id IS NULL
		GROUP BY
			address) b ON a.address = b.address
	AND a.association_sequence_id = b.max_association_sequence_id
`

type GetAddressLogsRow struct {
	Address               string
	InboxID               string
	AssociationSequenceID sql.NullInt64
}

func (q *Queries) GetAddressLogs(ctx context.Context, addresses []string) ([]GetAddressLogsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAddressLogs, pq.Array(addresses))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAddressLogsRow
	for rows.Next() {
		var i GetAddressLogsRow
		if err := rows.Scan(&i.Address, &i.InboxID, &i.AssociationSequenceID); err != nil {
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

const getBlocksInRange = `-- name: GetBlocksInRange :many
SELECT DISTINCT ON (block_number)
	block_number,
	block_hash
FROM
	blockchain_messages
WHERE
	block_number BETWEEN $1 AND $2
	AND block_hash IS NOT NULL
	AND is_canonical = TRUE
ORDER BY
	block_number ASC,
	block_hash
`

type GetBlocksInRangeParams struct {
	StartBlock uint64
	EndBlock   uint64
}

type GetBlocksInRangeRow struct {
	BlockNumber uint64
	BlockHash   []byte
}

// Returns blocks in ascending order (oldest to newest)
// StartBlock should be the lower bound (older block)
// EndBlock should be the upper bound (newer block)
// Example: GetBlocksInRange(1000, 2000), returns 1000, 1001, 1002, ..., 2000
func (q *Queries) GetBlocksInRange(ctx context.Context, arg GetBlocksInRangeParams) ([]GetBlocksInRangeRow, error) {
	rows, err := q.db.QueryContext(ctx, getBlocksInRange, arg.StartBlock, arg.EndBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBlocksInRangeRow
	for rows.Next() {
		var i GetBlocksInRangeRow
		if err := rows.Scan(&i.BlockNumber, &i.BlockHash); err != nil {
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

const getLatestBlock = `-- name: GetLatestBlock :one
SELECT
	block_number,
	block_hash
FROM
	latest_block
WHERE
	contract_address = $1
`

type GetLatestBlockRow struct {
	BlockNumber int64
	BlockHash   []byte
}

func (q *Queries) GetLatestBlock(ctx context.Context, contractAddress string) (GetLatestBlockRow, error) {
	row := q.db.QueryRowContext(ctx, getLatestBlock, contractAddress)
	var i GetLatestBlockRow
	err := row.Scan(&i.BlockNumber, &i.BlockHash)
	return i, err
}

const getLatestCursor = `-- name: GetLatestCursor :many
SELECT
	originator_node_id,
	MAX(originator_sequence_id)::BIGINT AS max_sequence_id
FROM
	gateway_envelopes
GROUP BY
	originator_node_id
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
SELECT
	COALESCE(max(originator_sequence_id), 0)::BIGINT AS originator_sequence_id
FROM
	gateway_envelopes
WHERE
	originator_node_id = $1
`

func (q *Queries) GetLatestSequenceId(ctx context.Context, originatorNodeID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLatestSequenceId, originatorNodeID)
	var originator_sequence_id int64
	err := row.Scan(&originator_sequence_id)
	return originator_sequence_id, err
}

const insertAddressLog = `-- name: InsertAddressLog :execrows
INSERT INTO address_log(address, inbox_id, association_sequence_id, revocation_sequence_id)
	VALUES ($1, decode($2, 'hex'), $3, NULL)
ON CONFLICT (address, inbox_id)
	DO UPDATE SET
		revocation_sequence_id = NULL, association_sequence_id = $3
	WHERE (address_log.revocation_sequence_id IS NULL
		OR address_log.revocation_sequence_id < $3)
		AND address_log.association_sequence_id < $3
`

type InsertAddressLogParams struct {
	Address               string
	InboxID               string
	AssociationSequenceID sql.NullInt64
}

func (q *Queries) InsertAddressLog(ctx context.Context, arg InsertAddressLogParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertAddressLog, arg.Address, arg.InboxID, arg.AssociationSequenceID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const insertBlockchainMessage = `-- name: InsertBlockchainMessage :exec
INSERT INTO blockchain_messages(block_number, block_hash, originator_node_id, originator_sequence_id, is_canonical)
	VALUES ($1, $2, $3, $4, $5)
ON CONFLICT
	DO NOTHING
`

type InsertBlockchainMessageParams struct {
	BlockNumber          uint64
	BlockHash            []byte
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	IsCanonical          bool
}

func (q *Queries) InsertBlockchainMessage(ctx context.Context, arg InsertBlockchainMessageParams) error {
	_, err := q.db.ExecContext(ctx, insertBlockchainMessage,
		arg.BlockNumber,
		arg.BlockHash,
		arg.OriginatorNodeID,
		arg.OriginatorSequenceID,
		arg.IsCanonical,
	)
	return err
}

const insertGatewayEnvelope = `-- name: InsertGatewayEnvelope :execrows
INSERT INTO gateway_envelopes(originator_node_id, originator_sequence_id, topic, originator_envelope)
	VALUES ($1, $2, $3, $4)
ON CONFLICT
	DO NOTHING
`

type InsertGatewayEnvelopeParams struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	OriginatorEnvelope   []byte
}

func (q *Queries) InsertGatewayEnvelope(ctx context.Context, arg InsertGatewayEnvelopeParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertGatewayEnvelope,
		arg.OriginatorNodeID,
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

const revokeAddressFromLog = `-- name: RevokeAddressFromLog :execrows
UPDATE
	address_log
SET
	revocation_sequence_id = $1
WHERE
	address = $2
	AND inbox_id = decode($3, 'hex')
`

type RevokeAddressFromLogParams struct {
	RevocationSequenceID sql.NullInt64
	Address              string
	InboxID              string
}

func (q *Queries) RevokeAddressFromLog(ctx context.Context, arg RevokeAddressFromLogParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, revokeAddressFromLog, arg.RevocationSequenceID, arg.Address, arg.InboxID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const selectGatewayEnvelopes = `-- name: SelectGatewayEnvelopes :many
SELECT
	gateway_time, originator_node_id, originator_sequence_id, topic, originator_envelope
FROM
	select_gateway_envelopes($1::INT[], $2::BIGINT[], $3::BYTEA[], $4::INT[], $5::INT)
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

const selectVectorClock = `-- name: SelectVectorClock :many
SELECT DISTINCT ON (originator_node_id)
	originator_node_id,
	originator_sequence_id,
	originator_envelope
FROM
	gateway_envelopes
ORDER BY
	originator_node_id,
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

const setLatestBlock = `-- name: SetLatestBlock :exec
INSERT INTO latest_block(contract_address, block_number, block_hash)
	VALUES ($1, $2, $3)
ON CONFLICT (contract_address)
	DO UPDATE SET
		block_number = $2, block_hash = $3
	WHERE
		$2 > latest_block.block_number
		AND $3 != latest_block.block_hash
`

type SetLatestBlockParams struct {
	ContractAddress string
	BlockNumber     int64
	BlockHash       []byte
}

func (q *Queries) SetLatestBlock(ctx context.Context, arg SetLatestBlockParams) error {
	_, err := q.db.ExecContext(ctx, setLatestBlock, arg.ContractAddress, arg.BlockNumber, arg.BlockHash)
	return err
}

const updateBlocksCanonicalityInRange = `-- name: UpdateBlocksCanonicalityInRange :exec
UPDATE
	blockchain_messages AS bm
SET
	is_canonical = FALSE
FROM (
	SELECT
		block_number
	FROM
		blockchain_messages
	WHERE
		bm.block_number BETWEEN $1 AND $2
	FOR UPDATE) AS locked_rows
WHERE
	bm.block_number = locked_rows.block_number
`

type UpdateBlocksCanonicalityInRangeParams struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

func (q *Queries) UpdateBlocksCanonicalityInRange(ctx context.Context, arg UpdateBlocksCanonicalityInRangeParams) error {
	_, err := q.db.ExecContext(ctx, updateBlocksCanonicalityInRange, arg.StartBlockNumber, arg.EndBlockNumber)
	return err
}

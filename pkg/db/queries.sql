-- name: InsertNodeInfo :execrows
INSERT INTO node_info(node_id, public_key)
	VALUES (@node_id, @public_key)
ON CONFLICT
	DO NOTHING;

-- name: SelectNodeInfo :one
SELECT
	*
FROM
	node_info
WHERE
	singleton_id = 1;

-- name: InsertGatewayEnvelope :execrows
INSERT INTO gateway_envelopes(originator_node_id, originator_sequence_id, topic, originator_envelope)
	VALUES (@originator_node_id, @originator_sequence_id, @topic, @originator_envelope)
ON CONFLICT
	DO NOTHING;

-- name: SelectGatewayEnvelopes :many
SELECT
	*
FROM
	select_gateway_envelopes(@cursor_node_ids::INT[], @cursor_sequence_ids::BIGINT[], @topics::BYTEA[], @originator_node_ids::INT[], @row_limit::INT);

-- name: InsertStagedOriginatorEnvelope :one
SELECT
	*
FROM
	insert_staged_originator_envelope(@topic, @payer_envelope);

-- name: SelectStagedOriginatorEnvelopes :many
SELECT
	*
FROM
	staged_originator_envelopes
WHERE
	id > @last_seen_id
ORDER BY
	id ASC
LIMIT @num_rows;

-- name: DeleteStagedOriginatorEnvelope :execrows
DELETE FROM staged_originator_envelopes
WHERE id = @id;

-- name: SelectVectorClock :many
SELECT DISTINCT ON (originator_node_id)
	originator_node_id,
	originator_sequence_id,
	originator_envelope
FROM
	gateway_envelopes
ORDER BY
	originator_node_id,
	originator_sequence_id DESC;

-- name: GetAddressLogs :many
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
			address = ANY (@addresses::TEXT[])
			AND revocation_sequence_id IS NULL
		GROUP BY
			address) b ON a.address = b.address
	AND a.association_sequence_id = b.max_association_sequence_id;

-- name: InsertAddressLog :execrows
INSERT INTO address_log(address, inbox_id, association_sequence_id, revocation_sequence_id)
	VALUES (@address, decode(@inbox_id, 'hex'), @association_sequence_id, NULL)
ON CONFLICT (address, inbox_id)
	DO UPDATE SET
		revocation_sequence_id = NULL, association_sequence_id = @association_sequence_id
	WHERE (address_log.revocation_sequence_id IS NULL
		OR address_log.revocation_sequence_id < @association_sequence_id)
		AND address_log.association_sequence_id < @association_sequence_id;

-- name: RevokeAddressFromLog :execrows
UPDATE
	address_log
SET
	revocation_sequence_id = @revocation_sequence_id
WHERE
	address = @address
	AND inbox_id = decode(@inbox_id, 'hex');

-- name: GetLatestSequenceId :one
SELECT
	COALESCE(max(originator_sequence_id), 0)::BIGINT AS originator_sequence_id
FROM
	gateway_envelopes
WHERE
	originator_node_id = @originator_node_id;

-- name: SetLatestBlock :exec
INSERT INTO latest_block(contract_address, block_number, block_hash)
	VALUES (@contract_address, @block_number, @block_hash)
ON CONFLICT (contract_address)
	DO UPDATE SET
		block_number = @block_number, block_hash = @block_hash
	WHERE
		@block_number > latest_block.block_number
		AND @block_hash != latest_block.block_hash;

-- name: GetLatestBlock :one
SELECT
	block_number,
	block_hash
FROM
	latest_block
WHERE
	contract_address = @contract_address;

-- name: GetLatestCursor :many
SELECT
    originator_node_id,
    MAX(originator_sequence_id)::BIGINT AS max_sequence_id
FROM
    gateway_envelopes
GROUP BY
    originator_node_id;
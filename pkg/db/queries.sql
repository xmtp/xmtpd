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
SELECT
	insert_gateway_envelope(@originator_id, @originator_sequence_id, @topic, @originator_envelope);

-- name: SelectGatewayEnvelopes :many
SELECT
	*
FROM
	gateway_envelopes
WHERE (sqlc.narg('topic')::BYTEA IS NULL
	OR topic = @topic)
AND (sqlc.narg('originator_node_id')::INT IS NULL
	OR originator_node_id = @originator_node_id)
AND (sqlc.narg('originator_sequence_id')::BIGINT IS NULL
	OR originator_sequence_id > @originator_sequence_id)
AND (sqlc.narg('gateway_sequence_id')::BIGINT IS NULL
	OR id > @gateway_sequence_id)
LIMIT sqlc.narg('row_limit')::INT;

-- name: InsertStagedOriginatorEnvelope :one
SELECT
	*
FROM
	insert_staged_originator_envelope(@payer_envelope);

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


-- name: InsertGatewayEnvelope :execrows
INSERT INTO gateway_envelopes(originator_node_id, originator_sequence_id, topic, originator_envelope, payer_id, gateway_time)
	VALUES (@originator_node_id, @originator_sequence_id, @topic, @originator_envelope, @payer_id, COALESCE(@gateway_time, NOW()))
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

-- name: GetLatestSequenceId :one
SELECT
	COALESCE(max(originator_sequence_id), 0)::BIGINT AS originator_sequence_id
FROM
	gateway_envelopes
WHERE
	originator_node_id = @originator_node_id;

-- name: GetLatestCursor :many
SELECT
	originator_node_id,
	MAX(originator_sequence_id)::BIGINT AS max_sequence_id
FROM
	gateway_envelopes
GROUP BY
	originator_node_id;
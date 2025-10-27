-- name: InsertStagedOriginatorEnvelope :one
SELECT *
FROM insert_staged_originator_envelope(@topic, @payer_envelope);

-- name: SelectStagedOriginatorEnvelopes :many
SELECT *
FROM staged_originator_envelopes
WHERE id > @last_seen_id
ORDER BY id ASC
LIMIT @num_rows;

-- name: DeleteStagedOriginatorEnvelope :execrows
DELETE FROM staged_originator_envelopes
WHERE id = @id;

-- name: SelectVectorClock :many
SELECT
    originator_node_id,
    originator_sequence_id,
    gateway_time
FROM gateway_envelopes_latest
ORDER BY originator_node_id;

-- name: GetLatestSequenceId :one
SELECT COALESCE((
    SELECT originator_sequence_id
    FROM gateway_envelopes_latest
    WHERE originator_node_id = @originator_node_id
), 0)::BIGINT AS originator_sequence_id;
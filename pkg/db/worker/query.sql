
-- name: ListPartitions :many
SELECT table_name
FROM information_schema.tables
WHERE table_name LIKE 'gateway_envelopes_meta%'
ORDER BY table_name;

-- name: GetLastSequenceID :one
SELECT MAX(originator_sequence_id)
FROM sqlc.narg('table_name');

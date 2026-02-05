
-- name: ListPartitions :many
SELECT table_name
FROM information_schema.tables
WHERE table_name ~ '^gateway_envelopes_meta_o\d+_s\d+_\d+$'
ORDER BY table_name;

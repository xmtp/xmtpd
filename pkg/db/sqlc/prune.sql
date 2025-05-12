-- name: CountExpiredEnvelopes :one
SELECT COUNT(*)::bigint AS expired_count
FROM public.gateway_envelopes
WHERE expiry IS NOT NULL
  AND expiry < EXTRACT(EPOCH FROM now())::bigint;

-- name: DeleteExpiredEnvelopesBatch :many
WITH to_delete AS (
    SELECT originator_node_id, originator_sequence_id
    FROM gateway_envelopes
    WHERE expiry IS NOT NULL
      AND expiry < EXTRACT(EPOCH FROM now())::bigint
    ORDER BY expiry
    LIMIT 1000
        FOR UPDATE SKIP LOCKED
)
DELETE FROM gateway_envelopes ge
    USING to_delete td
WHERE ge.originator_node_id = td.originator_node_id
  AND ge.originator_sequence_id = td.originator_sequence_id
RETURNING ge.originator_node_id, ge.originator_sequence_id, ge.expiry;

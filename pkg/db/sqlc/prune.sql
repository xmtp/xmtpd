-- name: CountExpiredEnvelopes :one
WITH max_prunable AS (SELECT originator_node_id,
                             COALESCE(MAX(end_sequence_id), 0) AS max_end_sequence_id
                      FROM payer_reports
                      WHERE submission_status = 1 -- SUBMITTED
                         OR submission_status = 2 -- SETTLED
                      GROUP BY originator_node_id)
SELECT COUNT(*)::bigint AS expired_count
FROM gateway_envelopes_meta ge
         LEFT JOIN max_prunable mp
                   ON ge.originator_node_id = mp.originator_node_id
WHERE ge.expiry IS NOT NULL
  AND ge.expiry < EXTRACT(EPOCH FROM now())::bigint
  AND ge.originator_sequence_id <= COALESCE(mp.max_end_sequence_id, 0);

-- name: DeleteExpiredEnvelopesBatch :many
WITH max_prunable AS (SELECT originator_node_id,
                             COALESCE(MAX(end_sequence_id), 0) AS max_end_sequence_id
                      FROM payer_reports
                      WHERE submission_status = 1 -- SUBMITTED
                         OR submission_status = 2 -- SETTLED
                      GROUP BY originator_node_id),
     to_delete AS (SELECT ge.originator_node_id,
                          ge.originator_sequence_id
                   FROM gateway_envelopes_meta ge
                            LEFT JOIN max_prunable mp
                                      ON ge.originator_node_id = mp.originator_node_id
                   WHERE ge.expiry IS NOT NULL
                     AND ge.expiry < EXTRACT(EPOCH FROM now())::bigint
                     AND ge.originator_sequence_id <= COALESCE(mp.max_end_sequence_id, 0)
                   ORDER BY ge.expiry, ge.originator_node_id, ge.originator_sequence_id
                   LIMIT @batch_size FOR UPDATE SKIP LOCKED)
DELETE
FROM gateway_envelopes_meta ge
    USING to_delete td
WHERE ge.originator_node_id = td.originator_node_id
  AND ge.originator_sequence_id = td.originator_sequence_id
RETURNING ge.originator_node_id, ge.originator_sequence_id, ge.expiry;

-- name: CountExpiredMigratedEnvelopes :one
SELECT COUNT(*)::bigint AS expired_count
FROM gateway_envelopes_meta
WHERE expiry IS NOT NULL
  AND expiry < EXTRACT(EPOCH FROM now())::bigint
  AND originator_node_id BETWEEN 10 AND 14;

-- name: DeleteExpiredMigratedEnvelopesBatch :many
WITH to_delete AS (SELECT originator_node_id,
                          originator_sequence_id
                   FROM gateway_envelopes_meta
                   WHERE expiry IS NOT NULL
                     AND expiry < EXTRACT(EPOCH FROM now())::bigint
                     AND originator_node_id BETWEEN 10 AND 14
                   ORDER BY expiry, originator_node_id, originator_sequence_id
                   LIMIT @batch_size FOR UPDATE SKIP LOCKED)
DELETE
FROM gateway_envelopes_meta ge
    USING to_delete td
WHERE ge.originator_node_id = td.originator_node_id
  AND ge.originator_sequence_id = td.originator_sequence_id
RETURNING ge.originator_node_id, ge.originator_sequence_id, ge.expiry;

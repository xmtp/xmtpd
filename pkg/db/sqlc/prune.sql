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

-- name: GetPrunableCeiling :many
SELECT originator_node_id,
       COALESCE(MAX(end_sequence_id), 0)::bigint AS max_end_sequence_id
FROM payer_reports
WHERE submission_status = 1 -- SUBMITTED
   OR submission_status = 2 -- SETTLED
GROUP BY originator_node_id;

-- name: GetPrunableMetaPartitions :many
SELECT
    p.originator_node_id::int   AS originator_node_id,
    p.schemaname::text           AS schemaname,
    p.tablename::text            AS tablename,
    p.band_start::int           AS band_start,
    p.band_end::int            AS band_end
FROM get_prunable_meta_partitions() AS p(
    originator_node_id,
    schemaname,
    tablename,
    band_start,
    band_end
    )
ORDER BY p.originator_node_id, p.band_start;
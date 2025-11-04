-- name: InsertGatewayEnvelopeV2 :one
WITH m AS (
    INSERT INTO gateway_envelopes_meta_v2 (
                                           originator_node_id,
                                           originator_sequence_id,
                                           topic,
                                           payer_id,
                                           gateway_time,
                                           expiry
        )
        VALUES (@originator_node_id,
                @originator_sequence_id,
                @topic,
                @payer_id,
                COALESCE(@gateway_time, NOW()),
                @expiry)
        ON CONFLICT DO NOTHING
        RETURNING 1),
     b AS (
         INSERT INTO gateway_envelope_blobs_v2 (
                                                originator_node_id,
                                                originator_sequence_id,
                                                originator_envelope
             )
             VALUES (@originator_node_id,
                     @originator_sequence_id,
                     @originator_envelope)
             ON CONFLICT DO NOTHING
             RETURNING 1)
SELECT (SELECT COUNT(*) FROM m)                            AS inserted_meta_rows,
       (SELECT COUNT(*) FROM b)                            AS inserted_blob_rows,
       (SELECT COUNT(*) FROM m) + (SELECT COUNT(*) FROM b) AS total_inserted_rows;


-- name: SelectNewestFromTopicsV2 :many
WITH latest AS (SELECT DISTINCT ON (m.topic) m.originator_node_id,
                                             m.originator_sequence_id,
                                             m.gateway_time,
                                             m.topic
                FROM gateway_envelopes_meta_v2 m
                WHERE m.topic = ANY (@topics::BYTEA[])
                ORDER BY m.topic, m.gateway_time DESC)
SELECT l.originator_node_id,
       l.originator_sequence_id,
       l.gateway_time,
       l.topic,
       b.originator_envelope
FROM latest l
         JOIN gateway_envelope_blobs_v2 b
              ON b.originator_node_id = l.originator_node_id
                  AND b.originator_sequence_id = l.originator_sequence_id
ORDER BY l.topic;

-- name: SelectGatewayEnvelopesV2ByOriginators :many
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
             JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                  USING (ord)
)
SELECT v.originator_node_id,
       v.originator_sequence_id,
       v.gateway_time,
       v.topic,
       v.originator_envelope
FROM gateway_envelopes_v2_view v
         LEFT JOIN cursors c
                   ON v.originator_node_id = c.cursor_node_id
WHERE v.originator_node_id = ANY(@originator_node_ids::INT[])
  AND v.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
ORDER BY v.gateway_time, v.originator_node_id, v.originator_sequence_id
LIMIT NULLIF(@row_limit::INT, 0);

-- name: SelectGatewayEnvelopesV2ByTopics :many
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
             JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                  USING (ord)
)
SELECT v.originator_node_id,
       v.originator_sequence_id,
       v.gateway_time,
       v.topic,
       v.originator_envelope
FROM gateway_envelopes_v2_view v
         LEFT JOIN cursors c
                   ON v.originator_node_id = c.cursor_node_id
WHERE v.topic = ANY(@topics::BYTEA[])
  AND v.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
ORDER BY v.gateway_time, v.originator_node_id, v.originator_sequence_id
LIMIT NULLIF(@row_limit::INT, 0);

-- name: SelectGatewayEnvelopesV2Unfiltered :many
WITH cursors AS (SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
                 FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
                          JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                               USING (ord))
SELECT v.originator_node_id,
       v.originator_sequence_id,
       v.gateway_time,
       v.topic,
       v.originator_envelope
FROM gateway_envelopes_v2_view v
         LEFT JOIN cursors c
                   ON v.originator_node_id = c.cursor_node_id
WHERE v.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
ORDER BY v.gateway_time,
         v.originator_node_id,
         v.originator_sequence_id
LIMIT NULLIF(@row_limit::INT, 0);
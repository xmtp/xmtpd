-- name: InsertGatewayEnvelope :one
WITH m AS (
    INSERT INTO gateway_envelopes_meta (
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
         INSERT INTO gateway_envelope_blobs (
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


-- TODO(mkysel) -- sorting by gateway time can lead to wrong results, this query needs to be redone
-- name: SelectNewestFromTopics :many
WITH latest AS (SELECT DISTINCT ON (m.topic) m.originator_node_id,
                                             m.originator_sequence_id,
                                             m.gateway_time,
                                             m.topic
                FROM gateway_envelopes_meta m
                WHERE m.topic = ANY (@topics::BYTEA[])
                ORDER BY m.topic, m.gateway_time DESC)
SELECT l.originator_node_id,
       l.originator_sequence_id,
       l.gateway_time,
       l.topic,
       b.originator_envelope
FROM latest l
         JOIN gateway_envelope_blobs b
              ON b.originator_node_id = l.originator_node_id
                  AND b.originator_sequence_id = l.originator_sequence_id
ORDER BY l.topic;

-- name: SelectGatewayEnvelopesByOriginators :many
WITH cursors AS (
    SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
filtered AS (
    SELECT m.originator_node_id,
           m.originator_sequence_id,
           m.gateway_time,
           m.topic
    FROM gateway_envelopes_meta AS m
    LEFT JOIN cursors AS c
        ON m.originator_node_id = c.cursor_node_id
    WHERE m.originator_node_id = ANY (@originator_node_ids::INT[])
    AND m.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
    ORDER BY m.originator_node_id, m.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
)
SELECT f.originator_node_id,
       f.originator_sequence_id,
       f.gateway_time,
       f.topic,
       b.originator_envelope
FROM filtered AS f
JOIN gateway_envelope_blobs AS b
    ON b.originator_node_id = f.originator_node_id
    AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id;

-- name: SelectGatewayEnvelopesByTopics :many
WITH cursors AS (SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
                 FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
                          JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                               USING (ord)),
     filtered AS (
         -- A) topic + cursor
         SELECT m.originator_node_id,
                m.originator_sequence_id,
                m.gateway_time,
                m.topic
         FROM gateway_envelopes_meta AS m
                  JOIN cursors AS c
                       ON m.originator_node_id = c.cursor_node_id
                           AND m.originator_sequence_id > c.cursor_sequence_id
         WHERE m.topic = ANY (@topics::BYTEA[])

         UNION ALL

         -- B) topic + no-cursor
         SELECT m.originator_node_id,
                m.originator_sequence_id,
                m.gateway_time,
                m.topic
         FROM gateway_envelopes_meta AS m
         WHERE m.topic = ANY (@topics::BYTEA[])
           AND m.originator_sequence_id > 0
           AND NOT EXISTS (SELECT 1
                           FROM cursors AS c
                           WHERE c.cursor_node_id = m.originator_node_id)

         -- Do the ordering/limit on meta rows before touching blobs
         ORDER BY originator_node_id, originator_sequence_id
         LIMIT NULLIF(@row_limit::INT, 0))
SELECT f.originator_node_id,
       f.originator_sequence_id,
       f.gateway_time,
       f.topic,
       b.originator_envelope
FROM filtered AS f
         JOIN gateway_envelope_blobs AS b
              ON b.originator_node_id = f.originator_node_id
                  AND b.originator_sequence_id = f.originator_sequence_id
ORDER BY f.originator_node_id, f.originator_sequence_id;


-- name: SelectGatewayEnvelopesUnfiltered :many
WITH cursors AS (SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
                 FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
                          JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                               USING (ord))
SELECT v.originator_node_id,
       v.originator_sequence_id,
       v.gateway_time,
       v.topic,
       v.originator_envelope
FROM gateway_envelopes_view v
         LEFT JOIN cursors c
                   ON v.originator_node_id = c.cursor_node_id
WHERE v.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
ORDER BY v.originator_node_id,
         v.originator_sequence_id
LIMIT NULLIF(@row_limit::INT, 0);
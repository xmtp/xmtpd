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
WITH cursors AS (SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
                 FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
                          JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
                               USING (ord)),
     filtered AS (SELECT m.originator_node_id,
                         m.originator_sequence_id,
                         m.gateway_time,
                         m.topic
                  FROM gateway_envelopes_meta AS m
                           LEFT JOIN cursors AS c
                                     ON m.originator_node_id = c.cursor_node_id
                  WHERE m.originator_node_id = ANY (@originator_node_ids::INT[])
                    AND m.originator_sequence_id > COALESCE(c.cursor_sequence_id, 0)
                  ORDER BY m.originator_node_id, m.originator_sequence_id
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

-- name: InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage :one
WITH input AS (
  SELECT
    a.originator_node_id,
    b.originator_sequence_id,
    c.topic,
    d.payer_id,
    e.gateway_time,
    f.expiry,
    g.originator_envelope,
    h.spend_picodollars
  FROM unnest(@originator_node_ids::int[]) WITH ORDINALITY AS a(originator_node_id, ord)
  JOIN unnest(@originator_sequence_ids::bigint[]) WITH ORDINALITY AS b(originator_sequence_id, ord) USING (ord)
  JOIN unnest(@topics::bytea[]) WITH ORDINALITY AS c(topic, ord) USING (ord)
  JOIN unnest(@payer_ids::int[]) WITH ORDINALITY AS d(payer_id, ord) USING (ord)
  JOIN unnest(@gateway_times::timestamp[]) WITH ORDINALITY AS e(gateway_time, ord) USING (ord)
  JOIN unnest(@expiries::bigint[]) WITH ORDINALITY AS f(expiry, ord) USING (ord)
  JOIN unnest(@originator_envelopes::bytea[]) WITH ORDINALITY AS g(originator_envelope, ord) USING (ord)
  JOIN unnest(@spend_picodollars::bigint[]) WITH ORDINALITY AS h(spend_picodollars, ord) USING (ord)
),

m AS (
  INSERT INTO gateway_envelopes_meta (
    originator_node_id,
    originator_sequence_id,
    topic,
    payer_id,
    gateway_time,
    expiry
  )
  SELECT originator_node_id, originator_sequence_id, topic, payer_id, gateway_time, expiry
  FROM input
  ON CONFLICT DO NOTHING
  RETURNING originator_node_id, originator_sequence_id, payer_id, gateway_time
),

b AS (
  INSERT INTO gateway_envelope_blobs (
    originator_node_id,
    originator_sequence_id,
    originator_envelope
  )
  SELECT originator_node_id, originator_sequence_id, originator_envelope
  FROM input
  ON CONFLICT DO NOTHING
  RETURNING originator_node_id, originator_sequence_id
),

m_with_spend AS (
  SELECT
    m.originator_node_id,
    m.originator_sequence_id,
    m.payer_id,
    m.gateway_time,
    i.spend_picodollars
  FROM m
  JOIN b USING (originator_node_id, originator_sequence_id)
  JOIN input i USING (originator_node_id, originator_sequence_id)
),

u_prep AS (
  SELECT
    payer_id,
    originator_node_id AS originator_id,
    floor(extract(epoch from gateway_time) / 60)::int AS minutes_since_epoch,
    sum(spend_picodollars)::bigint AS spend_picodollars,
    max(originator_sequence_id)::bigint AS last_sequence_id,
    count(*)::int AS message_count
  FROM m_with_spend
  WHERE payer_id IS NOT NULL
  GROUP BY 1,2,3
),

u AS (
  INSERT INTO unsettled_usage (
    payer_id,
    originator_id,
    minutes_since_epoch,
    spend_picodollars,
    last_sequence_id,
    message_count
  )
  SELECT payer_id, originator_id, minutes_since_epoch, spend_picodollars, last_sequence_id, message_count
  FROM u_prep
  ON CONFLICT (payer_id, originator_id, minutes_since_epoch) DO UPDATE
  SET
    spend_picodollars = unsettled_usage.spend_picodollars + EXCLUDED.spend_picodollars,
    message_count     = unsettled_usage.message_count + EXCLUDED.message_count,
    last_sequence_id  = GREATEST(unsettled_usage.last_sequence_id, EXCLUDED.last_sequence_id)
  RETURNING 1
)

SELECT
  (SELECT COUNT(*) FROM m) AS inserted_meta_rows,
  (SELECT COUNT(*) FROM b) AS inserted_blob_rows,
  (SELECT COUNT(*) FROM u) AS affected_usage_rows;



-- name: InsertGatewayEnvelopeV3 :one
-- Single-row envelope insert targeting the renamed gateway_envelopes_blobs
-- table. The pre-rename version (which referenced gateway_envelope_blobs in
-- inline SQL) cannot survive at HEAD because sqlc validates inline-SQL
-- query bodies against the live schema; the renamed table no longer exists
-- under the old name. Migration-behavior tests instead use
-- InsertGatewayEnvelopeBatchV2, which goes through the v2 stored function
-- (whose body is opaque to sqlc validation).
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
         INSERT INTO gateway_envelopes_blobs (
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
         JOIN gateway_envelopes_blobs b
              ON b.originator_node_id = l.originator_node_id
                  AND b.originator_sequence_id = l.originator_sequence_id
ORDER BY l.topic;

-- name: SelectGatewayEnvelopesBySingleOriginator :many
-- Optimized query for a single originator - uses direct index scan
SELECT m.originator_node_id,
       m.originator_sequence_id,
       m.gateway_time,
       m.topic,
       b.originator_envelope
FROM gateway_envelopes_meta AS m
JOIN gateway_envelopes_blobs AS b
    ON b.originator_node_id = m.originator_node_id
   AND b.originator_sequence_id = m.originator_sequence_id
   AND b.originator_node_id = @originator_node_id::INT
WHERE m.originator_node_id = @originator_node_id::INT
  AND m.originator_sequence_id > @cursor_sequence_id::BIGINT
ORDER BY m.originator_sequence_id
LIMIT NULLIF(@row_limit::INT, 0);

-- name: SelectGatewayEnvelopesByOriginators :many
-- LATERAL per originator with per-originator blob join.
-- Requires callers to include all desired originators in cursor arrays (use seq_id=0 for unseen).
WITH cursor_entries AS (
    SELECT x.node_id AS node_id, y.seq_id AS seq_id
    FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
    JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
    USING (ord)
),
filtered AS (
    SELECT sub.originator_node_id,
           sub.originator_sequence_id,
           sub.gateway_time,
           sub.topic
    FROM cursor_entries AS ce
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id,
               m.originator_sequence_id,
               m.gateway_time,
               m.topic
        FROM gateway_envelopes_meta AS m
        WHERE m.originator_node_id = ce.node_id
          AND m.originator_sequence_id > ce.seq_id
        ORDER BY m.originator_sequence_id
        LIMIT NULLIF(@row_limit::INT, 0)
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
),
originator_ids AS (
    SELECT DISTINCT originator_node_id FROM filtered
)
SELECT bl.originator_node_id,
       bl.originator_sequence_id,
       bl.gateway_time,
       bl.topic,
       bl.originator_envelope
FROM originator_ids AS oi
CROSS JOIN LATERAL (
    SELECT f.originator_node_id,
           f.originator_sequence_id,
           f.gateway_time,
           f.topic,
           b.originator_envelope
    FROM filtered AS f
    JOIN gateway_envelopes_blobs AS b
        ON b.originator_node_id = oi.originator_node_id
       AND b.originator_sequence_id = f.originator_sequence_id
    WHERE f.originator_node_id = oi.originator_node_id
) AS bl
ORDER BY bl.originator_node_id, bl.originator_sequence_id;

-- name: SelectGatewayEnvelopesByTopics :many
-- V3b LATERAL per (topic, originator) with per-originator blob join.
-- Requires callers to include ALL originators in cursor arrays (use seq_id=0 for unseen).
-- Uses gem_topic_orig_seq_idx for index-only scans.
-- row_limit is required and caps total rows returned.
WITH cursors AS (
	SELECT x.node_id AS cursor_node_id, y.seq_id AS cursor_sequence_id
	FROM unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS x(node_id, ord)
	JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS y(seq_id, ord)
	USING (ord)
),
cursor_entries AS (
	SELECT t.topic, c.cursor_node_id AS node_id, c.cursor_sequence_id AS seq_id
	FROM unnest(@topics::BYTEA[]) AS t(topic)
	CROSS JOIN cursors AS c
),
filtered AS (
	SELECT sub.originator_node_id,
	       sub.originator_sequence_id,
	       sub.gateway_time,
	       sub.topic
	FROM cursor_entries AS ce
	CROSS JOIN LATERAL (
		SELECT m.originator_node_id,
		       m.originator_sequence_id,
		       m.gateway_time,
		       m.topic
		FROM gateway_envelopes_meta AS m
		WHERE m.topic = ce.topic
		  AND m.originator_node_id = ce.node_id
		  AND m.originator_sequence_id > ce.seq_id
		ORDER BY m.originator_sequence_id
		LIMIT @row_limit::INT
	) AS sub
	ORDER BY sub.originator_node_id, sub.originator_sequence_id
	LIMIT @row_limit::INT
),
originator_ids AS (
	SELECT DISTINCT originator_node_id FROM filtered
)
SELECT bl.originator_node_id,
       bl.originator_sequence_id,
       bl.gateway_time,
       bl.topic,
       bl.originator_envelope
FROM originator_ids AS oi
CROSS JOIN LATERAL (
	SELECT f.originator_node_id,
	       f.originator_sequence_id,
	       f.gateway_time,
	       f.topic,
	       b.originator_envelope
	FROM filtered AS f
	JOIN gateway_envelopes_blobs AS b
	    ON b.originator_node_id = oi.originator_node_id
	   AND b.originator_sequence_id = f.originator_sequence_id
	WHERE f.originator_node_id = oi.originator_node_id
) AS bl
ORDER BY bl.originator_node_id, bl.originator_sequence_id;


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

-- name: SelectGatewayEnvelopesByPerTopicCursors :many
-- Per-topic cursor variant: accepts pre-flattened (topic, node_id, seq_id) triples
-- instead of CROSS JOINing topics × cursors.
-- Uses gem_topic_orig_seq_idx for index-only scans.
WITH cursor_entries AS (
    SELECT t.topic, n.node_id, s.seq_id
    FROM unnest(@cursor_topics::BYTEA[]) WITH ORDINALITY AS t(topic, ord)
    JOIN unnest(@cursor_node_ids::INT[]) WITH ORDINALITY AS n(node_id, ord) USING (ord)
    JOIN unnest(@cursor_sequence_ids::BIGINT[]) WITH ORDINALITY AS s(seq_id, ord) USING (ord)
),
filtered AS (
    SELECT sub.originator_node_id,
           sub.originator_sequence_id,
           sub.gateway_time,
           sub.topic
    FROM cursor_entries AS ce
    CROSS JOIN LATERAL (
        SELECT m.originator_node_id,
               m.originator_sequence_id,
               m.gateway_time,
               m.topic
        FROM gateway_envelopes_meta AS m
        WHERE m.topic = ce.topic
          AND m.originator_node_id = ce.node_id
          AND m.originator_sequence_id > ce.seq_id
        ORDER BY m.originator_sequence_id
        LIMIT @rows_per_entry::INT
    ) AS sub
    ORDER BY sub.originator_node_id, sub.originator_sequence_id
    LIMIT NULLIF(@row_limit::INT, 0)
),
originator_ids AS (
    SELECT DISTINCT originator_node_id FROM filtered
)
SELECT bl.originator_node_id,
       bl.originator_sequence_id,
       bl.gateway_time,
       bl.topic,
       bl.originator_envelope
FROM originator_ids AS oi
CROSS JOIN LATERAL (
    SELECT f.originator_node_id,
           f.originator_sequence_id,
           f.gateway_time,
           f.topic,
           b.originator_envelope
    FROM filtered AS f
    JOIN gateway_envelopes_blobs AS b
        ON b.originator_node_id = oi.originator_node_id
       AND b.originator_sequence_id = f.originator_sequence_id
    WHERE f.originator_node_id = oi.originator_node_id
) AS bl
ORDER BY bl.originator_node_id, bl.originator_sequence_id;

-- name: InsertGatewayEnvelopeBatchV2 :one
-- Pre-rename batch insert. Calls the v2 stored function which still
-- references the legacy gateway_envelope_blobs table. No production code
-- uses this query — it survives only so migration-behavior tests can
-- populate the database at the pre-rename schema version. Production code
-- uses InsertGatewayEnvelopeBatchV3.
SELECT
    inserted_meta_rows::bigint,
    inserted_blob_rows::bigint,
    affected_usage_rows::bigint,
    affected_congestion_rows::bigint
FROM insert_gateway_envelope_batch_v2(
    @originator_node_ids::int[],
    @originator_sequence_ids::bigint[],
    @topics::bytea[],
    @payer_ids::int[],
    @gateway_times::timestamp[],
    @expiries::bigint[],
    @originator_envelopes::bytea[],
    @spend_picodollars::bigint[],
    @count_usage::boolean[],
    @count_congestion::boolean[]
);

-- name: InsertGatewayEnvelopeBatchV3 :one
-- Batch envelope insert calling the renamed insert_gateway_envelope_batch_v3
-- stored function (which targets gateway_envelopes_blobs). This is the
-- production version.
SELECT
    inserted_meta_rows::bigint,
    inserted_blob_rows::bigint,
    affected_usage_rows::bigint,
    affected_congestion_rows::bigint
FROM insert_gateway_envelope_batch_v3(
    @originator_node_ids::int[],
    @originator_sequence_ids::bigint[],
    @topics::bytea[],
    @payer_ids::int[],
    @gateway_times::timestamp[],
    @expiries::bigint[],
    @originator_envelopes::bytea[],
    @spend_picodollars::bigint[],
    @count_usage::boolean[],
    @count_congestion::boolean[]
);

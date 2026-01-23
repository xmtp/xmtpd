CREATE OR REPLACE FUNCTION insert_gateway_envelope_batch(
    p_originator_node_ids     int[],
    p_originator_sequence_ids bigint[],
    p_topics                  bytea[],
    p_payer_ids               int[],
    p_gateway_times           timestamp[],
    p_expiries                bigint[],
    p_originator_envelopes    bytea[],
    p_spend_picodollars       bigint[]
)
RETURNS TABLE (
    inserted_meta_rows  bigint,
    inserted_blob_rows  bigint,
    affected_usage_rows bigint
)
LANGUAGE SQL
AS $$
WITH

input_meta AS MATERIALIZED (
    SELECT
        originator_node_id,
        originator_sequence_id,
        topic,
        NULLIF(payer_id, 0) AS payer_id,
        gateway_time,
        expiry,
        spend_picodollars
    FROM unnest(
        p_originator_node_ids,
        p_originator_sequence_ids,
        p_topics,
        p_payer_ids,
        p_gateway_times,
        p_expiries,
        p_spend_picodollars
    ) AS t(
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry,
        spend_picodollars
    )
),

-- Wide input (payload only). Do NOT materialize.
input_blob AS (
    SELECT
        originator_node_id,
        originator_sequence_id,
        originator_envelope
    FROM unnest(
        p_originator_node_ids,
        p_originator_sequence_ids,
        p_originator_envelopes
    ) AS t(
        originator_node_id,
        originator_sequence_id,
        originator_envelope
    )
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
    SELECT
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry
    FROM input_meta
    ON CONFLICT DO NOTHING
    RETURNING originator_node_id, originator_sequence_id, payer_id, gateway_time
),

-- Drive blob insert from m to avoid work for rows conflict-skipped in meta.
b AS (
    INSERT INTO gateway_envelope_blobs (
        originator_node_id,
        originator_sequence_id,
        originator_envelope
    )
    SELECT
        m.originator_node_id,
        m.originator_sequence_id,
        ib.originator_envelope
    FROM m
    JOIN input_blob ib
      USING (originator_node_id, originator_sequence_id)
    ON CONFLICT DO NOTHING
    RETURNING originator_node_id, originator_sequence_id
),

m_with_spend AS (
    SELECT
        m.originator_node_id,
        m.originator_sequence_id,
        m.payer_id,
        m.gateway_time,
        im.spend_picodollars
    FROM m
    JOIN b
      USING (originator_node_id, originator_sequence_id)
    JOIN input_meta im
      USING (originator_node_id, originator_sequence_id)
),

u_prep AS (
    SELECT
        payer_id,
        originator_node_id AS originator_id,
        (extract(epoch from gateway_time)::bigint / 60)::int AS minutes_since_epoch,
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
    SELECT
        payer_id,
        originator_id,
        minutes_since_epoch,
        spend_picodollars,
        last_sequence_id,
        message_count
    FROM u_prep
    ORDER BY payer_id, originator_id, minutes_since_epoch
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
$$;

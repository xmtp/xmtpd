CREATE OR REPLACE FUNCTION insert_gateway_envelope_batch_v2(
    p_originator_node_ids     int[],
    p_originator_sequence_ids bigint[],
    p_topics                  bytea[],
    p_payer_ids               int[],
    p_gateway_times           timestamp[],
    p_expiries                bigint[],
    p_originator_envelopes    bytea[],
    p_spend_picodollars       bigint[],
    p_count_usage             boolean[],
    p_count_congestion        boolean[]
)
RETURNS TABLE (
    inserted_meta_rows       bigint,
    inserted_blob_rows       bigint,
    affected_usage_rows      bigint,
    affected_congestion_rows bigint
)
LANGUAGE SQL
AS $$
WITH input AS (
    SELECT
        originator_node_id,
        originator_sequence_id,
        topic,
        NULLIF(payer_id, 0) AS payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars,
        count_usage,
        count_congestion
    FROM unnest(
        p_originator_node_ids,
        p_originator_sequence_ids,
        p_topics,
        p_payer_ids,
        p_gateway_times,
        p_expiries,
        p_originator_envelopes,
        p_spend_picodollars,
        p_count_usage,
        p_count_congestion
    ) AS t(
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars,
        count_usage,
        count_congestion
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
        i.spend_picodollars,
        i.count_usage,
        i.count_congestion
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
    WHERE payer_id IS NOT NULL AND count_usage
    GROUP BY 1, 2, 3
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
    ORDER BY payer_id, originator_id, minutes_since_epoch
    ON CONFLICT (payer_id, originator_id, minutes_since_epoch) DO UPDATE
    SET
        spend_picodollars = unsettled_usage.spend_picodollars + EXCLUDED.spend_picodollars,
        message_count     = unsettled_usage.message_count + EXCLUDED.message_count,
        last_sequence_id  = GREATEST(unsettled_usage.last_sequence_id, EXCLUDED.last_sequence_id)
    RETURNING 1
),

c_prep AS (
    SELECT
        originator_node_id AS originator_id,
        floor(extract(epoch from gateway_time) / 60)::int AS minutes_since_epoch,
        count(*)::int AS num_messages
    FROM m_with_spend
    WHERE count_congestion
    GROUP BY 1, 2
),

c AS (
    INSERT INTO originator_congestion (originator_id, minutes_since_epoch, num_messages)
    SELECT originator_id, minutes_since_epoch, num_messages
    FROM c_prep
    ORDER BY originator_id, minutes_since_epoch
    ON CONFLICT (originator_id, minutes_since_epoch) DO UPDATE
    SET num_messages = originator_congestion.num_messages + EXCLUDED.num_messages
    RETURNING 1
)

SELECT
    (SELECT COUNT(*) FROM m) AS inserted_meta_rows,
    (SELECT COUNT(*) FROM b) AS inserted_blob_rows,
    (SELECT COUNT(*) FROM u) AS affected_usage_rows,
    (SELECT COUNT(*) FROM c) AS affected_congestion_rows;
$$;

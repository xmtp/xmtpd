CREATE TYPE gateway_envelope_row AS (
    originator_node_id     int,
    originator_sequence_id bigint,
    topic                  bytea,
    payer_id               int,
    gateway_time           timestamp,
    expiry                 bigint,
    originator_envelope    bytea,
    spend_picodollars      bigint
);

-- Create the batch insert function
CREATE OR REPLACE FUNCTION insert_gateway_envelope_batch(
    envelopes gateway_envelope_row[]
)
RETURNS TABLE (
    inserted_meta_rows  bigint,
    inserted_blob_rows  bigint,
    affected_usage_rows bigint
)
LANGUAGE SQL
AS $$
WITH input AS (
    SELECT
        e.originator_node_id,
        e.originator_sequence_id,
        e.topic,
        NULLIF(e.payer_id, 0) AS payer_id,
        e.gateway_time,
        e.expiry,
        e.originator_envelope,
        e.spend_picodollars
    FROM unnest(envelopes) AS e
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

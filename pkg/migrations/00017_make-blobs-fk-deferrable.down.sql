-- Revert to non-deferrable FK constraints and original function

-- Step 1: Revert all existing leaf partitions to non-deferrable FK constraints
DO $$
DECLARE
    leaf_partition RECORD;
    constraint_name TEXT;
BEGIN
    -- Find all leaf partitions of gateway_envelope_blobs (2 levels deep: list -> range)
    FOR leaf_partition IN
        SELECT c.relname AS partition_name
        FROM pg_inherits i1
        JOIN pg_class c1 ON i1.inhrelid = c1.oid
        JOIN pg_inherits i2 ON c1.oid = i2.inhparent
        JOIN pg_class c ON i2.inhrelid = c.oid
        WHERE i1.inhparent = 'gateway_envelope_blobs'::regclass
    LOOP
        -- Find the FK constraint name on this partition
        SELECT con.conname INTO constraint_name
        FROM pg_constraint con
        JOIN pg_class rel ON con.conrelid = rel.oid
        WHERE rel.relname = leaf_partition.partition_name
          AND con.contype = 'f'
          AND con.confrelid IN (
              SELECT oid FROM pg_class WHERE relname LIKE 'gateway_envelopes_meta%'
          );

        IF constraint_name IS NOT NULL THEN
            -- Drop the deferrable constraint
            EXECUTE format(
                'ALTER TABLE %I DROP CONSTRAINT IF EXISTS %I',
                leaf_partition.partition_name,
                constraint_name
            );

            -- Recreate as non-deferrable
            EXECUTE format(
                'ALTER TABLE %I ADD CONSTRAINT %I
                 FOREIGN KEY (originator_node_id, originator_sequence_id)
                 REFERENCES gateway_envelopes_meta(originator_node_id, originator_sequence_id)
                 ON DELETE CASCADE',
                leaf_partition.partition_name,
                leaf_partition.partition_name || '_meta_fkey'
            );
        END IF;
    END LOOP;
END $$;

-- Step 2: Revert the blob subpartition creation function
CREATE OR REPLACE FUNCTION make_blob_seq_subpart_v2(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    -- gateway_envelope_blobs_oXXX
    parent text := format('gateway_envelope_blobs_o%s', _oid);
    -- gateway_envelope_blobs_oXXX_sN0_N1
    subname       text := format('gateway_envelope_blobs_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    -- Since it's a standalone table - setup a constraint.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelope_blobs INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        )',
        subname,
        _oid::text,
        _start::text,
        _end::text
    );

    EXECUTE format(
        'ALTER TABLE %I ATTACH PARTITION %I
            FOR VALUES FROM (%s) TO (%s)',
        parent,
        subname,
        _start::text,
        _end::text
    );

    -- Now we can drop the constraint.
    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT seq_id_check;',
        subname
    );
EXCEPTION
    WHEN OTHERS THEN
        IF SQLERRM ~ 'is already a partition' THEN
            -- Do nothing.
            NULL;
        ELSE
            RAISE;
        END IF;
END;
$$ LANGUAGE plpgsql;

-- Step 3: Revert to original SQL-based batch insert function
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

WITH input AS (
    SELECT
        originator_node_id,
        originator_sequence_id,
        topic,
        NULLIF(payer_id, 0) AS payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars
    FROM unnest(
        p_originator_node_ids,
        p_originator_sequence_ids,
        p_topics,
        p_payer_ids,
        p_gateway_times,
        p_expiries,
        p_originator_envelopes,
        p_spend_picodollars
    ) AS t(
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars
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

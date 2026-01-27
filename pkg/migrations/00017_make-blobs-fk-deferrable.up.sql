-- Make the FK constraint on gateway_envelope_blobs deferrable for batch performance.
-- This allows PostgreSQL to batch FK validation checks at commit time instead of per-row.

-- Step 1: Update all existing leaf partitions to have deferrable FK constraints
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
            -- Drop the existing non-deferrable constraint
            EXECUTE format(
                'ALTER TABLE %I DROP CONSTRAINT IF EXISTS %I',
                leaf_partition.partition_name,
                constraint_name
            );

            -- Recreate as deferrable
            EXECUTE format(
                'ALTER TABLE %I ADD CONSTRAINT %I
                 FOREIGN KEY (originator_node_id, originator_sequence_id)
                 REFERENCES gateway_envelopes_meta(originator_node_id, originator_sequence_id)
                 ON DELETE CASCADE
                 DEFERRABLE INITIALLY DEFERRED',
                leaf_partition.partition_name,
                leaf_partition.partition_name || '_meta_fkey'
            );
        END IF;
    END LOOP;
END $$;

-- Step 2: Update the blob subpartition creation function to use deferrable FK
CREATE OR REPLACE FUNCTION make_blob_seq_subpart_v2(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    -- gateway_envelope_blobs_oXXX
    parent text := format('gateway_envelope_blobs_o%s', _oid);
    -- gateway_envelope_blobs_oXXX_sN0_N1
    subname       text := format('gateway_envelope_blobs_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    -- Create standalone table WITHOUT inheriting FK (we'll add deferrable one)
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            originator_node_id     int    NOT NULL,
            originator_sequence_id bigint NOT NULL,
            originator_envelope    bytea  NOT NULL,
            PRIMARY KEY (originator_node_id, originator_sequence_id),
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        )',
        subname,
        _start::text,
        _end::text
    );

    -- Add deferrable FK constraint
    EXECUTE format(
        'ALTER TABLE %I ADD CONSTRAINT %I
         FOREIGN KEY (originator_node_id, originator_sequence_id)
         REFERENCES gateway_envelopes_meta(originator_node_id, originator_sequence_id)
         ON DELETE CASCADE
         DEFERRABLE INITIALLY DEFERRED',
        subname,
        subname || '_meta_fkey'
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

-- Step 3: Update the batch insert function to defer constraints
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
LANGUAGE plpgsql
AS $$
BEGIN
    -- Defer FK constraint checks to end of statement for batch efficiency
    SET CONSTRAINTS ALL DEFERRED;

    RETURN QUERY
    WITH input_small AS MATERIALIZED (
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
    input_env AS (
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

    m_ins AS (
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
        FROM input_small
        ON CONFLICT DO NOTHING
        RETURNING originator_node_id, originator_sequence_id, payer_id, gateway_time
    ),

    b_ins AS (
        INSERT INTO gateway_envelope_blobs (
            originator_node_id,
            originator_sequence_id,
            originator_envelope
        )
        SELECT
            m.originator_node_id,
            m.originator_sequence_id,
            e.originator_envelope
        FROM m_ins m
        JOIN input_env e USING (originator_node_id, originator_sequence_id)
        ON CONFLICT DO NOTHING
        RETURNING originator_node_id, originator_sequence_id
    ),

    u_prep AS (
        SELECT
            m.payer_id,
            m.originator_node_id AS originator_id,
            (extract(epoch from m.gateway_time)::bigint / 60)::int AS minutes_since_epoch,
            sum(s.spend_picodollars)::bigint AS spend_picodollars,
            max(m.originator_sequence_id)::bigint AS last_sequence_id,
            count(*)::int AS message_count
        FROM m_ins m
        JOIN b_ins b USING (originator_node_id, originator_sequence_id)
        JOIN input_small s USING (originator_node_id, originator_sequence_id)
        WHERE m.payer_id IS NOT NULL
        GROUP BY 1, 2, 3
    ),

    u AS (
        INSERT INTO unsettled_usage (
            payer_id, originator_id, minutes_since_epoch, spend_picodollars, last_sequence_id, message_count
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
        (SELECT COUNT(*) FROM m_ins) AS inserted_meta_rows,
        (SELECT COUNT(*) FROM b_ins) AS inserted_blob_rows,
        (SELECT COUNT(*) FROM u)     AS affected_usage_rows;
END;
$$;

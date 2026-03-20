-- Remove orphan blobs: blobs that have no matching meta row.
-- This can happen when the ON DELETE CASCADE FK fails to propagate due to the
-- partition attachment order in ensure_gateway_parts_v2 (meta partitions are
-- attached before blob partitions, leaving the FK state inconsistent).
DELETE FROM gateway_envelope_blobs b
WHERE NOT EXISTS (
    SELECT 1
    FROM gateway_envelopes_meta m
    WHERE m.originator_node_id     = b.originator_node_id
      AND m.originator_sequence_id = b.originator_sequence_id
);

-- Ensure the FK constraint exists on gateway_envelope_blobs.
-- It may have been lost or never properly established during partition management.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'gateway_envelope_blobs'::regclass
          AND contype  = 'f'
    ) THEN
        ALTER TABLE gateway_envelope_blobs
            ADD CONSTRAINT gateway_envelope_blobs_meta_fk
            FOREIGN KEY (originator_node_id, originator_sequence_id)
            REFERENCES gateway_envelopes_meta (originator_node_id, originator_sequence_id)
            ON DELETE CASCADE;
    END IF;
END$$;

-- Fix argument order in make_meta_seq_subpart_v2.
-- Previously _oid and _start were passed instead of _start and _end for the
-- seq_id_check constraint, making it impossible to satisfy (e.g. seq_id >= 100
-- AND seq_id < 0).  The constraint is dropped immediately after ATTACH so there
-- was no functional impact, but correcting it avoids confusion.
CREATE OR REPLACE FUNCTION make_meta_seq_subpart_v2(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    parent  text := format('gateway_envelopes_meta_o%s', _oid);
    subname text := format('gateway_envelopes_meta_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_meta INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        )',
        subname,
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

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT seq_id_check;',
        subname
    );
EXCEPTION
    WHEN OTHERS THEN
        IF SQLERRM ~ 'is already a partition' THEN
            NULL;
        ELSE
            RAISE;
        END IF;
END;
$$ LANGUAGE plpgsql;

-- Fix argument order in make_blob_seq_subpart_v2 (same bug as above).
CREATE OR REPLACE FUNCTION make_blob_seq_subpart_v2(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    parent  text := format('gateway_envelope_blobs_o%s', _oid);
    subname text := format('gateway_envelope_blobs_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelope_blobs INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        )',
        subname,
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

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT seq_id_check;',
        subname
    );
EXCEPTION
    WHEN OTHERS THEN
        IF SQLERRM ~ 'is already a partition' THEN
            NULL;
        ELSE
            RAISE;
        END IF;
END;
$$ LANGUAGE plpgsql;

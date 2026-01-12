-- META: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_meta_originator_part(_oid int)
    RETURNS void AS $$
DECLARE
    -- gateway_envelopes_meta_oXXX
    subname text := format(
        'gateway_envelopes_meta_o%s', _oid
    );
BEGIN
    -- Since it's a standalone table - setup a constraint.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_meta INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT oid_check CHECK (originator_node_id = %s)
        ) PARTITION BY RANGE (originator_sequence_id);
        ',
        subname,
        _oid::text
    );

    EXECUTE format('
        ALTER TABLE gateway_envelopes_meta ATTACH PARTITION %I
            FOR VALUES IN (%s);
        ',
        subname,
        _oid::text
    );

    -- Now we can drop the constraint.
    EXECUTE format('
        ALTER TABLE %I DROP CONSTRAINT oid_check;',
        subname
    );

END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION make_meta_seq_subpart(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    -- gateway_envelopes_meta_oXXX
    parent text := format('gateway_envelopes_meta_o%s', _oid);
    -- gateway_envelopes_meta_oXXX_sN0_N1
    subname       text := format('gateway_envelopes_meta_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    -- Since it's a standalone table - setup a constraint.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_meta INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
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

END;
$$ LANGUAGE plpgsql;


-- BLOBS: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_blob_originator_part(_oid int)
    RETURNS void AS $$
DECLARE
    -- gateway_envelope_blobs_oXXX
    subname text := format(
        'gateway_envelope_blobs_o%s', _oid
    );
BEGIN
    -- Since it's a standalone table - setup a constraint.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelope_blobs INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT oid_check CHECK (originator_node_id = %s)
        ) PARTITION BY RANGE (originator_sequence_id);
        ',
        subname,
        _oid::text
    );

    EXECUTE format('
        ALTER TABLE gateway_envelope_blobs ATTACH PARTITION %I
            FOR VALUES IN (%s);
        ',
        subname,
        _oid::text
    );

    -- Now we can drop the constraint.
    EXECUTE format('
        ALTER TABLE %I DROP CONSTRAINT oid_check;',
        subname
    );

END;
$$ LANGUAGE plpgsql;


-- BLOBS: create a RANGE subpartition [start, end)
CREATE OR REPLACE FUNCTION make_blob_seq_subpart(_oid int, _start bigint, _end bigint)
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

END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ensure_gateway_parts(
    p_originator_node_id     int,
    p_originator_sequence_id bigint,
    p_band_width             bigint DEFAULT 1000000
) RETURNS void LANGUAGE plpgsql AS $$
DECLARE
    v_band_start bigint := (p_originator_sequence_id / p_band_width) * p_band_width;
BEGIN
    PERFORM make_meta_originator_part(p_originator_node_id);
    PERFORM make_blob_originator_part(p_originator_node_id);
    PERFORM make_meta_seq_subpart(p_originator_node_id, v_band_start, v_band_start + p_band_width);
    PERFORM make_blob_seq_subpart(p_originator_node_id, v_band_start, v_band_start + p_band_width);
END$$;

-- META: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_meta_originator_part(_oid int)
    RETURNS void AS $$
BEGIN
    EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES IN (%s) PARTITION BY RANGE (originator_sequence_id)',
            format('gateway_envelopes_meta_o%s', _oid),
            'gateway_envelopes_meta',
            _oid::text
            );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION make_meta_seq_subpart(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    subname       text := format('gateway_envelopes_meta_o%s_s%s_%s', _oid, _start, _end);
    leaf_time_idx text := subname || '_time_node_seq_idx';
    leaf_exp_idx  text := subname || '_expiry_idx';
BEGIN
    EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I
               FOR VALUES FROM (%s) TO (%s)',
            subname,
            format('gateway_envelopes_meta_o%s', _oid),
            _start::text, _end::text
            );
END;
$$ LANGUAGE plpgsql;


-- BLOBS: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_blob_originator_part(_oid int)
    RETURNS void AS $$
BEGIN
    EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES IN (%s) PARTITION BY RANGE (originator_sequence_id)',
            format('gateway_envelope_blobs_o%s', _oid),
            'gateway_envelope_blobs',
            _oid::text
            );
END;
$$ LANGUAGE plpgsql;


-- BLOBS: create a RANGE subpartition [start, end)
CREATE OR REPLACE FUNCTION make_blob_seq_subpart(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    subname text := format('gateway_envelope_blobs_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%s) TO (%s)',
            subname,
            format('gateway_envelope_blobs_o%s', _oid),
            _start::text, _end::text
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
-- META: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_meta_originator_part(_oid int) RETURNS void AS
$$
BEGIN
    BEGIN
        EXECUTE format(
                'CREATE TABLE %I PARTITION OF %I FOR VALUES IN (%s) PARTITION BY RANGE (originator_sequence_id)',
                format('gateway_envelopes_meta_v2_o%s', _oid),
                'gateway_envelopes_meta_v2',
                _oid::text
                );
    EXCEPTION
        WHEN duplicate_table THEN
        -- ok, already created
    END;
END
$$ LANGUAGE plpgsql;

-- META: create a RANGE subpartition [start, end)
CREATE OR REPLACE FUNCTION make_meta_seq_subpart(_oid int, _start bigint, _end bigint) RETURNS void AS
$$
DECLARE
    subname text := format('gateway_envelopes_meta_v2_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    BEGIN
        EXECUTE format(
                'CREATE TABLE %I PARTITION OF %I FOR VALUES FROM (%s) TO (%s)',
                subname,
                format('gateway_envelopes_meta_v2_o%s', _oid),
                _start::text, _end::text
                );
    EXCEPTION
        WHEN duplicate_table THEN
    END;

    -- per-child index
    EXECUTE format(
            'CREATE INDEX IF NOT EXISTS %I ON %I (originator_sequence_id, gateway_time)',
            subname || '_seq_time_idx', subname
            );
END
$$ LANGUAGE plpgsql;

-- BLOBS: create LIST child for one originator, then make it RANGE-partitioned
CREATE OR REPLACE FUNCTION make_blob_originator_part(_oid int) RETURNS void AS
$$
BEGIN
    BEGIN
        EXECUTE format(
                'CREATE TABLE %I PARTITION OF %I FOR VALUES IN (%s) PARTITION BY RANGE (originator_sequence_id)',
                format('gateway_envelope_blobs_v2_o%s', _oid),
                'gateway_envelope_blobs_v2',
                _oid::text
                );
    EXCEPTION
        WHEN duplicate_table THEN
    END;
END
$$ LANGUAGE plpgsql;

-- BLOBS: create a RANGE subpartition [start, end)
CREATE OR REPLACE FUNCTION make_blob_seq_subpart(_oid int, _start bigint, _end bigint) RETURNS void AS
$$
DECLARE
    subname text := format('gateway_envelope_blobs_v2_o%s_s%s_%s', _oid, _start, _end);
BEGIN
    BEGIN
        EXECUTE format(
                'CREATE TABLE %I PARTITION OF %I FOR VALUES FROM (%s) TO (%s)',
                subname,
                format('gateway_envelope_blobs_v2_o%s', _oid),
                _start::text, _end::text
                );
    EXCEPTION
        WHEN duplicate_table THEN
    END;

    -- per-child index for PK lookups
    EXECUTE format(
            'CREATE INDEX IF NOT EXISTS %I ON %I (originator_sequence_id)',
            subname || '_seq_idx', subname
            );
END
$$ LANGUAGE plpgsql;

DO
$$
    DECLARE
        o int;
    BEGIN
        FOREACH o IN ARRAY ARRAY [0,1,10,11,13,100,200,300,400,500,600,700,800,900,1000,1100,1200]
            LOOP
                PERFORM make_meta_originator_part(o);
                PERFORM make_blob_originator_part(o);

                EXECUTE format(
                        'CREATE TABLE IF NOT EXISTS gateway_envelopes_meta_v2_o%s_default
                           PARTITION OF gateway_envelopes_meta_v2_o%s DEFAULT', o, o);
                EXECUTE format(
                        'CREATE INDEX IF NOT EXISTS gem_v2_o%s_def_seq_time
                           ON gateway_envelopes_meta_v2_o%s_default (originator_sequence_id, gateway_time)', o, o);

                EXECUTE format(
                        'CREATE TABLE IF NOT EXISTS gateway_envelope_blobs_v2_o%s_default
                           PARTITION OF gateway_envelope_blobs_v2_o%s DEFAULT', o, o);
                EXECUTE format(
                        'CREATE INDEX IF NOT EXISTS geb_v2_o%s_def_seq
                           ON gateway_envelope_blobs_v2_o%s_default (originator_sequence_id)', o, o);

                PERFORM make_meta_seq_subpart(o, 0, 1000000);
                PERFORM make_meta_seq_subpart(o, 1000000, 2000000);
                PERFORM make_meta_seq_subpart(o, 2000000, 3000000);

                PERFORM make_blob_seq_subpart(o, 0, 1000000);
                PERFORM make_blob_seq_subpart(o, 1000000, 2000000);
                PERFORM make_blob_seq_subpart(o, 2000000, 3000000);
            END LOOP;
    END
$$;

-- Harden gateway-envelopes partition creation against races and a
-- format-string bug in the v2 helpers.
--
-- Motivation (issue #1967):
--   * `make_meta_seq_subpart_v2` builds a CHECK predicate with `format()`,
--     passes four arguments into a string with three placeholders, and relies
--     on PostgreSQL silently dropping the extra. The resulting CHECK is
--     `originator_sequence_id >= _oid AND originator_sequence_id < _start`
--     (e.g. `>= 100 AND < 0` for the first band), not the intended
--     `>= _start AND < _end`. The CHECK is dropped right after ATTACH so the
--     bug is normally benign, but it is still an objective defect.
--   * Every `make_*_part_v*` helper ends in
--         EXCEPTION WHEN OTHERS THEN
--             IF SQLERRM ~ 'is already a partition' THEN NULL;
--             ELSE RAISE; END IF;
--     This regex-matches the PostgreSQL error text and, because the handler
--     sits inside a PL/pgSQL sub-transaction, rolls back the preceding
--     CREATE TABLE together with the failed ATTACH. Any other error whose
--     message happens to contain that substring — or any future change in
--     PostgreSQL's error text — turns into a silent no-op that leaves the
--     partition unattached while the caller sees success.
--   * `ensure_gateway_parts_v3` does not serialize concurrent callers. Two
--     callers racing on the same `(originator_node_id, band_start)` can
--     interleave CREATE / ATTACH / DROP CONSTRAINT, and PostgreSQL's
--     per-statement locks do not guarantee atomicity across the function.
--
-- Fix strategy (append-only; the v2/v3 helpers remain in `pg_proc`):
--   * New `make_*_part_v3`/`_v4` helpers that take a transaction-scoped
--     advisory lock, short-circuit via `pg_inherits` when the partition is
--     already attached, build a CORRECT CHECK predicate, and let any ATTACH
--     error propagate to the caller.
--   * New `ensure_gateway_parts_v4` that wraps the four helpers.
--
-- The legacy helpers are left in place so that migration-behavior tests
-- (e.g. `migration_00023_test.go`) continue to populate pre-rename databases
-- through their existing code paths.

-- META: create LIST child for one originator, idempotently and with no
-- exception swallowing.
CREATE FUNCTION make_meta_originator_part_v3(_oid int)
    RETURNS void AS $$
DECLARE
    subname text := format('gateway_envelopes_meta_o%s', _oid);
    already_attached boolean;
BEGIN
    -- Serialize concurrent callers for this originator on the meta side.
    -- `pg_advisory_xact_lock(int, int)` is the two-argument form; the first
    -- int is a namespace discriminator and the second is the resource id.
    PERFORM pg_advisory_xact_lock(hashtext('xmtpd.gateway_envelopes_meta_l1'), _oid);

    -- Short-circuit if the partition is already attached to the expected
    -- parent. This is the authoritative check — not a regex on SQLERRM.
    SELECT EXISTS (
        SELECT 1
        FROM pg_inherits i
                 JOIN pg_class c ON c.oid = i.inhrelid
                 JOIN pg_class p ON p.oid = i.inhparent
        WHERE c.relname = subname
          AND p.relname = 'gateway_envelopes_meta'
    ) INTO already_attached;

    IF already_attached THEN
        -- Defensive cleanup: older callers may have left seed constraints
        -- behind if ATTACH raised while the constraint was present. A
        -- successful attach means the constraint is no longer needed.
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT IF EXISTS oid_check;',
            subname
        );
        RETURN;
    END IF;

    -- Create the child with a validating CHECK so PostgreSQL can skip the
    -- full-scan validation during ATTACH. The CHECK is dropped immediately
    -- after the successful attach.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_meta INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT oid_check CHECK (originator_node_id = %s)
        ) PARTITION BY RANGE (originator_sequence_id);',
        subname,
        _oid::text
    );

    EXECUTE format(
        'ALTER TABLE gateway_envelopes_meta ATTACH PARTITION %I
            FOR VALUES IN (%s);',
        subname,
        _oid::text
    );

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT IF EXISTS oid_check;',
        subname
    );
END;
$$ LANGUAGE plpgsql;


-- META: create a RANGE subpartition [start, end), idempotently.
CREATE FUNCTION make_meta_seq_subpart_v3(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    parent  text := format('gateway_envelopes_meta_o%s', _oid);
    subname text := format('gateway_envelopes_meta_o%s_s%s_%s', _oid, _start, _end);
    already_attached boolean;
BEGIN
    -- Serialize concurrent callers per (originator, band_start).
    PERFORM pg_advisory_xact_lock(
        hashtext('xmtpd.gateway_envelopes_meta_l2'),
        hashtext(format('%s:%s', _oid, _start))
    );

    SELECT EXISTS (
        SELECT 1
        FROM pg_inherits i
                 JOIN pg_class c ON c.oid = i.inhrelid
                 JOIN pg_class p ON p.oid = i.inhparent
        WHERE c.relname = subname
          AND p.relname = parent
    ) INTO already_attached;

    IF already_attached THEN
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT IF EXISTS seq_id_check;',
            subname
        );
        RETURN;
    END IF;

    -- Correct CHECK predicate: bounds come from (_start, _end). The v2 helper
    -- had a format() arity bug that produced `>= _oid AND < _start` instead.
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_meta INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        );',
        subname,
        _start::text,
        _end::text
    );

    EXECUTE format(
        'ALTER TABLE %I ATTACH PARTITION %I
            FOR VALUES FROM (%s) TO (%s);',
        parent,
        subname,
        _start::text,
        _end::text
    );

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT IF EXISTS seq_id_check;',
        subname
    );
END;
$$ LANGUAGE plpgsql;


-- BLOB: create LIST child for one originator, idempotently.
CREATE FUNCTION make_blob_originator_part_v4(_oid int)
    RETURNS void AS $$
DECLARE
    subname text := format('gateway_envelopes_blob_o%s', _oid);
    already_attached boolean;
BEGIN
    PERFORM pg_advisory_xact_lock(hashtext('xmtpd.gateway_envelopes_blob_l1'), _oid);

    SELECT EXISTS (
        SELECT 1
        FROM pg_inherits i
                 JOIN pg_class c ON c.oid = i.inhrelid
                 JOIN pg_class p ON p.oid = i.inhparent
        WHERE c.relname = subname
          AND p.relname = 'gateway_envelopes_blob'
    ) INTO already_attached;

    IF already_attached THEN
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT IF EXISTS oid_check;',
            subname
        );
        RETURN;
    END IF;

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_blob INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT oid_check CHECK (originator_node_id = %s)
        ) PARTITION BY RANGE (originator_sequence_id);',
        subname,
        _oid::text
    );

    EXECUTE format(
        'ALTER TABLE gateway_envelopes_blob ATTACH PARTITION %I
            FOR VALUES IN (%s);',
        subname,
        _oid::text
    );

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT IF EXISTS oid_check;',
        subname
    );
END;
$$ LANGUAGE plpgsql;


-- BLOB: create a RANGE subpartition [start, end), idempotently.
CREATE FUNCTION make_blob_seq_subpart_v4(_oid int, _start bigint, _end bigint)
    RETURNS void AS $$
DECLARE
    parent  text := format('gateway_envelopes_blob_o%s', _oid);
    subname text := format('gateway_envelopes_blob_o%s_s%s_%s', _oid, _start, _end);
    already_attached boolean;
BEGIN
    PERFORM pg_advisory_xact_lock(
        hashtext('xmtpd.gateway_envelopes_blob_l2'),
        hashtext(format('%s:%s', _oid, _start))
    );

    SELECT EXISTS (
        SELECT 1
        FROM pg_inherits i
                 JOIN pg_class c ON c.oid = i.inhrelid
                 JOIN pg_class p ON p.oid = i.inhparent
        WHERE c.relname = subname
          AND p.relname = parent
    ) INTO already_attached;

    IF already_attached THEN
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT IF EXISTS seq_id_check;',
            subname
        );
        RETURN;
    END IF;

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I (
            LIKE gateway_envelopes_blob INCLUDING DEFAULTS INCLUDING CONSTRAINTS,
            CONSTRAINT seq_id_check CHECK ( originator_sequence_id >= %s AND originator_sequence_id < %s )
        );',
        subname,
        _start::text,
        _end::text
    );

    EXECUTE format(
        'ALTER TABLE %I ATTACH PARTITION %I
            FOR VALUES FROM (%s) TO (%s);',
        parent,
        subname,
        _start::text,
        _end::text
    );

    EXECUTE format(
        'ALTER TABLE %I DROP CONSTRAINT IF EXISTS seq_id_check;',
        subname
    );
END;
$$ LANGUAGE plpgsql;


-- Production partition ensure. Calls the hardened helpers above.
CREATE FUNCTION ensure_gateway_parts_v4(
    p_originator_node_id     int,
    p_originator_sequence_id bigint,
    p_band_width             bigint DEFAULT 1000000
) RETURNS void LANGUAGE plpgsql AS $$
DECLARE
    v_band_start bigint := (p_originator_sequence_id / p_band_width) * p_band_width;
BEGIN
    PERFORM make_meta_originator_part_v3(p_originator_node_id);
    PERFORM make_blob_originator_part_v4(p_originator_node_id);
    PERFORM make_meta_seq_subpart_v3(
        p_originator_node_id,
        v_band_start,
        v_band_start + p_band_width
    );
    PERFORM make_blob_seq_subpart_v4(
        p_originator_node_id,
        v_band_start,
        v_band_start + p_band_width
    );
END;
$$;

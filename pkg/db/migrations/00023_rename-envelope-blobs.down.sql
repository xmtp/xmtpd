-- Revert the rename of gateway_envelopes_blobs back to gateway_envelope_blobs.
-- The corresponding up migration is append-only: it did NOT modify the v1/v2
-- partition or batch-insert functions, so this down migration only needs to
-- (a) drop the new _v3 functions added by the up, and (b) rename the table
-- and its child partitions back to the old name. The pre-existing v1/v2
-- functions (which still reference gateway_envelope_blobs) will become
-- functional again as soon as the table is renamed back.

-- Step 1: Drop the v3 functions added by the up migration.
DROP FUNCTION IF EXISTS insert_gateway_envelope_batch_v3(
    int[], bigint[], bytea[], int[], timestamp[], bigint[], bytea[], bigint[], boolean[], boolean[]
);
DROP FUNCTION IF EXISTS ensure_gateway_parts_v3(int, bigint, bigint);
DROP FUNCTION IF EXISTS make_blob_seq_subpart_v3(int, bigint, bigint);
DROP FUNCTION IF EXISTS make_blob_originator_part_v3(int);

-- Step 2: Rename all child partitions back (deepest first)
DO $$
DECLARE
    r RECORD;
    old_name text;
    new_name text;
BEGIN
    FOR r IN
        SELECT c.relname
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'public'
          AND c.relkind IN ('r', 'p')
          AND c.relname LIKE 'gateway_envelopes_blobs_%'
        ORDER BY length(c.relname) DESC
    LOOP
        old_name := r.relname;
        new_name := 'gateway_envelope_blobs' || substring(old_name FROM length('gateway_envelopes_blobs') + 1);
        EXECUTE format('ALTER TABLE %I RENAME TO %I', old_name, new_name);
    END LOOP;
END$$;

-- Step 3: Rename the parent table back
ALTER TABLE gateway_envelopes_blobs RENAME TO gateway_envelope_blobs;

-- Step 4: Recreate the view to reference the old table name.
CREATE OR REPLACE VIEW gateway_envelopes_view AS
SELECT
    m.originator_node_id,
    m.originator_sequence_id,
    m.gateway_time,
    m.topic,
    b.originator_envelope
FROM gateway_envelopes_meta m
         JOIN gateway_envelope_blobs b
              ON b.originator_node_id     = m.originator_node_id
                  AND b.originator_sequence_id = m.originator_sequence_id;

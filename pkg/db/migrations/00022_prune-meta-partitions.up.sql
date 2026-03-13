CREATE OR REPLACE FUNCTION get_prunable_meta_partitions()
    RETURNS TABLE (
                      originator_node_id int,
                      schemaname text,
                      tablename text,
                      band_start bigint,
                      band_end bigint
                  )
    LANGUAGE plpgsql
AS $$
DECLARE
    part RECORD;
    has_rows BOOLEAN;
BEGIN
    FOR part IN
        WITH leaf_parts AS (
            SELECT
                cn.nspname AS schemaname,
                c.relname  AS tablename,
                ((regexp_match(
                        c.relname,
                        '^gateway_envelopes_meta_o([0-9]+)_s([0-9]+)_([0-9]+)$'
                  ))[1])::int    AS originator_node_id,
                ((regexp_match(
                        c.relname,
                        '^gateway_envelopes_meta_o([0-9]+)_s([0-9]+)_([0-9]+)$'
                  ))[2])::bigint AS band_start,
                ((regexp_match(
                        c.relname,
                        '^gateway_envelopes_meta_o([0-9]+)_s([0-9]+)_([0-9]+)$'
                  ))[3])::bigint AS band_end
            FROM pg_inherits i
                     JOIN pg_class c      ON c.oid = i.inhrelid
                     JOIN pg_namespace cn ON cn.oid = c.relnamespace
                     JOIN pg_class p      ON p.oid = i.inhparent
            WHERE c.relkind = 'r'
              AND p.relkind = 'p'
              AND c.relname ~ '^gateway_envelopes_meta_o[0-9]+_s[0-9]+_[0-9]+$'
              AND NOT EXISTS (
                SELECT 1
                FROM pg_inherits i2
                WHERE i2.inhparent = c.oid
            )
        ),
             ranked AS (
                 SELECT
                     lp.*,
                     row_number() OVER (
                         PARTITION BY lp.originator_node_id
                         ORDER BY lp.band_start DESC
                         ) AS rn
                 FROM leaf_parts lp
             )
        SELECT
            r.originator_node_id,
            r.schemaname,
            r.tablename,
            r.band_start,
            r.band_end
        FROM ranked r
        WHERE r.rn > 1
        ORDER BY r.originator_node_id, r.band_start
        LOOP
            EXECUTE format(
                    'SELECT EXISTS (SELECT 1 FROM %I.%I LIMIT 1)',
                    part.schemaname,
                    part.tablename
                    )
                INTO has_rows;

            IF NOT has_rows THEN
                originator_node_id := part.originator_node_id;
                schemaname := part.schemaname;
                tablename := part.tablename;
                band_start := part.band_start;
                band_end := part.band_end;
                RETURN NEXT;
            END IF;
        END LOOP;

    RETURN;
END;
$$;

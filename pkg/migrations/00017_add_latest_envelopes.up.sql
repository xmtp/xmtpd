-- Run without timeout, needed as the insert can take a few seconds
SET statement_timeout = 0;

CREATE OR REPLACE FUNCTION update_latest_envelope_v2()
RETURNS trigger AS $$
BEGIN
    INSERT INTO gateway_envelopes_latest as g
    SELECT originator_node_id, originator_sequence_id, gateway_time
    FROM (
        SELECT originator_node_id, originator_sequence_id, gateway_time, ROW_NUMBER() OVER (PARTITION BY originator_node_id ORDER BY originator_sequence_id DESC, gateway_time DESC) as rn
        FROM new
    ) ranked
    WHERE rn = 1
    ON CONFLICT (originator_node_id)
    DO UPDATE
        SET originator_sequence_id = EXCLUDED.originator_sequence_id,
            gateway_time = EXCLUDED.gateway_time
        WHERE EXCLUDED.originator_sequence_id > g.originator_sequence_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create a new trigger that's per statement.
CREATE TRIGGER gateway_latest_upd_v2
    AFTER INSERT ON gateway_envelopes_meta
    REFERENCING NEW TABLE AS new
    FOR EACH STATEMENT EXECUTE FUNCTION update_latest_envelope_v2();

-- Remove old trigger,
DROP TRIGGER IF EXISTS gateway_latest_upd ON gateway_envelopes_meta;
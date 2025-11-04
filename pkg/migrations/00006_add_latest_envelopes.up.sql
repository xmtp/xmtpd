-- Run without timeout, needed as the insert can take a few seconds
SET statement_timeout = 0;

CREATE TABLE gateway_envelopes_latest (
    originator_node_id int PRIMARY KEY,
    originator_sequence_id bigint NOT NULL,
    gateway_time timestamp NOT NULL
);

CREATE OR REPLACE FUNCTION update_latest_envelope() RETURNS trigger AS $$
BEGIN
    INSERT INTO gateway_envelopes_latest as g
    VALUES (NEW.originator_node_id, NEW.originator_sequence_id, NEW.gateway_time)
    ON CONFLICT (originator_node_id)
        DO UPDATE
        SET originator_sequence_id = EXCLUDED.originator_sequence_id,
            gateway_time = EXCLUDED.gateway_time
        WHERE EXCLUDED.originator_sequence_id > g.originator_sequence_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER gateway_latest_upd
    AFTER INSERT ON gateway_envelopes_meta
    FOR EACH ROW EXECUTE FUNCTION update_latest_envelope();
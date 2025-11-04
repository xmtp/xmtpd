CREATE TABLE IF NOT EXISTS gateway_envelopes_meta_v2
(
    gateway_time           timestamp NOT NULL DEFAULT now(),
    originator_node_id     int       NOT NULL,
    originator_sequence_id bigint    NOT NULL,
    topic                  bytea     NOT NULL,
    -- Leave column nullable since blockchain originated messages won't have a payer_id
    payer_id               int REFERENCES payers(id),
    expiry                 bigint NOT NULL,
    PRIMARY KEY (originator_node_id, originator_sequence_id)
) PARTITION BY LIST (originator_node_id);

-- BLOBS (cold path)
CREATE TABLE IF NOT EXISTS gateway_envelope_blobs_v2
(
    originator_node_id     int    NOT NULL,
    originator_sequence_id bigint NOT NULL,
    originator_envelope    bytea  NOT NULL,
    PRIMARY KEY (originator_node_id, originator_sequence_id),
    FOREIGN KEY (originator_node_id, originator_sequence_id) REFERENCES gateway_envelopes_meta_v2(originator_node_id, originator_sequence_id) ON DELETE CASCADE
) PARTITION BY LIST (originator_node_id);

CREATE TRIGGER gateway_v2_latest_upd
    AFTER INSERT ON gateway_envelopes_meta_v2
    FOR EACH ROW EXECUTE FUNCTION update_latest_envelope();

CREATE OR REPLACE VIEW gateway_envelopes_v2_view AS
SELECT
    m.originator_node_id,
    m.originator_sequence_id,
    m.gateway_time,
    m.topic,
    b.originator_envelope
FROM gateway_envelopes_meta_v2 m
         JOIN gateway_envelope_blobs_v2 b
              ON b.originator_node_id     = m.originator_node_id
                  AND b.originator_sequence_id = m.originator_sequence_id;
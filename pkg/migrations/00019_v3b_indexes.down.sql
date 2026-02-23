SET statement_timeout = 0;

-- Restore original indexes.
CREATE INDEX IF NOT EXISTS gem_time_node_seq_idx
    ON gateway_envelopes_meta (gateway_time, originator_node_id, originator_sequence_id);

CREATE INDEX IF NOT EXISTS gem_topic_time_idx
    ON gateway_envelopes_meta (topic, gateway_time, originator_node_id, originator_sequence_id);

CREATE INDEX IF NOT EXISTS gem_originator_node_id
    ON gateway_envelopes_meta(originator_node_id);

-- Drop V3b index.
DROP INDEX IF EXISTS gem_topic_orig_seq_idx;

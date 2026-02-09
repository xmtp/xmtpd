-- drop new index
DROP INDEX IF EXISTS gem_topic_idx;

-- revert to old index
CREATE INDEX IF NOT EXISTS gem_topic_time_idx
    ON gateway_envelopes_meta (topic, gateway_time, originator_node_id, originator_sequence_id);
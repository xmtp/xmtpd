SET statement_timeout = 0;

DROP INDEX IF EXISTS gem_topic_time_idx;

-- required for SelectGatewayEnvelopesByTopics
CREATE INDEX IF NOT EXISTS gem_topic_idx
    ON gateway_envelopes_meta (topic, originator_node_id, originator_sequence_id);
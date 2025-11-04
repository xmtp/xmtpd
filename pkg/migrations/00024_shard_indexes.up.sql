SET statement_timeout = 0;

-- required for most gateway sorted selects
CREATE INDEX IF NOT EXISTS gem_time_node_seq_idx
    ON gateway_envelopes_meta (gateway_time, originator_node_id, originator_sequence_id);

-- required for SelectGatewayEnvelopesByTopics
CREATE INDEX IF NOT EXISTS gem_topic_time_idx
    ON gateway_envelopes_meta (topic, gateway_time, originator_node_id, originator_sequence_id);

-- required for SelectNewestFromTopics
CREATE INDEX IF NOT EXISTS gem_topic_time_desc_idx
    ON gateway_envelopes_meta (topic, gateway_time DESC)
    INCLUDE (originator_node_id, originator_sequence_id);

-- required for pruning
CREATE INDEX gem_expiry_idx
    ON gateway_envelopes_meta (expiry)
    INCLUDE (originator_node_id, originator_sequence_id)
    WHERE expiry IS NOT NULL;

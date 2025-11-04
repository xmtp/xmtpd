SET statement_timeout = 0;

-- required for most gateway sorted selects
CREATE INDEX IF NOT EXISTS gem_v2_time_node_seq_idx
    ON gateway_envelopes_meta_v2 (gateway_time, originator_node_id, originator_sequence_id);

-- required for SelectGatewayEnvelopesV2ByTopics
CREATE INDEX IF NOT EXISTS gem_v2_topic_time_idx
    ON gateway_envelopes_meta_v2 (topic, gateway_time, originator_node_id, originator_sequence_id);

-- required for SelectNewestFromTopicsV2
CREATE INDEX IF NOT EXISTS gem_v2_topic_time_desc_idx
    ON gateway_envelopes_meta_v2 (topic, gateway_time DESC)
    INCLUDE (originator_node_id, originator_sequence_id);

-- required for pruning
CREATE INDEX gem_expiry_idx
    ON gateway_envelopes_meta_v2 (expiry)
    INCLUDE (originator_node_id, originator_sequence_id)
    WHERE expiry IS NOT NULL;

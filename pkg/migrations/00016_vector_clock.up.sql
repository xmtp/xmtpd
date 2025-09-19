-- Cover index for vector clock queries
CREATE INDEX idx_gateway_envelopes_vector_clock ON gateway_envelopes(originator_node_id, originator_sequence_id DESC, gateway_time);

-- Index efficient for distinct originator_node_id
CREATE INDEX idx_gateway_envelopes_originator_node_id 
ON gateway_envelopes(originator_node_id);

SET statement_timeout = 0;

-- Covering index for V3b LATERAL per-(topic, originator) query.
-- Enables index-only scans with zero heap fetches.
CREATE INDEX IF NOT EXISTS gem_topic_orig_seq_idx
    ON gateway_envelopes_meta (topic, originator_node_id, originator_sequence_id)
    INCLUDE (gateway_time);

-- Drop indexes superseded by gem_topic_orig_seq_idx or redundant with partitioning.
DROP INDEX IF EXISTS gem_topic_time_idx;
DROP INDEX IF EXISTS gem_time_node_seq_idx;
DROP INDEX IF EXISTS gem_originator_node_id;

-- Run without timeout, needed as the table can contain up to millions of rows
SET statement_timeout = 0;

CREATE INDEX IF NOT EXISTS idx_gw_by_originator_cover
    ON gateway_envelopes (originator_node_id, originator_sequence_id, gateway_time)
    INCLUDE (topic);
-- Run without timeout, needed as the table can contain up to millions of rows
SET statement_timeout = 0;

CREATE INDEX gateway_envelopes_expiry_idx
    ON gateway_envelopes (expiry)
    INCLUDE (originator_node_id, originator_sequence_id)
    WHERE expiry IS NOT NULL;
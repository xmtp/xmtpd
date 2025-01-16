-- Add blockchain-related columns and constraint
ALTER TABLE gateway_envelopes
	ADD COLUMN block_number BIGINT,
	ADD COLUMN block_hash BYTEA,
	ADD COLUMN version INT,
	ADD COLUMN is_canonical BOOLEAN;

ALTER TABLE gateway_envelopes
	ADD CONSTRAINT blockchain_message_constraint CHECK ((block_number IS NULL AND block_hash IS NULL AND version IS NULL AND is_canonical IS NULL) OR (block_number IS NOT NULL AND block_hash IS NOT NULL AND version IS NOT NULL AND is_canonical IS NOT NULL));

CREATE INDEX idx_gateway_envelopes_reorg ON gateway_envelopes(block_number DESC, block_hash)
WHERE
	block_number IS NOT NULL AND is_canonical = TRUE;


-- Drop everything in reverse order
DROP INDEX IF EXISTS idx_gateway_envelopes_reorg;

ALTER TABLE gateway_envelopes
	DROP CONSTRAINT IF EXISTS blockchain_message_constraint;

ALTER TABLE gateway_envelopes
	DROP COLUMN IF EXISTS block_number,
	DROP COLUMN IF EXISTS block_hash,
	DROP COLUMN IF EXISTS version,
	DROP COLUMN IF EXISTS is_canonical;


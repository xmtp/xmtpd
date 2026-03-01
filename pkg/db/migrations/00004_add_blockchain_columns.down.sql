-- Drop index first
DROP INDEX IF EXISTS idx_blockchain_messages_canonical;

-- Then drop the table
DROP TABLE IF EXISTS blockchain_messages;

-- Drop newly added column
ALTER TABLE latest_block DROP COLUMN IF EXISTS block_hash;
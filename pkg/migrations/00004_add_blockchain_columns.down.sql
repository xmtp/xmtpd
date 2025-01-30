-- Drop index first
DROP INDEX IF EXISTS idx_blockchain_messages_canonical;

-- Then drop the table
DROP TABLE IF EXISTS blockchain_messages;


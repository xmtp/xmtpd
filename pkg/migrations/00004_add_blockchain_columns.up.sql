CREATE TABLE blockchain_messages(
	block_number BIGINT NOT NULL,
	block_hash BYTEA NOT NULL,
	originator_node_id INT NOT NULL,
	originator_sequence_id BIGINT NOT NULL,
	is_canonical BOOLEAN NOT NULL DEFAULT TRUE,
	PRIMARY KEY (block_number, block_hash, originator_node_id, originator_sequence_id),
	FOREIGN KEY (originator_node_id, originator_sequence_id) REFERENCES gateway_envelopes(originator_node_id, originator_sequence_id)
);

CREATE INDEX idx_blockchain_messages_block_canonical ON blockchain_messages(block_number, is_canonical);

ALTER TABLE latest_block ADD COLUMN block_hash BYTEA;
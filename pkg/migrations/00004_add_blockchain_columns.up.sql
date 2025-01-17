-- Create table to track blockchain messages
CREATE TABLE blockchain_messages(
	block_number BIGINT NOT NULL,
	block_hash BYTEA NOT NULL,
	originator_node_id INT NOT NULL,
	originator_sequence_id BIGINT NOT NULL,
	is_canonical BOOLEAN NOT NULL DEFAULT TRUE,
	PRIMARY KEY (block_number, block_hash),
	FOREIGN KEY (originator_node_id, originator_sequence_id) REFERENCES gateway_envelopes(originator_node_id, originator_sequence_id)
);

CREATE INDEX idx_blockchain_messages_canonical ON blockchain_messages(block_number DESC, block_hash)
WHERE
	is_canonical = TRUE;


-- name: SetLatestBlock :exec
INSERT INTO latest_block(contract_address, block_number, block_hash)
	VALUES (@contract_address, @block_number, @block_hash)
ON CONFLICT (contract_address)
	DO UPDATE SET
		block_number = @block_number, block_hash = @block_hash
	WHERE
		@block_number > latest_block.block_number
		AND @block_hash != latest_block.block_hash;

-- name: GetLatestBlock :one
SELECT
	block_number,
	block_hash
FROM
	latest_block
WHERE
	contract_address = @contract_address;

-- name: InsertBlockchainMessage :exec
INSERT INTO blockchain_messages(block_number, block_hash, originator_node_id, originator_sequence_id, is_canonical)
	VALUES (@block_number, @block_hash, @originator_node_id, @originator_sequence_id, @is_canonical)
ON CONFLICT
	DO NOTHING;

-- name: GetBlocksInRange :many
-- Returns blocks in ascending order (oldest to newest)
-- StartBlock should be the lower bound (older block)
-- EndBlock should be the upper bound (newer block)
-- Example: GetBlocksInRange(1000, 2000), returns 1000, 1001, 1002, ..., 2000
SELECT DISTINCT ON (block_number)
	block_number,
	block_hash
FROM
	blockchain_messages
WHERE
	block_number BETWEEN @start_block AND @end_block
	AND block_hash IS NOT NULL
	AND is_canonical = TRUE
ORDER BY
	block_number ASC,
	block_hash;

-- name: UpdateBlocksCanonicalityInRange :exec
UPDATE
	blockchain_messages AS bm
SET
	is_canonical = FALSE
FROM (
	SELECT
		block_number
	FROM
		blockchain_messages
	WHERE
		bm.block_number BETWEEN @start_block_number AND @end_block_number
	FOR UPDATE) AS locked_rows
WHERE
	bm.block_number = locked_rows.block_number;
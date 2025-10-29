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
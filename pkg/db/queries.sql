-- name: InsertNodeInfo :execrows
INSERT INTO node_info(node_id, public_key)
	VALUES (@node_id, @public_key)
ON CONFLICT
	DO NOTHING;

-- name: SelectNodeInfo :one
SELECT
	*
FROM
	node_info
WHERE
	singleton_id = 1;

-- name: InsertGatewayEnvelope :execrows
INSERT INTO gateway_envelopes(originator_node_id, originator_sequence_id, topic, originator_envelope, payer_id, gateway_time)
	VALUES (@originator_node_id, @originator_sequence_id, @topic, @originator_envelope, @payer_id, COALESCE(@gateway_time, NOW()))
ON CONFLICT
	DO NOTHING;

-- name: SelectGatewayEnvelopes :many
SELECT
	*
FROM
	select_gateway_envelopes(@cursor_node_ids::INT[], @cursor_sequence_ids::BIGINT[], @topics::BYTEA[], @originator_node_ids::INT[], @row_limit::INT);

-- name: InsertStagedOriginatorEnvelope :one
SELECT
	*
FROM
	insert_staged_originator_envelope(@topic, @payer_envelope);

-- name: SelectStagedOriginatorEnvelopes :many
SELECT
	*
FROM
	staged_originator_envelopes
WHERE
	id > @last_seen_id
ORDER BY
	id ASC
LIMIT @num_rows;

-- name: DeleteStagedOriginatorEnvelope :execrows
DELETE FROM staged_originator_envelopes
WHERE id = @id;

-- name: SelectVectorClock :many
SELECT DISTINCT ON (originator_node_id)
	originator_node_id,
	originator_sequence_id,
	originator_envelope
FROM
	gateway_envelopes
ORDER BY
	originator_node_id,
	originator_sequence_id DESC;

-- name: GetAddressLogs :many
SELECT
	a.address,
	encode(a.inbox_id, 'hex') AS inbox_id,
	a.association_sequence_id
FROM
	address_log a
	INNER JOIN (
		SELECT
			address,
			MAX(association_sequence_id) AS max_association_sequence_id
		FROM
			address_log
		WHERE
			address = ANY (@addresses::TEXT[])
			AND revocation_sequence_id IS NULL
		GROUP BY
			address) b ON a.address = b.address
	AND a.association_sequence_id = b.max_association_sequence_id;

-- name: InsertAddressLog :execrows
INSERT INTO address_log(address, inbox_id, association_sequence_id, revocation_sequence_id)
	VALUES (@address, decode(@inbox_id, 'hex'), @association_sequence_id, NULL)
ON CONFLICT (address, inbox_id)
	DO UPDATE SET
		revocation_sequence_id = NULL, association_sequence_id = @association_sequence_id
	WHERE (address_log.revocation_sequence_id IS NULL
		OR address_log.revocation_sequence_id < @association_sequence_id)
		AND address_log.association_sequence_id < @association_sequence_id;

-- name: RevokeAddressFromLog :execrows
UPDATE
	address_log
SET
	revocation_sequence_id = @revocation_sequence_id
WHERE
	address = @address
	AND inbox_id = decode(@inbox_id, 'hex');

-- name: GetLatestSequenceId :one
SELECT
	COALESCE(max(originator_sequence_id), 0)::BIGINT AS originator_sequence_id
FROM
	gateway_envelopes
WHERE
	originator_node_id = @originator_node_id;

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

-- name: GetLatestCursor :many
SELECT
	originator_node_id,
	MAX(originator_sequence_id)::BIGINT AS max_sequence_id
FROM
	gateway_envelopes
GROUP BY
	originator_node_id;

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

-- name: FindOrCreatePayer :one
INSERT INTO payers(address)
	VALUES (@address)
ON CONFLICT (address)
	DO UPDATE SET
		address = @address
	RETURNING
		id;

-- name: IncrementUnsettledUsage :exec
INSERT INTO unsettled_usage(payer_id, originator_id, minutes_since_epoch, spend_picodollars)
	VALUES (@payer_id, @originator_id, @minutes_since_epoch, @spend_picodollars)
ON CONFLICT (payer_id, originator_id, minutes_since_epoch)
	DO UPDATE SET
		spend_picodollars = unsettled_usage.spend_picodollars + @spend_picodollars;

-- name: GetPayerUnsettledUsage :one
SELECT
	COALESCE(SUM(spend_picodollars), 0)::BIGINT AS total_spend_picodollars
FROM
	unsettled_usage
WHERE
	payer_id = @payer_id
	AND (@minutes_since_epoch_gt::BIGINT = 0
		OR minutes_since_epoch > @minutes_since_epoch_gt::BIGINT)
	AND (@minutes_since_epoch_lt::BIGINT = 0
		OR minutes_since_epoch < @minutes_since_epoch_lt::BIGINT);

-- name: FillNonceSequence :exec
SELECT fill_nonce_gap(@pending_nonce, @num_elements);

-- name: GetNextAvailableNonce :one
SELECT
    nonce
FROM
    nonce_table
ORDER BY
    nonce
    ASC LIMIT 1
    FOR UPDATE SKIP LOCKED;

-- name: DeleteAvailableNonce :execrows
DELETE FROM nonce_table
WHERE nonce = @nonce;

-- name: DeleteObsoleteNonces :execrows
WITH deletable AS (
    SELECT n.nonce
    FROM nonce_table n
    WHERE n.nonce < @nonce
    FOR UPDATE SKIP LOCKED
)
DELETE FROM nonce_table
    USING deletable
WHERE nonce_table.nonce = deletable.nonce;
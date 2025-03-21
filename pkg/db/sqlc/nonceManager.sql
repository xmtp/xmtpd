-- name: FillNonceSequence :one
SELECT COALESCE(
		fill_nonce_gap(@pending_nonce, @num_elements),
		@num_elements
	)::INT AS inserted_rows;

-- name: GetNextAvailableNonce :one
SELECT nonce
FROM nonce_table
ORDER BY nonce ASC
LIMIT 1 FOR
UPDATE SKIP LOCKED;

-- name: DeleteAvailableNonce :execrows
DELETE FROM nonce_table
WHERE nonce = @nonce;

-- name: DeleteObsoleteNonces :execrows
WITH deletable AS (
	SELECT
		n.nonce
	FROM
		nonce_table n
	WHERE
		n.nonce < @nonce
	FOR UPDATE
		SKIP LOCKED)
DELETE FROM nonce_table USING deletable
WHERE nonce_table.nonce = deletable.nonce;
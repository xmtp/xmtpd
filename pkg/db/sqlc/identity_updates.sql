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

-- name: InsertAddressLogsBatch :execrows
WITH input AS (
    SELECT unnest(@addresses::TEXT[]) AS address, 
           decode(@inbox_id, 'hex') AS inbox_id, 
           @association_sequence_id::BIGINT AS association_sequence_id
)
INSERT INTO address_log(address, inbox_id, association_sequence_id, revocation_sequence_id)
SELECT address, inbox_id, association_sequence_id, NULL
FROM input
ON CONFLICT (address, inbox_id)
    DO UPDATE SET
        revocation_sequence_id = NULL, 
        association_sequence_id = EXCLUDED.association_sequence_id
    WHERE (address_log.revocation_sequence_id IS NULL
        OR address_log.revocation_sequence_id < EXCLUDED.association_sequence_id)
        AND address_log.association_sequence_id < EXCLUDED.association_sequence_id;

-- name: RevokeAddressFromLog :execrows
UPDATE
	address_log
SET
	revocation_sequence_id = @revocation_sequence_id
WHERE
	address = @address
	AND inbox_id = decode(@inbox_id, 'hex');

-- name: RevokeAddressFromLogBatch :execrows
WITH input AS (
    SELECT unnest(@addresses::TEXT[]) AS address,
           decode(@inbox_id, 'hex') AS inbox_id,
           @revocation_sequence_id::BIGINT AS revocation_sequence_id
)
UPDATE address_log AS al
SET revocation_sequence_id = input.revocation_sequence_id
FROM input
WHERE al.address = input.address
  AND al.inbox_id = input.inbox_id;

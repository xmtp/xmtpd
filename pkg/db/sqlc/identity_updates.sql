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

-- name: AdvisoryLockIdentityUpdateInsert :exec
SELECT pg_advisory_xact_lock(
-- only take the lowest 32 bits (mod 2^32-1)
    (@node_id::bigint << 32) | (@sequence_id & 4294967295)
);
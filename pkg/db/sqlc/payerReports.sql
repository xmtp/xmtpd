-- name: FindOrCreatePayer :one
INSERT INTO payers(address)
VALUES (@address) ON CONFLICT (address) DO
UPDATE
SET address = @address
RETURNING id;

-- name: IncrementUnsettledUsage :exec
INSERT INTO unsettled_usage(
		payer_id,
		originator_id,
		minutes_since_epoch,
		spend_picodollars,
		last_sequence_id
	)
VALUES (
		@payer_id,
		@originator_id,
		@minutes_since_epoch,
		@spend_picodollars,
		@sequence_id
	) ON CONFLICT (payer_id, originator_id, minutes_since_epoch) DO
UPDATE
SET spend_picodollars = unsettled_usage.spend_picodollars + @spend_picodollars,
	last_sequence_id = GREATEST(unsettled_usage.last_sequence_id, @sequence_id);

-- name: GetPayerUnsettledUsage :one
SELECT COALESCE(SUM(spend_picodollars), 0)::BIGINT AS total_spend_picodollars,
	COALESCE(MAX(last_sequence_id), 0)::BIGINT AS last_sequence_id
FROM unsettled_usage
WHERE payer_id = @payer_id
	AND (
		@minutes_since_epoch_gt::BIGINT = 0
		OR minutes_since_epoch > @minutes_since_epoch_gt::BIGINT
	)
	AND (
		@minutes_since_epoch_lt::BIGINT = 0
		OR minutes_since_epoch < @minutes_since_epoch_lt::BIGINT
	);

-- name: BuildPayerReport :many
SELECT payers.address AS payer_address,
	SUM(spend_picodollars)::BIGINT AS total_spend_picodollars
FROM unsettled_usage
	JOIN payers ON payers.id = unsettled_usage.payer_id
WHERE originator_id = @originator_id
	AND minutes_since_epoch > @start_minutes_since_epoch
	AND minutes_since_epoch <= @end_minutes_since_epoch
GROUP BY payers.address;

-- name: GetGatewayEnvelopeByID :one
SELECT *
FROM gateway_envelopes
WHERE originator_sequence_id = @originator_sequence_id -- Include the node ID to take advantage of the primary key index
	AND originator_node_id = @originator_node_id;

-- name: GetSecondNewestMinute :one
WITH second_newest_minute AS (
	SELECT minutes_since_epoch
	FROM unsettled_usage
	WHERE originator_id = @originator_id
		AND unsettled_usage.minutes_since_epoch > @minimum_minutes_since_epoch
	GROUP BY unsettled_usage.minutes_since_epoch
	ORDER BY unsettled_usage.minutes_since_epoch DESC
	LIMIT 1 OFFSET 1
)
SELECT coalesce(max(last_sequence_id), 0)::BIGINT AS max_sequence_id,
	coalesce(max(unsettled_usage.minutes_since_epoch), 0)::INT AS minutes_since_epoch
FROM unsettled_usage
	JOIN second_newest_minute ON second_newest_minute.minutes_since_epoch = unsettled_usage.minutes_since_epoch
WHERE unsettled_usage.originator_id = @originator_id;

-- name: InsertOrIgnorePayerReport :exec
INSERT INTO payer_reports (
		id,
		originator_node_id,
		start_sequence_id,
		end_sequence_id,
		payers_merkle_root,
		payers_leaf_count,
		nodes_hash,
		nodes_count
	)
VALUES (
		@id,
		@originator_node_id,
		@start_sequence_id,
		@end_sequence_id,
		@payers_merkle_root,
		@payers_leaf_count,
		@nodes_hash,
		@nodes_count
	) ON CONFLICT (id) DO NOTHING;

-- name: InsertOrIgnorePayerReportAttestation :exec
INSERT INTO payer_report_attestations (payer_report_id, node_id, signature)
VALUES (@payer_report_id, @node_id, @signature) ON CONFLICT (payer_report_id, node_id) DO NOTHING;

-- name: FetchPayerReport :one
SELECT *
FROM payer_reports
WHERE id = @id;

-- name: FetchPayerReports :many
SELECT *
FROM payer_reports
WHERE (
		sqlc.narg(attestation_status)::SMALLINT IS NULL
		OR attestation_status = sqlc.narg(attestation_status)::SMALLINT
	)
	AND (
		sqlc.narg(submission_status)::SMALLINT IS NULL
		OR submission_status = sqlc.narg(submission_status)::SMALLINT
	)
	AND (
		sqlc.narg(created_after)::TIMESTAMP IS NULL
		OR created_at > sqlc.narg(created_after)
	)
	AND (
		sqlc.narg(end_sequence_id)::BIGINT IS NULL
		OR sqlc.narg(end_sequence_id) = end_sequence_id
	)
	AND (
		sqlc.narg(start_sequence_id)::BIGINT IS NULL
		OR sqlc.narg(start_sequence_id) = start_sequence_id
	);

-- name: SetReportAttestationStatus :exec
UPDATE payer_reports
SET attestation_status = @new_status
WHERE id = @report_id
	AND attestation_status IN (sqlc.slice(prev_status));
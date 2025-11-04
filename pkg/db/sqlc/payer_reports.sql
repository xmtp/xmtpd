-- name: FindOrCreatePayer :one
INSERT INTO payers(address)
VALUES (@address) ON CONFLICT (address) DO
UPDATE
SET address = @address
RETURNING id;

-- name: GetPayerByAddress :one
SELECT id
FROM payers
WHERE address = @address;

-- name: IncrementUnsettledUsage :exec
INSERT INTO unsettled_usage(
		payer_id,
		originator_id,
		minutes_since_epoch,
		spend_picodollars,
		last_sequence_id,
		message_count
	)
VALUES (
		@payer_id,
		@originator_id,
		@minutes_since_epoch,
		@spend_picodollars,
		@sequence_id,
		@message_count
	) ON CONFLICT (payer_id, originator_id, minutes_since_epoch) DO
UPDATE
SET spend_picodollars = unsettled_usage.spend_picodollars + @spend_picodollars,
	message_count = unsettled_usage.message_count + @message_count,
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

-- name: GetPayerInfoReport :many
SELECT EXTRACT(
		EPOCH
		FROM DATE_TRUNC(
				CASE
					WHEN @group_by = 'hour' THEN 'hour'
					ELSE 'day'
				END,
				TO_TIMESTAMP(minutes_since_epoch * 60)
			)
	)::BIGINT AS time_period,
	COALESCE(SUM(spend_picodollars), 0)::BIGINT AS total_spend_picodollars,
	COALESCE(SUM(message_count), 0)::INTEGER AS total_message_count
FROM unsettled_usage
WHERE payer_id = @payer_id
GROUP BY time_period
ORDER BY time_period;

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
FROM gateway_envelopes_view
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

-- name: GetLastSequenceIDForOriginatorMinute :one
SELECT COALESCE(MAX(last_sequence_id), 0)::BIGINT AS last_sequence_id
FROM unsettled_usage
WHERE originator_id = @originator_id
	AND minutes_since_epoch = @minutes_since_epoch;

-- name: InsertOrIgnorePayerReport :execrows
INSERT INTO payer_reports (
		id,
		originator_node_id,
		start_sequence_id,
		end_sequence_id,
		end_minute_since_epoch,
		payers_merkle_root,
		active_node_ids
	)
VALUES (
		@id,
		@originator_node_id,
		@start_sequence_id,
		@end_sequence_id,
		@end_minute_since_epoch,
		@payers_merkle_root,
		@active_node_ids
	) ON CONFLICT (id) DO NOTHING;

-- name: InsertOrIgnorePayerReportAttestation :exec
INSERT INTO payer_report_attestations (payer_report_id, node_id, signature)
VALUES (@payer_report_id, @node_id, @signature) ON CONFLICT (payer_report_id, node_id) DO NOTHING;

-- name: FetchPayerReports :many
WITH rpt AS (
	SELECT pr.*,
		pra.node_id,
		pra.signature,
		COUNT(pra.node_id) OVER (PARTITION BY pr.id) AS attestations_count
	FROM payer_reports AS pr
		LEFT JOIN payer_report_attestations AS pra ON pra.payer_report_id = pr.id
	WHERE (
			sqlc.narg(attestation_status_in)::SMALLINT [] IS NULL
			OR pr.attestation_status = ANY(sqlc.narg(attestation_status_in)::SMALLINT [])
		)
		AND (
			sqlc.narg(submission_status_in)::SMALLINT [] IS NULL
			OR pr.submission_status = ANY(sqlc.narg(submission_status_in)::SMALLINT [])
		)
		AND (
			sqlc.narg(created_after)::TIMESTAMP IS NULL
			OR pr.created_at > sqlc.narg(created_after)
		)
		AND (
			sqlc.narg(end_sequence_id)::BIGINT IS NULL
			OR sqlc.narg(end_sequence_id) = pr.end_sequence_id
		)
		AND (
			sqlc.narg(start_sequence_id)::BIGINT IS NULL
			OR sqlc.narg(start_sequence_id) = pr.start_sequence_id
		)
		AND (
			sqlc.narg(originator_node_id)::INT IS NULL
			OR sqlc.narg(originator_node_id) = pr.originator_node_id
		)
		AND (
			sqlc.narg(payer_report_id)::BYTEA IS NULL
			OR sqlc.narg(payer_report_id)::BYTEA = pr.id
		)
)
SELECT *
FROM rpt
WHERE sqlc.narg(min_attestations)::INT IS NULL
	OR attestations_count >= sqlc.narg(min_attestations)::INT
ORDER BY created_at ASC;

-- name: FetchPayerReport :one
SELECT *
FROM payer_reports
WHERE id = @id;

-- name: FetchPayerReportLocked :one
SELECT *
FROM payer_reports
WHERE id = @id
FOR UPDATE;


-- name: SetReportAttestationStatus :exec
UPDATE payer_reports
SET attestation_status = @new_status
WHERE id = @report_id
	AND attestation_status = ANY(sqlc.arg(prev_status)::SMALLINT []);

-- name: SetReportSubmissionStatus :exec
UPDATE payer_reports
SET submission_status = @new_status
WHERE id = @report_id
	AND submission_status = ANY(sqlc.arg(prev_status)::SMALLINT []);

-- name: SetReportSubmitted :exec
UPDATE payer_reports
SET submission_status = @new_status,
	submitted_report_index = sqlc.arg(submitted_report_index)::INTEGER
WHERE id = @report_id
	AND submission_status = ANY(sqlc.arg(prev_status)::SMALLINT []);

-- name: FetchAttestations :many
SELECT *
FROM payer_report_attestations
	LEFT JOIN payer_reports ON payer_reports.id = payer_report_attestations.payer_report_id
WHERE (
		sqlc.narg(payer_report_id)::BYTEA IS NULL
		OR sqlc.narg(payer_report_id)::BYTEA = payer_report_id
	)
	AND (
		sqlc.narg(attester_node_id)::INT IS NULL
		OR sqlc.narg(attester_node_id)::INT = node_id
	);

-- name: ClearUnsettledUsage :exec
DELETE FROM unsettled_usage
WHERE originator_id = @originator_id
	AND minutes_since_epoch > @prev_report_end_minute_since_epoch
	AND minutes_since_epoch <= @end_minute_since_epoch;

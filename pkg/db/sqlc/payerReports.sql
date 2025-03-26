-- name: FindOrCreatePayer :one
INSERT INTO payers(address)
	VALUES (@address)
ON CONFLICT (address)
	DO UPDATE SET
		address = @address
	RETURNING
		id;

-- name: IncrementUnsettledUsage :exec
INSERT INTO unsettled_usage(payer_id, originator_id, minutes_since_epoch, spend_picodollars, last_sequence_id)
	VALUES (@payer_id, @originator_id, @minutes_since_epoch, @spend_picodollars, @sequence_id)
ON CONFLICT (payer_id, originator_id, minutes_since_epoch)
	DO UPDATE SET
		spend_picodollars = unsettled_usage.spend_picodollars + @spend_picodollars,
		last_sequence_id = GREATEST(unsettled_usage.last_sequence_id, @sequence_id);

-- name: GetPayerUnsettledUsage :one
SELECT
	COALESCE(SUM(spend_picodollars), 0)::BIGINT AS total_spend_picodollars,
	COALESCE(MAX(last_sequence_id), 0)::BIGINT AS last_sequence_id
FROM
	unsettled_usage
WHERE
	payer_id = @payer_id
	AND (@minutes_since_epoch_gt::BIGINT = 0
		OR minutes_since_epoch > @minutes_since_epoch_gt::BIGINT)
	AND (@minutes_since_epoch_lt::BIGINT = 0
		OR minutes_since_epoch < @minutes_since_epoch_lt::BIGINT);

-- name: BuildPayerReport :many
SELECT
	payers.address as payer_address,
	SUM(spend_picodollars)::BIGINT AS total_spend_picodollars
FROM
	unsettled_usage
JOIN payers on payers.id = unsettled_usage.payer_id
WHERE
	originator_id = @originator_id
	AND minutes_since_epoch > @start_minutes_since_epoch
	AND minutes_since_epoch <= @end_minutes_since_epoch
GROUP BY
	payers.address;

-- name: GetGatewayEnvelopeByID :one
SELECT * FROM gateway_envelopes
WHERE originator_sequence_id = @originator_sequence_id
-- Include the node ID to take advantage of the primary key index
AND originator_node_id = @originator_node_id;

-- name: GetSecondNewestMinute :one
WITH second_newest_minute
AS
  (
           SELECT minutes_since_epoch
           FROM     unsettled_usage
           WHERE    originator_id = @originator_id
           AND      unsettled_usage.minutes_since_epoch > @minimum_minutes_since_epoch
           GROUP BY unsettled_usage.minutes_since_epoch
           ORDER BY unsettled_usage.minutes_since_epoch DESC
           LIMIT    1
           OFFSET   1)
  SELECT coalesce(max(last_sequence_id), 0)::BIGINT as max_sequence_id,
         coalesce(max(unsettled_usage.minutes_since_epoch), 0)::INT as minutes_since_epoch
  FROM   unsettled_usage
  JOIN   second_newest_minute
  ON     second_newest_minute.minutes_since_epoch = unsettled_usage.minutes_since_epoch
  WHERE  unsettled_usage.originator_id = @originator_id;
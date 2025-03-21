-- name: IncrementOriginatorCongestion :exec
INSERT INTO originator_congestion(originator_id, minutes_since_epoch, num_messages)
	VALUES (@originator_id, @minutes_since_epoch, 1)
ON CONFLICT (originator_id, minutes_since_epoch)
	DO UPDATE SET
		num_messages = originator_congestion.num_messages + 1;


-- name: GetRecentOriginatorCongestion :many
SELECT 
	minutes_since_epoch,
	num_messages
FROM
	originator_congestion
WHERE
	originator_id = $1
	AND minutes_since_epoch <= sqlc.arg(end_minute)::INT
	AND minutes_since_epoch > sqlc.arg(end_minute)::INT - sqlc.arg(num_minutes)::INT
ORDER BY minutes_since_epoch DESC;

-- name: SumOriginatorCongestion :one
SELECT
	COALESCE(SUM(num_messages), 0)::BIGINT AS num_messages
FROM
	originator_congestion
WHERE
	originator_id = @originator_id
	AND (@minutes_since_epoch_gt::BIGINT = 0
		OR minutes_since_epoch > @minutes_since_epoch_gt::BIGINT)
	AND (@minutes_since_epoch_lt::BIGINT = 0
		OR minutes_since_epoch < @minutes_since_epoch_lt::BIGINT);
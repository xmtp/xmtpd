-- name: AdvisoryLockWithKey :exec
SELECT pg_advisory_xact_lock(@locking_key);



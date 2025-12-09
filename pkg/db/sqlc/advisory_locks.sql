-- name: AdvisoryLockWithKey :exec
SELECT pg_advisory_xact_lock(@locking_key);

-- name: TryAdvisoryLockWithKey :one
SELECT pg_try_advisory_xact_lock(@locking_key) as lock_succeeded;



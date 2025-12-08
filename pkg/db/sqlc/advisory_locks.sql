-- name: AdvisoryTryLockWithKey :one
SELECT pg_try_advisory_xact_lock(@locking_key) as lock_succeeded;



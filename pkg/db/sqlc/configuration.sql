-- name: SetLocalWorkMem :one
SELECT set_config('work_mem', $1::TEXT, true);
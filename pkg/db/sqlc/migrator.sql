-- Migration Tracker Operations

-- name: GetMigrationProgress :one
SELECT last_migrated_id 
FROM migration_tracker 
WHERE source_table = @source_table;

-- name: UpdateMigrationProgress :exec
UPDATE migration_tracker 
SET last_migrated_id = @last_migrated_id,
    updated_at = NOW()
WHERE source_table = @source_table;

-- name: InsertMigrationDeadLetterBox :one
SELECT *
FROM insert_migration_dead_letter_box(@source_table, @sequence_id, @payload, @reason, @retryable);

-- name: DeleteMigrationDeadLetterBox :one
SELECT delete_migration_dead_letter_box(@source_table, @sequence_id);

-- name: GetRetryableMigrationDeadLetterBoxes :many
SELECT *
FROM migration_dead_letter_box
WHERE retryable = TRUE
ORDER BY retried_at ASC
LIMIT @row_limit;

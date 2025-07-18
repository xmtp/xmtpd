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

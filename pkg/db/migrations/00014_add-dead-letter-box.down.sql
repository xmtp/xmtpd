DROP TABLE IF EXISTS migration_dead_letter_box;
DROP INDEX IF EXISTS migration_dead_letter_box_source_table_added_at_idx;
DROP INDEX IF EXISTS migration_dead_letter_box_retryable_retried_at_idx;
DROP FUNCTION IF EXISTS insert_migration_dead_letter_box(TEXT, BIGINT, BYTEA, TEXT, BOOLEAN);
DROP FUNCTION IF EXISTS delete_migration_dead_letter_box(TEXT, BIGINT);

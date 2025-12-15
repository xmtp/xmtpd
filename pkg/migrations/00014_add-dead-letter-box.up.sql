CREATE TABLE IF NOT EXISTS migration_dead_letter_box(
	source_table TEXT NOT NULL,
	sequence_id BIGINT NOT NULL,
    payload BYTEA NOT NULL,
    reason TEXT NOT NULL,
	retryable BOOLEAN NOT NULL DEFAULT FALSE,
	added_at TIMESTAMP NOT NULL DEFAULT NOW(),
	retried_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (source_table, sequence_id)
);

-- Index for reports: query all failures for a source_table, ordered by added_at.
CREATE INDEX IF NOT EXISTS migration_dead_letter_box_source_table_added_at_idx
    ON migration_dead_letter_box (source_table, added_at);

-- Index for retry worker: query retryable records ordered by retried_at (oldest first).
CREATE INDEX IF NOT EXISTS migration_dead_letter_box_retryable_retried_at_idx
    ON migration_dead_letter_box (retried_at)
    WHERE retryable = TRUE;

CREATE FUNCTION insert_migration_dead_letter_box(source_table TEXT, sequence_id BIGINT, payload BYTEA, reason TEXT, retryable BOOLEAN)
	RETURNS SETOF migration_dead_letter_box
	AS $$
BEGIN
	PERFORM
		pg_advisory_xact_lock(hashtext('migration_dead_letter_box_sequence'));
	RETURN QUERY INSERT INTO migration_dead_letter_box(source_table, sequence_id, payload, reason, retryable)
		VALUES(source_table, sequence_id, payload, reason, retryable)
	ON CONFLICT (source_table, sequence_id)
		DO UPDATE SET
		    reason = EXCLUDED.reason,
        	payload = EXCLUDED.payload,
        	retryable = EXCLUDED.retryable,
        	retried_at = NOW()
	RETURNING
		*;
END;
$$
LANGUAGE plpgsql;

CREATE FUNCTION delete_migration_dead_letter_box(source_table TEXT, sequence_id BIGINT)
	RETURNS BOOLEAN
	AS $$
DECLARE
	deleted_count INTEGER;
BEGIN
	PERFORM
		pg_advisory_xact_lock(hashtext('migration_dead_letter_box_sequence'));
	DELETE FROM migration_dead_letter_box
	WHERE migration_dead_letter_box.source_table = delete_migration_dead_letter_box.source_table
		AND migration_dead_letter_box.sequence_id = delete_migration_dead_letter_box.sequence_id;
	GET DIAGNOSTICS deleted_count = ROW_COUNT;
	RETURN deleted_count > 0;
END;
$$
LANGUAGE plpgsql;

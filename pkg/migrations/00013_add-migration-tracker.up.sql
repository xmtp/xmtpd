CREATE TABLE migration_tracker(
	source_table TEXT NOT NULL PRIMARY KEY,
	last_migrated_id BIGINT NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO migration_tracker (source_table, last_migrated_id) VALUES
	('group_messages', 0),
	('inbox_log', 0),
	('installations', 0),
	('welcome_messages', 0); 
package migrator

import (
	"context"
	"database/sql"
)

// DBReader provides a generic implementation for fetching records from database tables.
type DBReader[T ISourceRecord] struct {
	db      *sql.DB
	query   string
	factory func() T
}

// NewDBReader creates a new reader for the specified table and type.
func NewDBReader[T ISourceRecord](
	db *sql.DB,
	query string,
	factory func() T,
) *DBReader[T] {
	return &DBReader[T]{
		db:      db,
		query:   query,
		factory: factory,
	}
}

// Fetch rows from the database and return a slice of records.
func (r *DBReader[T]) Fetch(
	ctx context.Context,
	lastID int64,
	limit int32,
) ([]ISourceRecord, int64, error) {
	rows, err := r.db.QueryContext(ctx, r.query, lastID, limit)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var (
		records = make([]ISourceRecord, 0, limit)
		maxID   int64
	)

	for rows.Next() {
		record := r.factory()

		if err := record.Scan(rows); err != nil {
			return nil, 0, err
		}

		records = append(records, record)

		if record.GetID() > maxID {
			maxID = record.GetID()
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, maxID, nil
}

type GroupMessageReader struct {
	*DBReader[*GroupMessage]
}

func NewGroupMessageReader(db *sql.DB) *GroupMessageReader {
	query := `
		SELECT id, created_at, group_id, data, group_id_data_hash
		FROM group_messages
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2
	`
	return &GroupMessageReader{
		DBReader: NewDBReader(
			db,
			query,
			func() *GroupMessage { return &GroupMessage{} },
		),
	}
}

type InboxLogReader struct {
	*DBReader[*InboxLog]
}

func NewInboxLogReader(db *sql.DB) *InboxLogReader {
	query := `
		SELECT sequence_id, inbox_id, server_timestamp_ns, identity_update_proto
		FROM inbox_log
		WHERE sequence_id > $1
		ORDER BY sequence_id ASC
		LIMIT $2
	`
	return &InboxLogReader{
		DBReader: NewDBReader(
			db,
			query,
			func() *InboxLog { return &InboxLog{} },
		),
	}
}

type InstallationReader struct {
	*DBReader[*Installation]
}

func NewInstallationReader(db *sql.DB) *InstallationReader {
	query := `
		SELECT id, created_at, updated_at, key_package
		FROM installations
		WHERE created_at > $1
		ORDER BY created_at ASC
		LIMIT $2
	`
	return &InstallationReader{
		DBReader: NewDBReader(
			db,
			query,
			func() *Installation { return &Installation{} },
		),
	}
}

type WelcomeMessageReader struct {
	*DBReader[*WelcomeMessage]
}

func NewWelcomeMessageReader(db *sql.DB) *WelcomeMessageReader {
	query := `
		SELECT id, created_at, installation_key, data, hpke_public_key, installation_key_data_hash, wrapper_algorithm
		FROM welcome_messages
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2
	`
	return &WelcomeMessageReader{
		DBReader: NewDBReader(
			db,
			query,
			func() *WelcomeMessage { return &WelcomeMessage{} },
		),
	}
}

package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

// migrateTableBatch processes a batch of records for a specific table.
func (s *dbMigrator) nextRecords(
	ctx context.Context,
	log *zap.Logger,
	dstQueries *queries.Queries,
	tableName string,
) ([]ISourceRecord, int64, error) {
	// Get migration progress for current table.
	lastMigratedID, err := dstQueries.GetMigrationProgress(ctx, tableName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get migration progress: %w", err)
	}

	log.Debug(
		"getting next batch of records",
		zap.Int64("lastMigratedID", lastMigratedID),
	)

	// Get reader for current table.
	reader, ok := s.readers[tableName]
	if !ok {
		return nil, 0, fmt.Errorf("unknown table: %s", tableName)
	}

	// Fetch next batch of records from source database.
	records, newLastID, err := reader.Fetch(ctx, lastMigratedID, s.batchSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrNoRows
		}

		return nil, 0, fmt.Errorf("failed to fetch batch from source database: %w", err)
	}

	log.Debug(
		"fetched batch of records",
		zap.Int("count", len(records)),
		zap.Int64("lastID", newLastID),
	)

	if len(records) == 0 {
		return nil, 0, ErrNoRows
	}

	return records, newLastID, nil
}

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

// TODO: Review all DB queries and make sure they are correct.

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

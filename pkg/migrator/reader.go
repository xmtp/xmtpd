package migrator

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
)

// DBReader provides a generic implementation for fetching records from database tables.
type DBReader[T ISourceRecord] struct {
	db          *sql.DB
	query       string
	queryHeight string
	factory     func() T
	startDate   int64

	// height metric throttling
	heightEvery  time.Duration
	heightMu     sync.RWMutex
	lastHeightAt time.Time
}

// NewDBReader creates a new reader for the specified table and type.
func NewDBReader[T ISourceRecord](
	db *sql.DB,
	query string,
	queryHeight string,
	factory func() T,
) *DBReader[T] {
	return &DBReader[T]{
		db:          db,
		query:       query,
		queryHeight: queryHeight,
		factory:     factory,
		heightEvery: 10 * time.Minute,
	}
}

// maybeEmitHeight updates the source height metric at most once per r.heightEvery.
func (r *DBReader[T]) maybeEmitHeight(ctx context.Context) {
	now := time.Now()

	// 1) Fast path: read-lock to check if we should emit.
	r.heightMu.RLock()
	last := r.lastHeightAt
	shouldEmit := last.IsZero() || now.Sub(last) >= r.heightEvery
	r.heightMu.RUnlock()

	if !shouldEmit {
		return
	}

	// 2) Slow path: write-lock to "claim" the emit, re-check to avoid races.
	r.heightMu.Lock()
	prev := r.lastHeightAt
	shouldEmit = prev.IsZero() || now.Sub(prev) >= r.heightEvery
	if shouldEmit {
		r.lastHeightAt = now
	}
	r.heightMu.Unlock()

	if !shouldEmit {
		return
	}

	// 3) Do the DB work outside locks.
	var heightID int64
	if err := r.db.QueryRowContext(ctx, r.queryHeight).Scan(&heightID); err != nil {
		// Roll back claim so a later call can retry.
		r.heightMu.Lock()
		// Only roll back if nobody else advanced it since we claimed.
		// (This should almost always be true, but it's cheap to be safe.)
		if r.lastHeightAt.Equal(now) {
			r.lastHeightAt = prev
		}
		r.heightMu.Unlock()
		return
	}

	metrics.EmitMigratorSourceLastSequenceID(r.factory().TableName(), heightID)
}

// Fetch rows from the database and return a slice of records.
func (r *DBReader[T]) Fetch(
	ctx context.Context,
	lastID int64,
	limit int32,
) ([]ISourceRecord, error) {
	r.maybeEmitHeight(ctx)

	var (
		rows *sql.Rows
		err  error
	)

	if r.startDate != 0 {
		rows, err = r.db.QueryContext(ctx, r.query, lastID, limit, r.startDate)
	} else {
		rows, err = r.db.QueryContext(ctx, r.query, lastID, limit)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	records := make([]ISourceRecord, 0, limit)
	for rows.Next() {
		record := r.factory()
		if err := record.Scan(rows); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

type GroupMessageReader struct {
	*DBReader[*GroupMessage]
}

func NewGroupMessageReader(db *sql.DB) *GroupMessageReader {
	query := `
		SELECT id, created_at, group_id, data, group_id_data_hash, is_commit, sender_hmac, should_push
		FROM group_messages
		WHERE id > $1 AND is_commit = false
		ORDER BY id ASC
		LIMIT $2
	`

	queryHeight := `
		SELECT id
		FROM group_messages
		WHERE is_commit = false
		ORDER BY id DESC
		LIMIT 1
	`
	return &GroupMessageReader{
		DBReader: NewDBReader(
			db,
			query,
			queryHeight,
			func() *GroupMessage { return &GroupMessage{} },
		),
	}
}

type CommitMessageReader struct {
	*DBReader[*CommitMessage]
}

func NewCommitMessageReader(db *sql.DB) *CommitMessageReader {
	query := `
		SELECT id, created_at, group_id, data, group_id_data_hash, is_commit, sender_hmac, should_push
		FROM group_messages
		WHERE id > $1 AND is_commit = true
		ORDER BY id ASC
		LIMIT $2
	`

	queryHeight := `
		SELECT id
		FROM group_messages
		WHERE is_commit = true
		ORDER BY id DESC
		LIMIT 1
	`

	return &CommitMessageReader{
		DBReader: NewDBReader(
			db,
			query,
			queryHeight,
			func() *CommitMessage { return &CommitMessage{} },
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

	queryHeight := `
		SELECT sequence_id
		FROM inbox_log
		ORDER BY sequence_id DESC
		LIMIT 1
	`

	return &InboxLogReader{
		DBReader: NewDBReader(
			db,
			query,
			queryHeight,
			func() *InboxLog { return &InboxLog{} },
		),
	}
}

type KeyPackageReader struct {
	*DBReader[*KeyPackage]
}

func NewKeyPackageReader(db *sql.DB) *KeyPackageReader {
	query := `
		SELECT sequence_id, installation_id, key_package, created_at
		FROM key_packages
		WHERE sequence_id > $1
		ORDER BY sequence_id ASC
		LIMIT $2
	`

	queryHeight := `
		SELECT sequence_id
		FROM key_packages
		ORDER BY sequence_id DESC
		LIMIT 1
	`

	return &KeyPackageReader{
		DBReader: NewDBReader(
			db,
			query,
			queryHeight,
			func() *KeyPackage { return &KeyPackage{} },
		),
	}
}

type WelcomeMessageReader struct {
	*DBReader[*WelcomeMessage]
}

func NewWelcomeMessageReader(db *sql.DB) *WelcomeMessageReader {
	query := `
		SELECT id, created_at, installation_key, data, hpke_public_key, installation_key_data_hash, wrapper_algorithm, welcome_metadata
		FROM welcome_messages
		WHERE id > 150000000 AND id > $1
		ORDER BY id ASC
		LIMIT $2
	`

	queryHeight := `
		SELECT id
		FROM welcome_messages
		ORDER BY id DESC
		LIMIT 1
	`

	return &WelcomeMessageReader{
		DBReader: NewDBReader(
			db,
			query,
			queryHeight,
			func() *WelcomeMessage { return &WelcomeMessage{} },
		),
	}
}

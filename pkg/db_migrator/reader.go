package db_migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

// TODO: Refine querying and scanning of records.

// migrateTableBatch processes a batch of records for a specific table.
func (s *dbMigrator) nextRecords(
	ctx context.Context,
	log *zap.Logger,
	dstQueries *queries.Queries,
	tableName string,
) ([]Record, int64, error) {
	// Get migration progress for current table.
	lastMigratedID, err := dstQueries.GetMigrationProgress(ctx, tableName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get migration progress: %w", err)
	}

	// Get next batch of records from source database.
	records, newLastID, err := s.getNextBatch(ctx, log, tableName, lastMigratedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrNoRows
		}

		return nil, 0, fmt.Errorf("failed to fetch batch from source database: %w", err)
	}

	if len(records) == 0 {
		return nil, 0, ErrNoRows
	}

	return records, newLastID, nil
}

// getNextBatch fetches a batch of records from the source database.
func (s *dbMigrator) getNextBatch(
	ctx context.Context,
	log *zap.Logger,
	tableName string,
	lastID int64,
) ([]Record, int64, error) {
	var (
		records = make([]Record, 0, s.batchSize)
		maxID   int64
	)

	log.Debug("fetching next batch", zap.String("table", tableName), zap.Int64("lastID", lastID))

	switch tableName {
	case addressLogTableName:
		query := "SELECT id, address, inbox_id, association_sequence_id, revocation_sequence_id FROM address_log WHERE id > $1 ORDER BY id ASC LIMIT $2"
		rows, err := s.src.QueryContext(ctx, query, lastID, s.batchSize)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			var record AddressLog
			err := rows.Scan(
				&record.ID,
				&record.Address,
				&record.InboxID,
				&record.AssociationSequenceID,
				&record.RevocationSequenceID,
			)
			if err != nil {
				return nil, 0, err
			}
			records = append(records, record)
			if record.ID > maxID {
				maxID = record.ID
			}
		}

	case groupMessagesTableName:
		query := "SELECT id, created_at, group_id, data, group_id_data_hash FROM group_messages WHERE id > $1 ORDER BY id ASC LIMIT $2"
		rows, err := s.src.QueryContext(ctx, query, lastID, s.batchSize)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			var record GroupMessage
			err := rows.Scan(
				&record.ID,
				&record.CreatedAt,
				&record.GroupID,
				&record.Data,
				&record.GroupIDDataHash,
			)
			if err != nil {
				return nil, 0, err
			}
			records = append(records, record)
			if record.ID > maxID {
				maxID = record.ID
			}
		}

	case inboxLogTableName:
		query := "SELECT sequence_id, inbox_id, server_timestamp_ns, identity_update_proto FROM inbox_log WHERE sequence_id > $1 ORDER BY sequence_id ASC LIMIT $2"
		rows, err := s.src.QueryContext(ctx, query, lastID, s.batchSize)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			var record InboxLog
			err := rows.Scan(
				&record.SequenceID,
				&record.InboxID,
				&record.ServerTimestampNs,
				&record.IdentityUpdateProto,
			)
			if err != nil {
				return nil, 0, err
			}
			records = append(records, record)
			if record.SequenceID > maxID {
				maxID = record.SequenceID
			}
		}

	case installationsTableName:
		query := "SELECT id, created_at, updated_at, key_package FROM installations WHERE created_at > $1 ORDER BY created_at ASC LIMIT $2"
		rows, err := s.src.QueryContext(ctx, query, lastID, s.batchSize)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			var record Installation
			err := rows.Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt, &record.KeyPackage)
			if err != nil {
				return nil, 0, err
			}
			records = append(records, record)
			if record.CreatedAt > maxID {
				maxID = record.CreatedAt
			}
		}

	case welcomeMessagesTableName:
		query := "SELECT id, created_at, installation_key, data, hpke_public_key, installation_key_data_hash, wrapper_algorithm FROM welcome_messages WHERE id > $1 ORDER BY id ASC LIMIT $2"
		rows, err := s.src.QueryContext(ctx, query, lastID, s.batchSize)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			var record WelcomeMessage
			err := rows.Scan(
				&record.ID,
				&record.CreatedAt,
				&record.InstallationKey,
				&record.Data,
				&record.HpkePublicKey,
				&record.InstallationKeyDataHash,
				&record.WrapperAlgorithm,
			)
			if err != nil {
				return nil, 0, err
			}
			records = append(records, record)
			if record.ID > maxID {
				maxID = record.ID
			}
		}

	default:
		return nil, 0, fmt.Errorf("unknown table: %s", tableName)
	}

	return records, maxID, nil
}

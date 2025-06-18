package db_migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// AddressLog represents the address_log table from the source database.
type AddressLog struct {
	ID                    int64         `db:"id"`
	Address               string        `db:"address"`
	InboxID               []byte        `db:"inbox_id"`
	AssociationSequenceID sql.NullInt64 `db:"association_sequence_id"`
	RevocationSequenceID  sql.NullInt64 `db:"revocation_sequence_id"`
}

// GroupMessage represents the group_messages table from the source database.
type GroupMessage struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	GroupID         []byte    `db:"group_id"`
	Data            []byte    `db:"data"`
	GroupIDDataHash []byte    `db:"group_id_data_hash"`
}

// InboxLog represents the inbox_log table from the source database.
type InboxLog struct {
	SequenceID          int64  `db:"sequence_id"`
	InboxID             []byte `db:"inbox_id"`
	ServerTimestampNs   int64  `db:"server_timestamp_ns"`
	IdentityUpdateProto []byte `db:"identity_update_proto"`
}

// Installation represents the installations table from the source database.
type Installation struct {
	ID         []byte `db:"id"`
	CreatedAt  int64  `db:"created_at"`
	UpdatedAt  int64  `db:"updated_at"`
	KeyPackage []byte `db:"key_package"`
}

// WelcomeMessage represents the welcome_messages table from the source database.
type WelcomeMessage struct {
	ID                      int64     `db:"id"`
	CreatedAt               time.Time `db:"created_at"`
	InstallationKey         []byte    `db:"installation_key"`
	Data                    []byte    `db:"data"`
	HpkePublicKey           []byte    `db:"hpke_public_key"`
	InstallationKeyDataHash []byte    `db:"installation_key_data_hash"`
	WrapperAlgorithm        int16     `db:"wrapper_algorithm"`
}

// TODO: Refine querying and scanning of records.

// getNextBatch fetches a batch of records from the source database.
func (s *dbMigrator) getNextBatch(
	ctx context.Context,
	log *zap.Logger,
	tableName string,
	lastID int64,
) ([]interface{}, int64, error) {
	var (
		records []interface{}
		maxID   int64
	)

	log.Debug("fetching next batch", zap.String("table", tableName), zap.Int64("lastID", lastID))

	switch tableName {
	case "address_log":
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

	case "group_messages":
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

	case "inbox_log":
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

	case "installations":
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

	case "welcome_messages":
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

package db_migrator

import (
	"context"
	"database/sql"
	"time"

	"github.com/xmtp/xmtpd/pkg/envelopes"
)

const (
	addressLogTableName      = "address_log"
	groupMessagesTableName   = "group_messages"
	inboxLogTableName        = "inbox_log"
	installationsTableName   = "installations"
	welcomeMessagesTableName = "welcome_messages"
)

// Reader defines the interface for reading records from the source database.
type Reader interface {
	Fetch(ctx context.Context, lastID int64, limit int32) ([]Record, int64, error)
}

// Transformer defines the interface for transforming external data to xmtpd OriginatorEnvelope format.
type Transformer interface {
	Transform(record Record) (*envelopes.OriginatorEnvelope, error)
}

// Record defines a record from the source database,
// that can scanned and ordered by some ID.
type Record interface {
	GetID() int64
	TableName() string
	Scan(rows *sql.Rows) error
}

// TODO: Probably not needed if they can be derived from InboxLog (IdentityUpdates).
// AddressLog represents the address_log table from the source database.
// Order by association_sequence_id ASC.
type AddressLog struct {
	ID                    int64         `db:"id"`
	Address               string        `db:"address"`
	InboxID               []byte        `db:"inbox_id"`
	AssociationSequenceID sql.NullInt64 `db:"association_sequence_id"`
	RevocationSequenceID  sql.NullInt64 `db:"revocation_sequence_id"`
}

func (a AddressLog) GetID() int64 {
	return a.AssociationSequenceID.Int64
}

func (a AddressLog) TableName() string {
	return addressLogTableName
}

func (a *AddressLog) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&a.ID,
		&a.Address,
		&a.InboxID,
		&a.AssociationSequenceID,
		&a.RevocationSequenceID,
	)
}

// GroupMessage represents the group_messages table from the source database.
// Order by ID ASC.
type GroupMessage struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	GroupID         []byte    `db:"group_id"`
	Data            []byte    `db:"data"`
	GroupIDDataHash []byte    `db:"group_id_data_hash"`
}

func (g GroupMessage) GetID() int64 {
	return g.ID
}

func (g GroupMessage) TableName() string {
	return groupMessagesTableName
}

func (g *GroupMessage) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&g.ID,
		&g.CreatedAt,
		&g.GroupID,
		&g.Data,
		&g.GroupIDDataHash,
	)
}

// InboxLog represents the inbox_log table from the source database.
// Order by SequenceID ASC.
type InboxLog struct {
	SequenceID          int64  `db:"sequence_id"`
	InboxID             []byte `db:"inbox_id"`
	ServerTimestampNs   int64  `db:"server_timestamp_ns"`
	IdentityUpdateProto []byte `db:"identity_update_proto"`
}

func (i InboxLog) GetID() int64 {
	return i.SequenceID
}

func (i InboxLog) TableName() string {
	return inboxLogTableName
}

func (i *InboxLog) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&i.SequenceID,
		&i.InboxID,
		&i.ServerTimestampNs,
		&i.IdentityUpdateProto,
	)
}

// Installation represents the installations table from the source database.
// Order by CreatedAt ASC.
type Installation struct {
	ID         []byte `db:"id"`
	CreatedAt  int64  `db:"created_at"`
	UpdatedAt  int64  `db:"updated_at"`
	KeyPackage []byte `db:"key_package"`
}

func (i Installation) GetID() int64 {
	return i.CreatedAt
}

func (i Installation) TableName() string {
	return installationsTableName
}

func (i *Installation) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.KeyPackage,
	)
}

// WelcomeMessage represents the welcome_messages table from the source database.
// Order by ID ASC.
type WelcomeMessage struct {
	ID                      int64     `db:"id"`
	CreatedAt               time.Time `db:"created_at"`
	InstallationKey         []byte    `db:"installation_key"`
	Data                    []byte    `db:"data"`
	HpkePublicKey           []byte    `db:"hpke_public_key"`
	InstallationKeyDataHash []byte    `db:"installation_key_data_hash"`
	WrapperAlgorithm        int16     `db:"wrapper_algorithm"`
}

func (w WelcomeMessage) GetID() int64 {
	return w.ID
}

func (w WelcomeMessage) TableName() string {
	return welcomeMessagesTableName
}

func (w *WelcomeMessage) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&w.ID,
		&w.CreatedAt,
		&w.InstallationKey,
		&w.Data,
		&w.HpkePublicKey,
		&w.InstallationKeyDataHash,
		&w.WrapperAlgorithm,
	)
}

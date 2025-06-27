package migrator

import (
	"context"
	"database/sql"
	"time"

	"github.com/xmtp/xmtpd/pkg/envelopes"
)

const (
	groupMessagesTableName          = "group_messages"
	groupMessageOriginatorID uint32 = 10

	welcomeMessagesTableName          = "welcome_messages"
	welcomeMessageOriginatorID uint32 = 11

	// IdentityUpdates in xmtpd.
	inboxLogTableName           = "inbox_log"
	inboxLogOriginatorID uint32 = 12

	// KeyPackages in xmtpd.
	installationsTableName          = "installations"
	installationOriginatorID uint32 = 13
)

var originatorIDToTableName = map[uint32]string{
	groupMessageOriginatorID:   groupMessagesTableName,
	welcomeMessageOriginatorID: welcomeMessagesTableName,
	inboxLogOriginatorID:       inboxLogTableName,
	installationOriginatorID:   installationsTableName,
}

// IDataTransformer defines the interface for transforming external data to xmtpd OriginatorEnvelope format.
type IDataTransformer interface {
	Transform(record ISourceRecord) (*envelopes.OriginatorEnvelope, error)
}

type IDestinationWriter interface {
	Write(ctx context.Context, env *envelopes.OriginatorEnvelope) error
}

// ISourceReader defines the interface for reading records from the source database.
type ISourceReader interface {
	Fetch(ctx context.Context, lastID int64, limit int32) ([]ISourceRecord, int64, error)
}

// ISourceRecord defines a record from the source database,
// that can scanned and ordered by some ID.
type ISourceRecord interface {
	GetID() int64
	TableName() string
	Scan(rows *sql.Rows) error
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

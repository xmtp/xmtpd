package migrator

import (
	"context"
	"database/sql"
	"errors"
	"iter"
	"time"

	"github.com/xmtp/xmtpd/pkg/envelopes"
)

const (
	groupMessagesTableName   = "group_messages"
	welcomeMessagesTableName = "welcome_messages"
	inboxLogTableName        = "inbox_log"
	keyPackagesTableName     = "key_packages"
	commitMessagesTableName  = "commit_messages"

	GroupMessageOriginatorID   uint32 = 10
	WelcomeMessageOriginatorID uint32 = 11
	InboxLogOriginatorID       uint32 = 12
	KeyPackagesOriginatorID    uint32 = 13
	CommitMessageOriginatorID  uint32 = 14
)

var originatorIDToTableName = map[uint32]string{
	GroupMessageOriginatorID:   groupMessagesTableName,
	WelcomeMessageOriginatorID: welcomeMessagesTableName,
	InboxLogOriginatorID:       inboxLogTableName,
	KeyPackagesOriginatorID:    keyPackagesTableName,
	CommitMessageOriginatorID:  commitMessagesTableName,
}

func isValidOriginatorID(originatorID uint32) bool {
	return originatorID == GroupMessageOriginatorID ||
		originatorID == WelcomeMessageOriginatorID ||
		originatorID == InboxLogOriginatorID ||
		originatorID == KeyPackagesOriginatorID ||
		originatorID == CommitMessageOriginatorID
}

func isDatabaseDestination(originatorID uint32) bool {
	return originatorID == GroupMessageOriginatorID ||
		originatorID == WelcomeMessageOriginatorID ||
		originatorID == KeyPackagesOriginatorID
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
	Fetch(ctx context.Context, lastID int64, limit int32) ([]ISourceRecord, error)
}

// ISourceRecord defines a record from the source database,
// that can scanned and ordered by some ID.
type ISourceRecord interface {
	GetID() int64
	TableName() string
	Scan(rows *sql.Rows) error
}

// GroupMessage represents the group_messages (with is_commit = false) from the source database.
// Order by ID ASC.
type GroupMessage struct {
	ID              int64        `db:"id"`
	CreatedAt       time.Time    `db:"created_at"`
	GroupID         []byte       `db:"group_id"`
	Data            []byte       `db:"data"`
	GroupIDDataHash []byte       `db:"group_id_data_hash"`
	IsCommit        sql.NullBool `db:"is_commit"`
	SenderHmac      []byte       `db:"sender_hmac"`
	ShouldPush      sql.NullBool `db:"should_push"`
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
		&g.IsCommit,
		&g.SenderHmac,
		&g.ShouldPush,
	)
}

// CommitMessage represents the group_messages (with is_commit = true) from the source database.
// Order by ID ASC.
type CommitMessage struct {
	GroupMessage
}

func (c CommitMessage) GetID() int64 {
	return c.ID
}

func (c CommitMessage) TableName() string {
	return commitMessagesTableName
}

func (c *CommitMessage) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&c.ID,
		&c.CreatedAt,
		&c.GroupID,
		&c.Data,
		&c.GroupIDDataHash,
		&c.IsCommit,
		&c.SenderHmac,
		&c.ShouldPush,
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

// KeyPackage represents the key_packages table from the source database.
// Order by CreatedAt ASC.
type KeyPackage struct {
	SequenceID     int64  `db:"sequence_id"`
	InstallationID []byte `db:"installation_id"`
	KeyPackage     []byte `db:"key_package"`
}

func (i KeyPackage) GetID() int64 {
	return i.SequenceID
}

func (i KeyPackage) TableName() string {
	return keyPackagesTableName
}

func (i *KeyPackage) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&i.SequenceID,
		&i.InstallationID,
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
	WelcomeMetadata         []byte    `db:"welcome_metadata"`
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
		&w.WelcomeMetadata,
	)
}

// FailureReason defines the reason for inserting a record into the dead letter box.
type FailureReason string

const (
	FailureTransformerError       FailureReason = "transformer error"
	FailureOversizedChainMessage  FailureReason = "oversized chain message"
	FailureBlockchainUndetermined FailureReason = "blockchain undetermined error"
)

var ErrDeadLetterBox = errors.New("skipped and added to dead letter box")

func (f FailureReason) String() string {
	return string(f)
}

func (f FailureReason) ShouldRetry() bool {
	switch f {
	case FailureTransformerError:
		return true
	case FailureOversizedChainMessage:
		return false
	default:
		return false
	}
}

// IdentityUpdateBatch defines a batch of identity updates to be published to the blockchain.
type IdentityUpdateBatch struct {
	inboxIDs        [][32]byte
	identityUpdates [][]byte
	sequenceIDs     []uint64
}

// IdentityUpdate represents a single element from the batch.
type IdentityUpdate struct {
	InboxID        [32]byte
	IdentityUpdate []byte
	SequenceID     uint64
}

// All returns an iterator over all items in the batch.
func (i *IdentityUpdateBatch) All() iter.Seq[IdentityUpdate] {
	return func(yield func(IdentityUpdate) bool) {
		for index := range i.inboxIDs {
			if !yield(IdentityUpdate{
				InboxID:        i.inboxIDs[index],
				IdentityUpdate: i.identityUpdates[index],
				SequenceID:     i.sequenceIDs[index],
			}) {
				return
			}
		}
	}
}

func (i *IdentityUpdateBatch) Size() int64 {
	size := 0

	for _, identityUpdate := range i.identityUpdates {
		size += len(identityUpdate)
	}

	return int64(len(i.inboxIDs)*32 + size + len(i.sequenceIDs)*8)
}

func (i *IdentityUpdateBatch) LastSequenceID() uint64 {
	return i.sequenceIDs[len(i.sequenceIDs)-1]
}

func (i *IdentityUpdateBatch) Add(inboxID [32]byte, payload []byte, sequenceID uint64) {
	i.inboxIDs = append(i.inboxIDs, inboxID)
	i.identityUpdates = append(i.identityUpdates, payload)
	i.sequenceIDs = append(i.sequenceIDs, sequenceID)
}

func (i *IdentityUpdateBatch) Len() int {
	return len(i.inboxIDs)
}

func (i *IdentityUpdateBatch) Reset() {
	i.inboxIDs = nil
	i.identityUpdates = nil
	i.sequenceIDs = nil
}

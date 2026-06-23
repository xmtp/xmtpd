// Package db provides DB access and advisory-lock helpers for coordinating HA workers.
// It includes two styles:
//   - AdvisoryLocker: inherits the caller's transaction/connection.
//   - TransactionScopedAdvisoryLocker: creates and owns its own transaction.
package db

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// LockKind marks the lowest 8 bits of the advisory lock key.
type LockKind uint8

const (
	LockKindIdentityUpdateInsert LockKind = 0x00
	LockKindAttestationWorker    LockKind = 0x01
	LockKindSubmitterWorker      LockKind = 0x02
	LockKindSettlementWorker     LockKind = 0x03
	LockKindGeneratorWorker      LockKind = 0x04
	LockKindPartitionCreation    LockKind = 0x05
)

// partitionCreationLockKey is the single, global advisory-lock key coordinating lazy
// gateway-envelope partition creation against concurrent inserts. It is not per-originator:
// the deadlock it prevents is between transactions touching different originators' partitions,
// so they must contend on the same key.
const partitionCreationLockKey = int64(LockKindPartitionCreation)

// AdvisoryLocker builds advisory-lock keys and acquires locks using the caller’s
// connection/transaction via the provided *queries.Queries
type AdvisoryLocker struct{}

func NewAdvisoryLocker() *AdvisoryLocker {
	return &AdvisoryLocker{}
}

func (a *AdvisoryLocker) LockIdentityUpdateInsert(
	ctx context.Context,
	queries *queries.Queries,
	nodeID uint32,
) error {
	key := int64((uint64(nodeID) << 8) | uint64(LockKindIdentityUpdateInsert))
	return queries.AdvisoryLockWithKey(ctx, key)
}

// SharedLockPartitionCreation takes the shared (reader) side of the partition-creation
// reader/writer advisory lock. Ordinary gateway-envelope inserts hold it so they run
// concurrently with each other, but block (and are blocked by) exclusive partition creation.
func (a *AdvisoryLocker) SharedLockPartitionCreation(
	ctx context.Context,
	queries *queries.Queries,
) error {
	return queries.SharedAdvisoryLockWithKey(ctx, partitionCreationLockKey)
}

// LockPartitionCreation takes the exclusive (writer) side of the partition-creation
// reader/writer advisory lock. It must be acquired in a dedicated transaction (see
// EnsureGatewayPartitions), never upgraded from the shared lock within an insert transaction,
// or two upgraders would deadlock.
//
// ensure_gateway_parts_v3 ATTACHes partitions to the shared gateway_envelopes_meta and
// gateway_envelopes_blob parents. Because of the blob->meta and meta->payers foreign keys, each
// ATTACH validates the FK and takes ShareRowExclusiveLock on both parents once they hold data.
// That conflicts with the RowExclusiveLock ordinary inserts take, and concurrent transactions
// acquire the parents' locks in opposite orders and deadlock (SQLSTATE 40P01). Running partition
// creation as the exclusive writer, while inserts hold the shared lock, guarantees DDL never
// overlaps DML, removing the conflict.
func (a *AdvisoryLocker) LockPartitionCreation(
	ctx context.Context,
	queries *queries.Queries,
) error {
	return queries.AdvisoryLockWithKey(ctx, partitionCreationLockKey)
}

func (a *AdvisoryLocker) TryLockGeneratorWorker(
	ctx context.Context,
	queries *queries.Queries,
) (bool, error) {
	key := int64(LockKindGeneratorWorker)
	return queries.TryAdvisoryLockWithKey(ctx, key)
}

func (a *AdvisoryLocker) TryLockAttestationWorker(
	ctx context.Context,
	queries *queries.Queries,
) (bool, error) {
	key := int64(LockKindAttestationWorker)
	return queries.TryAdvisoryLockWithKey(ctx, key)
}

func (a *AdvisoryLocker) TryLockSubmitterWorker(
	ctx context.Context,
	queries *queries.Queries,
) (bool, error) {
	key := int64(LockKindSubmitterWorker)
	return queries.TryAdvisoryLockWithKey(ctx, key)
}

func (a *AdvisoryLocker) TryLockSettlementWorker(
	ctx context.Context,
	queries *queries.Queries,
) (bool, error) {
	key := int64(LockKindSettlementWorker)
	return queries.TryAdvisoryLockWithKey(ctx, key)
}

type ITransactionScopedAdvisoryLocker interface {
	Release() error
	TryLockGeneratorWorker() (bool, error)
	TryLockAttestationWorker() (bool, error)
	TryLockSubmitterWorker() (bool, error)
	TryLockSettlementWorker() (bool, error)
	LockIdentityUpdateInsert(nodeID uint32) error
}

// TransactionScopedAdvisoryLocker creates and owns a transaction; methods acquire
// advisory locks within that tx. Release() rolls back the tx
type TransactionScopedAdvisoryLocker struct {
	tx     *sql.Tx
	ctx    context.Context
	locker *AdvisoryLocker
}

var _ ITransactionScopedAdvisoryLocker = (*TransactionScopedAdvisoryLocker)(nil)

func NewTransactionScopedAdvisoryLocker(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
) (*TransactionScopedAdvisoryLocker, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TransactionScopedAdvisoryLocker{tx: tx, ctx: ctx, locker: NewAdvisoryLocker()}, nil
}

func (a *TransactionScopedAdvisoryLocker) Release() error {
	return a.tx.Rollback()
}

func (a *TransactionScopedAdvisoryLocker) TryLockGeneratorWorker() (bool, error) {
	return a.locker.TryLockGeneratorWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) TryLockAttestationWorker() (bool, error) {
	return a.locker.TryLockAttestationWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) TryLockSubmitterWorker() (bool, error) {
	return a.locker.TryLockSubmitterWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) TryLockSettlementWorker() (bool, error) {
	return a.locker.TryLockSettlementWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) LockIdentityUpdateInsert(nodeID uint32) error {
	return a.locker.LockIdentityUpdateInsert(a.ctx, queries.New(a.tx), nodeID)
}

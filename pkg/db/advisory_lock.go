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
)

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

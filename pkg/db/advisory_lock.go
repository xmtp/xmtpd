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
	nodeId uint32,
) error {
	key := int64((uint64(nodeId) << 8) | uint64(LockKindIdentityUpdateInsert))
	return queries.AdvisoryLockWithKey(ctx, key)
}

func (a *AdvisoryLocker) LockAttestationWorker(
	ctx context.Context,
	queries *queries.Queries,
) error {
	key := int64(LockKindAttestationWorker)
	return queries.AdvisoryLockWithKey(ctx, key)
}

func (a *AdvisoryLocker) LockSubmitterWorker(
	ctx context.Context,
	queries *queries.Queries,
) error {
	key := int64(LockKindSubmitterWorker)
	return queries.AdvisoryLockWithKey(ctx, key)
}

type ITransactionScopedAdvisoryLocker interface {
	Release() error
	LockAttestationWorker() error
	LockSubmitterWorker() error
	LockIdentityUpdateInsert(nodeId uint32) error
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

func (a *TransactionScopedAdvisoryLocker) LockAttestationWorker() error {
	return a.locker.LockAttestationWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) LockSubmitterWorker() error {
	return a.locker.LockSubmitterWorker(a.ctx, queries.New(a.tx))
}

func (a *TransactionScopedAdvisoryLocker) LockIdentityUpdateInsert(nodeId uint32) error {
	return a.locker.LockIdentityUpdateInsert(a.ctx, queries.New(a.tx), nodeId)
}

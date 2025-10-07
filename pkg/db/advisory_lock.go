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
)

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

type ITransactionScopedAdvisoryLocker interface {
	Release() error
	LockAttestationWorker() error
	LockIdentityUpdateInsert(nodeId uint32) error
}

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

func (a *TransactionScopedAdvisoryLocker) LockIdentityUpdateInsert(nodeId uint32) error {
	return a.locker.LockIdentityUpdateInsert(a.ctx, queries.New(a.tx), nodeId)
}

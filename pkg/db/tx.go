package db

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

func RunInTxRaw(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, tx *sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	done = true
	return tx.Commit()
}

func RunInTx(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) error,
) error {
	return RunInTxRaw(ctx, db, opts, func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, queries.New(tx))
	})
}

func RunInTxWithResult[T any](
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) (T, error),
) (result T, err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return result, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	result, err = fn(ctx, queries.New(tx))
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

// RunInTxRaw runs fn inside a transaction using the handler's write DB.
func (h *Handler) RunInTxRaw(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(ctx context.Context, tx *sql.Tx) error,
) error {
	return RunInTxRaw(ctx, h.write, opts, fn)
}

// RunInTx runs fn inside a transaction, providing prepared queries bound to the transaction.
// This allows PostgreSQL to reuse cached query plans for all statements in the transaction.
func (h *Handler) RunInTx(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) error,
) error {
	return RunInTxRaw(ctx, h.write, opts, func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, h.query.WithTx(tx))
	})
}

// HandlerRunInTxWithResult runs fn inside a transaction on h and returns a result.
// It provides prepared queries bound to the transaction so PostgreSQL can reuse cached query plans.
// Note: Go does not allow generic methods on non-generic types, so this is a package-level function.
func HandlerRunInTxWithResult[T any](
	ctx context.Context,
	h *Handler,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) (T, error),
) (result T, err error) {
	tx, err := h.write.BeginTx(ctx, opts)
	if err != nil {
		return result, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	result, err = fn(ctx, h.query.WithTx(tx))
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

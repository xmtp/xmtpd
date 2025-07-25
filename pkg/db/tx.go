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

	result, err = fn(ctx, queries.New(db).WithTx(tx))
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

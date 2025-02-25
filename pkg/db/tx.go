package db

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

func RunInTx(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) error,
) error {
	querier := queries.New(db)
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

	if err := fn(ctx, querier.WithTx(tx)); err != nil {
		return err
	}

	done = true
	return tx.Commit()
}

func RunInTxWithResult[T any](
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) (T, error),
) (T, error) {
	querier := queries.New(db)
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		var zero T
		return zero, err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	result, err := fn(ctx, querier.WithTx(tx))
	if err != nil {
		var zero T
		return zero, err
	}

	done = true
	if err := tx.Commit(); err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}

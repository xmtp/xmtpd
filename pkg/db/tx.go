package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type txConfig struct {
	onCommit   []func()
	onRollback []func()
}

func newTxConfig() *txConfig {
	return &txConfig{
		onCommit:   make([]func(), 0),
		onRollback: make([]func(), 0),
	}
}

type txOption func(*txConfig)

func OnCommit(hook func()) txOption {
	return func(cfg *txConfig) {
		cfg.onCommit = append(cfg.onCommit, hook)
	}
}

func OnRollback(hook func()) txOption {
	return func(cfg *txConfig) {
		cfg.onRollback = append(cfg.onRollback, hook)
	}
}

func RunInTxRaw(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, tx *sql.Tx) error,
	options ...txOption,
) error {
	cfg := newTxConfig()
	for _, opt := range options {
		opt(cfg)
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}

	defer func() {
		// If commit succeeded, rollback fails so it's a no-op.
		err := tx.Rollback()
		if err == nil {
			// Rollback succeeded - invoke rollback callbacks.
			for _, cb := range cfg.onRollback {
				cb()
			}
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		return fmt.Errorf("could not execute transaction function: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	// Commit succeeded - invoke commit callbacks.
	for _, cb := range cfg.onCommit {
		cb()
	}

	return nil
}

func RunInTx(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) error,
	options ...txOption,
) error {
	return RunInTxRaw(ctx, db, opts, func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, queries.New(tx))
	},
		options...,
	)
}

func RunInTxWithResult[T any](
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, txQueries *queries.Queries) (T, error),
	options ...txOption,
) (result T, err error) {
	cfg := newTxConfig()
	for _, opt := range options {
		opt(cfg)
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return result, fmt.Errorf("could not start a transaction: %w", err)
	}

	defer func() {
		// If commit succeeded, rollback fails so it's a no-op.
		err := tx.Rollback()
		if err == nil {
			// Rollback succeeded - invoke rollback callbacks.
			for _, cb := range cfg.onRollback {
				cb()
			}
		}
	}()

	result, err = fn(ctx, queries.New(db).WithTx(tx))
	if err != nil {
		return result, fmt.Errorf("could not execute transaction function: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return result, fmt.Errorf("could not commit transaction: %w", err)
	}

	// Commit succeeded - invoke commit callbacks.
	for _, cb := range cfg.onCommit {
		cb()
	}

	return result, nil
}

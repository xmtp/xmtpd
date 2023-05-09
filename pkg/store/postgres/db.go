package postgresstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
)

const (
	// Default in postgres is typically configured to be 100.
	maxOpenConnections = 80
)

type DB struct {
	*sql.DB
	DSN string
}

func NewDB(opts *Options) (*DB, error) {
	db, err := sql.Open("pgx", opts.DSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	return &DB{
		DB:  db,
		DSN: opts.DSN,
	}, nil
}

func executeTx(ctx context.Context, db *DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				if rollbackErr == context.Canceled {
					return
				}
				ctx.Logger().Error("error rolling back", zap.Error(err))
			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			rollbackErr := tx.Rollback() // err is non-nil; don't change it
			if rollbackErr != nil {
				if rollbackErr == context.Canceled {
					return
				}
				ctx.Logger().Error("error rolling back", zap.Error(err))
			}
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = fn(tx)
	return err
}

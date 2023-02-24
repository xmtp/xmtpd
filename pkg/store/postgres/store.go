package postgresstore

import (
	"context"

	"github.com/xmtp/xmtpd/db/migrations"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type PostgresStore struct {
	log *zap.Logger
	db  *DB

	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context, log *zap.Logger, db *DB) (*PostgresStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	s := &PostgresStore{
		ctx:    ctx,
		cancel: cancel,

		log: log.Named("pgstore"),
		db:  db,
	}

	err := migrations.Run(db.DSN)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *PostgresStore) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *PostgresStore) Scoped(topic string) *ScopedPostgresStore {
	return newScoped(s.ctx, s.log, s.db, topic)
}

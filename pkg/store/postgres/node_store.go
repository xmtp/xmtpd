package postgresstore

import (
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/store/postgres/migrations"
)

type NodeStore struct {
	ctx context.Context
	db  *DB
}

func NewNodeStore(ctx context.Context, db *DB) (*NodeStore, error) {
	log := ctx.Logger().Named("store")
	ctx = context.WithLogger(ctx, log)
	s := &NodeStore{
		ctx: ctx,
		db:  db,
	}

	err := migrations.Run(db.DSN)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *NodeStore) Close() error {
	return nil
}

func (s *NodeStore) NewTopic(topic string) (crdt.Store, error) {
	return New(s.ctx, s.db, topic), nil
}

func (s *NodeStore) Topics() (topics []string, err error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(s.ctx, "SELECT topic FROM events GROUP BY topic")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var topic string
		err := rows.Scan(&topic)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func (s *NodeStore) DeleteTopic(topic string) error {
	return executeTx(s.ctx, s.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(s.ctx, "DELETE FROM heads WHERE topic = $1", topic)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(s.ctx, "DELETE FROM events WHERE topic = $1", topic)
		return err
	})
}

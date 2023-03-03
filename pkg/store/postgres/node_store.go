package postgresstore

import (
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/store/postgres/migrations"
)

type NodeStore struct {
	ctx context.Context
	db  *DB
}

func NewNodeStore(ctx context.Context, db *DB) (*NodeStore, error) {
	s := &NodeStore{
		ctx: context.WithLogger(ctx, ctx.Logger().Named("pgstore")),
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

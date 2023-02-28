package postgresstore

import (
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/store/postgres/migrations"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type NodeStore struct {
	log *zap.Logger
	db  *DB
}

func NewNodeStore(log *zap.Logger, db *DB) (*NodeStore, error) {
	s := &NodeStore{
		log: log.Named("pgstore"),
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
	return New(s.log, s.db, topic), nil
}

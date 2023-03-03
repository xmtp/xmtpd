package bolt

import (
	"bytes"
	"errors"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	bolt "go.etcd.io/bbolt"
)

var (
	// MetaBucket Keys
	VersionKey     = []byte("version")
	CurrentVersion = []byte{0x1, 0x0} // 1.0

	// Top Level Buckets, besides individual topic buckets
	MetaBucket = []byte("META")

	ErrVersionMismatch = errors.New("unexpected DB version")
)

// NodeStore using embedded BoltDB.
type NodeStore struct {
	ctx context.Context
	db  *bolt.DB
}

func NewNodeStore(ctx context.Context, opts *Options) (*NodeStore, error) {
	db, err := bolt.Open(opts.DataPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err = db.Update(func(tx *bolt.Tx) error {
		meta := tx.Bucket(MetaBucket)
		if meta != nil {
			v := meta.Get(VersionKey)
			if bytes.Equal(v, CurrentVersion) {
				return nil
			}
			return ErrVersionMismatch
		}
		if meta, err = tx.CreateBucket(MetaBucket); err != nil {
			return err
		}
		return meta.Put(VersionKey, CurrentVersion)
	}); err != nil {
		return nil, err
	}

	return &NodeStore{db: db, ctx: ctx}, nil
}

// NewTopic returns a store for a pre-existing topic or creates a new one.
func (s *NodeStore) NewTopic(name string) (crdt.Store, error) {
	nameBytes := []byte(name)
	// Make sure the topic structure is in place in the DB
	if err := s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(nameBytes) != nil {
			return nil
		}
		topic, err := tx.CreateBucket(nameBytes)
		if err != nil {
			return err
		}
		if _, err = topic.CreateBucket(EventsBucket); err != nil {
			return err
		}
		if _, err = topic.CreateBucket(ByCIDBucket); err != nil {
			return err
		}
		_, err = topic.CreateBucket(HeadsBucket)
		return err
	}); err != nil {
		return nil, err
	}
	return New(s.ctx, s.db, name), nil
}

// Topics lists all topics in the store.
func (s *NodeStore) Topics() (topics []string, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			// skip the META bucket
			if bytes.Equal(k, MetaBucket) {
				continue
			}
			topics = append(topics, string(k))
		}
		return nil
	})
	return topics, err
}

func (s *NodeStore) Close() error {
	return s.db.Close()
}

package bolt

import (
	"bytes"
	"errors"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/zap"
	bolt "go.etcd.io/bbolt"
)

var (
	// MetaBucket Keys
	VersionKey     = []byte("version")
	CurrentVersion = []byte{0x1, 0x0} // 1.0

	// Top Level Buckets, besides individual topic buckets
	MetaBucket = []byte("META")
	// Topic Level Buckets
	EventsBucket = []byte("EVENTS")
	HeadsBucket  = []byte("HEADS")
	ByTimeBucket = []byte("BY_TIME")

	ErrByTimeCorrupted = errors.New("By time index references unknown event")
	ErrVersionMismatch = errors.New("Unexpected DB version")
)

// Store using embedded BoltDB.
type Store struct {
	db *bolt.DB
}

func NewStore(fn string) (*Store, error) {
	db, err := bolt.Open(fn, 0600, nil)
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

	return &Store{db: db}, nil
}

// NewTopic returns a store for a pre-existing topic or creates a new one.
func (s *Store) NewTopic(name string, n *crdt.Node) (crdt.TopicStore, error) {
	nameBytes := []byte(name)
	// Make sure the topic structure is in place in the DB
	if err := s.db.Batch(func(tx *bolt.Tx) error {
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
		if _, err = topic.CreateBucket(ByTimeBucket); err != nil {
			return err
		}
		_, err = topic.CreateBucket(HeadsBucket)
		return err
	}); err != nil {
		return nil, err
	}
	return &TopicStore{
		name:  nameBytes,
		node:  n,
		Store: s,
		log:   n.LogNamed(name),
	}, nil
}

// Topics lists all topics in the store.
func (s *Store) Topics() (topics []string, err error) {
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

func (s *Store) Close() error {
	return s.db.Close()
}

// TopicStore is a topic focused adaptor on top of Store
type TopicStore struct {
	name []byte
	node *crdt.Node
	*Store
	log *zap.Logger
}

func (s *TopicStore) AddEvent(ev *crdt.Event) (added bool, err error) {
	err = s.db.Batch(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		if err = s.addEvent(topic, ev); err != nil {
			return err
		}
		added = true
		return nil
	})
	if err == nil && added {
		s.log.Debug("added event", zap.Cid("event", ev.Cid))
	}
	return added, err
}

func (s *TopicStore) AddHead(ev *crdt.Event) (added bool, err error) {
	err = s.db.Batch(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		if err = s.addEvent(topic, ev); err != nil {
			return err
		}
		heads := topic.Bucket(HeadsBucket)
		if err = heads.Put(ev.Cid, nil); err != nil {
			return err
		}
		added = true
		return nil
	})
	if err == nil && added {
		s.log.Debug("added head", zap.Cid("event", ev.Cid))
	}
	return added, err
}

func (s *TopicStore) RemoveHead(cid mh.Multihash) (have bool, err error) {
	var isHead bool
	err = s.db.Batch(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		if events.Get(cid) == nil {
			return nil
		}
		have = true
		heads := topic.Bucket(HeadsBucket)
		isHead = heads.Get(cid) != nil
		if isHead {
			if err = heads.Delete(cid); err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil && isHead {
		s.log.Debug("removing head", zap.Cid("event", cid))
	}
	return have, err
}

func (s *TopicStore) NewEvent(env *messagev1.Envelope) (ev *crdt.Event, err error) {
	err = s.db.Batch(func(tx *bolt.Tx) error {
		var allHeads []mh.Multihash
		topic := tx.Bucket(s.name)
		heads := topic.Bucket(HeadsBucket)
		c := heads.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cid := make(mh.Multihash, len(k))
			copy(cid, k)
			allHeads = append(allHeads, cid)
			if err = c.Delete(); err != nil {
				return err
			}
		}
		if ev, err = crdt.NewEvent(env, allHeads); err != nil {
			return err
		}
		if err = s.addEvent(topic, ev); err != nil {
			return err
		}
		return heads.Put(ev.Cid, nil)
	})
	if err == nil {
		s.log.Debug("creating event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))
	}
	return ev, err
}

func (s *TopicStore) FindMissingLinks() (links []mh.Multihash, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		c := events.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var cids []mh.Multihash
			if cids, err = crdt.LinksFromBytes(v); err != nil {
				return err
			}
			for _, cid := range cids {
				if events.Get(cid) == nil {
					link := make(mh.Multihash, len(cid))
					copy(link, cid)
					links = append(links, link)
				}
			}
		}
		return nil
	})
	return links, err
}

func (s *TopicStore) Get(cid mh.Multihash) (ev *crdt.Event, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		evBytes := events.Get(cid)
		if evBytes != nil {
			ev, err = crdt.EventFromBytes(cid, evBytes)
		}
		return err
	})
	return ev, err
}

func (s *TopicStore) Count() (count int, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		c := events.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})
	return count, err
}

// private functions

func (s *TopicStore) addEvent(topic *bolt.Bucket, ev *crdt.Event) error {
	events := topic.Bucket(EventsBucket)
	if events.Get(ev.Cid) != nil {
		return nil
	}
	evBytes, err := crdt.EventToBytes(ev)
	if err != nil {
		return err
	}
	if err = events.Put(ev.Cid, evBytes); err != nil {
		return err
	}
	byTime := topic.Bucket(ByTimeBucket)
	return byTime.Put(crdt.ToByTimeKey(ev.TimestampNs, ev.Cid), ev.Cid)
}

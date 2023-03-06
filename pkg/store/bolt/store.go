package bolt

import (
	"errors"
	"fmt"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
	bolt "go.etcd.io/bbolt"
)

var (
	// Topic Level Buckets
	EventsBucket = []byte("EVENTS")
	HeadsBucket  = []byte("HEADS")
	ByCIDBucket  = []byte("BY_CID")

	ErrByCIDCorrupted = errors.New("by CID index references unknown event time key")
)

// Store is a topic focused adaptor on top of Store
type Store struct {
	name []byte
	db   *bolt.DB
	log  *zap.Logger
}

func New(ctx context.Context, db *bolt.DB, topic string) *Store {
	return &Store{
		name: []byte(topic),
		db:   db,
		log:  ctx.Logger().With(zap.String("topic", topic)),
	}
}

func (s *Store) InsertEvent(ctx context.Context, ev *types.Event) (added bool, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		added, err = addEvent(topic, ev)
		return err
	})
	if err == nil && added {
		s.log.Debug("added event", zap.Cid("event", ev.Cid))
	}
	return added, err
}

func (s *Store) InsertHead(ctx context.Context, ev *types.Event) (added bool, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		if added, err = addEvent(topic, ev); err != nil {
			return err
		}
		heads := topic.Bucket(HeadsBucket)
		return heads.Put(ev.Cid, nil)
	})
	if err == nil && added {
		s.log.Debug("added head", zap.Cid("event", ev.Cid))
	}
	return added, err
}

func (s *Store) RemoveHead(ctx context.Context, cid multihash.Multihash) (have bool, err error) {
	var isHead bool
	err = s.db.Update(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		byCID := topic.Bucket(ByCIDBucket)
		if byCID.Get(cid) == nil {
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

func (s *Store) AppendEvent(ctx context.Context, env *messagev1.Envelope) (ev *types.Event, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		var allHeads []multihash.Multihash
		topic := tx.Bucket(s.name)
		heads := topic.Bucket(HeadsBucket)
		c := heads.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cid := make(multihash.Multihash, len(k))
			copy(cid, k)
			allHeads = append(allHeads, cid)
			if err = c.Delete(); err != nil {
				return err
			}
		}
		if ev, err = types.NewEvent(env, allHeads); err != nil {
			return err
		}
		if _, err = addEvent(topic, ev); err != nil {
			return err
		}
		return heads.Put(ev.Cid, nil)
	})
	if err == nil {
		s.log.Debug("creating event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))
	}
	return ev, err
}

func (s *Store) FindMissingLinks(ctx context.Context) (links []multihash.Multihash, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		byCID := topic.Bucket(ByCIDBucket)
		c := events.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var cids []multihash.Multihash
			if cids, err = types.LinksFromBytes(v); err != nil {
				return err
			}
			for _, cid := range cids {
				if byCID.Get(cid) == nil {
					link := make(multihash.Multihash, len(cid))
					copy(link, cid)
					links = append(links, link)
				}
			}
		}
		return nil
	})
	return links, err
}

func (s *Store) GetEvents(ctx context.Context, cids ...multihash.Multihash) (evs []*types.Event, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		byCID := topic.Bucket(ByCIDBucket)
		for _, cid := range cids {
			orderKey := byCID.Get(cid)
			evBytes := events.Get(orderKey)
			if evBytes != nil {
				ev, err := types.EventFromBytes(evBytes)
				if err != nil {
					return err
				}
				evs = append(evs, ev)
			}
		}
		return nil
	})
	return evs, err
}

func (s *Store) Heads(ctx context.Context) (cids []multihash.Multihash, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(HeadsBucket)
		c := events.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cids = append(cids, k)
		}
		return nil
	})
	return cids, err
}

func (s *Store) Events(ctx context.Context) (evs []*types.Event, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		events := topic.Bucket(EventsBucket)
		c := events.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ev, err := types.EventFromBytes(v)
			if err != nil {
				return err
			}
			evs = append(evs, ev)
		}
		return nil
	})
	return evs, err
}

func (s *Store) NewCursor(ev *types.Event) *messagev1.Cursor {
	return &messagev1.Cursor{
		Cursor: &messagev1.Cursor_Index{
			Index: &messagev1.IndexCursor{
				SenderTimeNs: ev.TimestampNs,
				Digest:       ev.Cid,
			},
		},
	}
}

func (s *Store) InsertNewEvents(ctx context.Context, evs []*types.Event) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		topic := tx.Bucket([]byte(s.name))
		for _, ev := range evs {
			added, err := addEvent(topic, ev)
			if err != nil {
				return err
			}
			if !added {
				return fmt.Errorf("event already exists: %s", ev.Cid)
			}
		}
		return nil
	})
}

// private functions

func addEvent(topic *bolt.Bucket, ev *types.Event) (added bool, err error) {
	byCid := topic.Bucket(ByCIDBucket)
	if byCid.Get(ev.Cid) != nil {
		return false, nil
	}
	evBytes, err := ev.ToBytes()
	if err != nil {
		return false, err
	}
	orderKey := types.ToByTimeKey(ev.TimestampNs, ev.Cid)
	events := topic.Bucket(EventsBucket)
	if err = events.Put(orderKey, evBytes); err != nil {
		return false, err
	}
	return true, byCid.Put(ev.Cid, orderKey)
}

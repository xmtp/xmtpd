package postgresstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	ErrTODO          = errors.New("TODO")
	ErrTopicMismatch = errors.New("topic mismatch")
	ErrMissingTopic  = errors.New("missing topic")
	ErrTooManyTopics = errors.New("too many topics")
)

type Store struct {
	log   *zap.Logger
	db    *DB
	topic string
}

func New(ctx context.Context, db *DB, topic string) *Store {
	return &Store{
		log:   ctx.Logger().With(zap.String("topic", topic)),
		db:    db,
		topic: topic,
	}
}

func (s *Store) InsertEvent(ctx context.Context, ev *types.Event) (bool, error) {
	s.log.Debug("inserting event", zap.Cid("event", ev.Cid))

	var added bool
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		var err error
		added, err = s.insertEvent(ctx, tx, ev)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return added, nil
}

func (s *Store) AppendEvent(ctx context.Context, env *messagev1.Envelope) (*types.Event, error) {
	if env.ContentTopic != s.topic {
		return nil, ErrTopicMismatch
	}

	var ev *types.Event
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		heads, err := s.heads(ctx, tx)
		if err != nil {
			return err
		}
		ev, err = types.NewEvent(env, heads)
		if err != nil {
			return err
		}
		s.log.Debug("appending event", zap.Cid("event", ev.Cid), zap.Int("links", len(ev.Links)))

		eventAdded, err := s.insertEvent(ctx, tx, ev)
		if err != nil {
			return err
		}
		if !eventAdded {
			return nil
		}

		_, err = s.insertHead(ctx, tx, ev)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *Store) InsertHead(ctx context.Context, ev *types.Event) (bool, error) {
	s.log.Debug("inserting head", zap.Cid("event", ev.Cid))

	var headAdded bool
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		eventAdded, err := s.insertEvent(ctx, tx, ev)
		if err != nil {
			return err
		}
		if !eventAdded {
			return nil
		}

		headAdded, err = s.insertHead(ctx, tx, ev)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return headAdded, nil
}

func (s *Store) RemoveHead(ctx context.Context, cid multihash.Multihash) (bool, error) {
	var deleted bool
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, "DELETE FROM heads WHERE topic = $1 AND cid = $2", s.topic, cid.HexString())
		if err != nil {
			return err
		}
		count, err := res.RowsAffected()
		if err != nil {
			return err
		}
		deleted = count == 1
		return nil
	})
	return deleted, err
}

func (s *Store) FindMissingLinks(ctx context.Context) ([]multihash.Multihash, error) {
	// TODO: this is a very inefficient way of finding missing links. If it's
	// a regular operation then we should track then separately in a
	// missing_links table at insert-time, or something more efficient than
	// pulling down the whole set of events.
	events, err := s.Events(ctx)
	if err != nil {
		return nil, err
	}
	cidSet := map[string]multihash.Multihash{}
	for _, ev := range events {
		for _, cid := range ev.Links {
			cidSet[cid.HexString()] = cid
		}
	}
	cids := make([]multihash.Multihash, 0, len(cidSet))
	for _, cid := range cidSet {
		cids = append(cids, cid)
	}
	return cids, nil
}

func (s *Store) GetEvents(ctx context.Context, cids ...multihash.Multihash) ([]*types.Event, error) {
	var events []*types.Event
	ids := make([]string, 0, len(cids))
	for _, cid := range cids {
		ids = append(ids, cid.HexString())
	}
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, "SELECT cid, links, topic, timestamp_ns, message FROM events WHERE topic = $1 AND cid = ANY($2)", s.topic, pq.StringArray(ids))
		if err != nil {
			return err
		}
		defer rows.Close()
		events, err = eventsFromRows(rows)
		return err
	})
	return events, err
}

func (s *Store) Events(ctx context.Context) ([]*types.Event, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT cid, links, topic, timestamp_ns, message FROM events WHERE topic = $1 ORDER BY topic, timestamp_ns, cid", s.topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return eventsFromRows(rows)
}

func (s *Store) Heads(ctx context.Context) ([]multihash.Multihash, error) {
	var cids []multihash.Multihash
	err := s.executeTx(ctx, func(tx *sql.Tx) error {
		var err error
		cids, err = s.heads(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cids, nil
}

func (s *Store) InsertNewEvents(ctx context.Context, evs []*types.Event) error {
	return s.executeTx(ctx, func(tx *sql.Tx) error {
		for _, ev := range evs {
			added, err := s.insertEvent(ctx, tx, ev)
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

func (s *Store) insertEvent(ctx context.Context, tx *sql.Tx, ev *types.Event) (bool, error) {
	if ev.ContentTopic != s.topic {
		return false, ErrTopicMismatch
	}

	linksJSON, err := json.Marshal(ev.Links)
	if err != nil {
		return false, err
	}

	message := ev.Message
	if message == nil {
		message = []byte{}
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO events (cid, links, topic, timestamp_ns, message) VALUES ($1, $2, $3, $4, $5)`, ev.Cid.HexString(), string(linksJSON), ev.ContentTopic, ev.TimestampNs, message)
	if err != nil {
		if isDuplicateKeyError(err) {
			s.log.Debug("ignoring duplicate key error", zap.Error(err), zap.Cid("event", ev.Cid))
			return false, nil
		}
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (s *Store) insertHead(ctx context.Context, tx *sql.Tx, ev *types.Event) (bool, error) {
	if ev.ContentTopic != s.topic {
		return false, ErrTopicMismatch
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO heads (topic, cid) VALUES ($1, $2)`, ev.ContentTopic, ev.Cid.HexString())
	if err != nil {
		if isDuplicateKeyError(err) {
			s.log.Debug("ignoring duplicate key error", zap.Error(err), zap.Cid("event", ev.Cid))
			return false, nil
		}
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (s *Store) heads(ctx context.Context, tx *sql.Tx) ([]multihash.Multihash, error) {
	rows, err := tx.QueryContext(ctx, "SELECT cid FROM heads WHERE topic = $1", s.topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cids []multihash.Multihash
	for rows.Next() {
		var cidHex string
		err := rows.Scan(&cidHex)
		if err != nil {
			return nil, err
		}
		cid, err := multihash.FromHexString(cidHex)
		if err != nil {
			return nil, err
		}
		cids = append(cids, cid)
	}

	return cids, nil
}

func (s *Store) executeTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return executeTx(ctx, s.db, fn)
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.Code == pgerrcode.UniqueViolation
}

func eventsFromRows(rows *sql.Rows) (events []*types.Event, err error) {
	for rows.Next() {
		var (
			cidHex      string
			linksJSON   string
			topic       string
			timestampNS uint64
			message     []byte
		)
		err := rows.Scan(&cidHex, &linksJSON, &topic, &timestampNS, &message)
		if err != nil {
			return nil, err
		}
		cid, err := multihash.FromHexString(cidHex)
		if err != nil {
			return nil, err
		}
		var links []multihash.Multihash
		err = json.Unmarshal([]byte(linksJSON), &links)
		if err != nil {
			return nil, err
		}
		events = append(events, &types.Event{
			Cid:   cid,
			Links: links,
			Envelope: &messagev1.Envelope{
				ContentTopic: topic,
				TimestampNs:  timestampNS,
				Message:      message,
			},
		})
	}

	return events, nil
}

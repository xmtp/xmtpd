package postgresstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

		_, err = s.insertEvent(ctx, tx, ev)
		if err != nil {
			return err
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
		_, err := s.insertEvent(ctx, tx, ev)
		if err != nil {
			return err
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
	res, err := s.db.ExecContext(ctx, "DELETE FROM heads WHERE topic = $1 AND cid = $2", s.topic, cid)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
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
	return nil, ErrTODO
}

func (s *Store) NewCursor(ev *types.Event) *messagev1.Cursor {
	return &messagev1.Cursor{
		Cursor: &messagev1.Cursor_Index{
			Index: &messagev1.IndexCursor{
				SenderTimeNs: ev.TimestampNs,
				Digest:       []byte(ev.Cid.HexString()),
			},
		},
	}
}

func (s *Store) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	if len(req.ContentTopics) == 0 {
		req.ContentTopics = []string{s.topic}
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}
	topic := req.ContentTopics[0]

	if topic != s.topic {
		return nil, ErrTopicMismatch
	}

	baseSQL := "SELECT cid, topic, timestamp_ns, message FROM events WHERE topic = $1"
	args := []any{topic}

	timeFilterSQL := ""
	if req.StartTimeNs > 0 {
		timeFilterSQL = fmt.Sprintf(" AND timestamp_ns >= $%d", len(args)+1)
		args = append(args, req.StartTimeNs)
	}
	if req.EndTimeNs > 0 {
		timeFilterSQL += fmt.Sprintf(" AND timestamp_ns <= $%d", len(args)+1)
		args = append(args, req.EndTimeNs)
	}

	sortSQL := " ORDER BY topic, timestamp_ns, cid"
	limitSQL := ""
	cursorFilterSQL := ""
	if req.PagingInfo != nil {
		if req.PagingInfo.Cursor != nil && req.PagingInfo.Cursor.GetIndex() != nil {
			cursor := req.PagingInfo.Cursor.GetIndex()
			if req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING {
				cursorFilterSQL = fmt.Sprintf(" AND timestamp_ns < $%d OR (timestamp_ns = $%d AND cid < $%d)", len(args)+1, len(args)+2, len(args)+3)
				args = append(args, cursor.SenderTimeNs, cursor.SenderTimeNs, string(cursor.Digest))
			} else {
				cursorFilterSQL = fmt.Sprintf(" AND timestamp_ns > $%d OR (timestamp_ns = $%d AND cid > $%d)", len(args)+1, len(args)+2, len(args)+3)
				args = append(args, cursor.SenderTimeNs, cursor.SenderTimeNs, string(cursor.Digest))
			}
		}

		if req.PagingInfo.Direction == messagev1.SortDirection_SORT_DIRECTION_DESCENDING {
			sortSQL = " ORDER BY topic, timestamp_ns DESC, cid"
		}

		if req.PagingInfo.Limit > 0 {
			limitSQL = fmt.Sprintf(" LIMIT $%d", len(args)+1)
			args = append(args, req.PagingInfo.Limit)
		}
	}

	sql := baseSQL + timeFilterSQL + cursorFilterSQL + sortSQL + limitSQL

	s.log.Debug("querying", zap.String("sql", sql))
	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []*messagev1.Envelope
	for rows.Next() {
		var (
			cidHex      string
			topic       string
			timestampNS uint64
			message     []byte
		)
		err := rows.Scan(&cidHex, &topic, &timestampNS, &message)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &messagev1.Envelope{
			ContentTopic: topic,
			TimestampNs:  timestampNS,
			Message:      message,
		})
	}

	return &messagev1.QueryResponse{
		Envelopes: envs,
	}, nil
}

func (s *Store) Events(ctx context.Context) ([]*types.Event, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT cid, links, topic, timestamp_ns, message FROM events WHERE topic = $1 ORDER BY topic, timestamp_ns, cid", s.topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*types.Event
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				s.log.Error("error rolling back", zap.Error(err))
			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			rollbackErr := tx.Rollback() // err is non-nil; don't change it
			if rollbackErr != nil {
				s.log.Error("error rolling back", zap.Error(err))
			}
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = fn(tx)
	return err
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.Code == pgerrcode.UniqueViolation
}

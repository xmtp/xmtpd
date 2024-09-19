package api

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	subscriptionBufferSize    = 1024
	maxSubscriptionsPerClient = 10000
	SubscribeWorkerPollTime   = 100 * time.Millisecond
	subscribeWorkerPollRows   = 10000
	maxTopicLength            = 128
)

type listener struct {
	ctx         context.Context
	ch          chan<- []*message_api.OriginatorEnvelope
	closed      bool
	topics      map[string]struct{}
	originators map[uint32]struct{}
	isGlobal    bool
}

type listenerSet struct {
	sync.Map // map[*listener]struct{}
}

func (ls *listenerSet) addListener(l *listener) {
	ls.Store(l, struct{}{})
}

func (ls *listenerSet) removeListener(l *listener) {
	ls.Delete(l)
}

func (ls *listenerSet) isEmpty() bool {
	empty := true
	ls.Range(func(_, _ interface{}) bool {
		empty = false
		return false // stop iteration
	})
	return empty
}

// Maps from a key to a set of listeners
type listenersMap[K comparable] struct {
	data sync.Map     // map[K]*listenerSet
	mu   sync.RWMutex // ensures mutations are consistent
}

func (lm *listenersMap[K]) addListener(keys map[K]struct{}, l *listener) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	for key := range keys {
		value, _ := lm.data.LoadOrStore(key, &listenerSet{})
		set := value.(*listenerSet)
		set.addListener(l)
	}
}

func (lm *listenersMap[K]) removeListener(keys map[K]struct{}, l *listener) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for key := range keys {
		value, ok := lm.data.Load(key)
		if !ok || value == nil {
			return // Key doesn't exist, nothing to do
		}
		set := value.(*listenerSet)
		set.removeListener(l)
		if set.isEmpty() {
			lm.data.Delete(key)
		}
	}
}

func (lm *listenersMap[K]) getListeners(key K) *listenerSet {
	// No lock needed, because we are not mutating lm.data
	if value, ok := lm.data.Load(key); ok {
		return value.(*listenerSet)
	}
	return nil
}

// A worker that listens for new envelopes in the DB and sends them to subscribers
// Assumes that there are many listeners - non-blocking updates are sent on buffered channels
// and may be dropped if full
type subscribeWorker struct {
	ctx context.Context
	log *zap.Logger

	dbSubscription <-chan []queries.GatewayEnvelope
	// Assumption: listeners cannot be in multiple slices
	globalListeners     listenerSet
	originatorListeners listenersMap[uint32]
	topicListeners      listenersMap[string]
}

func startSubscribeWorker(
	ctx context.Context,
	log *zap.Logger,
	store *sql.DB,
) (*subscribeWorker, error) {
	log = log.With(zap.String("method", "subscribeWorker"))
	q := queries.New(store)
	pollableQuery := func(ctx context.Context, lastSeen db.VectorClock, numRows int32) ([]queries.GatewayEnvelope, db.VectorClock, error) {
		envs, err := q.
			SelectGatewayEnvelopes(
				ctx,
				*db.SetVectorClock(&queries.SelectGatewayEnvelopesParams{}, lastSeen),
			)
		if err != nil {
			return nil, lastSeen, err
		}
		for _, env := range envs {
			// TODO(rich) Handle out-of-order envelopes
			lastSeen[uint32(env.OriginatorNodeID)] = uint64(env.OriginatorSequenceID)
		}
		return envs, lastSeen, nil
	}

	vc, err := q.SelectVectorClock(ctx)
	if err != nil {
		return nil, err
	}

	subscription := db.NewDBSubscription(
		ctx,
		log,
		pollableQuery,
		db.ToVectorClock(vc),
		db.PollingOptions{
			Interval: SubscribeWorkerPollTime,
			NumRows:  subscribeWorkerPollRows,
		},
	)
	dbChan, err := subscription.Start()
	if err != nil {
		return nil, err
	}
	worker := &subscribeWorker{
		ctx:                 ctx,
		log:                 log,
		dbSubscription:      dbChan,
		globalListeners:     listenerSet{},
		originatorListeners: listenersMap[uint32]{},
		topicListeners:      listenersMap[string]{},
	}

	go worker.start()

	return worker, nil
}

func (s *subscribeWorker) start() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case new_batch := <-s.dbSubscription:
			for _, row := range new_batch {
				s.dispatch(&row)
			}
		}
	}
}

func (s *subscribeWorker) dispatch(row *queries.GatewayEnvelope) {
	env := &message_api.OriginatorEnvelope{}
	err := proto.Unmarshal(row.OriginatorEnvelope, env)
	if err != nil {
		s.log.Error("Failed to unmarshal envelope", zap.Error(err))
		return
	}

	originatorListeners := s.originatorListeners.getListeners(uint32(row.OriginatorNodeID))
	topicListeners := s.topicListeners.getListeners(hex.EncodeToString(row.Topic))
	s.dispatchToListeners(originatorListeners, env)
	s.dispatchToListeners(topicListeners, env)
	s.dispatchToListeners(&s.globalListeners, env)
}

func (s *subscribeWorker) dispatchToListeners(
	listeners *listenerSet,
	env *message_api.OriginatorEnvelope,
) {
	if listeners == nil {
		return
	}
	listeners.Range(func(key, _ any) bool {
		l := key.(*listener)
		if l.closed {
			return true
		}
		// Assumption: listener channel is never closed by a different goroutine
		select {
		case <-l.ctx.Done():
			s.log.Debug("Stream closed, removing listener", zap.Any("listener", l.ch))
			s.closeListener(l)
		default:
			select {
			case l.ch <- []*message_api.OriginatorEnvelope{env}:
			default:
				s.log.Info("Channel full, removing listener", zap.Any("listener", l.ch))
				s.closeListener(l)
			}
		}
		return true
	})
}

func (s *subscribeWorker) closeListener(l *listener) {
	// Assumption: this method may not be called across multiple goroutines
	l.closed = true
	close(l.ch)

	go func() {
		if l.isGlobal {
			s.globalListeners.Delete(l)
		}
		s.topicListeners.removeListener(l.topics, l)
		s.originatorListeners.removeListener(l.originators, l)
	}()
}

func (s *subscribeWorker) listen(
	ctx context.Context,
	requests []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest,
) (<-chan []*message_api.OriginatorEnvelope, error) {
	ch := make(chan []*message_api.OriginatorEnvelope, subscriptionBufferSize)
	l := &listener{
		ctx:         ctx,
		ch:          ch,
		topics:      make(map[string]struct{}),
		originators: make(map[uint32]struct{}),
		isGlobal:    false,
	}

	if len(requests) > maxSubscriptionsPerClient {
		return nil, fmt.Errorf(
			"too many subscriptions: %d, consider subscribing to fewer topics or subscribing without a filter",
			len(requests),
		)
	}
	for _, req := range requests {
		enum := req.GetQuery().GetFilter()
		if enum == nil {
			l.isGlobal = true
		}
		switch filter := enum.(type) {
		case *message_api.EnvelopesQuery_Topic:
			topic := hex.EncodeToString(filter.Topic)
			if len(filter.Topic) == 0 || len(filter.Topic) > maxTopicLength {
				return nil, fmt.Errorf("invalid topic: %s", topic)
			}
			if _, exists := l.topics[topic]; exists {
				return nil, fmt.Errorf("multiple requests for same topic: %s", topic)
			}
			l.topics[topic] = struct{}{}
		case *message_api.EnvelopesQuery_OriginatorNodeId:
			if _, exists := l.originators[filter.OriginatorNodeId]; exists {
				return nil, fmt.Errorf("multiple requests for same originator: %d", filter.OriginatorNodeId)
			}
			l.originators[filter.OriginatorNodeId] = struct{}{}
		default:
			l.isGlobal = true
		}
	}

	if l.isGlobal {
		if len(l.topics) > 0 || len(l.originators) > 0 {
			return nil, fmt.Errorf(
				"cannot filter by topic or originator when subscribing to all",
			)
		}
		s.globalListeners.Store(l, struct{}{})
	} else if len(l.topics) > 0 {
		if len(l.originators) > 0 {
			return nil, fmt.Errorf("cannot filter by both topic and originator in same subscription request")
		}
		s.topicListeners.addListener(l.topics, l)
	} else if len(l.originators) > 0 {
		s.originatorListeners.addListener(l.originators, l)
	}

	return ch, nil
}

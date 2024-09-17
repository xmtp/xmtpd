package api

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ch          chan<- []*message_api.OriginatorEnvelope
	topics      map[string]struct{}
	originators map[uint32]struct{}
	isGlobal    bool
	closed      atomic.Bool
}

// Maps from a key to a set of listeners
type listenerMap[K comparable] struct {
	sync.Map // map[K]*sync.Map, where inner sync.Map is map[*listener]struct{}
}

func (lm *listenerMap[K]) addListener(key K, l *listener) {
	value, _ := lm.LoadOrStore(key, &sync.Map{})
	innerMap := value.(*sync.Map)
	innerMap.Store(l, struct{}{})
}

func (lm *listenerMap[K]) removeListener(key K, l *listener) {
	for {
		value, ok := lm.Load(key)
		if !ok || value == nil {
			return // Key doesn't exist, nothing to do
		}
		innerMap := value.(*sync.Map)
		innerMap.Delete(l)

		if !isEmptySyncMap(innerMap) || lm.CompareAndDelete(key, value) {
			return
		}
	}
}

// isEmptySyncMap checks if a sync.Map is empty
func isEmptySyncMap(m *sync.Map) bool {
	empty := true
	m.Range(func(_, _ interface{}) bool {
		empty = false
		return false // stop iteration
	})
	return empty
}

// A worker that listens for new envelopes in the DB and sends them to subscribers
// Assumes that there are many listeners - non-blocking updates are sent on buffered channels
// and may be dropped if full
type subscribeWorker struct {
	ctx context.Context
	log *zap.Logger

	dbSubscription <-chan []queries.GatewayEnvelope
	// Assumption: listeners cannot be in multiple slices
	topicListeners      listenerMap[string]
	originatorListeners listenerMap[uint32]
	globalListeners     sync.Map // map[*listener]struct{}
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
		globalListeners:     sync.Map{},
		originatorListeners: listenerMap[uint32]{},
		topicListeners:      listenerMap[string]{},
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

	originatorID := uint32(row.OriginatorNodeID)
	if listenersMap, ok := s.originatorListeners.Load(originatorID); ok {
		s.dispatchToListeners(listenersMap.(*sync.Map), env)
	}

	topic := hex.EncodeToString(row.Topic)
	if listenersMap, ok := s.topicListeners.Load(topic); ok {
		s.dispatchToListeners(listenersMap.(*sync.Map), env)
	}

	s.dispatchToListeners(&s.globalListeners, env)
}

func (s *subscribeWorker) dispatchToListeners(
	listeners *sync.Map,
	env *message_api.OriginatorEnvelope,
) {
	listeners.Range(func(key, _ any) bool {
		l := key.(*listener)
		select {
		case l.ch <- []*message_api.OriginatorEnvelope{env}:
			// Successfully sent
		default:
			// Channel is full or closed
			s.log.Info("Channel full or closed, removing listener", zap.Any("listener", l.ch))
			go s.removeListener(l)
		}
		return true
	})
}

func (s *subscribeWorker) listen(
	requests []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest,
) (<-chan []*message_api.OriginatorEnvelope, func(), error) {
	ch := make(chan []*message_api.OriginatorEnvelope, subscriptionBufferSize)
	l := &listener{
		ch:          ch,
		topics:      make(map[string]struct{}),
		originators: make(map[uint32]struct{}),
		isGlobal:    false,
	}

	if len(requests) > maxSubscriptionsPerClient {
		return nil, nil, fmt.Errorf(
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
			if len(filter.Topic) == 0 || len(filter.Topic) > maxTopicLength {
				return nil, nil, status.Errorf(codes.InvalidArgument, "invalid topic")
			}
			l.topics[hex.EncodeToString(filter.Topic)] = struct{}{}
		case *message_api.EnvelopesQuery_OriginatorNodeId:
			l.originators[filter.OriginatorNodeId] = struct{}{}
		default:
			l.isGlobal = true
		}
	}

	if l.isGlobal {
		if len(l.topics) > 0 || len(l.originators) > 0 {
			return nil, nil, fmt.Errorf(
				"cannot filter by topic or originator when subscribing to all",
			)
		}
		s.globalListeners.Store(l, struct{}{})
	} else if len(l.topics) > 0 {
		if len(l.originators) > 0 {
			return nil, nil, fmt.Errorf("cannot filter by both topic and originator in same subscription request")
		}
		for topic := range l.topics {
			s.topicListeners.addListener(topic, l)
		}
	} else if len(l.originators) > 0 {
		for originator := range l.originators {
			s.originatorListeners.addListener(originator, l)
		}
	}

	return ch, func() { s.removeListener(l) }, nil
}

func (s *subscribeWorker) removeListener(l *listener) {
	if l.closed.CompareAndSwap(false, true) {
		close(l.ch)

		if l.isGlobal {
			s.globalListeners.Delete(l)
		}

		for topic := range l.topics {
			s.topicListeners.removeListener(topic, l)
		}

		for origin := range l.originators {
			s.originatorListeners.removeListener(origin, l)
		}
	}
}

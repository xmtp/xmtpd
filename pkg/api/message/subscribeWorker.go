package message

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
)

const (
	subscriptionBufferSize  = 1024
	SubscribeWorkerPollTime = 100 * time.Millisecond
	subscribeWorkerPollRows = 10000
)

type listener struct {
	ctx         context.Context
	ch          chan<- []*envelopes.OriginatorEnvelope
	closed      bool
	topics      map[string]struct{}
	originators map[uint32]struct{}
	isGlobal    bool
}

func newListener(
	ctx context.Context,
	logger *zap.Logger,
	query *message_api.EnvelopesQuery,
	ch chan<- []*envelopes.OriginatorEnvelope,
) *listener {
	l := &listener{
		ctx:         ctx,
		ch:          ch,
		topics:      make(map[string]struct{}),
		originators: make(map[uint32]struct{}),
		isGlobal:    false,
	}
	topics := query.GetTopics()
	originators := query.GetOriginatorNodeIds()

	if len(topics) == 0 && len(originators) == 0 {
		l.isGlobal = true
		return l
	}

	for _, t := range topics {
		validatedTopic, err := topic.ParseTopic(t)
		if err != nil {
			logger.Warn("Skipping invalid topic", zap.Binary("topicBytes", t))
			continue
		}
		logger.Debug("Adding topic listener", zap.String("topic", validatedTopic.String()))
		l.topics[validatedTopic.String()] = struct{}{}
	}

	for _, originator := range originators {
		l.originators[originator] = struct{}{}
	}

	return l
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

func (lm *listenersMap[K]) rangeKeys(fn func(key K, listeners *listenerSet) bool) {
	lm.data.Range(func(key, value any) bool {
		return fn(key.(K), value.(*listenerSet))
	})
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
	log = log.Named("subscribeWorker")
	log.Info("Starting subscribe worker")
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
			s.log.Debug("Received new batch", zap.Int("numEnvelopes", len(new_batch)))
			envs := make([]*envelopes.OriginatorEnvelope, 0, len(new_batch))
			for _, row := range new_batch {
				env, err := envelopes.NewOriginatorEnvelopeFromBytes(row.OriginatorEnvelope)
				if err != nil {
					s.log.Error("Failed to unmarshal envelope", zap.Error(err))
					continue
				}
				envs = append(envs, env)
			}
			s.dispatchToOriginators(envs)
			s.dispatchToTopics(envs)
			s.dispatchToGlobals(envs)
		}
	}
}

func (s *subscribeWorker) dispatchToOriginators(envs []*envelopes.OriginatorEnvelope) {
	// We use nested loops here because the number of originators is expected to be small
	// Possible future optimization: Set up set up multiple DB subscriptions instead of one,
	// and have the DB group by originator, topic, and global.
	s.originatorListeners.rangeKeys(func(originator uint32, listeners *listenerSet) bool {
		filteredEnvs := make([]*envelopes.OriginatorEnvelope, 0, len(envs))
		for _, env := range envs {
			if env.OriginatorNodeID() == originator {
				filteredEnvs = append(filteredEnvs, env)
			}
		}
		s.dispatchToListeners(listeners, filteredEnvs)
		return true
	})
}

func (s *subscribeWorker) dispatchToTopics(envs []*envelopes.OriginatorEnvelope) {
	// We iterate envelopes one-by-one, because we expect the number of envelopers
	// per-topic to be small in each tick
	for _, env := range envs {
		listeners := s.topicListeners.getListeners(env.TargetTopic().String())
		s.dispatchToListeners(listeners, []*envelopes.OriginatorEnvelope{env})
	}
}

func (s *subscribeWorker) dispatchToGlobals(envs []*envelopes.OriginatorEnvelope) {
	s.dispatchToListeners(&s.globalListeners, envs)
}

func (s *subscribeWorker) dispatchToListeners(
	listeners *listenerSet,
	envs []*envelopes.OriginatorEnvelope,
) {
	if listeners == nil || len(envs) == 0 {
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
			case l.ch <- envs:
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
		} else if len(l.topics) > 0 {
			s.topicListeners.removeListener(l.topics, l)
		} else if len(l.originators) > 0 {
			s.originatorListeners.removeListener(l.originators, l)
		}
	}()
}

func (s *subscribeWorker) listen(
	ctx context.Context,
	query *message_api.EnvelopesQuery,
) <-chan []*envelopes.OriginatorEnvelope {
	ch := make(chan []*envelopes.OriginatorEnvelope, subscriptionBufferSize)
	l := newListener(ctx, s.log, query, ch)

	if l.isGlobal {
		s.globalListeners.Store(l, struct{}{})
	} else if len(l.topics) > 0 {
		s.topicListeners.addListener(l.topics, l)
	} else if len(l.originators) > 0 {
		s.originatorListeners.addListener(l.originators, l)
	}

	return ch
}

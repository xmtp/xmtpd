package message

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	subscriptionBufferSize  = 1024
	SubscribeWorkerPollTime = 100 * time.Millisecond
	// based on measurements in testnet using PG, we can poll at most 1000 elements in a large DB
	// this gives us sufficient throughput if being run continually
	subscribeWorkerPollRows = 1000
)

type listener struct {
	ctx         context.Context
	ch          chan<- []*envelopes.OriginatorEnvelope
	closed      bool
	topics      map[string]struct{}
	originators map[uint32]struct{}
	isEmpty     bool
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
		isEmpty:     false,
	}
	topics := query.GetTopics()
	originators := query.GetOriginatorNodeIds()

	if len(topics) == 0 && len(originators) == 0 {
		l.isEmpty = true
		return l
	}

	for _, t := range topics {
		validatedTopic, err := topic.ParseTopic(t)
		if err != nil {
			logger.Warn("skipping invalid topic", zap.Binary("topic_bytes", t))
			continue
		}
		logger.Debug("adding topic listener", zap.String("topic", validatedTopic.String()))
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

	ls.Range(func(_, _ any) bool {
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
	ctx    context.Context
	logger *zap.Logger

	// TODO: Check - queries.GatewayEnvelopesView and queries.SelectGatewayEnvelopesByOriginatorsRow
	// models are identical and they're overly verbose.
	// dbSubscription  <-chan []queries.GatewayEnvelopesView
	dbSubscriptions map[uint32]<-chan []queries.SelectGatewayEnvelopesByOriginatorsRow
	// Assumption: listeners cannot be in multiple slices
	emptyListeners      listenerSet
	originatorListeners listenersMap[uint32]
	topicListeners      listenersMap[string]
}

func shortEnvelopeID(e queries.GatewayEnvelopesView) string {
	return fmt.Sprintf("%v-%v", e.OriginatorNodeID, e.OriginatorSequenceID)
}

func startSubscribeWorker(
	ctx context.Context,
	logger *zap.Logger,
	store *db.Handler,
) (*subscribeWorker, error) {
	logger = logger.Named(utils.SubscribeWorkerLoggerName)

	// TODO: Use some actual IDs.
	nodeIDs := []uint32{100, 200, 300}

	latestEnvelopes, err := store.ReadQuery().SelectVectorClock(ctx)
	if err != nil {
		logger.Error("failed to get vector clock", zap.Error(err))
		return nil, fmt.Errorf("could not create subscribe worker: %w", err)
	}
	vc := db.ToVectorClock(latestEnvelopes)

	logger.Debug("queried vector clock", zap.Any("vector_clock", vc))

	subscriptions := make(map[uint32]<-chan []queries.SelectGatewayEnvelopesByOriginatorsRow)

	// TODO: Paralelize this.
	for _, nodeID := range nodeIDs {
		logger := logger.With(utils.OriginatorIDField(nodeID))

		query := func(ctx context.Context, lastSeen int64, numRows int32) ([]queries.SelectGatewayEnvelopesByOriginatorsRow, int64, error) {
			logger.Debug("running pollable query",
				zap.Int64("last_seen", lastSeen))

			envs, err := store.ReadQuery().SelectGatewayEnvelopesByOriginators(ctx,
				// TODO: Check this query - what are originator node IDs and what are cursor node IDs?
				queries.SelectGatewayEnvelopesByOriginatorsParams{
					OriginatorNodeIds: []int32{int32(nodeID)},
					CursorSequenceIds: []int64{lastSeen},
					RowLimit:          numRows,
				})
			if err != nil {
				logger.Error("failed to get envelopes",
					zap.Error(err))
				return nil, 0, fmt.Errorf("could not get envelopes: %w", err)
			}

			last := lastSeen

			for i, env := range envs {

				seqID := uint64(env.OriginatorSequenceID)

				logger.Debug("processing envelope",
					zap.String("env_no", fmt.Sprintf("%v/%v", i+1, len(envs))),
					// zap.String("short_id", shortEnvelopeID(env)),
					utils.SequenceIDField(int64(seqID)),
				)

				logger.Debug("checking last seen for this node",
					// zap.String("short_id", shortEnvelopeID(env)),
					zap.Uint64("sequence_id", seqID),
					zap.Int64("last_seen", last),
				)

				if env.OriginatorSequenceID < last {
					logger.Fatal("system invariant broken: unsorted envelope stream",
						utils.SequenceIDField(env.OriginatorSequenceID),
						utils.LastSequenceIDField(last))
				}

				last = env.OriginatorSequenceID
			}

			return envs, last, nil
		}

		sub := db.NewDBSubscription(ctx, logger, query, int64(vc[nodeID]),
			db.PollingOptions{
				Interval: SubscribeWorkerPollTime,
				NumRows:  subscribeWorkerPollRows,
			})

		ch, err := sub.Start()
		if err != nil {
			logger.Error("failed to start subscription", zap.Error(err))
			return nil, fmt.Errorf("could not start subscription: %w", err)
		}

		subscriptions[nodeID] = ch
	}

	worker := &subscribeWorker{
		ctx:                 ctx,
		logger:              logger,
		emptyListeners:      listenerSet{},
		dbSubscriptions:     subscriptions,
		originatorListeners: listenersMap[uint32]{},
		topicListeners:      listenersMap[string]{},
	}

	go worker.start()
	logger.Debug("started")

	return worker, nil
}

func merge[T any](ch ...<-chan T) <-chan T {
	var (
		wg  sync.WaitGroup
		out = make(chan T)
	)

	wg.Add(len(ch))

	// Function will forward entries from its channel to the common channel.
	fw := func(c <-chan T) {
		for e := range c {
			out <- e
		}

		// Once our channel is done, signal that we completed.
		wg.Done()
	}

	// Start a forwarding goroutine for each channel.
	for _, c := range ch {
		go fw(c)
	}

	go func() {
		// When all goroutines have completed (all channels done) - close our common, merged channel.
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *subscribeWorker) start() {
	var subs []<-chan []queries.SelectGatewayEnvelopesByOriginatorsRow
	for _, sub := range s.dbSubscriptions {
		subs = append(subs, sub)
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case batch, ok := <-merge(subs...):
			if !ok {
				s.logger.Error("database subscription is closed")
				return
			}

			s.logger.Debug("received new batch", utils.NumEnvelopesField(len(batch)))
			envs := make([]*envelopes.OriginatorEnvelope, 0, len(batch))
			for _, row := range batch {
				env, err := envelopes.NewOriginatorEnvelopeFromBytes(row.OriginatorEnvelope)
				if err != nil {
					s.logger.Error("failed to unmarshal envelope", zap.Error(err))
					continue
				}
				envs = append(envs, env)
			}
			s.dispatchToOriginators(envs)
			s.dispatchToTopics(envs)
			s.dispatchToEmpties()
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
	// We iterate envelopes one-by-one, because we expect the number of envelopes
	// per-topic to be small in each tick
	for _, env := range envs {
		listeners := s.topicListeners.getListeners(env.TargetTopic().String())
		s.dispatchToListeners(listeners, []*envelopes.OriginatorEnvelope{env})
	}
}

func (s *subscribeWorker) dispatchToEmpties() {
	// only keep this to possibly close listeners
	s.dispatchToListeners(&s.emptyListeners, []*envelopes.OriginatorEnvelope{})
}

func (s *subscribeWorker) dispatchToListeners(
	listeners *listenerSet,
	envs []*envelopes.OriginatorEnvelope,
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
			if s.logger.Core().Enabled(zap.DebugLevel) {
				s.logger.Debug("stream closed, removing listener", utils.BodyField(l.ch))
			}

			s.closeListener(l)

		default:
			if len(envs) == 0 {
				return true
			}

			select {
			case l.ch <- envs:
				if s.logger.Core().Enabled(zap.DebugLevel) {
					s.logger.Debug("sent envelopes to listener", utils.NumEnvelopesField(len(envs)))
				}

			default:
				if s.logger.Core().Enabled(zap.DebugLevel) {
					s.logger.Debug("channel full, removing listener", utils.BodyField(l.ch))
				}

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
		if l.isEmpty {
			s.emptyListeners.Delete(l)
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
	l := newListener(ctx, s.logger, query, ch)

	if l.isEmpty {
		s.emptyListeners.Store(l, struct{}{})
	} else if len(l.topics) > 0 {
		s.topicListeners.addListener(l.topics, l)
	} else if len(l.originators) > 0 {
		s.originatorListeners.addListener(l.originators, l)
	}

	return ch
}

package message

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	subscriptionBufferSize  = 1024
	SubscribeWorkerPollTime = 100 * time.Millisecond
	// based on measurements in testnet using PG, we can poll at most 1000 elements in a large DB
	// this gives us sufficient throughput if being run continually
	subscribeWorkerPollRows = 1000
)

// A worker that listens for new envelopes in the DB and sends them to subscribers
// Assumes that there are many listeners - non-blocking updates are sent on buffered channels
// and may be dropped if full
type subscribeWorker struct {
	ctx    context.Context
	logger *zap.Logger
	store  *db.Handler

	// Keep track of known originators.
	registry registry.NodeRegistry

	// Poll envelopes per originator
	subscriptions *subscriptionHandler

	// Assumption: listeners cannot be in multiple slices
	emptyListeners      listenerSet
	originatorListeners listenersMap[uint32]
	topicListeners      listenersMap[string]
}

func (s *subscribeWorker) getOriginatorNodeIds() ([]uint32, error) {
	// Get initial list of nodes.
	nodes, err := s.registry.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("could not get list of originators: %w", err)
	}

	var nodeIDs []uint32
	for _, node := range nodes {
		if !node.IsCanonical {
			s.logger.Debug(
				"skipping non-canonical node",
				utils.OriginatorIDField(node.NodeID),
			)
			continue
		}

		nodeIDs = append(nodeIDs, node.NodeID)
	}

	return nodeIDs, nil
}

func startSubscribeWorker(
	ctx context.Context,
	logger *zap.Logger,
	store *db.Handler,
	registry registry.NodeRegistry,
) (*subscribeWorker, error) {
	logger = logger.Named(utils.SubscribeWorkerLoggerName)

	latestEnvelopes, err := store.ReadQuery().SelectVectorClock(ctx)
	if err != nil {
		logger.Error("failed to get vector clock", zap.Error(err))
		return nil, fmt.Errorf("could not create subscribe worker: %w", err)
	}
	vc := db.ToVectorClock(latestEnvelopes)

	worker := &subscribeWorker{
		ctx:                 ctx,
		logger:              logger,
		emptyListeners:      listenerSet{},
		store:               store,
		registry:            registry,
		subscriptions:       newSubscriptionHandler(logger, store, vc),
		originatorListeners: listenersMap[uint32]{},
		topicListeners:      listenersMap[string]{},
	}

	nodeIDs, err := worker.getOriginatorNodeIds()
	if err != nil {
		logger.Error("failed to get list of originators", zap.Error(err))
		return nil, fmt.Errorf("could not get list of originators: %w", err)
	}

	// NOTE: This can be done in parallel.
	for _, id := range nodeIDs {
		err = worker.subscriptions.newSubscription(ctx, id)
		if err != nil {
			logger.Error(
				"could not create new subscription",
				utils.OriginatorIDField(id),
				zap.Error(err),
			)
			return nil, fmt.Errorf(
				"could not create new subscription (originator: %v): %w",
				id,
				err,
			)
		}
	}

	go worker.monitorNodeChanges()
	go worker.start()
	logger.Debug("started")

	return worker, nil
}

func (s *subscribeWorker) start() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case batch, ok := <-s.subscriptions.allSubscriptions():
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

func (s *subscribeWorker) monitorNodeChanges() {
	newNodes := s.registry.OnNewNodes()
	removedNodes := s.registry.OnRemovedNodes()
	for {
		select {
		case <-s.ctx.Done():
			return

		case nodes := <-newNodes:
			for _, node := range nodes {

				if !node.IsCanonical {
					s.logger.Debug(
						"skipping non-canonical node",
						utils.OriginatorIDField(node.NodeID),
					)
					continue
				}

				err := s.subscriptions.newSubscription(s.ctx, node.NodeID)
				if err != nil {
					s.logger.Error(
						"could not add subscription for new node",
						utils.OriginatorIDField(node.NodeID),
						zap.Error(err),
					)
				}
			}
		case ids := <-removedNodes:
			for _, id := range ids {
				s.subscriptions.removeSubscription(id)
			}
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

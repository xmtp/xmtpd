package message

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type subscriptionHandler struct {
	*sync.Mutex

	logger *zap.Logger
	store  *db.Handler

	subs       map[uint32]*envelopePoller
	mergedSubs *funnel[[]queries.SelectGatewayEnvelopesBySingleOriginatorRow]

	cursor db.VectorClock
}

type envelopePoller struct {
	cancel context.CancelFunc
	// TODO: Check - queries.GatewayEnvelopesView and queries.SelectGatewayEnvelopesByOriginatorsRow
	// models are identical and they're overly verbose.
	ch <-chan []queries.SelectGatewayEnvelopesBySingleOriginatorRow
}

func newSubscriptionHandler(
	logger *zap.Logger,
	store *db.Handler,
	cursor db.VectorClock,
) *subscriptionHandler {
	s := &subscriptionHandler{
		Mutex:      &sync.Mutex{},
		subs:       make(map[uint32]*envelopePoller),
		logger:     logger,
		store:      store,
		cursor:     cursor,
		mergedSubs: newFunnel[[]queries.SelectGatewayEnvelopesBySingleOriginatorRow](),
	}

	return s
}

func (s *subscriptionHandler) newSubscription(ctx context.Context, id uint32) (retErr error) {
	// NOTE: DB Subscription currently supports stopping only via context cancellation, so we need per-sub cancel function
	// so that we can stop polling for any originator.
	childCtx, cancel := context.WithCancel(ctx)
	defer func() {
		// NOTE: In case this function fails, cancel the above context.
		if retErr != nil {
			cancel()
		}
	}()

	// TODO: Check handling of lastSeen - does it need to be saved to the outside of the sub, or does it get saved in the subscription?
	// NOTE: I think it's handled in the subscription, so the cursor can be removed here.

	query := func(ctx context.Context, lastSeen int64, numRows int32) ([]queries.SelectGatewayEnvelopesBySingleOriginatorRow, int64, error) {
		envs, err := s.store.ReadQuery().SelectGatewayEnvelopesBySingleOriginator(ctx,
			queries.SelectGatewayEnvelopesBySingleOriginatorParams{
				OriginatorNodeID: int32(id),
				CursorSequenceID: lastSeen,
				RowLimit:         numRows,
			})
		if err != nil {
			s.logger.Error("failed to get envelopes",
				zap.Error(err), utils.OriginatorIDField(id))
			return nil, 0, fmt.Errorf("could not get envelopes: %w", err)
		}

		last := lastSeen

		if len(envs) > 0 {
			s.logger.Debug("pollable query returned results",
				zap.Int64("last_seen", lastSeen),
				zap.Int("count", len(envs)),
			)
		}

		for _, env := range envs {

			if env.OriginatorSequenceID < last {
				s.logger.Fatal("system invariant broken: unsorted envelope stream",
					utils.SequenceIDField(env.OriginatorSequenceID),
					utils.LastSequenceIDField(last))
			}

			last = env.OriginatorSequenceID
		}

		return envs, last, nil
	}

	sub := db.NewDBSubscription(childCtx, s.logger, query, int64(s.cursor[id]),
		db.PollingOptions{
			Interval: SubscribeWorkerPollTime,
			NumRows:  subscribeWorkerPollRows,
		})

	ch, err := sub.Start()
	if err != nil {
		s.logger.Error(
			"failed to create new subscription",
			utils.OriginatorIDField(id),
			zap.Error(err),
		)
		return fmt.Errorf("could not start subscription (id: %v): %w", id, err)
	}

	// Per node/originator poller and cancellation.
	e := &envelopePoller{
		cancel: cancel,
		ch:     ch,
	}

	// Save the poller in the subscription handler.
	s.Lock()
	defer s.Unlock()
	s.subs[id] = e

	s.mergedSubs.addChannel(ch)

	return nil
}

// allSubscriptions returns a channel merging all individual subscription channels.
func (s *subscriptionHandler) allSubscriptions() <-chan []queries.SelectGatewayEnvelopesBySingleOriginatorRow {
	s.Lock()
	defer s.Unlock()

	return s.mergedSubs.output()
}

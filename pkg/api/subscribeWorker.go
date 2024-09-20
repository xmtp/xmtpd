package api

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
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

type listener = chan<- []*message_api.OriginatorEnvelope

// A worker that listens for new envelopes in the DB and sends them to subscribers
// Assumes that there are many listeners - non-blocking updates are sent on buffered channels
// and may be dropped if full
type subscribeWorker struct {
	ctx context.Context
	log *zap.Logger

	dbSubscription <-chan []queries.GatewayEnvelope
	// Assumption: listeners cannot be in multiple slices
	globalListeners     []listener
	originatorListeners map[uint32][]listener
	topicListeners      map[string][]listener
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
		globalListeners:     make([]listener, 0),
		originatorListeners: make(map[uint32][]listener),
		topicListeners:      make(map[string][]listener),
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

func (s *subscribeWorker) dispatch(
	row *queries.GatewayEnvelope,
) {
	bytes := row.OriginatorEnvelope
	env := &message_api.OriginatorEnvelope{}
	err := proto.Unmarshal(bytes, env)
	if err != nil {
		s.log.Error("Failed to unmarshal envelope", zap.Error(err))
		return
	}
	for _, listener := range s.originatorListeners[uint32(row.OriginatorNodeID)] {
		select {
		case listener <- []*message_api.OriginatorEnvelope{env}:
		default: // TODO(rich) Close and clean up channel
		}
	}
	for _, listener := range s.topicListeners[hex.EncodeToString(row.Topic)] {
		select {
		case listener <- []*message_api.OriginatorEnvelope{env}:
		default:
		}
	}
	for _, listener := range s.globalListeners {
		select {
		case listener <- []*message_api.OriginatorEnvelope{env}:
		default:
		}
	}
}

func (s *subscribeWorker) listen(
	requests []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest,
) (<-chan []*message_api.OriginatorEnvelope, error) {
	subscribeAll := false
	topics := make(map[string]bool, len(requests))
	originators := make(map[uint32]bool, len(requests))

	if len(requests) > maxSubscriptionsPerClient {
		return nil, fmt.Errorf(
			"too many subscriptions: %d, consider subscribing to fewer topics or subscribing without a filter",
			len(requests),
		)
	}
	for _, req := range requests {
		enum := req.GetQuery().GetFilter()
		if enum == nil {
			subscribeAll = true
		}
		switch filter := enum.(type) {
		case *message_api.EnvelopesQuery_Topic:
			if len(filter.Topic) == 0 || len(filter.Topic) > maxTopicLength {
				return nil, status.Errorf(codes.InvalidArgument, "invalid topic")
			}
			topics[hex.EncodeToString(filter.Topic)] = true
		case *message_api.EnvelopesQuery_OriginatorNodeId:
			originators[filter.OriginatorNodeId] = true
		default:
			subscribeAll = true
		}
	}

	ch := make(chan []*message_api.OriginatorEnvelope, subscriptionBufferSize)

	if subscribeAll {
		if len(topics) > 0 || len(originators) > 0 {
			return nil, fmt.Errorf("cannot filter by topic or originator when subscribing to all")
		}
		// TODO(rich) thread safety
		s.globalListeners = append(s.globalListeners, ch)
	} else if len(topics) > 0 {
		if len(originators) > 0 {
			return nil, fmt.Errorf("cannot filter by both topic and originator in same subscription request")
		}
		for topic := range topics {
			s.topicListeners[topic] = append(s.topicListeners[topic], ch)
		}
	} else if len(originators) > 0 {
		for originator := range originators {
			s.originatorListeners[originator] = append(s.originatorListeners[originator], ch)
		}
	}

	return ch, nil
}

package messagev1

import (
	"context"
	"errors"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/node/subscribers"
	"github.com/xmtp/xmtpd/pkg/node/topics"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	ErrTODO          = errors.New("TODO")
	ErrMissingTopic  = errors.New("missing topic")
	ErrTooManyTopics = errors.New("too many topics")
)

type Service struct {
	messagev1.UnimplementedMessageApiServer

	// Configured as constructor options.
	log         *zap.Logger
	topics      topics.Manager
	subs        subscribers.Manager
	store       crdt.Store
	broadcaster crdt.Broadcaster
	syncer      crdt.Syncer

	// Configured internally.
	ctx       context.Context
	ctxCancel func()
}

func New(log *zap.Logger, topics topics.Manager, subs subscribers.Manager, store crdt.Store, bc crdt.Broadcaster, syncer crdt.Syncer) (*Service, error) {
	log = log.Named("message/v1")
	s := &Service{
		log:         log,
		topics:      topics,
		subs:        subs,
		store:       store,
		broadcaster: bc,
		syncer:      syncer,
	}
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	return s, nil
}

func (s *Service) Close() {
	if s.ctxCancel != nil {
		s.ctxCancel()
	}
}

func (s *Service) Publish(ctx context.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	for _, env := range req.Envelopes {
		topic, err := s.topics.GetOrCreateTopic(ctx, env.ContentTopic)
		if err != nil {
			return nil, err
		}
		ev, err := topic.BroadcastAppend(ctx, env)
		if err != nil {
			return nil, err
		}
		s.log.Debug("envelope published", zap.Cid("event", ev.Cid))
	}
	return &messagev1.PublishResponse{}, nil
}

func (s *Service) Subscribe(req *messagev1.SubscribeRequest, stream messagev1.MessageApi_SubscribeServer) error {
	if len(req.ContentTopics) == 0 {
		return ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return ErrTooManyTopics
	}
	topicName := req.ContentTopics[0]

	ctx := context.Background()
	eventsCh := s.subs.Subscribe(ctx, topicName, 100)
	defer s.subs.Unsubscribe(topicName, eventsCh)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-stream.Context().Done():
			return nil
		case ev := <-eventsCh:
			err := stream.Send(ev.Envelope)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Service) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	if len(req.ContentTopics) == 0 {
		return nil, ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}

	return s.store.Query(ctx, req)
}

func (s *Service) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	return ErrTODO
}

func (s *Service) BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	return &messagev1.BatchQueryResponse{}, ErrTODO
}

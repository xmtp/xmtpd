package messagev1

import (
	"context"
	"errors"
	"sync"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	crdtmemstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	crdttypes "github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/store"
	"github.com/xmtp/xmtpd/pkg/zap"
	"google.golang.org/protobuf/proto"
)

var (
	ErrTODO               = errors.New("TODO")
	ErrMissingTopic       = errors.New("missing topic")
	ErrTooManyTopics      = errors.New("too many topics")
	ErrTopicAlreadyExists = errors.New("topic already exists")
)

type Service struct {
	messagev1.UnimplementedMessageApiServer

	// Configured as constructor options.
	log   *zap.Logger
	store store.Store

	// Configured internally.
	ctx       context.Context
	ctxCancel func()

	topicSubs     map[string]map[chan *crdttypes.Event]struct{}
	topicSubsLock sync.RWMutex

	topics     map[string]*crdt.Replica
	topicsLock sync.RWMutex
}

func New(log *zap.Logger, store store.Store) (*Service, error) {
	log = log.Named("message/v1")
	s := &Service{
		log:   log,
		store: store,

		topicSubs: map[string]map[chan *crdttypes.Event]struct{}{},
		topics:    map[string]*crdt.Replica{},
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
		topic, err := s.getOrCreateTopic(ctx, env.ContentTopic)
		if err != nil {
			return nil, err
		}
		envB, err := proto.Marshal(env)
		if err != nil {
			return nil, err
		}
		ev, err := topic.BroadcastAppend(ctx, envB)
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
	eventsCh := s.subscribe(ctx, topicName, 100)
	defer s.unsubscribe(topicName, eventsCh)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-stream.Context().Done():
			return nil
		case ev := <-eventsCh:
			var env messagev1.Envelope
			err := proto.Unmarshal(ev.Payload, &env)
			if err != nil {
				return err
			}
			err = stream.Send(&env)
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

	return s.store.QueryEnvelopes(ctx, req)
}

func (s *Service) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	return ErrTODO
}

func (s *Service) BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	return &messagev1.BatchQueryResponse{}, ErrTODO
}

func (s *Service) getOrCreateTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	topic, err := s.getTopic(ctx, topicId)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		topic, err = s.createTopic(ctx, topicId)
		if err != nil {
			return nil, err
		}
	}
	return topic, nil
}

func (s *Service) getTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	s.topicsLock.RLock()
	defer s.topicsLock.RUnlock()
	topic, ok := s.topics[topicId]
	if !ok {
		return nil, nil
	}
	return topic, nil
}

func (s *Service) createTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()
	if _, ok := s.topics[topicId]; ok {
		return nil, ErrTopicAlreadyExists
	}
	store := crdtmemstore.New(s.log)
	return crdt.NewReplica(
		ctx,
		s.log,
		// TODO: these factories/makers should be passed in as options/config
		store,
		membroadcaster.New(s.log),
		memsyncer.New(s.log, store),
		func(ev *crdttypes.Event) {
			s.topicSubsLock.RLock()
			defer s.topicSubsLock.RUnlock()

			if _, ok := s.topicSubs[topicId]; !ok {
				return
			}
			for ch := range s.topicSubs[topicId] {
				// TODO: what to do if channel is full here
				ch <- ev
			}
		},
	)
}

func (s *Service) subscribe(ctx context.Context, topicId string, buffer int) chan *crdttypes.Event {
	s.topicSubsLock.Lock()
	defer s.topicSubsLock.Unlock()

	ch := make(chan *crdttypes.Event, buffer)
	if _, ok := s.topicSubs[topicId]; !ok {
		s.topicSubs[topicId] = map[chan *crdttypes.Event]struct{}{}
	}
	s.topicSubs[topicId][ch] = struct{}{}
	return ch
}

func (s *Service) unsubscribe(topicId string, ch chan *crdttypes.Event) {
	s.topicSubsLock.Lock()
	defer s.topicSubsLock.Unlock()

	if _, ok := s.topicSubs[topicId]; !ok {
		return
	}
	delete(s.topicSubs[topicId], ch)
	close(ch)
	if len(s.topicSubs[topicId]) == 0 {
		delete(s.topicSubs, topicId)
	}
}

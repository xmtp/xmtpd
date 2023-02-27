package node

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"
)

var (
	ErrTODO               = errors.New("TODO")
	ErrMissingTopic       = errors.New("missing topic")
	ErrTooManyTopics      = errors.New("too many topics")
	ErrTopicAlreadyExists = errors.New("topic already exists")
)

type StoreMakerFunc func(topic string) (crdt.Store, error)

type Node struct {
	messagev1.UnimplementedMessageApiServer

	log        *zap.Logger
	storeMaker StoreMakerFunc

	api *apigateway.Server

	ctx       context.Context
	ctxCancel context.CancelFunc

	topics      map[string]*crdt.Replica
	topicStores map[string]crdt.Store
	topicsLock  sync.RWMutex

	host  host.Host
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	broadcasters     map[string]chan *types.Event
	broadcastersLock sync.RWMutex

	ns *server.Server
	nc *nats.Conn
}

func New(ctx context.Context, log *zap.Logger, storeMaker StoreMakerFunc, opts *Options) (*Node, error) {
	n := &Node{
		log:        log,
		storeMaker: storeMaker,

		topics:      map[string]*crdt.Replica{},
		topicStores: map[string]crdt.Store{},

		broadcasters: map[string]chan *types.Event{},
	}
	n.ctx, n.ctxCancel = context.WithCancel(ctx)
	var err error

	// Initialize API server/gateway.
	n.api, err = apigateway.New(n.ctx, log, n, &opts.API)
	if err != nil {
		return nil, errors.Wrap(err, "initializing api")
	}

	// Initialize libp2p host.
	n.host, err = libp2p.New()
	if err != nil {
		return nil, err
	}

	// Initialize libp2p pubsub.
	gs, err := pubsub.NewGossipSub(ctx, n.host)
	if err != nil {
		return nil, err
	}

	// Initialize libp2p pubsub topic.
	n.topic, err = gs.Join("/xmtp/0")
	if err != nil {
		return nil, err
	}

	// Initialize libp2p pubsub topic subscription.
	n.sub, err = n.topic.Subscribe()
	if err != nil {
		return nil, err
	}

	// Initialize libp2p events consumer.
	go n.p2pEventConsumerLoop()

	// Initialize nats for API subscribers.
	n.ns, err = server.NewServer(&server.Options{
		Port: server.RANDOM_PORT,
	})
	if err != nil {
		return nil, err
	}
	go n.ns.Start()
	if !n.ns.ReadyForConnections(4 * time.Second) {
		return nil, errors.New("nats not ready")
	}
	n.nc, err = nats.Connect(n.ns.ClientURL())
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) Close() {
	if n.api != nil {
		n.api.Close()
	}

	if n.nc != nil {
		n.nc.Close()
	}

	if n.ns != nil {
		n.ns.Shutdown()
	}

	if n.host != nil {
		n.host.Close()
	}

	if n.sub != nil {
		n.sub.Cancel()
	}

	if n.topic != nil {
		n.topic.Close()
	}

	if n.ctxCancel != nil {
		n.ctxCancel()
	}

	for _, store := range n.topicStores {
		store.Close()
	}
}

func (n *Node) APIHTTPListenPort() uint {
	return n.api.HTTPListenPort()
}

func (n *Node) Connect(ctx context.Context, addr peer.AddrInfo) error {
	return n.host.Connect(ctx, addr)
}

func (n *Node) Address() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    n.host.ID(),
		Addrs: n.host.Addrs(),
	}
}

func (n *Node) Publish(ctx context.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	for _, env := range req.Envelopes {
		topic, err := n.getOrCreateTopic(ctx, env.ContentTopic)
		if err != nil {
			return nil, err
		}
		ev, err := topic.BroadcastAppend(ctx, env)
		if err != nil {
			return nil, err
		}
		n.log.Debug("envelope published", zap.Cid("event", ev.Cid))
	}
	return &messagev1.PublishResponse{}, nil
}

func (n *Node) Subscribe(req *messagev1.SubscribeRequest, stream messagev1.MessageApi_SubscribeServer) error {
	if len(req.ContentTopics) == 0 {
		return ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return ErrTooManyTopics
	}
	topic := req.ContentTopics[0]

	// Send subscribe confirmation.
	n.log.Debug("sending subscribe confirmation", zap.String("topic", topic))
	err := stream.Send(&messagev1.Envelope{})
	if err != nil {
		return err
	}

	sub, err := n.nc.Subscribe(topic, func(msg *nats.Msg) {
		ev, err := types.EventFromBytes(msg.Data)
		if err != nil {
			n.log.Error("error parsing event from bytes", zap.Error(err))
			return
		}
		err = stream.Send(ev.Envelope)
		if err != nil {
			n.log.Error("error emitting new event", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()

	select {
	case <-n.ctx.Done():
		return nil
	case <-stream.Context().Done():
		return nil
	}
}

func (n *Node) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	if len(req.ContentTopics) == 0 {
		return nil, ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}
	topic := req.ContentTopics[0]

	replica, err := n.getOrCreateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}

	return replica.Query(ctx, req)
}

func (n *Node) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	// Subscribe to all nats subjects via wildcard
	// https://docs.nats.io/nats-concepts/subjects#wildcards
	return n.Subscribe(&messagev1.SubscribeRequest{
		ContentTopics: []string{"*"},
	}, stream)
}

func (n *Node) BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	res := &messagev1.BatchQueryResponse{}
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	for _, r := range req.Requests {
		r := r
		g.Go(func() error {
			rs, err := n.Query(ctx, r)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			res.Responses = append(res.Responses, rs)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *Node) getOrCreateTopic(ctx context.Context, topic string) (*crdt.Replica, error) {
	replica, err := n.getTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	if replica == nil {
		replica, err = n.createTopic(ctx, topic)
		if err != nil {
			return nil, err
		}
	}
	return replica, nil
}

func (n *Node) getTopic(ctx context.Context, topic string) (*crdt.Replica, error) {
	n.log.Debug("getting topic", zap.String("topic", topic))
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	replica, ok := n.topics[topic]
	if !ok {
		return nil, nil
	}
	return replica, nil
}

func (n *Node) createTopic(ctx context.Context, topic string) (*crdt.Replica, error) {
	n.log.Debug("creating topic", zap.String("topic", topic))
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	if _, ok := n.topics[topic]; ok {
		return nil, ErrTopicAlreadyExists
	}
	bc, err := n.getOrCreateBroadcaster(topic)
	if err != nil {
		return nil, err
	}
	store, err := n.storeMaker(topic)
	if err != nil {
		return nil, err
	}
	replica, err := crdt.NewReplica(ctx, n.log, store, bc, nil, func(ev *types.Event) {
		evB, err := ev.ToBytes()
		if err != nil {
			n.log.Error("error converting event to bytes", zap.Error(err))
			return
		}
		err = n.nc.Publish(ev.ContentTopic, evB)
		if err != nil {
			n.log.Error("error publishing replicated event")
		}
	})
	if err != nil {
		return nil, err
	}
	n.topics[topic] = replica
	n.topicStores[topic] = store
	return replica, nil
}

func (n *Node) p2pEventConsumerLoop() {
	for {
		msg, err := n.sub.Next(n.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			n.log.Error("error getting next event", zap.Error(err))
			continue
		}
		ev, err := types.EventFromBytes(msg.Data)
		if err != nil {
			n.log.Error("error unmarshaling event", zap.Error(err))
			continue
		}

		_, err = n.getOrCreateTopic(n.ctx, ev.ContentTopic)
		if err != nil {
			n.log.Error("error getting or creating topic", zap.Error(err))
			continue
		}

		// Push onto broadcaster channel to be consumed by it's replica via Next.
		bc, err := n.getOrCreateBroadcaster(ev.ContentTopic)
		if err != nil {
			n.log.Error("error getting broadcaster", zap.Error(err))
		}
		bc.C <- ev
	}
}

func (n *Node) getOrCreateBroadcaster(topic string) (*broadcaster, error) {
	n.broadcastersLock.Lock()
	defer n.broadcastersLock.Unlock()

	if _, ok := n.broadcasters[topic]; !ok {
		n.broadcasters[topic] = make(chan *types.Event, 100)
	}
	return newBroadcaster(n.topic, n.broadcasters[topic])
}

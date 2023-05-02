package node

import (
	gocontext "context"
	"fmt"
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
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"
)

var (
	ErrUnknownTopic       = errors.New("topic does not exist")
	ErrMissingTopic       = errors.New("missing topic")
	ErrTooManyTopics      = errors.New("too many topics")
	ErrTopicAlreadyExists = errors.New("topic already exists")
)

const (
	pubsubTopic = "/xmtp/0"
)

type Node struct {
	messagev1.UnimplementedMessageApiServer

	log *zap.Logger
	ctx context.Context

	topics     map[string]*crdt.Replica
	topicsLock sync.RWMutex

	host  host.Host
	gs    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	api *apigateway.Server

	store            NodeStore
	broadcasters     map[string]chan *types.Event
	broadcastersLock sync.RWMutex

	ns *server.Server
	nc *nats.Conn

	metrics *Metrics
	peers   *persistentPeers
}

func New(ctx context.Context, metrics *Metrics, store NodeStore, opts *Options) (*Node, error) {
	n := &Node{
		ctx:          ctx,
		log:          ctx.Logger(),
		metrics:      metrics,
		store:        store,
		topics:       map[string]*crdt.Replica{},
		broadcasters: map[string]chan *types.Event{},
	}

	var err error

	// Initialize API server/gateway.
	var apiMetrics *apigateway.Metrics
	if metrics != nil {
		apiMetrics = metrics.API
	}
	n.api, err = apigateway.New(n.ctx, n, apiMetrics, &opts.API)
	if err != nil {
		return nil, errors.Wrap(err, "initializing api")
	}

	// Initialize libp2p host.
	privKey, err := getOrCreatePrivateKey(opts.P2P.NodeKey)
	if err != nil {
		return nil, err
	}
	n.host, err = libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", opts.P2P.Port)),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return nil, err
	}
	n.log = n.log.With(zap.PeerID("node", n.host.ID()))
	n.ctx = context.WithLogger(n.ctx, n.log)
	n.log.Info("p2p listening", zap.Strings("addresses", n.P2PListenAddresses()))

	// Initialize libp2p pubsub.
	n.gs, err = pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return nil, err
	}

	// Initialize libp2p pubsub topic.
	n.topic, err = n.gs.Join(pubsubTopic)
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

	// Find pre-existing topics
	topics, err := store.Topics()
	if err != nil {
		return nil, err
	}
	// Bootstrap all the topics with some parallelization.
	grp, _ := errgroup.WithContext(ctx)
	grp.SetLimit(1000) // up to 1000 topic bootstraps in parallel
	for _, name := range topics {
		topic := name
		grp.Go(func() (err error) {
			_, err = n.createTopic(topic)
			return err
		})
	}
	// Do not return until all topics are bootstrapped successfully.
	// If any bootstrap fails, bail out.
	if err := grp.Wait(); err != nil {
		return nil, err
	}

	n.setSyncHandler()

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

	n.peers, err = newPersistentPeers(n.ctx, n.log, n.host, opts.P2P.PersistentPeers)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) Close() {
	// Shut off the clients
	if n.api != nil {
		n.api.Close()
	}

	if n.nc != nil {
		n.nc.Close()
	}

	if n.ns != nil {
		n.ns.Shutdown()
	}

	// Shut down the topics
	n.ctx.Close()

	// Shut down all the infrastructure
	if n.sub != nil {
		n.sub.Cancel()
	}

	if n.topic != nil {
		n.topic.Close()
	}

	if n.host != nil {
		n.host.Close()
	}

	if n.store != nil {
		n.store.Close()
	}

}

func (n *Node) APIHTTPListenPort() uint {
	return n.api.HTTPListenPort()
}

func (n *Node) P2PListenAddresses() []string {
	exclude := map[string]bool{
		"/p2p-circuit": true,
	}
	addrs := []string{}
	for _, ma := range n.host.Network().ListenAddresses() {
		addr := ma.String()
		if exclude[addr] {
			continue
		}
		addrs = append(addrs, addr+"/p2p/"+n.host.ID().Pretty())
	}
	return addrs
}

func (n *Node) ID() peer.ID {
	return n.host.ID()
}

func (n *Node) Connect(ctx context.Context, addr peer.AddrInfo) error {
	return n.host.Connect(ctx, addr)
}

func (n *Node) ConnectedPeers() map[peer.ID]*peer.AddrInfo {
	return connectedPeers(n.host)
}

func (n *Node) PubSubPeers() []peer.ID {
	return n.topic.ListPeers()
}

func (n *Node) Disconnect(ctx context.Context, peer peer.ID) error {
	return n.host.Network().ClosePeer(peer)
}

func (n *Node) Address() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    n.host.ID(),
		Addrs: n.host.Addrs(),
	}
}

func (n *Node) Publish(gctx gocontext.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	ctx := context.New(gctx, n.log)
	for _, env := range req.Envelopes {
		start := time.Now()
		topic, err := n.getOrCreateTopic(env.ContentTopic)
		if err != nil {
			return nil, err
		}
		ev, err := topic.BroadcastAppend(ctx, env)
		if err != nil {
			return nil, err
		}
		n.log.Debug("envelope published",
			zap.Cid("event", ev.Cid),
			zap.String("topic", env.ContentTopic),
			zap.Int("timestamp_ns", int(env.TimestampNs)),
			zap.Duration("duration", time.Since(start)))
	}
	return &messagev1.PublishResponse{}, nil
}

func (n *Node) Subscribe(req *messagev1.SubscribeRequest, stream messagev1.MessageApi_SubscribeServer) error {
	if len(req.ContentTopics) == 0 {
		return ErrMissingTopic
	}

	var streamLock sync.Mutex
	for _, topic := range req.ContentTopics {
		sub, err := n.nc.Subscribe(topic, func(msg *nats.Msg) {
			ev, err := types.EventFromBytes(msg.Data)
			if err != nil {
				n.log.Error("error parsing event from bytes", zap.Error(err))
				return
			}
			func() {
				streamLock.Lock()
				defer streamLock.Unlock()
				err := stream.Send(ev.Envelope)
				if err != nil {
					n.log.Error("error emitting new event", zap.Error(err))
				}
			}()
		})
		if err != nil {
			return err
		}
		defer func() {
			_ = sub.Unsubscribe()
		}()

		// Send subscribe confirmation.
		func() {
			streamLock.Lock()
			defer streamLock.Unlock()
			n.log.Debug("sending subscribe confirmation", zap.String("topic", topic))
			err = stream.Send(&messagev1.Envelope{})
			if err != nil {
				n.log.Error("error emitting subscribe confirmation", zap.Error(err))
			}
		}()
	}

	select {
	case <-n.ctx.Done():
		return nil
	case <-stream.Context().Done():
		return nil
	}
}

func (n *Node) Query(gctx gocontext.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	start := time.Now()
	n.log.Debug("query", zap.Strings("topics", req.ContentTopics))
	if len(req.ContentTopics) == 0 {
		return nil, ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}
	topic := req.ContentTopics[0]

	replica, err := n.getTopic(topic)
	if err != nil {
		if err == ErrUnknownTopic {
			return &messagev1.QueryResponse{}, nil
		}
		return nil, err
	}

	resp, err := replica.Query(context.New(gctx, n.log), req)
	if err != nil {
		n.log.Debug("query error", zap.String("topic", topic), zap.Error(err), zap.Duration("duration", time.Since(start)))
	}
	n.log.Debug("query result", zap.String("topic", topic), zap.Int("envelopes", len(resp.Envelopes)), zap.Duration("duration", time.Since(start)))
	return resp, nil
}

func (n *Node) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	// Subscribe to all nats subjects via wildcard
	// https://docs.nats.io/nats-concepts/subjects#wildcards
	return n.Subscribe(&messagev1.SubscribeRequest{
		ContentTopics: []string{"*"},
	}, stream)
}

func (n *Node) BatchQuery(gctx gocontext.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	n.log.Debug("batch query", zap.Int("req-count", len(req.Requests)))
	res := &messagev1.BatchQueryResponse{}
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(gctx)
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

func (n *Node) getOrCreateTopic(topic string) (*crdt.Replica, error) {
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	if replica, ok := n.topics[topic]; ok {
		return replica, nil
	}
	return n.addTopic(topic)
}

func (n *Node) getTopic(topic string) (*crdt.Replica, error) {
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	replica, ok := n.topics[topic]
	if !ok {
		return nil, ErrUnknownTopic
	}
	return replica, nil
}

func (n *Node) createTopic(topic string) (*crdt.Replica, error) {
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	return n.addTopic(topic)
}

func (n *Node) addTopic(topic string) (*crdt.Replica, error) {
	log := n.log.With(zap.String("topic", topic))
	ctx := context.WithLogger(n.ctx, log)
	log.Debug("adding topic")
	if _, ok := n.topics[topic]; ok {
		return nil, ErrTopicAlreadyExists
	}
	bc, err := n.getOrCreateBroadcaster(topic)
	if err != nil {
		return nil, err
	}
	syn, err := n.getOrCreateSyncer(topic)
	if err != nil {
		return nil, err
	}
	store, err := n.store.NewTopic(topic)
	if err != nil {
		return nil, err
	}
	var crdtMetrics *crdt.Metrics
	if n.metrics != nil {
		crdtMetrics = n.metrics.Replicas
	}
	replica, err := crdt.NewReplica(ctx, crdtMetrics, store, bc, syn, func(ev *types.Event) {
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

		_, err = n.getOrCreateTopic(ev.ContentTopic)
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

func (n *Node) getOrCreateSyncer(topic string) (*syncer, error) {
	return &syncer{
		metrics: n.metrics,
		host:    n.host,
		topic:   topic,
	}, nil
}

package orbitdbnode

import (
	gocontext "context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/iface"
	orbitstores "berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/operation"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	ipfsconfig "github.com/ipfs/kubo/config"
	ipfscore "github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/bootstrap"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/repo"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/multiformats/go-multiaddr"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/otel"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

var (
	ErrTODO               = errors.New("TODO")
	ErrMissingTopic       = errors.New("missing topic")
	ErrTooManyTopics      = errors.New("too many topics")
	ErrTopicAlreadyExists = errors.New("topic already exists")

	infinity         = -1
	accessController = &accesscontroller.CreateAccessControllerOptions{
		Access: map[string][]string{
			"write": {"*"},
		},
	}
)

type Node struct {
	messagev1.UnimplementedMessageApiServer

	log *zap.Logger
	ctx context.Context

	topicsDB   orbitdb.KeyValueStore
	topics     map[string]orbitdb.EventLogStore
	topicsLock sync.RWMutex

	api *apigateway.Server

	ns *server.Server
	nc *nats.Conn

	ot *otel.OpenTelemetry

	ipfs  *ipfscore.IpfsNode
	orbit orbitdb.OrbitDB
}

func New(ctx context.Context, opts *Options) (*Node, error) {
	n := &Node{
		ctx:    ctx,
		log:    ctx.Logger(),
		topics: map[string]orbitdb.EventLogStore{},
	}

	var err error

	// Initialize open telemetry.
	n.ot, err = otel.New(n.ctx, &opts.OpenTelemetry)
	if err != nil {
		return nil, errors.Wrap(err, "initializing open telemetry")
	}

	// Initialize API server/gateway.
	n.api, err = apigateway.New(n.ctx, n, &opts.API)
	if err != nil {
		return nil, errors.Wrap(err, "initializing api")
	}

	// Initialize IPFS node.
	ipfsRepo, err := newIPFSRepo(ctx, &opts.P2P)
	if err != nil {
		return nil, err
	}
	n.ipfs, err = ipfscore.NewNode(ctx, &ipfscore.BuildCfg{
		// TODO: pass in bootstrap peers, or at least don't try to use default peers while starting up, is that online = false?
		Online: true,
		Repo:   ipfsRepo,
		// Host:   mock.MockHostOption(m),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, err
	}
	n.log = n.log.With(zap.String("node", n.ipfs.Identity.Pretty()))
	n.log.Info("ipfs node listening", zap.Strings("addresses", n.P2PListenAddresses()))
	peers := make([]peer.AddrInfo, 0, len(opts.P2P.BootstrapPeers))
	for _, addr := range opts.P2P.BootstrapPeers {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, errors.Wrap(err, "parsing persistent peer address")
		}
		peer, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return nil, errors.Wrap(err, "getting persistent peer address info")
		}
		if peer == nil {
			return nil, fmt.Errorf("persistent peer address info is nil: %s", addr)
		}
		if peer.ID == n.ipfs.Identity {
			continue
		}
		peers = append(peers, *peer)
	}
	err = n.ipfs.Bootstrap(bootstrap.BootstrapConfigWithPeers(peers))
	if err != nil {
		return nil, err
	}
	n.ipfs.IsOnline = true
	ipfsAPI, err := coreapi.NewCoreAPI(n.ipfs)
	if err != nil {
		return nil, err
	}
	// TODO: fix this data path
	dataPath := filepath.Join("orbitdb", n.ipfs.Identity.Pretty())
	n.orbit, err = orbitdb.NewOrbitDB(ctx, ipfsAPI, &orbitdb.NewOrbitDBOptions{
		Directory: &dataPath,
	})
	if err != nil {
		return nil, err
	}

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

	// Initialize topics registry replica.
	addr, err := n.orbit.DetermineAddress(ctx, "/xmtp/0/topics/0", "keyvalue", &iface.DetermineAddressOptions{
		AccessController: accessController,
	})
	if err != nil {
		return nil, err
	}

	// Open the data store, or creating it if necessary.
	n.topicsDB, err = n.orbit.KeyValue(n.ctx, addr.String(), &iface.CreateDBOptions{
		AccessController: accessController,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		// TODO: subscribe to new/replicated events so we can create the topic replica in real-time from here too
		go func() {
			sub, err := n.topicsDB.EventBus().Subscribe([]interface{}{
				new(orbitstores.EventReplicateProgress),
				new(orbitstores.EventReplicated),
			}, eventbus.BufSize(infinity+32))
			if err != nil {
				n.log.Error("error streaming from topics registry", zap.Error(err))
				// TODO: retry?
			}
			for {
				select {
				case <-n.ctx.Done():
					n.log.Debug("topics registry stream context closed")
					return
				case obj, ok := <-sub.Out():
					if !ok {
						return
					}
					switch ev := obj.(type) {
					case orbitstores.EventReplicateProgress:
						// n.log.Debug("topics registry replication progress", zap.Any("event", ev))
					case orbitstores.EventReplicated:
						for _, entry := range ev.Entries {
							op, err := operation.ParseOperation(entry)
							if err != nil {
								n.log.Error("error parsing topics registry replicated operation", zap.Error(err))
								continue
							}
							var topic string
							if op.GetKey() != nil {
								topic = *op.GetKey()
							}
							if topic == "" {
								continue
							}

							_, err = n.getOrCreateTopicReplica(ctx, topic)
							if err != nil {
								n.log.Error("error creating topic replica for replicated registry event", zap.Error(err))
								continue
							}

							n.log.Debug("received topics registry replicated event", zap.String("topic", topic))
						}
					}
				}
			}
		}()

		// Create local topic replicas for all existing topics in the registry.
		topics := n.topicsDB.All()
		fmt.Println("TOPICS", topics)
		for topic := range topics {
			_, err := n.getOrCreateTopicReplica(ctx, topic)
			if err != nil {
				n.log.Error("error initializing topic replica", zap.Error(err), zap.String("topic", topic))
			}
		}

		// Load log store data in case there are existing entries.
		// TODO: re-add this
		// err = n.topicsDB.Load(n.ctx, infinity)
		// if err != nil {
		// 	n.log.Error("error loading topics registry replica", zap.Error(err))
		// }
	}()

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

	n.ctx.Close()

	n.topicsDB.Close()

	for _, replica := range n.topics {
		replica.Close()
	}

	if n.orbit != nil {
		n.orbit.Close()
	}

	if n.ipfs != nil {
		n.ipfs.Close()
	}

	// Shut down telemetry
	if n.ot != nil {
		n.ot.Close()
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
	for _, ma := range n.ipfs.PeerHost.Network().ListenAddresses() {
		addr := ma.String()
		if exclude[addr] {
			continue
		}
		addrs = append(addrs, addr+"/p2p/"+n.ipfs.Identity.Pretty())
	}
	return addrs
}

func (n *Node) Publish(gctx gocontext.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	ctx := context.New(gctx, n.log)

	for _, env := range req.Envelopes {
		replica, err := n.getOrCreateTopicReplica(ctx, env.ContentTopic)
		if err != nil {
			return nil, err
		}

		envB, err := proto.Marshal(env)
		if err != nil {
			return nil, err
		}

		op, err := replica.Add(ctx, envB)
		if err != nil {
			return nil, err
		}

		err = n.nc.Publish(env.ContentTopic, envB)
		if err != nil {
			n.log.Error("error publishing published event")
		}

		n.log.Debug("envelope published", zap.String("event", op.GetEntry().GetHash().String()), zap.String("operation", string(op.GetEntry().GetPayload())))
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
		var env messagev1.Envelope
		err := proto.Unmarshal(msg.Data, &env)
		if err != nil {
			n.log.Error("error parsing event from bytes", zap.Error(err))
			return
		}
		err = stream.Send(&env)
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

func (n *Node) Query(gctx gocontext.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	if len(req.ContentTopics) == 0 {
		return nil, ErrMissingTopic
	} else if len(req.ContentTopics) > 1 {
		return nil, ErrTooManyTopics
	}
	// topic := req.ContentTopics[0]

	// replica, err := n.getOrCreateTopicReplica(topic)
	// if err != nil {
	// 	return nil, err
	// }

	// return replica.Query(context.New(gctx, n.log), req)
	return nil, ErrTODO
}

func (n *Node) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	// Subscribe to all nats subjects via wildcard
	// https://docs.nats.io/nats-concepts/subjects#wildcards
	return n.Subscribe(&messagev1.SubscribeRequest{
		ContentTopics: []string{"*"},
	}, stream)
}

func (n *Node) BatchQuery(gctx gocontext.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
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

func (n *Node) getOrCreateTopicReplica(ctx context.Context, topic string) (orbitdb.EventLogStore, error) {
	fmt.Println("HERE.getOrCreateTopicReplica.1")
	replica, err := n.getTopicReplica(ctx, topic)
	if err != nil {
		return nil, err
	}
	fmt.Println("HERE.getOrCreateTopicReplica.2")

	if replica == nil {
		val, err := n.topicsDB.Get(ctx, topic)
		fmt.Println("HERE.getOrCreateTopicReplica.3")
		if err != nil {
			return nil, err
		}

		fmt.Println("HERE.getOrCreateTopicReplica.4")

		if val == nil {
			fmt.Println("HERE.getOrCreateTopicReplica.5")
			_, err := n.topicsDB.Put(ctx, topic, []byte{1})
			if err != nil {
				return nil, err
			}
		}
		fmt.Println("HERE.getOrCreateTopicReplica.6")

		replica, err = n.createTopicReplica(ctx, topic)
		if err != nil {
			return nil, err
		}
		fmt.Println("HERE.getOrCreateTopicReplica.7")
	}
	fmt.Println("HERE.getOrCreateTopicReplica.8")

	return replica, nil
}

func (n *Node) getTopicReplica(ctx context.Context, topic string) (orbitdb.EventLogStore, error) {
	n.log.Debug("getting topic", zap.String("topic", topic))
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	replica, ok := n.topics[topic]
	if !ok {
		return nil, nil
	}
	return replica, nil
}

func (n *Node) createTopicReplica(ctx context.Context, topic string) (orbitdb.EventLogStore, error) {
	n.log.Debug("creating topic", zap.String("topic", topic))
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	if _, ok := n.topics[topic]; ok {
		return nil, ErrTopicAlreadyExists
	}

	// Determine the full address so we attempt to open as existing before creating.
	// TODO: namespace the topic?
	addr, err := n.orbit.DetermineAddress(ctx, topic, "eventlog", &iface.DetermineAddressOptions{
		AccessController: accessController,
	})
	if err != nil {
		return nil, err
	}

	// Open the data store, or creating it if necessary.
	replica, err := n.orbit.Log(n.ctx, addr.String(), &iface.CreateDBOptions{
		AccessController: accessController,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		// Load log store data in case there are existing entries.
		// TODO: re-add this
		// err = replica.Load(n.ctx, infinity)
		// if err != nil {
		// 	n.log.Error("error loading topic replica", zap.Error(err))
		// }
	}()

	go func() {
		sub, err := replica.EventBus().Subscribe([]interface{}{
			new(orbitstores.EventReplicateProgress),
			new(orbitstores.EventReplicated),
		}, eventbus.BufSize(infinity+32))
		defer func() {
			// TODO: log error
			_ = sub.Close()
		}()
		if err != nil {
			n.log.Error("error streaming from topic replica", zap.Error(err))
			// TODO: retry?
		}
		for {
			select {
			case <-n.ctx.Done():
				n.log.Debug("replica stream node context closed")
				return
			case obj, ok := <-sub.Out():
				if !ok {
					return
				}
				switch ev := obj.(type) {
				case orbitstores.EventReplicateProgress:
				// n.log.Debug("replication progress event", zap.Any("event", ev))
				case orbitstores.EventReplicated:
					fmt.Println("START", ev)
					for _, entry := range ev.Entries {
						fmt.Println("HERE", entry.GetHash(), string(entry.GetPayload()))
						op, err := operation.ParseOperation(entry)
						if err != nil {
							n.log.Error("error parsing replicated operation", zap.Error(err))
							continue
						}
						envB := op.GetValue()

						var env messagev1.Envelope
						err = proto.Unmarshal(envB, &env)
						if err != nil {
							n.log.Error("error unmarshaling replicated event", zap.Error(err), zap.String("event", string(envB)))
							continue
						}

						err = n.nc.Publish(env.ContentTopic, envB)
						if err != nil {
							n.log.Error("error publishing replicated event")
						}

						n.log.Debug("received replicated envelope", zap.String("env_topic", env.ContentTopic), zap.String("env_message", string(env.Message)), zap.Int("env_timestamp", int(env.TimestampNs)))
					}
				}
			}
		}
	}()

	n.topics[topic] = replica
	return replica, nil
}

func newIPFSRepo(ctx context.Context, opts *P2POptions) (repo.Repo, error) {
	c := ipfsconfig.Config{}

	// TODO: DRY up between the different node types
	privKey, err := getOrCreatePrivateKey(opts.NodeKey)
	if err != nil {
		return nil, err
	}

	pid, err := peer.IDFromPublicKey(privKey.GetPublic())
	if err != nil {
		return nil, err
	}

	privKeyB, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	c.Pubsub.Enabled = ipfsconfig.True
	// c.Swarm.ResourceMgr.Enabled = cfg.False
	c.Bootstrap = []string{}
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", opts.Port)
	c.Addresses.Swarm = []string{
		listenAddr,
		listenAddr + "/quic",
	}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privKeyB)

	return &repo.Mock{
		D: dsync.MutexWrap(ds.NewMapDatastore()),
		C: c,
	}, nil
}

func getOrCreatePrivateKey(key string) (crypto.PrivKey, error) {
	if key == "" {
		priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 1)
		if err != nil {
			return nil, err
		}

		return priv, nil
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, errors.Wrap(err, "decoding private key")
	}
	return crypto.UnmarshalPrivateKey(keyBytes)
}

func privateKeyToHex(key crypto.PrivKey) (string, error) {
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}

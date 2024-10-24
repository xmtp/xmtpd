package sync

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	clientInterceptors "github.com/xmtp/xmtpd/pkg/interceptors/client"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type syncWorker struct {
	ctx                context.Context
	log                *zap.Logger
	nodeRegistry       registry.NodeRegistry
	registrant         *registrant.Registrant
	store              *sql.DB
	wg                 sync.WaitGroup
	subscriptions      map[uint32]struct{}
	subscriptionsMutex sync.RWMutex
}

type ExitLoopError struct {
	Message string
}

func (e *ExitLoopError) Error() string {
	return e.Message
}

func startSyncWorker(
	ctx context.Context,
	log *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	registrant *registrant.Registrant,
	store *sql.DB,
) (*syncWorker, error) {
	s := &syncWorker{
		ctx:           ctx,
		log:           log.Named("syncWorker"),
		nodeRegistry:  nodeRegistry,
		registrant:    registrant,
		store:         store,
		wg:            sync.WaitGroup{},
		subscriptions: make(map[uint32]struct{}),
	}
	if err := s.start(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *syncWorker) start() error {
	// NOTE: subscriptions can be internally de-duplicated
	// to avoid race conditions, we first set up the listener for new nodes and then all the existing ones
	s.subscribeToRegistry()

	nodes, err := s.nodeRegistry.GetNodes()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		s.subscribeToNode(node.NodeID)
	}

	return nil
}

func (s *syncWorker) close() {
	s.log.Debug("Closing sync worker")
	s.wg.Wait()
	s.log.Debug("Closed sync worker")
}

func (s *syncWorker) subscribeToRegistry() {
	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"node-registry-listener",
		func(ctx context.Context) {
			newNodesCh, cancelNewNodes := s.nodeRegistry.OnNewNodes()
			defer cancelNewNodes()
			for {
				select {
				case <-ctx.Done():
					return
				case newNodes, ok := <-newNodesCh:
					if !ok {
						// data channel closed
						return
					}
					s.log.Info("New nodes received:", zap.Any("nodes", newNodes))
					for _, node := range newNodes {
						s.subscribeToNode(node.NodeID)
					}
				}
			}

		})
}

func (s *syncWorker) subscribeToNode(nodeid uint32) {
	if nodeid == s.registrant.NodeID() {
		return
	}

	s.subscriptionsMutex.Lock()
	defer s.subscriptionsMutex.Unlock()

	if _, exists := s.subscriptions[nodeid]; exists {
		// we already have a subscription to this node
		return
	}

	s.subscriptions[nodeid] = struct{}{}

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d", nodeid),
		func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					s.subscribeToNodeInner(ctx, nodeid)
				}
			}
		})
}

func (s *syncWorker) subscribeToNodeInner(ctx context.Context, nodeid uint32) {

	notifierCtx, notifierCancel := context.WithCancel(ctx)

	node, err := s.setupNodeRegistration(notifierCancel, nodeid)
	if err != nil {
		return
	}

	if !node.IsHealthy || !node.IsValidConfig {
		s.handleUnhealthyNode(notifierCtx)
		return
	}

	s.handleNodeConnection(notifierCtx, node)
}

func (s *syncWorker) handleNodeConnection(
	notifierCtx context.Context,
	node *registry.Node,
) {
	var conn *grpc.ClientConn
	var stream message_api.ReplicationApi_SubscribeEnvelopesClient
	var err error

	//mkysel we should eventually implement a better backoff strategy
	var backoff = time.Second
	for {
		select {
		case <-notifierCtx.Done():
			// either registry has changed or we are shutting down
			s.log.Debug("Context is done. Closing stream and connection")
			return
		default:
			if err != nil {
				s.log.Error(fmt.Sprintf("Error: %v, retrying...", err))
				time.Sleep(backoff)
				backoff = min(backoff*2, 30*time.Second)
			} else {
				backoff = time.Second
			}

			conn, err = s.connectToNode(*node)
			if err != nil {
				continue
			}
			stream, err = s.setupStream(notifierCtx, *node, conn)
			if err != nil {
				continue
			}
			err = s.listenToStream(stream)
		}
	}
}

func (s *syncWorker) handleUnhealthyNode(notifierCtx context.Context) {
	// keep the goroutine idle
	// this will exit the goroutine during shutdown or if the config changed
	<-notifierCtx.Done()
	s.log.Debug("Node configuration has changed. Closing stream and connection")
}

func (s *syncWorker) setupNodeRegistration(
	notifierCancel context.CancelFunc,
	nodeid uint32,
) (*registry.Node, error) {
	registryChan, cancelSub := s.nodeRegistry.OnChangedNode(nodeid)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d-notifier", nodeid),
		func(ctx context.Context) {
			defer cancelSub()
			select {
			case <-ctx.Done():
				// this indicates that the node is shutting down
				// the notifierCtx should have been shut down already,but it can't hurt to cancel it just in case
				notifierCancel()
			case <-registryChan:
				// this indicates that the registry has changed, and we need to rebuild the connection
				s.log.Info(
					"Node has been updated in the registry, terminating and rebuilding...",
				)
				notifierCancel()
			}
		},
	)

	// the nodeRegistry lock gets release between OnChangedNode and GetNode so we might end up in situation
	// where a change gets processed in between
	// this would lead to an unnecessary rebuild of the connection, which is infrequent and okay

	return s.nodeRegistry.GetNode(nodeid)
}

func (s *syncWorker) connectToNode(node registry.Node) (*grpc.ClientConn, error) {
	s.log.Info(fmt.Sprintf("Attempting to connect to %s...", node.HttpAddress))

	interceptor := clientInterceptors.NewAuthInterceptor(s.registrant.TokenFactory(), node.NodeID)
	dialOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	}
	conn, err := node.BuildClient(dialOpts...)
	if err != nil {
		s.log.Error(
			"Failed to connect to peer",
			zap.String("peer", node.HttpAddress),
			zap.Error(err),
		)
		return nil, fmt.Errorf("Failed to connect to peer at %s: %v", node.HttpAddress, err)
	}

	s.log.Debug(fmt.Sprintf("Successfully connected to peer at %s", node.HttpAddress))
	return conn, nil
}

func (s *syncWorker) setupStream(
	ctx context.Context,
	node registry.Node,
	conn *grpc.ClientConn,
) (message_api.ReplicationApi_SubscribeEnvelopesClient, error) {
	result, err := queries.New(s.store).SelectVectorClock(ctx)
	if err != nil {
		return nil, err
	}
	vc := db.ToVectorClock(result)
	s.log.Info(
		"Vector clock for sync subscription",
		zap.Any("nodeID", node.NodeID),
		zap.Any("vc", vc),
	)
	client := message_api.NewReplicationApiClient(conn)
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{node.NodeID},
				LastSeen: &envelopes.VectorClock{
					NodeIdToSequenceId: vc,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to batch subscribe to peer: %v",
			err,
		)
	}
	return stream, nil
}

func (s *syncWorker) listenToStream(
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
) error {
	for {
		// Recv() is a blocking operation that can only be interrupted by cancelling ctx
		envs, err := stream.Recv()
		if err == io.EOF {
			return fmt.Errorf("Stream closed with EOF")
		}
		if err != nil {
			return fmt.Errorf(
				"Stream closed with error: %v",
				err)
		}
		for _, env := range envs.Envelopes {
			s.insertEnvelope(env)
		}
	}

}

func (s *syncWorker) insertEnvelope(env *envelopes.OriginatorEnvelope) {
	s.log.Debug("Replication server received envelope", zap.Any("envelope", env))
	// TODO(nm) Validation logic - share code with API service and publish worker
	originatorBytes, err := proto.Marshal(env)
	if err != nil {
		s.log.Error("Failed to marshal originator envelope", zap.Error(err))
		return
	}

	unsignedEnvelope := &envelopes.UnsignedOriginatorEnvelope{}
	err = proto.Unmarshal(env.GetUnsignedOriginatorEnvelope(), unsignedEnvelope)
	if err != nil {
		s.log.Error(
			"Failed to unmarshal unsigned originator envelope",
			zap.Error(err),
		)
		return
	}

	clientEnvelope := &envelopes.ClientEnvelope{}
	err = proto.Unmarshal(
		unsignedEnvelope.GetPayerEnvelope().GetUnsignedClientEnvelope(),
		clientEnvelope,
	)
	if err != nil {
		s.log.Error(
			"Failed to unmarshal client envelope",
			zap.Error(err),
		)
		return
	}

	q := queries.New(s.store)

	inserted, err := q.InsertGatewayEnvelope(
		s.ctx,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID: int32(unsignedEnvelope.GetOriginatorNodeId()),
			OriginatorSequenceID: int64(
				unsignedEnvelope.GetOriginatorSequenceId(),
			),
			Topic:              clientEnvelope.GetAad().GetTargetTopic(),
			OriginatorEnvelope: originatorBytes,
		},
	)
	if err != nil {
		s.log.Error("Failed to insert gateway envelope", zap.Error(err))
		return
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		s.log.Warn("Envelope already inserted")
		return
	}
}

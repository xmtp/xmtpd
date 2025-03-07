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
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	clientInterceptors "github.com/xmtp/xmtpd/pkg/interceptors/client"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	cancel             context.CancelFunc
}

type originatorStream struct {
	nodeID       uint32
	lastEnvelope *envUtils.OriginatorEnvelope
	stream       message_api.ReplicationApi_SubscribeEnvelopesClient
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

	ctx, cancel := context.WithCancel(ctx)

	s := &syncWorker{
		ctx:           ctx,
		log:           log.Named("syncWorker"),
		nodeRegistry:  nodeRegistry,
		registrant:    registrant,
		store:         store,
		wg:            sync.WaitGroup{},
		subscriptions: make(map[uint32]struct{}),
		cancel:        cancel,
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
	s.cancel()
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
					config := s.setupNodeRegistration(ctx, nodeid)
					s.subscribeToNodeRegistration(config)
				}
			}
		})
}

func (s *syncWorker) subscribeToNodeRegistration(
	registration NodeRegistration,
) {

	node, err := s.nodeRegistry.GetNode(registration.nodeid)
	if err != nil {
		// this should never happen
		s.log.Error(
			"Unexpected state: Failed to get node from registry",
			zap.Uint32("nodeid", registration.nodeid),
			zap.Error(err),
		)
		s.handleUnhealthyNode(registration)
		return
	}

	if !node.IsValidConfig {
		s.handleUnhealthyNode(registration)
		return
	}

	err = nil

	// TODO(mkysel) we should eventually implement a better backoff strategy
	var backoff = time.Second
	for {
		select {
		case <-registration.ctx.Done():
			// either registry has changed or we are shutting down
			s.log.Debug(
				"Context is done. Closing stream and connection",
				zap.String("address", node.HttpAddress),
			)
			return
		default:
			if err != nil {
				s.log.Error(
					"Error connecting to node. Retrying...",
					zap.String("address", node.HttpAddress),
					zap.Error(err),
				)
				time.Sleep(backoff)
				backoff = min(backoff*2, 30*time.Second)
			} else {
				backoff = time.Second
			}

			var conn *grpc.ClientConn
			conn, err = s.connectToNode(*node)
			if err != nil {
				continue
			}

			var stream *originatorStream
			stream, err = s.setupStream(registration.ctx, *node, conn)
			if err != nil {
				_ = conn.Close()
				continue
			}
			err = s.listenToStream(registration.ctx, *node, stream)
			_ = stream.stream.CloseSend()
			_ = conn.Close()
		}
	}
}

func (s *syncWorker) handleUnhealthyNode(registration NodeRegistration) {
	// keep the goroutine idle
	// this will exit the goroutine during shutdown or if the config changed
	<-registration.ctx.Done()
	s.log.Debug("Node configuration has changed. Closing stream and connection")
}

type NodeRegistration struct {
	ctx    context.Context
	cancel context.CancelFunc
	nodeid uint32
}

func (s *syncWorker) setupNodeRegistration(
	ctx context.Context,
	nodeid uint32,
) NodeRegistration {
	notifierCtx, notifierCancel := context.WithCancel(ctx)
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

	return NodeRegistration{ctx: notifierCtx, cancel: notifierCancel, nodeid: nodeid}
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
		return nil, fmt.Errorf("failed to connect to peer at %s: %v", node.HttpAddress, err)
	}

	s.log.Debug(fmt.Sprintf("Successfully opened a connection to peer at %s", node.HttpAddress))
	return conn, nil
}

func (s *syncWorker) setupStream(
	ctx context.Context,
	node registry.Node,
	conn *grpc.ClientConn,
) (*originatorStream, error) {
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
	nodeID := node.NodeID
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{nodeID},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: vc,
				},
			},
		},
	)
	if err != nil {
		s.log.Error(
			"Failed to batch subscribe to peer",
			zap.String("peer", node.HttpAddress),
			zap.Error(err),
		)
		return nil, fmt.Errorf(
			"failed to batch subscribe to peer at %s: %v",
			node.HttpAddress,
			err,
		)
	}
	originatorStream := &originatorStream{nodeID: nodeID, stream: stream}
	for _, row := range result {
		if uint32(row.OriginatorNodeID) == nodeID {
			lastEnvelope, err := envUtils.NewOriginatorEnvelopeFromBytes(row.OriginatorEnvelope)
			if err != nil {
				return nil, err
			}
			originatorStream.lastEnvelope = lastEnvelope
		}
	}
	return originatorStream, nil
}

func (s *syncWorker) listenToStream(
	_ context.Context,
	node registry.Node,
	originatorStream *originatorStream,
) error {
	recvChan := make(chan *message_api.SubscribeEnvelopesResponse)
	errChan := make(chan error)

	go func() {
		for {
			envs, err := originatorStream.stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			recvChan <- envs
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("Context canceled, stopping stream listener")
			return nil

		case envs := <-recvChan:
			s.log.Debug(
				"Received envelopes",
				zap.String("peer", node.HttpAddress),
				zap.Any("numEnvelopes", len(envs.Envelopes)),
			)
			for _, env := range envs.Envelopes {
				s.validateAndInsertEnvelope(originatorStream, env)
			}

		case err := <-errChan:
			if err == io.EOF {
				s.log.Info("Stream closed with EOF")
				// let the caller rebuild the stream if required
				return nil
			}
			s.log.Error(
				"Stream closed with error",
				zap.String("peer", node.HttpAddress),
				zap.Error(err),
			)
			return err
		}
	}
}

func (s *syncWorker) validateAndInsertEnvelope(
	stream *originatorStream,
	envProto *envelopes.OriginatorEnvelope,
) {
	env, err := envUtils.NewOriginatorEnvelope(envProto)
	if err != nil {
		s.log.Error("Failed to unmarshal originator envelope", zap.Error(err))
		return
	}

	if env.OriginatorNodeID() != stream.nodeID {
		s.log.Error("Received envelope from wrong node", zap.Any("nodeID", env.OriginatorNodeID()))
		return
	}

	var lastSequenceID uint64 = 0
	var lastNs int64 = 0
	if stream.lastEnvelope != nil {
		lastSequenceID = stream.lastEnvelope.OriginatorSequenceID()
		lastNs = stream.lastEnvelope.OriginatorNs()
	}
	if env.OriginatorSequenceID() != lastSequenceID+1 || env.OriginatorNs() < lastNs {
		// TODO(rich) Submit misbehavior report and continue
		s.log.Error("Received out of order envelope")
	}

	if env.OriginatorSequenceID() > lastSequenceID {
		stream.lastEnvelope = env
	}

	// TODO Validation logic - share code with API service and publish worker
	// Signatures, topic type, etc
	s.insertEnvelope(env)
}

func (s *syncWorker) insertEnvelope(env *envUtils.OriginatorEnvelope) {
	s.log.Debug("Replication server received envelope", zap.Any("envelope", env))
	originatorBytes, err := env.Bytes()
	if err != nil {
		s.log.Error("Failed to marshal originator envelope", zap.Error(err))
		return
	}

	payerId, err := s.getPayerID(env)
	if err != nil {
		s.log.Error("Failed to get payer ID", zap.Error(err))
		return
	}

	originatorID := int32(env.OriginatorNodeID())
	originatorTime := utils.NsToDate(env.OriginatorNs())

	inserted, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		s.ctx,
		s.store,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(env.OriginatorNodeID()),
			OriginatorSequenceID: int64(env.OriginatorSequenceID()),
			Topic:                env.TargetTopic().Bytes(),
			OriginatorEnvelope:   originatorBytes,
			PayerID:              db.NullInt32(payerId),
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerId,
			OriginatorID:      originatorID,
			MinutesSinceEpoch: utils.MinutesSinceEpoch(originatorTime),
			// TODO:(nm) Independently calculate fees
			SpendPicodollars: int64(
				env.UnsignedOriginatorEnvelope.BaseFee(),
			) + int64(
				env.UnsignedOriginatorEnvelope.CongestionFee(),
			),
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

func (s *syncWorker) getPayerID(env *envUtils.OriginatorEnvelope) (int32, error) {
	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		return 0, err
	}

	q := queries.New(s.store)
	payerId, err := q.FindOrCreatePayer(s.ctx, payerAddress.Hex())
	if err != nil {
		return 0, err
	}

	return payerId, nil
}

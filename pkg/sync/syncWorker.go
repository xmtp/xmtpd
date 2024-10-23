package sync

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	clientInterceptors "github.com/xmtp/xmtpd/pkg/interceptors/client"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
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
	subscriptions      map[interface{}]struct{}
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
		subscriptions: make(map[interface{}]struct{}),
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
					err := s.subscribeToNodeInner(ctx, nodeid)
					if err != nil {
						return
					}
					s.log.Debug(fmt.Sprintf("Reloading configuration for %d", nodeid))
				}
			}
		})
}

func (s *syncWorker) subscribeToNodeInner(ctx context.Context, nodeid uint32) error {

	notifierCtx, notifierCancel := context.WithCancel(ctx)

	registrationRoutine := func(node registry.Node, registryChan <-chan registry.Node, cancelSub registry.CancelSubscription) error {
		if node.NodeID == s.registrant.NodeID() || !node.IsHealthy || !node.IsValidConfig {
			return fmt.Errorf("Invalid Node. No need for subscription")
		}

		tracing.GoPanicWrap(
			s.ctx,
			&s.wg,
			fmt.Sprintf("node-subscribe-%d-notifier", node.NodeID),
			func(ctx context.Context) {
				defer cancelSub()
				select {
				case <-ctx.Done():
					notifierCancel()
				case <-registryChan:
					s.log.Info(
						"Node has been updated in the registry, terminating and rebuilding",
					)
					notifierCancel()
				}
			},
		)

		return nil
	}
	node, err := s.nodeRegistry.RegisterNode(nodeid, registrationRoutine)
	if err != nil {
		return err
	}

	var conn *grpc.ClientConn
	var stream message_api.ReplicationApi_SubscribeEnvelopesClient
	for {
		select {
		case <-notifierCtx.Done():
			s.log.Debug("Node configuration has changed. Closing stream and connection")
			return nil
		default:
			if err != nil {
				s.log.Error(fmt.Sprintf("Error: %v, retrying...", err))
				time.Sleep(1 * time.Second)
			}

			conn, err = s.connectToNode(*node)
			if err != nil {
				continue
			}
			stream, err = s.setupStream(notifierCtx, *node, conn)
			if err != nil {
				continue
			}
			err = s.listenToStream(notifierCtx, stream)
		}
	}
}

func (s *syncWorker) connectToNode(node registry.Node) (*grpc.ClientConn, error) {
	s.log.Info(fmt.Sprintf("Attempting to connect to %s...", node.HttpAddress))
	target, isTLS, err := utils.HttpAddressToGrpcTarget(node.HttpAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert HTTP address to gRPC target: %v", err)
	}
	s.log.Debug(fmt.Sprintf("Mapped %s to %s", node.HttpAddress, target))

	creds, err := utils.GetCredentialsForAddress(isTLS)
	if err != nil {
		return nil, fmt.Errorf("Failed to get credentials: %v", err)
	}

	interceptor := clientInterceptors.NewAuthInterceptor(s.registrant.TokenFactory(), node.NodeID)
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to peer at %s: %v", target, err)
	}
	s.log.Debug(fmt.Sprintf("Successfully connected to peer at %s", target))
	return conn, nil
}

func (s *syncWorker) setupStream(
	ctx context.Context,
	node registry.Node,
	conn *grpc.ClientConn,
) (message_api.ReplicationApi_SubscribeEnvelopesClient, error) {
	client := message_api.NewReplicationApiClient(conn)
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{node.NodeID},
				LastSeen:          nil,
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
	ctx context.Context,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
) error {
	for {
		// Recv is a blocking operation that can only be interrupted by cancelling ctx
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

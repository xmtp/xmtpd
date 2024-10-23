package sync

import (
	"context"
	"database/sql"
	"errors"
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
	ctx          context.Context
	log          *zap.Logger
	nodeRegistry registry.NodeRegistry
	registrant   *registrant.Registrant
	store        *sql.DB
	wg           sync.WaitGroup
	channels     map[uint32]chan struct{}
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
		ctx:          ctx,
		log:          log.Named("syncWorker"),
		nodeRegistry: nodeRegistry,
		registrant:   registrant,
		store:        store,
		wg:           sync.WaitGroup{},
		channels:     make(map[uint32]chan struct{}),
	}
	if err := s.start(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *syncWorker) start() error {
	nodes, err := s.nodeRegistry.GetNodes()
	if err != nil {
		return err
	}
	// NOTE: subscriptions can be internally de-duplicated
	// to avoid race conditions, we first set up the listener for new nodes and then all the existing ones
	s.subscribeToRegistry()

	for _, node := range nodes {
		s.subscribeToNode(node)
	}

	return nil
}

func (s *syncWorker) close() {
	s.log.Info("Closing sync worker")
	for _, ch := range s.channels {
		close(ch)
	}

	//TODO mkysel:stream.Recv() does not seem to get cancelled properly via ctx
	//s.wg.Wait()

	s.log.Info("Closed sync worker")
}

func (s *syncWorker) subscribeToRegistry() {
	newNodesCh, cancelNewNodes := s.nodeRegistry.OnNewNodes()
	go func() {
		defer cancelNewNodes() // Ensure to clean up resources when done
		for newNodes := range newNodesCh {
			s.log.Info("New nodes received:", zap.Any("nodes", newNodes))
			for _, node := range newNodes {
				s.subscribeToNode(node)
			}
		}
	}()
}

func (s *syncWorker) subscribeToNode(node registry.Node) {
	if node.NodeID == s.registrant.NodeID() || !node.IsHealthy || !node.IsValidConfig {
		return
	}
	if _, exists := s.channels[node.NodeID]; exists {
		// we already have a subscription to this node
		return
	}

	cancel := make(chan struct{})
	s.channels[node.NodeID] = cancel

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d", node.NodeID),
		func(ctx context.Context) {
			var err error
			var conn *grpc.ClientConn
			var stream message_api.ReplicationApi_SubscribeEnvelopesClient
			for {
				select {
				case <-cancel:
					s.log.Info(
						fmt.Sprintf(
							"Received cancel signal, exiting loop for node %d",
							node.NodeID,
						),
					)
					return
				default:
					if err != nil {
						var exitError *ExitLoopError
						if errors.As(err, &exitError) {
							s.log.Info("Terminating listener")
							return
						}
						//s.log.Error(fmt.Sprintf("Error: %v, retrying...", err))
						time.Sleep(1 * time.Second)
					}

					conn, err = s.connectToNode(node)
					if err != nil {
						continue
					}
					stream, err = s.setupStream(node, conn)
					if err != nil {
						continue
					}
					err = s.listenToStream(stream, cancel)
					if err != nil {
						continue
					}
				}

			}
		},
	)
}

func (s *syncWorker) connectToNode(node registry.Node) (*grpc.ClientConn, error) {
	//s.log.Info(fmt.Sprintf("Attempting to connect to %s", node.HttpAddress))
	target, isTLS, err := utils.HttpAddressToGrpcTarget(node.HttpAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert HTTP address to gRPC target: %v", err)
	}
	//s.log.Info(fmt.Sprintf("Mapped %s to %s", node.HttpAddress, target))

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
	//s.log.Info(fmt.Sprintf("Successfully connected to peer at %s", target))
	return conn, nil
}

func (s *syncWorker) setupStream(
	node registry.Node,
	conn *grpc.ClientConn,
) (message_api.ReplicationApi_SubscribeEnvelopesClient, error) {
	client := message_api.NewReplicationApiClient(conn)
	stream, err := client.SubscribeEnvelopes(
		s.ctx,
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
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	cancel chan struct{},
) error {
	for {
		select {
		case <-cancel:
			s.log.Info("Received cancel signal, terminating stream listener")
			return &ExitLoopError{}
		default:
			envs, err := stream.Recv()
			if err == io.EOF {
				return nil
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

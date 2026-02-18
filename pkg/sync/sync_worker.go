package sync

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum/common"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/xmtp/xmtpd/pkg/db"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	clientInterceptors "github.com/xmtp/xmtpd/pkg/interceptors/client"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/payerreport"
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
	ctx                        context.Context
	logger                     *zap.Logger
	nodeRegistry               registry.NodeRegistry
	registrant                 *registrant.Registrant
	store                      *db.Handler
	wg                         sync.WaitGroup
	subscriptions              map[uint32]struct{}
	subscriptionsMutex         sync.RWMutex
	cancel                     context.CancelFunc
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
	migration                  MigrationConfig
	clientMetrics              *grpcprom.ClientMetrics
}

func startSyncWorker(
	cfg *SyncServerConfig,
) (*syncWorker, error) {
	ctx, cancel := context.WithCancel(cfg.Ctx)

	s := &syncWorker{
		ctx:                        ctx,
		logger:                     cfg.Logger.Named(utils.SyncWorkerName),
		nodeRegistry:               cfg.NodeRegistry,
		registrant:                 cfg.Registrant,
		store:                      cfg.DB,
		feeCalculator:              cfg.FeeCalculator,
		subscriptions:              make(map[uint32]struct{}),
		payerReportStore:           cfg.PayerReportStore,
		payerReportDomainSeparator: cfg.PayerReportDomainSeparator,
		migration:                  cfg.Migration,
		cancel:                     cancel,
		clientMetrics:              cfg.ClientMetrics,
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

	if s.migration.Enable {
		s.logger.Info(
			"Migration client is enabled. Will migrate from migration originator",
			utils.OriginatorIDField(s.migration.FromNodeID),
		)
	}

	for _, node := range nodes {
		s.subscribeToNode(node.NodeID)
	}

	return nil
}

func (s *syncWorker) close() {
	s.logger.Debug("closing")
	s.cancel()
	s.wg.Wait()
	s.logger.Debug("closed")
}

func (s *syncWorker) subscribeToRegistry() {
	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"node-registry-listener",
		func(ctx context.Context) {
			newNodesCh := s.nodeRegistry.OnNewNodes()
			for {
				select {
				case <-ctx.Done():
					return
				case newNodes, ok := <-newNodesCh:
					if !ok {
						// data channel closed
						return
					}

					if s.logger.Core().Enabled(zap.DebugLevel) {
						s.logger.Debug("new nodes received", utils.BodyField(newNodes))
					}

					for _, node := range newNodes {
						s.subscribeToNode(node.NodeID)
					}
				}
			}
		})
}

func (s *syncWorker) subscribeToNode(nodeID uint32) {
	if nodeID == s.registrant.NodeID() {
		return
	}

	s.subscriptionsMutex.Lock()
	defer s.subscriptionsMutex.Unlock()

	if _, exists := s.subscriptions[nodeID]; exists {
		// we already have a subscription to this node
		return
	}

	s.subscriptions[nodeID] = struct{}{}

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d-db", nodeID),
		func(ctx context.Context) {
			newEnvelopeSink(
				ctx,
				s.store,
				s.logger,
				s.feeCalculator,
				s.payerReportStore,
				s.payerReportDomainSeparator,
				writeQueue,
				1*time.Second,
			).Start()
		})

	changeListener := NewNodeRegistryWatcher(s.ctx, s.logger, nodeID, s.nodeRegistry)
	changeListener.Watch()

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d", nodeID),
		func(ctx context.Context) {
			defer close(writeQueue)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					notifierCtx, notifierCancel := context.WithCancel(ctx)
					changeListener.RegisterCancelFunction(notifierCancel)
					s.subscribeToNodeRegistration(
						NodeRegistration{
							ctx:    notifierCtx,
							cancel: notifierCancel,
							nodeID: nodeID,
						},
						writeQueue,
					)
				}
			}
		})
}

func (s *syncWorker) subscribeToNodeRegistration(
	registration NodeRegistration,
	writeQueue chan *envUtils.OriginatorEnvelope,
) {
	connectionsStatusCounter := metrics.NewSyncConnectionsStatusCounter(registration.nodeID)
	defer connectionsStatusCounter.Close()

	node, err := s.nodeRegistry.GetNode(registration.nodeID)
	if err != nil {
		// this should never happen
		s.logger.Error(
			"unexpected state: failed to get node from registry",
			utils.OriginatorIDField(registration.nodeID),
			zap.Error(err),
		)
		connectionsStatusCounter.MarkFailure()
		s.handleUnhealthyNode(registration)
		return
	}

	if !node.IsValidConfig {
		connectionsStatusCounter.MarkFailure()
		s.handleUnhealthyNode(registration)
		return
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 1 * time.Second

	operation := func() (string, error) {
		var (
			conn   *grpc.ClientConn
			stream *originatorStream
			err    error
		)

		defer func() {
			if stream != nil {
				_ = stream.stream.CloseSend()
			}
			if conn != nil {
				_ = conn.Close()
			}
		}()

		defer func() {
			if err != nil && s.ctx.Err() == nil {
				s.logger.Error(
					"error connecting to node, retrying",
					utils.NodeHTTPAddressField(node.HTTPAddress),
					zap.Error(err),
				)
				connectionsStatusCounter.MarkFailure()
			}
		}()

		if registration.ctx.Err() != nil {
			return "", backoff.Permanent(registration.ctx.Err())
		}

		conn, err = s.connectToNode(*node)
		if err != nil {
			return "", err
		}

		stream, err = s.setupStream(registration.ctx, *node, conn, writeQueue)
		if err != nil {
			return "", err
		}

		connectionsStatusCounter.MarkSuccess()

		err = stream.listen()
		return "", err
	}

	_, _ = backoff.Retry(
		registration.ctx,
		operation,
		backoff.WithBackOff(expBackoff),
		backoff.WithMaxElapsedTime(0),
	)
}

func (s *syncWorker) handleUnhealthyNode(registration NodeRegistration) {
	// keep the goroutine idle
	// this will exit the goroutine during shutdown or if the config changed
	<-registration.ctx.Done()
	s.logger.Debug("node configuration has changed, closing stream and connection")
}

type NodeRegistration struct {
	ctx    context.Context
	cancel context.CancelFunc
	nodeID uint32
}

// connectToNode connects to a node and returns a gRPC client connection.
// Note that this is a gRPC native connection, not a Connect-based connection.
// The server side uses Connect-based connections, which supports gRPC as well.
func (s *syncWorker) connectToNode(
	node registry.Node,
) (*grpc.ClientConn, error) {
	s.logger.Info("attempting to connect to node",
		utils.OriginatorIDField(node.NodeID),
		utils.NodeHTTPAddressField(node.HTTPAddress),
	)

	interceptor := clientInterceptors.NewClientAuthInterceptor(
		s.registrant.TokenFactory(),
		node.NodeID,
	)

	// Execute first the auth interceptor, then metrics.
	dialOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			interceptor.Unary(),
			s.clientMetrics.UnaryClientInterceptor(),
		),
		grpc.WithChainStreamInterceptor(
			interceptor.Stream(),
			s.clientMetrics.StreamClientInterceptor(),
		),
	}

	conn, err := node.BuildConn(dialOpts...)
	if err != nil {
		s.logger.Error(
			"failed to connect to node",
			utils.OriginatorIDField(node.NodeID),
			utils.NodeHTTPAddressField(node.HTTPAddress),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to connect to node at %s: %w", node.HTTPAddress, err)
	}

	s.logger.Debug("successfully opened a connection to node",
		utils.OriginatorIDField(node.NodeID),
		utils.NodeHTTPAddressField(node.HTTPAddress),
	)

	return conn, nil
}

func (s *syncWorker) setupStream(
	ctx context.Context,
	node registry.Node,
	conn *grpc.ClientConn,
	writeQueue chan *envUtils.OriginatorEnvelope,
) (*originatorStream, error) {
	result, err := s.store.ReadQuery().SelectVectorClock(ctx)
	if err != nil {
		return nil, err
	}

	var (
		client            = message_api.NewReplicationApiClient(conn)
		vc                = db.ToVectorClock(result)
		localNodeID       = s.registrant.NodeID()
		syncNodeID        = node.NodeID
		migratorNodeID    = s.migration.FromNodeID
		originatorNodeIDs = []uint32{syncNodeID}
	)

	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug(
			"vector clock for sync subscription",
			utils.OriginatorIDField(node.NodeID),
			utils.BodyField(vc),
		)
	}

	if s.migration.Enable && syncNodeID == migratorNodeID && migratorNodeID != localNodeID {
		originatorNodeIDs = []uint32{
			syncNodeID,
			migrator.GroupMessageOriginatorID,
			migrator.WelcomeMessageOriginatorID,
			migrator.KeyPackagesOriginatorID,
		}
		if s.logger.Core().Enabled(zap.DebugLevel) {
			s.logger.Debug(
				"requesting additional migrated payloads from originator node",
				utils.OriginatorIDField(syncNodeID),
				zap.Any("originators", originatorNodeIDs),
			)
		}
	}

	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: originatorNodeIDs,
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: vc,
				},
			},
		},
	)
	if err != nil {
		s.logger.Error(
			"failed to batch subscribe to node",
			utils.OriginatorIDField(node.NodeID),
			utils.NodeHTTPAddressField(node.HTTPAddress),
			zap.Error(err),
		)
		return nil, fmt.Errorf(
			"failed to batch subscribe to peer at %s: %w",
			node.HTTPAddress,
			err,
		)
	}

	lastSequenceID := uint64(0)
	for _, row := range result {
		if slices.Contains(originatorNodeIDs, uint32(row.OriginatorNodeID)) {
			lastSequenceID = uint64(row.OriginatorSequenceID)
		}
	}

	permittedOriginators := utils.SliceToSet(originatorNodeIDs)

	return newOriginatorStream(
		s.ctx,
		s.logger,
		&node,
		lastSequenceID,
		permittedOriginators,
		stream,
		writeQueue,
	), nil
}

package sync

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
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
	store                      *sql.DB
	wg                         sync.WaitGroup
	subscriptions              map[uint32]struct{}
	subscriptionsMutex         sync.RWMutex
	cancel                     context.CancelFunc
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
	migration                  MigrationConfig
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
	s.logger.Debug("stopping")
	s.cancel()
	s.wg.Wait()
	s.logger.Debug("stopped")
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
					s.logger.Info("new nodes received", zap.Any("nodes", newNodes))
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

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d-db", nodeid),
		func(ctx context.Context) {
			newEnvelopeSink(
				ctx,
				s.store,
				s.logger,
				s.feeCalculator,
				s.payerReportStore,
				s.payerReportDomainSeparator,
				writeQueue,
			).Start()
		})

	changeListener := NewNodeRegistryWatcher(s.ctx, s.logger, nodeid, s.nodeRegistry)
	changeListener.Watch()

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d", nodeid),
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
							nodeid: nodeid,
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
	connectionsStatusCounter := metrics.NewSyncConnectionsStatusCounter(registration.nodeid)
	defer connectionsStatusCounter.Close()

	node, err := s.nodeRegistry.GetNode(registration.nodeid)
	if err != nil {
		// this should never happen
		s.logger.Error(
			"unexpected state: failed to get node from registry",
			utils.OriginatorIDField(registration.nodeid),
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
		// Ensure cleanup of resources, defer works here since we are using a named function
		var conn *grpc.ClientConn
		var stream *originatorStream
		defer func() {
			if stream != nil {
				_ = stream.stream.CloseSend()
			}
			if conn != nil {
				_ = conn.Close()
			}
		}()

		var err error
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
	nodeid uint32
}

func (s *syncWorker) connectToNode(node registry.Node) (*grpc.ClientConn, error) {
	s.logger.Info("attempting to connect to node",
		utils.OriginatorIDField(node.NodeID),
		utils.NodeHTTPAddressField(node.HTTPAddress),
	)

	interceptor := clientInterceptors.NewAuthInterceptor(s.registrant.TokenFactory(), node.NodeID)
	dialOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	}
	conn, err := node.BuildClient(dialOpts...)
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
	result, err := queries.New(s.store).SelectVectorClock(ctx)
	if err != nil {
		return nil, err
	}

	var (
		vc                = db.ToVectorClock(result)
		client            = message_api.NewReplicationApiClient(conn)
		nodeID            = node.NodeID
		originatorNodeIDs = []uint32{nodeID}
	)

	s.logger.Info(
		"vector clock for sync subscription",
		utils.OriginatorIDField(node.NodeID),
		zap.Any("vector_clock", vc),
	)

	if s.migration.Enable && nodeID == s.migration.FromNodeID {
		originatorNodeIDs = []uint32{
			nodeID,
			migrator.GroupMessageOriginatorID,
			migrator.WelcomeMessageOriginatorID,
			migrator.KeyPackagesOriginatorID,
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

	var c *cursor
	for _, row := range result {
		if slices.Contains(originatorNodeIDs, uint32(row.OriginatorNodeID)) {
			c = &cursor{
				sequenceID:  uint64(row.OriginatorSequenceID),
				timestampNS: row.GatewayTime.UnixNano(),
			}
		}
	}

	return newOriginatorStream(
		s.ctx,
		s.logger,
		&node,
		c,
		stream,
		writeQueue,
	), nil
}

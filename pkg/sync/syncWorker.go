package sync

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type syncWorker struct {
	ctx                        context.Context
	log                        *zap.Logger
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
}

func startSyncWorker(
	cfg *SyncServerConfig,
) (*syncWorker, error) {
	ctx, cancel := context.WithCancel(cfg.Ctx)

	s := &syncWorker{
		ctx:                        ctx,
		log:                        cfg.Log.Named("syncWorker"),
		nodeRegistry:               cfg.NodeRegistry,
		registrant:                 cfg.Registrant,
		store:                      cfg.DB,
		feeCalculator:              cfg.FeeCalculator,
		subscriptions:              make(map[uint32]struct{}),
		payerReportStore:           cfg.PayerReportStore,
		payerReportDomainSeparator: cfg.PayerReportDomainSeparator,
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

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d-db", nodeid),
		func(ctx context.Context) {
			newEnvelopeSink(
				ctx,
				s.store,
				s.log,
				s.feeCalculator,
				s.payerReportStore,
				s.payerReportDomainSeparator,
				writeQueue,
			).Start()
		})

	changeListener := NewNodeRegistryWatcher(s.ctx, s.log, nodeid, s.nodeRegistry)
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
		s.log.Error(
			"Unexpected state: Failed to get node from registry",
			zap.Uint32("nodeid", registration.nodeid),
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
				s.log.Error(
					"Error connecting to node. Retrying...",
					zap.String("address", node.HttpAddress),
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
	s.log.Debug("Node configuration has changed. Closing stream and connection")
}

type NodeRegistration struct {
	ctx    context.Context
	cancel context.CancelFunc
	nodeid uint32
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
	writeQueue chan *envUtils.OriginatorEnvelope,
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

	var lastEnvelope *envUtils.OriginatorEnvelope
	for _, row := range result {
		if uint32(row.OriginatorNodeID) == nodeID {
			lastEnvelope, err = envUtils.NewOriginatorEnvelopeFromBytes(row.OriginatorEnvelope)
			if err != nil {
				return nil, err
			}
		}
	}

	return newOriginatorStream(
		s.ctx,
		s.log,
		&node,
		lastEnvelope,
		stream,
		writeQueue,
	), nil
}

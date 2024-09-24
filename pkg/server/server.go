package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type ReplicationServer struct {
	apiServer    *api.ApiServer
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   *registrant.Registrant
	nodeRegistry registry.NodeRegistry
	options      config.ServerOptions
	metrics      *metrics.Server
	writerDB     *sql.DB
	Wg           sync.WaitGroup
	// Can add reader DB later if needed
}

func NewReplicationServer(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
	nodeRegistry registry.NodeRegistry,
	writerDB *sql.DB,
) (*ReplicationServer, error) {
	var err error

	var mtcs *metrics.Server
	if options.Metrics.Enable {
		promReg := prometheus.NewRegistry()
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())

		mtcs, err = metrics.NewMetricsServer(ctx,
			options.Metrics.Address,
			options.Metrics.Port,
			log,
			promReg,
		)
		if err != nil {

			log.Error("initializing metrics server", zap.Error(err))
			return nil, err
		}
	}

	s := &ReplicationServer{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
		writerDB:     writerDB,
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	s.registrant, err = registrant.NewRegistrant(
		s.ctx,
		queries.New(s.writerDB),
		nodeRegistry,
		options.Signer.PrivateKey,
	)
	if err != nil {
		return nil, err
	}

	s.apiServer, err = api.NewAPIServer(
		s.ctx,
		s.writerDB,
		log,
		options.API.Port,
		s.registrant,
		options.Reflection.Enable,
	)
	if err != nil {
		return nil, err
	}

	log.Info("Replication server started", zap.Int("port", options.API.Port))

	nodes, err := nodeRegistry.GetNodes()
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if node.NodeID == s.registrant.NodeID() || node.NodeID == 0 {
			continue
		}
		subscribeToNode(node, log, s)
	}
	return s, nil
}

func subscribeToNode(node registry.Node, log *zap.Logger, s *ReplicationServer) {
	tracing.GoPanicWrap(
		s.ctx,
		&s.Wg,
		fmt.Sprintf("node-subscribe-%d", node.NodeID),
		func(ctx context.Context) {
			for {
				addr := node.HttpAddress
				log.Info(fmt.Sprintf("attempting to connect to %s", addr))
				conn, err := s.apiServer.DialGRPC(addr)
				if err != nil {
					time.Sleep(1000 * time.Millisecond)
					log.Info("Replication server failed to connect to peer. Retrying...")
					continue
				}
				client := message_api.NewReplicationApiClient(conn)
				stream, err := client.BatchSubscribeEnvelopes(
					s.ctx,
					&message_api.BatchSubscribeEnvelopesRequest{
						Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
							{
								Query: &message_api.EnvelopesQuery{
									Filter: &message_api.EnvelopesQuery_OriginatorNodeId{
										OriginatorNodeId: node.NodeID,
									},
									LastSeen: nil,
								},
							},
						},
					},
				)
				if err != nil {
					time.Sleep(1000 * time.Millisecond)
					log.Info(fmt.Sprintf(
						"Replication server failed to batch subscribe to peer. Retrying... %v",
						err),
					)
					continue
				}

				log.Info(fmt.Sprintf("Successfully connected to peer at %s", addr))
				for {
					envs, err := stream.Recv()
					if err != nil {
						log.Info(fmt.Sprintf(
							"Replication server subscription closed. Retrying... %v",
							err),
						)
						break
					}
					for _, env := range envs.Envelopes {
						log.Info(fmt.Sprintf("Replication server received envelope %s", env))
					}
				}
			}
		},
	)
}

func (s *ReplicationServer) Addr() net.Addr {
	return s.apiServer.Addr()
}

func (s *ReplicationServer) WaitForShutdown() {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(
		termChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	<-termChannel
	s.Shutdown()
}

func (s *ReplicationServer) Shutdown() {
	// Close metrics server.
	if s.metrics != nil {
		if err := s.metrics.Close(); err != nil {
			s.log.Error("stopping metrics", zap.Error(err))
		}
	}

	if s.apiServer != nil {
		s.apiServer.Close()
	}
	s.cancel()
}

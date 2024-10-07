package sync

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type syncWorker struct {
	ctx          context.Context
	log          *zap.Logger
	nodeRegistry registry.NodeRegistry
	registrant   *registrant.Registrant
	store        *sql.DB
	wg           sync.WaitGroup
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
		log:          log.With(zap.String("method", "syncWorker")),
		nodeRegistry: nodeRegistry,
		registrant:   registrant,
		store:        store,
		wg:           sync.WaitGroup{},
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
	for _, node := range nodes {
		if node.NodeID == s.registrant.NodeID() {
			continue
		}
		s.subscribeToNode(node)
	}
	return nil
}

func (s *syncWorker) close() {
	s.wg.Wait()
}

func (s *syncWorker) subscribeToNode(node registry.Node) {
	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("node-subscribe-%d", node.NodeID),
		func(ctx context.Context) {
			for {
				addr := node.HttpAddress
				log.Info(fmt.Sprintf("attempting to connect to %s", addr))
				conn, err := grpc.NewClient(
					addr,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithDefaultCallOptions(),
				)
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
						originatorBytes, err := proto.Marshal(env)
						if err != nil {
							log.Error("Failed to marshal originator envelope", zap.Error(err))
						}

						unsignedEnvelope := &message_api.UnsignedOriginatorEnvelope{}
						err = proto.Unmarshal(env.GetUnsignedOriginatorEnvelope(), unsignedEnvelope)
						if err != nil {
							log.Error(
								"Failed to unmarshal unsigned originator envelope",
								zap.Error(err),
							)
						}

						clientEnvelope := &message_api.ClientEnvelope{}
						err = proto.Unmarshal(
							unsignedEnvelope.GetPayerEnvelope().GetUnsignedClientEnvelope(),
							clientEnvelope,
						)
						if err != nil {
							log.Error(
								"Failed to unmarshal client envelope",
								zap.Error(err),
							)
						}

						q := queries.New(s.store)

						// On unique constraint conflicts, no error is thrown, but numRows is 0
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
							log.Error("Failed to insert gateway envelope", zap.Error(err))
						} else if inserted == 0 {
							// Envelope was already inserted by another worker
							log.Error("Envelope already inserted")
						}
					}
				}
			}
		},
	)
}

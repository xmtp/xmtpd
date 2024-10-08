package api

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

const (
	maxRequestedRows     uint32 = 1000
	maxVectorClockLength int    = 100
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

	ctx              context.Context
	log              *zap.Logger
	registrant       *registrant.Registrant
	store            *sql.DB
	publishWorker    *publishWorker
	subscribeWorker  *subscribeWorker
	messagePublisher blockchain.IBlockchainPublisher
}

func NewReplicationApiService(
	ctx context.Context,
	log *zap.Logger,
	registrant *registrant.Registrant,
	store *sql.DB,
	messagePublisher blockchain.IBlockchainPublisher,

) (*Service, error) {
	publishWorker, err := startPublishWorker(ctx, log, registrant, store)
	if err != nil {
		return nil, err
	}
	subscribeWorker, err := startSubscribeWorker(ctx, log, store)
	if err != nil {
		return nil, err
	}

	return &Service{
		ctx:              ctx,
		log:              log,
		registrant:       registrant,
		store:            store,
		publishWorker:    publishWorker,
		subscribeWorker:  subscribeWorker,
		messagePublisher: messagePublisher,
	}, nil
}

func (s *Service) Close() {
	s.log.Info("closed")
}

func (s *Service) BatchSubscribeEnvelopes(
	req *message_api.BatchSubscribeEnvelopesRequest,
	stream message_api.ReplicationApi_BatchSubscribeEnvelopesServer,
) error {
	log := s.log.With(zap.String("method", "batchSubscribe"))

	// Send a header (any header) to fix an issue with Tonic based GRPC clients.
	// See: https://github.com/xmtp/libxmtp/pull/58
	err := stream.SendHeader(metadata.Pairs("subscribed", "true"))
	if err != nil {
		return status.Errorf(codes.Internal, "could not send header: %v", err)
	}

	requests := req.GetRequests()
	if len(requests) == 0 {
		return status.Errorf(codes.InvalidArgument, "missing requests")
	}

	ch, err := s.subscribeWorker.listen(stream.Context(), requests)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid subscription request: %v", err)
	}

	for {
		select {
		case envs, open := <-ch:
			if open {
				err := stream.Send(&message_api.BatchSubscribeEnvelopesResponse{
					Envelopes: envs,
				})
				if err != nil {
					return status.Errorf(codes.Internal, "error sending envelope: %v", err)
				}
			} else {
				// TODO(rich) Recover from backpressure
				log.Debug("channel closed by worker")
				return nil
			}
		case <-stream.Context().Done():
			log.Debug("stream closed")
			return nil
		case <-s.ctx.Done():
			log.Info("service closed")
			return nil
		}
	}
}

func (s *Service) QueryEnvelopes(
	ctx context.Context,
	req *message_api.QueryEnvelopesRequest,
) (*message_api.QueryEnvelopesResponse, error) {
	log := s.log.With(zap.String("method", "query"))
	params, err := s.queryReqToDBParams(req)
	if err != nil {
		return nil, err
	}

	rows, err := queries.New(s.store).SelectGatewayEnvelopes(ctx, *params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
	}

	envs := make([]*message_api.OriginatorEnvelope, 0, len(rows))
	for _, row := range rows {
		originatorEnv := &message_api.OriginatorEnvelope{}
		err := proto.Unmarshal(row.OriginatorEnvelope, originatorEnv)
		if err != nil {
			// We expect to have already validated the envelope when it was inserted
			log.Error("could not unmarshal originator envelope", zap.Error(err))
			continue
		}
		envs = append(envs, originatorEnv)
	}

	return &message_api.QueryEnvelopesResponse{
		Envelopes: envs,
	}, nil
}

func (s *Service) queryReqToDBParams(
	req *message_api.QueryEnvelopesRequest,
) (*queries.SelectGatewayEnvelopesParams, error) {
	params := queries.SelectGatewayEnvelopesParams{
		Topic:             nil,
		OriginatorNodeID:  sql.NullInt32{},
		RowLimit:          sql.NullInt32{},
		CursorNodeIds:     nil,
		CursorSequenceIds: nil,
	}

	query := req.GetQuery()
	if query == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing query")
	}

	switch filter := query.GetFilter().(type) {
	case *message_api.EnvelopesQuery_Topic:
		if len(filter.Topic) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "missing topic")
		}
		params.Topic = filter.Topic
	case *message_api.EnvelopesQuery_OriginatorNodeId:
		params.OriginatorNodeID = db.NullInt32(int32(filter.OriginatorNodeId))
	default:
	}

	vc := query.GetLastSeen().GetNodeIdToSequenceId()
	if len(vc) > maxVectorClockLength {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"vector clock length exceeds maximum of %d",
			maxVectorClockLength,
		)
	}
	db.SetVectorClock(&params, vc)

	limit := req.GetLimit()
	if limit > 0 && limit <= maxRequestedRows {
		params.RowLimit = db.NullInt32(int32(limit))
	}

	return &params, nil
}

func (s *Service) PublishEnvelope(
	ctx context.Context,
	req *message_api.PublishEnvelopeRequest,
) (*message_api.PublishEnvelopeResponse, error) {
	clientEnv, err := s.validatePayerInfo(req.GetPayerEnvelope())
	if err != nil {
		return nil, err
	}

	topic, err := s.validateClientInfo(clientEnv)
	if err != nil {
		return nil, err
	}

	didPublish, err := s.maybePublishToBlockchain(ctx, clientEnv)
	if err != nil {
		return nil, err
	}
	if didPublish {
		return &message_api.PublishEnvelopeResponse{}, nil
	}

	payerBytes, err := proto.Marshal(req.GetPayerEnvelope())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal envelope: %v", err)
	}

	stagedEnv, err := queries.New(s.store).
		InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         topic,
			PayerEnvelope: payerBytes,
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert staged envelope: %v", err)
	}
	s.publishWorker.notifyStagedPublish()

	originatorEnv, err := s.registrant.SignStagedEnvelope(stagedEnv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not sign envelope: %v", err)
	}

	return &message_api.PublishEnvelopeResponse{OriginatorEnvelope: originatorEnv}, nil
}

func (s *Service) maybePublishToBlockchain(
	ctx context.Context,
	clientEnv *message_api.ClientEnvelope,
) (didPublish bool, err error) {
	payload, ok := clientEnv.GetPayload().(*message_api.ClientEnvelope_IdentityUpdate)
	if ok && payload.IdentityUpdate != nil {
		if err = s.publishIdentityUpdate(ctx, payload.IdentityUpdate); err != nil {
			s.log.Error("could not publish identity update", zap.Error(err))
			return false, status.Errorf(
				codes.Internal,
				"could not publish identity update: %v",
				err,
			)
		}
		return true, nil
	}

	return false, nil
}

func (s *Service) publishIdentityUpdate(
	ctx context.Context,
	identityUpdate *associations.IdentityUpdate,
) error {
	identityUpdateBytes, err := proto.Marshal(identityUpdate)
	if err != nil {
		return err
	}
	inboxId, err := utils.ParseInboxId(identityUpdate.InboxId)
	if err != nil {
		return err
	}
	return s.messagePublisher.PublishIdentityUpdate(
		ctx,
		inboxId,
		identityUpdateBytes,
	)
}

func (s *Service) GetInboxIds(
	ctx context.Context,
	req *message_api.GetInboxIdsRequest,
) (*message_api.GetInboxIdsResponse, error) {
	logger := s.log.With(zap.String("method", "GetInboxIds"))
	queries := queries.New(s.store)
	addresses := []string{}
	for _, request := range req.Requests {
		addresses = append(addresses, request.GetAddress())
	}

	addressLogEntries, err := queries.GetAddressLogs(ctx, addresses)
	if err != nil {
		return nil, err
	}

	out := make([]*message_api.GetInboxIdsResponse_Response, len(addresses))

	for index, address := range addresses {
		resp := message_api.GetInboxIdsResponse_Response{}
		resp.Address = address

		for _, logEntry := range addressLogEntries {
			if logEntry.Address == address {
				inboxId := logEntry.InboxID
				resp.InboxId = &inboxId
			}
		}
		out[index] = &resp
	}

	logger.Info("got inbox ids", zap.Int("numResponses", len(out)))

	return &message_api.GetInboxIdsResponse{
		Responses: out,
	}, nil
}

func (s *Service) validatePayerInfo(
	payerEnv *message_api.PayerEnvelope,
) (*message_api.ClientEnvelope, error) {
	clientBytes := payerEnv.GetUnsignedClientEnvelope()
	sig := payerEnv.GetPayerSignature()
	if clientBytes == nil || sig == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing envelope or signature")
	}
	// TODO(rich): Verify payer signature

	clientEnv := &message_api.ClientEnvelope{}
	err := proto.Unmarshal(clientBytes, clientEnv)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"could not unmarshal client envelope: %v",
			err,
		)
	}

	return clientEnv, nil
}

func (s *Service) validateClientInfo(clientEnv *message_api.ClientEnvelope) ([]byte, error) {
	if clientEnv.GetAad().GetTargetOriginator() != s.registrant.NodeID() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target originator")
	}

	topic := clientEnv.GetAad().GetTargetTopic()
	if len(topic) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing target topic")
	}

	// TODO(rich): Verify all originators have synced past `last_seen`
	// TODO(rich): Check that the blockchain sequence ID is equal to the latest on the group
	// TODO(rich): Perform any payload-specific validation (e.g. identity updates)

	return topic, nil
}

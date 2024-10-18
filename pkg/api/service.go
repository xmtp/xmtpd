package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	maxRequestedRows     int32 = 1000
	maxQueriesPerRequest int   = 10000
	maxTopicLength       int   = 128
	maxVectorClockLength int   = 100
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

// Pulls from DB and sends to client, updating the query's last seen cursor, until
// the stream has caught up to the latest in the database.
func (s *Service) catchUpFromCursor(
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
	query *message_api.EnvelopesQuery,
	logger *zap.Logger,
) error {
	// TODO(rich): Pull one more time after first tick.
	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen = &message_api.VectorClock{
			NodeIdToSequenceId: cursor,
		}
		return nil
	}

	for {
		rows, err := s.fetchEnvelopes(stream.Context(), query, maxRequestedRows)
		if err != nil {
			return err
		}
		payloads := make([]*OriginatorEnvelopeWithInfo, 0, len(rows))
		for _, env := range rows {
			p := &OriginatorEnvelopeWithInfo{
				Envelope:             &message_api.OriginatorEnvelope{},
				OriginatorNodeID:     uint32(env.OriginatorNodeID),
				OriginatorSequenceID: uint64(env.OriginatorSequenceID),
			}
			err := proto.Unmarshal(env.OriginatorEnvelope, p.Envelope)
			if err != nil {
				// We expect to have already validated the envelope when it was inserted
				logger.Error("could not unmarshal originator envelope", zap.Error(err))
				continue
			}
			payloads = append(payloads, p)
		}
		err = s.sendEnvelopes(stream, query, payloads)
		if err != nil {
			return status.Errorf(codes.Internal, "error sending envelopes: %v", err)
		}
		if len(rows) < int(maxRequestedRows) {
			// There were no more envelopes in DB at time of fetch
			break
		}
		// TODO(rich): Determine default rate limits and set query interval to match
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func (s *Service) sendEnvelopes(
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
	query *message_api.EnvelopesQuery,
	payloads []*OriginatorEnvelopeWithInfo,
) error {
	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	for _, p := range payloads {
		if cursor[uint32(p.OriginatorNodeID)] >= p.OriginatorSequenceID {
			continue
		}

		// TODO(rich): Either batch send envelopes, or modify stream proto to
		// send one envelope at a time.
		err := stream.Send(&message_api.SubscribeEnvelopesResponse{
			Envelopes: []*message_api.OriginatorEnvelope{p.Envelope},
		})
		if err != nil {
			return status.Errorf(codes.Internal, "error sending envelope: %v", err)
		}
		cursor[uint32(p.OriginatorNodeID)] = p.OriginatorSequenceID
	}
	return nil
}

func (s *Service) SubscribeEnvelopes(
	req *message_api.SubscribeEnvelopesRequest,
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
) error {
	log := s.log.With(zap.String("method", "subscribe"))

	// Send a header (any header) to fix an issue with Tonic based GRPC clients.
	// See: https://github.com/xmtp/libxmtp/pull/58
	err := stream.SendHeader(metadata.Pairs("subscribed", "true"))
	if err != nil {
		return status.Errorf(codes.Internal, "could not send header: %v", err)
	}

	query := req.GetQuery()
	if err := s.validateQuery(query); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid subscription request: %v", err)
	}

	ch := s.subscribeWorker.listen(stream.Context(), query)
	err = s.catchUpFromCursor(stream, query, log)
	if err != nil {
		return err
	}

	for {
		select {
		case payloads, open := <-ch:
			if open {
				err := s.sendEnvelopes(stream, query, payloads)
				if err != nil {
					return status.Errorf(codes.Internal, "error sending envelope: %v", err)
				}
			} else {
				// TODO(rich) Reset whole subscribe flow from new cursor
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

	limit := int32(req.GetLimit())
	if limit == 0 {
		limit = maxRequestedRows
	}
	rows, err := s.fetchEnvelopes(ctx, req.GetQuery(), limit)
	if err != nil {
		return nil, err
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

func (s *Service) validateQuery(
	query *message_api.EnvelopesQuery,
) error {
	if query == nil {
		return fmt.Errorf("missing query")
	}

	topics := query.GetTopics()
	originators := query.GetOriginatorNodeIds()
	if len(topics) != 0 && len(originators) != 0 {
		return fmt.Errorf(
			"cannot filter by both topic and originator in same subscription request",
		)
	}

	numQueries := len(topics) + len(originators)
	if numQueries > maxQueriesPerRequest {
		return fmt.Errorf(
			"too many subscriptions: %d, consider subscribing to fewer topics or subscribing without a filter",
			numQueries,
		)
	}

	for _, topic := range topics {
		if len(topic) == 0 || len(topic) > maxTopicLength {
			return fmt.Errorf("invalid topic: %s", topic)
		}
	}

	vc := query.GetLastSeen().GetNodeIdToSequenceId()
	if len(vc) > maxVectorClockLength {
		return fmt.Errorf(
			"vector clock length exceeds maximum of %d",
			maxVectorClockLength,
		)
	}

	return nil
}

func (s *Service) fetchEnvelopes(
	ctx context.Context,
	query *message_api.EnvelopesQuery,
	rowLimit int32,
) ([]queries.GatewayEnvelope, error) {
	// TODO(rich) make this efficient when querying multiple times
	if err := s.validateQuery(query); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid query: %v", err)
	}

	params := queries.SelectGatewayEnvelopesParams{
		Topics:            query.GetTopics(),
		OriginatorNodeIds: make([]int32, 0, len(query.GetOriginatorNodeIds())),
		RowLimit:          rowLimit,
		CursorNodeIds:     nil,
		CursorSequenceIds: nil,
	}

	for _, o := range query.GetOriginatorNodeIds() {
		params.OriginatorNodeIds = append(params.OriginatorNodeIds, int32(o))
	}

	db.SetVectorClock(&params, query.GetLastSeen().GetNodeIdToSequenceId())

	rows, err := queries.New(s.store).SelectGatewayEnvelopes(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
	}

	return rows, nil
}

func (s *Service) PublishEnvelopes(
	ctx context.Context,
	req *message_api.PublishEnvelopesRequest,
) (*message_api.PublishEnvelopesResponse, error) {
	if len(req.GetPayerEnvelopes()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing payer envelope")
	}
	clientEnv, err := s.validatePayerInfo(req.GetPayerEnvelopes()[0])
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
		return &message_api.PublishEnvelopesResponse{}, nil
	}

	// TODO(rich): Properly support batch publishing
	payerBytes, err := proto.Marshal(req.GetPayerEnvelopes()[0])
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

	return &message_api.PublishEnvelopesResponse{
		OriginatorEnvelopes: []*message_api.OriginatorEnvelope{originatorEnv},
	}, nil
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

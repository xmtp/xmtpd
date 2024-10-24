package message

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

const (
	maxRequestedRows     int32         = 1000
	maxQueriesPerRequest int           = 10000
	maxTopicLength       int           = 128
	maxVectorClockLength int           = 100
	pagingInterval       time.Duration = 100 * time.Millisecond
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

	ctx             context.Context
	log             *zap.Logger
	registrant      *registrant.Registrant
	store           *sql.DB
	publishWorker   *publishWorker
	subscribeWorker *subscribeWorker
}

func NewReplicationApiService(
	ctx context.Context,
	log *zap.Logger,
	registrant *registrant.Registrant,
	store *sql.DB,

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
		ctx:             ctx,
		log:             log,
		registrant:      registrant,
		store:           store,
		publishWorker:   publishWorker,
		subscribeWorker: subscribeWorker,
	}, nil
}

func (s *Service) Close() {
	s.log.Info("closed")
}

func (s *Service) SubscribeEnvelopes(
	req *message_api.SubscribeEnvelopesRequest,
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
) error {
	log := s.log.With(zap.String("method", "subscribe"))
	log.Debug("SubscribeEnvelopes", zap.Any("request", req))

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
		case envs, open := <-ch:
			if open {
				err := s.sendEnvelopes(stream, query, envs)
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

// Pulls from DB and sends to client, updating the query's last seen cursor, until
// the stream has caught up to the latest in the database.
func (s *Service) catchUpFromCursor(
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
	query *message_api.EnvelopesQuery,
	logger *zap.Logger,
) error {
	log := logger.With(zap.String("stage", "catchUpFromCursor"))
	if query.GetLastSeen() == nil {
		log.Debug("Skipping catch up")
		// Requester only wants new envelopes
		return nil
	}

	cursor := query.LastSeen.GetNodeIdToSequenceId()
	// GRPC does not distinguish between empty map and nil
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen.NodeIdToSequenceId = cursor
	}

	log.Debug("Catching up from cursor", zap.Any("cursor", cursor))
	for {
		rows, err := s.fetchEnvelopes(stream.Context(), query, maxRequestedRows)
		log.Debug("Fetched envelopes", zap.Any("rows", rows))
		if err != nil {
			return err
		}
		envs := make([]*envelopes.OriginatorEnvelope, 0, len(rows))
		for _, r := range rows {
			env, err := envelopes.NewOriginatorEnvelopeFromBytes(r.OriginatorEnvelope)
			if err != nil {
				// We expect to have already validated the envelope when it was inserted
				logger.Error("could not unmarshal originator envelope", zap.Error(err))
				continue
			}
			envs = append(envs, env)
		}
		err = s.sendEnvelopes(stream, query, envs)
		if err != nil {
			return status.Errorf(codes.Internal, "error sending envelopes: %v", err)
		}
		if len(rows) < int(maxRequestedRows) {
			// There were no more envelopes in DB at time of fetch
			break
		}
		time.Sleep(pagingInterval)
	}
	return nil
}

func (s *Service) sendEnvelopes(
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
	query *message_api.EnvelopesQuery,
	envs []*envelopes.OriginatorEnvelope,
) error {
	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen = &envelopesProto.VectorClock{
			NodeIdToSequenceId: cursor,
		}
	}
	envsToSend := make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))
	for _, env := range envs {
		if cursor[uint32(env.OriginatorNodeID())] >= env.OriginatorSequenceID() {
			continue
		}

		envsToSend = append(envsToSend, env.Proto())
		cursor[uint32(env.OriginatorNodeID())] = env.OriginatorSequenceID()
	}
	if len(envsToSend) == 0 {
		return nil
	}
	err := stream.Send(&message_api.SubscribeEnvelopesResponse{
		Envelopes: envsToSend,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error sending envelopes: %v", err)
	}
	return nil
}

func (s *Service) QueryEnvelopes(
	ctx context.Context,
	req *message_api.QueryEnvelopesRequest,
) (*message_api.QueryEnvelopesResponse, error) {
	log := s.log.With(zap.String("method", "query"))
	if err := s.validateQuery(req.GetQuery()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid query: %v", err)
	}

	limit := int32(req.GetLimit())
	if limit == 0 {
		limit = maxRequestedRows
	}
	rows, err := s.fetchEnvelopes(ctx, req.GetQuery(), limit)
	if err != nil {
		return nil, err
	}

	envs := make([]*envelopesProto.OriginatorEnvelope, 0, len(rows))
	for _, row := range rows {
		originatorEnv := &envelopesProto.OriginatorEnvelope{}
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

func (s *Service) PublishPayerEnvelopes(
	ctx context.Context,
	req *message_api.PublishPayerEnvelopesRequest,
) (*message_api.PublishPayerEnvelopesResponse, error) {
	if len(req.GetPayerEnvelopes()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing payer envelope")
	}

	payerEnv, err := s.validatePayerEnvelope(req.GetPayerEnvelopes()[0])
	if err != nil {
		return nil, err
	}

	// TODO(rich): Properly support batch publishing
	payerBytes, err := payerEnv.Bytes()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal envelope: %v", err)
	}

	targetTopic := payerEnv.ClientEnvelope.TargetTopic()

	stagedEnv, err := queries.New(s.store).
		InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         targetTopic.Bytes(),
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

	return &message_api.PublishPayerEnvelopesResponse{
		OriginatorEnvelopes: []*envelopesProto.OriginatorEnvelope{originatorEnv},
	}, nil
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

func (s *Service) validatePayerEnvelope(
	rawEnv *envelopesProto.PayerEnvelope,
) (*envelopes.PayerEnvelope, error) {
	payerEnv, err := envelopes.NewPayerEnvelope(rawEnv)
	if err != nil {
		return nil, err
	}

	if err := s.validateClientInfo(&payerEnv.ClientEnvelope); err != nil {
		return nil, err
	}

	return payerEnv, nil
}

func (s *Service) validateClientInfo(clientEnv *envelopes.ClientEnvelope) error {
	aad := clientEnv.Aad()
	if aad.GetTargetOriginator() != s.registrant.NodeID() {
		return status.Errorf(codes.InvalidArgument, "invalid target originator")
	}

	if !clientEnv.TopicMatchesPayload() {
		return status.Errorf(codes.InvalidArgument, "topic does not match payload")
	}

	// TODO(rich): Verify all originators have synced past `last_seen`
	// TODO(rich): Check that the blockchain sequence ID is equal to the latest on the group
	// TODO(rich): Perform any payload-specific validation (e.g. identity updates)

	return nil
}

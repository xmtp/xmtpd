// Package message implements the replication API service.
package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/deserializer"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/xmtp/xmtpd/pkg/config"

	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/fees"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	message_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

const (
	maxRequestedRows     int32         = 1000
	maxQueriesPerRequest int           = 10000
	maxTopicLength       int           = 128
	maxVectorClockLength int           = 100
	pagingInterval       time.Duration = 100 * time.Millisecond

	requestMissingMessageError = "missing request message"
)

type Service struct {
	message_apiconnect.UnimplementedReplicationApiHandler

	ctx               context.Context
	logger            *zap.Logger
	registrant        *registrant.Registrant
	store             *sql.DB
	publishWorker     *publishWorker
	subscribeWorker   *subscribeWorker
	validationService mlsvalidate.MLSValidationService
	cu                metadata.CursorUpdater
	feeCalculator     fees.IFeeCalculator
	options           config.APIOptions
	migrationEnabled  bool
}

var _ message_apiconnect.ReplicationApiHandler = (*Service)(nil)

func NewReplicationAPIService(
	ctx context.Context,
	logger *zap.Logger,
	registrant *registrant.Registrant,
	store *sql.DB,
	validationService mlsvalidate.MLSValidationService,
	updater metadata.CursorUpdater,
	feeCalculator fees.IFeeCalculator,
	options config.APIOptions,
	migrationEnabled bool,
	sleepOnFailureTime time.Duration,
) (*Service, error) {
	if validationService == nil {
		return nil, errors.New("validation service must not be nil")
	}

	publishWorker, err := startPublishWorker(
		ctx,
		logger,
		registrant,
		store,
		feeCalculator,
		sleepOnFailureTime,
	)
	if err != nil {
		logger.Error("could not start publish worker", zap.Error(err))
		return nil, err
	}

	subscribeWorker, err := startSubscribeWorker(ctx, logger, store)
	if err != nil {
		logger.Error("could not start subscribe worker", zap.Error(err))
		return nil, err
	}

	return &Service{
		ctx:               ctx,
		logger:            logger,
		registrant:        registrant,
		store:             store,
		publishWorker:     publishWorker,
		subscribeWorker:   subscribeWorker,
		validationService: validationService,
		cu:                updater,
		feeCalculator:     feeCalculator,
		options:           options,
		migrationEnabled:  migrationEnabled,
	}, nil
}

func (s *Service) Close() {
	s.logger.Debug("closed")
}

func (s *Service) SubscribeEnvelopes(
	ctx context.Context,
	req *connect.Request[message_api.SubscribeEnvelopesRequest],
	stream *connect.ServerStream[message_api.SubscribeEnvelopesResponse],
) error {
	if req.Msg == nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request",
			utils.BodyField(req),
		)
	}

	// Send a keepalive immediately, so wasm based clients maintain the connection open.
	err := stream.Send(&message_api.SubscribeEnvelopesResponse{})
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not send keepalive: %w", err),
		)
	}

	query := req.Msg.GetQuery()

	if err := s.validateQuery(query); err != nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("invalid subscription request: %w", err),
		)
	}

	envelopesCh := s.subscribeWorker.listen(ctx, query)

	err = s.catchUpFromCursor(ctx, stream, query, logger)
	if err != nil {
		return err
	}

	// GRPC keep-alives are not sufficient in some load balanced environments.
	// We need to send an actual payload: https://github.com/xmtp/xmtpd/issues/669
	ticker := time.NewTicker(s.options.SendKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send a keepalive at the interval specified in the config.
			err = stream.Send(&message_api.SubscribeEnvelopesResponse{})
			if err != nil {
				return connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("could not send keepalive: %w", err),
				)
			}

		case envs, open := <-envelopesCh:
			ticker.Reset(s.options.SendKeepAliveInterval)

			if !open {
				logger.Debug("channel closed by worker")
				return nil
			}

			err = s.sendEnvelopes(stream, query, envs)
			if err != nil {
				return connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("error sending envelope: %w", err),
				)
			}

		case <-ctx.Done():
			logger.Debug("message subscription stream closed")
			return nil

		case <-s.ctx.Done():
			logger.Debug("message service closed")
			return nil
		}
	}
}

// Pulls from DB and sends to client, updating the query's last seen cursor, until
// the stream has caught up to the latest in the database.
func (s *Service) catchUpFromCursor(
	ctx context.Context,
	stream *connect.ServerStream[message_api.SubscribeEnvelopesResponse],
	query *message_api.EnvelopesQuery,
	logger *zap.Logger,
) error {
	if query.GetLastSeen() == nil {
		logger.Debug("skipping catch up")
		// Requester only wants new envelopes
		return nil
	}

	cursor := query.LastSeen.GetNodeIdToSequenceId()
	// GRPC does not distinguish between empty map and nil
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen.NodeIdToSequenceId = cursor
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("catching up from cursor", utils.BodyField(cursor))
	}

	for {
		rows, err := s.fetchEnvelopes(ctx, query, maxRequestedRows)
		if err != nil {
			return err
		}

		if s.logger.Core().Enabled(zap.DebugLevel) {
			logger.Debug("fetched envelopes", utils.BodyField(rows))
		}

		envs := make([]*envelopes.OriginatorEnvelope, 0, len(rows))
		for _, r := range rows {
			env, err := envelopes.NewOriginatorEnvelopeFromBytes(r.OriginatorEnvelope)
			if err != nil {
				// We expect to have already validated the envelope when it was inserted
				s.logger.Error("could not unmarshal originator envelope", zap.Error(err))
				continue
			}
			envs = append(envs, env)
		}

		err = s.sendEnvelopes(stream, query, envs)
		if err != nil {
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error sending envelopes: %w", err),
			)
		}
		if len(rows) < int(maxRequestedRows) {
			// There were no more envelopes in DB at time of fetch
			break
		}
		time.Sleep(pagingInterval)
	}

	return nil
}

// TODO: Make this method context aware.
func (s *Service) sendEnvelopes(
	stream *connect.ServerStream[message_api.SubscribeEnvelopesResponse],
	query *message_api.EnvelopesQuery,
	envs []*envelopes.OriginatorEnvelope,
) error {
	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen = &envelopesProto.Cursor{
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
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error sending envelopes: %w", err),
		)
	}

	return nil
}

func (s *Service) QueryEnvelopes(
	ctx context.Context,
	req *connect.Request[message_api.QueryEnvelopesRequest],
) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	if err := s.validateQuery(req.Msg.GetQuery()); err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("invalid query: %w", err),
		)
	}

	var limit int32
	if req.Msg.GetLimit() > uint32(maxRequestedRows) || req.Msg.GetLimit() == 0 {
		limit = maxRequestedRows
	} else {
		limit = int32(req.Msg.GetLimit())
	}

	rows, err := s.fetchEnvelopes(ctx, req.Msg.GetQuery(), limit)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&message_api.QueryEnvelopesResponse{
		Envelopes: make([]*envelopesProto.OriginatorEnvelope, 0, len(rows)),
	})

	for _, row := range rows {
		originatorEnv := &envelopesProto.OriginatorEnvelope{}
		err := proto.Unmarshal(row.OriginatorEnvelope, originatorEnv)
		if err != nil {
			// We expect to have already validated the envelope when it was inserted
			logger.Error("could not unmarshal originator envelope", zap.Error(err))
			continue
		}
		response.Msg.Envelopes = append(response.Msg.Envelopes, originatorEnv)
	}

	return response, nil
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
) ([]queries.GatewayEnvelopesView, error) {
	if len(query.GetTopics()) != 0 {
		params := queries.SelectGatewayEnvelopesByTopicsParams{
			Topics:            query.GetTopics(),
			RowLimit:          rowLimit,
			CursorNodeIds:     nil,
			CursorSequenceIds: nil,
		}

		db.SetVectorClockByTopics(&params, query.GetLastSeen().GetNodeIdToSequenceId())

		rows, err := queries.New(s.store).SelectGatewayEnvelopesByTopics(ctx, params)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("could not select envelopes: %w", err),
			)
		}

		return db.TransformRowsByTopic(rows), nil
	}
	if len(query.GetOriginatorNodeIds()) != 0 {
		params := queries.SelectGatewayEnvelopesByOriginatorsParams{
			OriginatorNodeIds: make([]int32, 0, len(query.GetOriginatorNodeIds())),
			RowLimit:          rowLimit,
			CursorNodeIds:     nil,
			CursorSequenceIds: nil,
		}
		for _, o := range query.GetOriginatorNodeIds() {
			params.OriginatorNodeIds = append(params.OriginatorNodeIds, int32(o))
		}

		db.SetVectorClockByOriginators(&params, query.GetLastSeen().GetNodeIdToSequenceId())

		rows, err := queries.New(s.store).SelectGatewayEnvelopesByOriginators(ctx, params)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("could not select envelopes: %w", err),
			)
		}

		return db.TransformRowsByOriginator(rows), nil
	}

	params := queries.SelectGatewayEnvelopesUnfilteredParams{
		RowLimit:          rowLimit,
		CursorNodeIds:     nil,
		CursorSequenceIds: nil,
	}
	db.SetVectorClockUnfiltered(&params, query.GetLastSeen().GetNodeIdToSequenceId())

	rows, err := queries.New(s.store).SelectGatewayEnvelopesUnfiltered(ctx, params)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not select envelopes: %w", err),
		)
	}

	return rows, nil
}

func (s *Service) PublishPayerEnvelopes(
	ctx context.Context,
	req *connect.Request[message_api.PublishPayerEnvelopesRequest],
) (*connect.Response[message_api.PublishPayerEnvelopesResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	if s.migrationEnabled {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("D14N API is read-only while migration is enabled"),
		)
	}

	payerEnvelopes := req.Msg.GetPayerEnvelopes()

	if len(payerEnvelopes) == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("missing payer envelope"),
		)
	}

	payerEnvelope, err := s.validatePayerEnvelope(payerEnvelopes[0])
	if err != nil {
		return nil, err
	}

	// TODO(rich): Properly support batch publishing
	payerBytes, err := payerEnvelope.Bytes()
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not marshal envelope: %w", err),
		)
	}

	targetTopic := payerEnvelope.ClientEnvelope.TargetTopic()
	topicKind := targetTopic.Kind()

	if targetTopic.IsReserved() {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("reserved topics cannot be published to by gateways"),
		)
	}

	if topicKind == topic.TopicKindIdentityUpdatesV1 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("identity updates must be published via the blockchain"),
		)
	}

	if topicKind == topic.TopicKindGroupMessagesV1 {
		if err = s.validateGroupMessage(&payerEnvelope.ClientEnvelope); err != nil {
			return nil, err
		}
	}

	if topicKind == topic.TopicKindKeyPackagesV1 {
		if err = s.validateKeyPackage(ctx, &payerEnvelope.ClientEnvelope); err != nil {
			return nil, err
		}
	}

	stagedEnvelope, err := queries.New(s.store).
		InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         targetTopic.Bytes(),
			PayerEnvelope: payerBytes,
		})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not insert staged envelope: %w", err),
		)
	}
	s.publishWorker.notifyStagedPublish()

	baseFee, congestionFee, err := s.publishWorker.calculateFees(
		&stagedEnvelope,
		payerEnvelope.Proto().GetMessageRetentionDays(),
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not calculate fees: %w", err),
		)
	}

	originatorEnvelope, err := s.registrant.SignStagedEnvelope(
		stagedEnvelope,
		baseFee,
		congestionFee,
		payerEnvelope.Proto().GetMessageRetentionDays(),
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not sign envelope: %w", err),
		)
	}

	s.waitForGatewayPublish(ctx, stagedEnvelope, logger)

	return connect.NewResponse(&message_api.PublishPayerEnvelopesResponse{
		OriginatorEnvelopes: []*envelopesProto.OriginatorEnvelope{originatorEnvelope},
	}), nil
}

func (s *Service) validateGroupMessage(
	clientEnv *envelopes.ClientEnvelope,
) error {
	payload, ok := clientEnv.Payload().(*envelopesProto.ClientEnvelope_GroupMessage)
	if !ok {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid payload type"),
		)
	}

	isCommit, err := deserializer.IsGroupMessageCommit(payload)
	if err != nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid group message"),
		)
	}

	if isCommit {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("commit messages must be published via the blockchain"),
		)
	}

	return nil
}

func (s *Service) GetInboxIds(
	ctx context.Context,
	req *connect.Request[message_api.GetInboxIdsRequest],
) (*connect.Response[message_api.GetInboxIdsResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	var (
		queries   = queries.New(s.store)
		addresses = []string{}
	)

	for _, request := range req.Msg.Requests {
		addresses = append(addresses, request.GetIdentifier())
	}

	addressLogEntries, err := queries.GetAddressLogs(ctx, addresses)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&message_api.GetInboxIdsResponse{
		Responses: make([]*message_api.GetInboxIdsResponse_Response, len(addresses)),
	})

	for index, address := range addresses {
		resp := message_api.GetInboxIdsResponse_Response{}
		resp.Identifier = address

		for _, logEntry := range addressLogEntries {
			if logEntry.Address == address {
				inboxID := logEntry.InboxID
				resp.InboxId = &inboxID
			}
		}
		response.Msg.Responses[index] = &resp
	}

	logger.Debug("got inbox ids", utils.NumResponsesField(len(response.Msg.Responses)))

	return response, nil
}

func (s *Service) GetNewestEnvelope(
	ctx context.Context,
	req *connect.Request[message_api.GetNewestEnvelopeRequest],
) (*connect.Response[message_api.GetNewestEnvelopeResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	var (
		queries      = queries.New(s.store)
		topics       = req.Msg.Topics
		originalSort = make(map[string]int)
	)

	for idx, topic := range topics {
		originalSort[string(topic)] = idx
	}

	rows, err := queries.SelectNewestFromTopics(ctx, topics)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not select envelopes: %w", err),
		)
	}

	logger.Debug(
		"received newest envelopes for topics",
		utils.NumEnvelopesField(len(rows)),
		utils.NumTopicsField(len(topics)),
	)

	response := connect.NewResponse(&message_api.GetNewestEnvelopeResponse{
		Results: make([]*message_api.GetNewestEnvelopeResponse_Response, len(topics)),
	})

	for _, row := range rows {
		idx, ok := originalSort[string(row.Topic)]
		if !ok {
			// We will leave the index empty if there are no envelopes for that topic
			continue
		}
		originatorEnv := &envelopesProto.OriginatorEnvelope{}
		err := proto.Unmarshal(row.OriginatorEnvelope, originatorEnv)
		if err != nil {
			// We expect to have already validated the envelope when it was inserted
			logger.Error("could not unmarshal originator envelope", zap.Error(err))
			continue
		}

		response.Msg.Results[idx] = &message_api.GetNewestEnvelopeResponse_Response{
			OriginatorEnvelope: originatorEnv,
		}
	}

	return response, nil
}

func (s *Service) validatePayerEnvelope(
	rawEnv *envelopesProto.PayerEnvelope,
) (*envelopes.PayerEnvelope, error) {
	payerEnv, err := envelopes.NewPayerEnvelope(rawEnv)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("could not unmarshal payer envelope: %w", err),
		)
	}

	if payerEnv.TargetOriginator != s.registrant.NodeID() {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid target originator"),
		)
	}

	if _, err = payerEnv.RecoverSigner(); err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("could not recover signer: %w", err),
		)
	}

	if err = s.validateClientInfo(&payerEnv.ClientEnvelope); err != nil {
		return nil, err
	}

	err = s.validateExpiry(payerEnv)
	if err != nil {
		return nil, err
	}

	return payerEnv, nil
}

func (s *Service) validateExpiry(payerEnv *envelopes.PayerEnvelope) error {
	// the payload should be valid for at least for 2 days
	if payerEnv.RetentionDays() < 2 {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid expiry retention days. Must be >= 2"),
		)
	}

	// more than a ~year sounds like a mistake
	if payerEnv.RetentionDays() != math.MaxUint32 && payerEnv.RetentionDays() > 365 {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid expiry retention days. Must be <= 365"),
		)
	}

	return nil
}

func (s *Service) validateKeyPackage(
	ctx context.Context,
	clientEnv *envelopes.ClientEnvelope,
) error {
	payload, ok := clientEnv.Payload().(*envelopesProto.ClientEnvelope_UploadKeyPackage)
	if !ok {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid payload type"),
		)
	}

	validationResult, err := s.validationService.ValidateKeyPackages(
		ctx,
		[][]byte{payload.UploadKeyPackage.KeyPackage.KeyPackageTlsSerialized},
	)
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not validate key package: %w", err),
		)
	}

	if len(validationResult) == 0 {
		return connect.NewError(
			connect.CodeInternal,
			errors.New("no validation results"),
		)
	}

	if !validationResult[0].IsOk {
		return connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("key package validation failed: %s", validationResult[0].ErrorMessage),
		)
	}

	return nil
}

func (s *Service) validateClientInfo(clientEnv *envelopes.ClientEnvelope) error {
	aad := clientEnv.Aad()

	if aad == nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("authenticated data is missing"),
		)
	}

	if !clientEnv.TopicMatchesPayload() {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("topic does not match payload"),
		)
	}

	if aad.GetDependsOn() != nil {
		lastSeenCursor := s.cu.GetCursor()
		for nodeID, seqID := range aad.GetDependsOn().NodeIdToSequenceId {
			lastSeqID, exists := lastSeenCursor.NodeIdToSequenceId[nodeID]
			if nodeID >= 100 {
				// The failure scenarios of non-commits are different from the blockchain path
				// and as such should be prevented
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"node ID %d specified in DependsOn is not a valid node ID, a message can not depend on a non-commit",
						nodeID,
					),
				)
			} else if !exists {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"node ID %d specified in DependsOn has not been seen by this node",
						nodeID,
					),
				)
			} else if seqID > lastSeqID {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"sequence ID %d for node ID %d specified in DependsOn exceeds last seen sequence ID %d",
						seqID,
						nodeID,
						lastSeqID,
					),
				)
			}
		}
	}
	// TODO(rich): Check that the blockchain sequence ID is equal to the latest on the group
	// TODO(rich): Perform any payload-specific validation (e.g. identity updates)

	return nil
}

func (s *Service) waitForGatewayPublish(
	ctx context.Context,
	stagedEnv queries.StagedOriginatorEnvelope,
	logger *zap.Logger,
) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger = logger.With(
			utils.SequenceIDField(stagedEnv.ID),
			utils.EnvelopeIDField(stagedEnv.ID),
		)
	}

	startTime := time.Now()
	timeout := time.After(30 * time.Second)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			if s.logger.Core().Enabled(zap.DebugLevel) {
				logger.Debug(
					"timeout waiting for publisher",
					utils.LastProcessedField(s.publishWorker.lastProcessed.Load()),
				)
			}
			return

		case <-ctx.Done():
			if s.logger.Core().Enabled(zap.DebugLevel) {
				logger.Debug(
					"context cancelled while waiting for publisher",
					utils.LastProcessedField(s.publishWorker.lastProcessed.Load()),
				)
			}
			return

		case <-ticker.C:
			// Check if the last processed ID has reached or exceeded the current ID
			if s.publishWorker.lastProcessed.Load() >= stagedEnv.ID {
				if s.logger.Core().Enabled(zap.DebugLevel) {
					logger.Debug(
						"finished waiting for publisher",
						utils.DurationMsField(time.Since(startTime)),
					)
				}

				return
			}
		}
	}
}

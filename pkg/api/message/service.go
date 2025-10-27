// Package message implements the replication API service.
package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

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
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/grpc/codes"
	metaProtos "google.golang.org/grpc/metadata"
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

	subscribeEnvelopesMethod    = "SubscribeEnvelopes"
	queryEnvelopesMethod        = "QueryEnvelopes"
	publishPayerEnvelopesMethod = "PublishPayerEnvelopes"
	getInboxIdsMethod           = "GetInboxIds"
	getNewestEnvelopeMethod     = "GetNewestEnvelope"
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

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
) (*Service, error) {
	if validationService == nil {
		return nil, errors.New("validation service must not be nil")
	}

	publishWorker, err := startPublishWorker(ctx, logger, registrant, store, feeCalculator)
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
	req *message_api.SubscribeEnvelopesRequest,
	stream message_api.ReplicationApi_SubscribeEnvelopesServer,
) error {
	logger := s.logger.With(utils.MethodField(subscribeEnvelopesMethod))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request",
			utils.BodyField(req),
		)
	}

	// Send a header (any header) to fix an issue with Tonic based GRPC clients.
	// See: https://github.com/xmtp/libxmtp/pull/58
	err := stream.SendHeader(metaProtos.Pairs("subscribed", "true"))
	if err != nil {
		return status.Errorf(codes.Internal, "could not send header: %v", err)
	}

	query := req.GetQuery()
	if err := s.validateQuery(query); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid subscription request: %v", err)
	}

	ch := s.subscribeWorker.listen(stream.Context(), query)
	err = s.catchUpFromCursor(stream, query, logger)
	if err != nil {
		return err
	}

	// GRPC keep-alives are not sufficient in some load balanced environments
	// we need to send an actual payload
	// see https://github.com/xmtp/xmtpd/issues/669
	ticker := time.NewTicker(s.options.SendKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = stream.Send(&message_api.SubscribeEnvelopesResponse{})
			if err != nil {
				return status.Errorf(codes.Internal, "could not send keepalive: %v", err)
			}
		case envs, open := <-ch:
			ticker.Reset(s.options.SendKeepAliveInterval)
			if open {
				err := s.sendEnvelopes(stream, query, envs)
				if err != nil {
					return status.Errorf(codes.Internal, "error sending envelope: %v", err)
				}
			} else {
				// TODO(rich) Reset whole subscribe flow from new cursor
				logger.Debug("channel closed by worker")
				return nil
			}
		case <-stream.Context().Done():
			logger.Debug("stream closed")
			return nil
		case <-s.ctx.Done():
			logger.Debug("service closed")
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
		rows, err := s.fetchEnvelopes(stream.Context(), query, maxRequestedRows)
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
		return status.Errorf(codes.Internal, "error sending envelopes: %v", err)
	}
	return nil
}

func (s *Service) QueryEnvelopes(
	ctx context.Context,
	req *message_api.QueryEnvelopesRequest,
) (*message_api.QueryEnvelopesResponse, error) {
	logger := s.logger.With(utils.MethodField(queryEnvelopesMethod))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

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
			logger.Error("could not unmarshal originator envelope", zap.Error(err))
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
) ([]queries.GatewayEnvelopesV2View, error) {
	if len(query.GetTopics()) != 0 {
		params := queries.SelectGatewayEnvelopesV2ByTopicsParams{
			Topics:            query.GetTopics(),
			RowLimit:          rowLimit,
			CursorNodeIds:     nil,
			CursorSequenceIds: nil,
		}

		db.SetVectorClockByTopics(&params, query.GetLastSeen().GetNodeIdToSequenceId())

		rows, err := queries.New(s.store).SelectGatewayEnvelopesV2ByTopics(ctx, params)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
		}

		return rows, nil
	}
	if len(query.GetOriginatorNodeIds()) != 0 {
		params := queries.SelectGatewayEnvelopesV2ByOriginatorsParams{
			OriginatorNodeIds: make([]int32, 0, len(query.GetOriginatorNodeIds())),
			RowLimit:          rowLimit,
			CursorNodeIds:     nil,
			CursorSequenceIds: nil,
		}
		for _, o := range query.GetOriginatorNodeIds() {
			params.OriginatorNodeIds = append(params.OriginatorNodeIds, int32(o))
		}

		db.SetVectorClockByOriginators(&params, query.GetLastSeen().GetNodeIdToSequenceId())

		rows, err := queries.New(s.store).SelectGatewayEnvelopesV2ByOriginators(ctx, params)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
		}

		return rows, nil
	}

	params := queries.SelectGatewayEnvelopesV2UnfilteredParams{
		RowLimit:          rowLimit,
		CursorNodeIds:     nil,
		CursorSequenceIds: nil,
	}
	db.SetVectorClockUnfiltered(&params, query.GetLastSeen().GetNodeIdToSequenceId())

	rows, err := queries.New(s.store).SelectGatewayEnvelopesV2Unfiltered(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
	}

	return rows, nil
}

func (s *Service) PublishPayerEnvelopes(
	ctx context.Context,
	req *message_api.PublishPayerEnvelopesRequest,
) (*message_api.PublishPayerEnvelopesResponse, error) {
	logger := s.logger.With(utils.MethodField(publishPayerEnvelopesMethod))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	if s.migrationEnabled {
		return nil, status.Errorf(
			codes.Internal,
			"D14N API is read-only while migration is enabled",
		)
	}

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
	topicKind := targetTopic.Kind()

	if targetTopic.IsReserved() {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"reserved topics cannot be published to by Payers",
		)
	}

	if topicKind == topic.TopicKindIdentityUpdatesV1 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"identity updates must be published via the blockchain",
		)
	}

	if topicKind == topic.TopicKindGroupMessagesV1 {
		if err = s.validateGroupMessage(&payerEnv.ClientEnvelope); err != nil {
			return nil, err
		}
	}

	if topicKind == topic.TopicKindKeyPackagesV1 {
		if err = s.validateKeyPackage(ctx, &payerEnv.ClientEnvelope); err != nil {
			return nil, err
		}
	}

	stagedEnv, err := queries.New(s.store).
		InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         targetTopic.Bytes(),
			PayerEnvelope: payerBytes,
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert staged envelope: %v", err)
	}
	s.publishWorker.notifyStagedPublish()

	baseFee, congestionFee, err := s.publishWorker.calculateFees(
		&stagedEnv,
		payerEnv.Proto().GetMessageRetentionDays(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not calculate fees: %v", err)
	}

	originatorEnv, err := s.registrant.SignStagedEnvelope(
		stagedEnv,
		baseFee,
		congestionFee,
		payerEnv.Proto().GetMessageRetentionDays(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not sign envelope: %v", err)
	}

	s.waitForGatewayPublish(ctx, stagedEnv, logger)

	return &message_api.PublishPayerEnvelopesResponse{
		OriginatorEnvelopes: []*envelopesProto.OriginatorEnvelope{originatorEnv},
	}, nil
}

func (s *Service) validateGroupMessage(
	clientEnv *envelopes.ClientEnvelope,
) error {
	payload, ok := clientEnv.Payload().(*envelopesProto.ClientEnvelope_GroupMessage)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid payload type")
	}

	isCommit, err := deserializer.IsGroupMessageCommit(payload)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid group message")
	}

	if isCommit {
		return status.Errorf(
			codes.InvalidArgument,
			"commit messages must be published via the blockchain",
		)
	}

	return nil
}

func (s *Service) GetInboxIds(
	ctx context.Context,
	req *message_api.GetInboxIdsRequest,
) (*message_api.GetInboxIdsResponse, error) {
	logger := s.logger.With(utils.MethodField(getInboxIdsMethod))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	queries := queries.New(s.store)
	addresses := []string{}
	for _, request := range req.Requests {
		addresses = append(addresses, request.GetIdentifier())
	}

	addressLogEntries, err := queries.GetAddressLogs(ctx, addresses)
	if err != nil {
		return nil, err
	}

	out := make([]*message_api.GetInboxIdsResponse_Response, len(addresses))

	for index, address := range addresses {
		resp := message_api.GetInboxIdsResponse_Response{}
		resp.Identifier = address

		for _, logEntry := range addressLogEntries {
			if logEntry.Address == address {
				inboxID := logEntry.InboxID
				resp.InboxId = &inboxID
			}
		}
		out[index] = &resp
	}

	logger.Debug("got inbox ids", utils.NumResponsesField(len(out)))

	return &message_api.GetInboxIdsResponse{
		Responses: out,
	}, nil
}

func (s *Service) GetNewestEnvelope(
	ctx context.Context,
	req *message_api.GetNewestEnvelopeRequest,
) (*message_api.GetNewestEnvelopeResponse, error) {
	logger := s.logger.With(utils.MethodField(getNewestEnvelopeMethod))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	queries := queries.New(s.store)
	topics := req.GetTopics()
	originalSort := make(map[string]int)

	for idx, topic := range topics {
		originalSort[string(topic)] = idx
	}

	rows, err := queries.SelectNewestFromTopicsV2(ctx, topics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not select envelopes: %v", err)
	}

	logger.Debug(
		"received newest envelopes for topics",
		utils.NumEnvelopesField(len(rows)),
		utils.NumTopicsField(len(topics)),
	)

	results := make([]*message_api.GetNewestEnvelopeResponse_Response, len(topics))
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

		results[idx] = &message_api.GetNewestEnvelopeResponse_Response{
			OriginatorEnvelope: originatorEnv,
		}
	}

	return &message_api.GetNewestEnvelopeResponse{
		Results: results,
	}, nil
}

func (s *Service) validatePayerEnvelope(
	rawEnv *envelopesProto.PayerEnvelope,
) (*envelopes.PayerEnvelope, error) {
	payerEnv, err := envelopes.NewPayerEnvelope(rawEnv)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if payerEnv.TargetOriginator != s.registrant.NodeID() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target originator")
	}

	if _, err = payerEnv.RecoverSigner(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
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
		return status.Errorf(codes.InvalidArgument, "invalid expiry retention days. Must be >= 2")
	}

	// more than a ~year sounds like a mistake
	if payerEnv.RetentionDays() != math.MaxUint32 && payerEnv.RetentionDays() > 365 {
		return status.Errorf(codes.InvalidArgument, "invalid expiry retention days. Must be <= 365")
	}

	return nil
}

func (s *Service) validateKeyPackage(
	ctx context.Context,
	clientEnv *envelopes.ClientEnvelope,
) error {
	payload, ok := clientEnv.Payload().(*envelopesProto.ClientEnvelope_UploadKeyPackage)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid payload type")
	}

	validationResult, err := s.validationService.ValidateKeyPackages(
		ctx,
		[][]byte{payload.UploadKeyPackage.KeyPackage.KeyPackageTlsSerialized},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "could not validate key package: %v", err)
	}

	if len(validationResult) == 0 {
		return status.Errorf(codes.Internal, "no validation results")
	}

	if !validationResult[0].IsOk {
		return status.Errorf(
			codes.InvalidArgument,
			"key package validation failed: %s",
			validationResult[0].ErrorMessage,
		)
	}

	return nil
}

func (s *Service) validateClientInfo(clientEnv *envelopes.ClientEnvelope) error {
	aad := clientEnv.Aad()

	if aad == nil {
		return status.Errorf(codes.InvalidArgument, "aad is missing")
	}

	if !clientEnv.TopicMatchesPayload() {
		return status.Errorf(codes.InvalidArgument, "topic does not match payload")
	}

	if aad.GetDependsOn() != nil {
		lastSeenCursor := s.cu.GetCursor()
		for nodeID, seqID := range aad.GetDependsOn().NodeIdToSequenceId {
			lastSeqID, exists := lastSeenCursor.NodeIdToSequenceId[nodeID]
			if nodeID >= 100 {
				// The failure scenarios of non-commits are different from the blockchain path
				// and as such should be prevented
				return status.Errorf(
					codes.InvalidArgument,
					"node ID %d specified in DependsOn is not a valid node ID. A message can not depend on a non-commit.",
					nodeID,
				)
			} else if !exists {
				return status.Errorf(codes.InvalidArgument,
					"node ID %d specified in DependsOn has not been seen by this node",
					nodeID,
				)
			} else if seqID > lastSeqID {
				return status.Errorf(codes.InvalidArgument,
					"sequence ID %d for node ID %d specified in DependsOn exceeds last seen sequence ID %d",
					seqID,
					nodeID,
					lastSeqID,
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
					"Timeout waiting for publisher",
					utils.LastProcessedField(s.publishWorker.lastProcessed.Load()),
				)
			}
			return

		case <-ctx.Done():
			if s.logger.Core().Enabled(zap.DebugLevel) {
				logger.Debug(
					"Context cancelled while waiting for publisher",
					utils.LastProcessedField(s.publishWorker.lastProcessed.Load()),
				)
			}
			return

		case <-ticker.C:
			// Check if the last processed ID has reached or exceeded the current ID
			if s.publishWorker.lastProcessed.Load() >= stagedEnv.ID {
				if s.logger.Core().Enabled(zap.DebugLevel) {
					logger.Debug(
						"Finished waiting for publisher",
						utils.DurationMsField(time.Since(startTime)),
					)
				}

				return
			}
		}
	}
}

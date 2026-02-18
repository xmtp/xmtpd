package message

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"connectrpc.com/connect"
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/utils/retryerrors"
)

const (
	maxTopicsPerChunk int   = 500
	topicPageLimit    int32 = 500
	maxTopicFilters   int   = 10000
)

// SubscribeTopics implements the per-topic cursor subscribe API.
func (s *Service) SubscribeTopics(
	ctx context.Context,
	req *connect.Request[message_api.SubscribeTopicsRequest],
	stream *connect.ServerStream[message_api.SubscribeTopicsResponse],
) error {
	if req.Msg == nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	logger := s.logger.With(utils.MethodField(req.Spec().Procedure))

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("received request", utils.BodyField(req))
	}

	// Send STARTED status so wasm-based clients maintain the connection open.
	err := stream.Send(newSubscriptionStatusMessage(
		message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_STARTED,
	))
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not send status: %w", err),
		)
	}

	filters := req.Msg.GetFilters()

	knownOriginators, err := s.originatorList.GetOriginatorNodeIDs(ctx)
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not get originator list: %w", err),
		)
	}

	if err := validateTopicFilters(filters, knownOriginators); err != nil {
		return err
	}

	cursors, topics, catchUpKeys := buildTopicCursors(filters)

	envelopesCh := s.subscribeWorker.listen(ctx, &message_api.EnvelopesQuery{
		Topics: topics,
	})

	err = s.catchUpTopics(ctx, stream, cursors, catchUpKeys, logger)
	if err != nil {
		return err
	}

	err = stream.Send(newSubscriptionStatusMessage(
		message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_CATCHUP_COMPLETE,
	))
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not send status: %w", err),
		)
	}

	// GRPC keep-alives are not sufficient in some load balanced environments.
	ticker := time.NewTicker(s.options.SendKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = stream.Send(newSubscriptionStatusMessage(
				message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_WAITING,
			))
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

			// Advance cursors to filter duplicates between catch-up and live delivery.
			envsToSend := advanceTopicCursors(cursors, envs, logger)
			err = s.sendTopicEnvelopes(stream, envsToSend)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			logger.Debug("topic subscription stream closed")
			return nil

		case <-s.ctx.Done():
			logger.Debug("message service closed")
			return nil
		}
	}
}

// validateTopicFilters validates the topic filters in a SubscribeTopicsRequest.
func validateTopicFilters(
	filters []*message_api.SubscribeTopicsRequest_TopicFilter,
	knownOriginators []int32,
) error {
	if len(filters) == 0 {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("filters must not be empty"),
		)
	}

	if len(filters) > maxTopicFilters {
		return connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("too many filters: %d, maximum is %d", len(filters), maxTopicFilters),
		)
	}

	referencedOriginators := make(map[uint32]struct{})
	for _, f := range filters {
		if err := validateTopicFilter(f); err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		for origID := range f.GetLastSeen().GetNodeIdToSequenceId() {
			referencedOriginators[origID] = struct{}{}
		}
	}

	if err := validateOriginatorIDs(referencedOriginators, knownOriginators); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	return nil
}

// validateTopicFilter validates a single topic filter's field lengths.
func validateTopicFilter(f *message_api.SubscribeTopicsRequest_TopicFilter) error {
	topicBytes := f.GetTopic()
	if len(topicBytes) == 0 || len(topicBytes) > maxTopicLength {
		return fmt.Errorf("invalid topic length: %d", len(topicBytes))
	}

	vc := f.GetLastSeen().GetNodeIdToSequenceId()
	if len(vc) > maxVectorClockLength {
		return fmt.Errorf("vector clock length exceeds maximum of %d", maxVectorClockLength)
	}

	return nil
}

// validateOriginatorIDs checks that all referenced originator IDs are known.
func validateOriginatorIDs(
	referenced map[uint32]struct{},
	knownOriginators []int32,
) error {
	if len(referenced) == 0 {
		return nil
	}

	known := make(map[uint32]struct{}, len(knownOriginators))
	for _, id := range knownOriginators {
		known[uint32(id)] = struct{}{}
	}

	for origID := range referenced {
		if _, ok := known[origID]; !ok {
			return fmt.Errorf("unknown originator node ID in cursor: %d", origID)
		}
	}

	return nil
}

// buildTopicCursors converts topic filters into a TopicCursors map,
// the deduplicated topic list (as [][]byte), and the keys that need
// catch-up (those with non-nil LastSeen).
func buildTopicCursors(
	filters []*message_api.SubscribeTopicsRequest_TopicFilter,
) (db.TopicCursors, [][]byte, []string) {
	cursors := make(db.TopicCursors, len(filters))
	topics := make([][]byte, 0, len(filters))
	catchUpKeys := make([]string, 0, len(filters))

	for _, f := range filters {
		key := string(f.GetTopic())

		if _, exists := cursors[key]; exists {
			continue
		}

		topics = append(topics, f.GetTopic())

		lastSeen := f.GetLastSeen()

		if lastSeen != nil {
			// Copy the cursor map for catch-up.
			vc := lastSeen.GetNodeIdToSequenceId()
			cursorCopy := make(db.VectorClock, len(vc))
			maps.Copy(cursorCopy, vc)
			cursors[key] = cursorCopy
			catchUpKeys = append(catchUpKeys, key)
		} else {
			// No catch-up needed — live only. Create empty VectorClock for dedup.
			cursors[key] = make(db.VectorClock)
		}
	}

	return cursors, topics, catchUpKeys
}

// fillMissingOriginatorsForTopics calls FillMissingOriginators on each
// topic's VectorClock in the given keys list.
func fillMissingOriginatorsForTopics(
	cursors db.TopicCursors,
	keys []string,
	allOriginators []int32,
) {
	for _, key := range keys {
		if vc, ok := cursors[key]; ok {
			db.FillMissingOriginators(vc, allOriginators)
		}
	}
}

// fetchTopicEnvelopesWithRetry fetches envelopes using exponential backoff.
func (s *Service) fetchTopicEnvelopesWithRetry(
	ctx context.Context,
	subCursors db.TopicCursors,
	rowLimit int32,
	rowsPerEntry int32,
) ([]queries.GatewayEnvelopesView, error) {
	boCtx := backoff.WithContext(
		utils.NewBackoff(50*time.Millisecond, 300*time.Millisecond, 2*time.Second), ctx,
	)

	var result []queries.GatewayEnvelopesView

	operation := func() error {
		params := queries.SelectGatewayEnvelopesByPerTopicCursorsParams{
			RowsPerEntry: rowsPerEntry,
			RowLimit:     rowLimit,
		}
		db.SetPerTopicCursors(&params, subCursors)

		rows, err := s.store.ReadQuery().SelectGatewayEnvelopesByPerTopicCursors(ctx, params)
		if err == nil {
			result = db.TransformRowsByPerTopicCursors(rows)
			return nil
		}
		if !retryerrors.IsRetryableSQLError(err) {
			return backoff.Permanent(err)
		}
		return err
	}

	err := backoff.Retry(operation, boCtx)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not select envelopes: %w", err),
		)
	}

	return result, nil
}

// advanceTopicCursors filters envelopes against subscribed topics, deduplicates
// using current cursor state, and advances cursors in-place for each envelope
// that passes the filter. Returns the proto envelopes ready to send.
//
// Cursor advancement is critical: in catchUpTopics it drives pagination (the next
// query uses the advanced cursors), and in the live loop it prevents duplicates
// between catch-up and live delivery.
func advanceTopicCursors(
	cursors db.TopicCursors,
	envs []*envelopes.OriginatorEnvelope,
	logger *zap.Logger,
) []*envelopesProto.OriginatorEnvelope {
	result := make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))

	for _, env := range envs {
		vc, ok := cursors[string(env.TargetTopic().Bytes())]
		if !ok {
			logger.Warn(
				"received envelope for unsubscribed topic",
				zap.Binary("topic", env.TargetTopic().Bytes()),
			)
			continue
		}

		origID := uint32(env.OriginatorNodeID())
		seqID := env.OriginatorSequenceID()

		lastSeq, seen := vc[origID]
		if seen && lastSeq >= seqID {
			// Already seen.
			continue
		}

		result = append(result, env.Proto())
		vc[origID] = seqID
	}

	return result
}

func newSubscriptionStatusMessage(
	status message_api.SubscribeTopicsResponse_SubscriptionStatus,
) *message_api.SubscribeTopicsResponse {
	return &message_api.SubscribeTopicsResponse{
		Response: &message_api.SubscribeTopicsResponse_StatusUpdate_{
			StatusUpdate: &message_api.SubscribeTopicsResponse_StatusUpdate{
				Status: status,
			},
		},
	}
}

func newEnvelopesMessage(
	envs []*envelopesProto.OriginatorEnvelope,
) *message_api.SubscribeTopicsResponse {
	return &message_api.SubscribeTopicsResponse{
		Response: &message_api.SubscribeTopicsResponse_Envelopes_{
			Envelopes: &message_api.SubscribeTopicsResponse_Envelopes{
				Envelopes: envs,
			},
		},
	}
}

// sendTopicEnvelopes sends the given envelopes to the stream.
// No-ops if the slice is empty.
func (s *Service) sendTopicEnvelopes(
	stream *connect.ServerStream[message_api.SubscribeTopicsResponse],
	envs []*envelopesProto.OriginatorEnvelope,
) error {
	if len(envs) == 0 {
		return nil
	}

	err := stream.Send(newEnvelopesMessage(envs))
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error sending envelopes: %w", err),
		)
	}

	return nil
}

// catchUpTopics performs the catch-up phase for topics that have cursors.
func (s *Service) catchUpTopics(
	ctx context.Context,
	stream *connect.ServerStream[message_api.SubscribeTopicsResponse],
	cursors db.TopicCursors,
	catchUpKeys []string,
	logger *zap.Logger,
) error {
	if len(catchUpKeys) == 0 {
		logger.Debug("skipping catch up, no topics with cursors")
		return nil
	}

	allOriginators, err := s.originatorList.GetOriginatorNodeIDs(ctx)
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not get originator list: %w", err),
		)
	}

	fillMissingOriginatorsForTopics(cursors, catchUpKeys, allOriginators)

	chunks := utils.ChunkSlice(catchUpKeys, maxTopicsPerChunk)

	for _, chunkKeys := range chunks {
		rowsPerEntry := db.CalculateRowsPerEntry(len(chunkKeys), topicPageLimit)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Build sub-cursor map for this chunk from the shared cursors.
			subCursors := make(db.TopicCursors, len(chunkKeys))
			for _, key := range chunkKeys {
				subCursors[key] = cursors[key]
			}

			rows, err := s.fetchTopicEnvelopesWithRetry(
				ctx, subCursors, topicPageLimit, rowsPerEntry,
			)
			if err != nil {
				return err
			}

			if s.logger.Core().Enabled(zap.DebugLevel) {
				logger.Debug("topic catch-up fetched envelopes", utils.CountField(int64(len(rows))))
			}

			envs := unmarshalEnvelopes(rows, s.logger)

			// Advance cursors so the next query page starts after these envelopes.
			envsToSend := advanceTopicCursors(cursors, envs, logger)
			err = s.sendTopicEnvelopes(stream, envsToSend)
			if err != nil {
				return err
			}

			// Compare against rowsPerEntry, not topicPageLimit. The LATERAL query
			// distributes topicPageLimit across (topic, originator) pairs via
			// per-entry sub-limits, so total rows returned can be less than
			// topicPageLimit even when more data exists. Using topicPageLimit
			// here would cause premature termination.
			if int32(len(rows)) < rowsPerEntry {
				break
			}
		}
	}

	return nil
}

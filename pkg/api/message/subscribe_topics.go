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

// SubscribeTopicEnvelopes implements the per-topic cursor subscribe API.
func (s *Service) SubscribeTopicEnvelopes(
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

	// Send a keepalive immediately, so wasm based clients maintain the connection open.
	err := stream.Send(&message_api.SubscribeTopicsResponse{})
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not send keepalive: %w", err),
		)
	}

	filters := req.Msg.GetFilters()
	if err := s.validateTopicFilters(ctx, filters); err != nil {
		return err
	}

	cursors, catchUpKeys := buildTopicCursors(filters)

	// Build a synthetic EnvelopesQuery for the listener using deduplicated topics.
	topics := make([][]byte, 0, len(cursors))
	for key := range cursors {
		topics = append(topics, []byte(key))
	}
	syntheticQuery := &message_api.EnvelopesQuery{
		Topics: topics,
	}
	envelopesCh := s.subscribeWorker.listen(ctx, syntheticQuery)

	err = s.catchUpTopics(ctx, stream, cursors, catchUpKeys, logger)
	if err != nil {
		return err
	}

	// GRPC keep-alives are not sufficient in some load balanced environments.
	ticker := time.NewTicker(s.options.SendKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = stream.Send(&message_api.SubscribeTopicsResponse{})
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
			envsToSend := advanceTopicCursors(cursors, envs)
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
func (s *Service) validateTopicFilters(
	ctx context.Context,
	filters []*message_api.SubscribeTopicsRequest_TopicFilter,
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

	// Collect all originator IDs referenced by cursors so we can validate them.
	needOriginatorValidation := false
	referencedOriginators := make(map[uint32]struct{})

	for _, f := range filters {
		topicBytes := f.GetTopic()
		if len(topicBytes) == 0 || len(topicBytes) > maxTopicLength {
			return connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("invalid topic length: %d", len(topicBytes)),
			)
		}

		vc := f.GetLastSeen().GetNodeIdToSequenceId()
		if len(vc) > maxVectorClockLength {
			return connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("vector clock length exceeds maximum of %d", maxVectorClockLength),
			)
		}
		if len(vc) > 0 {
			needOriginatorValidation = true
			for origID := range vc {
				referencedOriginators[origID] = struct{}{}
			}
		}
	}

	// Validate that all referenced originator IDs are known.
	if needOriginatorValidation {
		allOriginators, err := s.originatorList.GetOriginatorNodeIDs(ctx)
		if err != nil {
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("could not get originator list: %w", err),
			)
		}
		known := make(map[uint32]struct{}, len(allOriginators))
		for _, id := range allOriginators {
			known[uint32(id)] = struct{}{}
		}
		for origID := range referencedOriginators {
			if _, ok := known[origID]; !ok {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("unknown originator node ID in cursor: %d", origID),
				)
			}
		}
	}

	return nil
}

// buildTopicCursors converts topic filters into a TopicCursors map and
// returns the keys that need catch-up (those with non-nil LastSeen).
func buildTopicCursors(
	filters []*message_api.SubscribeTopicsRequest_TopicFilter,
) (db.TopicCursors, []string) {
	cursors := make(db.TopicCursors, len(filters))
	catchUpKeys := make([]string, 0, len(filters))

	for _, f := range filters {
		key := string(f.GetTopic())

		if _, exists := cursors[key]; exists {
			continue
		}

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

	return cursors, catchUpKeys
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
) []*envelopesProto.OriginatorEnvelope {
	result := make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))

	for _, env := range envs {
		vc, ok := cursors[string(env.TargetTopic().Bytes())]
		if !ok {
			// Not subscribed to this topic.
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

// sendTopicEnvelopes sends the given envelopes to the stream.
// No-ops if the slice is empty.
func (s *Service) sendTopicEnvelopes(
	stream *connect.ServerStream[message_api.SubscribeTopicsResponse],
	envs []*envelopesProto.OriginatorEnvelope,
) error {
	if len(envs) == 0 {
		return nil
	}

	err := stream.Send(&message_api.SubscribeTopicsResponse{
		Envelopes: envs,
	})
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
			envsToSend := advanceTopicCursors(cursors, envs)
			err = s.sendTopicEnvelopes(stream, envsToSend)
			if err != nil {
				return err
			}

			// If fewer rows than rowsPerEntry, chunk is exhausted.
			if int32(len(rows)) < rowsPerEntry {
				break
			}
		}
	}

	return nil
}

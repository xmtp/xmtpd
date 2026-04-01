package message

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// batchAndSendEnvelopes handles cursor-based deduplication, batching by gRPC
// payload size, and sending via the provided flushFn callback. It advances the
// query's LastSeen cursor as envelopes are processed.
//
// This is the shared core used by both sendEnvelopes (SubscribeEnvelopes) and
// sendOriginatorEnvelopes (SubscribeOriginators).
func batchAndSendEnvelopes(
	query *message_api.EnvelopesQuery,
	envs []*envelopes.OriginatorEnvelope,
	flushFn func(batch []*envelopesProto.OriginatorEnvelope) error,
) error {
	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen = &envelopesProto.Cursor{
			NodeIdToSequenceId: cursor,
		}
	}

	var (
		batch          = make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))
		batchWireBytes = 0
	)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := flushFn(batch); err != nil {
			return err
		}
		batchWireBytes = 0
		batch = batch[:0]
		return nil
	}

	for _, env := range envs {
		var (
			origID = env.OriginatorNodeID()
			seqID  = env.OriginatorSequenceID()
		)

		// Skip if we've already seen this envelope.
		if cursor[origID] >= seqID {
			continue
		}

		var (
			envProto     = env.Proto()
			envProtoSize = proto.Size(envProto)
			envWireSize  = envProtoSize + envelopeOverhead(uint64(envProtoSize))
		)

		if len(batch) > 0 && batchWireBytes+envWireSize > constants.GRPCPayloadLimit {
			if err := flush(); err != nil {
				return err
			}
		}

		batch = append(batch, envProto)
		batchWireBytes += envWireSize
		cursor[origID] = seqID
	}

	return flush()
}

// catchUpWithSendFn performs the cursor-based catch-up phase, paginating
// through the database and sending batches via sendFn. It is the shared core
// for both catchUpFromCursor (SubscribeEnvelopes) and catchUpOriginators
// (SubscribeOriginators).
func (s *Service) catchUpWithSendFn(
	ctx context.Context,
	query *message_api.EnvelopesQuery,
	logger *zap.Logger,
	sendFn func(envs []*envelopes.OriginatorEnvelope) error,
) error {
	if query.GetLastSeen() == nil {
		logger.Debug("skipping catch up")
		return nil
	}

	cursor := query.GetLastSeen().GetNodeIdToSequenceId()
	if cursor == nil {
		cursor = make(map[uint32]uint64)
		query.LastSeen.NodeIdToSequenceId = cursor
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("catching up from cursor", utils.BodyField(cursor))
	}

	for {
		rows, err := s.fetchEnvelopesWithRetry(ctx, query, maxRequestedRows)
		if err != nil {
			return err
		}

		if s.logger.Core().Enabled(zap.DebugLevel) {
			logger.Debug("fetched envelopes", utils.CountField(int64(len(rows))))
		}

		envs := unmarshalEnvelopes(rows, s.logger)

		if err := sendFn(envs); err != nil {
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error sending envelopes: %w", err),
			)
		}

		if len(rows) < int(maxRequestedRows) {
			break
		}
		time.Sleep(pagingInterval)
	}

	return nil
}

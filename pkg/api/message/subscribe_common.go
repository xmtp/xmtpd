package message

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type catchUpMode int

const (
	catchUpNone       catchUpMode = iota // stream new envelopes only
	catchUpFromStart                     // catch up from the very beginning
	catchUpFromCursor                    // catch up from specified cursor position
)

// subscribeFilter is the internal representation of a subscription filter,
// decoupled from the deprecated message_api.EnvelopesQuery proto.
type subscribeFilter struct {
	topics            [][]byte
	originatorNodeIDs []uint32
	catchUpMode       catchUpMode
	cursor            map[uint32]uint64 // always initialized; starting position + progress tracker
}

func cursorFromProto(c *envelopesProto.Cursor) map[uint32]uint64 {
	if c == nil || c.GetNodeIdToSequenceId() == nil {
		return make(map[uint32]uint64)
	}
	return c.GetNodeIdToSequenceId()
}

// batchAndSendEnvelopes handles cursor-based deduplication, batching by gRPC
// payload size, and sending via the provided flushFn callback. It advances the
// cursor as envelopes are processed.
//
// This is the shared core used by both sendEnvelopes (SubscribeEnvelopes) and
// sendOriginatorsResponse (SubscribeOriginators).
func batchAndSendEnvelopes(
	cursor map[uint32]uint64,
	envs []*envelopes.OriginatorEnvelope,
	flushFn func(batch []*envelopesProto.OriginatorEnvelope) error,
) error {
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
	query *subscribeFilter,
	logger *zap.Logger,
	sendFn func(envs []*envelopes.OriginatorEnvelope) error,
) error {
	switch query.catchUpMode {
	case catchUpNone:
		logger.Debug("skipping catch up")
		return nil
	case catchUpFromStart, catchUpFromCursor:
		// fall through to perform catch-up
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug("catching up from cursor", utils.BodyField(query.cursor))
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

package message

import (
	"context"
	"fmt"
	"maps"
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
	src := c.GetNodeIdToSequenceId()
	if src == nil {
		return make(map[uint32]uint64)
	}
	// Clone to avoid mutating the proto's internal map.
	m := make(map[uint32]uint64, len(src))
	maps.Copy(m, src)
	return m
}

// originatorResponseOverhead is the fixed per-batch wire overhead from the
// SubscribeOriginatorsResponse wrapper. The oneof + nested message adds ~10
// bytes per batch (2 tag+length pairs with varints). SubscribeEnvelopesResponse
// has envelopes directly on the message (0 overhead).
const originatorResponseOverhead = 10

// batchAndSendEnvelopes handles cursor-based deduplication, batching by gRPC
// payload size, and sending via the provided flushFn callback. It advances the
// cursor as envelopes are processed.
func batchAndSendEnvelopes(
	logger *zap.Logger,
	cursor map[uint32]uint64,
	envs []*envelopes.OriginatorEnvelope,
	wrapperOverhead int,
	flushFn func(batch []*envelopesProto.OriginatorEnvelope) error,
) error {
	maxEnvelopeSize := constants.GRPCPayloadLimit - wrapperOverhead

	var (
		batch          = make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))
		batchWireBytes = wrapperOverhead
	)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := flushFn(batch); err != nil {
			return err
		}
		batchWireBytes = wrapperOverhead
		batch = batch[:0]
		return nil
	}

	for _, env := range envs {
		var (
			origID = env.OriginatorNodeID()
			seqID  = env.OriginatorSequenceID()
		)

		if cursor[origID] >= seqID {
			continue
		}

		var (
			envProto     = env.Proto()
			envProtoSize = proto.Size(envProto)
			envWireSize  = envProtoSize + envelopeOverhead(uint64(envProtoSize))
		)

		// A single envelope that exceeds the gRPC payload limit cannot be sent.
		// Skip it and advance the cursor so pagination doesn't get stuck.
		if envWireSize > maxEnvelopeSize {
			logger.Warn(
				"skipping oversized envelope",
				zap.Uint32("originator_node_id", origID),
				zap.Uint64("originator_sequence_id", seqID),
				zap.Int("wire_bytes", envWireSize),
				zap.Int("limit", maxEnvelopeSize),
			)
			cursor[origID] = seqID
			continue
		}

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
	case catchUpFromCursor:
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

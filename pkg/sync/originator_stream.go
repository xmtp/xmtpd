package sync

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/cenkalti/backoff/v5"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type originatorStream struct {
	ctx                  context.Context
	logger               *zap.Logger
	node                 *registry.Node
	lastSequenceIds      map[uint32]uint64
	permittedOriginators map[uint32]struct{}
	stream               message_api.ReplicationApi_SubscribeEnvelopesClient
	writeQueue           chan *envUtils.OriginatorEnvelope
}

func newOriginatorStream(
	ctx context.Context,
	logger *zap.Logger,
	node *registry.Node,
	lastSequenceIds map[uint32]uint64,
	permittedOriginators map[uint32]struct{},
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	writeQueue chan *envUtils.OriginatorEnvelope,
) *originatorStream {
	return &originatorStream{
		ctx: ctx,
		logger: logger.With(
			utils.OriginatorIDField(node.NodeID),
			utils.NodeHTTPAddressField(node.HTTPAddress),
		),
		node:                 node,
		lastSequenceIds:      lastSequenceIds,
		permittedOriginators: permittedOriginators,
		stream:               stream,
		writeQueue:           writeQueue,
	}
}

func (s *originatorStream) listen() error {
	var (
		recvChan = make(chan *message_api.SubscribeEnvelopesResponse)
		errChan  = make(chan error, 1)
	)

	// Reader routine, responsible for reading from a blocking GRPC channel
	// TODO: Use tracing.GoWrap and waitgroup.
	go func() {
		for {
			envs, err := s.stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			recvChan <- envs
		}
	}()

	// main routine, responsible for processing and validating messages
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("context canceled, stopping stream listener")
			return backoff.Permanent(s.ctx.Err())

		case envs, ok := <-recvChan:
			if !ok {
				s.logger.Error("recvChan is closed")
				return backoff.Permanent(errors.New("recvChan is closed"))
			}

			if envs == nil || len(envs.GetEnvelopes()) == 0 {
				continue
			}

			// Create span for processing this batch of envelopes
			batchSpan := tracing.StartSpan(tracing.SpanSyncReceiveBatch)
			tracing.SpanTag(batchSpan, "source_node", s.node.NodeID)
			tracing.SpanTag(batchSpan, "num_envelopes", len(envs.Envelopes))

			s.logger.Debug(
				"received envelopes",
				utils.NumEnvelopesField(len(envs.GetEnvelopes())),
			)

			validCount := 0
			invalidCount := 0
			for _, env := range envs.GetEnvelopes() {
				// Any message that fails validation here will be dropped permanently
				parsedEnv, err := s.validateEnvelope(env)
				if err != nil {
					s.logger.Error("discarding envelope after validation failed", zap.Error(err))
					invalidCount++
					continue
				}
				validCount++
				s.writeQueue <- parsedEnv
			}

			tracing.SpanTag(batchSpan, "valid_count", validCount)
			tracing.SpanTag(batchSpan, "invalid_count", invalidCount)
			batchSpan.Finish()

		case err, ok := <-errChan:
			if !ok {
				s.logger.Error("errChan is closed")
				return backoff.Permanent(errors.New("errChan is closed"))
			}

			if errors.Is(err, io.EOF) {
				s.logger.Info("stream closed with EOF")
				// reset backoff to 1 second
				return backoff.RetryAfter(1)
			}
			s.logger.Error(
				"stream closed with error",
				zap.Error(err),
			)

			if strings.Contains(err.Error(), "is not compatible") {
				// the node won't accept our version
				// try again in an hour in case their config has changed
				return backoff.RetryAfter(3600)
			}

			// keep existing backoff
			return err
		}
	}
}

// validateEnvelope performs all static validation on an envelope
// if an error is encountered, the envelope will be dropped and the stream will continue
func (s *originatorStream) validateEnvelope(
	envProto *envelopes.OriginatorEnvelope,
) (*envUtils.OriginatorEnvelope, error) {
	// Create span for envelope validation
	span := tracing.StartSpan(tracing.SpanSyncValidateEnvelope)
	defer span.Finish()

	tracing.SpanTag(span, "source_node", s.node.NodeID)

	var err error
	defer func() {
		if err != nil {
			metrics.EmitSyncOriginatorErrorMessages(s.node.NodeID, 1)
			span.Finish(tracing.WithError(err))
		}
	}()

	var env *envUtils.OriginatorEnvelope
	env, err = envUtils.NewOriginatorEnvelope(envProto)
	if err != nil {
		s.logger.Error("failed to unmarshal originator envelope", zap.Error(err))
		return nil, err
	}

	// Add envelope details to span
	tracing.SpanTag(span, "sequence_id", env.OriginatorSequenceID())
	tracing.SpanTag(span, "topic", hex.EncodeToString(env.TargetTopic().Bytes()))

	// TODO:(nm) Handle fetching envelopes from other nodes
	originatorID := env.OriginatorNodeID()
	seqID := env.OriginatorSequenceID()
	if _, permitted := s.permittedOriginators[originatorID]; !permitted {
		permittedIDs := make([]uint32, 0, len(s.permittedOriginators))
		for id := range s.permittedOriginators {
			permittedIDs = append(permittedIDs, id)
		}
		slices.Sort(permittedIDs)

		err = fmt.Errorf(
			"invalid envelope originator: got=%d permitted=%v",
			originatorID,
			permittedIDs,
		)
		tracing.SpanTag(span, "wrong_originator", env.OriginatorNodeID())

		s.logger.Error("received envelope from invalid originator",
			zap.Uint32("originator_id", originatorID),
			zap.Uint32s("permitted_originator_ids", permittedIDs),
			zap.Error(err),
		)

		return nil, err
	}

	metrics.EmitSyncLastSeenOriginatorSequenceID(env.OriginatorNodeID(), env.OriginatorSequenceID())
	metrics.EmitSyncOriginatorReceivedMessagesCount(env.OriginatorNodeID(), 1)

	lastSeq := s.lastSequenceIds[originatorID]

	if seqID != lastSeq+1 {
		s.logger.Error(
			"received out-of-order envelope",
			utils.OriginatorIDField(originatorID),
			utils.SequenceIDField(int64(seqID)),
			zap.Uint64("expected_sequence_id", lastSeq+1),
		)
		tracing.SpanTag(span, "out_of_order", true)
		tracing.SpanTag(span, "expected_sequence_id", s.lastSequenceId+1)
	}

	if seqID > lastSeq {
		s.lastSequenceIds[originatorID] = seqID
	}

	// Validate that there is a valid payer signature
	_, err = env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		s.logger.Error("failed to recover payer address", zap.Error(err))
		return nil, err
	}

	return env, nil
}

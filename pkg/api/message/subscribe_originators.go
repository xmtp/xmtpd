package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Service) SubscribeOriginators(
	ctx context.Context,
	req *connect.Request[message_api.SubscribeOriginatorsRequest],
	stream *connect.ServerStream[message_api.SubscribeOriginatorsResponse],
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

	filter := req.Msg.GetFilter()
	if filter == nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("filter must not be nil"),
		)
	}

	if len(filter.GetOriginatorNodeIds()) == 0 {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("filter must contain at least one originator node id"),
		)
	}

	query := &message_api.EnvelopesQuery{
		OriginatorNodeIds: filter.GetOriginatorNodeIds(),
		LastSeen:          filter.GetLastSeen(),
	}

	if err := s.validateQuery(query, false); err != nil {
		return connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("invalid subscription request: %w", err),
		)
	}

	// Send a keepalive immediately so wasm-based clients maintain the connection open.
	if err := stream.Send(&message_api.SubscribeOriginatorsResponse{}); err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("could not send keepalive: %w", err),
		)
	}

	ch := s.subscribeWorker.listen(ctx, query)

	if err := s.catchUpOriginators(ctx, stream, query, logger); err != nil {
		return err
	}

	ticker := time.NewTicker(s.options.SendKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := stream.Send(&message_api.SubscribeOriginatorsResponse{}); err != nil {
				return connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("could not send keepalive: %w", err),
				)
			}

		case envs, open := <-ch:
			ticker.Reset(s.options.SendKeepAliveInterval)

			if !open {
				logger.Debug("channel closed by worker")
				return nil
			}

			if err := s.sendOriginatorEnvelopes(stream, query, envs); err != nil {
				return connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("error sending envelope: %w", err),
				)
			}

		case <-ctx.Done():
			logger.Debug("originator subscription stream closed")
			return nil

		case <-s.ctx.Done():
			logger.Debug("message service closed")
			return nil
		}
	}
}

func (s *Service) catchUpOriginators(
	ctx context.Context,
	stream *connect.ServerStream[message_api.SubscribeOriginatorsResponse],
	query *message_api.EnvelopesQuery,
	logger *zap.Logger,
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

		if err := s.sendOriginatorEnvelopes(stream, query, envs); err != nil {
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

func (s *Service) sendOriginatorEnvelopes(
	stream *connect.ServerStream[message_api.SubscribeOriginatorsResponse],
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

	var (
		batch          = make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))
		batchWireBytes = 0
	)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}

		if err := stream.Send(&message_api.SubscribeOriginatorsResponse{
			Response: &message_api.SubscribeOriginatorsResponse_Envelopes_{
				Envelopes: &message_api.SubscribeOriginatorsResponse_Envelopes{
					Envelopes: batch,
				},
			},
		}); err != nil {
			return connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error sending envelopes: %w", err),
			)
		}

		metrics.EmitAPIOutgoingEnvelopes(stream.Conn().Spec().Procedure, len(batch))

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

		// If the batch is not empty and the total bytes would exceed the limit, flush first.
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

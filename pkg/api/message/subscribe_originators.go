package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
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
	if err := validateOriginatorFilter(filter); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	query := &subscribeFilter{
		originatorNodeIDs: filter.GetOriginatorNodeIds(),
		catchUpMode:       catchUpFromCursor,
		cursor:            cursorFromProto(filter.GetLastSeen()),
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

	sendFn := func(envs []*envelopes.OriginatorEnvelope) error {
		return s.sendOriginatorsResponse(stream, query.cursor, envs)
	}

	if err := s.catchUpWithSendFn(ctx, query, logger, sendFn); err != nil {
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

			if err := s.sendOriginatorsResponse(stream, query.cursor, envs); err != nil {
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

// validateOriginatorFilter validates the filter for SubscribeOriginators.
// Returns a plain error (caller wraps with connect.CodeInvalidArgument).
func validateOriginatorFilter(
	filter *message_api.SubscribeOriginatorsRequest_OriginatorFilter,
) error {
	if filter == nil {
		return errors.New("filter must not be nil")
	}
	if len(filter.GetOriginatorNodeIds()) == 0 {
		return errors.New("filter must contain at least one originator node id")
	}
	if filter.GetLastSeen() == nil {
		return errors.New("last_seen cursor is required")
	}
	return nil
}

func (s *Service) sendOriginatorsResponse(
	stream *connect.ServerStream[message_api.SubscribeOriginatorsResponse],
	cursor map[uint32]uint64,
	envs []*envelopes.OriginatorEnvelope,
) error {
	return batchAndSendEnvelopes(
		s.logger,
		cursor,
		envs,
		originatorResponseOverhead,
		func(batch []*envelopesProto.OriginatorEnvelope) error {
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
			return nil
		},
	)
}

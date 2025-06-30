package sync

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/cenkalti/backoff/v5"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type originatorStream struct {
	ctx          context.Context
	log          *zap.Logger
	node         *registry.Node
	lastEnvelope *envUtils.OriginatorEnvelope
	stream       message_api.ReplicationApi_SubscribeEnvelopesClient
	writeQueue   chan *envUtils.OriginatorEnvelope
}

func newOriginatorStream(
	ctx context.Context,
	log *zap.Logger,
	node *registry.Node,
	lastEnvelope *envUtils.OriginatorEnvelope,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	writeQueue chan *envUtils.OriginatorEnvelope,
) *originatorStream {
	return &originatorStream{
		ctx: ctx,
		log: log.With(
			zap.Uint32("originator_id", node.NodeID),
			zap.String("http_address", node.HttpAddress),
		),
		node:         node,
		lastEnvelope: lastEnvelope,
		stream:       stream,
		writeQueue:   writeQueue,
	}
}

func (s *originatorStream) listen() error {
	recvChan := make(chan *message_api.SubscribeEnvelopesResponse)
	errChan := make(chan error, 1)

	// Reader routine, responsible for reading from a blocking GRPC channel
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
			s.log.Info("Context canceled, stopping stream listener")
			return backoff.Permanent(s.ctx.Err())

		case envs := <-recvChan:
			if envs == nil || len(envs.Envelopes) == 0 {
				continue
			}
			s.log.Debug(
				"Received envelopes",
				zap.Any("numEnvelopes", len(envs.Envelopes)),
			)

			for _, env := range envs.Envelopes {
				// Any message that fails validation here will be dropped permanently
				parsedEnv, err := s.validateEnvelope(env)
				if err != nil {
					s.log.Error("discarding envelope after validation failed", zap.Error(err))
					continue
				}
				s.writeQueue <- parsedEnv
			}

		case err := <-errChan:
			if err == io.EOF {
				s.log.Info("Stream closed with EOF")
				// reset backoff to 1 second
				return backoff.RetryAfter(1)
			}
			s.log.Error(
				"Stream closed with error",
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
	var err error
	defer func() {
		if err != nil {
			metrics.EmitSyncOriginatorErrorMessages(s.node.NodeID, 1)
		}
	}()

	var env *envUtils.OriginatorEnvelope
	env, err = envUtils.NewOriginatorEnvelope(envProto)
	if err != nil {
		s.log.Error("Failed to unmarshal originator envelope", zap.Error(err))
		return nil, err
	}

	// TODO:(nm) Handle fetching envelopes from other nodes
	if env.OriginatorNodeID() != s.node.NodeID {
		s.log.Error("Received envelope from wrong node",
			zap.Any("nodeID", env.OriginatorNodeID()),
			zap.Any("expectedNodeId", s.node.NodeID),
		)
		err = errors.New("originator ID does not match envelope")
		return nil, err
	}

	metrics.EmitSyncLastSeenOriginatorSequenceId(env.OriginatorNodeID(), env.OriginatorSequenceID())
	metrics.EmitSyncOriginatorReceivedMessagesCount(env.OriginatorNodeID(), 1)

	var lastSequenceID uint64 = 0
	var lastNs int64 = 0
	if s.lastEnvelope != nil {
		lastSequenceID = s.lastEnvelope.OriginatorSequenceID()
		lastNs = s.lastEnvelope.OriginatorNs()
	}

	if env.OriginatorSequenceID() != lastSequenceID+1 || env.OriginatorNs() < lastNs {
		// TODO(rich) Submit misbehavior report and continue
		s.log.Error(
			"Received out of order envelope",
			zap.Any("envelope", env),
			zap.Any("lastEnvelope", s.lastEnvelope),
		)
	}

	if env.OriginatorSequenceID() > lastSequenceID {
		s.lastEnvelope = env
	}

	// Validate that there is a valid payer signature
	_, err = env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		s.log.Error("Failed to recover payer address", zap.Error(err))
		return nil, err
	}

	return env, nil
}

package api

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

	ctx        context.Context
	log        *zap.Logger
	registrant *registrant.Registrant
	store      *sql.DB
	worker     *PublishWorker
}

func NewReplicationApiService(
	ctx context.Context,
	log *zap.Logger,
	registrant *registrant.Registrant,
	store *sql.DB,
) (*Service, error) {
	worker, err := StartPublishWorker(ctx, log, registrant, store)
	if err != nil {
		return nil, err
	}
	return &Service{
		ctx:        ctx,
		log:        log,
		registrant: registrant,
		store:      store,
		worker:     worker,
	}, nil
}

func (s *Service) Close() {
	s.log.Info("closed")
}

func (s *Service) BatchSubscribeEnvelopes(
	req *message_api.BatchSubscribeEnvelopesRequest,
	server message_api.ReplicationApi_BatchSubscribeEnvelopesServer,
) error {
	return status.Errorf(codes.Unimplemented, "method BatchSubscribeEnvelopes not implemented")
}

func (s *Service) QueryEnvelopes(
	ctx context.Context,
	req *message_api.QueryEnvelopesRequest,
) (*message_api.QueryEnvelopesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryEnvelopes not implemented")
}

func (s *Service) PublishEnvelope(
	ctx context.Context,
	req *message_api.PublishEnvelopeRequest,
) (*message_api.PublishEnvelopeResponse, error) {
	clientEnv, err := s.validatePayerInfo(req.GetPayerEnvelope())
	if err != nil {
		return nil, err
	}

	topic, err := s.validateClientInfo(clientEnv)
	if err != nil {
		return nil, err
	}

	// TODO(rich): If it is a commit, publish it to blockchain instead

	payerBytes, err := proto.Marshal(req.GetPayerEnvelope())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal envelope: %v", err)
	}

	stagedEnv, err := queries.New(s.store).
		InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         topic,
			PayerEnvelope: payerBytes,
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert staged envelope: %v", err)
	}
	s.worker.NotifyStagedPublish()

	originatorEnv, err := s.registrant.SignStagedEnvelope(stagedEnv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not sign envelope: %v", err)
	}

	return &message_api.PublishEnvelopeResponse{OriginatorEnvelope: originatorEnv}, nil
}

func (s *Service) validatePayerInfo(
	payerEnv *message_api.PayerEnvelope,
) (*message_api.ClientEnvelope, error) {
	clientBytes := payerEnv.GetUnsignedClientEnvelope()
	sig := payerEnv.GetPayerSignature()
	if (clientBytes == nil) || (sig == nil) {
		return nil, status.Errorf(codes.InvalidArgument, "missing envelope or signature")
	}
	// TODO(rich): Verify payer signature

	clientEnv := &message_api.ClientEnvelope{}
	err := proto.Unmarshal(clientBytes, clientEnv)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"could not unmarshal client envelope: %v",
			err,
		)
	}

	return clientEnv, nil
}

func (s *Service) validateClientInfo(clientEnv *message_api.ClientEnvelope) ([]byte, error) {
	if clientEnv.GetAad().GetTargetOriginator() != uint32(s.registrant.NodeID()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target originator")
	}

	topic := clientEnv.GetAad().GetTargetTopic()
	if len(topic) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing target topic")
	}

	// TODO(rich): Verify all originators have synced past `last_originator_sids`
	// TODO(rich): Check that the blockchain sequence ID is equal to the latest on the group
	// TODO(rich): Perform any payload-specific validation (e.g. identity updates)

	return topic, nil
}

package api

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/node"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

	ctx     context.Context
	log     *zap.Logger
	node    *node.Node
	queries *queries.Queries
}

func NewReplicationApiService(
	ctx context.Context,
	log *zap.Logger,
	node *node.Node,
	writerDB *sql.DB,
) (*Service, error) {
	return &Service{ctx: ctx, log: log, node: node, queries: queries.New(writerDB)}, nil
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
	payerEnv := req.GetPayerEnvelope()
	clientBytes := payerEnv.GetUnsignedClientEnvelope()
	sig := payerEnv.GetPayerSignature()
	if (clientBytes == nil) || (sig == nil) {
		return nil, status.Errorf(codes.InvalidArgument, "missing envelope or signature")
	}
	// TODO(rich): Verify payer signature
	// TODO(rich): Verify all originators have synced past `last_originator_sids`
	// TODO(rich): Check that the blockchain sequence ID is equal to the latest on the group
	// TODO(rich): Perform any payload-specific validation (e.g. identity updates)
	// TODO(rich): If it is a commit, publish it to blockchain instead

	payerBytes, err := proto.Marshal(payerEnv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal envelope: %v", err)
	}

	stagedEnv, err := s.queries.InsertStagedOriginatorEnvelope(ctx, payerBytes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert staged envelope: %v", err)
	}

	originatorEnv, err := utils.SignStagedEnvelope(stagedEnv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not sign envelope: %v", err)
	}

	return &message_api.PublishEnvelopeResponse{OriginatorEnvelope: originatorEnv}, nil
}

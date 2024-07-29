package api

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
)

type Service struct {
	message_api.UnimplementedReplicationApiServer

	ctx context.Context
	log *zap.Logger
}

func NewReplicationApiService(ctx context.Context, log *zap.Logger) (message_api.ReplicationApiServer, error) {
	return &Service{ctx: ctx, log: log}, nil
}

func (s *Service) SubscribeEnvelopes(req *message_api.BatchSubscribeEnvelopesRequest, server message_api.ReplicationApi_SubscribeEnvelopesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeEnvelopes not implemented")
}

func (s *Service) QueryEnvelopes(ctx context.Context, req *message_api.QueryEnvelopesRequest) (*message_api.QueryEnvelopesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryEnvelopes not implemented")
}

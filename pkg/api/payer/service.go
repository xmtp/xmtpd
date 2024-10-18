package payer

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	payer_api.UnimplementedPayerApiServer

	ctx context.Context
	log *zap.Logger
}

func NewPayerApiService(ctx context.Context, log *zap.Logger) (*Service, error) {
	return &Service{
		ctx: ctx,
		log: log,
	}, nil
}

func (s *Service) PublishClientEnvelopes(
	ctx context.Context,
	req *payer_api.PublishClientEnvelopesRequest,
) (*payer_api.PublishClientEnvelopesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishClientEnvelopes not implemented")
}

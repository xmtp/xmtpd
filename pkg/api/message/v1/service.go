package messagev1

import (
	"context"
	"errors"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	ErrTODO = errors.New("TODO")
)

type Service struct {
	messagev1.UnimplementedMessageApiServer

	// Configured as constructor options.
	log *zap.Logger

	// Configured internally.
	ctx       context.Context
	ctxCancel func()
}

func NewService(log *zap.Logger) (s *Service, err error) {
	s = &Service{
		log: log.Named("message/v1"),
	}
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	return s, nil
}

func (s *Service) Close() {
	if s.ctxCancel != nil {
		s.ctxCancel()
	}
}

func (s *Service) Publish(ctx context.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	return &messagev1.PublishResponse{}, ErrTODO
}

func (s *Service) Subscribe(req *messagev1.SubscribeRequest, stream messagev1.MessageApi_SubscribeServer) error {
	return ErrTODO
}

func (s *Service) SubscribeAll(req *messagev1.SubscribeAllRequest, stream messagev1.MessageApi_SubscribeAllServer) error {
	return ErrTODO
}

func (s *Service) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	return &messagev1.QueryResponse{}, ErrTODO
}

func (s *Service) BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	return &messagev1.BatchQueryResponse{}, ErrTODO
}

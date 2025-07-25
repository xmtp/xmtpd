package metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service struct {
	metadata_api.UnimplementedMetadataApiServer
	ctx     context.Context
	log     *zap.Logger
	cu      CursorUpdater
	version *semver.Version
}

func NewMetadataApiService(
	ctx context.Context,
	log *zap.Logger,
	updater CursorUpdater,
	version *semver.Version,
) (*Service, error) {
	return &Service{
		ctx:     ctx,
		log:     log,
		cu:      updater,
		version: version,
	}, nil
}

func (s *Service) GetSyncCursor(
	_ context.Context,
	_ *metadata_api.GetSyncCursorRequest,
) (*metadata_api.GetSyncCursorResponse, error) {
	return &metadata_api.GetSyncCursorResponse{
		LatestSync: s.cu.GetCursor(),
	}, nil
}

func (s *Service) SubscribeSyncCursor(
	_ *metadata_api.GetSyncCursorRequest,
	stream metadata_api.MetadataApi_SubscribeSyncCursorServer,
) error {
	err := stream.SendHeader(metadata.Pairs("subscribed", "true"))
	if err != nil {
		return status.Errorf(codes.Internal, "could not send header: %v", err)
	}

	// send the initial cursor
	// the subscriber will only send a new message if there was a change
	cursor := s.cu.GetCursor()
	err = stream.Send(&metadata_api.GetSyncCursorResponse{
		LatestSync: cursor,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error sending cursor: %v", err)
	}

	clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
	updateChan := make(chan struct{}, 1)
	s.cu.AddSubscriber(clientID, updateChan)
	defer s.cu.RemoveSubscriber(clientID)

	for {
		select {
		case _, open := <-updateChan:
			if open {
				cursor := s.cu.GetCursor()
				err := stream.Send(&metadata_api.GetSyncCursorResponse{
					LatestSync: cursor,
				})
				if err != nil {
					return status.Errorf(codes.Internal, "error sending cursor: %v", err)
				}
			} else {
				s.log.Debug("channel closed by worker")
				return nil
			}
		case <-stream.Context().Done():
			s.log.Debug("stream closed")
			return nil
		case <-s.ctx.Done():
			s.log.Debug("service closed")
			return nil
		}
	}
}

func (s *Service) GetVersion(
	_ context.Context,
	_ *metadata_api.GetVersionRequest,
) (*metadata_api.GetVersionResponse, error) {
	if s.version == nil {
		s.log.Error("version is not set")
		return nil, status.Errorf(codes.Internal, "version is not set")
	}

	return &metadata_api.GetVersionResponse{
		Version: s.version.String(),
	}, nil
}

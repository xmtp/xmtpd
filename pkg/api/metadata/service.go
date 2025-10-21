// Package metadata implements the metadata API service.
package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	getSyncCursorMethod       = "GetSyncCursor"
	subscribeSyncCursorMethod = "SubscribeSyncCursor"
	getVersionMethod          = "GetVersion"
	getPayerInfoMethod        = "GetPayerInfo"
)

type Service struct {
	metadata_api.UnimplementedMetadataApiServer
	ctx              context.Context
	logger           *zap.Logger
	cu               CursorUpdater
	version          *semver.Version
	payerInfoFetcher IPayerInfoFetcher
}

func NewMetadataAPIService(
	ctx context.Context,
	logger *zap.Logger,
	updater CursorUpdater,
	version *semver.Version,
	payerInfoFetcher IPayerInfoFetcher,
) (*Service, error) {
	return &Service{
		ctx:              ctx,
		logger:           logger,
		cu:               updater,
		version:          version,
		payerInfoFetcher: payerInfoFetcher,
	}, nil
}

func (s *Service) GetSyncCursor(
	_ context.Context,
	_ *metadata_api.GetSyncCursorRequest,
) (*metadata_api.GetSyncCursorResponse, error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(getSyncCursorMethod))
	}

	return &metadata_api.GetSyncCursorResponse{
		LatestSync: s.cu.GetCursor(),
	}, nil
}

func (s *Service) SubscribeSyncCursor(
	_ *metadata_api.GetSyncCursorRequest,
	stream metadata_api.MetadataApi_SubscribeSyncCursorServer,
) error {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(subscribeSyncCursorMethod))
	}

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
				s.logger.Debug("channel closed by worker")
				return nil
			}
		case <-stream.Context().Done():
			s.logger.Debug("stream closed")
			return nil
		case <-s.ctx.Done():
			s.logger.Debug("service closed")
			return nil
		}
	}
}

func (s *Service) GetVersion(
	_ context.Context,
	_ *metadata_api.GetVersionRequest,
) (*metadata_api.GetVersionResponse, error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(getVersionMethod))
	}

	if s.version == nil {
		s.logger.Error("version is not set")
		return nil, status.Errorf(codes.Internal, "version is not set")
	}

	return &metadata_api.GetVersionResponse{
		Version: s.version.String(),
	}, nil
}

func (s *Service) GetPayerInfo(
	ctx context.Context,
	req *metadata_api.GetPayerInfoRequest,
) (*metadata_api.GetPayerInfoResponse, error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug(
			"received request",
			utils.MethodField(getPayerInfoMethod),
			utils.BodyField(req),
		)
	}

	if len(req.PayerAddresses) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "payer_addresses cannot be empty")
	}

	// Map the granularity enum to the internal type
	var groupBy PayerInfoGroupBy
	switch req.Granularity {
	case metadata_api.PayerInfoGranularity_PAYER_INFO_GRANULARITY_HOUR:
		groupBy = PayerInfoGroupByHour
	case metadata_api.PayerInfoGranularity_PAYER_INFO_GRANULARITY_DAY:
		groupBy = PayerInfoGroupByDay
	default:
		// Default to hour granularity if unspecified
		groupBy = PayerInfoGroupByHour
	}

	// Initialize response
	response := &metadata_api.GetPayerInfoResponse{
		PayerInfo: make(map[string]*metadata_api.GetPayerInfoResponse_PayerInfo),
	}

	// Look up each payer address and fetch their info
	for _, address := range req.PayerAddresses {
		// Look up payer ID from the database
		payerID, err := s.payerInfoFetcher.GetPayerByAddress(ctx, address)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, status.Errorf(codes.NotFound, "payer address not found: %s", address)
			}
			s.logger.Error("failed to find payer",
				utils.MethodField(getPayerInfoMethod),
				utils.PayerAddressField(address),
				zap.Error(err),
			)
			return nil, status.Errorf(codes.Internal, "failed to look up payer")
		}

		// Fetch payer info from the fetcher
		payerInfo, err := s.payerInfoFetcher.GetPayerInfo(
			ctx,
			payerID,
			groupBy,
		)
		if err != nil {
			s.logger.Error("failed to get payer info",
				utils.MethodField(getPayerInfoMethod),
				utils.PayerAddressField(address),
				utils.PayerIDField(payerID),
				zap.Error(err),
			)
			return nil, status.Errorf(
				codes.Internal,
				"failed to get payer info for address: %s",
				address,
			)
		}

		response.PayerInfo[address] = payerInfo
	}

	return response, nil
}

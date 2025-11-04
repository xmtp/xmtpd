// Package metadata implements the metadata API service.
package metadata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/Masterminds/semver/v3"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	metadata_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	getSyncCursorMethod       = "GetSyncCursor"
	subscribeSyncCursorMethod = "SubscribeSyncCursor"
	getVersionMethod          = "GetVersion"
	getPayerInfoMethod        = "GetPayerInfo"

	requestMissingMessageError = "missing request message"
)

type Service struct {
	metadata_apiconnect.UnimplementedMetadataApiHandler

	ctx              context.Context
	logger           *zap.Logger
	cu               CursorUpdater
	version          *semver.Version
	payerInfoFetcher IPayerInfoFetcher
}

var _ metadata_apiconnect.MetadataApiHandler = (*Service)(nil)

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
	_ *connect.Request[metadata_api.GetSyncCursorRequest],
) (*connect.Response[metadata_api.GetSyncCursorResponse], error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(getSyncCursorMethod))
	}

	response := connect.NewResponse(&metadata_api.GetSyncCursorResponse{
		LatestSync: s.cu.GetCursor(),
	})

	return response, nil
}

func (s *Service) SubscribeSyncCursor(
	ctx context.Context,
	_ *connect.Request[metadata_api.GetSyncCursorRequest],
	stream *connect.ServerStream[metadata_api.GetSyncCursorResponse],
) error {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(subscribeSyncCursorMethod))
	}

	// Send the initial cursor. The subscriber will only send a new message if there was a change.
	cursor := s.cu.GetCursor()

	err := stream.Send(&metadata_api.GetSyncCursorResponse{
		LatestSync: cursor,
	})
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error sending cursor: %w", err),
		)
	}

	var (
		clientID   = fmt.Sprintf("client-%d", time.Now().UnixNano())
		updateChan = make(chan struct{}, 1)
	)

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
					return connect.NewError(
						connect.CodeInternal,
						fmt.Errorf("error sending cursor: %w", err),
					)
				}
			} else {
				s.logger.Debug("channel closed by worker")
				return nil
			}

		case <-ctx.Done():
			s.logger.Debug("metadata subscription stream closed")
			return nil

		case <-s.ctx.Done():
			s.logger.Debug("metadata service closed")
			return nil
		}
	}
}

func (s *Service) GetVersion(
	_ context.Context,
	_ *connect.Request[metadata_api.GetVersionRequest],
) (*connect.Response[metadata_api.GetVersionResponse], error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(getVersionMethod))
	}

	if s.version == nil {
		s.logger.Error("version is not set")
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("version is not set"),
		)
	}

	response := connect.NewResponse(&metadata_api.GetVersionResponse{
		Version: s.version.String(),
	})

	return response, nil
}

func (s *Service) GetPayerInfo(
	ctx context.Context,
	req *connect.Request[metadata_api.GetPayerInfoRequest],
) (*connect.Response[metadata_api.GetPayerInfoResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New(requestMissingMessageError),
		)
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug(
			"received request",
			utils.MethodField(getPayerInfoMethod),
			utils.BodyField(req),
		)
	}

	if len(req.Msg.PayerAddresses) == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("payer_addresses cannot be empty"),
		)
	}

	// Map the granularity enum to the internal type
	var groupBy PayerInfoGroupBy
	switch req.Msg.Granularity {
	case metadata_api.PayerInfoGranularity_PAYER_INFO_GRANULARITY_HOUR:
		groupBy = PayerInfoGroupByHour
	case metadata_api.PayerInfoGranularity_PAYER_INFO_GRANULARITY_DAY:
		groupBy = PayerInfoGroupByDay
	default:
		// Default to hour granularity if unspecified
		groupBy = PayerInfoGroupByHour
	}

	// Initialize response
	response := connect.NewResponse(&metadata_api.GetPayerInfoResponse{
		PayerInfo: make(map[string]*metadata_api.GetPayerInfoResponse_PayerInfo),
	})

	// Look up each payer address and fetch their info
	for _, address := range req.Msg.PayerAddresses {
		// Look up payer ID from the database
		payerID, err := s.payerInfoFetcher.GetPayerByAddress(ctx, address)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, connect.NewError(
					connect.CodeNotFound,
					fmt.Errorf("payer address not found: %s", address),
				)
			}
			s.logger.Error("failed to find payer",
				utils.MethodField(getPayerInfoMethod),
				utils.PayerAddressField(address),
				zap.Error(err),
			)
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("failed to look up payer"),
			)
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
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("failed to get payer info for address: %s", address),
			)
		}

		response.Msg.PayerInfo[address] = payerInfo
	}

	return response, nil
}

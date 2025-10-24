// Package mlsvalidate implements the MLS validation service interface.
package mlsvalidate

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	svc "github.com/xmtp/xmtpd/pkg/proto/mls_validation/v1"
	"github.com/xmtp/xmtpd/pkg/proto/mls_validation/v1/mls_validationv1connect"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type MLSValidationServiceImpl struct {
	grpcClient mls_validationv1connect.ValidationApiClient
}

var _ MLSValidationService = (*MLSValidationServiceImpl)(nil)

func NewMLSValidationService(
	ctx context.Context,
	logger *zap.Logger,
	cfg config.MlsValidationOptions,
	clientMetrics *grpcprom.ClientMetrics,
) (*MLSValidationServiceImpl, error) {
	target, isTLS, err := utils.HTTPAddressToGRPCTarget(cfg.GrpcAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTTP address to gRPC target: %w", err)
	}

	logger.Info(
		"connecting to mls validation service",
		zap.String("url", cfg.GrpcAddress),
		zap.String("target", target),
	)

	httpClient, err := utils.BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP client: %w", err)
	}

	validationClient := mls_validationv1connect.NewValidationApiClient(
		httpClient,
		target,
		utils.BuildGRPCDialOptions()...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validation client: %w", err)
	}

	return &MLSValidationServiceImpl{
		grpcClient: validationClient,
	}, nil
}

func (s *MLSValidationServiceImpl) GetAssociationState(
	ctx context.Context,
	oldUpdates []*associations.IdentityUpdate,
	newUpdates []*associations.IdentityUpdate,
) (*AssociationStateResult, error) {
	request := &svc.GetAssociationStateRequest{
		OldUpdates: oldUpdates,
		NewUpdates: newUpdates,
	}

	response, err := s.grpcClient.GetAssociationState(ctx, connect.NewRequest(request))
	if err != nil {
		return nil, err
	}

	return &AssociationStateResult{
		AssociationState: response.Msg.GetAssociationState(),
		StateDiff:        response.Msg.GetStateDiff(),
	}, nil
}

func (s *MLSValidationServiceImpl) GetAssociationStateFromEnvelopes(
	ctx context.Context,
	oldUpdateEnvelopes []queries.GatewayEnvelopesView,
	newUpdate *associations.IdentityUpdate,
) (*AssociationStateResult, error) {
	oldUpdates := make([]*associations.IdentityUpdate, len(oldUpdateEnvelopes))
	for i, update := range oldUpdateEnvelopes {
		originatorEnvelope, err := envelopes.NewOriginatorEnvelopeFromBytes(
			update.OriginatorEnvelope,
		)
		if err != nil {
			return nil, err
		}

		payloadInterface := originatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()

		payload, ok := payloadInterface.(*envelopesProto.ClientEnvelope_IdentityUpdate)
		if !ok || payload.IdentityUpdate == nil {
			return nil, fmt.Errorf("identity update is nil")
		}

		oldUpdates[i] = payload.IdentityUpdate
	}

	return s.GetAssociationState(ctx, oldUpdates, []*associations.IdentityUpdate{newUpdate})
}

func (s *MLSValidationServiceImpl) ValidateKeyPackages(
	ctx context.Context,
	keyPackages [][]byte,
) ([]KeyPackageValidationResult, error) {
	request := makeValidateKeyPackageRequest(keyPackages)

	response, err := s.grpcClient.ValidateInboxIdKeyPackages(ctx, connect.NewRequest(request))
	if err != nil {
		return nil, err
	}

	out := make([]KeyPackageValidationResult, len(response.Msg.Responses))
	for i, response := range response.Msg.Responses {
		if !response.IsOk {
			out[i] = KeyPackageValidationResult{
				IsOk:            false,
				InstallationKey: nil,
				Credential:      nil,
				Expiration:      0,
				ErrorMessage:    response.ErrorMessage,
			}
		} else {
			out[i] = KeyPackageValidationResult{
				IsOk:            true,
				InstallationKey: response.InstallationPublicKey,
				Credential:      nil,
				Expiration:      response.Expiration,
				ErrorMessage:    response.ErrorMessage,
			}
		}
	}
	return out, nil
}

func makeValidateKeyPackageRequest(
	keyPackageBytes [][]byte,
) *svc.ValidateKeyPackagesRequest {
	keyPackageRequests := make(
		[]*svc.ValidateKeyPackagesRequest_KeyPackage,
		len(keyPackageBytes),
	)

	for i, keyPackage := range keyPackageBytes {
		keyPackageRequests[i] = &svc.ValidateKeyPackagesRequest_KeyPackage{
			KeyPackageBytesTlsSerialized: keyPackage,
			IsInboxIdCredential:          true,
		}
	}

	return &svc.ValidateKeyPackagesRequest{
		KeyPackages: keyPackageRequests,
	}
}

func (s *MLSValidationServiceImpl) ValidateGroupMessages(
	ctx context.Context,
	groupMessages []*svc.ValidateGroupMessagesRequest_GroupMessage,
) ([]GroupMessageValidationResult, error) {
	request := makeValidateGroupMessagesRequest(groupMessages)

	response, err := s.grpcClient.ValidateGroupMessages(
		ctx,
		connect.NewRequest(request),
	)
	if err != nil {
		return nil, err
	}

	out := make([]GroupMessageValidationResult, len(response.Msg.Responses))
	for i, response := range response.Msg.Responses {
		if !response.IsOk {
			return nil, fmt.Errorf("validation failed with error %s", response.ErrorMessage)
		}
		out[i] = GroupMessageValidationResult{
			GroupID: response.GroupId,
		}
	}

	return out, nil
}

func makeValidateGroupMessagesRequest(
	groupMessages []*svc.ValidateGroupMessagesRequest_GroupMessage,
) *svc.ValidateGroupMessagesRequest {
	groupMessageRequests := make(
		[]*svc.ValidateGroupMessagesRequest_GroupMessage,
		len(groupMessages),
	)
	for i, groupMessage := range groupMessages {
		groupMessageRequests[i] = &svc.ValidateGroupMessagesRequest_GroupMessage{
			GroupMessageBytesTlsSerialized: groupMessage.GetGroupMessageBytesTlsSerialized(),
		}
	}
	return &svc.ValidateGroupMessagesRequest{
		GroupMessages: groupMessageRequests,
	}
}

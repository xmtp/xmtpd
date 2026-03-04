// Package mlsvalidate implements the MLS validation service interface.
package mlsvalidate

import (
	"context"
	"errors"
	"fmt"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	svc "github.com/xmtp/xmtpd/pkg/proto/mls_validation/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type MLSValidationServiceImpl struct {
	grpcClient svc.ValidationApiClient
}

var _ MLSValidationService = (*MLSValidationServiceImpl)(nil)

func NewMLSValidationService(
	ctx context.Context,
	logger *zap.Logger,
	cfg config.MlsValidationOptions,
	clientMetrics *grpcprom.ClientMetrics,
) (*MLSValidationServiceImpl, error) {
	logger.Info(
		"connecting to mls validation service",
		zap.String("url", cfg.GrpcAddress),
	)

	conn, err := utils.NewGRPCConn(cfg.GrpcAddress,
		grpc.WithUnaryInterceptor(clientMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(clientMetrics.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return &MLSValidationServiceImpl{
		grpcClient: svc.NewValidationApiClient(conn),
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

	response, err := s.grpcClient.GetAssociationState(ctx, request)
	if err != nil {
		return nil, err
	}

	return &AssociationStateResult{
		AssociationState: response.GetAssociationState(),
		StateDiff:        response.GetStateDiff(),
	}, nil
}

func (s *MLSValidationServiceImpl) GetAssociationStateFromEnvelopes(
	ctx context.Context,
	oldUpdateEnvelopes []queries.SelectGatewayEnvelopesByTopicsRow,
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
			return nil, errors.New("identity update is nil")
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

	response, err := s.grpcClient.ValidateInboxIdKeyPackages(ctx, request)
	if err != nil {
		return nil, err
	}

	out := make([]KeyPackageValidationResult, len(response.GetResponses()))
	for i, response := range response.GetResponses() {
		if !response.GetIsOk() {
			out[i] = KeyPackageValidationResult{
				IsOk:            false,
				InstallationKey: nil,
				Credential:      nil,
				Expiration:      0,
				ErrorMessage:    response.GetErrorMessage(),
			}
		} else {
			out[i] = KeyPackageValidationResult{
				IsOk:            true,
				InstallationKey: response.GetInstallationPublicKey(),
				Credential:      nil,
				Expiration:      response.GetExpiration(),
				ErrorMessage:    response.GetErrorMessage(),
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
	groupMessages []*mlsv1.GroupMessageInput,
) ([]GroupMessageValidationResult, error) {
	request := makeValidateGroupMessagesRequest(groupMessages)

	response, err := s.grpcClient.ValidateGroupMessages(
		ctx,
		request,
	)
	if err != nil {
		return nil, err
	}

	out := make([]GroupMessageValidationResult, len(response.GetResponses()))
	for i, response := range response.GetResponses() {
		if !response.GetIsOk() {
			return nil, fmt.Errorf("validation failed with error %s", response.GetErrorMessage())
		}
		out[i] = GroupMessageValidationResult{
			GroupID: response.GetGroupId(),
		}
	}

	return out, nil
}

func makeValidateGroupMessagesRequest(
	groupMessages []*mlsv1.GroupMessageInput,
) *svc.ValidateGroupMessagesRequest {
	groupMessageRequests := make(
		[]*svc.ValidateGroupMessagesRequest_GroupMessage,
		len(groupMessages),
	)

	for i, groupMessage := range groupMessages {
		groupMessageRequests[i] = &svc.ValidateGroupMessagesRequest_GroupMessage{
			GroupMessageBytesTlsSerialized: groupMessage.GetV1().GetData(),
		}
	}

	return &svc.ValidateGroupMessagesRequest{
		GroupMessages: groupMessageRequests,
	}
}

package mlsvalidate

import (
	"context"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	svc "github.com/xmtp/xmtpd/pkg/proto/mls_validation/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MLSValidationServiceImpl struct {
	grpcClient svc.ValidationApiClient
}

func NewMlsValidationService(
	ctx context.Context,
	cfg config.MlsValidationOptions,
) (*MLSValidationServiceImpl, error) {
	conn, err := grpc.NewClient(
		cfg.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		conn.Close()
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
	req := &svc.GetAssociationStateRequest{
		OldUpdates: oldUpdates,
		NewUpdates: newUpdates,
	}
	response, err := s.grpcClient.GetAssociationState(ctx, req)
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
	oldUpdateEnvelopes []queries.GatewayEnvelope,
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
	req := makeValidateKeyPackageRequest(keyPackages)

	response, err := s.grpcClient.ValidateInboxIdKeyPackages(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make([]KeyPackageValidationResult, len(response.Responses))
	for i, response := range response.Responses {
		if !response.IsOk {
			return nil, fmt.Errorf("validation failed with error %s", response.ErrorMessage)
		}
		out[i] = KeyPackageValidationResult{
			InstallationKey: response.InstallationPublicKey,
			Credential:      nil,
			Expiration:      response.Expiration,
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
	req := makeValidateGroupMessagesRequest(groupMessages)

	response, err := s.grpcClient.ValidateGroupMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make([]GroupMessageValidationResult, len(response.Responses))
	for i, response := range response.Responses {
		if !response.IsOk {
			return nil, fmt.Errorf("validation failed with error %s", response.ErrorMessage)
		}
		out[i] = GroupMessageValidationResult{
			GroupId: response.GroupId,
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
			GroupMessageBytesTlsSerialized: groupMessage.GetV1().Data,
		}
	}
	return &svc.ValidateGroupMessagesRequest{
		GroupMessages: groupMessageRequests,
	}
}

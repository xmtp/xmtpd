package mlsvalidate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/mls_validationv1"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	proto "github.com/xmtp/xmtpd/pkg/proto/mls_validation/v1"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestValidateKeyPackages(t *testing.T) {
	apiClient := mocks.NewMockValidationApiClient(t)
	svc := &MLSValidationServiceImpl{
		grpcClient: apiClient,
	}

	mockResponse := proto.ValidateInboxIdKeyPackagesResponse_Response{
		IsOk:                  true,
		InstallationPublicKey: testutils.RandomBytes(32),
		Credential:            nil,
		Expiration:            1,
	}

	apiClient.EXPECT().
		ValidateInboxIdKeyPackages(mock.Anything, mock.Anything).
		Times(1).
		Return(&proto.ValidateInboxIdKeyPackagesResponse{
			Responses: []*proto.ValidateInboxIdKeyPackagesResponse_Response{&mockResponse},
		},
			nil,
		)

	res, err := svc.ValidateKeyPackages(context.Background(), [][]byte{testutils.RandomBytes(32)})
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, mockResponse.InstallationPublicKey, res[0].InstallationKey)
	require.Nil(t, res[0].Credential)
}

func TestGetAssociationState(t *testing.T) {
	apiClient := mocks.NewMockValidationApiClient(t)
	svc := &MLSValidationServiceImpl{
		grpcClient: apiClient,
	}

	inboxID := testutils.RandomInboxIDString()
	address := testutils.RandomAddress().String()

	mockResponse := proto.GetAssociationStateResponse{
		AssociationState: &associations.AssociationState{
			InboxId: inboxID,
		},
		StateDiff: &associations.AssociationStateDiff{
			NewMembers: []*associations.MemberIdentifier{{
				Kind: &associations.MemberIdentifier_EthereumAddress{EthereumAddress: address},
			}},
		},
	}

	apiClient.EXPECT().
		GetAssociationState(mock.Anything, mock.Anything).
		Times(1).
		Return(&mockResponse, nil)

	res, err := svc.GetAssociationState(
		context.Background(),
		[]*associations.IdentityUpdate{},
		[]*associations.IdentityUpdate{},
	)
	require.NoError(t, err)
	require.Equal(t, inboxID, res.AssociationState.InboxId)
	require.Equal(t, address, res.StateDiff.NewMembers[0].GetEthereumAddress())
}

package storer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	mlsvalidateMock "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func buildIdentityUpdateStorer(
	t *testing.T,
) (*IdentityUpdateStorer, *mlsvalidateMock.MockMLSValidationService, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	db, _, cleanup := testutils.NewDB(t, ctx)
	appChainOptions, _ := testutils.GetContractsOptions(t)
	contractAddress := appChainOptions.IdentityUpdatesContractAddress

	client, err := blockchain.NewClient(ctx, appChainOptions.RpcUrl)
	require.NoError(t, err)
	contract, err := identityupdates.NewIdentityUpdates(
		common.HexToAddress(contractAddress),
		client,
	)
	validationService := mlsvalidateMock.NewMockMLSValidationService(t)

	require.NoError(t, err)
	storer := NewIdentityUpdateStorer(db, testutils.NewLog(t), contract, validationService)

	return storer, validationService, func() {
		cancel()
		cleanup()
	}
}

func TestStoreIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	storer, validationService, cleanup := buildIdentityUpdateStorer(t)
	defer cleanup()
	newAddress := "0x12345"
	validationService.EXPECT().
		GetAssociationStateFromEnvelopes(mock.Anything, mock.Anything, mock.Anything).
		Return(&mlsvalidate.AssociationStateResult{
			StateDiff: &associations.AssociationStateDiff{
				NewMembers: []*associations.MemberIdentifier{{
					Kind: &associations.MemberIdentifier_Address{Address: newAddress},
				}},
			},
		}, nil)

	// Using the RandomInboxId function, since they are both 32 bytes and we treat inbox IDs as
	// strings outside the blockchain
	inboxId := testutils.RandomGroupID()
	identityUpdate := associations.IdentityUpdate{
		InboxId: utils.HexEncode(inboxId[:]),
	}

	sequenceID := uint64(1)

	logMessage := testutils.BuildIdentityUpdateLog(
		t,
		inboxId,
		envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxId, &identityUpdate),
		sequenceID,
	)

	err := storer.StoreLog(
		ctx,
		logMessage,
	)
	require.NoError(t, err)

	querier := queries.New(storer.db)

	envelopes, queryErr := querier.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{
			OriginatorNodeIds: []int32{IDENTITY_UPDATE_ORIGINATOR_ID},
			RowLimit:          10,
		},
	)
	require.NoError(t, queryErr)

	require.Equal(t, len(envelopes), 1)

	firstEnvelope := envelopes[0]
	deserializedEnvelope := envelopesProto.OriginatorEnvelope{}
	require.NoError(t, proto.Unmarshal(firstEnvelope.OriginatorEnvelope, &deserializedEnvelope))
	require.Greater(t, len(deserializedEnvelope.UnsignedOriginatorEnvelope), 0)

	getInboxIdResult, logsErr := querier.GetAddressLogs(ctx, []string{newAddress})
	require.NoError(t, logsErr)
	require.Equal(t, getInboxIdResult[0].InboxID, utils.HexEncode(inboxId[:]))
}

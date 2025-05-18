package storer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	mlsvalidateMock "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func buildIdentityUpdateStorer(
	t *testing.T,
) (*IdentityUpdateStorer, *mlsvalidateMock.MockMLSValidationService) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	db, _ := testutils.NewDB(t, ctx)
	rpcUrl := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(rpcUrl)
	config.AppChain.IdentityUpdateBroadcasterAddress = testutils.DeployIdentityUpdatesContract(
		t,
		rpcUrl,
	)

	client, err := blockchain.NewChainClient(
		ctx,
		false, /*useDatabaseClient*/
		config.AppChain.RpcURL,
		common.HexToAddress(config.AppChain.GroupMessageBroadcasterAddress),
		common.HexToAddress(config.AppChain.IdentityUpdateBroadcasterAddress),
	)
	require.NoError(t, err)
	validationService := mlsvalidateMock.NewMockMLSValidationService(t)

	require.NoError(t, err)
	storer := NewIdentityUpdateStorer(db, testutils.NewLog(t), client, validationService)

	return storer, validationService
}

func TestStoreIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	storer, validationService := buildIdentityUpdateStorer(t)
	newAddress := "0x12345"
	validationService.EXPECT().
		GetAssociationStateFromEnvelopes(mock.Anything, mock.Anything, mock.Anything).
		Return(&mlsvalidate.AssociationStateResult{
			StateDiff: &associations.AssociationStateDiff{
				NewMembers: []*associations.MemberIdentifier{{
					Kind: &associations.MemberIdentifier_EthereumAddress{
						EthereumAddress: newAddress,
					},
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

	require.NoError(t, storer.StoreLog(
		ctx,
		logMessage,
	))

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

func TestStoreSequential(t *testing.T) {
	ctx := context.Background()
	storer, validationService := buildIdentityUpdateStorer(t)
	newAddress := "0x12345"

	numCalls := 0
	validationService.EXPECT().
		GetAssociationStateFromEnvelopes(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, prevEnvs []queries.GatewayEnvelope, newUpdate *associations.IdentityUpdate) (*mlsvalidate.AssociationStateResult, error) {
			numCalls++
			if numCalls > 1 {
				require.Len(t, prevEnvs, 1)

				return &mlsvalidate.AssociationStateResult{
					StateDiff: &associations.AssociationStateDiff{},
				}, nil
			}
			return &mlsvalidate.AssociationStateResult{
				StateDiff: &associations.AssociationStateDiff{
					NewMembers: []*associations.MemberIdentifier{{
						Kind: &associations.MemberIdentifier_EthereumAddress{
							EthereumAddress: newAddress,
						},
					}},
				},
			}, nil
		})

	// Using the RandomGroupID function, since they are both 32 bytes and we treat inbox IDs as
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

	require.NoError(t, storer.StoreLog(
		ctx,
		logMessage,
	))

	logMessage = testutils.BuildIdentityUpdateLog(
		t,
		inboxId,
		envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxId, &identityUpdate),
		sequenceID+1, // Increment the sequence ID by 1
	)

	require.NoError(t, storer.StoreLog(
		ctx,
		logMessage,
	))
}

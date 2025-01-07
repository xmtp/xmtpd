package payer_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	registryMocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func buildPayerService(
	t *testing.T,
) (*payer.Service, *blockchainMocks.MockIBlockchainPublisher, *registryMocks.MockNodeRegistry, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	log := testutils.NewLog(t)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	mockRegistry := registryMocks.NewMockNodeRegistry(t)

	require.NoError(t, err)
	mockMessagePublisher := blockchainMocks.NewMockIBlockchainPublisher(t)

	payerService, err := payer.NewPayerApiService(
		ctx,
		log,
		mockRegistry,
		privKey,
		mockMessagePublisher,
	)
	require.NoError(t, err)

	return payerService, mockMessagePublisher, mockRegistry, func() {
		cancel()
	}
}

func TestPublishIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	svc, mockMessagePublisher, _, cleanup := buildPayerService(t)
	defer cleanup()

	inboxId := testutils.RandomInboxId()
	inboxIdBytes, err := utils.ParseInboxId(inboxId)
	require.NoError(t, err)

	txnHash := common.Hash{1, 2, 3}
	sequenceId := uint64(99)

	identityUpdate := &associations.IdentityUpdate{
		InboxId: inboxId,
	}

	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxIdBytes, identityUpdate)
	envelopeBytes, err := proto.Marshal(envelope)
	require.NoError(t, err)

	mockMessagePublisher.EXPECT().
		PublishIdentityUpdate(mock.Anything, mock.Anything, mock.Anything).
		Return(&abis.IdentityUpdatesIdentityUpdateCreated{
			Raw: types.Log{
				TxHash: txnHash,
			},
			SequenceId: sequenceId,
			Update:     envelopeBytes,
		}, nil)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{envelope},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.OriginatorEnvelopes, 1)

	responseEnvelope := publishResponse.OriginatorEnvelopes[0]
	parsedOriginatorEnvelope, err := envelopes.NewOriginatorEnvelope(responseEnvelope)
	require.NoError(t, err)

	proof := parsedOriginatorEnvelope.Proto().Proof.(*envelopesProto.OriginatorEnvelope_BlockchainProof)

	require.Equal(t, proof.BlockchainProof.TransactionHash, txnHash[:])
	require.Equal(
		t,
		parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.OriginatorSequenceID(),
		sequenceId,
	)
}

func TestPublishToNodes(t *testing.T) {
	originatorServer, _, _, originatorCleanup := apiTestUtils.NewTestAPIServer(t)
	defer originatorCleanup()

	ctx := context.Background()
	svc, _, mockRegistry, cleanup := buildPayerService(t)
	defer cleanup()

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		HttpAddress: formatAddress(originatorServer.Addr().String()),
	}, nil)

	groupId := testutils.RandomGroupID()
	testGroupMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupId,
		[]byte("test message"),
		100, // This is the expected originator ID of the test server
	)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{testGroupMessage},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.OriginatorEnvelopes, 1)

	responseEnvelope := publishResponse.OriginatorEnvelopes[0]
	parsedOriginatorEnvelope, err := envelopes.NewOriginatorEnvelope(responseEnvelope)
	require.NoError(t, err)

	targetTopic := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.TargetTopic()

	require.Equal(t, targetTopic.Identifier(), groupId[:])
}

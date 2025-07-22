package payer_test

import (
	"context"
	"testing"
	"time"

	mlsvalidate2 "github.com/xmtp/xmtpd/pkg/mlsvalidate"

	"github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"

	"github.com/xmtp/xmtpd/pkg/constants"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	metadataMocks "github.com/xmtp/xmtpd/pkg/mocks/metadata_api"
	registryMocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type FixedMetadataApiClientConstructor struct {
	mockClient *metadataMocks.MockMetadataApiClient
}

func (c *FixedMetadataApiClientConstructor) NewMetadataApiClient(
	nodeId uint32,
) (metadata_api.MetadataApiClient, error) {
	return c.mockClient, nil
}

type MockSubscribeSyncCursorClient struct {
	metadata_api.MetadataApi_SubscribeSyncCursorClient
	ctx     context.Context
	updates []*metadata_api.GetSyncCursorResponse
	err     error
	index   int
}

func (m *MockSubscribeSyncCursorClient) CloseSend() error {
	return nil // No-op for the mock
}

// Recv simulates receiving cursor updates over time.
func (m *MockSubscribeSyncCursorClient) Recv() (*metadata_api.GetSyncCursorResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.index < len(m.updates) {
		resp := m.updates[m.index]
		m.index++
		return resp, nil
	}
	<-m.ctx.Done()
	return nil, m.ctx.Err()
}

func buildPayerService(
	t *testing.T,
) (*payer.Service, *blockchainMocks.MockIBlockchainPublisher, *registryMocks.MockNodeRegistry, *metadataMocks.MockMetadataApiClient, *mlsvalidate.MockMLSValidationService) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	log := testutils.NewLog(t)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	mockRegistry := registryMocks.NewMockNodeRegistry(t)

	require.NoError(t, err)
	mockMessagePublisher := blockchainMocks.NewMockIBlockchainPublisher(t)

	mockMLSvalidation := mlsvalidate.NewMockMLSValidationService(t)

	metaMocks := metadataMocks.NewMockMetadataApiClient(t)
	payerService, err := payer.NewPayerApiService(
		ctx,
		log,
		mockRegistry,
		privKey,
		mockMessagePublisher,
		&FixedMetadataApiClientConstructor{
			mockClient: metaMocks,
		},
		nil,
		mockMLSvalidation,
	)
	require.NoError(t, err)

	return payerService, mockMessagePublisher, mockRegistry, metaMocks, mockMLSvalidation
}

func TestPublishIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	svc, mockMessagePublisher, registryMocks, metaMocks, _ := buildPayerService(t)

	inboxId := testutils.RandomInboxId()
	inboxIdBytes, err := utils.ParseInboxId(inboxId)
	require.NoError(t, err)

	txnHash := common.Hash{1, 2, 3}
	sequenceId := uint64(99)

	identityUpdate := &associations.IdentityUpdate{
		InboxId: inboxId,
	}

	mockStream := &MockSubscribeSyncCursorClient{
		updates: []*metadata_api.GetSyncCursorResponse{
			{
				LatestSync: &envelopesProto.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: sequenceId},
				},
			},
		},
	}

	metaMocks.On("SubscribeSyncCursor", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			// Capture the context from the caller
			capturedCtx := args.Get(0).(context.Context)
			mockStream.ctx = capturedCtx // Store the captured context in the mock
		}).
		Return(mockStream, nil).
		Once()

	registryMocks.On("GetNodes").Return([]registry.Node{
		testutils.GetHealthyNode(100),
	}, nil)

	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxIdBytes, identityUpdate)
	envelopeBytes, err := proto.Marshal(envelope)
	require.NoError(t, err)

	mockMessagePublisher.EXPECT().
		PublishIdentityUpdate(mock.Anything, mock.Anything, mock.Anything).
		Return(&iu.IdentityUpdateBroadcasterIdentityUpdateCreated{
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
	originatorServer, _, mocks := apiTestUtils.NewTestAPIServer(t)

	ctx := context.Background()
	svc, _, mockRegistry, _, mockValidation := buildPayerService(t)

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		HttpAddress: formatAddress(originatorServer.Addr().String()),
	}, nil)

	mockRegistry.On("GetNodes").Return([]registry.Node{
		testutils.GetHealthyNode(100),
	}, nil)

	mocks.MockValidationService.On(
		"ValidateGroupMessages",
		mock.Anything,
		mock.Anything,
	).Return([]mlsvalidate2.GroupMessageValidationResult{
		{GroupId: "", IsCommit: false},
	}, nil)

	mockValidation.On(
		"ValidateGroupMessages",
		mock.Anything,
		mock.Anything,
	).Return([]mlsvalidate2.GroupMessageValidationResult{
		{GroupId: "", IsCommit: false},
	}, nil)

	groupId := testutils.RandomGroupID()
	testGroupMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupId,
		[]byte("test message"),
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

	targetOriginator := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator
	require.EqualValues(t, 100, targetOriginator)

	// expiry assumptions
	require.EqualValues(
		t,
		constants.DEFAULT_STORAGE_DURATION_DAYS,
		parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)

	expiryTime := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()
	expectedExpiry := time.Now().
		Add(time.Duration(constants.DEFAULT_STORAGE_DURATION_DAYS) * 24 * time.Hour).
		Unix()
	require.InDelta(
		t,
		expectedExpiry,
		expiryTime,
		10,
		"expiry time should be roughly now + DEFAULT_STORAGE_DURATION_DAYS",
	)
}

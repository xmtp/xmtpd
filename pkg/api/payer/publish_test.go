package payer_test

import (
	"context"
	"net"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"go.uber.org/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/blockchain"
	metadataMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/metadata_api"
	registryMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registry"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type FixedMetadataAPIClientConstructor struct {
	mockClient *metadataMocks.MockMetadataApiClient
}

func (c *FixedMetadataAPIClientConstructor) NewMetadataAPIClient(
	nodeID uint32,
) (metadata_api.MetadataApiClient, error) {
	return c.mockClient, nil
}

type MockSubscribeSyncCursorClient struct {
	metadata_api.MetadataApiClient
	ctx     context.Context
	updates []*metadata_api.GetSyncCursorResponse
	err     error
	index   int
}

var _ grpc.ServerStreamingClient[metadata_api.GetSyncCursorResponse] = (*MockSubscribeSyncCursorClient)(
	nil,
)

func (m *MockSubscribeSyncCursorClient) Context() context.Context {
	return m.ctx
}

func (m *MockSubscribeSyncCursorClient) Header() (metadata.MD, error) {
	return nil, nil
}

func (m *MockSubscribeSyncCursorClient) CloseSend() error {
	return nil
}

func (m *MockSubscribeSyncCursorClient) Trailer() metadata.MD {
	return nil
}

func (m *MockSubscribeSyncCursorClient) RecvMsg(req any) error {
	return nil
}

func (m *MockSubscribeSyncCursorClient) SendMsg(req any) error {
	return nil
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
	opts ...payer.Option,
) (*payer.Service, *blockchainMocks.MockIBlockchainPublisher, *registryMocks.MockNodeRegistry, *metadataMocks.MockMetadataApiClient) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	log := testutils.NewLog(t)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	mockRegistry := registryMocks.NewMockNodeRegistry(t)

	require.NoError(t, err)
	mockMessagePublisher := blockchainMocks.NewMockIBlockchainPublisher(t)

	metaMocks := metadataMocks.NewMockMetadataApiClient(t)
	payerService, err := payer.NewPayerAPIService(
		ctx,
		log,
		mockRegistry,
		privKey,
		mockMessagePublisher,
		nil,
		0,
		nil,
		opts...,
	)
	require.NoError(t, err)

	return payerService, mockMessagePublisher, mockRegistry, metaMocks
}

func TestPublishIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	svc, mockMessagePublisher, _, _ := buildPayerService(t)

	inboxID := testutils.RandomInboxIDBytes()
	txnHash := common.Hash{1, 2, 3}
	sequenceID := uint64(99)

	identityUpdate := &associations.IdentityUpdate{
		InboxId: utils.HexEncode(inboxID[:]),
	}

	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxID, identityUpdate)
	envelopeBytes, err := proto.Marshal(envelope)
	require.NoError(t, err)

	mockMessagePublisher.EXPECT().
		PublishIdentityUpdate(mock.Anything, mock.Anything, mock.Anything).
		Return(&iu.IdentityUpdateBroadcasterIdentityUpdateCreated{
			Raw: types.Log{
				TxHash: txnHash,
			},
			SequenceId: sequenceID,
			Update:     envelopeBytes,
		}, nil)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{envelope},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.Msg.GetOriginatorEnvelopes(), 1)

	responseEnvelope := publishResponse.Msg.GetOriginatorEnvelopes()[0]
	parsedOriginatorEnvelope, err := envelopes.NewOriginatorEnvelope(responseEnvelope)
	require.NoError(t, err)

	proof := parsedOriginatorEnvelope.Proto().GetProof().(*envelopesProto.OriginatorEnvelope_BlockchainProof)

	require.Equal(t, proof.BlockchainProof.GetTransactionHash(), txnHash[:])
	require.Equal(
		t,
		parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.OriginatorSequenceID(),
		sequenceID,
	)
}

func TestPublishToNodes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(t)

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
	}, nil)

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil)

	groupID := testutils.RandomGroupID()
	testGroupMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(false),
	)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{testGroupMessage},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.Msg.GetOriginatorEnvelopes(), 1)

	responseEnvelope := publishResponse.Msg.GetOriginatorEnvelopes()[0]
	parsedOriginatorEnvelope, err := envelopes.NewOriginatorEnvelope(responseEnvelope)
	require.NoError(t, err)

	targetTopic := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.TargetTopic()
	require.Equal(t, targetTopic.Identifier(), groupID[:])

	targetOriginator := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator
	require.EqualValues(t, 100, targetOriginator)

	// expiry assumptions
	require.EqualValues(
		t,
		constants.DefaultStorageDurationDays,
		parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)

	expiryTime := parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()
	expectedExpiry := time.Now().
		Add(time.Duration(constants.DefaultStorageDurationDays) * 24 * time.Hour).
		Unix()
	require.InDelta(
		t,
		expectedExpiry,
		expiryTime,
		10,
		"expiry time should be roughly now + DEFAULT_STORAGE_DURATION_DAYS.\nExpected: %v\nActual: %v",
		time.Unix(expectedExpiry, 0).Local().Format(time.RFC3339),
		time.Unix(int64(expiryTime), 0).Local().Format(time.RFC3339),
	)
}

type slowServer struct {
	delay atomic.Duration
	message_api.UnimplementedReplicationApiServer
}

func (s *slowServer) PublishPayerEnvelopes(
	context.Context,
	*message_api.PublishPayerEnvelopesRequest,
) (*message_api.PublishPayerEnvelopesResponse, error) {
	time.Sleep(s.delay.Load())

	res := &message_api.PublishPayerEnvelopesResponse{}
	return res, nil
}

func TestPublishToNodesExpires(t *testing.T) {
	const (
		publishTimeout = 2 * time.Second
	)

	// Create a GRPC server for which we can control the response time.

	listen, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	defer func() {
		_ = listen.Close()
	}()

	s := grpc.NewServer()

	var srv slowServer
	message_api.RegisterReplicationApiServer(s, &srv)

	go func() {
		err = s.Serve(listen)
		// Will happen *eventually* on test close.
		// If it happens before we will be alerted by other things
		assert.ErrorIs(t, err, net.ErrClosed)
	}()

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(t, payer.WithPublishTimeout(publishTimeout))

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		HTTPAddress: formatAddress(listen.Addr().String()),
	}, nil)

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil)

	groupID := testutils.RandomGroupID()
	testGroupMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(false),
	)

	// Make the server take longer than the service is willing to wait.
	srv.delay.Store(publishTimeout + time.Second)
	_, err = svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{testGroupMessage},
		}),
	)
	// Publish rpc should fail if server fails to respond within an appropriate amount of time.
	require.Error(t, err)

	// Publish rpc should succeed if completed within the deadline.
	srv.delay.Store(0)
	_, err = svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{testGroupMessage},
		}),
	)
	require.NoError(t, err)
}

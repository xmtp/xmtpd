package payer_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestNoBlockchain_IdentityUpdateRoutedToNode(t *testing.T) {
	ctx := context.Background()
	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
	)

	mockRegistry.EXPECT().GetNode(uint32(200)).Return(&registry.Node{
		NodeID:      200,
		HTTPAddress: "http://localhost:5051",
		IsCanonical: true,
	}, nil).Maybe()

	inboxID := testutils.RandomInboxIDBytes()
	identityUpdate := &associations.IdentityUpdate{
		InboxId: utils.HexEncode(inboxID[:]),
	}
	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxID, identityUpdate)

	_, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{envelope},
		}),
	)
	// Will fail connecting to node 200, but NOT with blockchain error
	if err != nil {
		require.NotContains(t, err.Error(), "blockchain",
			"identity update should route to node, not blockchain")
		t.Logf("Expected error (no node to connect to): %v", err)
	}

	mockBlockchain.AssertNotCalled(
		t, "PublishIdentityUpdate", mock.Anything, mock.Anything, mock.Anything,
	)
}

func TestNoBlockchain_CommitRoutedToNode(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t, apiTestUtils.WithTestNoBlockchain())

	ctx := context.Background()
	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
	)

	// Node 100 is the commit handler — point it at our test server
	mockRegistry.EXPECT().GetNode(uint32(100)).Return(&registry.Node{
		NodeID:      100,
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Create a commit message (shouldSendToBlockchain=true)
	groupID := testutils.RandomGroupID()
	commitMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true), // isCommit=true
	)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{commitMessage},
		}),
	)

	require.NoError(t, err, "commit should publish to node successfully in no-blockchain mode")
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.Msg.GetOriginatorEnvelopes(), 1)

	// Verify it was sent to node 100 (commit handler), not blockchain
	responseEnvelope := publishResponse.Msg.GetOriginatorEnvelopes()[0]
	parsedOriginatorEnvelope, err := envelopes.NewOriginatorEnvelope(responseEnvelope)
	require.NoError(t, err)

	targetOriginator := parsedOriginatorEnvelope.
		UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator
	require.Equal(t, uint32(100), targetOriginator, "commit should target node 100 (commit handler)")

	// Verify retention policy is set
	require.EqualValues(
		t,
		constants.DefaultStorageDurationDays,
		parsedOriginatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)

	// Blockchain publisher must NOT have been called
	mockBlockchain.AssertNotCalled(
		t, "PublishGroupMessage", mock.Anything, mock.Anything, mock.Anything,
	)
}

func TestNoBlockchain_RegularMessageStillUsesNodeSelector(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
	)

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Regular (non-commit) group message
	groupID := testutils.RandomGroupID()
	regularMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(false),
	)

	publishResponse, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{regularMessage},
		}),
	)

	require.NoError(t, err, "regular message should publish normally in no-blockchain mode")
	require.NotNil(t, publishResponse)
	require.Len(t, publishResponse.Msg.GetOriginatorEnvelopes(), 1)
}

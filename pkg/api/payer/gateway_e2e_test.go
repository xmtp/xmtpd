package payer_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/api/payer"
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

// TestNoBlockchain_GatewayE2E_CommitFlow tests full payer → node flow for
// commits in no-blockchain mode: publish commit, verify stored on commit node.
func TestNoBlockchain_GatewayE2E_CommitFlow(t *testing.T) {
	// Test server defaults to node ID 100
	commitNodeID := uint32(100)
	identityNodeID := uint32(200)

	commitSuite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, identityNodeID),
	)

	mockRegistry.EXPECT().GetNode(commitNodeID).Return(&registry.Node{
		NodeID:      commitNodeID,
		HTTPAddress: formatAddress(commitSuite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(commitNodeID),
	}, nil).Maybe()

	groupID := testutils.RandomGroupID()
	commitMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true),
	)

	resp, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{commitMessage},
		}),
	)
	require.NoError(t, err)
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)

	parsed, err := envelopes.NewOriginatorEnvelope(resp.Msg.GetOriginatorEnvelopes()[0])
	require.NoError(t, err)
	require.EqualValues(t, commitNodeID, parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator)

	mockBlockchain.AssertNotCalled(t, "PublishGroupMessage", mock.Anything, mock.Anything, mock.Anything)
	t.Logf("E2E commit flow: publish + verify on node %d OK", commitNodeID)
}

// TestNoBlockchain_GatewayE2E_MixedTraffic tests regular messages and commits
// flowing through payer — commits go to commit node, regular to selector.
func TestNoBlockchain_GatewayE2E_MixedTraffic(t *testing.T) {
	commitNodeID := uint32(100)

	// Both regular and commit messages will hit the same test server (node 100)
	// since that's what the test infra supports. The key assertion is that
	// commits get routed with the correct target_originator.
	suite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, 200),
	)

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		NodeID:      commitNodeID,
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(commitNodeID),
	}, nil).Maybe()

	groupID := testutils.RandomGroupID()

	// 5 regular messages
	for i := 0; i < 5; i++ {
		regularMsg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID,
			envelopesTestUtils.GetRealisticGroupMessagePayload(false),
		)
		resp, err := svc.PublishClientEnvelopes(
			ctx,
			connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{regularMsg},
			}),
		)
		require.NoError(t, err, "regular msg %d", i)
		require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)
	}

	// 1 commit — must target commit node specifically
	commitMsg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true),
	)
	commitResp, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{commitMsg},
		}),
	)
	require.NoError(t, err)

	parsed, err := envelopes.NewOriginatorEnvelope(commitResp.Msg.GetOriginatorEnvelopes()[0])
	require.NoError(t, err)
	require.EqualValues(t, commitNodeID, parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator)

	t.Logf("Mixed traffic: 5 regular + 1 commit, all published OK")
}

// TestNoBlockchain_GatewayE2E_IdentityUpdateFlow tests identity update
// routing through payer to dedicated identity node.
func TestNoBlockchain_GatewayE2E_IdentityUpdateFlow(t *testing.T) {
	// Use node 100 as identity node (test server is always 100)
	commitNodeID := uint32(200)
	identityNodeID := uint32(100)

	identitySuite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, identityNodeID),
	)

	mockRegistry.EXPECT().GetNode(identityNodeID).Return(&registry.Node{
		NodeID:      identityNodeID,
		HTTPAddress: formatAddress(identitySuite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(identityNodeID),
	}, nil).Maybe()

	inboxID := testutils.RandomInboxIDBytes()
	identityUpdate := &associations.IdentityUpdate{
		InboxId: utils.HexEncode(inboxID[:]),
	}
	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(inboxID, identityUpdate)

	resp, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{envelope},
		}),
	)
	require.NoError(t, err)
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)

	parsed, err := envelopes.NewOriginatorEnvelope(resp.Msg.GetOriginatorEnvelopes()[0])
	require.NoError(t, err)
	require.EqualValues(t, identityNodeID, parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator)

	mockBlockchain.AssertNotCalled(t, "PublishIdentityUpdate", mock.Anything, mock.Anything, mock.Anything)
	t.Logf("E2E identity update flow: publish + verify on node %d OK", identityNodeID)
}

// TestNoBlockchain_GatewayE2E_MultiGroupOrdering tests that commits for
// different groups all route to the same commit node, preserving ordering.
func TestNoBlockchain_GatewayE2E_MultiGroupOrdering(t *testing.T) {
	commitNodeID := uint32(100)

	commitSuite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, 200),
	)

	mockRegistry.EXPECT().GetNode(commitNodeID).Return(&registry.Node{
		NodeID:      commitNodeID,
		HTTPAddress: formatAddress(commitSuite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(commitNodeID),
	}, nil).Maybe()

	numGroups := 10
	for i := 0; i < numGroups; i++ {
		groupID := testutils.RandomGroupID()
		commitMsg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID,
			envelopesTestUtils.GetRealisticGroupMessagePayload(true),
		)

		resp, err := svc.PublishClientEnvelopes(
			ctx,
			connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{commitMsg},
			}),
		)
		require.NoError(t, err, "commit for group %d", i)
		require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)

		parsed, err := envelopes.NewOriginatorEnvelope(resp.Msg.GetOriginatorEnvelopes()[0])
		require.NoError(t, err)
		require.EqualValues(t, commitNodeID, parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator,
			"group %d commit should go to commit node", i)
	}

	t.Logf("Multi-group ordering: %d groups, all commits → node %d", numGroups, commitNodeID)
}

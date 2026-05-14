package payer_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// TestNoBlockchain_GatewayE2E_CommitThroughPayer exercises the full path:
// client → payer service → commit routed to dedicated node → stored in DB.
// This simulates what the gateway does when --no-blockchain is enabled.
func TestNoBlockchain_GatewayE2E_CommitThroughPayer(t *testing.T) {
	// Test API server always uses node ID 100
	commitNodeID := uint32(100)
	identityNodeID := uint32(200)

	commitSuite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx := context.Background()
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

	// Send 5 commits in sequence, verify all succeed
	for i := 0; i < 5; i++ {
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

		require.NoError(t, err, "commit %d should succeed", i)
		require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)

		parsed, err := envelopes.NewOriginatorEnvelope(
			resp.Msg.GetOriginatorEnvelopes()[0],
		)
		require.NoError(t, err)
		require.EqualValues(t, commitNodeID,
			parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator)
	}

	mockBlockchain.AssertNotCalled(t, "PublishGroupMessage",
		mock.Anything, mock.Anything, mock.Anything)
}

// TestNoBlockchain_GatewayE2E_IdentityUpdateThroughPayer exercises identity
// update routing through the payer to a dedicated identity node.
// Note: test server always uses node ID 100, so we configure identity node
// to also be 100 for this test (both can coexist in no-blockchain mode).
func TestNoBlockchain_GatewayE2E_IdentityUpdateThroughPayer(t *testing.T) {
	// Use node 100 for both commit and identity in this test
	// (test server is always node 100)
	identityNodeID := uint32(100)

	identitySuite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx := context.Background()
	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(200, identityNodeID), // commitNode=200, identityNode=100
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
	envelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(
		inboxID, identityUpdate,
	)

	resp, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{envelope},
		}),
	)

	require.NoError(t, err, "identity update should succeed through dedicated node")
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)

	parsed, err := envelopes.NewOriginatorEnvelope(
		resp.Msg.GetOriginatorEnvelopes()[0],
	)
	require.NoError(t, err)
	require.EqualValues(t, identityNodeID,
		parsed.UnsignedOriginatorEnvelope.PayerEnvelope.TargetOriginator)

	mockBlockchain.AssertNotCalled(t, "PublishIdentityUpdate",
		mock.Anything, mock.Anything, mock.Anything)
}

// TestNoBlockchain_MultiGroupOrdering verifies that commits for different
// groups all route to the same commit node (deterministic) and get sequential
// originator sequence IDs.
func TestNoBlockchain_MultiGroupOrdering(t *testing.T) {
	commitNodeID := uint32(100)

	suite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, 200),
	)

	mockRegistry.EXPECT().GetNode(commitNodeID).Return(&registry.Node{
		NodeID:      commitNodeID,
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(commitNodeID),
	}, nil).Maybe()

	// Publish commits for 10 different groups
	var seqIDs []uint64
	for i := 0; i < 10; i++ {
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

		parsed, err := envelopes.NewOriginatorEnvelope(
			resp.Msg.GetOriginatorEnvelopes()[0],
		)
		require.NoError(t, err)
		seqIDs = append(seqIDs, parsed.OriginatorSequenceID())
	}

	// Verify monotonically increasing sequence IDs (total ordering at commit node)
	for i := 1; i < len(seqIDs); i++ {
		require.Greater(t, seqIDs[i], seqIDs[i-1],
			"sequence IDs should be monotonically increasing: seq[%d]=%d should > seq[%d]=%d",
			i, seqIDs[i], i-1, seqIDs[i-1])
	}

	t.Logf("10 commits across 10 groups, all sequenced by node %d: %v", commitNodeID, seqIDs)
}

// TestNoBlockchain_MixedWorkload_CommitsAndRegularMessages simulates realistic
// traffic: regular messages + commits + identity updates all flowing through
// the same payer in no-blockchain mode.
func TestNoBlockchain_MixedWorkload_CommitsAndRegularMessages(t *testing.T) {
	// All test servers are node 100 internally. We use node 100 as commit
	// node and identity node target (both going to separate test server
	// instances). Regular messages also go to node 100 via GetNodes.
	commitNodeID := uint32(100)
	identityNodeID := uint32(100) // same node for simplicity

	// Single test server (node 100) handles all message types
	suite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, identityNodeID),
	)

	mockRegistry.EXPECT().GetNode(mock.Anything).Return(&registry.Node{
		NodeID:      100,
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// 1. Regular message
	groupID := testutils.RandomGroupID()
	regularMsg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(false),
	)
	resp, err := svc.PublishClientEnvelopes(ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{regularMsg},
		}),
	)
	require.NoError(t, err, "regular message")
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)
	t.Log("Regular message OK")

	// 2. Commit
	commitMsg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true),
	)
	resp, err = svc.PublishClientEnvelopes(ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{commitMsg},
		}),
	)
	require.NoError(t, err, "commit")
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)
	t.Log("Commit OK")

	// 3. Identity update
	inboxID := testutils.RandomInboxIDBytes()
	idUpdate := &associations.IdentityUpdate{
		InboxId: utils.HexEncode(inboxID[:]),
	}
	idEnvelope := envelopesTestUtils.CreateIdentityUpdateClientEnvelope(
		inboxID, idUpdate,
	)
	resp, err = svc.PublishClientEnvelopes(ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{idEnvelope},
		}),
	)
	require.NoError(t, err, "identity update")
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)
	t.Log("Identity update OK")

	// 4. Welcome message (should route normally)
	welcomeTopicID := make([]byte, 16)
	_, _ = rand.Read(welcomeTopicID)
	welcomeEnv := &envelopesProto.ClientEnvelope{
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindWelcomeMessagesV1, welcomeTopicID).Bytes(),
		},
		Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
			WelcomeMessage: &apiv1.WelcomeMessageInput{
				Version: &apiv1.WelcomeMessageInput_V1_{
					V1: &apiv1.WelcomeMessageInput_V1{
						Data: []byte("welcome-test-payload"),
					},
				},
			},
		},
	}
	resp, err = svc.PublishClientEnvelopes(ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{welcomeEnv},
		}),
	)
	require.NoError(t, err, "welcome message")
	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 1)
	t.Log("Welcome OK")
}

// TestNoBlockchain_SequentialCommitLatency measures per-commit latency through
// the full payer → node path, giving an accurate picture of gateway-level
// commit performance.
func TestNoBlockchain_SequentialCommitLatency(t *testing.T) {
	commitNodeID := uint32(100)

	suite := apiTestUtils.NewTestAPIServer(t,
		apiTestUtils.WithTestNoBlockchain(),
	)

	ctx := context.Background()
	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(commitNodeID, 200),
	)

	mockRegistry.EXPECT().GetNode(commitNodeID).Return(&registry.Node{
		NodeID:      commitNodeID,
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(commitNodeID),
	}, nil).Maybe()

	iterations := 50
	var totalDuration time.Duration

	for i := 0; i < iterations; i++ {
		groupID := testutils.RandomGroupID()
		commitMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID,
			envelopesTestUtils.GetRealisticGroupMessagePayload(true),
		)

		start := time.Now()
		_, err := svc.PublishClientEnvelopes(
			ctx,
			connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{commitMessage},
			}),
		)
		elapsed := time.Since(start)
		require.NoError(t, err)
		totalDuration += elapsed
	}

	avgMs := float64(totalDuration.Milliseconds()) / float64(iterations)
	t.Logf("Sequential commit latency: %d iterations, avg=%.1fms, total=%v",
		iterations, avgMs, totalDuration)

	// Commit through no-blockchain should be under 100ms avg
	require.Less(t, avgMs, 100.0,
		"average commit latency should be under 100ms")
}

package payer_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/api/payer"
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

// TestNoBlockchain_CommitNodeDown verifies that when the dedicated commit node
// is unreachable, the payer returns an error rather than silently dropping.
func TestNoBlockchain_CommitNodeDown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
		payer.WithPublishTimeout(2*time.Second),
	)

	// Point commit node 100 at a bogus address — simulates node down
	mockRegistry.EXPECT().GetNode(uint32(100)).Return(&registry.Node{
		NodeID:      100,
		HTTPAddress: "http://localhost:19999", // nothing listening
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Create a commit envelope
	groupID := testutils.RandomGroupID()
	commitMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true), // isCommit=true
	)

	_, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{commitMessage},
		}),
	)

	// Should fail because node 100 is unreachable
	require.Error(t, err, "commit to unreachable node should fail, not silently drop")
	t.Logf("Expected error when commit node down: %v", err)

	// Blockchain publisher must NOT have been called
	mockBlockchain.AssertNotCalled(t, "PublishGroupMessage", mock.Anything, mock.Anything, mock.Anything)
}

// TestNoBlockchain_IdentityNodeDown verifies identity update routing fails
// gracefully when the identity node is unreachable.
func TestNoBlockchain_IdentityNodeDown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
		payer.WithPublishTimeout(2*time.Second),
	)

	// Point identity node 200 at a bogus address
	mockRegistry.EXPECT().GetNode(uint32(200)).Return(&registry.Node{
		NodeID:      200,
		HTTPAddress: "http://localhost:19999",
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(200),
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

	require.Error(t, err, "identity update to unreachable node should fail")
	t.Logf("Expected error when identity node down: %v", err)

	mockBlockchain.AssertNotCalled(t, "PublishIdentityUpdate", mock.Anything, mock.Anything, mock.Anything)
}

// TestNoBlockchain_MixedBatch_PartialFailure verifies that a batch with both
// regular messages and commits handles partial failure correctly when the
// commit node is down but regular nodes are up.
func TestNoBlockchain_MixedBatch_PartialFailure(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc, _, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
		payer.WithPublishTimeout(2*time.Second),
	)

	// Regular messages route to a healthy node via selector
	mockRegistry.EXPECT().GetNode(mock.MatchedBy(func(id uint32) bool {
		return id != 100 && id != 200
	})).Return(&registry.Node{
		HTTPAddress: formatAddress(suite.APIServer.Addr()),
	}, nil).Maybe()

	// Commit node 100 is down
	mockRegistry.EXPECT().GetNode(uint32(100)).Return(&registry.Node{
		NodeID:      100,
		HTTPAddress: "http://localhost:19999",
		IsCanonical: true,
	}, nil).Maybe()

	mockRegistry.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Create mixed batch: one regular message + one commit
	groupID := testutils.RandomGroupID()
	regularMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(false),
	)
	commitMessage := envelopesTestUtils.CreateGroupMessageClientEnvelope(
		groupID,
		envelopesTestUtils.GetRealisticGroupMessagePayload(true),
	)

	_, err := svc.PublishClientEnvelopes(
		ctx,
		connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopesProto.ClientEnvelope{regularMessage, commitMessage},
		}),
	)

	// Batch should fail because one destination is unreachable
	require.Error(t, err, "mixed batch with unreachable commit node should fail")
	t.Logf("Mixed batch partial failure: %v", err)
}

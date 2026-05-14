package payer_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestNoBlockchain_IdentityUpdateRoutedToNode(t *testing.T) {
	ctx := context.Background()
	svc, mockBlockchain, mockRegistry, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
	)

	// Mock GetNode so client manager can look up node 200
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

	// Will fail trying to connect to node 200 (not running in test), but
	// critically should NOT try blockchain path.
	if err != nil {
		require.NotContains(t, err.Error(), "blockchain",
			"identity update should route to node, not blockchain")
		t.Logf("Expected error (no node to connect to): %v", err)
	}

	// Blockchain publisher must NOT have been called
	mockBlockchain.AssertNotCalled(t, "PublishIdentityUpdate", mock.Anything, mock.Anything, mock.Anything)
}

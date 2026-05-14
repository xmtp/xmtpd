package payer_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"

	"github.com/xmtp/xmtpd/pkg/api/payer"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
)

func TestNoBlockchain_LatencyComparison(t *testing.T) {
	const iterations = 50

	// --- No-blockchain commit path ---
	suiteNoChain := apiTestUtils.NewTestAPIServer(t, apiTestUtils.WithTestNoBlockchain())
	ctx := context.Background()
	svcNoChain, _, mockRegNoChain, _ := buildPayerService(
		t,
		payer.WithNoBlockchain(100, 200),
	)

	mockRegNoChain.EXPECT().GetNode(uint32(100)).Return(&registry.Node{
		NodeID:      100,
		HTTPAddress: formatAddress(suiteNoChain.APIServer.Addr()),
		IsCanonical: true,
	}, nil).Maybe()
	mockRegNoChain.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Warmup
	for range 3 {
		groupID := testutils.RandomGroupID()
		msg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID, envelopesTestUtils.GetRealisticGroupMessagePayload(true),
		)
		_, _ = svcNoChain.PublishClientEnvelopes(ctx, connect.NewRequest(
			&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{msg},
			}))
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		groupID := testutils.RandomGroupID()
		msg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID, envelopesTestUtils.GetRealisticGroupMessagePayload(true),
		)
		_, err := svcNoChain.PublishClientEnvelopes(ctx, connect.NewRequest(
			&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{msg},
			}))
		if err != nil {
			t.Fatalf("no-blockchain commit publish failed: %v", err)
		}
	}
	noChainDuration := time.Since(start)
	noChainAvg := noChainDuration / time.Duration(iterations)

	// --- Regular message path (baseline) ---
	suiteRegular := apiTestUtils.NewTestAPIServer(t)
	svcRegular, _, mockRegRegular, _ := buildPayerService(t)

	mockRegRegular.EXPECT().GetNode(uint32(100)).Return(&registry.Node{
		HTTPAddress: formatAddress(suiteRegular.APIServer.Addr()),
	}, nil).Maybe()
	mockRegRegular.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
	}, nil).Maybe()

	// Warmup
	for range 3 {
		groupID := testutils.RandomGroupID()
		msg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID, envelopesTestUtils.GetRealisticGroupMessagePayload(false),
		)
		_, _ = svcRegular.PublishClientEnvelopes(ctx, connect.NewRequest(
			&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{msg},
			}))
	}

	start = time.Now()
	for i := 0; i < iterations; i++ {
		groupID := testutils.RandomGroupID()
		msg := envelopesTestUtils.CreateGroupMessageClientEnvelope(
			groupID, envelopesTestUtils.GetRealisticGroupMessagePayload(false),
		)
		_, err := svcRegular.PublishClientEnvelopes(ctx, connect.NewRequest(
			&payer_api.PublishClientEnvelopesRequest{
				Envelopes: []*envelopesProto.ClientEnvelope{msg},
			}))
		if err != nil {
			t.Fatalf("regular publish failed: %v", err)
		}
	}
	regularDuration := time.Since(start)
	regularAvg := regularDuration / time.Duration(iterations)

	t.Logf("=== LATENCY COMPARISON (%d iterations) ===", iterations)
	t.Logf("No-blockchain commit path: %v total, %v avg per publish", noChainDuration, noChainAvg)
	t.Logf("Regular message path:      %v total, %v avg per publish", regularDuration, regularAvg)
	t.Logf("Ratio (nochain/regular):   %.2fx", float64(noChainAvg)/float64(regularAvg))
	t.Logf("Note: blockchain path would add ~2-15s per commit (Arbitrum L2 finality)")
}

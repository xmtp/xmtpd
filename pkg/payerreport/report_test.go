package payerreport_test

import (
	"crypto/rand"
	"testing"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"

	"github.com/xmtp/xmtpd/pkg/payerreport"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// Temporary function until we have a real merkle root
func randomBytes32() [32]byte {
	var b [32]byte
	//nolint:errcheck
	rand.Read(b[:])
	return b
}

func TestBuildPayerReport(t *testing.T) {
	inputs := []struct {
		name        string
		params      payerreport.BuildPayerReportParams
		expectErr   bool
		errContains string
	}{
		{
			name: "full report",
			params: payerreport.BuildPayerReportParams{
				OriginatorNodeID:    1,
				StartSequenceID:     0,
				EndSequenceID:       10,
				EndMinuteSinceEpoch: 10,
				Payers: map[common.Address]currency.PicoDollar{
					testutils.RandomAddress(): currency.PicoDollar(10),
				},
				NodeIDs:         []uint32{1},
				DomainSeparator: testutils.RandomDomainSeparator(),
			},
			expectErr: false,
		},
		{
			name: "empty payers",
			params: payerreport.BuildPayerReportParams{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				DomainSeparator:  testutils.RandomDomainSeparator(),
			},
			expectErr: false,
		},
		{
			name: "empty domain separator",
			params: payerreport.BuildPayerReportParams{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
			},
			expectErr:   true,
			errContains: "domain separator",
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			_, err := payerreport.BuildPayerReport(input.params)
			if input.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), input.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReportPacksToSame(t *testing.T) {
	reportsManager := constructReportsManager(t)

	domainSeparator, err := reportsManager.GetDomainSeparator(t.Context())
	require.NoError(t, err)

	payerReport, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    100,
		StartSequenceID:     0,
		EndSequenceID:       1,
		EndMinuteSinceEpoch: 10,
		Payers:              map[common.Address]currency.PicoDollar{},
		NodeIDs:             []uint32{100, 200},
		DomainSeparator:     domainSeparator,
	})

	require.NoError(t, err)

	reportWithStatus := payerreport.PayerReportWithStatus{
		PayerReport: payerReport.PayerReport,
	}

	// Ensure the report ID matches the one we built
	reportID, err := reportsManager.GetReportID(t.Context(), &reportWithStatus)
	require.NoError(t, err)
	require.Equal(t, reportID, payerReport.ID)
}

func constructReportsManager(t *testing.T) *blockchain.ReportsManager {
	log := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(
		t.Context(),
		rpcURL,
	)
	require.NoError(t, err)

	reportsManager, err := blockchain.NewReportsManager(
		log, client, signer, contractsOptions.SettlementChain,
	)
	require.NoError(t, err)
	return reportsManager
}

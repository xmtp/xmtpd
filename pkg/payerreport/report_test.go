package payerreport

import (
	"crypto/rand"
	"testing"

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
		params      BuildPayerReportParams
		expectErr   bool
		errContains string
	}{
		{
			name: "full report",
			params: BuildPayerReportParams{
				OriginatorNodeID:    1,
				StartSequenceID:     1,
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
			params: BuildPayerReportParams{
				OriginatorNodeID: 1,
				StartSequenceID:  1,
				EndSequenceID:    10,
				DomainSeparator:  testutils.RandomDomainSeparator(),
			},
			expectErr: false,
		},
		{
			name: "empty domain separator",
			params: BuildPayerReportParams{
				OriginatorNodeID: 1,
				StartSequenceID:  1,
				EndSequenceID:    10,
			},
			expectErr:   true,
			errContains: "domain separator",
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			_, err := BuildPayerReport(input.params)
			if input.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), input.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

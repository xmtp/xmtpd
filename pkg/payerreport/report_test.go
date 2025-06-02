package payerreport

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
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

func TestGetDigest(t *testing.T) {
	expectedNodeIDsHash := common.Hex2Bytes(
		"ea13edf2a1dffdeb6f76acdbc46a352bd5b9071e7a3a5e6a63a498a9caa547fa",
	)
	expectedDigest := common.Hex2Bytes(
		"1ec269bb27455a17e615c98f34f05a635943526e8fddff7b6a81a73bb1468b9c",
	)
	require.Equal(t, len(expectedNodeIDsHash), 32)
	require.Equal(t, len(expectedDigest), 32)

	// Get the expected values
	payersMerkleRoot := common.HexToHash(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	)
	nodeIDs := []uint32{100, 200, 300, 400, 500}
	nodeIDsHash := utils.PackAndHashNodeIDs(nodeIDs)
	require.Equal(t, nodeIDsHash, common.BytesToHash(expectedNodeIDsHash))

	originatorNodeID := uint32(1)
	startSequenceID := uint64(2)
	endSequenceID := uint64(3)
	domainSeparator := common.HexToHash(
		"dbc3c9c77bfb8c8656e87b666d2b06300835634ecfb091e1925d30614ceb1e43",
	)

	builtID, err := BuildPayerReportID(
		originatorNodeID,
		startSequenceID,
		endSequenceID,
		payersMerkleRoot,
		nodeIDs,
		domainSeparator,
	)
	require.NoError(t, err)
	require.Equal(t, *builtID, ReportID(expectedDigest))
}

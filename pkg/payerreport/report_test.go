package payerreport_test

import (
	"crypto/rand"
	"testing"

	"github.com/xmtp/xmtpd/pkg/payerreport"

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

func TestGetDigest(t *testing.T) {
	expectedNodeIDsHash := common.Hex2Bytes(
		"ea13edf2a1dffdeb6f76acdbc46a352bd5b9071e7a3a5e6a63a498a9caa547fa",
	)
	expectedDigest := common.Hex2Bytes(
		"79f316f2836745161f3020e431db382ce57aab339df1429de068a62bf940295b",
	)
	require.Equal(t, len(expectedNodeIDsHash), 32)
	require.Equal(t, len(expectedDigest), 32)

	// Get the expected values
	payersMerkleRoot := common.HexToHash(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	)
	nodeIDs := []uint32{100, 200, 300, 400, 500}
	nodeIDsHash := utils.PackSortAndHashNodeIDs(nodeIDs)
	require.Equal(t, nodeIDsHash, common.BytesToHash(expectedNodeIDsHash))

	originatorNodeID := uint32(1)
	startSequenceID := uint64(2)
	endSequenceID := uint64(3)
	endMinuteSinceEpoch := uint32(4)
	domainSeparator := common.HexToHash(
		"dbc3c9c77bfb8c8656e87b666d2b06300835634ecfb091e1925d30614ceb1e43",
	)

	builtID, err := payerreport.BuildPayerReportID(
		originatorNodeID,
		startSequenceID,
		endSequenceID,
		endMinuteSinceEpoch,
		payersMerkleRoot,
		nodeIDs,
		domainSeparator,
	)
	require.NoError(t, err)
	require.Equal(t, *builtID, payerreport.ReportID(expectedDigest))
}

func TestNodeOrderPacksToSameHash(t *testing.T) {
	nodeIDs := []uint32{100, 200, 300, 400, 500}
	nodeIDsHash := utils.PackSortAndHashNodeIDs(nodeIDs)

	require.Equal(t, nodeIDsHash, utils.PackSortAndHashNodeIDs([]uint32{500, 100, 200, 300, 400}))
	require.Equal(t, nodeIDsHash, utils.PackSortAndHashNodeIDs([]uint32{300, 200, 500, 100, 400}))
}

package payerreport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayerReportID(t *testing.T) {
	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  1,
		EndSequenceID:    10,
		PayersMerkleRoot: []byte{1, 2, 3},
		PayersLeafCount:  1,
		NodesHash:        []byte{4, 5, 6},
		NodesCount:       1,
	}

	id, err := report.ID()
	require.NoError(t, err)
	require.NotNil(t, id)
	require.Len(t, id, 32)
}

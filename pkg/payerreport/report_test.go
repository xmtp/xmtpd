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
		ActiveNodeIds:    []uint32{1, 2, 3},
	}

	id, err := report.ID()
	require.NoError(t, err)
	require.NotNil(t, id)
	require.Len(t, id, 32)
}

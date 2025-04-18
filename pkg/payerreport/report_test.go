package payerreport

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

// Temporary function until we have a real merkle root
func randomBytes32() [32]byte {
	var b [32]byte
	//nolint:errcheck
	rand.Read(b[:])
	return b
}

func TestPayerReportID(t *testing.T) {
	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  1,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIds:    []uint32{1, 2, 3},
	}

	id, err := report.ID()
	require.NoError(t, err)
	require.NotNil(t, id)
	require.Len(t, id, 32)
}

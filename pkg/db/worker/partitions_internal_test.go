package worker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePartitionInfo(t *testing.T) {
	const (
		name = "gateway_envelopes_meta_o100_s0_1000000"
	)

	info, err := parsePartitionInfo(name)
	require.NoError(t, err)

	require.Equal(t, name, info.name)
	require.Equal(t, uint32(100), info.nodeID)
	require.Equal(t, uint64(0), info.start)
	require.Equal(t, uint64(1_000_000), info.end)
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidSequenceID(t *testing.T) {
	sequenceID := int64(-1)
	require.False(t, IsValidSequenceID(sequenceID))
	sequenceID = int64(0)
	require.False(t, IsValidSequenceID(sequenceID))
	sequenceID = int64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.False(t, IsValidSequenceID(sequenceID))
}

func TestValidSequenceID(t *testing.T) {
	sequenceID := int64(1)
	require.True(t, IsValidSequenceID(sequenceID))
	sequenceID = int64(0b0000000000000000111111111111111111111111111111111111111111111111)
	require.True(t, IsValidSequenceID(sequenceID))
}

func TestGetNodeID(t *testing.T) {
	sid := uint64(1)
	require.Equal(t, uint16(0), NodeID(sid))
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.Equal(t, uint16(1), NodeID(sid))
}

func TestGetSequenceID(t *testing.T) {
	sid := uint64(0b0000000000000001111111111111111111111111111111111111111111111111)
	require.Equal(
		t,
		int64(0b0000000000000000111111111111111111111111111111111111111111111111),
		SequenceID(sid),
	)
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.Equal(t, int64(0), SequenceID(sid))
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000001)
	require.Equal(t, int64(1), SequenceID(sid))
}

func TestGetSID(t *testing.T) {
	nodeID := uint16(1)
	sequenceID := int64(1)
	require.Equal(
		t,
		uint64(0b0000000000000001000000000000000000000000000000000000000000000001),
		SID(nodeID, sequenceID),
	)
	nodeID = uint16(1)
	sequenceID = int64(0)
	require.Equal(
		t,
		uint64(0b0000000000000001000000000000000000000000000000000000000000000000),
		SID(nodeID, sequenceID),
	)
	nodeID = uint16(0)
	sequenceID = int64(1)
	require.Equal(
		t,
		uint64(0b0000000000000000000000000000000000000000000000000000000000000001),
		SID(nodeID, sequenceID),
	)
}

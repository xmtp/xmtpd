package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidLocalID(t *testing.T) {
	localID := int64(-1)
	require.False(t, IsValidLocalID(localID))
	localID = int64(0)
	require.False(t, IsValidLocalID(localID))
	localID = int64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.False(t, IsValidLocalID(localID))
}

func TestValidLocalID(t *testing.T) {
	localID := int64(1)
	require.True(t, IsValidLocalID(localID))
	localID = int64(0b0000000000000000111111111111111111111111111111111111111111111111)
	require.True(t, IsValidLocalID(localID))
}

func TestGetNodeID(t *testing.T) {
	sid := uint64(1)
	require.Equal(t, uint16(0), NodeID(sid))
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.Equal(t, uint16(1), NodeID(sid))
}

func TestGetLocalID(t *testing.T) {
	sid := uint64(0b0000000000000001111111111111111111111111111111111111111111111111)
	require.Equal(t, int64(0b0000000000000000111111111111111111111111111111111111111111111111), LocalID(sid))
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000000)
	require.Equal(t, int64(0), LocalID(sid))
	sid = uint64(0b0000000000000001000000000000000000000000000000000000000000000001)
	require.Equal(t, int64(1), LocalID(sid))
}

func TestGetSID(t *testing.T) {
	nodeID := uint16(1)
	localID := int64(1)
	require.Equal(t, uint64(0b0000000000000001000000000000000000000000000000000000000000000001), SID(nodeID, localID))
	nodeID = uint16(1)
	localID = int64(0)
	require.Equal(t, uint64(0b0000000000000001000000000000000000000000000000000000000000000000), SID(nodeID, localID))
	nodeID = uint16(0)
	localID = int64(1)
	require.Equal(t, uint64(0b0000000000000000000000000000000000000000000000000000000000000001), SID(nodeID, localID))
}

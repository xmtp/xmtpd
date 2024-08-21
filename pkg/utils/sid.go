package utils

// SIDS are 64-bit numbers consisting of 16 bits for the node ID
// followed by 48 bits for the sequence ID. This file
// contains methods for reading and constructing sids.
//
// We also leverage type-checking throughout the repo to avoid confusion:
// - SIDs are uint64
// - node IDs are uint16
// - sequence IDs are int64

const (
	// Number of bits used for node ID
	nodeIDBits uint64 = 16
	// Number of bits used for local ID
	localIDBits uint64 = 64 - nodeIDBits
	// A mask made by setting every bit in the node ID
	nodeIDMask uint64 = 0xFFFF << localIDBits
	// A mask made by setting every bit in the local ID
	localIDMask uint64 = ^nodeIDMask
)

func IsValidSequenceID(localID int64) bool {
	return localID > 0 && localID>>localIDBits == 0
}

func NodeID(sid uint64) uint16 {
	return uint16(sid >> localIDBits)
}

func SequenceID(sid uint64) int64 {
	return int64(sid & localIDMask)
}

func SID(nodeID uint16, localID int64) uint64 {
	return uint64(nodeID)<<localIDBits | uint64(localID)
}

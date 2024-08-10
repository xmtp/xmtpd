package utils

// Operations to marshall to/from SIDs.
// It is recommended not to use these directly - instead use a wrapper type that
// performs the relevant error handling.

const (
	nodeIDBits  uint64 = 16
	localIDBits uint64 = 64 - nodeIDBits
	nodeIDMask  uint64 = 0xFFFF << localIDBits
	localIDMask uint64 = ^nodeIDMask
)

func IsValidNodeID(nodeID int32) bool {
	return nodeID >= 0 && nodeID>>nodeIDBits == 0
}

func IsValidLocalID(localID int64) bool {
	return localID > 0 && localID>>localIDBits == 0
}

func IsValidSID(sid uint64) bool {
	return IsValidNodeID(NodeID(sid)) && IsValidLocalID(LocalID(sid))
}

func NodeID(sid uint64) int32 {
	return int32(sid >> localIDBits)
}

func LocalID(sid uint64) int64 {
	return int64(sid & localIDMask)
}

func SID(nodeID int32, localID int64) uint64 {
	return uint64(nodeID)<<localIDBits | uint64(localID)
}

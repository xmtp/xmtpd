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

func IsValidLocalID(localID int64) bool {
	return localID > 0 && localID>>localIDBits == 0
}

func NodeID(sid uint64) uint16 {
	return uint16(sid >> localIDBits)
}

func LocalID(sid uint64) int64 {
	return int64(sid & localIDMask)
}

func SID(nodeID uint16, localID int64) uint64 {
	return uint64(nodeID)<<localIDBits | uint64(localID)
}

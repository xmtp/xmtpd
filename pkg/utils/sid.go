package utils

const (
	nodeIDMask  uint64 = 0xFFFF << 48
	localIDMask uint64 = ^nodeIDMask
)

// Converts a local serial ID from the database into a global SID with a node ID prefix
func SID(localID int64) uint64 {
	nodeMask := uint64(localID) & nodeIDMask
	if localID < 0 || nodeMask != 0 {
		// Either indicates ID exhaustion or developer error -
		// the service should not continue running either way
		panic("Invalid local ID")
	}
	// TODO(rich): Plumb through and set node ID
	return uint64(localID)
}

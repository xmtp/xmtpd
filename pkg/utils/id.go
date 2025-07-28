package utils

import "fmt"

func ParseInboxID(inboxID []byte) ([32]byte, error) {
	if len(inboxID) != 32 {
		return [32]byte{}, fmt.Errorf(
			"invalid bytes length: expected 32 bytes, got %d",
			len(inboxID),
		)
	}

	var result [32]byte
	copy(result[:], inboxID)

	return result, nil
}

func ParseGroupID(groupID []byte) ([16]byte, error) {
	if len(groupID) != 16 {
		return [16]byte{}, fmt.Errorf(
			"invalid bytes length: expected 16 bytes, got %d",
			len(groupID),
		)
	}

	var result [16]byte
	copy(result[:], groupID)

	return result, nil
}

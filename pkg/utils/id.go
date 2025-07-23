package utils

import "fmt"

func ParseInboxID(inboxIDString []byte) ([32]byte, error) {
	if len(inboxIDString) != 32 {
		return [32]byte{}, fmt.Errorf(
			"invalid bytes length: expected 32 bytes, got %d",
			len(inboxIDString),
		)
	}

	var result [32]byte
	copy(result[:], inboxIDString)

	return result, nil
}

func ParseGroupID(groupIDString []byte) ([16]byte, error) {
	if len(groupIDString) != 16 {
		return [16]byte{}, fmt.Errorf(
			"invalid bytes length: expected 16 bytes, got %d",
			len(groupIDString),
		)
	}

	var result [16]byte
	copy(result[:], groupIDString)

	return result, nil
}

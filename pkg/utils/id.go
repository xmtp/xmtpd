package utils

import (
	"encoding/hex"
	"fmt"
)

func ParseInboxId(inboxIdString string) ([32]byte, error) {
	// Decode the hex string
	decoded, err := hex.DecodeString(inboxIdString)
	if err != nil {
		return [32]byte{}, err
	}

	return BytesToId(decoded)
}

func BytesToId(bytes []byte) ([32]byte, error) {
	if len(bytes) != 32 {
		return [32]byte{}, fmt.Errorf("invalid bytes length: expected 32 bytes, got %d", len(bytes))
	}

	var result [32]byte
	copy(result[:], bytes)
	return result, nil
}

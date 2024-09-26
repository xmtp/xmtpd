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

	// Ensure the decoded bytes are exactly 32 bytes long
	if len(decoded) != 32 {
		return [32]byte{}, fmt.Errorf(
			"invalid inbox ID length: expected 32 bytes, got %d",
			len(decoded),
		)
	}

	// Convert to [32]byte
	var result [32]byte
	copy(result[:], decoded)

	return result, nil
}

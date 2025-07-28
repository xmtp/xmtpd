package deserializer

import (
	"bytes"
	"fmt"
	"io"
)

func readVariableOpaqueVec(r *bytes.Reader) ([]byte, error) {
	// Step 1: Read first byte to get length format
	firstByte, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read length discriminator: %w", err)
	}

	lengthPrefixLen := 1 << (firstByte >> 6) // 0b00→1, 0b01→2, 0b10→4, 0b11→8
	length := uint64(firstByte & 0x3F)       // lower 6 bits

	// Step 2: Read remaining length bytes
	for i := 1; i < lengthPrefixLen; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read length byte %d: %w", i, err)
		}
		length = (length << 8) | uint64(b)
	}

	// Step 3: Read payload
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("failed to read opaque vector of length %d: %w", length, err)
	}
	return buf, nil
}

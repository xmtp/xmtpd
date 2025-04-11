package merkle

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/bits"

	"golang.org/x/crypto/sha3"
)

const LEAF_PREFIX = "leaf|"

// LeftPad pads a string with a specified character up to a given size.
func LeftPad(num string, size int, char string) string {
	for len(num) < size {
		num = char + num
	}
	return num
}

// To32ByteBuffer converts a number to a 32-byte buffer (64 hex chars).
func To32ByteBuffer(number uint64) []byte {
	hexStr := fmt.Sprintf("%064x", number)
	buf, _ := hex.DecodeString(hexStr)
	return buf
}

// From32ByteBuffer reads a uint32 from the last 4 bytes of a buffer.
func From32ByteBuffer(buffer []byte) uint32 {
	return binary.BigEndian.Uint32(buffer[28:32])
}

// BitCount32 counts the number of set bits in a 32-bit integer.
func BitCount32(n uint32) int {
	return bits.OnesCount32(n)
}

// Hash computes the Keccak-256 hash of a buffer.
func Hash(buffer []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(buffer)
	return hash.Sum(nil)
}

// HashNode computes the hash of two nodes concatenated.
func HashNode(left, right []byte) []byte {
	return Hash(append(left, right...))
}

func HashLeaf(leaf []byte) []byte {
	return Hash(append([]byte(LEAF_PREFIX), leaf...))
}

// RoundUpToPowerOf2 rounds up a number to the next power of 2.
func RoundUpToPowerOf2(number uint32) uint32 {
	if BitCount32(number) == 1 {
		return number
	}

	number |= number >> 1
	number |= number >> 2
	number |= number >> 4
	number |= number >> 8
	number |= number >> 16

	return number + 1
}

// BytesEqual compares two byte slices for equality.
func BytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

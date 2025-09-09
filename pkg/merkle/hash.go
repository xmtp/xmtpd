package merkle

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"
)

const (
	leafPrefix = "leaf|"
	nodePrefix = "node|"
	rootPrefix = "root|"
)

var (
	leafPrefixBytes             = []byte(leafPrefix)
	nodePrefixBytes             = []byte(nodePrefix)
	rootPrefixBytes             = []byte(rootPrefix)
	ErrInvalidLeafCount         = errors.New("invalid leaf count")
	ErrInvalidBufferLength      = errors.New("invalid buffer length")
	ErrInvalidIntToBytes32Input = errors.New("invalid int to bytes32 input")
	ErrInvalidBytes32ToIntInput = errors.New("invalid bytes32 to int input")
)

// Hash computes the Keccak-256 hash of a buffer.
func Hash(buffer []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(buffer)
	return hash.Sum(nil)
}

func HashLeaf(leaf []byte) []byte {
	leafPrefixLen := len(leafPrefixBytes)
	buffer := make([]byte, leafPrefixLen+len(leaf))

	copy(buffer[:leafPrefixLen], leafPrefixBytes)
	copy(buffer[leafPrefixLen:], leaf)

	return Hash(buffer)
}

func HashNodePair(left, right []byte) []byte {
	nodePrefixLen := len(nodePrefixBytes)
	buffer := make([]byte, nodePrefixLen+len(left)+len(right))

	copy(buffer[:nodePrefixLen], nodePrefixBytes)
	copy(buffer[nodePrefixLen:], left)
	copy(buffer[nodePrefixLen+len(left):], right)

	return Hash(buffer)
}

func HashPairlessNode(node []byte) []byte {
	nodePrefixLen := len(nodePrefixBytes)
	buffer := make([]byte, nodePrefixLen+len(node))

	copy(buffer[:nodePrefixLen], nodePrefixBytes)
	copy(buffer[nodePrefixLen:], node)

	return Hash(buffer)
}

func HashRoot(leafCount int, root []byte) ([]byte, error) {
	if leafCount < 0 || leafCount > (1<<31)-1 {
		return nil, ErrInvalidLeafCount
	}

	leafCountBytes, err := IntToBytes32(leafCount)
	if err != nil {
		return nil, err
	}

	rootPrefixLen := len(rootPrefixBytes)
	leafCountLen := len(leafCountBytes) // Length of the byte representation
	rootLen := len(root)

	buffer := make([]byte, rootPrefixLen+leafCountLen+rootLen)

	copy(buffer[:rootPrefixLen], rootPrefixBytes)
	copy(buffer[rootPrefixLen:rootPrefixLen+leafCountLen], leafCountBytes) // Copy the bytes
	copy(buffer[rootPrefixLen+leafCountLen:], root)                        // Copy the root

	return Hash(buffer), nil
}

func IntToBytes32(value int) ([]byte, error) {
	if value < 0 || value > (1<<31)-1 {
		return nil, ErrInvalidIntToBytes32Input
	}

	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(value))

	buffer := make([]byte, 32)
	copy(buffer[24:], valueBytes)

	return buffer, nil
}

func Bytes32ToInt(buffer []byte) (int, error) {
	if len(buffer) != 32 {
		return 0, ErrInvalidBufferLength
	}

	// Check that all of the first 28 bytes are 0
	for i := 0; i < 28; i++ {
		if buffer[i] != 0 {
			return 0, ErrInvalidBytes32ToIntInput
		}
	}

	uint32Value := binary.BigEndian.Uint32(buffer[28:])

	if uint32Value > 1<<31-1 {
		return 0, ErrInvalidBytes32ToIntInput
	}

	return int(uint32Value), nil
}

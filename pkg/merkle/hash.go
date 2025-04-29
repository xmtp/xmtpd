package merkle

import (
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

const (
	LEAF_PREFIX = "leaf|"
	NODE_PREFIX = "node|"
	ROOT_PREFIX = "root|"
)

var (
	leafPrefixBytes = []byte(LEAF_PREFIX)
	nodePrefixBytes = []byte(NODE_PREFIX)
	rootPrefixBytes = []byte(ROOT_PREFIX)
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

func HashRoot(leafCount int, root []byte) []byte {
	leafCountBytes := IntTo32Bytes(leafCount)

	rootPrefixLen := len(rootPrefixBytes)
	leafCountLen := len(leafCountBytes) // Length of the byte representation
	rootLen := len(root)

	buffer := make([]byte, rootPrefixLen+leafCountLen+rootLen)

	copy(buffer[:rootPrefixLen], rootPrefixBytes)
	copy(buffer[rootPrefixLen:rootPrefixLen+leafCountLen], leafCountBytes) // Copy the bytes
	copy(buffer[rootPrefixLen+leafCountLen:], root)                        // Copy the root

	return Hash(buffer)
}

func IntTo32Bytes(value int) []byte {
	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(value))

	buffer := make([]byte, 32)
	copy(buffer[24:], valueBytes)

	return buffer
}

func BytesToBigInt(buffer []byte) int {
	return int(binary.BigEndian.Uint64(buffer[24:]))
}

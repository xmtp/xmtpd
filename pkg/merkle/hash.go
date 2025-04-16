package merkle

import (
	"golang.org/x/crypto/sha3"
)

const (
	LEAF_PREFIX = "leaf|"
	NODE_PREFIX = "node|"
)

// Pre-computed byte slices for prefixes to avoid repeated conversions
var (
	leafPrefixBytes = []byte(LEAF_PREFIX)
	nodePrefixBytes = []byte(NODE_PREFIX)
)

// Hash computes the Keccak-256 hash of a buffer.
func Hash(buffer []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(buffer)
	return hash.Sum(nil)
}

func HashNode(left, right []byte) []byte {
	nodePrefixLen := len(nodePrefixBytes)
	buffer := make([]byte, nodePrefixLen+len(left)+len(right))

	copy(buffer[:nodePrefixLen], nodePrefixBytes)
	copy(buffer[nodePrefixLen:], left)
	copy(buffer[nodePrefixLen+len(left):], right)

	return Hash(buffer)
}

// HashLeaf computes the hash of a leaf node.
func HashLeaf(leaf []byte) []byte {
	leafPrefixLen := len(leafPrefixBytes)
	buffer := make([]byte, leafPrefixLen+len(leaf))

	copy(buffer[:leafPrefixLen], leafPrefixBytes)
	copy(buffer[leafPrefixLen:], leaf)

	return Hash(buffer)
}

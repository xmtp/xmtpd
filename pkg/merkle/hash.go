package merkle

import (
	"golang.org/x/crypto/sha3"
)

const LEAF_PREFIX = "leaf|"

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

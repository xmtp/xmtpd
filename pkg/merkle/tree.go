package merkle

import (
	"errors"
	"math/bits"
)

// MerkleTree represents a Merkle tree data structure.
type MerkleTree struct {
	tree      [][]byte
	elements  [][]byte
	root      []byte
	leafCount int
}

var (
	ErrTreeEmpty   = errors.New("tree is empty")
	ErrTreeRootNil = errors.New("tree root is nil")
)

// NewMerkleTree creates a new Merkle tree from the given elements.
func NewMerkleTree(elements [][]byte) (*MerkleTree, error) {
	elementsDeepCopy := make([][]byte, len(elements))
	for i, element := range elements {
		elementsDeepCopy[i] = make([]byte, len(element))
		copy(elementsDeepCopy[i], element)
	}

	leaves, err := makeLeaves(elementsDeepCopy)
	if err != nil {
		return nil, err
	}

	tree, err := makeTree(leaves)
	if err != nil {
		return nil, err
	}

	return &MerkleTree{
		tree:      tree,
		elements:  elementsDeepCopy,
		root:      tree[1],
		leafCount: len(elementsDeepCopy),
	}, nil
}

// Tree returns the 1-indexed representation of the Merkle tree.
func (m *MerkleTree) Tree() [][]byte {
	return m.tree
}

// Elements returns the raw elements of the Merkle tree.
func (m *MerkleTree) Elements() [][]byte {
	return m.elements
}

// Root returns the root hash of the Merkle tree.
func (m *MerkleTree) Root() []byte {
	return m.root
}

// LeafCount returns the number of leaves in the Merkle tree.
func (m *MerkleTree) LeafCount() int {
	return m.leafCount
}

// BuildTree builds a serialized Merkle tree from an array of leaf nodes.
//
// The tree is 1-indexed, so root is at index 1.
// The internal nodes are at index 2 to N.
// The leaves are at index N+1 to 2N-1.
//
// For any node at index i:
// - left child is at index 2*i
// - right child is at index 2*i+1
// - parent is at floor(i/2)
func makeTree(leaves [][]byte) ([][]byte, error) {
	if len(leaves) == 0 {
		return nil, ErrTreeEmpty
	}

	leafCount := CalculateBalancedLeafCount(len(leaves))

	// Allocate 2N space for the tree. (N leaf nodes, N-1 internal nodes)
	tree := make([][]byte, leafCount<<1)

	lowerBound := leafCount
	upperBound := leafCount + len(leaves) - 1

	// Copy leaves into the tree, starting at index N.
	for i := 0; i < len(leaves); i++ {
		tree[lowerBound+i] = leaves[i]
	}

	// Build the tree.
	// Start from the last internal node and work our way up to the root.
	for i := leafCount - 1; i >= 0; i-- {
		leftChildIndex := GetLeftChild(i)

		if leftChildIndex > upperBound {
			continue
		}

		// Detect the level is processed.
		if leftChildIndex <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		rightChildIndex := GetRightChild(i)

		// Hash both children if they exist.
		// If only left children exists, use it.
		// If no children exist, leave it as nil.
		if tree[leftChildIndex] != nil && tree[rightChildIndex] != nil {
			tree[i] = HashNode(tree[leftChildIndex], tree[rightChildIndex])
		} else if tree[leftChildIndex] != nil {
			tree[i] = tree[leftChildIndex]
		}
	}

	if tree[1] == nil {
		return nil, ErrTreeRootNil
	}

	return tree, nil
}

// makeLeaves returns the leaves of the tree,
// ordered in the same order as the provided elements.
func makeLeaves(elements [][]byte) ([][]byte, error) {
	if len(elements) == 0 {
		return nil, errors.New("elements cannot be empty")
	}

	leaves := make([][]byte, len(elements))
	for i, element := range elements {
		leaves[i] = HashLeaf(element)
	}

	return leaves, nil
}

// CalculateBalancedLeafCount returns the number of leaves in a balanced tree.
// To calculate the number of leaves in a tree, we need to round up to the next power of 2.
func CalculateBalancedLeafCount(elementCount int) int {
	if elementCount == 0 {
		return 0
	}

	return int(roundUpToPowerOf2(uint32(elementCount)))
}

// roundUpToPowerOf2 rounds up a number to the next power of 2.
// Rounding up to the next power of 2 is necessary to ensure that the tree is balanced.
func roundUpToPowerOf2(n uint32) uint32 {
	if bits.OnesCount32(n) == 1 {
		return n
	}

	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	return n + 1
}

// GetLeftChild returns the index of the left child for a node at the given index
func GetLeftChild(index int) int {
	return index << 1 // index * 2
}

// GetRightChild returns the index of the right child for a node at the given index
func GetRightChild(index int) int {
	return (index << 1) + 1 // index * 2 + 1
}

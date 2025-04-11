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
	depth     int
	leafCount int
}

// NewMerkleTree creates a new Merkle tree from the given elements.
func NewMerkleTree(elements [][]byte) (*MerkleTree, error) {
	if len(elements) == 0 {
		return nil, errors.New("elements cannot be empty")
	}

	leaves := make([][]byte, len(elements))
	for i, element := range elements {
		leaves[i] = HashLeaf(element)
	}

	tree, depth := buildTree(leaves)

	return &MerkleTree{
		tree:      tree,
		elements:  elements,
		root:      tree[1],
		depth:     depth,
		leafCount: len(elements),
	}, nil
}

// Tree returns the tree of the Merkle tree.
func (m *MerkleTree) Tree() [][]byte {
	return m.tree
}

// Elements returns the elements of the Merkle tree.
func (m *MerkleTree) Elements() [][]byte {
	return m.elements
}

// Root returns the root hash of the Merkle tree.
func (m *MerkleTree) Root() []byte {
	return m.root
}

// Depth returns the depth of the Merkle tree.
func (m *MerkleTree) Depth() int {
	return m.depth
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
func buildTree(leaves [][]byte) ([][]byte, int) {
	depth := getDepth(uint32(len(leaves)))
	leafCount := getLeafCount(len(leaves))

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
		leftChildIndex := getLeftChild(i)

		if leftChildIndex > upperBound {
			continue
		}

		// Detect the level is processed.
		if leftChildIndex <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		rightChildIndex := getRightChild(i)

		// Hash both children if they exist.
		// If only left children exists, use it.
		// If no children exist, leave it as nil.
		if tree[leftChildIndex] != nil && tree[rightChildIndex] != nil {
			tree[i] = HashNode(tree[leftChildIndex], tree[rightChildIndex])
		} else if tree[leftChildIndex] != nil {
			tree[i] = tree[leftChildIndex]
		}
	}

	return tree, depth
}

// getDepth calculates the depth of a Merkle tree based on element count.
func getDepth(n uint32) int {
	if n <= 1 {
		return 0
	}
	// Round up to next power of 2 if not already a power of 2
	if bits.OnesCount32(n) > 1 {
		n = roundUpToPowerOf2(n)
	}
	// bits.Len32 gives position of highest bit + 1
	return bits.Len32(n) - 1
}

// getLeafCount returns the number of leaves in a tree.
func getLeafCount(elementCount int) int {
	return int(roundUpToPowerOf2(uint32(elementCount)))
}

// roundUpToPowerOf2 rounds up a number to the next power of 2.
// Rounding up to the next power of 2 is necessary to ensure that the tree is balanced.
func roundUpToPowerOf2(number uint32) uint32 {
	if bits.OnesCount32(number) == 1 {
		return number
	}

	number |= number >> 1
	number |= number >> 2
	number |= number >> 4
	number |= number >> 8
	number |= number >> 16

	return number + 1
}

// getLeftChild returns the index of the left child for a node at the given index
func getLeftChild(index int) int {
	return index << 1 // index * 2
}

// getRightChild returns the index of the right child for a node at the given index
func getRightChild(index int) int {
	return (index << 1) + 1 // index * 2 + 1
}

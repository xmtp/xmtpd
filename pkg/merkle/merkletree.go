package merkle

import (
	"errors"
	"math"
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
		root:      tree[0],
		depth:     depth,
		leafCount: len(elements),
	}, nil
}

// Root returns the root hash of the Merkle tree.
func (m *MerkleTree) Root() []byte {
	return m.root
}

// Elements returns the elements of the Merkle tree.
func (m *MerkleTree) Elements() [][]byte {
	return m.elements
}

// Depth returns the depth of the Merkle tree.
func (m *MerkleTree) Depth() int {
	return m.depth
}

// BuildTree builds a serialized Merkle tree from an array of leaf nodes.
func buildTree(leaves [][]byte) ([][]byte, int) {
	depth := getDepth(len(leaves))
	leafCount := getLeafCount(len(leaves))

	// Allocate 2N space for the tree. (N leaf nodes, N-1 internal nodes)
	tree := make([][]byte, leafCount<<1)

	// Copy leafs into the tree
	for i := 0; i < len(leaves); i++ {
		tree[leafCount+i] = leaves[i]
	}

	lowerBound := leafCount
	upperBound := leafCount + len(leaves) - 1

	// Build the tree
	for i := leafCount - 1; i >= 0; i-- {
		index := i << 1

		if index > upperBound {
			continue
		}

		if index <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		if index == upperBound {
			tree[i] = tree[index]
			continue
		}

		// Ensure both children are available before hashing
		if tree[index] != nil && tree[index+1] != nil {
			tree[i] = HashNode(tree[index], tree[index+1])
		} else if tree[index] != nil {
			// If right child is missing, use left child
			tree[i] = tree[index]
		}
	}

	// Ensure the root is properly set
	// The algorithm might not set tree[0] in some cases, especially for unbalanced trees
	// In those cases, tree[1] contains the correct root value
	if tree[0] == nil && tree[1] != nil {
		tree[0] = tree[1]
	}

	return tree, depth
}

// getDepth calculates the depth of a Merkle tree based on element count.
func getDepth(elementCount int) int {
	return int(math.Ceil(math.Log2(float64(elementCount))))
}

// getLeafCount returns the number of leafs in a balanced tree.
func getLeafCount(elementCount int) int {
	return int(RoundUpToPowerOf2(uint32(elementCount)))
}

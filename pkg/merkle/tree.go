// Package merkle implements the Merkle tree for the payer report.
package merkle

import (
	"errors"
	"fmt"
	"math/bits"
)

type Node []byte

type Leaf []byte

// MerkleTree is a binary Merkle tree.
//
// tree has the collection of nodes, where:
// - The tree is 1-indexed, so root is at index 1.
// - The internal nodes are at index 1 to N-1.
// - The leaves are at index N to 2N-1.
type MerkleTree struct {
	tree   []Node
	leaves []Leaf
}

var (
	EmptyTreeRoot         = make([]byte, 32)
	ErrNilLeaf            = errors.New("leaf is nil")
	ErrInvalidRange       = errors.New("invalid range")
	ErrIndicesOutOfBounds = errors.New("indices out of bounds")
)

// NewMerkleTree creates a new Merkle tree from the given leaves.
func NewMerkleTree(leaves []Leaf) (*MerkleTree, error) {
	nodes, err := makeLeafNodes(leaves)
	if err != nil {
		return nil, err
	}

	tree, err := makeTree(nodes)
	if err != nil {
		return nil, err
	}

	return &MerkleTree{tree, leaves}, nil
}

// GenerateMultiProofSequential generates a sequential multi-proof starting from the given index.
func (m *MerkleTree) GenerateMultiProofSequential(startingIndex, count int) (*MultiProof, error) {
	if startingIndex+count > m.LeafCount() {
		return nil, ErrIndicesOutOfBounds
	}

	indices, err := makeIndices(startingIndex, count)
	if err != nil {
		return nil, err
	}

	proof, err := m.makeProof(indices)
	if err != nil {
		return nil, err
	}

	if err := proof.validate(); err != nil {
		return nil, err
	}

	return &proof, nil
}

// makeProof returns a MultiProof for the given indices.
func (m *MerkleTree) makeProof(indices []int) (MultiProof, error) {
	var (
		balancedLeafCount = len(m.tree) >> 1
		proofElements     []ProofElement
		known             = make([]bool, len(m.tree))
		startingIndex     = 0
	)

	// Mark provided indices as known.
	for _, idx := range indices {
		known[balancedLeafCount+idx] = true
	}

	leafCountBytes, err := IntToBytes32(m.LeafCount())
	if err != nil {
		return MultiProof{}, err
	}

	// The first proof element is always the leaf count.
	proofElements = append(proofElements, ProofElement(leafCountBytes))

	// Calculate proofs to prove the existence of the indices.
	for i := balancedLeafCount - 1; i > 0; i-- {
		leftChildIdx := GetLeftChild(i)
		rightChildIdx := GetRightChild(i)

		left := known[leftChildIdx]
		right := known[rightChildIdx]

		// If at least one of the children is known, the parent is known.
		known[i] = left || right

		// If both children are known, we don't need to prove anything.
		// If neither child is known, we don't need to prove anything yet.
		if left == right {
			continue
		}

		// Only right child is known, so we need the left child as a proof.
		if right {
			proofElements = append(proofElements, ProofElement(m.tree[leftChildIdx]))
			continue
		}

		// If the right child is nil, we don't need to prove anything.
		if m.tree[rightChildIdx] == nil {
			continue
		}

		// Only left child is known, so we need the right child as a proof.
		proofElements = append(proofElements, ProofElement(m.tree[rightChildIdx]))
	}

	leaves := make([]Leaf, len(indices))
	for i, idx := range indices {
		leaves[i] = m.leaves[idx]
	}

	if balancedLeafCount != 0 && !known[1] {
		proofElements = append(proofElements, ProofElement(m.tree[1]))
	}

	if len(indices) != 0 {
		startingIndex = indices[0]
	}

	return MultiProof{
		startingIndex: startingIndex,
		leaves:        leaves,
		proofElements: proofElements,
	}, nil
}

// Tree returns the 1-indexed representation of the Merkle tree.
func (m *MerkleTree) Tree() []Node {
	return m.tree
}

// Leaves returns the raw leaves of the Merkle tree.
func (m *MerkleTree) Leaves() []Leaf {
	return m.leaves
}

// Root returns the root hash of the Merkle tree.
func (m *MerkleTree) Root() []byte {
	if len(m.tree) == 0 {
		return EmptyTreeRoot
	}

	return m.tree[0]
}

// LeafCount returns the number of leaves in the Merkle tree.
func (m *MerkleTree) LeafCount() int {
	return len(m.leaves)
}

// makeTree builds a serialized Merkle tree from an array of leaf nodes.
//
// For a tree that, when balanced, would have N leaves:
// The root is at index 0.
// The internal nodes are at index 1 to N-1.
// The leaves are at index N to 2N-1.
//
// For any node at index i:
// - left child is at index 2*i
// - right child is at index 2*i+1
// - parent is at floor(i/2)
func makeTree(leafNodes []Node) ([]Node, error) {
	leafCount := len(leafNodes)

	balancedLeafCount, err := CalculateBalancedNodesCount(leafCount)
	if err != nil {
		return nil, err
	}

	// Allocate 2N space for the tree. (N leaf nodes, N-1 internal nodes)
	tree := make([]Node, balancedLeafCount<<1)

	if leafCount == 0 {
		return tree, nil
	}

	lowerBound := balancedLeafCount
	upperBound := balancedLeafCount + leafCount - 1

	// Copy leaves into the tree, starting at index N.
	for i := 0; i < leafCount; i++ {
		tree[lowerBound+i] = leafNodes[i]
	}

	// Build the tree.
	for i := balancedLeafCount - 1; i > 0; i-- {
		leftChildIndex := GetLeftChild(i)

		if leftChildIndex > upperBound {
			continue
		}

		// If the left child is the last node in the level, we can use a pairless node.
		if leftChildIndex == upperBound {
			tree[i] = HashPairlessNode(tree[leftChildIndex])
			continue
		}

		// Detect the level is processed.
		if leftChildIndex <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		rightChildIndex := GetRightChild(i)

		tree[i] = HashNodePair(tree[leftChildIndex], tree[rightChildIndex])
	}

	tree[0], err = HashRoot(leafCount, tree[1])
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// makeLeafNodes returns the hashed leaves of the tree,
// ordered in the same order as the provided leaves.
func makeLeafNodes(leaves []Leaf) ([]Node, error) {
	nodes := make([]Node, len(leaves))
	for i, leaf := range leaves {
		if leaf == nil {
			return nil, ErrNilLeaf
		}

		nodes[i] = HashLeaf(leaf)
	}

	return nodes, nil
}

// CalculateBalancedNodesCount returns the number of nodes in a balanced tree.
// To calculate the number of nodes in a tree, we need to round up to the next power of 2.
// Returns an error if the leaf count is too large to be represented in a int32.
func CalculateBalancedNodesCount(count int) (int, error) {
	if count < 0 {
		return 0, fmt.Errorf("count cannot be negative")
	}

	if count > 1<<31-1 {
		return 0, fmt.Errorf("count must be less than or equal than max int32 (%d)", 1<<31-1)
	}

	if count == 0 {
		return 0, nil
	}

	// Despite 1 being a power of 2, a tree with 1 leaf is not a balanced tree in this implementation,
	// since a leaf is first hashed to a node using the leaf prefix, and then that node must be hashed into a node
	// using the node prefix, either with or without a paired node.
	if count == 1 {
		return 2, nil
	}

	return roundUpToPowerOf2(count), nil
}

// roundUpToPowerOf2 rounds up a number to the next power of 2.
func roundUpToPowerOf2(n int) int {
	if bits.OnesCount(uint(n)) == 1 {
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

// makeIndices returns a slice of ascending ordered indices for the given starting index and count.
func makeIndices(startingIndex, count int) ([]int, error) {
	if startingIndex < 0 || count < 0 {
		return nil, ErrInvalidRange
	}

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = startingIndex + i
	}

	return indices, nil
}

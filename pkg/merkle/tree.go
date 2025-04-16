package merkle

import (
	"errors"
	"fmt"
	"math/bits"
	"sort"
)

type Node struct {
	hash []byte
}

func (n *Node) Hash() []byte {
	return n.hash
}

func (n *Node) IsNil() bool {
	return n.hash == nil
}

type Leaf []byte

type MerkleTree struct {
	tree   []Node
	leaves []Leaf
}

var (
	ErrTreeEmpty          = errors.New("tree is empty")
	ErrTreeRootNil        = errors.New("tree root is nil")
	ErrTreeLeavesOverflow = errors.New("amount of leaves overflows uint32")
)

// NewMerkleTree creates a new Merkle tree from the given elements.
func NewMerkleTree(leaves []Leaf) (*MerkleTree, error) {
	if len(leaves) == 0 {
		return nil, ErrTreeEmpty
	}

	leavesDeepCopy := make([]Leaf, len(leaves))
	copy(leavesDeepCopy, leaves)

	nodes, err := makeNodes(leavesDeepCopy)
	if err != nil {
		return nil, err
	}

	tree, err := makeTree(nodes)
	if err != nil {
		return nil, err
	}

	return &MerkleTree{
		tree:   tree,
		leaves: leavesDeepCopy,
	}, nil
}

// GenerateMultiProofSequential generates a sequential multi-proof starting from the given index.
func (m *MerkleTree) GenerateMultiProofSequential(startingIndex, count int) (*MultiProof, error) {
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

// GenerateMultiProofWithIndices generates a multi-proof for the given indices.
func (m *MerkleTree) GenerateMultiProofWithIndices(indices []int) (*MultiProof, error) {
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	proof, err := m.makeProof(sortedIndices)
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
	if len(m.tree) == 0 {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, ErrEmptyTree)
	}

	if len(m.Root()) == 0 {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, ErrNilRoot)
	}

	if err := validateIndices(indices, m.LeafCount()); err != nil {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, err)
	}

	indexedValues := makeIndexedValues(m.leaves, indices)
	if len(indexedValues) != len(indices) {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, ErrElementMismatch)
	}

	// Handle single-element trees.
	if m.LeafCount() == 1 {
		return MultiProof{
			values:    indexedValues,
			proofs:    []Proof{m.Root()},
			leafCount: m.LeafCount(),
		}, nil
	}

	var (
		startLeafIdx = len(m.tree) >> 1
		proofs       []Proof
		known        = make([]bool, len(m.tree))
	)

	// Mark provided indices as known.
	for _, idx := range indices {
		known[startLeafIdx+idx] = true
	}

	// Calculate proofs to prove the existence of the indices.
	for i := startLeafIdx - 1; i > 0; i-- {
		leftChildIdx := GetLeftChild(i)
		rightChildIdx := GetRightChild(i)

		left := known[leftChildIdx]
		right := known[rightChildIdx]

		// Only one of children would be known, so we need the sibling as a proof
		if left != right {
			if right {
				// Only add non-nil sibling nodes
				if !m.tree[leftChildIdx].IsNil() {
					proofs = append(proofs, cloneBuffer(m.tree[leftChildIdx].Hash()))
				}
			} else {
				// Only add non-nil sibling nodes
				if !m.tree[rightChildIdx].IsNil() {
					proofs = append(proofs, cloneBuffer(m.tree[rightChildIdx].Hash()))
				}
			}
		}

		// If at least one of the children is known, the parent is known
		known[i] = left || right
	}

	return MultiProof{
		values:    indexedValues,
		proofs:    proofs,
		leafCount: len(m.leaves),
	}, nil
}

// Tree returns the 1-indexed representation of the Merkle tree.
func (m *MerkleTree) Tree() []Node {
	return m.tree
}

// Leaves returns the raw elements of the Merkle tree.
func (m *MerkleTree) Leaves() []Leaf {
	return m.leaves
}

// Root returns the root hash of the Merkle tree.
func (m *MerkleTree) Root() []byte {
	return m.tree[1].Hash()
}

// LeafCount returns the number of leaves in the Merkle tree.
func (m *MerkleTree) LeafCount() int {
	return len(m.leaves)
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
func makeTree(nodes []Node) ([]Node, error) {
	if len(nodes) == 0 {
		return nil, ErrTreeEmpty
	}

	leafCount, err := CalculateBalancedLeafCount(len(nodes))
	if err != nil {
		return nil, err
	}

	// Allocate 2N space for the tree. (N leaf nodes, N-1 internal nodes)
	tree := make([]Node, leafCount<<1)

	lowerBound := leafCount
	upperBound := leafCount + len(nodes) - 1

	// Copy leaves into the tree, starting at index N.
	for i := 0; i < len(nodes); i++ {
		tree[lowerBound+i] = nodes[i]
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
		if !tree[leftChildIndex].IsNil() && !tree[rightChildIndex].IsNil() {
			tree[i] = Node{
				hash: HashNode(tree[leftChildIndex].Hash(), tree[rightChildIndex].Hash()),
			}
		} else if !tree[leftChildIndex].IsNil() {
			tree[i] = tree[leftChildIndex]
		}
	}

	if tree[1].IsNil() {
		return nil, ErrTreeRootNil
	}

	return tree, nil
}

// makeNodes returns the leaves of the tree,
// ordered in the same order as the provided elements.
func makeNodes(leaves []Leaf) ([]Node, error) {
	if len(leaves) == 0 {
		return nil, ErrTreeEmpty
	}

	nodes := make([]Node, len(leaves))
	for i, element := range leaves {
		nodes[i] = Node{hash: HashLeaf(element)}
	}

	return nodes, nil
}

// CalculateBalancedLeafCount returns the number of leaves in a balanced tree.
// To calculate the number of leaves in a tree, we need to round up to the next power of 2.
// Returns an error if the element count is too large to be represented in a uint32.
func CalculateBalancedLeafCount(count int) (int, error) {
	if count <= 0 {
		return 0, nil
	}

	if count > int(1<<31) {
		return 0, ErrTreeLeavesOverflow
	}

	return int(roundUpToPowerOf2(uint32(count))), nil
}

// roundUpToPowerOf2 rounds up a number to the next power of 2.
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

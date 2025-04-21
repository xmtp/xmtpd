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
	if n.hash == nil {
		return nil
	}

	cp := make([]byte, len(n.hash))
	copy(cp, n.hash)

	return cp
}

type Leaf []byte

// MerkleTree is a binary Merkle tree.
//
// tree has the collection of nodes, where:
// - The tree is 1-indexed, so root is at index 1.
// - The internal nodes are at index 1 to N-1.
// - The leaves are at index N to 2N-1.
//
// leaves contains the raw elements of the tree.
type MerkleTree struct {
	tree   []Node
	leaves []Leaf
}

var (
	ErrNoLeaves  = errors.New("no leaves provided")
	ErrTreeEmpty = errors.New("tree is empty")
)

// NewMerkleTree creates a new Merkle tree from the given elements.
func NewMerkleTree(leaves []Leaf) (*MerkleTree, error) {
	if len(leaves) == 0 {
		return nil, ErrNoLeaves
	}

	balancedLeaves, err := makeLeaves(leaves)
	if err != nil {
		return nil, err
	}

	nodes, err := makeNodes(balancedLeaves)
	if err != nil {
		return nil, err
	}

	tree, err := makeTree(nodes)
	if err != nil {
		return nil, err
	}

	return &MerkleTree{tree, balancedLeaves}, nil
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
	err := validateIndices(indices, m.LeafCount())
	if err != nil {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, err)
	}

	indexedValues, err := makeIndexedValues(m.leaves, indices)
	if err != nil {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, err)
	}

	// Handle single-element trees.
	if m.LeafCount() == 1 {
		return MultiProof{
			values:    indexedValues,
			proofs:    []Node{{hash: m.Root()}},
			leafCount: m.LeafCount(),
		}, nil
	}

	var (
		startLeafIdx = len(m.tree) >> 1
		proofs       []Node
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

		// Only one of children is known, so we need the sibling as a proof.
		if left != right {
			if right {
				proofs = append(proofs, m.tree[leftChildIdx])
			} else {
				proofs = append(proofs, m.tree[rightChildIdx])
			}
		}

		// If at least one of the children is known, the parent is known.
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

// makeTree builds a serialized Merkle tree from an array of leaf nodes.
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

	balancedLeafCount, err := CalculateBalancedNodesCount(len(nodes))
	if err != nil {
		return nil, err
	}

	// Allocate 2N space for the tree. (N leaf nodes, N-1 internal nodes)
	tree := make([]Node, balancedLeafCount<<1)

	lowerBound := balancedLeafCount
	upperBound := balancedLeafCount + len(nodes) - 1

	// Copy leaves into the tree, starting at index N.
	for i := 0; i < len(nodes); i++ {
		tree[lowerBound+i] = nodes[i]
	}

	for i := len(nodes); i < balancedLeafCount; i++ {
		tree[lowerBound+i] = Node{hash: HashEmptyLeaf()}
	}

	// Build the tree.
	for i := balancedLeafCount - 1; i >= 1; i-- {
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

		tree[i] = Node{
			hash: HashNode(tree[leftChildIndex].Hash(), tree[rightChildIndex].Hash()),
		}

	}

	return tree, nil
}

// makeNodes returns the hashed leaves of the tree,
// ordered in the same order as the provided elements.
func makeNodes(leaves []Leaf) ([]Node, error) {
	if len(leaves) == 0 {
		return nil, ErrNoLeaves
	}

	nodes := make([]Node, len(leaves))
	for i, element := range leaves {
		nodes[i] = Node{hash: HashLeaf(element)}
	}

	return nodes, nil
}

// makeLeaves creates a balanced leaves count for the given leaves.
// It fills empty leaves with an empty leaf hash.
func makeLeaves(leaves []Leaf) ([]Leaf, error) {
	balancedLeafCount, err := CalculateBalancedNodesCount(len(leaves))
	if err != nil {
		return nil, err
	}

	balancedLeaves := make([]Leaf, balancedLeafCount)

	copy(balancedLeaves, leaves)

	for i := len(leaves); i < balancedLeafCount; i++ {
		balancedLeaves[i] = []byte{}
	}

	return balancedLeaves, nil
}

// CalculateBalancedNodesCount returns the number of nodes in a balanced tree.
// To calculate the number of nodes in a tree, we need to round up to the next power of 2.
// Returns an error if the element count is too large to be represented in a uint32.
func CalculateBalancedNodesCount(count int) (int, error) {
	if count <= 0 {
		return 0, fmt.Errorf("count must be greater than 0")
	}

	if count > int(1<<31) {
		return 0, fmt.Errorf("count must be less than or equal to 2^31")
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

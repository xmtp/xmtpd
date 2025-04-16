package merkle

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

const (
	ErrGenerateProof string = "cannot generate proof: %w"
	ErrVerifyProof   string = "cannot verify proof: %w"
)

var (
	ErrEmptyTree          = errors.New("empty tree")
	ErrDuplicateIndices   = errors.New("duplicate indices")
	ErrIndicesOutOfBounds = errors.New("indices out of bounds")
	ErrInvalidRange       = errors.New("invalid range")
	ErrInvalidLeafCount   = errors.New("invalid leaf count")
	ErrElementMismatch    = errors.New("element and indices mismatch")
	ErrNilProof           = errors.New("nil proof")
	ErrNilRoot            = errors.New("nil root")
	ErrNilElement         = errors.New("nil element")
	ErrNoElements         = errors.New("no elements")
	ErrNoIndices          = errors.New("no indices")
	ErrNoProofs           = errors.New("no proofs provided")
)

type MultiProof struct {
	values    IndexedValues
	proofs    []Proof
	leafCount int
}

type Proof []byte

type IndexedValues []IndexedValue

type IndexedValue struct {
	value []byte
	index int
}

func (iv IndexedValues) Indices() []int {
	indices := make([]int, len(iv))
	for i, v := range iv {
		indices[i] = v.index
	}
	return indices
}

func (iv IndexedValues) ToLeaves() []Leaf {
	leaves := make([]Leaf, len(iv))
	for i, v := range iv {
		leaves[i] = Leaf(v.value)
	}
	return leaves
}

// Verify verifies a MultiProof against the given tree root.
func Verify(root []byte, proof *MultiProof) (bool, error) {
	if len(root) == 0 {
		return false, fmt.Errorf(ErrVerifyProof, ErrNilRoot)
	}

	if err := proof.validate(); err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	// Handle single-element trees.
	if proof.leafCount == 1 {
		return bytes.Equal(root, HashLeaf(proof.values[0].value)), nil
	}

	// If all the elements are provided, we can directly reconstruct the tree.
	if len(proof.values) == proof.leafCount {
		leaves := proof.values.ToLeaves()

		tree, err := NewMerkleTree(leaves)
		if err != nil {
			return false, fmt.Errorf(ErrVerifyProof, err)
		}

		return bytes.Equal(tree.Root(), root), nil
	}

	computedRoot, err := proof.computeRoot()
	if err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	return bytes.Equal(computedRoot, root), nil
}

// validate performs common validation for Merkle proofs.
func (p *MultiProof) validate() error {
	if len(p.values) == 0 {
		return ErrNoElements
	}

	if p.leafCount <= 0 {
		return ErrInvalidLeafCount
	}

	if err := validateIndices(p.values.Indices(), p.leafCount); err != nil {
		return err
	}

	for _, elem := range p.values {
		if elem.value == nil {
			return ErrNilElement
		}
	}

	for _, proof := range p.proofs {
		if proof == nil {
			return ErrNilProof
		}
	}

	isPartialProof := len(p.values) < p.leafCount
	isNonTrivialTree := p.leafCount > 1
	needsProofs := isPartialProof && isNonTrivialTree

	if needsProofs && len(p.proofs) == 0 {
		return ErrNoProofs
	}

	return nil
}

// getNextProof safely retrieves the next proof and increments the index.
func (p *MultiProof) getNextProof(index *int) ([]byte, error) {
	if *index >= len(p.proofs) {
		return nil, ErrNilProof
	}
	proof := p.proofs[*index]
	*index++
	return proof, nil
}

// NodeQueue represents a node in the computation queue with its tree index and hash value.
// It's used during proof verification to track nodes as they are processed.
type NodeQueue struct {
	index int
	hash  []byte
}

// buildNodeQueue builds the node queue for the proof computation.
func (p *MultiProof) buildNodeQueue(balancedLeafCount int) ([]NodeQueue, error) {
	leaves := p.values.ToLeaves()

	nodes, err := makeNodes(leaves)
	if err != nil {
		return nil, err
	}

	indices := p.values.Indices()
	n := len(indices)

	queue := make([]NodeQueue, n)
	for i, idx := range indices {
		insertPos := n - 1 - i
		queue[insertPos] = NodeQueue{
			index: balancedLeafCount + idx,
			hash:  nodes[i].Hash(),
		}
	}

	return queue, nil
}

// computeRoot computes the root of the Merkle tree from the given proof.
func (p *MultiProof) computeRoot() ([]byte, error) {
	// 1. Prepare the queue.
	blc, err := CalculateBalancedLeafCount(p.leafCount)
	if err != nil {
		return nil, err
	}

	nodeQueue, err := p.buildNodeQueue(blc)
	if err != nil {
		return nil, err
	}

	var (
		head, proofIdx = 0, 0
		upperBound     = blc + p.leafCount - 1
		lowerBound     = nodeQueue[len(nodeQueue)-1].index
		left, right    []byte
	)

	// 2. Process queue until we hit the root.
	for head < len(nodeQueue) {
		current := nodeQueue[head]
		head++

		if current.index == 1 {
			return current.hash, nil
		}

		// Detect level-up.
		if current.index == lowerBound ||
			(head < len(nodeQueue) && nodeQueue[head].index == lowerBound) {
			lowerBound >>= 1
			upperBound >>= 1
		}

		// Detect solo-left at tree edge.
		if isLeftNodeAtUpperBound(current.index, upperBound) {
			nodeQueue = append(nodeQueue, NodeQueue{
				index: current.index >> 1,
				hash:  current.hash,
			})
			continue
		}

		if isLeftChild(current.index) {
			// Handle left-child branch.
			left = current.hash
			if head < len(nodeQueue) && nodeQueue[head].index == current.index+1 {
				right = nodeQueue[head].hash
				head++
			} else {
				right, err = p.getNextProof(&proofIdx)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// Handle right-child branch.
			right = current.hash
			if head < len(nodeQueue) && nodeQueue[head].index == current.index-1 {
				left = nodeQueue[head].hash
				head++
			} else {
				left, err = p.getNextProof(&proofIdx)
				if err != nil {
					return nil, err
				}
			}
		}

		parentHash := HashNode(left, right)

		nodeQueue = append(nodeQueue, NodeQueue{
			index: current.index >> 1,
			hash:  parentHash,
		})
	}

	return nil, ErrNilRoot
}

// makeIndexedValues creates indexed values from elements and their indices.
func makeIndexedValues(leaves []Leaf, indices []int) IndexedValues {
	result := make(IndexedValues, len(indices))
	for i, idx := range indices {
		result[i] = IndexedValue{
			value: leaves[idx],
			index: idx,
		}
	}
	return result
}

// makeIndices returns a slice of ascending ordered indices for the given starting index and count.
func makeIndices(startingIndex, count int) ([]int, error) {
	if startingIndex < 0 || count <= 0 {
		return nil, ErrInvalidRange
	}

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = startingIndex + i
	}

	return indices, nil
}

// validateIndices validates the indices slice of a proof.
func validateIndices(indices []int, leafCount int) error {
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	if len(sortedIndices) == 0 {
		return ErrNoIndices
	}

	if hasDuplicates(sortedIndices) {
		return ErrDuplicateIndices
	}

	if hasOutOfBounds(sortedIndices, leafCount) {
		return ErrIndicesOutOfBounds
	}

	return nil
}

func cloneBuffer(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}
	clone := make([]byte, len(buffer))
	copy(clone, buffer)
	return clone
}

// hasDuplicates checks if the sorted indices slice contains duplicates.
func hasDuplicates(indices []int) bool {
	for i := 1; i < len(indices); i++ {
		if indices[i] == indices[i-1] {
			return true
		}
	}
	return false
}

// hasOutOfBounds checks if all indices are within the valid range [0, leafCount).
func hasOutOfBounds(indices []int, leafCount int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= leafCount {
			return true
		}
	}
	return false
}

// isLeftChild returns true if the given index is odd (right child).
func isLeftChild(index int) bool {
	return index%2 == 0
}

// hasMissingSibling returns true if the node is a left child and its right sibling
// (index+1) is beyond the last real leaf.
// Note: this can happen because some leaves can be nil in unbalanced trees.
func hasMissingSibling(index int, upperBound int) bool {
	return isLeftChild(index) && index+1 > upperBound
}

// isLeftNodeAtUpperBound returns true if the node is a left child at the upper bound with no sibling.
// Left children have even indices in the 1-indexed tree representation.
func isLeftNodeAtUpperBound(index int, upperBound int) bool {
	return index == upperBound && hasMissingSibling(index, upperBound)
}

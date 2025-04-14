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
	elements  indexedValues
	proofs    [][]byte
	root      []byte
	leafCount int
}

func (p *MultiProof) Elements() indexedValues {
	return p.elements
}

func (p *MultiProof) Proofs() [][]byte {
	return p.proofs
}

func (p *MultiProof) Root() []byte {
	return p.root
}

func (p *MultiProof) LeafCount() int {
	return p.leafCount
}

type indexedValues []indexedValue

type indexedValue struct {
	value []byte
	index int
}

// Values returns all leaf values
func (iv indexedValues) Values() [][]byte {
	values := make([][]byte, len(iv))
	for i, v := range iv {
		values[i] = v.value
	}
	return values
}

// Indices returns all index values
func (iv indexedValues) Indices() []int {
	indices := make([]int, len(iv))
	for i, v := range iv {
		indices[i] = v.index
	}
	return indices
}

// Verify verifies a proof.
func Verify(proof *MultiProof) (bool, error) {
	if err := proof.validate(); err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	// Handle single-element trees.
	if proof.leafCount == 1 {
		return bytes.Equal(proof.root, HashLeaf(proof.elements[0].value)), nil
	}

	// If all the elements are provided, we can verify the proof by recalculating the root.
	if len(proof.elements) == proof.leafCount {
		// Extract just the leaves for tree creation
		tree, err := NewMerkleTree(proof.elements.Values())
		if err != nil {
			return false, fmt.Errorf(ErrVerifyProof, err)
		}

		return bytes.Equal(tree.Root(), proof.root), nil
	}

	result, err := proof.computeRoot()
	if err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	return bytes.Equal(result, proof.root), nil
}

// computeRoot computes the root given the leaves, their indices, and proofs.
func (p *MultiProof) computeRoot() ([]byte, error) {
	leaves, err := makeLeaves(p.elements.Values())
	if err != nil {
		return nil, err
	}

	balancedLeafCount := CalculateBalancedLeafCount(p.leafCount)

	// Prepare circular queues
	count := len(p.elements.Indices())
	treeIndices := make([]int, count)
	hashes := make([][]byte, count)

	// Initialize queues
	for i := 0; i < count; i++ {
		treeIndices[count-1-i] = balancedLeafCount + p.elements.Indices()[i]
		hashes[count-1-i] = cloneBuffer(leaves[i])
	}

	readIndex := 0
	writeIndex := 0
	proofIndex := 0
	upperBound := balancedLeafCount + p.leafCount - 1
	lowestTreeIndex := treeIndices[count-1]
	var nextNodeIndex int

	for {
		nodeIndex := treeIndices[readIndex]

		if nodeIndex == 1 {
			// Reached the root
			rootIndex := writeIndex - 1
			if writeIndex == 0 {
				rootIndex = count - 1
			}
			return hashes[rootIndex], nil
		}

		indexIsOdd := nodeIndex&1 == 1

		if nodeIndex == upperBound && !indexIsOdd {
			treeIndices[writeIndex] = nodeIndex >> 1
			hashes[writeIndex] = hashes[readIndex]
			writeIndex = (writeIndex + 1) % count
			readIndex = (readIndex + 1) % count
		} else {
			nextReadIndex := (readIndex + 1) % count
			if nextReadIndex < len(treeIndices) {
				nextNodeIndex = treeIndices[nextReadIndex]
			}

			// Check if the next node is a sibling
			nextIsPair := nextNodeIndex == nodeIndex-1

			var right, left []byte
			if indexIsOdd {
				right = hashes[readIndex]
				readIndex = (readIndex + 1) % count
				if !nextIsPair {
					if proofIndex >= len(p.proofs) {
						return nil, ErrNilProof
					}
					left = p.proofs[proofIndex]
					proofIndex++
				} else {
					left = hashes[readIndex]
					readIndex = (readIndex + 1) % count
				}
			} else {
				if proofIndex >= len(p.proofs) {
					return nil, ErrNilProof
				}
				right = p.proofs[proofIndex]
				proofIndex++
				left = hashes[readIndex]
				readIndex = (readIndex + 1) % count
			}

			if left == nil || right == nil {
				return nil, ErrNilProof
			}

			parentIndex := nodeIndex >> 1
			treeIndices[writeIndex] = parentIndex
			parentHash := HashNode(left, right)
			hashes[writeIndex] = parentHash
			writeIndex = (writeIndex + 1) % count
		}

		if nodeIndex == lowestTreeIndex || nextNodeIndex == lowestTreeIndex {
			lowestTreeIndex >>= 1
			upperBound >>= 1
		}
	}
}

// validate performs common validation for Merkle proofs.
func (p *MultiProof) validate() error {
	if p.root == nil {
		return ErrNilRoot
	}

	if len(p.elements) == 0 {
		return ErrNoElements
	}

	if err := validateIndices(p.elements.Indices(), p.leafCount); err != nil {
		return err
	}

	if p.leafCount <= 0 {
		return ErrInvalidLeafCount
	}

	for _, elem := range p.elements {
		if elem.value == nil {
			return ErrNilElement
		}
	}

	for _, proof := range p.proofs {
		if proof == nil {
			return ErrNilProof
		}
	}

	isPartialProof := len(p.elements) < p.leafCount
	isNonTrivialTree := p.leafCount > 1
	needsProofs := isPartialProof && isNonTrivialTree

	if needsProofs && len(p.proofs) == 0 {
		return ErrNoProofs
	}

	return nil
}

// makeIndexedValues creates indexed values from elements and their indices
func makeIndexedValues(elements [][]byte, indices []int) indexedValues {
	result := make(indexedValues, len(indices))
	for i, idx := range indices {
		result[i] = indexedValue{
			value: elements[idx],
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

// makeProof returns a MultiProof for the given indices.
func makeProof(
	tree [][]byte,
	root []byte,
	indices []int,
	leafCount int,
) (MultiProof, error) {
	if len(tree) == 0 {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, ErrEmptyTree)
	}

	if root == nil {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, ErrNilRoot)
	}

	if err := validateIndices(indices, leafCount); err != nil {
		return MultiProof{}, fmt.Errorf(ErrGenerateProof, err)
	}

	// Handle single-element trees.
	if leafCount == 1 {
		return MultiProof{
			root:      root,
			leafCount: leafCount,
			proofs:    [][]byte{root},
		}, nil
	}

	var (
		startLeafIdx = len(tree) >> 1
		proofs       [][]byte
		known        = make([]bool, len(tree))
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
				if tree[leftChildIdx] != nil {
					proofs = append(proofs, cloneBuffer(tree[leftChildIdx]))
				}
			} else {
				// Only add non-nil sibling nodes
				if tree[rightChildIdx] != nil {
					proofs = append(proofs, cloneBuffer(tree[rightChildIdx]))
				}
			}
		}

		// If at least one of the children is known, the parent is known
		known[i] = left || right
	}

	return MultiProof{
		root:      root,
		leafCount: leafCount,
		proofs:    proofs,
	}, nil
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

// hasDuplicates checks if the sorted indices slice contains duplicates.
func hasDuplicates(indices []int) bool {
	for i := 1; i < len(indices); i++ {
		if indices[i] == indices[i-1] {
			return true
		}
	}
	return false
}

// hasOutOfBounds checks if all indices are within the valid range [0, elementCount).
func hasOutOfBounds(indices []int, elementCount int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= elementCount {
			return true
		}
	}
	return false
}

func cloneBuffer(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}
	clone := make([]byte, len(buffer))
	copy(clone, buffer)
	return clone
}

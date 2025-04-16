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
	leafCount int
}

func (p *MultiProof) Elements() indexedValues {
	return p.elements
}

func (p *MultiProof) Proofs() [][]byte {
	return p.proofs
}

func (p *MultiProof) LeafCount() int {
	return p.leafCount
}

type indexedValues []indexedValue

type indexedValue struct {
	value []byte
	index int
}

func (iv indexedValues) Values() [][]byte {
	values := make([][]byte, len(iv))
	for i, v := range iv {
		values[i] = v.value
	}
	return values
}

func (iv indexedValues) Indices() []int {
	indices := make([]int, len(iv))
	for i, v := range iv {
		indices[i] = v.index
	}
	return indices
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
		return bytes.Equal(root, HashLeaf(proof.elements[0].value)), nil
	}

	// If all the elements are provided, we can directly reconstruct the tree.
	if len(proof.elements) == proof.leafCount {
		tree, err := NewMerkleTree(proof.elements.Values())
		if err != nil {
			return false, fmt.Errorf(ErrVerifyProof, err)
		}

		return bytes.Equal(tree.Root(), root), nil
	}

	result, err := proof.computeRoot()
	if err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	return bytes.Equal(result, root), nil
}

func (p *MultiProof) computeRoot() ([]byte, error) {
	balancedLeafCount, err := CalculateBalancedLeafCount(p.leafCount)
	if err != nil {
		return nil, err
	}

	leaves, err := makeLeaves(p.elements.Values())
	if err != nil {
		return nil, err
	}

	indices := p.elements.Indices()
	numElements := len(indices)

	// nodeQueue is populated with the indices we want to prove, from right to left.
	nodeQueue := make([]int, numElements)
	hashQueue := make([][]byte, numElements)
	for i := 0; i < numElements; i++ {
		reverseIdx := numElements - 1 - i
		nodeQueue[reverseIdx] = balancedLeafCount + indices[i]
		hashQueue[reverseIdx] = cloneBuffer(leaves[i])
	}

	var (
		readIndex, writeIndex, proofIndex = 0, 0, 0
		upperBound                        = balancedLeafCount + p.leafCount - 1
		lowestTreeIndex                   = nodeQueue[numElements-1]
		nextNodeIndex                     int
		leftHash, rightHash               []byte
	)

	for {
		currentIndex := nodeQueue[readIndex]

		// When the root is reached, return the hash.
		if currentIndex == 1 {
			rootIdx := (writeIndex + numElements - 1) % numElements
			return hashQueue[rootIdx], nil
		}

		// A left child is even, a right child is odd.
		isOdd := isOddIndex(currentIndex)

		// A left child at the boundary might not have a right sibling.
		// The sibling of a left child (even index) would be at index+1.
		shouldPromoteDirectly := isLeftNodeAtUpperBound(currentIndex, upperBound)

		// Promote node if it's a left child at the boundary with no sibling
		if shouldPromoteDirectly {
			nodeQueue[writeIndex] = currentIndex >> 1
			hashQueue[writeIndex] = hashQueue[readIndex]
			writeIndex = (writeIndex + 1) % numElements
			readIndex = (readIndex + 1) % numElements
		} else {
			nextReadIndex := (readIndex + 1) % numElements
			nextNodeIndex = nodeQueue[nextReadIndex]
			nextIsSibling := nextNodeIndex == currentIndex-1

			if isOdd {
				// Current node is right child
				rightHash = hashQueue[readIndex]
				readIndex = (readIndex + 1) % numElements

				if nextIsSibling {
					// Get left sibling from queue
					leftHash = hashQueue[readIndex]
					readIndex = (readIndex + 1) % numElements
				} else {
					// Get left sibling from proof
					leftHash, err = p.getNextProof(&proofIndex)
					if err != nil {
						return nil, err
					}
				}
			} else {
				// Current node is left child
				leftHash = hashQueue[readIndex]
				readIndex = (readIndex + 1) % numElements

				// Always get right sibling from proof
				rightHash, err = p.getNextProof(&proofIndex)
				if err != nil {
					return nil, err
				}
			}

			if leftHash == nil || rightHash == nil {
				return nil, ErrNilProof
			}

			// Compute parent hash and add to queue
			parentIndex := currentIndex >> 1
			parentHash := HashNode(leftHash, rightHash)
			nodeQueue[writeIndex] = parentIndex
			hashQueue[writeIndex] = parentHash
			writeIndex = (writeIndex + 1) % numElements
		}

		// Update level tracking when processing the lowest index.
		if currentIndex == lowestTreeIndex || nextNodeIndex == lowestTreeIndex {
			lowestTreeIndex >>= 1
			upperBound >>= 1
		}
	}
}

// getNextProof safely retrieves the next proof and increments the index.
func (p *MultiProof) getNextProof(proofIndex *int) ([]byte, error) {
	if *proofIndex >= len(p.proofs) {
		return nil, ErrNilProof
	}
	proof := p.proofs[*proofIndex]
	*proofIndex++
	return proof, nil
}

// validate performs common validation for Merkle proofs.
func (p *MultiProof) validate() error {
	if len(p.elements) == 0 {
		return ErrNoElements
	}

	if p.leafCount <= 0 {
		return ErrInvalidLeafCount
	}

	if err := validateIndices(p.elements.Indices(), p.leafCount); err != nil {
		return err
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

// makeIndexedValues creates indexed values from elements and their indices.
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

// hasOutOfBounds checks if all indices are within the valid range [0, elementCount).
func hasOutOfBounds(indices []int, elementCount int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= elementCount {
			return true
		}
	}
	return false
}

// isOddIndex returns true if the given index is odd (right child).
func isOddIndex(index int) bool {
	return index%2 == 1
}

// hasMissingSibling returns true if the node at the given index is missing its sibling.
// Only left children (even indices) at the upper bound might be missing their siblings.
func hasMissingSibling(index int, upperBound int) bool {
	return !isOddIndex(index) && index+1 > upperBound
}

// isLeftNodeAtUpperBound returns true if the node is a left child at the upper bound with no sibling.
// Left children have even indices in the 1-indexed tree representation.
func isLeftNodeAtUpperBound(index int, upperBound int) bool {
	return index == upperBound && hasMissingSibling(index, upperBound)
}

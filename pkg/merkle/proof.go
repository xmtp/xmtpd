package merkle

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type MultiProof struct {
	Elements     [][]byte
	Proofs       [][]byte
	Root         []byte
	Indices      []int
	ElementCount int
}

var (
	ErrProofEmptyTree            = errors.New("proof has empty tree")
	ErrProofEmptyIndices         = errors.New("proof has empty indices")
	ErrProofDuplicateIndices     = errors.New("proof has duplicate indices")
	ErrProofIndicesOutOfBounds   = errors.New("proof has indices out of bounds")
	ErrProofInvalidStartingIndex = errors.New("proof has invalid starting index")
	ErrProofInvalidElementCount  = errors.New("proof has invalid element count")
	ErrProofInvalidRange         = errors.New("proof has invalid range")
	ErrProofNil                  = errors.New("proof is nil")
	ErrProofNilRoot              = errors.New("proof root cannot be nil")
	ErrProofNoElements           = errors.New("proof has no elements")
	ErrProofNoIndices            = errors.New("proof has no indices")
	ErrProofNoProofs             = errors.New("proof has no proofs")
)

// GenerateMultiProofSequential generates a sequential multi-proof starting from the given index.
func (m *MerkleTree) GenerateMultiProofSequential(
	startingIndex, count int,
) (*MultiProof, error) {
	if startingIndex < 0 || startingIndex+count > m.leafCount {
		return nil, ErrProofInvalidRange
	}

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = startingIndex + i
	}

	proof, err := generateProof(m.tree, m.root, indices, m.leafCount)
	if err != nil {
		return nil, err
	}

	elements := make([][]byte, count)
	for i := 0; i < count; i++ {
		elements[i] = m.elements[startingIndex+i]
	}

	result := &MultiProof{
		Elements:     elements,
		Proofs:       proof.Proofs,
		Root:         proof.Root,
		Indices:      proof.Indices,
		ElementCount: proof.ElementCount,
	}

	return result, nil
}

// GenerateMultiProofWithIndices generates a multi-proof for the given indices.
func (m *MerkleTree) GenerateMultiProofWithIndices(indices []int) (*MultiProof, error) {
	proof, err := generateProof(m.tree, m.root, indices, m.leafCount)
	if err != nil {
		return nil, err
	}

	elements := make([][]byte, len(proof.Indices))
	for i, index := range proof.Indices {
		elements[i] = m.elements[index]
	}

	result := &MultiProof{
		Elements:     elements,
		Proofs:       proof.Proofs,
		Root:         proof.Root,
		Indices:      proof.Indices,
		ElementCount: proof.ElementCount,
	}

	return result, nil
}

// generateProof returns a MultiProof for the given indices.
func generateProof(
	tree [][]byte,
	root []byte,
	indices []int,
	elementCount int,
) (MultiProof, error) {
	if len(tree) == 0 {
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", ErrProofEmptyTree)
	}

	if root == nil {
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", ErrProofNilRoot)
	}

	// Do not modify the original indices slice.
	idxs := make([]int, len(indices))
	copy(idxs, indices)
	sort.Ints(idxs)

	if err := validateIndices(idxs, elementCount); err != nil {
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", err)
	}

	// Handle single-element trees.
	if elementCount == 1 {
		return MultiProof{
			Root:         root,
			Indices:      idxs,
			ElementCount: elementCount,
			Proofs:       [][]byte{root},
		}, nil
	}

	// Mark provided indices as known.
	var (
		leafCount = len(tree) >> 1
		proofs    [][]byte
		known     = make([]bool, len(tree))
	)

	for _, idx := range idxs {
		known[leafCount+idx] = true
	}

	// Calculate proofs to prove the existence of the indices.
	for i := leafCount - 1; i > 0; i-- {
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
		Root:         root,
		Indices:      idxs,
		ElementCount: elementCount,
		Proofs:       proofs,
	}, nil
}

// VerifyProof verifies a proof.
func VerifyProof(proof *MultiProof) (bool, error) {
	if err := validateProof(proof); err != nil {
		return false, fmt.Errorf("cannot verify proof: %w", err)
	}

	// Handle single-element trees.
	if proof.ElementCount == 1 {
		return bytes.Equal(proof.Root, HashLeaf(proof.Elements[0])), nil
	}

	// If all the elements are provided, we can verify the proof by recalculating the root.
	if len(proof.Elements) == proof.ElementCount {
		tree, err := NewMerkleTree(proof.Elements)
		if err != nil {
			return false, fmt.Errorf("cannot verify proof: %w", err)
		}

		return bytes.Equal(tree.Root(), proof.Root), nil
	}

	// If only some of the elements are provided, we need to calculate the root.
	leaves, err := makeLeaves(proof.Elements)
	if err != nil {
		return false, fmt.Errorf("cannot verify proof: %w", err)
	}

	result := getRoot(leaves, proof.Indices, proof.ElementCount, proof.Proofs)
	if result == nil {
		return false, fmt.Errorf("cannot verify proof: %w", ErrProofNilRoot)
	}

	return bytes.Equal(result, proof.Root), nil
}

// getRoot computes the root given the leaves, their indices, and proofs.
func getRoot(leaves [][]byte, indices []int, elementCount int, proofs [][]byte) []byte {
	// Ensure indices are valid
	for _, index := range indices {
		if index < 0 || index >= elementCount {
			return nil
		}
	}

	// Validate input
	if len(leaves) == 0 || len(indices) == 0 ||
		len(leaves) != len(indices) {
		return nil
	}

	// Sort indices and corresponding leaves
	indexLeafPairs := make([]struct {
		Index int
		Leaf  []byte
	}, len(indices))

	for i, index := range indices {
		indexLeafPairs[i] = struct {
			Index int
			Leaf  []byte
		}{Index: index, Leaf: leaves[i]}
	}

	sort.Slice(indexLeafPairs, func(i, j int) bool {
		return indexLeafPairs[i].Index < indexLeafPairs[j].Index
	})

	// Update sorted indices and leaves
	sortedIndices := make([]int, len(indexLeafPairs))
	sortedLeaves := make([][]byte, len(indexLeafPairs))
	for i, pair := range indexLeafPairs {
		sortedIndices[i] = pair.Index
		sortedLeaves[i] = pair.Leaf
	}

	// Original GetRoot implementation using balanced tree
	balancedLeafCount := int(roundUpToPowerOf2(uint32(elementCount)))

	// Prepare circular queues
	count := len(sortedIndices)
	treeIndices := make([]int, count)
	hashes := make([][]byte, count)

	// Initialize queues
	for i := 0; i < count; i++ {
		treeIndices[count-1-i] = balancedLeafCount + sortedIndices[i]
		hashes[count-1-i] = cloneBuffer(sortedLeaves[i])
	}

	readIndex := 0
	writeIndex := 0
	proofIndex := 0
	upperBound := balancedLeafCount + elementCount - 1
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
			return hashes[rootIndex]
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
					if proofIndex >= len(proofs) {
						return nil
					}
					left = proofs[proofIndex]
					proofIndex++
				} else {
					left = hashes[readIndex]
					readIndex = (readIndex + 1) % count
				}
			} else {
				if proofIndex >= len(proofs) {
					return nil
				}
				right = proofs[proofIndex]
				proofIndex++
				left = hashes[readIndex]
				readIndex = (readIndex + 1) % count
			}

			if left == nil || right == nil {
				return nil
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

// validateProof performs common validation for all types of Merkle proofs.
func validateProof(proof *MultiProof) error {
	if proof == nil {
		return ErrProofNil
	}

	if proof.Root == nil {
		return ErrProofNilRoot
	}

	if err := validateIndices(proof.Indices, proof.ElementCount); err != nil {
		return err
	}

	if len(proof.Indices) != len(proof.Elements) {
		return ErrProofInvalidElementCount
	}

	if len(proof.Elements) == 0 {
		return ErrProofNoElements
	}

	if proof.ElementCount <= 0 {
		return ErrProofInvalidElementCount
	}

	for i, element := range proof.Elements {
		if element == nil {
			return fmt.Errorf("nil element at index %d", i)
		}
	}

	if len(proof.Elements) < proof.ElementCount && proof.ElementCount > 1 {
		if len(proof.Proofs) == 0 {
			return ErrProofNoProofs
		}

		for i, decommitment := range proof.Proofs {
			if decommitment == nil {
				return fmt.Errorf("nil decommitment at index %d", i)
			}
		}
	}

	return nil
}

// validateIndices validates the indices slice of a proof.
func validateIndices(indices []int, elementCount int) error {
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	if len(sortedIndices) == 0 {
		return ErrProofEmptyIndices
	}

	if hasDuplicates(sortedIndices) {
		return ErrProofDuplicateIndices
	}

	if hasOutOfBounds(sortedIndices, elementCount) {
		return ErrProofIndicesOutOfBounds
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

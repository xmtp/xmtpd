package merkle

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type MultiProof struct {
	Elements  [][]byte
	Proofs    [][]byte
	Root      []byte
	Indices   []int
	LeafCount int
}

var (
	ErrProofEmptyTree           = errors.New("proof has empty tree")
	ErrProodIndicesDuplicated   = errors.New("proof has duplicate indices")
	ErrProofIndicesOutOfBounds  = errors.New("proof has indices out of bounds")
	ErrProofIndicesInvalidRange = errors.New("proof has invalid range")
	ErrProofInvalidElementCount = errors.New("proof has invalid element count")
	ErrProofElementMismatch     = errors.New("proof has a different number of indices and elements")
	ErrProofNil                 = errors.New("proof is nil")
	ErrProofNilRoot             = errors.New("proof root is nil")
	ErrProofNoElements          = errors.New("proof has no elements")
	ErrProofNoIndices           = errors.New("proof has no indices")
	ErrNoProofs                 = errors.New("no proofs provided")
)

// GenerateMultiProofSequential generates a sequential multi-proof starting from the given index.
func (m *MerkleTree) GenerateMultiProofSequential(startingIndex, count int) (*MultiProof, error) {
	indices, err := makeIndices(startingIndex, count)
	if err != nil {
		return nil, err
	}

	proof, err := makeProof(m.tree, m.root, indices, m.leafCount)
	if err != nil {
		return nil, err
	}

	proof.Elements = make([][]byte, count)
	for i := 0; i < count; i++ {
		proof.Elements[i] = m.elements[startingIndex+i]
	}

	if err := validateProof(&proof); err != nil {
		return nil, err
	}

	return &proof, nil
}

// GenerateMultiProofWithIndices generates a multi-proof for the given indices.
func (m *MerkleTree) GenerateMultiProofWithIndices(indices []int) (*MultiProof, error) {
	proof, err := makeProof(m.tree, m.root, indices, m.leafCount)
	if err != nil {
		return nil, err
	}

	proof.Elements = make([][]byte, len(proof.Indices))
	for i, index := range proof.Indices {
		proof.Elements[i] = m.elements[index]
	}

	if err := validateProof(&proof); err != nil {
		return nil, err
	}

	return &proof, nil
}

// VerifyProof verifies a proof.
func VerifyProof(proof *MultiProof) (bool, error) {
	if err := validateProof(proof); err != nil {
		return false, fmt.Errorf("cannot verify proof: %w", err)
	}

	// Handle single-element trees.
	if proof.LeafCount == 1 {
		return bytes.Equal(proof.Root, HashLeaf(proof.Elements[0])), nil
	}

	// If all the elements are provided, we can verify the proof by recalculating the root.
	if len(proof.Elements) == proof.LeafCount {
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

	result := computeRoot(leaves, proof.Indices, proof.LeafCount, proof.Proofs)
	if result == nil {
		return false, fmt.Errorf("cannot verify proof: %w", ErrProofNilRoot)
	}

	return bytes.Equal(result, proof.Root), nil
}

// computeRoot computes the root given the leaves, their indices, and proofs.
func computeRoot(leaves [][]byte, indices []int, elementCount int, proofs [][]byte) []byte {
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

// makeIndices returns a slice of indices for the given starting index and count.
func makeIndices(startingIndex, count int) ([]int, error) {
	if startingIndex < 0 || count <= 0 {
		return nil, ErrProofIndicesInvalidRange
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
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", ErrProofEmptyTree)
	}

	if root == nil {
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", ErrProofNilRoot)
	}

	// Do not modify the original indices slice.
	idxs := make([]int, len(indices))
	copy(idxs, indices)
	sort.Ints(idxs)

	if err := validateIndices(idxs, leafCount); err != nil {
		return MultiProof{}, fmt.Errorf("cannot generate proof: %w", err)
	}

	// Handle single-element trees.
	if leafCount == 1 {
		return MultiProof{
			Root:      root,
			Indices:   idxs,
			LeafCount: leafCount,
			Proofs:    [][]byte{root},
		}, nil
	}

	var (
		startLeafIdx = len(tree) >> 1
		proofs       [][]byte
		known        = make([]bool, len(tree))
	)

	// Mark provided indices as known.
	for _, idx := range idxs {
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
		Root:      root,
		Indices:   idxs,
		LeafCount: leafCount,
		Proofs:    proofs,
	}, nil
}

// validateProof performs common validation for Merkle proofs.
func validateProof(proof *MultiProof) error {
	if proof == nil {
		return ErrProofNil
	}

	if proof.Root == nil {
		return ErrProofNilRoot
	}

	if err := validateIndices(proof.Indices, proof.LeafCount); err != nil {
		return err
	}

	if len(proof.Elements) == 0 {
		return ErrProofNoElements
	}

	if proof.LeafCount <= 0 {
		return ErrProofInvalidElementCount
	}

	if len(proof.Indices) != len(proof.Elements) {
		return ErrProofElementMismatch
	}

	for i, e := range proof.Elements {
		if e == nil {
			return fmt.Errorf("nil element at index %d", i)
		}
	}

	for i, p := range proof.Proofs {
		if p == nil {
			return fmt.Errorf("nil proof at index %d", i)
		}
	}

	isPartialProof := len(proof.Elements) < proof.LeafCount
	isNonTrivialTree := proof.LeafCount > 1
	needsProofs := isPartialProof && isNonTrivialTree

	if needsProofs && len(proof.Proofs) == 0 {
		return ErrNoProofs
	}

	return nil
}

// validateIndices validates the indices slice of a proof.
func validateIndices(indices []int, elementCount int) error {
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	if len(sortedIndices) == 0 {
		return ErrProofNoIndices
	}

	if hasDuplicates(sortedIndices) {
		return ErrProodIndicesDuplicated
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

package merkle

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type MultiProof struct {
	Elements      [][]byte
	Proofs        [][]byte
	Root          []byte
	Indices       []int
	StartingIndex int
	ElementCount  int
}

var (
	ErrInputDuplicateIndices     = errors.New("input has duplicate indices")
	ErrInputIndicesOutOfBounds   = errors.New("input has indices out of bounds")
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

// generateProof returns a MultiProof for the given indices.
func generateProof(
	tree [][]byte,
	root []byte,
	indices []int,
	elementCount int,
) (MultiProof, error) {
	if len(tree) == 0 {
		return MultiProof{}, ErrProofEmptyTree
	}

	if root == nil {
		return MultiProof{}, ErrProofNilRoot
	}

	if len(indices) == 0 {
		return MultiProof{}, ErrProofEmptyIndices
	}

	// Do not modify the original indices slice.
	idxs := make([]int, len(indices))
	copy(idxs, indices)
	sort.Ints(idxs)

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
	leafCount := len(tree) >> 1
	known := make([]bool, len(tree))
	var proofs [][]byte

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

func verifyProof(
	proof *MultiProof,
	validateProof func(proof *MultiProof) error,
	getRoot func(leaves [][]byte, proofs [][]byte, startingIndex, elementCount int) []byte,
) (bool, error) {
	if err := validateProof(proof); err != nil {
		return false, err
	}

	// Handle single-element trees.
	if proof.ElementCount == 1 {
		return bytes.Equal(proof.Root, HashLeaf(proof.Elements[0])), nil
	}

	// If all the elements are provided, we can verify the proof by recalculating the root.
	if len(proof.Elements) == proof.ElementCount {
		tree, err := NewMerkleTree(proof.Elements)
		if err != nil {
			return false, err
		}

		return bytes.Equal(tree.Root(), proof.Root), nil
	}

	// If only some of the elements are provided, we need to calculate the root.
	leaves := make([][]byte, len(proof.Elements))
	for i, element := range proof.Elements {
		leaves[i] = HashLeaf(element)
	}

	result := getRoot(leaves, proof.Proofs, proof.StartingIndex, proof.ElementCount)
	if result == nil {
		return false, nil
	}

	return bytes.Equal(result, proof.Root), nil
}

// validateProofBase performs common validation for all types of Merkle proofs.
func validateProofBase(proof *MultiProof) error {
	if proof == nil {
		return ErrProofNil
	}

	if proof.Root == nil {
		return ErrProofNilRoot
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

// hasDuplicates checks if the sorted indices slice contains duplicates.
func hasDuplicates(indices []int) bool {
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	for i := 1; i < len(sortedIndices); i++ {
		if sortedIndices[i] == sortedIndices[i-1] {
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

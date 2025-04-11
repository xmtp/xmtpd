package merkle

import (
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
	ErrProofNilRoot              = errors.New("proof root cannot be nil")
	ErrProofNoElements           = errors.New("proof has no elements")
	ErrProofNoIndices            = errors.New("proof has no indices")
	ErrProofNoProofs             = errors.New("proof has no proofs")
	ErrProofInvalidStartingIndex = errors.New("proof has invalid starting index")
	ErrProofInvalidElementCount  = errors.New("proof has invalid element count")
	ErrProofInvalidRange         = errors.New("proof has invalid range")
)

// generateProof returns a MultiProof for the given indices.
func generateProof(tree [][]byte, indices []int, elementCount int) (MultiProof, error) {
	if len(tree) == 0 || tree[0] == nil {
		return MultiProof{}, fmt.Errorf("tree cannot be nil or empty")
	}

	if len(indices) == 0 {
		return MultiProof{}, fmt.Errorf("indices cannot be empty")
	}

	// Handle single-element trees.
	if elementCount == 1 {
		root := cloneBuffer(tree[0])
		return MultiProof{
			Root:         root,
			Indices:      indices,
			ElementCount: elementCount,
			Proofs:       [][]byte{root},
		}, nil
	}

	// Do not modify the original indices slice.
	idxs := make([]int, len(indices))
	copy(idxs, indices)
	sort.Ints(idxs)

	// Mark provided indices as known.
	leafCount := len(tree) >> 1
	known := make([]bool, len(tree))
	var proofs [][]byte

	for _, idx := range idxs {
		known[leafCount+idx] = true
	}

	// Calculate proofs to prove the existence of the indices.
	for i := leafCount - 1; i > 0; i-- {
		leftChildIdx := getLeftChild(i)
		rightChildIdx := getRightChild(i)

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
		Root:         cloneBuffer(tree[0]),
		Indices:      idxs,
		ElementCount: elementCount,
		Proofs:       proofs,
	}, nil
}

// validateProofBase performs common validation for all types of Merkle proofs.
// This covers the validation requirements shared between indices and sequential proofs.
func validateProofBase(proof *MultiProof) error {
	if proof.Root == nil {
		return fmt.Errorf("proof has nil root")
	}

	if len(proof.Elements) == 0 {
		return fmt.Errorf("proof has no elements")
	}

	if proof.ElementCount <= 0 {
		return fmt.Errorf("invalid element count: %d", proof.ElementCount)
	}

	for i, element := range proof.Elements {
		if element == nil {
			return fmt.Errorf("nil element at index %d", i)
		}
	}

	if len(proof.Elements) < proof.ElementCount && proof.ElementCount > 1 {
		if len(proof.Proofs) == 0 {
			return fmt.Errorf("partial proof has no decommitments")
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
func hasDuplicates(sortedIndices []int) bool {
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

// combineLeaves combines a set of leaf nodes into a single root hash.
func combineLeaves(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return nil
	}

	if len(leaves) == 1 {
		return leaves[0]
	}

	// Create a balanced tree
	level := leaves
	for len(level) > 1 {
		nextLevel := make([][]byte, 0, (len(level)+1)/2)

		for i := 0; i < len(level); i += 2 {
			if i+1 < len(level) {
				// Combine pairs
				nextLevel = append(nextLevel, HashNode(level[i], level[i+1]))
			} else {
				// Odd node out - propagate up
				nextLevel = append(nextLevel, level[i])
			}
		}

		level = nextLevel
	}

	return level[0]
}

func cloneBuffer(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}
	clone := make([]byte, len(buffer))
	copy(clone, buffer)
	return clone
}

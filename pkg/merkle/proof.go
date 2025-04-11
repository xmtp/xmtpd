package merkle

import (
	"errors"
	"fmt"
	"sort"
)

type MultiProof struct {
	// Common fields for all proofs
	Root         []byte
	Elements     [][]byte
	Indices      []int
	ElementCount int
	Proofs       [][]byte

	// Optional field for sequential proofs
	StartingIndex int
}

// generateProof returns a MultiProof for the given indices.
func generateProof(tree [][]byte, indices []int, elementCount int) (MultiProof, error) {
	if len(indices) == 0 {
		return MultiProof{}, fmt.Errorf("indices cannot be empty")
	}

	// Check for out-of-bounds indices
	for _, index := range indices {
		if index < 0 || index >= elementCount {
			return MultiProof{}, fmt.Errorf(
				"index %d is out of range [0, %d)",
				index,
				elementCount,
			)
		}
	}

	// Create a copy of indices to avoid modifying the original
	idxs := make([]int, len(indices))
	copy(idxs, indices)

	// Sort indices to ensure consistent processing
	sort.Ints(idxs)

	// Mark indices as known
	for i, idx := range idxs {
		if i > 0 && idxs[i-1] >= idx {
			return MultiProof{}, errors.New("indices must be in ascending order")
		}
	}

	leafCount := len(tree) >> 1
	known := make([]bool, len(tree))
	var proofs [][]byte

	// Mark indices as known
	for _, idx := range idxs {
		known[leafCount+idx] = true
	}

	// Calculate proofs
	for i := leafCount - 1; i > 0; i-- {
		leftChildIndex := i << 1
		left := known[leftChildIndex]
		right := known[leftChildIndex+1]

		// Only one of children would be known, so we need the sibling as a proof
		if left != right {
			if right {
				proofs = append(proofs, cloneBuffer(tree[leftChildIndex]))
			} else {
				proofs = append(proofs, cloneBuffer(tree[leftChildIndex+1]))
			}
		}

		// If at least one of the children is known, the parent is known
		known[i] = left || right
	}

	// Filter out nil proofs
	filteredProofs := make([][]byte, 0, len(proofs))
	for _, d := range proofs {
		if d != nil {
			filteredProofs = append(filteredProofs, d)
		}
	}

	// Special case: If we have no proofs (e.g., for a tree with a single element),
	// add a sentinel proof to allow verification to proceed
	if len(filteredProofs) == 0 && len(idxs) < elementCount {
		// Add the root itself as a sentinel
		if tree[0] != nil {
			filteredProofs = append(filteredProofs, cloneBuffer(tree[0]))
		}
	}

	// Get the root for verification
	var root []byte
	if tree[0] != nil {
		root = cloneBuffer(tree[0])
	}

	return MultiProof{
		Root:         root,
		Indices:      idxs,
		ElementCount: elementCount,
		Proofs:       filteredProofs,
	}, nil
}

// Helper function to clone a byte buffer
func cloneBuffer(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}
	clone := make([]byte, len(buffer))
	copy(clone, buffer)
	return clone
}

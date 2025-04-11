package merkle

import (
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

// generateProof returns a MultiProof for the given indices.
func generateProof(tree [][]byte, indices []int, elementCount int) (MultiProof, error) {
	if len(indices) == 0 {
		return MultiProof{}, fmt.Errorf("indices cannot be empty")
	}

	idxs := make([]int, len(indices))
	copy(idxs, indices)
	sort.Ints(idxs)

	if hasDuplicates(idxs) {
		return MultiProof{}, fmt.Errorf("found duplicate indices")
	}

	if hasOutOfBounds(idxs, elementCount) {
		return MultiProof{}, fmt.Errorf("found indices out of range")
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

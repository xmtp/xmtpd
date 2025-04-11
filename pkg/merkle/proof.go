package merkle

import (
	"errors"
	"fmt"
	"sort"
)

// GetRootResult holds the result of computing the root.
type GetRootResult struct {
	Root         []byte
	ElementCount int
}

// MultiProof represents both proof parameters and results
// It combines the fields from VerifyMultiProofParams and MultiProofResult
type MultiProof struct {
	// Common fields for all proofs
	Root          []byte
	Elements      [][]byte
	Indices       []int
	ElementCount  int
	Decommitments [][]byte

	// Optional field for sequential proofs
	StartingIndex int
}

// Common parameters for proof generation
type GenerateParams struct {
	Tree         [][]byte
	ElementCount int
	Indices      []int
}

// Generate creates decommitments for proving the existence of leaves at specified indices
func Generate(params GenerateParams) (MultiProof, error) {
	if len(params.Indices) == 0 {
		return MultiProof{}, fmt.Errorf("indices cannot be empty")
	}

	// Check for out-of-bounds indices
	for _, index := range params.Indices {
		if index < 0 || index >= params.ElementCount {
			return MultiProof{}, fmt.Errorf(
				"index %d is out of range [0, %d)",
				index,
				params.ElementCount,
			)
		}
	}

	// Create a copy of indices to avoid modifying the original
	indices := make([]int, len(params.Indices))
	copy(indices, params.Indices)

	// Sort indices to ensure consistent processing
	sort.Ints(indices)

	// Mark indices as known
	for i, idx := range indices {
		if i > 0 && indices[i-1] >= idx {
			return MultiProof{}, errors.New("indices must be in ascending order")
		}
	}

	leafCount := len(params.Tree) >> 1
	known := make([]bool, len(params.Tree))
	var decommitments [][]byte

	// Mark indices as known
	for _, idx := range indices {
		known[leafCount+idx] = true
	}

	// Calculate decommitments
	for i := leafCount - 1; i > 0; i-- {
		leftChildIndex := i << 1
		left := known[leftChildIndex]
		right := known[leftChildIndex+1]

		// Only one of children would be known, so we need the sibling as a decommitment
		if left != right {
			if right {
				decommitments = append(decommitments, cloneBuffer(params.Tree[leftChildIndex]))
			} else {
				decommitments = append(decommitments, cloneBuffer(params.Tree[leftChildIndex+1]))
			}
		}

		// If at least one of the children is known, the parent is known
		known[i] = left || right
	}

	// Filter out nil decommitments
	filteredDecommitments := make([][]byte, 0, len(decommitments))
	for _, d := range decommitments {
		if d != nil {
			filteredDecommitments = append(filteredDecommitments, d)
		}
	}

	// Special case: If we have no decommitments (e.g., for a tree with a single element),
	// add a sentinel decommitment to allow verification to proceed
	if len(filteredDecommitments) == 0 && len(indices) < params.ElementCount {
		// Add the root itself as a sentinel
		if params.Tree[0] != nil {
			filteredDecommitments = append(filteredDecommitments, cloneBuffer(params.Tree[0]))
		}
	}

	// Get the root for verification
	var root []byte
	if params.Tree[0] != nil {
		root = cloneBuffer(params.Tree[0])
	}

	return MultiProof{
		Root:          root,
		Indices:       indices,
		ElementCount:  params.ElementCount,
		Decommitments: filteredDecommitments,
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

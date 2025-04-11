package merkle

import (
	"bytes"
	"fmt"
)

// GetRootSequentiallyParams holds parameters for sequential root computation.
type GetRootSequentiallyParams struct {
	StartingIndex int
	Leafs         [][]byte
	ElementCount  int
	Decommitments [][]byte
}

// GenerateSequentialMultiProof generates a sequential multi-proof starting from the given index
func (m *MerkleTree) GenerateSequentialMultiProof(
	startingIndex, count int,
) (*MultiProof, error) {

	// Check if the range is valid
	if startingIndex < 0 || startingIndex+count > m.leafCount {
		return nil, fmt.Errorf(
			"invalid range: startingIndex=%d, count=%d, elementCount=%d",
			startingIndex,
			count,
			m.leafCount,
		)
	}

	// Extract elements at the specified indices
	elements := make([][]byte, count)
	for i := 0; i < count; i++ {
		elements[i] = m.elements[startingIndex+i]
	}

	// Generate sequential indices
	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = startingIndex + i
	}

	// Generate the proof
	proof, err := generateSequentialProof(
		m.tree,
		m.leafCount,
		startingIndex,
		count,
	)
	if err != nil {
		return nil, err
	}

	result := &MultiProof{
		Root:          proof.Root,
		Elements:      elements,
		Indices:       indices,
		StartingIndex: startingIndex,
		ElementCount:  m.leafCount,
		Decommitments: proof.Decommitments,
	}

	return result, nil
}

// VerifySequentialMultiProof verifies a sequential multi-proof
func VerifySequentialMultiProof(proof MultiProof) bool {
	if len(proof.Elements) == 0 || proof.StartingIndex < 0 || proof.ElementCount <= 0 {
		return false
	}

	// Check if the sequential range is valid (within bounds)
	count := len(proof.Elements)
	if proof.StartingIndex+count > proof.ElementCount {
		return false
	}

	// Special case: If this is a single-element tree or we're verifying all elements,
	// we don't need decommitments
	if len(proof.Elements) == proof.ElementCount || proof.ElementCount == 1 {
		// Just verify that the proof's root matches the recalculated root
		root := HashLeaf(proof.Elements[0])

		// For multiple elements, we need to combine them
		if len(proof.Elements) > 1 {
			leafs := make([][]byte, len(proof.Elements))
			for i, element := range proof.Elements {
				leafs[i] = HashLeaf(element)
			}

			// Combine the leaves into a root
			// This simplified approach only works when we have all elements
			root = combineLeaves(leafs)
		}

		if proof.Root == nil {
			return false
		}

		return bytes.Equal(root, proof.Root)
	}

	// If there's no decommitments for a normal case, it's invalid
	if len(proof.Decommitments) == 0 {
		return false
	}

	// Hash the elements with the prefix
	leafs := make([][]byte, len(proof.Elements))
	for i, element := range proof.Elements {
		leafs[i] = HashLeaf(element)
	}

	// Prepare parameters for GetRootSequentially
	getRootParams := GetRootSequentiallyParams{
		StartingIndex: proof.StartingIndex,
		Leafs:         leafs,
		ElementCount:  proof.ElementCount,
		Decommitments: proof.Decommitments,
	}

	// Compute the root
	result := getRootSequentially(getRootParams)

	// Handle nil cases
	if proof.Root == nil {
		return false
	}
	if result.Root == nil {
		return false
	}

	// Verify the root matches
	return bytes.Equal(result.Root, proof.Root)
}

// getRootSequentially computes the root given sequential leafs and decommitments.
func getRootSequentially(params GetRootSequentiallyParams) GetRootResult {
	// Validate input parameters
	if params.StartingIndex < 0 || len(params.Leafs) == 0 {
		return GetRootResult{Root: nil}
	}

	// Ensure starting index and count are within bounds
	count := len(params.Leafs)
	if params.StartingIndex+count > params.ElementCount {
		return GetRootResult{Root: nil}
	}

	elementCount := params.ElementCount
	decommitments := params.Decommitments

	// Validate decommitments
	if len(decommitments) == 0 {
		return GetRootResult{Root: nil}
	}

	balancedLeafCount := int(RoundUpToPowerOf2(uint32(elementCount)))

	// Prepare circular queues
	treeIndices := make([]int, count)
	hashes := make([][]byte, count)

	// Initialize hashes queue
	for i := 0; i < count; i++ {
		hashes[count-1-i] = cloneBuffer(params.Leafs[i])
	}

	readIndex := 0
	writeIndex := 0
	decommitmentIndex := 0
	upperBound := balancedLeafCount + elementCount - 1
	lowestTreeIndex := balancedLeafCount + params.StartingIndex

	var nodeIndex, nextNodeIndex int

	for {
		if upperBound >= elementCount {
			nodeIndex = lowestTreeIndex + count - 1 - readIndex
		} else {
			nodeIndex = treeIndices[readIndex]
		}

		if nodeIndex == 1 {
			// Reached the root
			rootIndex := writeIndex - 1
			if writeIndex == 0 {
				rootIndex = count - 1
			}
			return GetRootResult{
				Root:         hashes[rootIndex],
				ElementCount: elementCount,
			}
		}

		indexIsOdd := nodeIndex&1 == 1

		if nodeIndex == upperBound && !indexIsOdd {
			treeIndices[writeIndex] = nodeIndex >> 1
			hashes[writeIndex] = hashes[readIndex]
			writeIndex = (writeIndex + 1) % count
			readIndex = (readIndex + 1) % count
		} else {
			nextReadIndex := (readIndex + 1) % count
			if upperBound >= elementCount {
				nextNodeIndex = lowestTreeIndex + count - 1 - nextReadIndex
			} else {
				nextNodeIndex = treeIndices[nextReadIndex]
			}

			// Check if the next node is a sibling
			nextIsPair := nextNodeIndex == nodeIndex-1

			var right, left []byte
			if indexIsOdd {
				right = hashes[readIndex]
				readIndex = (readIndex + 1) % count
				if !nextIsPair {
					if decommitmentIndex >= len(decommitments) {
						// Not enough decommitments
						return GetRootResult{Root: nil}
					}
					left = decommitments[decommitmentIndex]
					decommitmentIndex++
				} else {
					left = hashes[readIndex]
					readIndex = (readIndex + 1) % count
				}
			} else {
				if decommitmentIndex >= len(decommitments) {
					// Not enough decommitments
					return GetRootResult{Root: nil}
				}
				right = decommitments[decommitmentIndex]
				decommitmentIndex++
				left = hashes[readIndex]
				readIndex = (readIndex + 1) % count
			}

			if left == nil || right == nil {
				return GetRootResult{Root: nil}
			}

			treeIndices[writeIndex] = nodeIndex >> 1
			hashes[writeIndex] = HashNode(left, right)
			writeIndex = (writeIndex + 1) % count
		}

		if nodeIndex == lowestTreeIndex || nextNodeIndex == lowestTreeIndex {
			lowestTreeIndex >>= 1
			upperBound >>= 1
		}
	}
}

// generateSequentialProof generates a proof for sequential leaves
func generateSequentialProof(
	tree [][]byte,
	elementCount, startingIndex, count int,
) (MultiProof, error) {
	// Validate parameters
	if startingIndex < 0 || count <= 0 || startingIndex+count > elementCount {
		return MultiProof{}, fmt.Errorf(
			"invalid range: startingIndex=%d, count=%d, elementCount=%d",
			startingIndex,
			count,
			elementCount,
		)
	}

	balancedLeafCount := int(RoundUpToPowerOf2(uint32(elementCount)))
	known := make([]bool, len(tree))

	// Mark indices as known
	for i := 0; i < count; i++ {
		idx := startingIndex + i
		if idx < elementCount {
			known[balancedLeafCount+idx] = true
		}
	}

	var decommitments [][]byte

	// Calculate decommitments
	for i := balancedLeafCount - 1; i > 0; i-- {
		leftChildIndex := i << 1
		rightChildIndex := leftChildIndex + 1

		// Check if children are valid indices
		if leftChildIndex >= len(tree) || rightChildIndex >= len(tree) {
			continue
		}

		left := known[leftChildIndex]
		right := known[rightChildIndex]

		// Only one of children would be known, so we need the sibling as a decommitment
		if left != right {
			if right {
				decommitments = append(decommitments, cloneBuffer(tree[leftChildIndex]))
			} else {
				decommitments = append(decommitments, cloneBuffer(tree[rightChildIndex]))
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

	// Special case for empty decommitments:
	// If all sequential elements are provided, we still need at least one decommitment for verification
	if len(filteredDecommitments) == 0 && count < elementCount {
		// Add a sentinel decommitment (using the root itself)
		// This isn't ideal but allows verification to proceed
		if tree[0] != nil {
			filteredDecommitments = append(filteredDecommitments, cloneBuffer(tree[0]))
		}
	}

	// Get the root for verification
	var root []byte
	if tree[0] != nil {
		root = cloneBuffer(tree[0])
	}

	return MultiProof{
		Root:          root,
		ElementCount:  elementCount,
		Decommitments: filteredDecommitments,
		StartingIndex: startingIndex,
	}, nil
}

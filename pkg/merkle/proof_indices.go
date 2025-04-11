package merkle

import (
	"bytes"
	"fmt"
	"sort"
)

// GenerateMultiProofWithIndices generates a multi-proof for the given indices
func (m *MerkleTree) GenerateMultiProofWithIndices(indices []int) (*MultiProof, error) {
	if hasDuplicates(indices) {
		return nil, fmt.Errorf("found duplicate indices")
	}

	if hasOutOfBounds(indices, m.leafCount) {
		return nil, fmt.Errorf("found indices out of range")
	}

	proof, err := generateProof(m.tree, indices, m.leafCount)
	if err != nil {
		return nil, err
	}

	// Extract elements at the specified indices.
	// They are provided in the proof.
	elements := make([][]byte, len(indices))
	for i, index := range indices {
		elements[i] = m.elements[index]
	}

	result := &MultiProof{
		Elements:      elements,
		Proofs:        proof.Proofs,
		Root:          proof.Root,
		Indices:       indices,
		StartingIndex: indices[0],
		ElementCount:  m.leafCount,
	}

	return result, nil
}

// TODO: Abstract VerifyMultiProofSequential and VerifyMultiProofWithIndices to use a common function.
func VerifyMultiProofWithIndices(proof *MultiProof) (bool, error) {
	if err := validateProofIndices(proof); err != nil {
		return false, err
	}

	// Special case: If this is a single-element tree or we're verifying all elements
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

		return bytes.Equal(root, proof.Root), nil
	}

	// Hash the elements with the prefix
	leafs := make([][]byte, len(proof.Elements))
	for i, element := range proof.Elements {
		leafs[i] = HashLeaf(element)
	}

	result := getRootIndices(leafs, proof.Indices, proof.ElementCount, proof.Proofs)
	if result == nil {
		return false, nil
	}

	return bytes.Equal(result, proof.Root), nil
}

// getRootIndices computes the root given the leaves, their indices, and proofs
func getRootIndices(leafs [][]byte, indices []int, elementCount int, proofs [][]byte) []byte {
	// Ensure indices are valid
	for _, index := range indices {
		if index < 0 || index >= elementCount {
			return nil
		}
	}

	// Validate input
	if len(leafs) == 0 || len(indices) == 0 ||
		len(leafs) != len(indices) {
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
		}{Index: index, Leaf: leafs[i]}
	}

	sort.Slice(indexLeafPairs, func(i, j int) bool {
		return indexLeafPairs[i].Index < indexLeafPairs[j].Index
	})

	// Update sorted indices and leaves
	sortedIndices := make([]int, len(indexLeafPairs))
	sortedLeafs := make([][]byte, len(indexLeafPairs))
	for i, pair := range indexLeafPairs {
		sortedIndices[i] = pair.Index
		sortedLeafs[i] = pair.Leaf
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
		hashes[count-1-i] = cloneBuffer(sortedLeafs[i])
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

// validateProofIndices validates a proof with arbitrary indices.
// It handles specific validation for non-sequential proofs.
func validateProofIndices(proof *MultiProof) error {
	// Run common validation first
	if err := validateProofBase(proof); err != nil {
		return err
	}

	// Indices validation
	if len(proof.Indices) != len(proof.Elements) {
		return fmt.Errorf("indices count doesn't match elements count")
	}

	// Check for out-of-bounds indices
	for _, idx := range proof.Indices {
		if idx < 0 || idx >= proof.ElementCount {
			return fmt.Errorf("index %d is out of range [0, %d)", idx, proof.ElementCount)
		}
	}

	// Check for duplicate indices
	// First, make a copy to avoid modifying the original
	sortedIndices := make([]int, len(proof.Indices))
	copy(sortedIndices, proof.Indices)
	sort.Ints(sortedIndices)

	for i := 1; i < len(sortedIndices); i++ {
		if sortedIndices[i] == sortedIndices[i-1] {
			return fmt.Errorf("duplicate index found: %d", sortedIndices[i])
		}
	}

	return nil
}

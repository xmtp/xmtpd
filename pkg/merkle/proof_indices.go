package merkle

import (
	"bytes"
	"fmt"
	"sort"
)

// GetRootParams holds parameters for computing the root
type GetRootParams struct {
	Leafs         [][]byte
	Indices       []int
	ElementCount  int
	Decommitments [][]byte
}

// GenerateIndicesMultiProof generates a multi-proof for the given indices
func (m *MerkleTree) GenerateIndicesMultiProof(indices []int) (*MultiProof, error) {
	// Ensure indices are valid
	for _, index := range indices {
		if index < 0 || index >= m.leafCount {
			return nil, fmt.Errorf("index %d is out of range [0, %d)", index, m.leafCount)
		}
	}

	// Extract elements at the specified indices
	elements := make([][]byte, len(indices))
	for i, index := range indices {
		elements[i] = m.elements[index]
	}

	// Generate the proof
	generateParams := GenerateParams{
		Tree:         m.tree,
		Indices:      indices,
		ElementCount: m.leafCount,
	}

	proof, err := Generate(generateParams)
	if err != nil {
		return nil, err
	}

	result := &MultiProof{
		Root:          proof.Root,
		Elements:      elements,
		Indices:       indices,
		ElementCount:  m.leafCount,
		Decommitments: proof.Decommitments,
	}

	return result, nil
}

// getRootIndices computes the root given the leaves, their indices, and decommitments
func getRootIndices(params GetRootParams) GetRootResult {
	elementCount := params.ElementCount
	decommitments := params.Decommitments

	// Ensure indices are valid
	for _, index := range params.Indices {
		if index < 0 || index >= elementCount {
			return GetRootResult{Root: nil}
		}
	}

	// Validate input
	if len(params.Leafs) == 0 || len(params.Indices) == 0 ||
		len(params.Leafs) != len(params.Indices) {
		return GetRootResult{Root: nil}
	}

	// Sort indices and corresponding leaves
	indexLeafPairs := make([]struct {
		Index int
		Leaf  []byte
	}, len(params.Indices))

	for i, index := range params.Indices {
		indexLeafPairs[i] = struct {
			Index int
			Leaf  []byte
		}{Index: index, Leaf: params.Leafs[i]}
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
	balancedLeafCount := int(RoundUpToPowerOf2(uint32(elementCount)))

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
	decommitmentIndex := 0
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
					if decommitmentIndex >= len(decommitments) {
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

// VerifyMultiProof verifies a multi-proof
func VerifyMultiProof(proof MultiProof) bool {
	if len(proof.Elements) == 0 || len(proof.Indices) == 0 || proof.ElementCount <= 0 {
		return false
	}

	// Check if indices are valid (within bounds)
	for _, index := range proof.Indices {
		if index < 0 || index >= proof.ElementCount {
			return false
		}
	}

	// Check that we have the same number of elements as indices
	if len(proof.Elements) != len(proof.Indices) {
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

	// Prepare parameters for GetRoot
	getRootParams := GetRootParams{
		Leafs:         leafs,
		Indices:       proof.Indices,
		ElementCount:  proof.ElementCount,
		Decommitments: proof.Decommitments,
	}

	// Compute the root
	result := getRootIndices(getRootParams)

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

// combineLeaves combines a set of leaf nodes into a single root hash
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

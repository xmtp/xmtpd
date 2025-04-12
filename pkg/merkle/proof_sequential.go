package merkle

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
		Elements:      elements,
		Proofs:        proof.Proofs,
		Root:          proof.Root,
		Indices:       proof.Indices,
		StartingIndex: proof.Indices[0],
		ElementCount:  proof.ElementCount,
	}

	return result, nil
}

// VerifyMultiProofSequential verifies a sequential multi-proof.
func VerifyMultiProofSequential(proof *MultiProof) (bool, error) {
	return verifyProof(
		proof,
		validateProofSequential,
		getRootSequentially,
	)
}

// getRootSequentially computes the root given sequential leafs and proofs.
func getRootSequentially(leafs [][]byte, proofs [][]byte, startingIndex, elementCount int) []byte {
	// Validate input parameters
	if startingIndex < 0 || len(leafs) == 0 {
		return nil
	}

	// Ensure starting index and count are within bounds
	count := len(leafs)
	if startingIndex+count > elementCount {
		return nil
	}

	// Validate proofs
	if len(proofs) == 0 {
		return nil
	}

	balancedLeafCount := int(roundUpToPowerOf2(uint32(elementCount)))

	// Prepare circular queues
	treeIndices := make([]int, count)
	hashes := make([][]byte, count)

	// Initialize hashes queue
	for i := 0; i < count; i++ {
		hashes[count-1-i] = cloneBuffer(leafs[i])
	}

	readIndex := 0
	writeIndex := 0
	proofIndex := 0
	upperBound := balancedLeafCount + elementCount - 1
	lowestTreeIndex := balancedLeafCount + startingIndex

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
					if proofIndex >= len(proofs) {
						// Not enough proofs
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
					// Not enough proofs
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

// validateProofSequential validates a sequential proof.
func validateProofSequential(proof *MultiProof) error {
	if err := validateProof(proof); err != nil {
		return err
	}

	if proof.StartingIndex < 0 {
		return ErrProofInvalidStartingIndex
	}

	if proof.StartingIndex+len(proof.Elements) > proof.ElementCount {
		return ErrProofInvalidRange
	}

	return nil
}

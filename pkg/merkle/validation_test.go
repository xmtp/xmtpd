package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnbalancedTrees tests the correctness of Merkle trees with various unbalanced structures
func TestUnbalancedTrees(t *testing.T) {
	// Test cases with different numbers of elements
	testCases := []struct {
		name     string
		elements [][]byte
	}{
		{"SingleElement", [][]byte{[]byte("A")}},
		{"TwoElements", [][]byte{[]byte("A"), []byte("B")}},
		{"ThreeElements", [][]byte{[]byte("A"), []byte("B"), []byte("C")}},
		{"FiveElements", [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D"), []byte("E")}},
		{"SevenElements", [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D"),
			[]byte("E"), []byte("F"), []byte("G")}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Merkle tree from the elements
			tree, err := NewMerkleTree(tc.elements)
			require.NoError(t, err)
			require.NotNil(t, tree)

			// Ensure the root is not nil
			require.NotNil(t, tree.Root())

			// Validate that we can generate proofs for each element
			for i := range tc.elements {
				proof, err := tree.GenerateIndicesMultiProof([]int{i})
				require.NoError(t, err)
				require.NotNil(t, proof)

				// Verify the proof
				result := VerifyMultiProof(MultiProof{
					Root:          proof.Root,
					Elements:      proof.Elements,
					Indices:       proof.Indices,
					ElementCount:  proof.ElementCount,
					Decommitments: proof.Decommitments,
				})

				assert.True(t, result, "Proof verification failed for element %d", i)
			}

			// Test sequential proof for all elements
			if len(tc.elements) > 0 {
				seqProof, err := tree.GenerateSequentialMultiProof(0, len(tc.elements))
				require.NoError(t, err)
				require.NotNil(t, seqProof)

				// Verify sequential proof
				result := VerifySequentialMultiProof(MultiProof{
					Root:          seqProof.Root,
					Elements:      seqProof.Elements,
					StartingIndex: seqProof.StartingIndex,
					ElementCount:  seqProof.ElementCount,
					Decommitments: seqProof.Decommitments,
				})

				assert.True(t, result, "Sequential proof verification failed")
			}
		})
	}
}

// TestDuplicateElements ensures the tree handles duplicate elements correctly
func TestDuplicateElements(t *testing.T) {
	// Create a tree with duplicate elements
	elements := [][]byte{
		[]byte("A"),
		[]byte("A"), // Duplicate
		[]byte("B"),
		[]byte("C"),
		[]byte("B"), // Duplicate
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)

	// Ensure we can generate and verify proofs for all elements
	for i := range elements {
		proof, err := tree.GenerateIndicesMultiProof([]int{i})
		require.NoError(t, err)

		result := VerifyMultiProof(MultiProof{
			Root:          proof.Root,
			Elements:      proof.Elements,
			Indices:       proof.Indices,
			ElementCount:  proof.ElementCount,
			Decommitments: proof.Decommitments,
		})

		assert.True(t, result, "Proof verification failed for element %d", i)
	}
}

// TestLargeUnbalancedTree tests a significantly unbalanced tree
func TestLargeUnbalancedTree(t *testing.T) {
	// Create a tree with elements that's not a power of 2 (e.g., 100 elements)
	elements := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		elements[i] = []byte{byte(i)}
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)
	require.NotNil(t, tree.Root())

	// Test a few random indices
	indices := []int{0, 50, 99}
	for _, idx := range indices {
		proof, err := tree.GenerateIndicesMultiProof([]int{idx})
		require.NoError(t, err)

		result := VerifyMultiProof(MultiProof{
			Root:          proof.Root,
			Elements:      proof.Elements,
			Indices:       proof.Indices,
			ElementCount:  proof.ElementCount,
			Decommitments: proof.Decommitments,
		})

		assert.True(t, result, "Proof verification failed for element %d", idx)
	}

	// Test multi-index proof
	multiProof, err := tree.GenerateIndicesMultiProof(indices)
	require.NoError(t, err)

	multiResult := VerifyMultiProof(MultiProof{
		Root:          multiProof.Root,
		Elements:      multiProof.Elements,
		Indices:       multiProof.Indices,
		ElementCount:  multiProof.ElementCount,
		Decommitments: multiProof.Decommitments,
	})

	assert.True(t, multiResult, "Multi-index proof verification failed")
}

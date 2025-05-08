package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidate tests the internal validate function for MultiProof.
func TestValidate(t *testing.T) {
	t.Run("valid proof", func(t *testing.T) {
		leaves := []Leaf{[]byte("test")}

		proofElements := []ProofElement{IntToBytes32(2), []byte("proof")}

		proof := MultiProof{
			startingIndex: 0,
			leaves:        leaves,
			proofElements: proofElements,
		}

		err := proof.validate()
		assert.NoError(t, err)
	})

	t.Run("invalid starting index", func(t *testing.T) {
		proof := MultiProof{
			startingIndex: -1,
			leaves:        []Leaf{[]byte("test")},
			proofElements: []ProofElement{IntToBytes32(2), []byte("proof")},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrInvalidStartingIndex)
	})

	t.Run("no proofs", func(t *testing.T) {
		proof := MultiProof{
			startingIndex: 0,
			leaves:        []Leaf{[]byte("test")},
			proofElements: []ProofElement{},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNoProofs)
	})

	t.Run("insufficient proofs", func(t *testing.T) {
		proof := MultiProof{
			startingIndex: 0,
			leaves:        []Leaf{},
			proofElements: []ProofElement{IntToBytes32(1)},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrInsufficientProofs)
	})

	t.Run("indices out of bounds for empty tree", func(t *testing.T) {
		leaves := []Leaf{[]byte("test")}

		proof := MultiProof{
			startingIndex: 0,
			leaves:        leaves,
			proofElements: []ProofElement{IntToBytes32(0), []byte("proof")},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrIndicesOutOfBounds)
	})

	t.Run("indices out of bounds", func(t *testing.T) {
		proof := MultiProof{
			startingIndex: 2,
			leaves:        []Leaf{[]byte("test")},
			proofElements: []ProofElement{IntToBytes32(2), []byte("proof")},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrIndicesOutOfBounds)
	})

	t.Run("nil leaf", func(t *testing.T) {
		leaves := []Leaf{nil}

		proof := MultiProof{
			startingIndex: 0,
			leaves:        leaves,
			proofElements: []ProofElement{IntToBytes32(2), []byte("proof")},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNilLeaf)
	})

	t.Run("nil proof", func(t *testing.T) {
		leaves := []Leaf{[]byte("test")}

		proof := MultiProof{
			startingIndex: 0,
			leaves:        leaves,
			proofElements: []ProofElement{IntToBytes32(2), nil},
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNilProof)
	})

	t.Run("only 1 proof element needed for single leaf tree", func(t *testing.T) {
		leaves := []Leaf{[]byte("test")}

		proof := MultiProof{
			startingIndex: 0,
			leaves:        leaves,
			proofElements: []ProofElement{IntToBytes32(1)},
		}

		err := proof.validate()
		assert.NoError(t, err)
	})
}

// TestGetNextProofElement tests the internal getNextProofElement function.
func TestGetNextProofElement(t *testing.T) {
	proofElements := []ProofElement{
		IntToBytes32(2),
		[]byte{0x01},
		[]byte{0x02},
	}

	proof := MultiProof{
		startingIndex: 0,
		leaves:        []Leaf{[]byte("test")},
		proofElements: proofElements,
	}

	idx := 0

	p, err := proof.getNextProofElement(&idx)
	require.NoError(t, err)
	assert.Equal(t, ProofElement(IntToBytes32(2)), p)
	assert.Equal(t, 1, idx)

	p, err = proof.getNextProofElement(&idx)
	require.NoError(t, err)
	assert.Equal(t, ProofElement([]byte{0x01}), p)
	assert.Equal(t, 2, idx)

	// Test getting the second proof element.
	p, err = proof.getNextProofElement(&idx)
	require.NoError(t, err)
	assert.Equal(t, ProofElement([]byte{0x02}), p)
	assert.Equal(t, 3, idx)

	// Test error when out of proof elements.
	p, err = proof.getNextProofElement(&idx)
	assert.ErrorIs(t, err, ErrNoProofs)
	assert.Nil(t, p)
}

// TestBuildNodeQueue tests the internal buildNodeQueue function.
func TestBuildNodeQueue(t *testing.T) {
	leaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	// We don't actually need the tree for the test, just the leaves.
	_, err := NewMerkleTree(leaves)
	require.NoError(t, err)

	proof := MultiProof{
		startingIndex: 0,
		leaves:        leaves[0:2],
		proofElements: []ProofElement{},
	}

	// Build the node queue with 4 leaves, as it's the balanced leaf count for 4 leaves.
	queue, err := proof.buildNodeQueue(4)
	require.NoError(t, err)
	assert.Equal(t, 2, len(queue))

	// Check queue order and values.
	assert.Equal(t, 5, queue[0].index) // 4+1 (balanced leaf count + index)
	assert.Equal(t, HashLeaf(leaves[1]), queue[0].hash)
	assert.Equal(t, 4, queue[1].index) // 4+0 (balanced leaf count + index)
	assert.Equal(t, HashLeaf(leaves[0]), queue[1].hash)
}

// TestComputeRoot tests the internal computeRoot function.
func TestComputeRoot(t *testing.T) {
	// Create a simple tree for testing.
	leaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree, err := NewMerkleTree(leaves)
	require.NoError(t, err)

	expectedRoot := tree.Root()

	// Test with valid proof.
	multiProof, err := tree.GenerateMultiProofSequential(0, 2)
	require.NoError(t, err)

	computedRoot, err := multiProof.computeRoot()
	require.NoError(t, err)
	assert.Equal(t, expectedRoot, computedRoot)
}

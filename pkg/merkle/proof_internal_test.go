package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidate tests the internal validate function for MultiProof.
func TestValidate(t *testing.T) {
	t.Run("valid proof", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: []byte("test"), index: 0}

		proofs := []Node{
			{hash: []byte("proof")},
		}

		proof := MultiProof{
			values:    values,
			proofs:    proofs,
			leafCount: 2,
		}

		err := proof.validate()
		assert.NoError(t, err)
	})

	t.Run("empty values", func(t *testing.T) {
		proof := MultiProof{
			values:    IndexedValues{},
			proofs:    []Node{{hash: []byte("proof")}},
			leafCount: 2,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNoElements)
	})

	t.Run("invalid leaf count", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: []byte("test"), index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{{hash: []byte("proof")}},
			leafCount: 0,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrInvalidLeafCount)
	})

	t.Run("duplicate indices", func(t *testing.T) {
		values := make(IndexedValues, 2)
		values[0] = IndexedValue{value: []byte("test1"), index: 0}
		values[1] = IndexedValue{value: []byte("test2"), index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{{hash: []byte("proof")}},
			leafCount: 2,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrDuplicateIndices)
	})

	t.Run("nil element", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: nil, index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{{hash: []byte("proof")}},
			leafCount: 2,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNilElement)
	})

	t.Run("nil proof", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: []byte("test"), index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{{hash: nil}},
			leafCount: 2,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNilProof)
	})

	t.Run("no proofs needed for single element tree", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: []byte("test"), index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{},
			leafCount: 1,
		}

		err := proof.validate()
		assert.NoError(t, err)
	})

	t.Run("no proofs when needed", func(t *testing.T) {
		values := make(IndexedValues, 1)
		values[0] = IndexedValue{value: []byte("test"), index: 0}

		proof := MultiProof{
			values:    values,
			proofs:    []Node{},
			leafCount: 2,
		}

		err := proof.validate()
		assert.ErrorIs(t, err, ErrNoProofs)
	})
}

// TestGetNextProof tests the internal getNextProof function.
func TestGetNextProof(t *testing.T) {
	proofs := []Node{
		{hash: []byte{0x01}},
		{hash: []byte{0x02}},
	}

	proof := MultiProof{
		values:    IndexedValues{{value: []byte("test"), index: 0}},
		proofs:    proofs,
		leafCount: 2,
	}

	idx := 0
	p, err := proof.getNextProof(&idx)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01}, p)
	assert.Equal(t, 1, idx)

	// Test getting the second proof.
	p, err = proof.getNextProof(&idx)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x02}, p)
	assert.Equal(t, 2, idx)

	// Test error when out of proofs.
	p, err = proof.getNextProof(&idx)
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

	indexedValues := make(IndexedValues, 2)
	indexedValues[0] = IndexedValue{value: leaves[0], index: 0}
	indexedValues[1] = IndexedValue{value: leaves[2], index: 2}

	proof := MultiProof{
		values:    indexedValues,
		proofs:    []Node{},
		leafCount: len(leaves),
	}

	// Build the node queue with 4 leaves, as it's the balanced leaf count for 4 leaves.
	queue, err := proof.buildNodeQueue(4)
	require.NoError(t, err)
	assert.Equal(t, 2, len(queue))

	// Check queue order and values.
	assert.Equal(t, 6, queue[0].index) // 4+2 (balanced leaf count + index)
	assert.Equal(t, HashLeaf(leaves[2]), queue[0].hash)
	assert.Equal(t, 4, queue[1].index) // 4+0 (balanced leaf count + index)
	assert.Equal(t, HashLeaf(leaves[0]), queue[1].hash)
}

// TestMakeIndexedValues tests the internal makeIndexedValues function.
func TestMakeIndexedValues(t *testing.T) {
	leaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
	}

	indices := []int{0, 2}

	indexedValues, err := makeIndexedValues(leaves, indices)
	require.NoError(t, err)

	assert.Equal(t, 2, len(indexedValues))

	assert.Equal(t, []byte(leaves[0]), []byte(indexedValues[0].value))
	assert.Equal(t, 0, indexedValues[0].index)
	assert.Equal(t, []byte(leaves[2]), []byte(indexedValues[1].value))
	assert.Equal(t, 2, indexedValues[1].index)
}

// TestMakeIndices tests the internal makeIndices function.
func TestMakeIndices(t *testing.T) {
	tests := []struct {
		name          string
		startingIndex int
		count         int
		expected      []int
		wantErr       error
	}{
		{
			name:          "valid indices",
			startingIndex: 2,
			count:         3,
			expected:      []int{2, 3, 4},
			wantErr:       nil,
		},
		{
			name:          "negative starting index",
			startingIndex: -1,
			count:         3,
			expected:      nil,
			wantErr:       ErrInvalidRange,
		},
		{
			name:          "zero count",
			startingIndex: 0,
			count:         0,
			expected:      nil,
			wantErr:       ErrInvalidRange,
		},
		{
			name:          "negative count",
			startingIndex: 0,
			count:         -1,
			expected:      nil,
			wantErr:       ErrInvalidRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indices, err := makeIndices(tt.startingIndex, tt.count)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, indices)
			}
		})
	}
}

// TestValidateIndices tests the internal validateIndices function.
func TestValidateIndices(t *testing.T) {
	tests := []struct {
		name      string
		indices   []int
		leafCount int
		wantErr   error
	}{
		{
			name:      "valid indices",
			indices:   []int{0, 1, 2},
			leafCount: 4,
			wantErr:   nil,
		},
		{
			name:      "empty indices",
			indices:   []int{},
			leafCount: 4,
			wantErr:   ErrNoIndices,
		},
		{
			name:      "duplicate indices",
			indices:   []int{0, 1, 1, 2},
			leafCount: 4,
			wantErr:   ErrDuplicateIndices,
		},
		{
			name:      "out of bounds indices",
			indices:   []int{0, 4, 2},
			leafCount: 4,
			wantErr:   ErrIndicesOutOfBounds,
		},
		{
			name:      "negative index",
			indices:   []int{0, -1, 2},
			leafCount: 4,
			wantErr:   ErrIndicesOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIndices(tt.indices, tt.leafCount)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestHasDuplicates tests the internal hasDuplicates function.
func TestHasDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		indices  []int
		expected bool
	}{
		{
			name:     "no duplicates",
			indices:  []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "has duplicates in middle",
			indices:  []int{1, 2, 2, 3},
			expected: true,
		},
		{
			name:     "has duplicates at start",
			indices:  []int{1, 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasDuplicates(tt.indices)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHasOutOfBounds tests the internal hasOutOfBounds function.
func TestHasOutOfBounds(t *testing.T) {
	tests := []struct {
		name      string
		indices   []int
		leafCount int
		expected  bool
	}{
		{
			name:      "all in bounds",
			indices:   []int{0, 1, 2, 3},
			leafCount: 4,
			expected:  false,
		},
		{
			name:      "upper bound exceeded",
			indices:   []int{0, 1, 4},
			leafCount: 4,
			expected:  true,
		},
		{
			name:      "negative index",
			indices:   []int{-1, 0, 1},
			leafCount: 4,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasOutOfBounds(tt.indices, tt.leafCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsLeftChild tests the internal isLeftChild function.
func TestIsLeftChild(t *testing.T) {
	// Even indices are left children.
	assert.True(t, isLeftChild(0))
	assert.True(t, isLeftChild(2))
	assert.True(t, isLeftChild(4))

	// Odd indices are right children.
	assert.False(t, isLeftChild(1))
	assert.False(t, isLeftChild(3))
	assert.False(t, isLeftChild(5))
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
	indices := []int{0, 2}
	multiProof, err := tree.GenerateMultiProofWithIndices(indices)
	require.NoError(t, err)

	computedRoot, err := multiProof.computeRoot()
	require.NoError(t, err)
	assert.Equal(t, expectedRoot, computedRoot)

	// Test with invalid leaf count.
	badProof := &MultiProof{
		values:    multiProof.values,
		proofs:    multiProof.proofs,
		leafCount: -1,
	}

	_, err = badProof.computeRoot()
	assert.Error(t, err)

	// Test with not enough proofs.
	noProofsProof := &MultiProof{
		values:    multiProof.values,
		proofs:    []Node{},
		leafCount: multiProof.leafCount,
	}

	_, err = noProofsProof.computeRoot()
	assert.Error(t, err)
}

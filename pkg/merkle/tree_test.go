package merkle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/merkle"
)

func TestEmptyTree(t *testing.T) {
	_, err := merkle.NewMerkleTree([]merkle.Leaf{})
	assert.Error(t, err, "Should error on empty elements")
	assert.ErrorAs(t, err, &merkle.ErrTreeEmpty)
}

func TestBalancedTrees(t *testing.T) {
	testCases := []struct {
		name          string
		numElements   int
		expectedDepth int
	}{
		{"TwoElements", 2, 1},
		{"FourElements", 4, 2},
		{"EightElements", 8, 3},
		{"SixteenElements", 16, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create elements.
			elements := make([]merkle.Leaf, tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				elements[i] = []byte(tc.name + "_element" + string(rune('A'+i)))
			}

			tree, err := merkle.NewMerkleTree(elements)
			require.NoError(t, err)

			// Check structure.
			assert.Equal(
				t,
				tc.numElements,
				tree.LeafCount(),
				"Leaf count should match element count",
			)
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			// For balanced trees, the array size is exactly 2*n where n is the next power of 2 >= numElements.
			leafCount, err := merkle.CalculateBalancedNodesCount(tc.numElements)
			require.NoError(t, err)
			expectedArraySize := leafCount * 2
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.Tree()),
				"Tree array size should be 2*leafCount",
			)

			// Check that all leaves are present.
			for i := 0; i < tc.numElements; i++ {
				leafIndex := leafCount + i
				assert.NotNil(t, tree.Tree()[leafIndex], "Leaf node should not be nil")
				assert.Equal(
					t,
					merkle.HashLeaf(elements[i]),
					tree.Tree()[leafIndex].Hash(),
					"Leaf hash should match",
				)
			}

			// Check that all internal nodes up to the root are not nil.
			for i := 1; i < leafCount; i++ {
				assert.NotNil(
					t,
					tree.Tree()[i].Hash(),
					"Internal node should not be nil in a balanced tree",
				)
			}
		})
	}
}

func TestUnbalancedTrees(t *testing.T) {
	testCases := []struct {
		name        string
		numElements int
	}{
		{"SingleElement", 1},
		{"ThreeElements", 3},
		{"FiveElements", 5},
		{"SevenElements", 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elements := make([]merkle.Leaf, tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				elements[i] = []byte(tc.name + "_element" + string(rune('A'+i)))
			}

			tree, err := merkle.NewMerkleTree(elements)
			require.NoError(t, err)

			leafCount, err := merkle.CalculateBalancedNodesCount(tc.numElements)
			require.NoError(t, err)
			for i := 0; i < tc.numElements; i++ {
				leafIdx := leafCount + i
				assert.Equal(t, merkle.HashLeaf(elements[i]), tree.Tree()[leafIdx].Hash(),
					"Leaf %d should be at tree[%d]", i, leafIdx)
			}

			verifyUnbalancedTreeStructure(t, tree.Tree(), leafCount, tc.numElements)
		})
	}
}

func TestLargeTrees(t *testing.T) {
	testCases := []struct {
		name           string
		providedLeaves int
		expectedLeaves int
	}{
		{"TreeSize100", 100, 128},
		{"TreeSize1023", 1023, 1024},
		{"TreeSize2048", 2048, 2048},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elements := make([]merkle.Leaf, tc.providedLeaves)
			for i := 0; i < tc.providedLeaves; i++ {
				elements[i] = []byte{
					byte(i & 0xFF),
					byte((i >> 8) & 0xFF),
					byte((i >> 16) & 0xFF),
					byte((i >> 24) & 0xFF),
				}
			}

			tree, err := merkle.NewMerkleTree(elements)
			require.NoError(t, err)

			// Verify basic properties.
			assert.Equal(
				t,
				tc.expectedLeaves,
				tree.LeafCount(),
				"Leaf count should match element count",
			)
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			// Verify tree structure size.
			leafCount, err := merkle.CalculateBalancedNodesCount(tc.providedLeaves)
			require.NoError(t, err)
			expectedArraySize := leafCount * 2
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.Tree()),
				"Tree array size should be 2*leafCount",
			)

			// Sample testing.
			for i := 0; i < 5; i++ {
				idx := i * (tc.providedLeaves / 5)
				if idx >= tc.providedLeaves {
					idx = tc.providedLeaves - 1
				}
				leafIndex := leafCount + idx
				assert.NotNil(t, tree.Tree()[leafIndex], "Sampled leaf should not be nil")
				assert.Equal(
					t,
					merkle.HashLeaf(elements[idx]),
					tree.Tree()[leafIndex].Hash(),
					"Sampled leaf hash should match",
				)
			}
		})
	}
}

func TestTreeWithDuplicateElements(t *testing.T) {
	elements := []merkle.Leaf{
		[]byte("same"),
		[]byte("same"),
		[]byte("same"),
		[]byte("different"),
	}

	tree, err := merkle.NewMerkleTree(elements)
	require.NoError(t, err)

	leafCount, err := merkle.CalculateBalancedNodesCount(len(elements))
	require.NoError(t, err)
	leafHash1 := tree.Tree()[leafCount].Hash()
	leafHash2 := tree.Tree()[leafCount+1].Hash()
	leafHash3 := tree.Tree()[leafCount+2].Hash()
	leafHash4 := tree.Tree()[leafCount+3].Hash()

	assert.Equal(t, leafHash1, leafHash2, "Identical elements should have identical leaf hashes")
	assert.Equal(t, leafHash2, leafHash3, "Identical elements should have identical leaf hashes")
	assert.NotEqual(t, leafHash3, leafHash4, "Different elements should have different leaf hashes")

	assert.NotNil(t, tree.Root(), "Tree with duplicate elements should have a valid root")
}

func TestTreeWithLargeElements(t *testing.T) {
	bigElement := make([]byte, 1024*1024)
	for i := range bigElement {
		bigElement[i] = byte(i & 0xFF)
	}

	elements := []merkle.Leaf{
		bigElement,
		bigElement[:len(bigElement)/2],
	}

	tree, err := merkle.NewMerkleTree(elements)
	require.NoError(t, err)
	assert.NotNil(
		t,
		tree.Root(),
		"Root should be calculated correctly even with large elements",
	)
}

func TestTreeWithEmptyElements(t *testing.T) {
	elements := []merkle.Leaf{
		{},
		{},
		[]byte("non-empty"),
	}

	tree, err := merkle.NewMerkleTree(elements)
	require.NoError(t, err)
	assert.NotNil(t, tree.Root(), "Root should be calculated correctly with empty elements")

	leafCount, err := merkle.CalculateBalancedNodesCount(len(elements))
	require.NoError(t, err)
	for i := 0; i < 2; i++ {
		leafIndex := leafCount + i
		assert.Equal(
			t,
			merkle.HashLeaf([]byte{}),
			tree.Tree()[leafIndex].Hash(),
			"Empty element should be properly hashed",
		)
	}
}

func TestTreeInternals(t *testing.T) {
	// Test with a 3-element tree (unbalanced)
	// Check everything "manually".
	//
	// Tree structure:
	//        [1]
	//       /   \
	//     [2]    [3]
	//    /  \    /
	//  [4]  [5] [6]
	//  A    B    C

	elements := []merkle.Leaf{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}

	tree, err := merkle.NewMerkleTree(elements)
	require.NoError(t, err)

	internalTree := tree.Tree()

	// For a 3-element tree, the balanced leaf count is 4
	// So the tree array should have size 8 (2*4)
	assert.Equal(t, 8, len(internalTree), "Tree array should have size 2*leafCount")
	assert.Equal(t, tree.Root(), internalTree[1].Hash(), "Root should be at index 1")

	// Check that all leaves are present.
	leafStartIdx := 4
	assert.Equal(
		t,
		merkle.HashLeaf(elements[0]),
		internalTree[leafStartIdx].Hash(),
		"Leaf 0 should match hash",
	)
	assert.Equal(
		t,
		merkle.HashLeaf(elements[1]),
		internalTree[leafStartIdx+1].Hash(),
		"Leaf 1 should match hash",
	)
	assert.Equal(
		t,
		merkle.HashLeaf(elements[2]),
		internalTree[leafStartIdx+2].Hash(),
		"Leaf 2 should match hash",
	)
	assert.Equal(
		t,
		internalTree[leafStartIdx+3].Hash(),
		merkle.HashLeaf([]byte{}),
		"Leaf 3 should match empty element hash",
	)

	verifyUnbalancedTreeStructure(t, internalTree, 4, len(elements))
}

// Helper function to verify the structure of an unbalanced tree
func verifyUnbalancedTreeStructure(
	t *testing.T,
	tree []merkle.Node,
	leafCount, actualElements int,
) {
	t.Helper()

	// Handle single element trees.
	if actualElements == 1 {
		assert.Nil(t, tree[0].Hash(), "Merkle Tree 1-indexed, index 0 should be nil")
		assert.Equal(
			t,
			tree[leafCount].Hash(),
			tree[1].Hash(),
			"For single element tree, root should equal leaf",
		)
		return
	}

	// Starting from the leaf parents, check nodes up to the root.
	for i := leafCount - 1; i > 0; i-- {
		leftIdx := merkle.GetLeftChild(i)
		rightIdx := merkle.GetRightChild(i)

		if leftIdx >= len(tree) {
			continue
		}

		// If both children exist
		if leftIdx < len(tree) && rightIdx < len(tree) {
			assert.NotNil(t, tree[i], "Parent node should not be nil when both children exist")
			assert.Equal(
				t,
				merkle.HashNode(tree[leftIdx].Hash(), tree[rightIdx].Hash()),
				tree[i].Hash(),
				"Node at index %d should be equal to its children at %d and %d",
				i,
				leftIdx,
				rightIdx,
			)
		}
	}

	// Check the root is valid.
	expectedRoot := tree[1].Hash()
	assert.NotNil(t, expectedRoot, "Root should not be nil")
	assert.NotNil(
		t,
		merkle.HashNode(tree[1].Hash(), tree[2].Hash()),
		"Root should be calculated when both children exist",
	)
}

func TestRoundUpToPowerOf2Values(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
		wantErr  bool
	}{
		{"zero", 0, 0, true},
		{"one", 1, 1, false},
		{"already power of 2", 4, 4, false},
		{"already power of 2 (large)", 16384, 16384, false},
		{"regular case", 5, 8, false},
		{"regular case (large)", 5000, 8192, false},
		{"large number", 1<<30 - 1, 1 << 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := merkle.CalculateBalancedNodesCount(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if result != tt.expected {
					t.Errorf(
						"Power of 2 rounding for %d = %d, expected %d",
						tt.input,
						result,
						tt.expected,
					)
				}
			}
		})
	}
}

func TestCalculateBalancedLeafCount(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
		wantErr  bool
	}{
		{"negative", -1, 0, true},
		{"zero", 0, 0, true},
		{"one", 1, 1, false},
		{"power of 2", 16, 16, false},
		{"not power of 2", 15, 16, false},
		{"large number", 1000000, 1048576, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := merkle.CalculateBalancedNodesCount(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if result != tt.expected {
					t.Errorf("CalculateBalancedLeafCount(%d) = %d, expected %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestCalculateBalancedLeafCountError tests that the function returns an error with large inputs
func TestCalculateBalancedLeafCountError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode - requires large values")
	}

	// This is larger than max uint32 and should cause an error
	massiveInput := int(^uint32(0)) + 1
	_, err := merkle.CalculateBalancedNodesCount(massiveInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "count must be less than or equal than max int32")
}

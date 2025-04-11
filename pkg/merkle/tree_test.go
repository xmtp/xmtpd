package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTree(t *testing.T) {
	_, err := NewMerkleTree([][]byte{})
	assert.Error(t, err, "Should error on empty elements")
	assert.Contains(
		t,
		err.Error(),
		"elements cannot be empty",
	)
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
			elements := make([][]byte, tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				elements[i] = []byte(tc.name + "_element" + string(rune('A'+i)))
			}

			tree, err := NewMerkleTree(elements)
			require.NoError(t, err)

			// Check structure.
			assert.Equal(
				t,
				tc.expectedDepth,
				tree.Depth(),
				"Tree depth should match expected value",
			)
			assert.Equal(t, tc.numElements, tree.leafCount, "Leaf count should match element count")
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			// For balanced trees, the array size is exactly 2*n where n is the next power of 2 >= numElements.
			leafCount := getLeafCount(tc.numElements)
			expectedArraySize := leafCount * 2
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.tree),
				"Tree array size should be 2*leafCount",
			)

			// Check that all leaves are present.
			for i := 0; i < tc.numElements; i++ {
				leafIndex := leafCount + i
				assert.NotNil(t, tree.tree[leafIndex], "Leaf node should not be nil")
				assert.Equal(
					t,
					HashLeaf(elements[i]),
					tree.tree[leafIndex],
					"Leaf hash should match",
				)
			}

			// Check that all internal nodes up to the root are not nil.
			for i := 1; i < leafCount; i++ {
				assert.NotNil(t, tree.tree[i], "Internal node should not be nil in a balanced tree")
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
			elements := make([][]byte, tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				elements[i] = []byte(tc.name + "_element" + string(rune('A'+i)))
			}

			tree, err := NewMerkleTree(elements)
			require.NoError(t, err)

			leafCount := getLeafCount(tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				leafIdx := leafCount + i
				assert.Equal(t, HashLeaf(elements[i]), tree.tree[leafIdx],
					"Leaf %d should be at tree[%d]", i, leafIdx)
			}

			verifyUnbalancedTreeStructure(t, tree.tree, leafCount, tc.numElements)
		})
	}
}

func TestLargeTrees(t *testing.T) {
	testCases := []struct {
		name        string
		numElements int
	}{
		{"TreeSize100", 100},
		{"TreeSize1023", 1023},
		{"TreeSize2048", 2048},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elements := make([][]byte, tc.numElements)
			for i := 0; i < tc.numElements; i++ {
				elements[i] = []byte{
					byte(i & 0xFF),
					byte((i >> 8) & 0xFF),
					byte((i >> 16) & 0xFF),
					byte((i >> 24) & 0xFF),
				}
			}

			tree, err := NewMerkleTree(elements)
			require.NoError(t, err)

			// Verify basic properties.
			assert.Equal(t, tc.numElements, tree.leafCount, "Leaf count should match element count")
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			// Verify tree structure size.
			leafCount := getLeafCount(tc.numElements)
			expectedArraySize := leafCount * 2
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.tree),
				"Tree array size should be 2*leafCount",
			)

			// Sample testing.
			for i := 0; i < 5; i++ {
				idx := i * (tc.numElements / 5)
				if idx >= tc.numElements {
					idx = tc.numElements - 1
				}
				leafIndex := leafCount + idx
				assert.NotNil(t, tree.tree[leafIndex], "Sampled leaf should not be nil")
				assert.Equal(
					t,
					HashLeaf(elements[idx]),
					tree.tree[leafIndex],
					"Sampled leaf hash should match",
				)
			}
		})
	}
}

func TestTreeWithDuplicateElements(t *testing.T) {
	elements := [][]byte{
		[]byte("same"),
		[]byte("same"),
		[]byte("same"),
		[]byte("different"),
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)

	leafCount := getLeafCount(len(elements))
	leafHash1 := tree.tree[leafCount]
	leafHash2 := tree.tree[leafCount+1]
	leafHash3 := tree.tree[leafCount+2]
	leafHash4 := tree.tree[leafCount+3]

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

	elements := [][]byte{
		bigElement,
		bigElement[:len(bigElement)/2],
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)
	assert.NotNil(
		t,
		tree.Root(),
		"Root should be calculated correctly even with large elements",
	)
}

func TestTreeWithEmptyElements(t *testing.T) {
	elements := [][]byte{
		{},
		{},
		[]byte("non-empty"),
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)
	assert.NotNil(t, tree.Root(), "Root should be calculated correctly with empty elements")

	leafCount := getLeafCount(len(elements))
	for i := 0; i < 2; i++ {
		leafIndex := leafCount + i
		assert.Equal(
			t,
			HashLeaf([]byte{}),
			tree.tree[leafIndex],
			"Empty element should be properly hashed",
		)
	}
}

func TestTreeInternals(t *testing.T) {
	// Test with a 3-element tree (unbalanced)
	// Check everything "manually".
	//
	// Tree structure:
	//        [0]
	//       /   \
	//     [1]    [2]
	//    /  \    /
	//  [3]  [4] [5]
	//  A    B    C

	elements := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}

	tree, err := NewMerkleTree(elements)
	require.NoError(t, err)

	internalTree := tree.tree

	// For a 3-element tree, the balanced leaf count is 4
	// So the tree array should have size 8 (2*4)
	assert.Equal(t, 8, len(internalTree), "Tree array should have size 2*leafCount")

	// Check that all nodes are present.
	assert.NotNil(t, internalTree[0], "Root should not be nil")
	assert.NotNil(t, internalTree[1], "Node 1 should not be nil")
	assert.NotNil(t, internalTree[2], "Node 2 should not be nil")
	assert.NotNil(t, internalTree[3], "Node 3 should not be nil")
	assert.NotNil(t, internalTree[4], "Node 4 should not be nil")
	assert.NotNil(t, internalTree[5], "Node 5 should not be nil")

	// Check that all leaves are present.
	leafStartIdx := 4
	assert.Equal(t, HashLeaf(elements[0]), internalTree[leafStartIdx], "Leaf 0 should match")
	assert.Equal(t, HashLeaf(elements[1]), internalTree[leafStartIdx+1], "Leaf 1 should match")
	assert.Equal(t, HashLeaf(elements[2]), internalTree[leafStartIdx+2], "Leaf 2 should match")
	assert.Nil(t, internalTree[leafStartIdx+3], "Leaf 3 should be nil (not in original elements)")

	verifyUnbalancedTreeStructure(t, internalTree, 4, len(elements))
}

// Helper function to verify the structure of an unbalanced tree
func verifyUnbalancedTreeStructure(t *testing.T, tree [][]byte, leafCount, actualElements int) {
	t.Helper()

	// Handle single element trees.
	if actualElements == 1 {
		assert.Equal(t, tree[leafCount], tree[0], "For single element tree, root should equal leaf")
		return
	}

	// Starting from the leaf parents, check nodes up to the root.
	for i := leafCount - 1; i > 0; i-- {
		leftIdx := getLeftChild(i)
		rightIdx := getRightChild(i)

		if leftIdx >= len(tree) {
			continue
		}

		// If both children exist
		if leftIdx < len(tree) && rightIdx < len(tree) &&
			tree[leftIdx] != nil && tree[rightIdx] != nil {
			assert.NotNil(t, tree[i], "Parent node should not be nil when both children exist")
			assert.Equal(
				t,
				HashNode(tree[leftIdx], tree[rightIdx]),
				tree[i],
				"Node at index %d should be equal to its children at %d and %d",
				i,
				leftIdx,
				rightIdx,
			)
		} else if leftIdx < len(tree) && tree[leftIdx] != nil {
			assert.Equal(t, tree[leftIdx], tree[i],
				"Node at index %d should be equal to its only child at %d", i, leftIdx)
		}
	}

	// Check the root is valid.
	assert.NotNil(t, tree[0], "Root should not be nil")
	if tree[1] != nil && tree[2] != nil {
		assert.NotNil(t, tree[0], "Root should be calculated when both children exist")
	} else {
		assert.Equal(t, tree[1], tree[0], "Root should equal its only child (node 1)")
	}
}

package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundUpToPowerOf2(t *testing.T) {
	testCases := []struct {
		input    uint32
		expected uint32
	}{
		{0, 1},       // Special case: 0 rounds up to 1
		{1, 1},       // Already a power of 2
		{2, 2},       // Already a power of 2
		{3, 4},       // Rounds up to 4
		{4, 4},       // Already a power of 2
		{5, 8},       // Rounds up to 8
		{7, 8},       // Rounds up to 8
		{8, 8},       // Already a power of 2
		{9, 16},      // Rounds up to 16
		{15, 16},     // Rounds up to 16
		{16, 16},     // Already a power of 2
		{31, 32},     // Rounds up to 32
		{32, 32},     // Already a power of 2
		{33, 64},     // Rounds up to 64
		{63, 64},     // Rounds up to 64
		{64, 64},     // Already a power of 2
		{65, 128},    // Rounds up to 128
		{127, 128},   // Rounds up to 128
		{128, 128},   // Already a power of 2
		{1023, 1024}, // Rounds up to 1024
		{1024, 1024}, // Already a power of 2
		{1025, 2048}, // Rounds up to 2048
	}

	for _, tc := range testCases {
		result := RoundUpToPowerOf2(tc.input)
		assert.Equal(t, tc.expected, result, "RoundUpToPowerOf2(%d) should be %d, got %d",
			tc.input, tc.expected, result)
	}
}

func TestGetBalancedLeafCount(t *testing.T) {
	testCases := []struct {
		elements int
		expected int
	}{
		{1, 1},   // Single element
		{2, 2},   // Already balanced
		{3, 4},   // Unbalanced -> rounds up to 4
		{4, 4},   // Already balanced
		{5, 8},   // Unbalanced -> rounds up to 8
		{7, 8},   // Unbalanced -> rounds up to 8
		{8, 8},   // Already balanced
		{9, 16},  // Unbalanced -> rounds up to 16
		{15, 16}, // Unbalanced -> rounds up to 16
		{16, 16}, // Already balanced
	}

	for _, tc := range testCases {
		result := getLeafCount(tc.elements)
		assert.Equal(t, tc.expected, result, "getBalancedLeafCount(%d) should be %d, got %d",
			tc.elements, tc.expected, result)
	}
}

func TestTreeCapacityVsActualElements(t *testing.T) {
	// Verify that trees with different element counts use the right capacity
	testCases := []struct {
		elements       int
		expectedDepth  int
		expectedLeaves int
	}{
		{1, 0, 1},  // Depth 0 (just root), 1 leaf
		{2, 1, 2},  // Depth 1, 2 leaves
		{3, 2, 4},  // Depth 2, capacity for 4 leaves
		{4, 2, 4},  // Depth 2, 4 leaves
		{5, 3, 8},  // Depth 3, capacity for 8 leaves
		{7, 3, 8},  // Depth 3, capacity for 8 leaves
		{8, 3, 8},  // Depth 3, 8 leaves
		{9, 4, 16}, // Depth 4, capacity for 16 leaves
	}

	for _, tc := range testCases {
		t.Run("Elements="+string(rune(tc.elements+'0')), func(t *testing.T) {
			// Create elements
			elements := make([][]byte, tc.elements)
			for i := range elements {
				elements[i] = []byte{byte(i + 1)}
			}

			// Create tree
			tree, err := NewMerkleTree(elements)
			assert.NoError(t, err)

			// Verify depth
			assert.Equal(t, tc.expectedDepth, tree.Depth(),
				"Tree with %d elements should have depth %d", tc.elements, tc.expectedDepth)

			// Verify internal structure
			actualLeaves := len(tree.tree) >> 1
			assert.Equal(
				t,
				tc.expectedLeaves,
				actualLeaves,
				"Tree with %d elements should have capacity for %d leaves",
				tc.elements,
				tc.expectedLeaves,
			)

			// Verify root is not nil
			assert.NotNil(t, tree.Root(), "Root should not be nil for %d elements", tc.elements)
		})
	}
}

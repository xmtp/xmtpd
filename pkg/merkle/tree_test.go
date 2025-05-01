package merkle_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/merkle"
)

func TestEmptyTree(t *testing.T) {
	_, err := merkle.NewMerkleTree([]merkle.Leaf{})
	assert.Error(t, err, "Should error on empty leaves")
	assert.ErrorAs(t, err, &merkle.ErrTreeEmpty)
}

func TestBalancedTrees(t *testing.T) {
	testCases := []struct {
		name      string
		leafCount int
	}{
		{"TwoLeaves", 2},
		{"FourLeaves", 4},
		{"EightLeaves", 8},
		{"SixteenLeaves", 16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create leaves.
			leaves := make([]merkle.Leaf, tc.leafCount)
			for i := 0; i < tc.leafCount; i++ {
				leaves[i] = []byte(tc.name + "_leaf" + string(rune('A'+i)))
			}

			tree, err := merkle.NewMerkleTree(leaves)
			require.NoError(t, err)

			// Check structure.
			assert.Equal(
				t,
				tc.leafCount,
				tree.LeafCount(),
				"Leaf count should match",
			)
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			expectedArraySize := tc.leafCount * 2
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.Tree()),
				"Tree array size should be 2*tc.leafCount",
			)

			// Check that all leaves are present.
			for i := 0; i < tc.leafCount; i++ {
				leafIndex := tc.leafCount + i
				assert.NotNil(t, tree.Tree()[leafIndex], "Leaf node should not be nil")
				assert.Equal(
					t,
					merkle.Node(merkle.HashLeaf(leaves[i])),
					tree.Tree()[leafIndex],
					"Leaf hash should match",
				)
			}

			// Check that all internal nodes up to the root are not nil.
			for i := 1; i < tc.leafCount; i++ {
				assert.NotNil(
					t,
					tree.Tree()[i],
					"Internal node should not be nil in a balanced tree",
				)
			}
		})
	}
}

func TestUnbalancedTrees(t *testing.T) {
	testCases := []struct {
		name      string
		leafCount int
	}{
		{"SingleLeaf", 1},
		{"ThreeLeaves", 3},
		{"FiveLeaves", 6},
		{"SevenLeaves", 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create leaves.
			leaves := make([]merkle.Leaf, tc.leafCount)
			for i := 0; i < tc.leafCount; i++ {
				leaves[i] = []byte(tc.name + "_leaf" + string(rune('A'+i)))
			}

			tree, err := merkle.NewMerkleTree(leaves)
			require.NoError(t, err)

			// Check structure.
			assert.Equal(
				t,
				tc.leafCount,
				tree.LeafCount(),
				"Leaf count should match",
			)
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			balancedLeafCount, err := merkle.CalculateBalancedNodesCount(tc.leafCount)
			require.NoError(t, err)

			expectedArraySize := balancedLeafCount << 1
			assert.Equal(
				t,
				expectedArraySize,
				len(tree.Tree()),
				"Tree array size should be balancedLeafCount + tc.leafCount",
			)

			// Check that all leaves are present.
			for i := 0; i < tc.leafCount; i++ {
				leafIndex := balancedLeafCount + i
				assert.NotNil(t, tree.Tree()[leafIndex], "Leaf node should not be nil")
				assert.Equal(
					t,
					merkle.Node(merkle.HashLeaf(leaves[i])),
					tree.Tree()[leafIndex],
					"Leaf hash should match",
				)
			}

			// Check that all internal nodes, let of the upperBound, up to the root are not nil.
			// TODO
		})
	}
}

func TestLargeTrees(t *testing.T) {
	testCases := []struct {
		name             string
		leafCount        int
		expectedTreeSize int
	}{
		{"TreeSize100", 100, 256},
		{"TreeSize1023", 1023, 2048},
		{"TreeSize2048", 2048, 4096},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			leaves := make([]merkle.Leaf, tc.leafCount)
			for i := 0; i < tc.leafCount; i++ {
				leaves[i] = []byte{
					byte(i & 0xFF),
					byte((i >> 8) & 0xFF),
					byte((i >> 16) & 0xFF),
					byte((i >> 24) & 0xFF),
				}
			}

			tree, err := merkle.NewMerkleTree(leaves)
			require.NoError(t, err)

			// Verify basic properties.
			assert.Equal(
				t,
				tc.leafCount,
				tree.LeafCount(),
				"Leaf count should match",
			)
			assert.NotNil(t, tree.Root(), "Root should not be nil")

			// Verify tree structure size.
			balancedLeafCount, err := merkle.CalculateBalancedNodesCount(tc.leafCount)
			require.NoError(t, err)
			assert.Equal(
				t,
				tc.expectedTreeSize,
				len(tree.Tree()),
				"Tree array size should be tc.expectedTreeSize",
			)

			// Sample testing.
			for i := 0; i < 5; i++ {
				idx := i * (tc.leafCount / 5)
				if idx >= tc.leafCount {
					idx = tc.leafCount - 1
				}
				leafIndex := balancedLeafCount + idx
				assert.NotNil(t, tree.Tree()[leafIndex], "Sampled leaf should not be nil")
				assert.Equal(
					t,
					merkle.Node(merkle.HashLeaf(leaves[idx])),
					tree.Tree()[leafIndex],
					"Sampled leaf hash should match",
				)
			}
		})
	}
}

func TestTreeWithDuplicateLeaves(t *testing.T) {
	leaves := []merkle.Leaf{
		[]byte("same"),
		[]byte("same"),
		[]byte("same"),
		[]byte("different"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	leafCount, err := merkle.CalculateBalancedNodesCount(len(leaves))
	require.NoError(t, err)
	leafHash1 := tree.Tree()[leafCount]
	leafHash2 := tree.Tree()[leafCount+1]
	leafHash3 := tree.Tree()[leafCount+2]
	leafHash4 := tree.Tree()[leafCount+3]

	assert.True(
		t,
		bytes.Equal(leafHash1, leafHash2),
		"Identical leaves should have identical leaf hashes",
	)
	assert.True(
		t,
		bytes.Equal(leafHash2, leafHash3),
		"Identical leaves should have identical leaf hashes",
	)
	assert.False(
		t,
		bytes.Equal(leafHash3, leafHash4),
		"Different leaves should have different leaf hashes",
	)

	assert.NotNil(t, tree.Root(), "Tree with duplicate leaves should have a valid root")
}

func TestTreeWithLargeLeaves(t *testing.T) {
	bigLeaf := make([]byte, 1024*1024)
	for i := range bigLeaf {
		bigLeaf[i] = byte(i & 0xFF)
	}

	leaves := []merkle.Leaf{
		bigLeaf,
		bigLeaf[:len(bigLeaf)/2],
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)
	assert.NotNil(
		t,
		tree.Root(),
		"Root should be calculated correctly even with large leaves",
	)
}

func TestTreeWithEmptyLeaves(t *testing.T) {
	leaves := []merkle.Leaf{
		{},
		{},
		[]byte("non-empty"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)
	assert.NotNil(t, tree.Root(), "Root should be calculated correctly with empty leaves")

	leafCount, err := merkle.CalculateBalancedNodesCount(len(leaves))
	require.NoError(t, err)
	for i := 0; i < 2; i++ {
		leafIndex := leafCount + i
		assert.Equal(
			t,
			merkle.Node(merkle.HashLeaf([]byte{})),
			tree.Tree()[leafIndex],
			"Empty leaves should be properly hashed",
		)
	}
}

func TestTreeWithNilLeaves(t *testing.T) {
	leaves := []merkle.Leaf{
		[]byte("non-empty"),
		nil,
		[]byte("non-empty"),
	}

	_, err := merkle.NewMerkleTree(leaves)
	assert.Error(t, err, "Should error on nil leaves")
	assert.ErrorAs(t, err, &merkle.ErrNilLeaf)
}

func TestTreeInternals(t *testing.T) {
	// Test with a 3-leaf tree (unbalanced)
	// Check everything "manually".
	//
	// Tree structure:
	//        [1]
	//       /   \
	//     [2]    [3]
	//    /  \    /
	//  [4]  [5] [6]
	//  A    B    C

	leaves := []merkle.Leaf{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	internalTree := tree.Tree()

	// For a 3-leaf tree, the balanced leaf count is 4
	// So the tree array should have size 8 (2*4)
	assert.Equal(
		t,
		8,
		len(internalTree),
		"Tree array should have size balancedLeafCount + leafCount",
	)
	assert.True(t, bytes.Equal(tree.Root(), internalTree[0]), "Root should be at index 0")

	// Check that all leaf nodes are present.
	assert.Equal(
		t,
		merkle.Node(merkle.HashLeaf(leaves[0])),
		internalTree[4],
		"Node 4 should match hash",
	)
	assert.Equal(
		t,
		merkle.Node(merkle.HashLeaf(leaves[1])),
		internalTree[5],
		"Node 5 should match hash",
	)
	assert.Equal(
		t,
		merkle.Node(merkle.HashLeaf(leaves[2])),
		internalTree[6],
		"Node 6 should match hash",
	)
	assert.Equal(
		t,
		merkle.Node(nil),
		internalTree[7],
		"Node 7 should be nil",
	)

	// Check that all nodes are present.
	assert.Equal(
		t,
		merkle.Node(merkle.HashNodePair(merkle.HashLeaf(leaves[0]), merkle.HashLeaf(leaves[1]))),
		internalTree[2],
		"Node 2 should match hash",
	)
	assert.Equal(
		t,
		merkle.Node(merkle.HashPairlessNode(merkle.HashLeaf(leaves[2]))),
		internalTree[3],
		"Node 2 should match hash",
	)

	assert.Equal(
		t,
		merkle.Node(merkle.HashNodePair(
			merkle.HashNodePair(merkle.HashLeaf(leaves[0]), merkle.HashLeaf(leaves[1])),
			merkle.HashPairlessNode(merkle.HashLeaf(leaves[2])),
		)),
		internalTree[1],
		"Node 1 should match hash",
	)

	assert.Equal(
		t,
		merkle.Node(merkle.HashRoot(
			3,
			merkle.HashNodePair(
				merkle.HashNodePair(merkle.HashLeaf(leaves[0]), merkle.HashLeaf(leaves[1])),
				merkle.HashPairlessNode(merkle.HashLeaf(leaves[2])),
			),
		)),
		internalTree[0],
		"Node 0 should match hash",
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

func TestBalancedSample(t *testing.T) {
	leaves := []merkle.Leaf{
		getBytesFromHexString("6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2"),
		getBytesFromHexString("a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31"),
		getBytesFromHexString("007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599"),
		getBytesFromHexString("7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828"),
		getBytesFromHexString("4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69"),
		getBytesFromHexString("51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31"),
		getBytesFromHexString("aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c"),
		getBytesFromHexString("e83870d75c6c4c4d1f6ba674481932301e0a1029b44c1407b6aea06cd56d4836"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	assert.Equal(
		t,
		8,
		tree.LeafCount(),
		"LeafCount should be 8",
	)

	assert.Equal(
		t,
		merkle.Node(
			getBytesFromHexString(
				"1f70e7dd11a042e3868e8b0992118a3d7bd301b029a3b967a5b2042466c5110c",
			),
		),
		tree.Tree()[1],
		"Sub Root should be asa expected",
	)

	assert.Equal(
		t,
		getBytesFromHexString("00f8c0ad3c60c727ededce5717c8baa64047b5c3f29e409085df14dc3bfda1a7"),
		merkle.HashRoot(8, tree.Tree()[1]),
		"Root should be asa expected",
	)

	assert.Equal(
		t,
		getBytesFromHexString("00f8c0ad3c60c727ededce5717c8baa64047b5c3f29e409085df14dc3bfda1a7"),
		tree.Root(),
		"Root should be asa expected",
	)
}

func TestUnbalancedSample(t *testing.T) {
	leaves := []merkle.Leaf{
		getBytesFromHexString("6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2"),
		getBytesFromHexString("a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31"),
		getBytesFromHexString("007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599"),
		getBytesFromHexString("7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828"),
		getBytesFromHexString("4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69"),
		getBytesFromHexString("51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31"),
		getBytesFromHexString("aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	assert.Equal(
		t,
		7,
		tree.LeafCount(),
		"LeafCount should be 7",
	)

	assert.Equal(
		t,
		merkle.Node(
			getBytesFromHexString(
				"a9a18d92fa458bf5d28a44d6c0fb4baaf5b4da5918ab7819d5a7d29d8b103205",
			),
		),
		tree.Tree()[1],
		"Sub Root should be asa expected",
	)

	assert.Equal(
		t,
		getBytesFromHexString("38631dd8b5081555ec3c51cc8db7918ee90158fa33a70674c1399234d23908b2"),
		merkle.HashRoot(7, tree.Tree()[1]),
		"Root should be asa expected",
	)

	assert.Equal(
		t,
		getBytesFromHexString("38631dd8b5081555ec3c51cc8db7918ee90158fa33a70674c1399234d23908b2"),
		tree.Root(),
		"Root should be asa expected",
	)
}

func getBytesFromHexString(s string) []byte {
	decoded, _ := hex.DecodeString(s)
	return decoded
}

package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakeLeafNodes tests the internal makeLeafNodes function.
func TestMakeLeafNodes(t *testing.T) {
	// Test with valid leaves.
	leaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
	}

	nodes, err := makeLeafNodes(leaves)
	require.NoError(t, err)
	assert.Equal(t, 3, len(nodes))

	for i, leaf := range leaves {
		assert.Equal(t, Node(HashLeaf(leaf)), nodes[i])
	}

	// Test with empty leaves.
	_, err = makeLeafNodes([]Leaf{})
	assert.ErrorIs(t, err, ErrNoLeaves)
}

// TestMakeTree tests the internal makeTree function.
func TestMakeTree(t *testing.T) {
	// Test with empty nodes.
	_, err := makeTree([]Node{})
	assert.ErrorIs(t, err, ErrTreeEmpty)

	// Test single node tree.
	singleLeaf := []Leaf{[]byte("single")}
	nodes, err := makeLeafNodes(singleLeaf)
	require.NoError(t, err)

	tree, err := makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 2, len(tree))
	assert.Equal(t, nodes[0], tree[1])

	// Test small balanced tree (2 nodes).
	twoLeaves := []Leaf{[]byte("leaf1"), []byte("leaf2")}
	nodes, err = makeLeafNodes(twoLeaves)
	require.NoError(t, err)

	tree, err = makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 4, len(tree))

	// Root should be parent of two leaves.
	expectedRoot := HashNodePair(nodes[0], nodes[1])
	assert.Equal(t, Node(expectedRoot), tree[1])

	// Leaves should be preserved.
	assert.Equal(t, nodes[0], tree[2])
	assert.Equal(t, nodes[1], tree[3])

	// Test unbalanced tree (3 nodes).
	threeLeaves := []Leaf{[]byte("leaf1"), []byte("leaf2"), []byte("leaf3")}
	nodes, err = makeLeafNodes(threeLeaves)
	require.NoError(t, err)

	tree, err = makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 8, len(tree))

	// Check internal nodes.
	leftSubtreeRoot := HashNodePair(nodes[0], nodes[1])
	rightSubtreeRoot := HashPairlessNode(nodes[2])
	expectedRoot = HashNodePair(leftSubtreeRoot, rightSubtreeRoot)
	assert.Equal(t, Node(expectedRoot), tree[1])
	assert.Equal(t, Node(leftSubtreeRoot), tree[2])
	assert.Equal(t, Node(rightSubtreeRoot), tree[3])
}

// TestLeaves tests the public Leaves method of MerkleTree.
func TestLeaves(t *testing.T) {
	// Test with normal leaves.
	originalLeaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
	}

	tree, err := NewMerkleTree(originalLeaves)
	require.NoError(t, err)

	// Test padding.
	leaves := tree.Leaves()
	assert.Equal(t, 3, len(leaves))

	// Verify original leaves are preserved.
	for i, leaf := range originalLeaves {
		assert.Equal(t, leaf, leaves[i])
	}

	// Test with exactly power of 2 leaves.
	powTwoLeaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree, err = NewMerkleTree(powTwoLeaves)
	require.NoError(t, err)

	leaves = tree.Leaves()
	assert.Equal(t, 4, len(leaves)) // No padding needed.
	assert.Equal(t, powTwoLeaves, leaves)
}

// TestNewMerkleTreeErrors tests more error cases for NewMerkleTree.
func TestNewMerkleTreeErrors(t *testing.T) {
	// Test with empty leaves.
	_, err := NewMerkleTree([]Leaf{})
	assert.ErrorIs(t, err, ErrNoLeaves)

	// Test with a successful case as well.
	_, err = NewMerkleTree([]Leaf{[]byte("test")})
	assert.NoError(t, err)
}

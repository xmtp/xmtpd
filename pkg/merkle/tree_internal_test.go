package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakeNodes tests the internal makeNodes function.
func TestMakeNodes(t *testing.T) {
	// Test with valid leaves.
	leaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
	}

	nodes, err := makeNodes(leaves)
	require.NoError(t, err)
	assert.Equal(t, 3, len(nodes))

	for i, leaf := range leaves {
		assert.Equal(t, HashLeaf(leaf), nodes[i].Hash())
	}

	// Test with empty leaves.
	_, err = makeNodes([]Leaf{})
	assert.ErrorIs(t, err, ErrNoLeaves)
}

// TestMakeLeaves tests the internal makeLeaves function.
func TestMakeLeaves(t *testing.T) {
	// Test normal case (3 leaves becomes 4 balanced leaves).
	originalLeaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
	}

	balancedLeaves, err := makeLeaves(originalLeaves)
	require.NoError(t, err)
	assert.Equal(t, 4, len(balancedLeaves))

	// Check original leaves are preserved.
	for i, leaf := range originalLeaves {
		assert.Equal(t, leaf, balancedLeaves[i])
	}

	// Check padding is with empty leaves.
	assert.Equal(t, Leaf{}, balancedLeaves[3])

	// Test with power of 2 leaves (4 leaves stays 4 leaves)
	powTwoLeaves := []Leaf{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	balancedLeaves, err = makeLeaves(powTwoLeaves)
	require.NoError(t, err)
	assert.Equal(t, 4, len(balancedLeaves))
	assert.Equal(t, powTwoLeaves, balancedLeaves)
}

// TestMakeTree tests the internal makeTree function.
func TestMakeTree(t *testing.T) {
	// Test with empty nodes.
	_, err := makeTree([]Node{})
	assert.ErrorIs(t, err, ErrTreeEmpty)

	// Test single node tree.
	singleLeaf := []Leaf{[]byte("single")}
	nodes, err := makeNodes(singleLeaf)
	require.NoError(t, err)

	tree, err := makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 2, len(tree))
	assert.Equal(t, nodes[0].Hash(), tree[1].Hash())

	// Test small balanced tree (2 nodes).
	twoLeaves := []Leaf{[]byte("leaf1"), []byte("leaf2")}
	nodes, err = makeNodes(twoLeaves)
	require.NoError(t, err)

	tree, err = makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 4, len(tree))

	// Root should be parent of two leaves.
	expectedRoot := HashNode(nodes[0].Hash(), nodes[1].Hash())
	assert.Equal(t, expectedRoot, tree[1].Hash())

	// Leaves should be preserved.
	assert.Equal(t, nodes[0].Hash(), tree[2].Hash())
	assert.Equal(t, nodes[1].Hash(), tree[3].Hash())

	// Test unbalanced tree (3 nodes).
	threeLeaves := []Leaf{[]byte("leaf1"), []byte("leaf2"), []byte("leaf3")}
	balancedLeaves, err := makeLeaves(threeLeaves)
	require.NoError(t, err)
	nodes, err = makeNodes(balancedLeaves)
	require.NoError(t, err)

	tree, err = makeTree(nodes)
	require.NoError(t, err)
	assert.Equal(t, 8, len(tree))

	// Check internal nodes.
	leftSubtreeRoot := HashNode(nodes[0].Hash(), nodes[1].Hash())
	rightSubtreeRoot := HashNode(nodes[2].Hash(), nodes[3].Hash())
	expectedRoot = HashNode(leftSubtreeRoot, rightSubtreeRoot)
	assert.Equal(t, expectedRoot, tree[1].Hash())
	assert.Equal(t, leftSubtreeRoot, tree[2].Hash())
	assert.Equal(t, rightSubtreeRoot, tree[3].Hash())
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
	assert.Equal(t, 4, len(leaves))

	// Verify original leaves are preserved.
	for i, leaf := range originalLeaves {
		assert.Equal(t, leaf, leaves[i])
	}

	assert.Equal(t, Leaf{}, leaves[3])

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

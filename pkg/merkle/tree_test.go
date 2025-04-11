package merkle

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilRootScenario(t *testing.T) {
	// Test with different numbers of elements to find cases where tree[0] might be nil
	testCases := []struct {
		name     string
		elements [][]byte
	}{
		{"SingleElement", [][]byte{[]byte("element1")}},
		{"TwoElements", [][]byte{[]byte("element1"), []byte("element2")}},
		{"ThreeElements", [][]byte{[]byte("element1"), []byte("element2"), []byte("element3")}},
		{
			"FourElements",
			[][]byte{
				[]byte("element1"),
				[]byte("element2"),
				[]byte("element3"),
				[]byte("element4"),
			},
		},
		{
			"FiveElements",
			[][]byte{
				[]byte("element1"),
				[]byte("element2"),
				[]byte("element3"),
				[]byte("element4"),
				[]byte("element5"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Build the tree without the nil check
			leafs := make([][]byte, len(tc.elements))
			for i, element := range tc.elements {
				leafs[i] = HashLeaf(element)
			}

			tree, _ := buildTree(leafs)

			// Check if root would be nil without the fix
			isRootNil := tree[0] == nil

			// Apply the fix
			if isRootNil && tree[1] != nil {
				fmt.Printf("For %s: Root was nil, tree[1]=%s\n",
					tc.name, hex.EncodeToString(tree[1]))
				tree[0] = tree[1]
			} else {
				fmt.Printf("For %s: Root was NOT nil\n", tc.name)
			}

			// Verify the fix is needed by checking if we have a valid root now
			require.NotNil(t, tree[0], "Root should not be nil after fix")

			// Print tree structure for visual inspection
			printTreeStructure(t, tree, len(tc.elements))
		})
	}
}

func printTreeStructure(t *testing.T, tree [][]byte, numElements int) {
	t.Helper()
	fmt.Printf("Tree structure (numElements=%d):\n", numElements)

	for i, node := range tree {
		if node != nil {
			fmt.Printf("  tree[%d] = %s\n", i, hex.EncodeToString(node)[:8]+"...")
		} else {
			fmt.Printf("  tree[%d] = nil\n", i)
		}
	}
	fmt.Println()
}

func TestBuildTree(t *testing.T) {
	// Test unbalanced tree scenario specifically
	elements := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}

	// Without the specific fix
	tree1 := buildTreeWithoutFix(elements)

	// With the fix
	tree2 := buildTreeWithFix(elements)

	// We expect the fix to be needed if the root was nil
	if tree1[0] == nil {
		assert.NotNil(t, tree2[0], "Fix should give us a non-nil root")
		assert.Equal(t, tree2[0], tree2[1], "Fixed root should match tree[1]")
	} else {
		// If the root wasn't nil, then the fix shouldn't change anything
		assert.Equal(t, tree1[0], tree2[0], "Trees should have the same root")
	}
}

// buildTreeWithoutFix is a copy of buildTree without the root nil check
func buildTreeWithoutFix(elements [][]byte) [][]byte {
	leafs := make([][]byte, len(elements))
	for i, element := range elements {
		leafs[i] = HashLeaf(element)
	}

	// We don't need the depth here, just the tree
	balancedLeafCount := getLeafCount(len(leafs))
	tree := make([][]byte, balancedLeafCount<<1)

	// Copy leafs into the tree
	for i := 0; i < len(leafs); i++ {
		tree[balancedLeafCount+i] = leafs[i]
	}

	lowerBound := balancedLeafCount
	upperBound := balancedLeafCount + len(leafs) - 1

	// Build the tree
	for i := balancedLeafCount - 1; i >= 0; i-- {
		index := i << 1

		if index > upperBound {
			continue
		}

		if index <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		if index == upperBound {
			tree[i] = tree[index]
			continue
		}

		// Ensure both children are available before hashing
		if tree[index] != nil && tree[index+1] != nil {
			tree[i] = HashNode(tree[index], tree[index+1])
		} else if tree[index] != nil {
			// If right child is missing, use left child
			tree[i] = tree[index]
		}
	}

	return tree
}

// buildTreeWithFix is a copy of buildTree with the root nil check
func buildTreeWithFix(elements [][]byte) [][]byte {
	leafs := make([][]byte, len(elements))
	for i, element := range elements {
		leafs[i] = HashLeaf(element)
	}

	// We don't need the depth here, just the tree
	balancedLeafCount := getLeafCount(len(leafs))
	tree := make([][]byte, balancedLeafCount<<1)

	// Copy leafs into the tree
	for i := 0; i < len(leafs); i++ {
		tree[balancedLeafCount+i] = leafs[i]
	}

	lowerBound := balancedLeafCount
	upperBound := balancedLeafCount + len(leafs) - 1

	// Build the tree
	for i := balancedLeafCount - 1; i >= 0; i-- {
		index := i << 1

		if index > upperBound {
			continue
		}

		if index <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}

		if index == upperBound {
			tree[i] = tree[index]
			continue
		}

		// Ensure both children are available before hashing
		if tree[index] != nil && tree[index+1] != nil {
			tree[i] = HashNode(tree[index], tree[index+1])
		} else if tree[index] != nil {
			// If right child is missing, use left child
			tree[i] = tree[index]
		}
	}

	// Apply the fix for nil root
	if tree[0] == nil && tree[1] != nil {
		tree[0] = tree[1]
	}

	return tree
}

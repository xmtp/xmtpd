package merkle_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/merkle"
)

func TestVerifyMultiProofSequential(t *testing.T) {
	cases := []struct {
		name            string
		leaves          []merkle.Leaf
		startIdx        int
		count           int
		wantErrCreate   bool
		wantErrGenerate bool
		wantErrVerify   bool
		wantProofCount  int
	}{
		{
			name:            "Balanced tree - 8 leaves - consecutive leaves from middle",
			leaves:          createLeaves(t, 8),
			startIdx:        2,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 8 leaves - consecutive leaves from start",
			leaves:          createLeaves(t, 8),
			startIdx:        0,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
		},
		{
			name:            "Balanced tree - 8 leaves - consecutive leaves from end",
			leaves:          createLeaves(t, 8),
			startIdx:        5,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
		},
		{
			name:            "Balanced tree - 16 leaves - large range",
			leaves:          createLeaves(t, 16),
			startIdx:        2,
			count:           8,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 4 leaves - all leaves",
			leaves:          createLeaves(t, 4),
			startIdx:        0,
			count:           4,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  0,
		},
		{
			name:            "Single leaf tree",
			leaves:          createLeaves(t, 1),
			startIdx:        0,
			count:           1,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  1,
		},
		{
			name:            "Unbalanced tree - 3 leaves - last leaf",
			leaves:          createLeaves(t, 3),
			startIdx:        2,
			count:           1,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
		},
		{
			name:            "Unbalanced tree - 7 leaves - middle leaves",
			leaves:          createLeaves(t, 7),
			startIdx:        2,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Out of bounds - start index negative",
			leaves:          createLeaves(t, 8),
			startIdx:        -1,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: true,
			wantErrVerify:   false,
			wantProofCount:  0,
		},
		{
			name:            "Out of bounds - count exceeds leaves",
			leaves:          createLeaves(t, 8),
			startIdx:        6,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: true,
			wantErrVerify:   false,
			wantProofCount:  0,
		},
		{
			name:            "Invalid input - zero count",
			leaves:          createLeaves(t, 8),
			startIdx:        2,
			count:           0,
			wantErrCreate:   false,
			wantErrGenerate: true,
			wantErrVerify:   false,
			wantProofCount:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := merkle.NewMerkleTree(tc.leaves)
			if tc.wantErrCreate {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			proof, err := tree.GenerateMultiProofSequential(tc.startIdx, tc.count)
			if tc.wantErrGenerate {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			ok, err := merkle.Verify(tree.Root(), proof)
			if tc.wantErrVerify {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, ok)

			assert.Equal(t, tc.wantProofCount, len(proof.GetProofs()), "Invalid proof count")
			assert.Equal(t, tc.count, len(proof.GetValues()), "Invalid values count")
		})
	}
}

func TestVerifyMultiProofWithIndices(t *testing.T) {
	cases := []struct {
		name            string
		leaves          []merkle.Leaf
		indices         []int
		wantErrCreate   bool
		wantErrGenerate bool
		wantErrVerify   bool
		wantProofCount  int
	}{
		{
			name:            "Balanced tree - 8 leaves - out-of-order indices",
			leaves:          createLeaves(t, 8),
			indices:         []int{3, 1, 6},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Balanced tree - 8 leaves - all indices",
			leaves:          createLeaves(t, 8),
			indices:         []int{0, 1, 2, 3, 4, 5, 6, 7},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  0, // No proofs needed when all leaves included
		},
		{
			name:            "Balanced tree - 8 leaves - single index",
			leaves:          createLeaves(t, 8),
			indices:         []int{4},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 8 leaves - sibling pairs",
			leaves:          createLeaves(t, 8),
			indices:         []int{2, 3}, // Siblings at the same parent
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
		},
		{
			name:            "Balanced tree - 8 leaves - different subtrees",
			leaves:          createLeaves(t, 8),
			indices:         []int{1, 6}, // From different parts of the tree
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Balanced tree - 16 leaves - various positions",
			leaves:          createLeaves(t, 16),
			indices:         []int{2, 7, 9, 15},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  8,
		},
		{
			name:            "Unbalanced tree - 5 leaves - odd indices",
			leaves:          createLeaves(t, 5),
			indices:         []int{1, 3},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Unbalanced tree - 6 leaves - edge positions",
			leaves:          createLeaves(t, 6),
			indices:         []int{0, 5},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Unbalanced tree - 7 leaves - first and last",
			leaves:          createLeaves(t, 7),
			indices:         []int{0, 6},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Unbalanced tree - 9 leaves - mixed positions",
			leaves:          createLeaves(t, 9),
			indices:         []int{0, 6, 8},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  7,
		},
		{
			name:            "Single leaf tree",
			leaves:          createLeaves(t, 1),
			indices:         []int{0},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  1,
		},
		{
			name:            "Invalid input - empty indices",
			leaves:          createLeaves(t, 8),
			indices:         []int{},
			wantErrCreate:   false,
			wantErrGenerate: true, // Should error on empty indices
			wantErrVerify:   false,
			wantProofCount:  0,
		},
		{
			name:            "Invalid input - duplicate indices",
			leaves:          createLeaves(t, 8),
			indices:         []int{1, 1, 3},
			wantErrCreate:   false,
			wantErrGenerate: true, // Should error on duplicate indices
			wantErrVerify:   false,
			wantProofCount:  0,
		},
		{
			name:            "Invalid input - out of bounds indices",
			leaves:          createLeaves(t, 8),
			indices:         []int{1, 9, 3},
			wantErrCreate:   false,
			wantErrGenerate: true, // Should error on out of bounds
			wantErrVerify:   false,
			wantProofCount:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := merkle.NewMerkleTree(tc.leaves)
			if tc.wantErrCreate {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			proof, err := tree.GenerateMultiProofWithIndices(tc.indices)
			if tc.wantErrGenerate {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			ok, err := merkle.Verify(tree.Root(), proof)
			if tc.wantErrVerify {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, ok)

			assert.Equal(t, tc.wantProofCount, len(proof.GetProofs()), "Invalid proof count")
			assert.Equal(t, len(tc.indices), len(proof.GetValues()), "Invalid values count")

			vals := proof.GetValues()
			for i := 1; i < len(vals); i++ {
				assert.LessOrEqual(t, vals[i-1].GetIndex(), vals[i].GetIndex())
			}
		})
	}
}

func TestVerifyEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		setupFunc func() ([]byte, *merkle.MultiProof, error)
		wantOK    bool
		wantErr   string
		errorIs   error
	}{
		{
			name: "NilRoot",
			setupFunc: func() ([]byte, *merkle.MultiProof, error) {
				tree, err := merkle.NewMerkleTree(createLeaves(t, 4))
				if err != nil {
					return nil, nil, err
				}
				multiProof, err := tree.GenerateMultiProofSequential(0, 1)
				return nil, multiProof, err
			},
			wantOK:  false,
			wantErr: "nil root",
			errorIs: merkle.ErrNilRoot,
		},
		{
			name: "TamperedRoot",
			setupFunc: func() ([]byte, *merkle.MultiProof, error) {
				tree, err := merkle.NewMerkleTree(createLeaves(t, 4))
				if err != nil {
					return nil, nil, err
				}
				multiProof, err := tree.GenerateMultiProofSequential(0, 1)
				if err != nil {
					return nil, nil, err
				}

				// Tamper with the root
				tamperedRoot := make([]byte, len(tree.Root()))
				copy(tamperedRoot, tree.Root())
				tamperedRoot[0] ^= 0xFF

				return tamperedRoot, multiProof, nil
			},
			wantOK: false,
		},
		{
			name: "TamperedProofValues",
			setupFunc: func() ([]byte, *merkle.MultiProof, error) {
				tree, err := merkle.NewMerkleTree(createLeaves(t, 4))
				if err != nil {
					return nil, nil, err
				}
				multiProof, err := tree.GenerateMultiProofSequential(0, 1)
				if err != nil {
					return nil, nil, err
				}

				return tree.Root(), multiProof, nil
			},
			wantOK: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root, multiProof, err := tc.setupFunc()
			require.NoError(t, err)

			verified, err := merkle.Verify(root, multiProof)

			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				if tc.errorIs != nil {
					assert.ErrorIs(t, err, tc.errorIs)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantOK, verified)
		})
	}
}

func TestMultiProofManipulation(t *testing.T) {
	leaves := createLeaves(t, 8)
	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	multiProof, err := tree.GenerateMultiProofWithIndices([]int{0, 3})
	require.NoError(t, err)

	verified, err := merkle.Verify(tree.Root(), multiProof)
	require.NoError(t, err)
	assert.True(t, verified)

	values := multiProof.GetValues()
	assert.Equal(t, 2, len(values))
	assert.True(t, bytes.Equal(leaves[0], values[0].GetValue()))
	assert.True(t, bytes.Equal(leaves[3], values[1].GetValue()))
}

func createLeaves(t *testing.T, count int) []merkle.Leaf {
	t.Helper()
	leaves := make([]merkle.Leaf, count)
	for i := 0; i < count; i++ {
		leaves[i] = []byte{byte(i + 1)}
	}
	return leaves
}

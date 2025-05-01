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
			wantProofCount:  4,
		},
		{
			name:            "Balanced tree - 8 leaves - consecutive leaves from start",
			leaves:          createLeaves(t, 8),
			startIdx:        0,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 8 leaves - consecutive leaves from end",
			leaves:          createLeaves(t, 8),
			startIdx:        5,
			count:           3,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 16 leaves - large range",
			leaves:          createLeaves(t, 16),
			startIdx:        2,
			count:           8,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Balanced tree - 4 leaves - all leaves",
			leaves:          createLeaves(t, 4),
			startIdx:        0,
			count:           4,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  1,
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
			wantProofCount:  4,
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

			assert.Equal(t, tc.startIdx, proof.GetStartingIndex(), "Invalid starting index")
			assert.Equal(t, len(tc.leaves), proof.GetLeafCount(), "Invalid leaf count")
			assert.Equal(t, tc.count, len(proof.GetLeaves()), "Invalid number of leaves in proof")
			assert.Equal(t, tc.wantProofCount, len(proof.GetProofElements()), "Invalid proof count")
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
			indices:         []int{2, 3, 4},
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
			wantProofCount:  1,
		},
		{
			name:            "Balanced tree - 8 leaves - single index",
			leaves:          createLeaves(t, 8),
			indices:         []int{4},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Balanced tree - 8 leaves - sibling pairs",
			leaves:          createLeaves(t, 8),
			indices:         []int{2, 3}, // Siblings at the same parent
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Balanced tree - 8 leaves - different subtrees",
			leaves:          createLeaves(t, 8),
			indices:         []int{3, 4}, // From different parts of the tree
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  5,
		},
		{
			name:            "Balanced tree - 16 leaves - different subtrees",
			leaves:          createLeaves(t, 16),
			indices:         []int{6, 7, 8, 9},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  5,
		},
		{
			name:            "Unbalanced tree - 5 leaves - left edge",
			leaves:          createLeaves(t, 5),
			indices:         []int{0},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  4,
		},
		{
			name:            "Unbalanced tree - 6 leaves - right edge",
			leaves:          createLeaves(t, 6),
			indices:         []int{5},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Unbalanced tree - 7 leaves - right edge",
			leaves:          createLeaves(t, 7),
			indices:         []int{6},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  3,
		},
		{
			name:            "Unbalanced tree - 10 leaves - left edge",
			leaves:          createLeaves(t, 10),
			indices:         []int{0},
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  5,
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
			indices:         []int{6, 7, 8},
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

			assert.Equal(t, tc.indices[0], proof.GetStartingIndex(), "Invalid starting index")
			assert.Equal(t, len(tc.leaves), proof.GetLeafCount(), "Invalid leaf count")
			assert.Equal(t, len(tc.indices), len(proof.GetLeaves()), "Invalid number of leaves in proof")
			assert.Equal(t, tc.wantProofCount, len(proof.GetProofElements()), "Invalid proof count")

			leaves := proof.GetLeaves()
			for i := 0; i < len(leaves); i++ {
				assert.Equal(t, tc.leaves[tc.indices[0]+i], leaves[i])
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

// TODO: Not sure what this is testing
func TestMultiProofManipulation(t *testing.T) {
	leaves := createLeaves(t, 8)
	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	multiProof, err := tree.GenerateMultiProofWithIndices([]int{0, 1})
	require.NoError(t, err)

	verified, err := merkle.Verify(tree.Root(), multiProof)
	require.NoError(t, err)
	assert.True(t, verified)

	assert.Equal(t, 2, len(multiProof.GetLeaves()))
	assert.True(t, bytes.Equal(leaves[0], multiProof.GetLeaves()[0]))
	assert.True(t, bytes.Equal(leaves[1], multiProof.GetLeaves()[1]))
}

func createLeaves(t *testing.T, count int) []merkle.Leaf {
	t.Helper()
	leaves := make([]merkle.Leaf, count)
	for i := 0; i < count; i++ {
		leaves[i] = []byte{byte(i + 1)}
	}
	return leaves
}

func TestGenerateSequentialProofBalancedSample(t *testing.T) {
	leaves := []merkle.Leaf{
		getBytesFromHexString("fe0403dd500c862ddbbb4736c071a575d9fd43fbdeea602833234e9cd198bb3f"),
		getBytesFromHexString("5129d77889ad52ef899ba3d85e01a6df5c230804ccc586faa7e8376b89aa0600"),
		getBytesFromHexString("f1d6a296210c112eda6ca002273abc7d21b933cf99ca0c8292e344ba2d62f750"),
		getBytesFromHexString("907f141e1702157a9ca13881b0ecb0cb10be9b588bd6499479496965ce933225"),
		getBytesFromHexString("c77ccc18e17ec5d50bbbf775f564debfa1f9da8031be2fc7ce618cea8a625226"),
		getBytesFromHexString("aad8eb018edb0690bbdcebe48db820a6a13e3f8102bb0e5535fdfed79eaa5c39"),
		getBytesFromHexString("374346003edf10ea2da568e8e973857ec4291bf10a8bf81612b36efaf6e93c23"),
		getBytesFromHexString("45a3c2e32a3c8ea5e34562806f8279e85cfad69f6015588b599dc2e774d50cee"),
		getBytesFromHexString("9322c7cbfb342280d72d4eba15c0d9711e3f730072375bc4b0cf538970014e7a"),
		getBytesFromHexString("b3017e1b1874f68c4508c0b147a4f2b96dd166274c81bb4a86082d391efad6fa"),
		getBytesFromHexString("76d7f81c6813544dc594db8b5f3d6cef292e91d02a50ad556e6bcddc1eb59f8c"),
		getBytesFromHexString("320cfee2fb238c49255d102213b395963c2a3ebd7a8ee424e8df02b01361e484"),
		getBytesFromHexString("697fa39a1e23106c54eeecf38d213bfe2fcf896f55b790815085241b7e7e7361"),
		getBytesFromHexString("510fde4a04ecbe4bf94a514ab836074270959d2b96492a56909706de0b0248be"),
		getBytesFromHexString("5fa7828c8c661f4243c97527d637bf318e5e801dcf3b488f0f302622abdf61e9"),
		getBytesFromHexString("6f7fafb63e26bafeb5fd21a76615d0f2b45e6684324e294f10078661c00b4026"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	assert.Equal(t, getBytesFromHexString("ebc81f451b32d23cb3dc568b79cf8abc5de0d4c1a2c16e4d98e83a5146889de2"), tree.Root())

	proof, err := tree.GenerateMultiProofSequential(9, 4)
	require.NoError(t, err)

	assert.Equal(t, 9, proof.GetStartingIndex())

	assert.Equal(t, 4, len(proof.GetLeaves()))
	assert.Equal(t, leaves[9], proof.GetLeaves()[0])
	assert.Equal(t, leaves[10], proof.GetLeaves()[1])
	assert.Equal(t, leaves[11], proof.GetLeaves()[2])
	assert.Equal(t, leaves[12], proof.GetLeaves()[3])

	assert.Equal(t, 5, len(proof.GetProofElements()))
	assert.Equal(t, merkle.ProofElement(merkle.IntToBytes32(16)), proof.GetProofElements()[0])

	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("3fe32d0f3b236fb4afa5481033e969b4fdf444abc1b0e34b365f88ab37231633")),
		proof.GetProofElements()[1],
	)
	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("280bab377a3bc33db3a5e33595d3436d7c000ebff79acc8d0c36d43855aa707e")),
		proof.GetProofElements()[2],
	)
	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("a51f72dc2f8caaf2547915aae8027a3447ed26d7c99ab8499225a9075522ab13")),
		proof.GetProofElements()[3],
	)
	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("6ea454f465ec9869bae0f4f1f75956a4382d35bb4b35d2ddb506e4aa54460e4c")),
		proof.GetProofElements()[4],
	)
}

func TestGenerateSequentialProofUnbalancedSample(t *testing.T) {
	leaves := []merkle.Leaf{
		getBytesFromHexString("fe0403dd500c862ddbbb4736c071a575d9fd43fbdeea602833234e9cd198bb3f"),
		getBytesFromHexString("5129d77889ad52ef899ba3d85e01a6df5c230804ccc586faa7e8376b89aa0600"),
		getBytesFromHexString("f1d6a296210c112eda6ca002273abc7d21b933cf99ca0c8292e344ba2d62f750"),
		getBytesFromHexString("907f141e1702157a9ca13881b0ecb0cb10be9b588bd6499479496965ce933225"),
		getBytesFromHexString("c77ccc18e17ec5d50bbbf775f564debfa1f9da8031be2fc7ce618cea8a625226"),
		getBytesFromHexString("aad8eb018edb0690bbdcebe48db820a6a13e3f8102bb0e5535fdfed79eaa5c39"),
		getBytesFromHexString("374346003edf10ea2da568e8e973857ec4291bf10a8bf81612b36efaf6e93c23"),
		getBytesFromHexString("45a3c2e32a3c8ea5e34562806f8279e85cfad69f6015588b599dc2e774d50cee"),
		getBytesFromHexString("9322c7cbfb342280d72d4eba15c0d9711e3f730072375bc4b0cf538970014e7a"),
		getBytesFromHexString("b3017e1b1874f68c4508c0b147a4f2b96dd166274c81bb4a86082d391efad6fa"),
		getBytesFromHexString("76d7f81c6813544dc594db8b5f3d6cef292e91d02a50ad556e6bcddc1eb59f8c"),
		getBytesFromHexString("320cfee2fb238c49255d102213b395963c2a3ebd7a8ee424e8df02b01361e484"),
		getBytesFromHexString("697fa39a1e23106c54eeecf38d213bfe2fcf896f55b790815085241b7e7e7361"),
	}

	tree, err := merkle.NewMerkleTree(leaves)
	require.NoError(t, err)

	assert.Equal(t, getBytesFromHexString("4b17cd88068bf71cf6930362daca6bc8c45f0f0629bfc6436e6eb518bd086d61"), tree.Root())

	proof, err := tree.GenerateMultiProofSequential(9, 4)
	require.NoError(t, err)

	assert.Equal(t, 9, proof.GetStartingIndex())

	assert.Equal(t, 4, len(proof.GetLeaves()))
	assert.Equal(t, leaves[9], proof.GetLeaves()[0])
	assert.Equal(t, leaves[10], proof.GetLeaves()[1])
	assert.Equal(t, leaves[11], proof.GetLeaves()[2])
	assert.Equal(t, leaves[12], proof.GetLeaves()[3])

	assert.Equal(t, 3, len(proof.GetProofElements()))
	assert.Equal(t, merkle.ProofElement(merkle.IntToBytes32(13)), proof.GetProofElements()[0])

	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("280bab377a3bc33db3a5e33595d3436d7c000ebff79acc8d0c36d43855aa707e")),
		proof.GetProofElements()[1],
	)
	assert.Equal(
		t,
		merkle.ProofElement(getBytesFromHexString("6ea454f465ec9869bae0f4f1f75956a4382d35bb4b35d2ddb506e4aa54460e4c")),
		proof.GetProofElements()[2],
	)
}

func TestVerifySample1(t *testing.T) {
	leaves := []merkle.Leaf{
		getBytesFromHexString("fe0403dd500c862ddbbb4736c071a575d9fd43fbdeea602833234e9cd198bb3f"),
		getBytesFromHexString("5129d77889ad52ef899ba3d85e01a6df5c230804ccc586faa7e8376b89aa0600"),
		getBytesFromHexString("f1d6a296210c112eda6ca002273abc7d21b933cf99ca0c8292e344ba2d62f750"),
		getBytesFromHexString("907f141e1702157a9ca13881b0ecb0cb10be9b588bd6499479496965ce933225"),
		getBytesFromHexString("c77ccc18e17ec5d50bbbf775f564debfa1f9da8031be2fc7ce618cea8a625226"),
		getBytesFromHexString("aad8eb018edb0690bbdcebe48db820a6a13e3f8102bb0e5535fdfed79eaa5c39"),
		getBytesFromHexString("374346003edf10ea2da568e8e973857ec4291bf10a8bf81612b36efaf6e93c23"),
		getBytesFromHexString("45a3c2e32a3c8ea5e34562806f8279e85cfad69f6015588b599dc2e774d50cee"),
		getBytesFromHexString("9322c7cbfb342280d72d4eba15c0d9711e3f730072375bc4b0cf538970014e7a"),
		getBytesFromHexString("b3017e1b1874f68c4508c0b147a4f2b96dd166274c81bb4a86082d391efad6fa"),
		getBytesFromHexString("76d7f81c6813544dc594db8b5f3d6cef292e91d02a50ad556e6bcddc1eb59f8c"),
		getBytesFromHexString("320cfee2fb238c49255d102213b395963c2a3ebd7a8ee424e8df02b01361e484"),
		getBytesFromHexString("697fa39a1e23106c54eeecf38d213bfe2fcf896f55b790815085241b7e7e7361"),
		getBytesFromHexString("510fde4a04ecbe4bf94a514ab836074270959d2b96492a56909706de0b0248be"),
		getBytesFromHexString("5fa7828c8c661f4243c97527d637bf318e5e801dcf3b488f0f302622abdf61e9"),
		getBytesFromHexString("6f7fafb63e26bafeb5fd21a76615d0f2b45e6684324e294f10078661c00b4026"),
		getBytesFromHexString("e80b189bd7a1e5d0a04d9ba778d11f1bbffec4e276f502780a97fc7006715ac5"),
		getBytesFromHexString("25e66cd82e8fa7b6072d7759d2b49d0979bad20b018c730f4d1fcac7a6959bd4"),
		getBytesFromHexString("0695e86d43e2c5097acd182042d09cda438fc1f773335ce0f190dbed668ff508"),
		getBytesFromHexString("783eb995b1a99bb4181dad1c386dff06c5c02469a4893ca8c0a5b5fd82f7747a"),
		getBytesFromHexString("44fb0e55a04f4d5accd7a9e00fd5bbbf8926083ef169c0711f503ae16b589c20"),
		getBytesFromHexString("2815ef424dd0613d69f70546821ebddf2fe1b7452510cd21d25f3e438863e8a3"),
	}

	proofElements := []merkle.ProofElement{
		merkle.ProofElement(merkle.IntToBytes32(324)),
		merkle.ProofElement(getBytesFromHexString("55f9b393403d39fdf3b35fe8f394e13a272509df05d5562e5da6937c11bb1214")),
		merkle.ProofElement(getBytesFromHexString("a4c0ad363d17cb84be9e3f757278de2b26843e918e48e33e1b00438764fefa2b")),
		merkle.ProofElement(getBytesFromHexString("afee522c5ba27318361dd426e355fa02e5821fbc8b2ba52dea3b36885d76ec94")),
		merkle.ProofElement(getBytesFromHexString("25171fda9d5059c4cc4b8a86af5e0634dd8217cd2167c8d868733eef3f23054f")),
		merkle.ProofElement(getBytesFromHexString("1e3106cd69a1fec956fa702eae218bdba94f7150d87c05210b282e334218d265")),
		merkle.ProofElement(getBytesFromHexString("d2b4a539a1349d2d145fe307c30ce5bd43fa10887ce20a8a71da53300c0a6150")),
	}

	proof := merkle.NewMerkleProof(200, leaves, proofElements)

	root := getBytesFromHexString("8f28d0f19d1805a3b539bcc2bb0e627ded4bf5b873bebc9199ed179a30ca312c")

	ok, err := merkle.Verify(root, proof)
	require.NoError(t, err)
	require.True(t, ok)
}

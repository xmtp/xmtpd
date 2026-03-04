package merkle_test

import (
	"encoding/hex"
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
			name:            "Balanced tree - 4 leaves - no leaves",
			leaves:          createLeaves(t, 4),
			startIdx:        0,
			count:           0,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
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
			name:            "Unbalanced tree - 7 leaves - all leaves",
			leaves:          createLeaves(t, 7),
			startIdx:        0,
			count:           7,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  1,
		},
		{
			name:            "Unbalanced tree - 7 leaves - no leaves",
			leaves:          createLeaves(t, 7),
			startIdx:        0,
			count:           0,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  2,
		},
		{
			name:            "Empty tree",
			leaves:          createLeaves(t, 0),
			startIdx:        0,
			count:           0,
			wantErrCreate:   false,
			wantErrGenerate: false,
			wantErrVerify:   false,
			wantProofCount:  1,
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

			leafCount, err := proof.GetLeafCount()
			require.NoError(t, err)
			assert.Equal(t, len(tc.leaves), leafCount, "Invalid leaf count")

			assert.Len(t, proof.GetLeaves(), tc.count, "Invalid number of leaves in proof")
			assert.Len(t, proof.GetProofElements(), tc.wantProofCount, "Invalid proof count")
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
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				if tc.errorIs != nil {
					require.ErrorIs(t, err, tc.errorIs)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantOK, verified)
		})
	}
}

func createLeaves(t *testing.T, count int) []merkle.Leaf {
	t.Helper()
	leaves := make([]merkle.Leaf, count)
	for i := range count {
		leaves[i] = []byte{byte(i + 1)}
	}
	return leaves
}

func TestGenerateSequentialProofBalancedSamples(t *testing.T) {
	tests := []struct {
		name                  string
		leaves                []string
		startingIndex         int
		count                 int
		expectedLeaves        []string
		expectedProofElements []string
		expectedRoot          string
	}{
		{
			name: "sample 1",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
			},
			startingIndex: 0,
			count:         1,
			expectedLeaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
			},
			expectedProofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000002",
				"b8283b9a33ca222061b7d8d5f85170289c5c4a7f96997ce2ae5bafc94fc8a59b",
			},
			expectedRoot: "eeef536868dc2c030bec2d3602cc13fbe660bd5d63deca6a0a4dfd201eb941c0",
		},
		{
			name: "sample 2",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
				"aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c",
				"e83870d75c6c4c4d1f6ba674481932301e0a1029b44c1407b6aea06cd56d4836",
			},
			startingIndex: 1,
			count:         4,
			expectedLeaves: []string{
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
			},
			expectedProofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000008",
				"b7f40118daf8bbb92b1759e810c6eb4b0b92e04d1b7e4be48a147721ea457d87",
				"06331c9b61be683c1819b2cd20c83f315499477f6253ae8d8bb02cbc6bd93c9f",
				"c3d064e0f1ace0f92e6c822d5346de5330fb23a6c31a0c84a80d9d5c8543cc0e",
			},
			expectedRoot: "00f8c0ad3c60c727ededce5717c8baa64047b5c3f29e409085df14dc3bfda1a7",
		},
		{
			name: "sample 3",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
				"aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c",
				"e83870d75c6c4c4d1f6ba674481932301e0a1029b44c1407b6aea06cd56d4836",
				"2815ef424dd0613d69f70546821ebddf2fe1b7452510cd21d25f3e438863e8a3",
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
				"9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba",
				"deadbeefcafebabe0123456789abcdefdeadbeefcafebabe0123456789abcdef",
				"cafebabebeefdeadabcdef0123456789cafebabebeefdeadabcdef0123456789",
			},
			startingIndex: 9,
			count:         4,
			expectedLeaves: []string{
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
			expectedProofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000010",
				"194f8cb59bc03a6e9dfcc3566bed0a85a069f50de609f51af6c74dee673450d4",
				"2642cf6153e4ddfd7ea9c1c99d86efe99c03d29bce0ce4adbf7c0162865aa93a",
				"9d8799ecccbca75f60304876c8426007565302fdb72449fe980917d7847d43e7",
				"1f70e7dd11a042e3868e8b0992118a3d7bd301b029a3b967a5b2042466c5110c",
			},
			expectedRoot: "31338b156e26447f0a3a965981b0be87957bf5606b44e0dcdc99eb4646048942",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leaves := make([]merkle.Leaf, len(test.leaves))
			for i, leaf := range test.leaves {
				leaves[i] = getBytesFromHexString(leaf)
			}

			tree, err := merkle.NewMerkleTree(leaves)
			require.NoError(t, err)

			assert.Equal(
				t,
				getBytesFromHexString(test.expectedRoot),
				tree.Root(),
			)

			proof, err := tree.GenerateMultiProofSequential(test.startingIndex, test.count)
			require.NoError(t, err)

			assert.Equal(t, test.startingIndex, proof.GetStartingIndex())

			leafCount, err := proof.GetLeafCount()
			require.NoError(t, err)
			assert.Equal(t, tree.LeafCount(), leafCount)

			assert.Len(t, proof.GetLeaves(), test.count)

			for i := range len(test.expectedLeaves) {
				assert.Equal(t, test.expectedLeaves[i], hex.EncodeToString(proof.GetLeaves()[i]))
			}

			assert.Len(t, proof.GetProofElements(), len(test.expectedProofElements))

			for i := range len(test.expectedProofElements) {
				assert.Equal(
					t,
					test.expectedProofElements[i],
					hex.EncodeToString(proof.GetProofElements()[i]),
				)
			}
		})
	}
}

func TestGenerateSequentialProofUnbalancedSamples(t *testing.T) {
	tests := []struct {
		name                  string
		leaves                []string
		startingIndex         int
		count                 int
		expectedLeaves        []string
		expectedProofElements []string
		expectedRoot          string
	}{
		{
			name: "sample 1",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
			},
			startingIndex: 0,
			count:         1,
			expectedLeaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
			},
			expectedProofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000001",
			},
			expectedRoot: "5b833bdf4f55e39d1838653841d4a2c651a71b5626b7936e1bedb5212cae96e3",
		},
		{
			name: "sample 2",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
				"aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c",
			},
			startingIndex: 4,
			count:         2,
			expectedLeaves: []string{
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
			},
			expectedProofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000007",
				"6b379ab3dd3b0b8df76cbe09e2a04f21ebf601f2932d3505ecfb89417f3de836",
				"52bd35868bccc1b1f4c3ac67dd7e3b3db4a24f60436411162d633e6d1118de89",
			},
			expectedRoot: "38631dd8b5081555ec3c51cc8db7918ee90158fa33a70674c1399234d23908b2",
		},
		{
			name: "sample 3",
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
				"aa79d134afbdcf008b487dbab5717dfc6518bffd2dc6ce71724a9e87200a086c",
				"e83870d75c6c4c4d1f6ba674481932301e0a1029b44c1407b6aea06cd56d4836",
				"2815ef424dd0613d69f70546821ebddf2fe1b7452510cd21d25f3e438863e8a3",
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
			startingIndex: 9,
			count:         4,
			expectedLeaves: []string{
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
			expectedProofElements: []string{
				"000000000000000000000000000000000000000000000000000000000000000d",
				"2642cf6153e4ddfd7ea9c1c99d86efe99c03d29bce0ce4adbf7c0162865aa93a",
				"1f70e7dd11a042e3868e8b0992118a3d7bd301b029a3b967a5b2042466c5110c",
			},
			expectedRoot: "f92d4ca528834b0350cecd9307bec2dd97d0a6bbb58b077ab51cdad36fc5c087",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leaves := make([]merkle.Leaf, len(test.leaves))
			for i, leaf := range test.leaves {
				leaves[i] = getBytesFromHexString(leaf)
			}

			tree, err := merkle.NewMerkleTree(leaves)
			require.NoError(t, err)

			assert.Equal(
				t,
				getBytesFromHexString(test.expectedRoot),
				tree.Root(),
			)

			proof, err := tree.GenerateMultiProofSequential(test.startingIndex, test.count)
			require.NoError(t, err)

			assert.Equal(t, test.startingIndex, proof.GetStartingIndex())

			leafCount, err := proof.GetLeafCount()
			require.NoError(t, err)
			assert.Equal(t, tree.LeafCount(), leafCount)

			assert.Len(t, proof.GetLeaves(), test.count)

			for i := range len(test.expectedLeaves) {
				assert.Equal(t, test.expectedLeaves[i], hex.EncodeToString(proof.GetLeaves()[i]))
			}

			assert.Len(t, proof.GetProofElements(), len(test.expectedProofElements))

			for i := range len(test.expectedProofElements) {
				assert.Equal(
					t,
					test.expectedProofElements[i],
					hex.EncodeToString(proof.GetProofElements()[i]),
				)
			}
		})
	}
}

func TestVerifySequentialProofBalancedSamples(t *testing.T) {
	tests := []struct {
		name          string
		startingIndex int
		leaves        []string
		proofElements []string
		root          string
	}{
		{
			name:          "sample 1",
			startingIndex: 0,
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000002",
				"b8283b9a33ca222061b7d8d5f85170289c5c4a7f96997ce2ae5bafc94fc8a59b",
			},
			root: "eeef536868dc2c030bec2d3602cc13fbe660bd5d63deca6a0a4dfd201eb941c0",
		},
		{
			name:          "sample 2",
			startingIndex: 1,
			leaves: []string{
				"a8152e7c56b62d9fcb8af361257a260b2b9481c8683e8df1651a31508cc6ee31",
				"007f47e1c51d53cab18977050347e8e8dc488bdd9590babe3e104fcb9a1ef599",
				"7cbe68a29af312d42c40e6d083bb64fe2ba0ac6bf1cac8e4b10f5356142e3828",
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000008",
				"b7f40118daf8bbb92b1759e810c6eb4b0b92e04d1b7e4be48a147721ea457d87",
				"06331c9b61be683c1819b2cd20c83f315499477f6253ae8d8bb02cbc6bd93c9f",
				"c3d064e0f1ace0f92e6c822d5346de5330fb23a6c31a0c84a80d9d5c8543cc0e",
			},
			root: "00f8c0ad3c60c727ededce5717c8baa64047b5c3f29e409085df14dc3bfda1a7",
		},
		{
			name:          "sample 3",
			startingIndex: 9,
			leaves: []string{
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000010",
				"194f8cb59bc03a6e9dfcc3566bed0a85a069f50de609f51af6c74dee673450d4",
				"2642cf6153e4ddfd7ea9c1c99d86efe99c03d29bce0ce4adbf7c0162865aa93a",
				"9d8799ecccbca75f60304876c8426007565302fdb72449fe980917d7847d43e7",
				"1f70e7dd11a042e3868e8b0992118a3d7bd301b029a3b967a5b2042466c5110c",
			},
			root: "31338b156e26447f0a3a965981b0be87957bf5606b44e0dcdc99eb4646048942",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leaves := make([]merkle.Leaf, len(test.leaves))
			for i, leaf := range test.leaves {
				leaves[i] = getBytesFromHexString(leaf)
			}

			proofElements := make([]merkle.ProofElement, len(test.proofElements))
			for i, proofElement := range test.proofElements {
				proofElements[i] = getBytesFromHexString(proofElement)
			}

			proof := merkle.NewMerkleProof(test.startingIndex, leaves, proofElements)

			root := getBytesFromHexString(test.root)

			ok, err := merkle.Verify(root, proof)
			require.NoError(t, err)
			require.True(t, ok)
		})
	}
}

func TestVerifySequentialProofUnbalancedSamples(t *testing.T) {
	tests := []struct {
		name          string
		startingIndex int
		leaves        []string
		proofElements []string
		root          string
	}{
		{
			name:          "sample 1",
			startingIndex: 0,
			leaves: []string{
				"6330b989705733cc5c1f7285b8a5b892e08be86ed6fbe9d254713a4277bc5bd2",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000001",
			},
			root: "5b833bdf4f55e39d1838653841d4a2c651a71b5626b7936e1bedb5212cae96e3",
		},
		{
			name:          "sample 2",
			startingIndex: 4,
			leaves: []string{
				"4a864e860c0d0247c6aa5ebcb2bc3f15fc4ddf86213258f4bf0b72e51c9d9c69",
				"51b7ae2bab96bd3fbb3b26e1efefb0b9b6a60054ed7ffcfa700374d58f315a31",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000007",
				"6b379ab3dd3b0b8df76cbe09e2a04f21ebf601f2932d3505ecfb89417f3de836",
				"52bd35868bccc1b1f4c3ac67dd7e3b3db4a24f60436411162d633e6d1118de89",
			},
			root: "38631dd8b5081555ec3c51cc8db7918ee90158fa33a70674c1399234d23908b2",
		},
		{
			name:          "sample 3",
			startingIndex: 9,
			leaves: []string{
				"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
				"112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
				"f0e1d2c3b4a5968778695a4b3c2d1e0f1f2e3d4c5b6a79887766554433221100",
				"abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
			proofElements: []string{
				"000000000000000000000000000000000000000000000000000000000000000d",
				"2642cf6153e4ddfd7ea9c1c99d86efe99c03d29bce0ce4adbf7c0162865aa93a",
				"1f70e7dd11a042e3868e8b0992118a3d7bd301b029a3b967a5b2042466c5110c",
			},
			root: "f92d4ca528834b0350cecd9307bec2dd97d0a6bbb58b077ab51cdad36fc5c087",
		},
		{
			name:          "sample 4",
			startingIndex: 200,
			leaves: []string{
				"fe0403dd500c862ddbbb4736c071a575d9fd43fbdeea602833234e9cd198bb3f",
				"5129d77889ad52ef899ba3d85e01a6df5c230804ccc586faa7e8376b89aa0600",
				"f1d6a296210c112eda6ca002273abc7d21b933cf99ca0c8292e344ba2d62f750",
				"907f141e1702157a9ca13881b0ecb0cb10be9b588bd6499479496965ce933225",
				"c77ccc18e17ec5d50bbbf775f564debfa1f9da8031be2fc7ce618cea8a625226",
				"aad8eb018edb0690bbdcebe48db820a6a13e3f8102bb0e5535fdfed79eaa5c39",
				"374346003edf10ea2da568e8e973857ec4291bf10a8bf81612b36efaf6e93c23",
				"45a3c2e32a3c8ea5e34562806f8279e85cfad69f6015588b599dc2e774d50cee",
				"9322c7cbfb342280d72d4eba15c0d9711e3f730072375bc4b0cf538970014e7a",
				"b3017e1b1874f68c4508c0b147a4f2b96dd166274c81bb4a86082d391efad6fa",
				"76d7f81c6813544dc594db8b5f3d6cef292e91d02a50ad556e6bcddc1eb59f8c",
				"320cfee2fb238c49255d102213b395963c2a3ebd7a8ee424e8df02b01361e484",
				"697fa39a1e23106c54eeecf38d213bfe2fcf896f55b790815085241b7e7e7361",
				"510fde4a04ecbe4bf94a514ab836074270959d2b96492a56909706de0b0248be",
				"5fa7828c8c661f4243c97527d637bf318e5e801dcf3b488f0f302622abdf61e9",
				"6f7fafb63e26bafeb5fd21a76615d0f2b45e6684324e294f10078661c00b4026",
				"e80b189bd7a1e5d0a04d9ba778d11f1bbffec4e276f502780a97fc7006715ac5",
				"25e66cd82e8fa7b6072d7759d2b49d0979bad20b018c730f4d1fcac7a6959bd4",
				"0695e86d43e2c5097acd182042d09cda438fc1f773335ce0f190dbed668ff508",
				"783eb995b1a99bb4181dad1c386dff06c5c02469a4893ca8c0a5b5fd82f7747a",
				"44fb0e55a04f4d5accd7a9e00fd5bbbf8926083ef169c0711f503ae16b589c20",
				"2815ef424dd0613d69f70546821ebddf2fe1b7452510cd21d25f3e438863e8a3",
			},
			proofElements: []string{
				"0000000000000000000000000000000000000000000000000000000000000144",
				"55f9b393403d39fdf3b35fe8f394e13a272509df05d5562e5da6937c11bb1214",
				"a4c0ad363d17cb84be9e3f757278de2b26843e918e48e33e1b00438764fefa2b",
				"afee522c5ba27318361dd426e355fa02e5821fbc8b2ba52dea3b36885d76ec94",
				"25171fda9d5059c4cc4b8a86af5e0634dd8217cd2167c8d868733eef3f23054f",
				"1e3106cd69a1fec956fa702eae218bdba94f7150d87c05210b282e334218d265",
				"d2b4a539a1349d2d145fe307c30ce5bd43fa10887ce20a8a71da53300c0a6150",
			},
			root: "8f28d0f19d1805a3b539bcc2bb0e627ded4bf5b873bebc9199ed179a30ca312c",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leaves := make([]merkle.Leaf, len(test.leaves))
			for i, leaf := range test.leaves {
				leaves[i] = getBytesFromHexString(leaf)
			}

			proofElements := make([]merkle.ProofElement, len(test.proofElements))
			for i, proofElement := range test.proofElements {
				proofElements[i] = getBytesFromHexString(proofElement)
			}

			proof := merkle.NewMerkleProof(test.startingIndex, leaves, proofElements)

			root := getBytesFromHexString(test.root)

			ok, err := merkle.Verify(root, proof)
			require.NoError(t, err)
			require.True(t, ok)
		})
	}
}

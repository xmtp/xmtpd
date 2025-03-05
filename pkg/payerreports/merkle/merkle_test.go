package merkle_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/payerreports/merkle"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildVerifier(t *testing.T) (*blockchain.MerkleCaller, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	contractsOptions := testutils.GetContractsOptions(t)
	client, err := blockchain.NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)

	caller, err := blockchain.NewMerkleCaller(logger, client, contractsOptions)
	require.NoError(t, err)

	return caller, func() {
		cancel()
		client.Close()
	}
}

func buildInputs(numInputs int) []merkle.PayerReportInput {
	inputs := make([]merkle.PayerReportInput, numInputs)
	for i := range numInputs {
		inputs[i] = merkle.PayerReportInput{
			PayerAddress: testutils.RandomAddress(),
			Amount:       big.NewInt(int64(testutils.RandomInt32())),
		}
	}

	return inputs
}

func TestBuildTree(t *testing.T) {
	inputs := buildInputs(10)
	tree, err := merkle.NewPayerReportTree(inputs)
	require.NoError(t, err)
	require.NotNil(t, tree)

	root := tree.Root()
	require.NotNil(t, root)
}

func TestBuildTreeEmpty(t *testing.T) {
	inputs := []merkle.PayerReportInput{}
	tree, err := merkle.NewPayerReportTree(inputs)
	require.Error(t, err)
	require.Nil(t, tree)
}

func TestGetMultiProof(t *testing.T) {
	inputs := buildInputs(10)
	tree, err := merkle.NewPayerReportTree(inputs)
	require.NoError(t, err)
	require.NotNil(t, tree)

	proof, err := tree.GetMultiProof(0, 1)
	require.NoError(t, err)
	require.NotNil(t, proof)

	proof2, err := tree.GetMultiProof(1, 1)
	require.NoError(t, err)
	require.NotNil(t, proof2)
	require.NotEqual(t, proof.Proof(), proof2.Proof())

	proof, err = tree.GetMultiProof(10, 1)
	require.Error(t, err)
	require.Nil(t, proof)
}

func TestVerifyMultiProof(t *testing.T) {
	inputs := buildInputs(10)
	tree, err := merkle.NewPayerReportTree(inputs)
	require.NoError(t, err)
	require.NotNil(t, tree)

	verifier, cleanup := buildVerifier(t)
	defer cleanup()

	proof, err := tree.GetMultiProof(0, 6)
	require.NoError(t, err)
	require.NotNil(t, proof)

	verified, err := verifier.VerifyMultiProof(context.Background(), proof)
	require.NoError(t, err)
	require.True(t, verified)
}

package builders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func RegisterNode(t *testing.T, network string, rpcHost string, xmtpdAlias string) {
	envVars := constructVariables(t)
	httpAddress := fmt.Sprintf("http://%s:5050", xmtpdAlias)
	signerPublicKey := envVars["XMTPD_SIGNER_PUBLIC_KEY"]
	require.NotEmpty(t, signerPublicKey)
	signerAddress := envVars["XMTPD_SIGNER_ADDRESS"]
	require.NotEmpty(t, signerAddress)

	registerNode := []string{
		"--config-file=/cfg/anvil.json",
		"--private-key=" + adminPrivateKey,
		"--settlement-rpc-url=" + rpcHost,
		"nodes", "register",
		"--owner-address=" + signerAddress,
		"--signing-key-pub=" + signerPublicKey,
		"--http-address=" + httpAddress,
	}
	err := NewCLIContainerBuilder(t).WithCmd(registerNode).WithNetwork(network).Build(t)
	require.NoError(t, err)
}

func EnableNode(t *testing.T, network string, rpcHost string, nodeID uint32) {
	enableNode := []string{
		"--config-file=/cfg/anvil.json",
		"--private-key=" + adminPrivateKey,
		"--settlement-rpc-url=" + rpcHost,
		"nodes", "canonical-network",
		"--add",
		fmt.Sprintf("--node-id=%d", nodeID),
	}
	err := NewCLIContainerBuilder(t).WithCmd(enableNode).WithNetwork(network).Build(t)
	require.NoError(t, err)
}

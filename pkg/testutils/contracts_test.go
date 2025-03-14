package testutils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func TestDeployContract(t *testing.T) {
	rpcUrl, cleanup := anvil.StartAnvil(t, false)
	defer cleanup()
	deployedTo := DeployNodesContract(t, rpcUrl)
	require.True(t, common.IsHexAddress(deployedTo), "invalid contract address")
}

func TestDeployGroupMessages(t *testing.T) {
	rpcUrl, cleanup := anvil.StartAnvil(t, false)
	defer cleanup()
	deployedTo := DeployGroupMessagesContract(t, rpcUrl)
	require.True(t, common.IsHexAddress(deployedTo), "invalid contract address")
}

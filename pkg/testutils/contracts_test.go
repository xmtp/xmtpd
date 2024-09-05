package testutils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDeployContract(t *testing.T) {
	deployedTo := DeployNodesContract(t)
	require.True(t, common.IsHexAddress(deployedTo), "invalid contract address")
}

func TestDeployGroupMessages(t *testing.T) {
	deployedTo := DeployGroupMessagesContract(t)
	require.True(t, common.IsHexAddress(deployedTo), "invalid contract address")
}

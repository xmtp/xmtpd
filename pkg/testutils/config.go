package testutils

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
)

const BLOCKCHAIN_RPC_URL = "http://localhost:7545"

// This is the private key that anvil has funded by default
// This is safe to commit
const TEST_PRIVATE_KEY = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

type contractInfo struct {
	DeployedTo string `json:"deployedTo"`
}

/*
*
In tests it's weirdly difficult to get the working directory of the project root.

Keep moving up the folder hierarchy until you find a go.mod
*
*/
func rootPath(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir { // reached the root directory
			t.Fatal("Could not find the root directory")
		}
		dir = parent
	}
}

/*
*
Parse the JSON file at this location to
*
*/
func getDeployedTo(t *testing.T, fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read GroupMessages.json: %v", err)
	}

	var info contractInfo

	if err := json.Unmarshal(data, &info); err != nil {
		t.Fatalf("Failed to parse GroupMessages.json: %v", err)
	}

	return info.DeployedTo
}

func GetContractsOptions(t *testing.T) config.ContractsOptions {
	rootDir := rootPath(t)

	return config.ContractsOptions{
		RpcUrl:                  BLOCKCHAIN_RPC_URL,
		MessagesContractAddress: getDeployedTo(t, path.Join(rootDir, "./build/GroupMessages.json")),
		NodesContractAddress:    getDeployedTo(t, path.Join(rootDir, "./build/Nodes.json")),
		RefreshInterval:         100 * time.Millisecond,
		ChainID:                 31337,
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TEST_PRIVATE_KEY,
	}
}

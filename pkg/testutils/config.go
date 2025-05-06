package testutils

import (
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
)

const BLOCKCHAIN_RPC_URL = "http://localhost:7545"

// This is the private key that anvil has funded by default
// This is safe to commit
const TEST_PRIVATE_KEY = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func NewContractsOptions(rpcUrl string) config.ContractsOptions {
	return config.ContractsOptions{
		AppChain: config.AppChainOptions{
			RpcURL:                 rpcUrl,
			ChainID:                31337,
			MaxChainDisconnectTime: 10 * time.Second,
		},
		SettlementChain: config.SettlementChainOptions{
			RpcURL:                      rpcUrl,
			NodeRegistryRefreshInterval: 100 * time.Millisecond,
			ChainID:                     31337,
			RateRegistryRefreshInterval: 10 * time.Second,
		},
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TEST_PRIVATE_KEY,
	}
}

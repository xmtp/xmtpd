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

const RATE_REGISTRY_ADDRESS = "0xE71ac6dE80392495eB52FB1dCa321f5dB8f51BAE"
const NODE_REGISTRY_ADDRESS = "0x8d69E9834f1e4b38443C638956F7D81CD04eBB2F"
const GROUP_MESSAGE_BROADCAST_ADDRESS = "0xD5b7B43B0e31112fF99Bd5d5C4f6b828259bedDE"
const IDENTITY_UPDATE_BROADCAST_ADDRESS = "0xe67104BC93003192ab78B797d120DBA6e9Ff4928"

func NewContractsOptions(rpcUrl string) config.ContractsOptions {
	return config.ContractsOptions{
		AppChain: config.AppChainOptions{
			RpcURL:                           rpcUrl,
			ChainID:                          31337,
			MaxChainDisconnectTime:           10 * time.Second,
			GroupMessageBroadcasterAddress:   GROUP_MESSAGE_BROADCAST_ADDRESS,
			IdentityUpdateBroadcasterAddress: IDENTITY_UPDATE_BROADCAST_ADDRESS,
		},
		SettlementChain: config.SettlementChainOptions{
			RpcURL:                      rpcUrl,
			NodeRegistryRefreshInterval: 100 * time.Millisecond,
			ChainID:                     31337,
			RateRegistryRefreshInterval: 10 * time.Second,
			RateRegistryAddress:         RATE_REGISTRY_ADDRESS,
			NodeRegistryAddress:         NODE_REGISTRY_ADDRESS,
		},
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TEST_PRIVATE_KEY,
	}
}

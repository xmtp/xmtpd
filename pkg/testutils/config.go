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

const (
	NODE_REGISTRY_ADDRESS        = "0xB47363AfbeAf04a8Fbaf6bE085050A979d2a9794"
	RATE_REGISTRY_ADDRESS        = "0xe9Fb03945475587B03eA28b8118E2bc5B753E3E9"
	PAYER_REGISTRY_ADDRESS       = "0x95E856F1542EB9Eb1BFB6019f0D438584b1652ed"
	PAYER_REPORT_MANAGER_ADDRESS = "0x87cbA0310D1f1a2bF15408b54d987Bf3ec45B2a5"
	DISTRIBUTION_MANAGER_ADDRESS = "0x73A4846B953EFcD9242A3e94666DEE3312EE8a5F"
	PARAMETER_REGISTRY_ADDRESS   = "0x866e7279B86a71093F2a601883b9c66EdB320ddD"

	GROUP_MESSAGE_BROADCAST_ADDRESS   = "0x8c5908AFbd1a5C25590D78eC7Bb0422262BDE6a1"
	IDENTITY_UPDATE_BROADCAST_ADDRESS = "0x2c7A0c3856ca0CC9bf339E19fE25ca4c1f57A567"
)

func NewContractsOptions(rpcUrl string) config.ContractsOptions {
	return config.ContractsOptions{
		AppChain: config.AppChainOptions{
			RpcURL:                           rpcUrl,
			ChainID:                          31337,
			MaxChainDisconnectTime:           10 * time.Second,
			GroupMessageBroadcasterAddress:   GROUP_MESSAGE_BROADCAST_ADDRESS,
			IdentityUpdateBroadcasterAddress: IDENTITY_UPDATE_BROADCAST_ADDRESS,
			BackfillBlockSize:                500,
		},
		SettlementChain: config.SettlementChainOptions{
			RpcURL:                      rpcUrl,
			NodeRegistryRefreshInterval: 100 * time.Millisecond,
			ChainID:                     31337,
			RateRegistryRefreshInterval: 10 * time.Second,
			RateRegistryAddress:         RATE_REGISTRY_ADDRESS,
			NodeRegistryAddress:         NODE_REGISTRY_ADDRESS,
			ParameterRegistryAddress:    PARAMETER_REGISTRY_ADDRESS,
			PayerRegistryAddress:        PAYER_REGISTRY_ADDRESS,
			PayerReportManagerAddress:   PAYER_REPORT_MANAGER_ADDRESS,
			MaxChainDisconnectTime:      10 * time.Second,
			BackfillBlockSize:           500,
		},
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TEST_PRIVATE_KEY,
	}
}

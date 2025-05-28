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
	RATE_REGISTRY_ADDRESS             = "0xc76Fdbd3b026b73fcC03eB1D5CD779D771eA03b6"
	NODE_REGISTRY_ADDRESS             = "0xf82fd040a85154B4914a6cF81aCc3F83316c72CD"
	GROUP_MESSAGE_BROADCAST_ADDRESS   = "0x7D7df4ed008d64b9b816007b69dcbf559b9C5494"
	IDENTITY_UPDATE_BROADCAST_ADDRESS = "0xDCbC334a97c6a8DBe7d673bd52fA56718708BC9C"
	PARAMETER_REGISTRY_ADDRESS        = "0x7FFc148AF5f00C56D78Bce732fC79a08007eC8be"
	PAYER_REGISTRY_ADDRESS            = "0x6ADCc064469C3b69ED6Fa4DbeFA21490FD200D7c"
	PAYER_REPORT_MANAGER_ADDRESS      = "0x2aA7BC557FF0b9B55FFD82274706DD2aD37E687B"
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

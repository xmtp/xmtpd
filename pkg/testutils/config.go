package testutils

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/config/environments"
)

// This is the private key that anvil has funded by default
// This is safe to commit
const (
	TestPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

func NewContractsOptions(t *testing.T, rpcURL, wsURL string) config.ContractsOptions {
	anvilJson, err := environments.GetEnvironmentConfig(environments.Anvil)
	require.NoError(t, err)
	var chainConfig config.ChainConfig
	err = json.Unmarshal(anvilJson, &chainConfig)
	require.NoError(t, err)

	return config.ContractsOptions{
		AppChain: config.AppChainOptions{
			RPCURL:                           rpcURL,
			WssURL:                           wsURL,
			ChainID:                          31337,
			MaxChainDisconnectTime:           10 * time.Second,
			GroupMessageBroadcasterAddress:   chainConfig.GroupMessageBroadcaster,
			IdentityUpdateBroadcasterAddress: chainConfig.IdentityUpdateBroadcaster,
			BackfillBlockPageSize:            500,
			GatewayAddress:                   chainConfig.AppChainGateway,
			ParameterRegistryAddress:         chainConfig.AppChainParameterRegistry,
		},
		SettlementChain: config.SettlementChainOptions{
			RPCURL:                      rpcURL,
			WssURL:                      wsURL,
			NodeRegistryRefreshInterval: 100 * time.Millisecond,
			ChainID:                     31337,
			RateRegistryRefreshInterval: 10 * time.Second,
			RateRegistryAddress:         chainConfig.RateRegistry,
			NodeRegistryAddress:         chainConfig.NodeRegistry,
			ParameterRegistryAddress:    chainConfig.SettlementChainParameterRegistry,
			PayerRegistryAddress:        chainConfig.PayerRegistry,
			PayerReportManagerAddress:   chainConfig.PayerReportManager,
			MaxChainDisconnectTime:      10 * time.Second,
			BackfillBlockPageSize:       500,
			GatewayAddress:              chainConfig.SettlementChainGateway,
			DistributionManagerAddress:  chainConfig.DistributionManager,
			UnderlyingFeeToken:          chainConfig.UnderlyingFeeToken,
			FeeToken:                    chainConfig.FeeToken,
		},
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TestPrivateKey,
	}
}

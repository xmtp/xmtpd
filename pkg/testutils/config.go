package testutils

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/config"
)

// This is the private key that anvil has funded by default
// This is safe to commit
const (
	TestPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	thisFileDepth  = "../.."
	pathToJSONFile = "dev/environments/anvil.json"
)

func NewContractsOptions(t *testing.T, rpcURL, wsURL string) config.ContractsOptions {
	file, err := os.Open(GetScriptPath(thisFileDepth, pathToJSONFile))
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	require.NoError(t, err)

	var chainConfig config.ChainConfig
	err = json.Unmarshal(data, &chainConfig)
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
		},
	}
}

func GetPayerOptions(t *testing.T) config.PayerOptions {
	return config.PayerOptions{
		PrivateKey: TestPrivateKey,
	}
}

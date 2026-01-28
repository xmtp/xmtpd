package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadContractsConfig_Environment(t *testing.T) {
	opts, err := LoadContractsConfig(ContractsSource{Environment: "testnet"})
	require.NoError(t, err)
	require.NotEmpty(t, opts.AppChain.GroupMessageBroadcasterAddress)
	require.NotEmpty(t, opts.SettlementChain.NodeRegistryAddress)
}

func TestLoadContractsConfig_ConfigURL(t *testing.T) {
	// config:// URL scheme should resolve to environment
	opts, err := LoadContractsConfig(ContractsSource{FilePath: "config://testnet"})
	require.NoError(t, err)
	require.NotEmpty(t, opts.AppChain.GroupMessageBroadcasterAddress)
}

func TestLoadContractsConfig_MutuallyExclusive(t *testing.T) {
	_, err := LoadContractsConfig(ContractsSource{
		Environment: "testnet",
		FilePath:    "/some/path",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")
}

func TestLoadContractsConfig_NoSource(t *testing.T) {
	_, err := LoadContractsConfig(ContractsSource{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "required")
}

func TestContractOptionsFromEnv_Deprecated(t *testing.T) {
	// The deprecated function should still work
	opts, err := ContractOptionsFromEnv("config://testnet")
	require.NoError(t, err)
	require.NotEmpty(t, opts.AppChain.GroupMessageBroadcasterAddress)
}

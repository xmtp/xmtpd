// file: commands/admin_setup.go
package commands

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

func setupAppChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.IAppChainAdmin, error) {
	var (
		rpcURL     = viper.GetString("rpc-url")
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc-url is required")
	}
	if configFile == "" {
		return nil, fmt.Errorf("config-file is required")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("private-key is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	client, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	signer, err := blockchain.NewPrivateKeySigner(privateKey, contracts.AppChain.ChainID)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	// NewAppChainAdmin expects a ParameterAdmin inside; this remains internal to the admin.
	paramAdmin, err := blockchain.NewParameterAdmin(logger, client, signer, contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	return blockchain.NewAppChainAdmin(logger, client, signer, contracts, paramAdmin)
}

func setupSettlementChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.ISettlementChainAdmin, error) {
	var (
		rpcURL     = viper.GetString("rpc-url")
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc-url is required")
	}
	if configFile == "" {
		return nil, fmt.Errorf("config-file is required")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("private-key is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	client, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	signer, err := blockchain.NewPrivateKeySigner(privateKey, contracts.SettlementChain.ChainID)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	paramAdmin, err := blockchain.NewParameterAdmin(logger, client, signer, contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	return blockchain.NewSettlementChainAdmin(logger, client, signer, contracts, paramAdmin)
}

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

func setupNodeRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.INodeRegistryAdmin, error) {
	var (
		privateKey = viper.GetString("private-key")
		rpcURL     = viper.GetString("rpc-url")
		configFile = viper.GetString("config-file")
	)

	if privateKey == "" {
		return nil, fmt.Errorf("private key is required")
	}

	if rpcURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}

	if configFile == "" {
		return nil, fmt.Errorf("config file is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	chainClient, err := blockchain.NewRPCClient(
		ctx,
		rpcURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	parameterAdmin, err := blockchain.NewParameterAdmin(logger, chainClient, signer, contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		contracts,
		parameterAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry admin: %w", err)
	}

	return registryAdmin, nil
}

func setupNodeRegistryCaller(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.INodeRegistryCaller, error) {
	var (
		rpcURL     = viper.GetString("rpc-url")
		configFile = viper.GetString("config-file")
	)

	if rpcURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}

	if configFile == "" {
		return nil, fmt.Errorf("config file is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	chainClient, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("could not create chain client: %w", err)
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry caller: %w", err)
	}

	return caller, nil
}

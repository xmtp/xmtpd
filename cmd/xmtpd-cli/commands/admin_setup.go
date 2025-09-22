package commands

import (
	"context"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/fees"

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
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)

	rpcURL, err := resolveAppRPCURL()
	if err != nil {
		return nil, err
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
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)
	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
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
		configFile = viper.GetString("config-file")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}
	if privateKey == "" {
		return nil, fmt.Errorf("private key is required")
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
	configFile := viper.GetString("config-file")

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
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

func setupRateRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (*blockchain.RatesAdmin, error) {
	var (
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
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

	chainClient, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("could not create chain client: %w", err)
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	paramAdmin, err := blockchain.NewParameterAdmin(logger, chainClient, signer, contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	registryAdmin, err := blockchain.NewRatesAdmin(
		logger,
		paramAdmin,
		chainClient,
		contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry admin: %w", err)
	}

	return registryAdmin, nil
}

func setupRatesFetcher(
	ctx context.Context,
	logger *zap.Logger,
) (*fees.ContractRatesFetcher, error) {
	configFile := viper.GetString("config-file")

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}
	if configFile == "" {
		return nil, fmt.Errorf("config-file is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	chainClient, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("could not create chain client: %w", err)
	}

	fetcher, err := fees.NewContractRatesFetcher(ctx, chainClient, logger, contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create rates fetcher: %w", err)
	}

	return fetcher, nil
}

func setupFundsAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.IFundsAdmin, error) {
	var (
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
	)

	settlementRpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}
	appRpcURL, err := resolveAppRPCURL()
	if err != nil {
		return nil, err
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

	chainClientSettlement, err := blockchain.NewRPCClient(ctx, settlementRpcURL)
	if err != nil {
		return nil, fmt.Errorf("could not create chain client: %w", err)
	}

	signerSettlement, err := blockchain.NewPrivateKeySigner(
		privateKey,
		contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	chainClientApp, err := blockchain.NewRPCClient(ctx, appRpcURL)
	if err != nil {
		return nil, fmt.Errorf("could not create chain client: %w", err)
	}

	signerApp, err := blockchain.NewPrivateKeySigner(
		privateKey,
		contracts.AppChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	fundsAdmin, err := blockchain.NewFundsAdmin(
		blockchain.FundsAdminOpts{
			Logger:          logger,
			ContractOptions: contracts,
			Settlement: blockchain.FundsAdminSettlementOpts{
				Client: chainClientSettlement,
				Signer: signerSettlement,
			},
			App: blockchain.FundsAdminAppOpts{
				Client: chainClientApp,
				Signer: signerApp,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry admin: %w", err)
	}

	return fundsAdmin, nil
}

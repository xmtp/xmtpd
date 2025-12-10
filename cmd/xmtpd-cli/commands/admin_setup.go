package commands

import (
	"context"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/fees"

	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"go.uber.org/zap"
)

func setupAppChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.IParameterAdmin, blockchain.IAppChainAdmin, error) {
	var (
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
		env        = viper.GetString("environment")
	)

	rpcURL, err := resolveAppRPCURL()
	if err != nil {
		return nil, nil, err
	}

	if privateKey == "" {
		return nil, nil, fmt.Errorf("private-key is required")
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, nil, err
	}

	client, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, nil, err
	}
	signer, err := blockchain.NewPrivateKeySigner(privateKey, contracts.AppChain.ChainID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create signer: %w", err)
	}

	// NewAppChainAdmin expects a ParameterAdmin inside; this remains internal to the admin.
	paramAdmin, err := blockchain.NewAppChainParameterAdmin(logger, client, signer, contracts)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	appAdmin, err := blockchain.NewAppChainAdmin(logger, client, signer, contracts, paramAdmin)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create app admin: %w", err)
	}

	return paramAdmin, appAdmin, nil
}

func setupSettlementChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.IParameterAdmin, blockchain.ISettlementChainAdmin, error) {
	var (
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
		env        = viper.GetString("environment")
	)
	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, nil, err
	}

	if privateKey == "" {
		return nil, nil, fmt.Errorf("private-key is required")
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, nil, err
	}

	client, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		return nil, nil, err
	}
	signer, err := blockchain.NewPrivateKeySigner(privateKey, contracts.SettlementChain.ChainID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create signer: %w", err)
	}

	paramAdmin, err := blockchain.NewSettlementParameterAdmin(logger, client, signer, contracts)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	settlementAdmin, err := blockchain.NewSettlementChainAdmin(
		logger,
		client,
		signer,
		contracts,
		paramAdmin,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create settlement admin: %w", err)
	}

	return paramAdmin, settlementAdmin, nil
}

func setupNodeRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (blockchain.INodeRegistryAdmin, error) {
	var (
		privateKey = viper.GetString("private-key")
		configFile = viper.GetString("config-file")
		env        = viper.GetString("environment")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}

	if privateKey == "" {
		return nil, fmt.Errorf("private key is required")
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, err
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

	parameterAdmin, err := blockchain.NewSettlementParameterAdmin(
		logger,
		chainClient,
		signer,
		contracts,
	)
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
		configFile = viper.GetString("config-file")
		env        = viper.GetString("environment")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, err
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
) (blockchain.IRatesAdmin, error) {
	var (
		configFile = viper.GetString("config-file")
		privateKey = viper.GetString("private-key")
		env        = viper.GetString("environment")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}

	if privateKey == "" {
		return nil, fmt.Errorf("private-key is required")
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, err
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

	paramAdmin, err := blockchain.NewSettlementParameterAdmin(
		logger,
		chainClient,
		signer,
		contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	registryAdmin, err := blockchain.NewRatesAdmin(
		logger,
		chainClient,
		signer,
		paramAdmin,
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
	var (
		configFile = viper.GetString("config-file")
		env        = viper.GetString("environment")
	)

	rpcURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, err
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
		env        = viper.GetString("environment")
	)

	settlementRPCURL, err := resolveSettlementRPCURL()
	if err != nil {
		return nil, err
	}

	appRPCURL, err := resolveAppRPCURL()
	if err != nil {
		return nil, err
	}

	if privateKey == "" {
		return nil, fmt.Errorf("private-key is required")
	}

	contracts, err := resolveConfig(configFile, env)
	if err != nil {
		return nil, err
	}

	chainClientSettlement, err := blockchain.NewRPCClient(ctx, settlementRPCURL)
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

	chainClientApp, err := blockchain.NewRPCClient(ctx, appRPCURL)
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

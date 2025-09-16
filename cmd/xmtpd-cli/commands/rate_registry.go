package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
	"go.uber.org/zap"
)

func rateRegistryCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "rates",
		Short:        "Manage Rate Registry",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		addRatesCommand(),
		getRatesCommand(),
	)
	return &cmd
}

func addRatesCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:          "add",
		Short:        "Add rates to the rate registry",
		SilenceUsage: true,
		RunE:         addRatesHandler,
		Example: `
Usage: xmtpd-cli rates add --message-fee <message-fee> --storage-fee <storage-fee> --congestion-fee <congestion-fee> --target-rate <target-rate>

Example:
xmtpd-cli rates add --message-fee 1000000000000000000 --storage-fee 1000000000000000000 --congestion-fee 1000000000000000000 --target-rate 1000000000000000000
`,
	}

	cmd.PersistentFlags().
		Int64("message-fee", 0, "message fee to use")
	cmd.PersistentFlags().
		Int64("storage-fee", 0, "storage fee to use")
	cmd.PersistentFlags().
		Int64("congestion-fee", 0, "congestion fee to use")
	cmd.PersistentFlags().
		Uint64("target-rate", 0, "target rate to use")

	_ = cmd.MarkFlagRequired("message-fee")
	_ = cmd.MarkFlagRequired("storage-fee")
	_ = cmd.MarkFlagRequired("congestion-fee")
	_ = cmd.MarkFlagRequired("target-rate")

	return &cmd
}

func addRatesHandler(cmd *cobra.Command, _ []string) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	messageFee, err := cmd.Flags().GetInt64("message-fee")
	if err != nil {
		logger.Error("could not get message fee", zap.Error(err))
		return fmt.Errorf("could not get message fee: %w", err)
	}

	storageFee, err := cmd.Flags().GetInt64("storage-fee")
	if err != nil {
		logger.Error("could not get storage fee", zap.Error(err))
		return fmt.Errorf("could not get storage fee: %w", err)
	}

	congestionFee, err := cmd.Flags().GetInt64("congestion-fee")
	if err != nil {
		logger.Error("could not get congestion fee", zap.Error(err))
		return fmt.Errorf("could not get congestion fee: %w", err)
	}

	targetRate, err := cmd.Flags().GetUint64("target-rate")
	if err != nil {
		logger.Error("could not get target rate", zap.Error(err))
		return fmt.Errorf("could not get target rate: %w", err)
	}

	registryAdmin, err := setupRateRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup rate registry admin", zap.Error(err))
		return err
	}

	rates := fees.Rates{
		MessageFee:          currency.PicoDollar(messageFee),
		StorageFee:          currency.PicoDollar(storageFee),
		CongestionFee:       currency.PicoDollar(congestionFee),
		TargetRatePerMinute: targetRate,
	}

	if err := registryAdmin.AddRates(ctx, rates); err != nil {
		logger.Error("could not add rates to rate registry", zap.Error(err))
		return err
	}

	logger.Info("rates added to rate registry", zap.Any("rates", rates))
	return nil
}

func getRatesCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:          "get",
		Short:        "Get rates from the rate registry",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return getRatesHandler()
		},
		Example: `
Usage: xmtpd-cli rates get

Example:
xmtpd-cli rates get
`,
	}
	return &cmd
}

func getRatesHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	fetcher, err := setupRatesFetcher(ctx, logger)
	if err != nil {
		logger.Error("could not setup rates fetcher", zap.Error(err))
		return err
	}

	if err := fetcher.Start(); err != nil {
		if strings.Contains(err.Error(), "no rates found") {
			logger.Info("no rates found")
			return nil
		}
		logger.Error("could not start rates fetcher", zap.Error(err))
		return fmt.Errorf("could not start rates fetcher: %w", err)
	}

	rates, err := fetcher.GetRates(time.Now())
	if err != nil {
		logger.Error("could not get rates", zap.Error(err))
		return err
	}

	logger.Info("rates fetched successfully", zap.Any("rates", rates))
	return nil
}

func setupRateRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
) (*blockchain.RatesAdmin, error) {
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
	var (
		rpcURL     = viper.GetString("rpc-url")
		configFile = viper.GetString("config-file")
	)

	if rpcURL == "" {
		return nil, fmt.Errorf("rpc-url is required")
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

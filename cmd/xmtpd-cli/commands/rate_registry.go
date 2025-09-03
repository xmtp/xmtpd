package commands

import (
	"context"
	"fmt"
	"log"
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
		Use:   "rates",
		Short: "Manage Rate Registry",
	}

	cmd.AddCommand(
		addRatesCommand(),
		getRatesCommand(),
	)

	return &cmd
}

func addRatesCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add",
		Short: "Add rates to the rate registry",
		Run:   addRatesHandler,
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

func addRatesHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()

	messageFee, err := cmd.Flags().GetInt64("message-fee")
	if err != nil {
		logger.Fatal("could not get message fee", zap.Error(err))
	}

	storageFee, err := cmd.Flags().GetInt64("storage-fee")
	if err != nil {
		logger.Fatal("could not get storage fee", zap.Error(err))
	}

	congestionFee, err := cmd.Flags().GetInt64("congestion-fee")
	if err != nil {
		logger.Fatal("could not get congestion fee", zap.Error(err))
	}

	targetRate, err := cmd.Flags().GetUint64("target-rate")
	if err != nil {
		logger.Fatal("could not get target rate", zap.Error(err))
	}

	registryAdmin, err := setupRateRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup rate registry admin", zap.Error(err))
	}

	rates := fees.Rates{
		MessageFee:          currency.PicoDollar(messageFee),
		StorageFee:          currency.PicoDollar(storageFee),
		CongestionFee:       currency.PicoDollar(congestionFee),
		TargetRatePerMinute: targetRate,
	}

	err = registryAdmin.AddRates(ctx, rates)
	if err != nil {
		logger.Fatal("could not add rates to rate registry", zap.Error(err))
	}

	logger.Info("rates added to rate registry", zap.Any("rates", rates))
}

func getRatesCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get rates from the rate registry",
		Run:   getRatesHandler,
		Example: `
Usage: xmtpd-cli rates get

Example:
xmtpd-cli rates get
`,
	}

	return &cmd
}

func getRatesHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()

	fetcher, err := setupRatesFetcher(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup rates fetcher", zap.Error(err))
	}

	err = fetcher.Start()
	if err != nil {
		if strings.Contains(err.Error(), "no rates found") {
			logger.Info("no rates found")
			return
		}
		logger.Fatal("could not start rates fetcher", zap.Error(err))
	}

	rates, err := fetcher.GetRates(time.Now())
	if err != nil {
		logger.Fatal("could not get rates", zap.Error(err))
	}

	logger.Info("rates fetched successfully", zap.Any("rates", rates))
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
		return nil, err
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
		logger.Fatal("could not load config from file", zap.Error(err))
	}

	chainClient, err := blockchain.NewRPCClient(ctx, rpcURL)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	fetcher, err := fees.NewContractRatesFetcher(ctx, chainClient, logger, contracts)
	if err != nil {
		logger.Fatal("could not create rates fetcher", zap.Error(err))
	}

	return fetcher, nil
}

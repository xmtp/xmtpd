package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
	"go.uber.org/zap"
)

// ---- Options ----

type AddRatesOpts struct {
	MessageFee    int64
	StorageFee    int64
	CongestionFee int64
	TargetRate    uint64
}

// ---- Root ----

func rateRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rates",
		Short:        "Manage Rate Registry",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		addRatesCommand(),
		getRatesCommand(),
	)
	return cmd
}

// ---- add ----

func addRatesCommand() *cobra.Command {
	var opts AddRatesOpts

	cmd := &cobra.Command{
		Use:          "add",
		Short:        "Add rates to the rate registry",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return addRatesHandler(opts)
		},
		Example: `
Usage: xmtpd-cli rates add --message-fee <message-fee> --storage-fee <storage-fee> --congestion-fee <congestion-fee> --target-rate <target-rate>

Example:
xmtpd-cli rates add --message-fee 1000000000000000000 --storage-fee 1000000000000000000 --congestion-fee 1000000000000000000 --target-rate 1000000000000000000
`,
	}

	cmd.Flags().Int64Var(&opts.MessageFee, "message-fee", 0, "message fee to use")
	cmd.Flags().Int64Var(&opts.StorageFee, "storage-fee", 0, "storage fee to use")
	cmd.Flags().Int64Var(&opts.CongestionFee, "congestion-fee", 0, "congestion fee to use")
	cmd.Flags().Uint64Var(&opts.TargetRate, "target-rate", 0, "target rate to use")

	_ = cmd.MarkFlagRequired("message-fee")
	_ = cmd.MarkFlagRequired("storage-fee")
	_ = cmd.MarkFlagRequired("congestion-fee")
	_ = cmd.MarkFlagRequired("target-rate")

	return cmd
}

func addRatesHandler(opts AddRatesOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	registryAdmin, err := setupRateRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup rate registry admin", zap.Error(err))
		return err
	}

	rates := fees.Rates{
		MessageFee:          currency.PicoDollar(opts.MessageFee),
		StorageFee:          currency.PicoDollar(opts.StorageFee),
		CongestionFee:       currency.PicoDollar(opts.CongestionFee),
		TargetRatePerMinute: opts.TargetRate,
	}

	if err := registryAdmin.AddRates(ctx, rates); err != nil {
		logger.Error("could not add rates to rate registry", zap.Error(err))
		return err
	}

	logger.Info("rates added to rate registry", zap.Any("rates", rates))
	return nil
}

// ---- get ----

func getRatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get rates from the rate registry",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return getRatesHandler()
		},
		Example: `
Usage: xmtpd-cli rates get

Example:
xmtpd-cli rates get
`,
	}
	return cmd
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

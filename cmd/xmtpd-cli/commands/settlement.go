package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
	"go.uber.org/zap"
)

func settlementChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "settlement",
		Short:        "Manage Settlement Chain (gateways, payer registry, distribution manager, rate registry)",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		settlePauseCmd(),
		settleDMFeesRecipientCmd(),
		settleNodeAdminCmd(),
		settleMinDepositCmd(),
		settleWithdrawLockCmd(),
		settlePRMFeeRateCmd(),
		settleRateMigratorCmd(),
	)
	return cmd
}

// --- pause ---
func settlePauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pause",
		Short:        "Get/Update pause statuses on settlement chain",
		SilenceUsage: true,
	}
	cmd.AddCommand(settlePauseGetCmd(), settlePauseUpdateCmd())
	return cmd
}

func settlePauseGetCmd() *cobra.Command {
	var target options.Target
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get pause status for target: settlement-chain-gateway|payer-registry|distribution-manager",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settlePauseGetHandler(target)
		},
	}
	cmd.Flags().
		Var(&target, "target", "settlement-chain-gateway|payer-registry|distribution-manager")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func settlePauseGetHandler(target options.Target) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	switch target {
	case options.TargetSettlementChainGateway:
		p, e := admin.GetSettlementChainGatewayPauseStatus(ctx)
		if e != nil {
			logger.Error("read", zap.Error(e))
			return e
		}
		logger.Info("settlement-chain gateway pause", zap.Bool("paused", p))
	case options.TargetPayerRegistry:
		p, e := admin.GetPayerRegistryPauseStatus(ctx)
		if e != nil {
			logger.Error("read", zap.Error(e))
			return e
		}
		logger.Info("payer registry pause", zap.Bool("paused", p))
	case options.TargetDistributionManager:
		p, e := admin.GetDistributionManagerPauseStatus(ctx)
		if e != nil {
			logger.Error("read", zap.Error(e))
			return e
		}
		logger.Info("distribution manager pause", zap.Bool("paused", p))
	default:
		return errors.New(
			"target must be settlement-chain-gateway|payer-registry|distribution-manager",
		)
	}
	return nil
}

func settlePauseUpdateCmd() *cobra.Command {
	var target options.Target
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update pause status for target: settlement-chain-gateway|payer-registry|distribution-manager",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settlePauseUpdateHandler(target)
		},
	}
	cmd.Flags().
		Var(&target, "target", "settlement-chain-gateway|payer-registry|distribution-manager")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func settlePauseUpdateHandler(target options.Target) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	switch target {
	case options.TargetSettlementChainGateway:
		if err := admin.UpdateSettlementChainGatewayPauseStatus(ctx); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("settlement-chain gateway pause updated")
	case options.TargetPayerRegistry:
		if err := admin.UpdatePayerRegistryPauseStatus(ctx); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("payer registry pause updated")
	case options.TargetDistributionManager:
		if err := admin.UpdateDistributionManagerPauseStatus(ctx); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("distribution manager pause set updated")
	default:
		return errors.New(
			"target must be settlement-chain-gateway|payer-registry|distribution-manager",
		)
	}
	return nil
}

// --- DistributionManager: protocol fees recipient ---

func settleDMFeesRecipientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dm-protocol-fees-recipient",
		Short:        "Get/Update DistributionManager protocol fees recipient",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleDMFeesRecipientGetCmd(), settleDMFeesRecipientUpdateCmd())
	return cmd
}

func settleDMFeesRecipientGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get protocol fees recipient",
		SilenceUsage: true,
		RunE:         func(*cobra.Command, []string) error { return settleDMFeesRecipientGetHandler() },
	}
}

func settleDMFeesRecipientUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update protocol fees recipient",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleDMFeesRecipientUpdateHandler()
		},
	}
	return cmd
}

func settleDMFeesRecipientGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	addr, err := admin.GetDistributionManagerProtocolFeesRecipient(ctx)
	if err != nil {
		return err
	}
	logger.Info("distribution manager protocol fees recipient", zap.String("address", addr.Hex()))
	return nil
}

func settleDMFeesRecipientUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.UpdateDistributionManagerProtocolFeesRecipient(ctx); err != nil {
		return err
	}
	logger.Info(
		"distribution manager protocol fees recipient updated")
	return nil
}

// --- NodeRegistry: admin address ---

func settleNodeAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "node-admin",
		Short:        "Get/Update NodeRegistry admin address",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleNodeAdminGetCmd(), settleNodeAdminUpdateCmd())
	return cmd
}

func settleNodeAdminGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get node registry admin",
		SilenceUsage: true,
		RunE:         func(*cobra.Command, []string) error { return settleNodeAdminGetHandler() },
	}
}

func settleNodeAdminUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update node registry admin",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleNodeAdminUpdateHandler()
		},
	}
	return cmd
}

func settleNodeAdminGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	addr, err := admin.GetNodeRegistryAdmin(ctx)
	if err != nil {
		return err
	}
	logger.Info("node registry admin", zap.String("address", addr.Hex()))
	return nil
}

func settleNodeAdminUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.UpdateNodeRegistryAdmin(ctx); err != nil {
		return err
	}
	logger.Info("node registry admin updated")
	return nil
}

func settleMinDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "payer-min-deposit",
		Short:        "Get/Update PayerRegistry minimum deposit (uint96 microdollars)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleMinDepositGetCmd(), settleMinDepositUpdateCmd())
	return cmd
}

func settleMinDepositGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get minimum deposit (microdollars)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleMinDepositGetHandler()
		},
	}
}

func settleMinDepositGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	v, err := admin.GetPayerRegistryMinimumDeposit(ctx)
	if err != nil {
		return err
	}
	logger.Info("payer registry minimum deposit (microdollars)", zap.String("value", v.String()))
	return nil
}

func settleMinDepositUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update minimum deposit (microdollars, uint96)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleMinDepositUpdateHandler()
		},
	}
	return cmd
}

func settleMinDepositUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.UpdatePayerRegistryMinimumDeposit(ctx); err != nil {
		return err
	}
	logger.Info(
		"payer registry minimum deposit set (microdollars)",
	)
	return nil
}

// --- PayerRegistry: withdraw lock period (uint32 seconds) ---

func settleWithdrawLockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "payer-withdraw-lock",
		Short:        "Get/Update PayerRegistry withdraw lock period (seconds)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleWithdrawLockGetCmd(), settleWithdrawLockUpdateCmd())
	return cmd
}

func settleWithdrawLockGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get withdraw lock period (seconds)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleWithdrawLockGetHandler()
		},
	}
}

func settleWithdrawLockGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	secs, err := admin.GetPayerRegistryWithdrawLockPeriod(ctx)
	if err != nil {
		return err
	}
	logger.Info("payer registry withdraw lock period", zap.Uint32("seconds", secs))
	return nil
}

func settleWithdrawLockUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update withdraw lock period (seconds, uint32)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleWithdrawLockUpdateHandler()
		},
	}
	return cmd
}

func settleWithdrawLockUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.UpdatePayerRegistryWithdrawLockPeriod(ctx); err != nil {
		return err
	}
	logger.Info("payer registry withdraw lock period updated")
	return nil
}

// --- PayerReportManager: protocol fee rate (uint16 bps) ---

func settlePRMFeeRateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "prm-fee-rate",
		Short:        "Get/Update PayerReportManager protocol fee rate (bps, uint16)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settlePRMFeeRateGetCmd(), settlePRMFeeRateUpdateCmd())
	return cmd
}

func settlePRMFeeRateGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get PRM protocol fee rate (bps)",
		SilenceUsage: true,
		RunE:         func(*cobra.Command, []string) error { return settlePRMFeeRateGetHandler() },
	}
}

func settlePRMFeeRateUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update PRM protocol fee rate (bps, uint16)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settlePRMFeeRateUpdateHandler()
		},
	}
	return cmd
}

func settlePRMFeeRateGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return err
	}
	v, err := admin.GetPayerReportManagerProtocolFeeRate(ctx)
	if err != nil {
		logger.Error("read", zap.Error(err))
		return err
	}
	logger.Info("payer report manager fee rate (bps)", zap.Uint16("bps", v))
	return nil
}

func settlePRMFeeRateUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return err
	}
	if err := admin.UpdatePayerReportManagerProtocolFeeRate(ctx); err != nil {
		logger.Error("write", zap.Error(err))
		return err
	}
	logger.Info("payer report manager protocol fee rate updated")
	return nil
}

// --- RateRegistry: migrator address ---

func settleRateMigratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rate-migrator",
		Short:        "Get/Update RateRegistry migrator (address)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleRateMigratorGetCmd(), settleRateMigratorUpdateCmd())
	return cmd
}

func settleRateMigratorGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get RateRegistry migrator address",
		SilenceUsage: true,
		RunE:         func(*cobra.Command, []string) error { return settleRateMigratorGetHandler() },
	}
}

func settleRateMigratorUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set RateRegistry migrator address",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleRateMigratorUpdateHandler()
		},
	}
	return cmd
}

func settleRateMigratorGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup rates admin", zap.Error(err))
		return err
	}
	addr, perr := admin.GetRateRegistryMigrator(ctx)
	if perr != nil {
		logger.Error("read", zap.Error(perr))
		return perr
	}
	logger.Info("rate registry migrator", zap.String("address", addr.Hex()))
	return nil
}

func settleRateMigratorUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup rates admin", zap.Error(err))
		return err
	}
	if err := admin.UpdateRateRegistryMigrator(ctx); err != nil {
		logger.Error("write", zap.Error(err))
		return err
	}
	logger.Info("rate registry migrator updated")
	return nil
}

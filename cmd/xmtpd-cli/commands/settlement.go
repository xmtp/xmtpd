package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/xmtp/xmtpd/pkg/currency"

	"github.com/ethereum/go-ethereum/common"
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
		settleUnderlyingMintCmd(),
	)
	return cmd
}

// --- pause ---
func settlePauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pause",
		Short:        "Get/Set pause statuses on settlement chain",
		SilenceUsage: true,
	}
	cmd.AddCommand(settlePauseGetCmd(), settlePauseSetCmd())
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

	admin, err := setupSettlementChainAdmin(ctx, logger)
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
		return fmt.Errorf(
			"target must be settlement-chain-gateway|payer-registry|distribution-manager",
		)
	}
	return nil
}

func settlePauseSetCmd() *cobra.Command {
	var target options.Target
	var paused bool
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set pause status for target: settlement-chain-gateway|payer-registry|distribution-manager",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settlePauseSetHandler(target, paused)
		},
	}
	cmd.Flags().
		Var(&target, "target", "settlement-chain-gateway|payer-registry|distribution-manager")
	cmd.Flags().BoolVar(&paused, "paused", false, "pause status")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func settlePauseSetHandler(target options.Target, paused bool) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	switch target {
	case options.TargetSettlementChainGateway:
		if err := admin.SetSettlementChainGatewayPauseStatus(ctx, paused); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("settlement-chain gateway pause set", zap.Bool("paused", paused))
	case options.TargetPayerRegistry:
		if err := admin.SetPayerRegistryPauseStatus(ctx, paused); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("payer registry pause set", zap.Bool("paused", paused))
	case options.TargetDistributionManager:
		if err := admin.SetDistributionManagerPauseStatus(ctx, paused); err != nil {
			logger.Error("write", zap.Error(err))
			return err
		}
		logger.Info("distribution manager pause set", zap.Bool("paused", paused))
	default:
		return fmt.Errorf(
			"target must be settlement-chain-gateway|payer-registry|distribution-manager",
		)
	}
	return nil
}

// --- DistributionManager: protocol fees recipient ---

func settleDMFeesRecipientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dm-protocol-fees-recipient",
		Short:        "Get/Set DistributionManager protocol fees recipient",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleDMFeesRecipientGetCmd(), settleDMFeesRecipientSetCmd())
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

func settleDMFeesRecipientSetCmd() *cobra.Command {
	var recipient options.AddressFlag
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set protocol fees recipient",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleDMFeesRecipientSetHandler(recipient.Address)
		},
	}
	cmd.Flags().Var(&recipient, "address", "recipient address (hex)")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func settleDMFeesRecipientGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settleDMFeesRecipientSetHandler(addr common.Address) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.SetDistributionManagerProtocolFeesRecipient(ctx, addr); err != nil {
		return err
	}
	logger.Info(
		"distribution manager protocol fees recipient set",
		zap.String("address", addr.Hex()),
	)
	return nil
}

// --- NodeRegistry: admin address ---

func settleNodeAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "node-admin",
		Short:        "Get/Set NodeRegistry admin address",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleNodeAdminGetCmd(), settleNodeAdminSetCmd())
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

func settleNodeAdminSetCmd() *cobra.Command {
	var adminAddr options.AddressFlag
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set node registry admin",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleNodeAdminSetHandler(adminAddr.Address)
		},
	}
	cmd.Flags().Var(&adminAddr, "address", "admin address (hex)")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func settleNodeAdminGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settleNodeAdminSetHandler(addr common.Address) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.SetNodeRegistryAdmin(ctx, addr); err != nil {
		return err
	}
	logger.Info("node registry admin set", zap.String("address", addr.Hex()))
	return nil
}

// --- PayerRegistry: minimum deposit (uint96 microdollars) ---

func settleMinDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "payer-min-deposit",
		Short:        "Get/Set PayerRegistry minimum deposit (uint96 microdollars)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleMinDepositGetCmd(), settleMinDepositSetCmd())
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
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settleMinDepositSetCmd() *cobra.Command {
	var amountStr string
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set minimum deposit (microdollars, uint96)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			bi, ok := new(big.Int).SetString(amountStr, 10)
			if !ok {
				return fmt.Errorf("invalid --amount, must be base-10 integer")
			}
			return settleMinDepositSetHandler(bi)
		},
	}
	cmd.Flags().StringVar(&amountStr, "amount", "", "amount in microdollars (decimal string)")
	_ = cmd.MarkFlagRequired("amount")
	return cmd
}

func settleMinDepositSetHandler(amount *big.Int) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.SetPayerRegistryMinimumDeposit(ctx, amount); err != nil {
		return err
	}
	logger.Info(
		"payer registry minimum deposit set (microdollars)",
		zap.String("value", amount.String()),
	)
	return nil
}

// --- PayerRegistry: withdraw lock period (uint32 seconds) ---

func settleWithdrawLockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "payer-withdraw-lock",
		Short:        "Get/Set PayerRegistry withdraw lock period (seconds)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleWithdrawLockGetCmd(), settleWithdrawLockSetCmd())
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
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settleWithdrawLockSetCmd() *cobra.Command {
	var secs uint32
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set withdraw lock period (seconds, uint32)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleWithdrawLockSetHandler(secs)
		},
	}
	cmd.Flags().Uint32Var(&secs, "seconds", 0, "seconds")
	_ = cmd.MarkFlagRequired("seconds")
	return cmd
}

func settleWithdrawLockSetHandler(secs uint32) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		return err
	}
	if err := admin.SetPayerRegistryWithdrawLockPeriod(ctx, secs); err != nil {
		return err
	}
	logger.Info("payer registry withdraw lock period set", zap.Uint32("seconds", secs))
	return nil
}

// --- PayerReportManager: protocol fee rate (uint16 bps) ---

func settlePRMFeeRateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "prm-fee-rate",
		Short:        "Get/Set PayerReportManager protocol fee rate (bps, uint16)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settlePRMFeeRateGetCmd(), settlePRMFeeRateSetCmd())
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

func settlePRMFeeRateSetCmd() *cobra.Command {
	var bps uint16
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set PRM protocol fee rate (bps, uint16)",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settlePRMFeeRateSetHandler(bps)
		},
	}
	cmd.Flags().Uint16Var(&bps, "bps", 0, "basis points (0..65535)")
	_ = cmd.MarkFlagRequired("bps")
	return cmd
}

func settlePRMFeeRateGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settlePRMFeeRateSetHandler(bps uint16) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return err
	}
	if err := admin.SetPayerReportManagerProtocolFeeRate(ctx, bps); err != nil {
		logger.Error("write", zap.Error(err))
		return err
	}
	logger.Info("payer report manager protocol fee rate set", zap.Uint16("bps", bps))
	return nil
}

// --- RateRegistry: migrator address ---

func settleRateMigratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rate-migrator",
		Short:        "Get/Set RateRegistry migrator (address)",
		SilenceUsage: true,
	}
	cmd.AddCommand(settleRateMigratorGetCmd(), settleRateMigratorSetCmd())
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

func settleRateMigratorSetCmd() *cobra.Command {
	var migratorAddr options.AddressFlag
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set RateRegistry migrator address",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return settleRateMigratorSetHandler(migratorAddr.Address)
		},
	}
	cmd.Flags().Var(&migratorAddr, "address", "migrator address (hex)")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func settleRateMigratorGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
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

func settleRateMigratorSetHandler(addr common.Address) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup rates admin", zap.Error(err))
		return err
	}
	if err := admin.SetRateRegistryMigrator(ctx, addr); err != nil {
		logger.Error("write", zap.Error(err))
		return err
	}
	logger.Info("rate registry migrator set", zap.String("address", addr.Hex()))
	return nil
}

func settleUnderlyingMintCmd() *cobra.Command {
	var toAddr options.AddressFlag
	var amountHuman string
	var raw bool
	cmd := &cobra.Command{
		Use:          "underlying-mint",
		Hidden:       true,
		Short:        "Mint mock underlying fee token to an address (max 10000 tokens)",
		SilenceUsage: true,
		Example: `
# Mint 1000 tokens to 0xabc... using the token's own decimals
xmtpd-cli settlement underlying-mint \
  --to 0xRecipient \
  --amount 1000

# If you already have the raw uint256 (no scaling), use --raw
xmtpd-cli settlement underlying-mint \
  --to 0xRecipient \
  --amount 1000000000 --raw
`,
		RunE: func(*cobra.Command, []string) error {
			return settleUnderlyingMintHandler(toAddr.Address, amountHuman, raw)
		},
	}
	cmd.Flags().Var(&toAddr, "to", "recipient address (hex)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().
		StringVar(&amountHuman, "amount", "", "amount; decimal string in token units unless --raw")
	_ = cmd.MarkFlagRequired("amount")
	cmd.Flags().
		BoolVar(&raw, "raw", false, "interpret --amount as raw uint256 (no decimals scaling)")

	return cmd
}

func settleUnderlyingMintHandler(to common.Address, amountStr string, raw bool) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	// Parse amount
	var amount *big.Int
	if raw {
		ai, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("invalid --amount (raw uint256) %q", amountStr)
		}
		amount = ai
	} else {
		// interpret --amount as whole tokens and scale to micro (6 decimals)
		ai, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("invalid --amount (decimal string) %q", amountStr)
		}
		scaled := new(big.Int).Mul(ai, big.NewInt(currency.MicroDollarsPerDollar))
		amount = scaled
	}

	if err := admin.MintMockUSDC(ctx, to, amount); err != nil {
		logger.Error("mint mock underlying fee token", zap.Error(err))
		return err
	}

	logger.Info("successfully minted mock underlying fee token",
		zap.String("to", to.Hex()),
		zap.String("amountRaw", amount.String()),
	)

	return nil
}

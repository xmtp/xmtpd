// file: commands/settlement.go
package commands

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func settlementChainCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "settlement",
		Short: "Manage Settlement Chain (gateways, payer registry, distribution manager, rate registry)",
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
	return &cmd
}

// --- pause ---

func settlePauseCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "pause",
		Short: "Get/Set pause statuses on settlement chain",
	}
	cmd.AddCommand(settlePauseGetCmd(), settlePauseSetCmd())
	return &cmd
}

func settlePauseGetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get pause status for target: gateway|payer-registry|distribution-manager",
		Run:   settlePauseGetHandler,
	}
	cmd.Flags().String("target", "", "gateway|payer-registry|distribution-manager")
	_ = cmd.MarkFlagRequired("target")
	return &cmd
}

func settlePauseGetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	ctx := context.Background()

	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	switch target {
	case "gateway":
		p, e := admin.GetSettlementChainGatewayPauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("settlement gateway pause", zap.Bool("paused", p))
	case "payer-registry":
		p, e := admin.GetPayerRegistryPauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("payer registry pause", zap.Bool("paused", p))
	case "distribution-manager":
		p, e := admin.GetDistributionManagerPauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("distribution manager pause", zap.Bool("paused", p))
	default:
		logger.Fatal("target must be gateway|payer-registry|distribution-manager")
	}
}

func settlePauseSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set pause status for target: gateway|payer-registry|distribution-manager",
		Run:   settlePauseSetHandler,
	}
	cmd.Flags().String("target", "", "gateway|payer-registry|distribution-manager")
	cmd.Flags().Bool("paused", false, "pause status")
	_ = cmd.MarkFlagRequired("target")
	return &cmd
}

func settlePauseSetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	paused, _ := cmd.Flags().GetBool("paused")
	ctx := context.Background()

	admin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	switch target {
	case "gateway":
		if err := admin.SetSettlementChainGatewayPauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("settlement gateway pause set", zap.Bool("paused", paused))
	case "payer-registry":
		if err := admin.SetPayerRegistryPauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("payer registry pause set", zap.Bool("paused", paused))
	case "distribution-manager":
		if err := admin.SetDistributionManagerPauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("distribution manager pause set", zap.Bool("paused", paused))
	default:
		logger.Fatal("target must be gateway|payer-registry|distribution-manager")
	}
}

// --- DistributionManager: protocol fees recipient ---

func settleDMFeesRecipientCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "dm-protocol-fees-recipient",
		Short: "Get/Set DistributionManager protocol fees recipient",
	}
	cmd.AddCommand(settleDMFeesRecipientGetCmd(), settleDMFeesRecipientSetCmd())
	return &cmd
}

func settleDMFeesRecipientGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get protocol fees recipient",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			addr, err := admin.GetDistributionManagerProtocolFeesRecipient(ctx)
			if err != nil {
				logger.Fatal("read", zap.Error(err))
			}
			logger.Info(
				"distribution manager protocol fees recipient",
				zap.String("address", addr.Hex()),
			)
		},
	}
}

func settleDMFeesRecipientSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set protocol fees recipient",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			addrStr, _ := cmd.Flags().GetString("address")
			if !common.IsHexAddress(addrStr) {
				logger.Fatal("invalid address")
			}
			addr := common.HexToAddress(addrStr)

			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			if err := admin.SetDistributionManagerProtocolFeesRecipient(ctx, addr); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info(
				"distribution manager protocol fees recipient set",
				zap.String("address", addr.Hex()),
			)
		},
	}
	cmd.Flags().String("address", "", "recipient address")
	_ = cmd.MarkFlagRequired("address")
	return &cmd
}

// --- NodeRegistry: admin address ---

func settleNodeAdminCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "node-admin",
		Short: "Get/Set NodeRegistry admin address",
	}
	cmd.AddCommand(settleNodeAdminGetCmd(), settleNodeAdminSetCmd())
	return &cmd
}

func settleNodeAdminGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get node registry admin",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			addr, err := admin.GetNodeRegistryAdmin(ctx)
			if err != nil {
				logger.Fatal("read", zap.Error(err))
			}
			logger.Info("node registry admin", zap.String("address", addr.Hex()))
		},
	}
}

func settleNodeAdminSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set node registry admin",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			addrStr, _ := cmd.Flags().GetString("address")
			if !common.IsHexAddress(addrStr) {
				logger.Fatal("invalid address")
			}
			addr := common.HexToAddress(addrStr)

			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			if err := admin.SetNodeRegistryAdmin(ctx, addr); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info("node registry admin set", zap.String("address", addr.Hex()))
		},
	}
	cmd.Flags().String("address", "", "admin address")
	_ = cmd.MarkFlagRequired("address")
	return &cmd
}

// --- PayerRegistry: minimum deposit (uint96 microdollars) ---

func settleMinDepositCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "payer-min-deposit",
		Short: "Get/Set PayerRegistry minimum deposit (uint96 microdollars)",
	}
	cmd.AddCommand(settleMinDepositGetCmd(), settleMinDepositSetCmd())
	return &cmd
}

func settleMinDepositGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get minimum deposit (microdollars)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			v, err := admin.GetPayerRegistryMinimumDeposit(ctx)
			if err != nil {
				logger.Fatal("read", zap.Error(err))
			}
			logger.Info(
				"payer registry minimum deposit (microdollars)",
				zap.String("value", v.String()),
			)
		},
	}
}

func settleMinDepositSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set minimum deposit (microdollars, uint96)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			amountStr, _ := cmd.Flags().GetString("amount")
			bi, ok := new(big.Int).SetString(amountStr, 10)
			if !ok {
				logger.Fatal("invalid --amount, must be base-10 integer")
			}

			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			if err := admin.SetPayerRegistryMinimumDeposit(ctx, bi); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info(
				"payer registry minimum deposit set (microdollars)",
				zap.String("value", bi.String()),
			)
		},
	}
	cmd.Flags().String("amount", "", "amount in microdollars (decimal string)")
	_ = cmd.MarkFlagRequired("amount")
	return &cmd
}

// --- PayerRegistry: withdraw lock period (uint32 seconds) ---

func settleWithdrawLockCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "payer-withdraw-lock",
		Short: "Get/Set PayerRegistry withdraw lock period (seconds)",
	}
	cmd.AddCommand(settleWithdrawLockGetCmd(), settleWithdrawLockSetCmd())
	return &cmd
}

func settleWithdrawLockGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get withdraw lock period (seconds)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			secs, err := admin.GetPayerRegistryWithdrawLockPeriod(ctx)
			if err != nil {
				logger.Fatal("read", zap.Error(err))
			}
			logger.Info("payer registry withdraw lock period", zap.Uint32("seconds", secs))
		},
	}
}

func settleWithdrawLockSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set withdraw lock period (seconds, uint32)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			secs, _ := cmd.Flags().GetUint32("seconds")
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			if err := admin.SetPayerRegistryWithdrawLockPeriod(ctx, secs); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info("payer registry withdraw lock period set", zap.Uint32("seconds", secs))
		},
	}
	cmd.Flags().Uint32("seconds", 0, "seconds")
	_ = cmd.MarkFlagRequired("seconds")
	return &cmd
}

// --- PayerReportManager: protocol fee rate (uint16 bps) ---

func settlePRMFeeRateCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "prm-fee-rate",
		Short: "Get/Set PayerReportManager protocol fee rate (bps, uint16)",
	}
	cmd.AddCommand(settlePRMFeeRateGetCmd(), settlePRMFeeRateSetCmd())
	return &cmd
}

func settlePRMFeeRateGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get PRM protocol fee rate (bps)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			v, err := admin.GetPayerReportManagerProtocolFeeRate(ctx)
			if err != nil {
				logger.Fatal("read", zap.Error(err))
			}
			logger.Info("payer report manager fee rate (bps)", zap.Uint16("bps", v))
		},
	}
}

func settlePRMFeeRateSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set PRM protocol fee rate (bps, uint16)",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			bps, _ := cmd.Flags().GetUint16("bps")
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}
			if err := admin.SetPayerReportManagerProtocolFeeRate(ctx, bps); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info("payer report manager protocol fee rate set", zap.Uint16("bps", bps))
		},
	}
	cmd.Flags().Uint16("bps", 0, "basis points (0..65535)")
	_ = cmd.MarkFlagRequired("bps")
	return &cmd
}

// --- RateRegistry: migrator address ---

func settleRateMigratorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "rate-migrator",
		Short: "Get/Set RateRegistry migrator (address)",
	}
	cmd.AddCommand(settleRateMigratorGetCmd(), settleRateMigratorSetCmd())
	return &cmd
}

func settleRateMigratorGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get RateRegistry migrator address",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup rates admin", zap.Error(err))
			}
			addr, perr := admin.GetRateRegistryMigrator(ctx)
			if perr != nil {
				logger.Fatal("read", zap.Error(perr))
			}
			logger.Info("rate registry migrator", zap.String("address", addr.Hex()))
		},
	}
}

func settleRateMigratorSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set RateRegistry migrator address",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			addrStr, _ := cmd.Flags().GetString("address")
			if !common.IsHexAddress(addrStr) {
				logger.Fatal("invalid address")
			}
			addr := common.HexToAddress(addrStr)
			ctx := context.Background()
			admin, err := setupSettlementChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup rates admin", zap.Error(err))
			}
			if err := admin.SetRateRegistryMigrator(ctx, addr); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
			logger.Info("rate registry migrator set", zap.String("address", addr.Hex()))
		},
	}
	cmd.Flags().String("address", "", "migrator address")
	_ = cmd.MarkFlagRequired("address")
	return &cmd
}

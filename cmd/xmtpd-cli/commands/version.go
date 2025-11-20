package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
	"go.uber.org/zap"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Get version of contract",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		getVersionCmd(),
	)
	return cmd
}

// --- pause ---
func getVersionCmd() *cobra.Command {
	var target options.Target
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "todo",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return getVersionHandler(target)
		},
	}
	cmd.Flags().
		Var(&target, "target", "settlement-chain-gateway|payer-registry|distribution-manager")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func getVersionHandler(target options.Target) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	_, settlementAdmin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	_, appAdmin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	switch target {
	case options.TargetSettlementChainGateway:
		version, err := settlementAdmin.GetSettlementChainGatewayVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("settlement-chain gateway version", zap.String("version", version))

	case options.TargetDistributionManager:
		version, err := settlementAdmin.GetDistributionManagerVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("distribution manager version", zap.String("version", version))

	case options.TargetPayerRegistry:
		version, err := settlementAdmin.GetPayerRegistryVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("payer registry version", zap.String("version", version))

	case options.TargetSettlementParameterRegistry:
		version, err := settlementAdmin.GetSettlementParameterRegistryVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("settlement-chain parameter registry version", zap.String("version", version))

	case options.TargetPayerReportManager:
		version, err := settlementAdmin.GetPayerReportManagerVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("payer report manager version", zap.String("version", version))

	case options.TargetRateRegistry:
		version, err := settlementAdmin.GetRateRegistryVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("rate registry version", zap.String("version", version))

	case options.TargetGroup:
		version, err := appAdmin.GetGroupMessageBroadcasterVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("group message broadcaster version", zap.String("version", version))

	case options.TargetIdentity:
		version, err := appAdmin.GetIdentityUpdateBroadcasterVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("identity update broadcaster version", zap.String("version", version))

	case options.TargetAppChainGateway:
		version, err := appAdmin.GetAppChainGatewayVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("app-chain gateway version", zap.String("version", version))

	case options.TargetAppParameterRegistry:
		version, err := appAdmin.GetAppParameterRegistryVersion(ctx)
		if err != nil {
			logger.Error("getting version", zap.Error(err))
			return err
		}
		logger.Info("app-chain parameter registry version", zap.String("version", version))
	default:
		return fmt.Errorf(
			"unknown target",
		)
	}
	return nil
}

//func (s settlementChainAdmin) GetSettlementChainGatewayVersion(
//	ctx context.Context,
//) (string, error) {
//	return s.settlementChainGateway.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetSettlementParameterRegistryVersion(
//	ctx context.Context,
//) (string, error) {
//	return s.settlementParameterRegistry.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetPayerReportManagerVersion(ctx context.Context) (string, error) {
//	return s.payerReportManager.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetRateRegistryVersion(ctx context.Context) (string, error) {
//	return s.rateRegistry.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetPayerRegistryVersion(ctx context.Context) (string, error) {
//	return s.payerRegistry.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetNodeRegistryVersion(ctx context.Context) (string, error) {
//	return s.nodeRegistry.Version(&bind.CallOpts{Context: ctx})
//}
//
//func (s settlementChainAdmin) GetDistributionManagerVersion(ctx context.Context) (string, error) {
//	return s.distributionManager.Version(&bind.CallOpts{Context: ctx})
//}

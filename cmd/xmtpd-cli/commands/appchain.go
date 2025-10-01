package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
	"go.uber.org/zap"
)

func appChainCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "appchain",
		Short:        "Manage App Chain (broadcasters & app-chain gateway)",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		appPauseCmd(),
		appBootstrapperCmd(),
		appPayloadSizeCmd(),
	)
	return &cmd
}

// --- pause ---

func appPauseCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "pause",
		Short:        "Get/Update app-chain pause statuses",
		SilenceUsage: true,
	}
	cmd.AddCommand(appPauseGetCmd(), appPauseUpdateCmd())
	return &cmd
}

func appPauseGetCmd() *cobra.Command {
	var target options.Target

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get pause status for target: identity|group|app-chain-gateway",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPauseGetHandler(target)
		},
	}

	cmd.Flags().Var(&target, "target", "identity|group|app-chain-gateway")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}

func appPauseGetHandler(target options.Target) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()

	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup appchain admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		p, e := admin.GetIdentityUpdatePauseStatus(ctx)
		if e != nil {
			logger.Error("read identity pause", zap.Error(e))
			return e
		}
		logger.Info("identity broadcaster pause", zap.Bool("paused", p))
	case options.TargetGroup:
		p, e := admin.GetGroupMessagePauseStatus(ctx)
		if e != nil {
			logger.Error("read group pause", zap.Error(e))
			return e
		}
		logger.Info("group broadcaster pause", zap.Bool("paused", p))
	case options.TargetAppChainGateway:
		p, e := admin.GetAppChainGatewayPauseStatus(ctx)
		if e != nil {
			logger.Error("read gateway pause", zap.Error(e))
			return e
		}
		logger.Info("app-chain gateway pause", zap.Bool("paused", p))
	default:
		return fmt.Errorf("target must be identity|group|app-chain-gateway")
	}

	return nil
}

func appPauseUpdateCmd() *cobra.Command {
	var target options.Target

	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update pause status for target: identity|group|app-chain-gateway",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPauseUpdateHandler(target)
		},
	}

	cmd.Flags().Var(&target, "target", "identity|group|app-chain-gateway")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func appPauseUpdateHandler(target options.Target) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		if err := admin.UpdateIdentityUpdatePauseStatus(ctx); err != nil {
			logger.Error("write identity pause", zap.Error(err))
			return err
		}
		logger.Info("identity broadcaster pause updated")
	case options.TargetGroup:
		if err := admin.UpdateGroupMessagePauseStatus(ctx); err != nil {
			logger.Error("write group pause", zap.Error(err))
			return err
		}
		logger.Info("group broadcaster pause updated")
	case options.TargetAppChainGateway:
		if err := admin.UpdateAppChainGatewayPauseStatus(ctx); err != nil {
			logger.Error("write gateway pause", zap.Error(err))
			return err
		}
		logger.Info("app-chain gateway pause updated")
	default:
		return fmt.Errorf("target must be identity|group|app-chain-gateway")
	}

	return nil
}

// --- bootstrapper ---

func appBootstrapperCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "bootstrapper",
		Short:        "Get/Update payload bootstrapper (identity & group)",
		SilenceUsage: true,
	}
	cmd.AddCommand(appBootstrapperGetCmd(), appBootstrapperUpdateCmd())
	return &cmd
}

func appBootstrapperGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "get",
		Short:        "Get bootstrapper addresses for identity & group",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appBootstrapperGetHandler()
		},
	}
}

func appBootstrapperGetHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()
	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	iu, err := admin.GetIdentityUpdateBootstrapper(ctx)
	if err != nil {
		logger.Error("read identity bootstrapper", zap.Error(err))
		return err
	}
	gm, err := admin.GetGroupMessageBootstrapper(ctx)
	if err != nil {
		logger.Error("read group bootstrapper", zap.Error(err))
		return err
	}

	logger.Info("bootstrapper",
		zap.String("identity", iu.Hex()),
		zap.String("group", gm.Hex()),
	)
	if iu != gm {
		logger.Warn("identity and group bootstrappers differ")
	}
	return nil
}

func appBootstrapperUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update bootstrapper address for BOTH identity & group",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appBootstrapperUpdateHandler()
		},
	}
	return cmd
}

func appBootstrapperUpdateHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	if err := admin.UpdateIdentityUpdateBootstrapper(ctx); err != nil {
		logger.Error("update identity bootstrapper", zap.Error(err))
		return err
	}
	if err := admin.UpdateGroupMessageBootstrapper(ctx); err != nil {
		logger.Error("update group bootstrapper", zap.Error(err))
		return err
	}
	logger.Info("bootstrapper updated")
	return nil
}

// --- payload-size ---

func appPayloadSizeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "payload-size",
		Short:        "Get/Update payload size bounds for broadcasters",
		SilenceUsage: true,
	}
	cmd.AddCommand(appPayloadSizeGetCmd(), appPayloadSizeUpdateCmd())
	return &cmd
}

func appPayloadSizeGetCmd() *cobra.Command {
	var target options.Target
	var bound options.PayloadBound

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get payload size for --target identity|group and --bound min|max",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPayloadSizeGetHandler(target, bound)
		},
	}
	cmd.Flags().Var(&target, "target", "identity|group")
	cmd.Flags().Var(&bound, "bound", "min|max")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("bound")
	return cmd
}

func appPayloadSizeGetHandler(target options.Target, bound options.PayloadBound) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		switch bound {
		case options.PayloadMin:
			v, e := admin.GetIdentityUpdateMinPayloadSize(ctx)
			if e != nil {
				logger.Error("read", zap.Error(e))
				return e
			}
			logger.Info(
				"payload size",
				zap.String("target", "identity"),
				zap.String("bound", "min"),
				zap.Uint32("bytes", v),
			)
		case options.PayloadMax:
			v, e := admin.GetIdentityUpdateMaxPayloadSize(ctx)
			if e != nil {
				logger.Error("read", zap.Error(e))
				return e
			}
			logger.Info(
				"payload size",
				zap.String("target", "identity"),
				zap.String("bound", "max"),
				zap.Uint32("bytes", v),
			)
		}
	case options.TargetGroup:
		switch bound {
		case options.PayloadMin:
			v, e := admin.GetGroupMessageMinPayloadSize(ctx)
			if e != nil {
				logger.Error("read", zap.Error(e))
				return e
			}
			logger.Info(
				"payload size",
				zap.String("target", "group"),
				zap.String("bound", "min"),
				zap.Uint32("bytes", v),
			)
		case options.PayloadMax:
			v, e := admin.GetGroupMessageMaxPayloadSize(ctx)
			if e != nil {
				logger.Error("read", zap.Error(e))
				return e
			}
			logger.Info(
				"payload size",
				zap.String("target", "group"),
				zap.String("bound", "max"),
				zap.Uint32("bytes", v),
			)
		}
	default:
		return fmt.Errorf("target must be identity|group")
	}
	return nil
}

func appPayloadSizeUpdateCmd() *cobra.Command {
	var target options.Target
	var bound options.PayloadBound

	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Update payload size for --target identity|group and --bound min|max",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPayloadSizeUpdateHandler(target, bound)
		},
	}
	cmd.Flags().Var(&target, "target", "identity|group")
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().Var(&bound, "bound", "min|max")
	_ = cmd.MarkFlagRequired("bound")
	return cmd
}

func appPayloadSizeUpdateHandler(
	target options.Target,
	bound options.PayloadBound,
) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	_, admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		if bound == options.PayloadMin {
			if err := admin.UpdateIdentityUpdateMinPayloadSize(ctx); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		} else {
			if err := admin.UpdateIdentityUpdateMaxPayloadSize(ctx); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		}
	case options.TargetGroup:
		if bound == options.PayloadMin {
			if err := admin.UpdateGroupMessageMinPayloadSize(ctx); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		} else {
			if err := admin.UpdateGroupMessageMaxPayloadSize(ctx); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		}
	default:
		return fmt.Errorf("target must be identity|group")
	}

	logger.Info("payload size updated",
		zap.String("target", string(target)),
		zap.String("bound", string(bound)),
	)
	return nil
}

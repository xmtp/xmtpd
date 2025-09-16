package commands

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

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
		Short:        "Get/Set app-chain pause statuses",
		SilenceUsage: true,
	}
	cmd.AddCommand(appPauseGetCmd(), appPauseSetCmd())
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

	admin, err := setupAppChainAdmin(ctx, logger)
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

func appPauseSetCmd() *cobra.Command {
	var target options.Target
	var paused bool

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set pause status for target: identity|group|app-chain-gateway",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPauseSetHandler(target, paused)
		},
	}

	cmd.Flags().Var(&target, "target", "identity|group|app-chain-gateway")
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().BoolVar(&paused, "paused", false, "pause status (true|false)")
	_ = cmd.MarkFlagRequired("paused")
	return cmd
}

func appPauseSetHandler(target options.Target, paused bool) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		if err := admin.SetIdentityUpdatePauseStatus(ctx, paused); err != nil {
			logger.Error("write identity pause", zap.Error(err))
			return err
		}
		logger.Info("identity broadcaster pause set", zap.Bool("paused", paused))
	case options.TargetGroup:
		if err := admin.SetGroupMessagePauseStatus(ctx, paused); err != nil {
			logger.Error("write group pause", zap.Error(err))
			return err
		}
		logger.Info("group broadcaster pause set", zap.Bool("paused", paused))
	case options.TargetAppChainGateway:
		if err := admin.SetAppChainGatewayPauseStatus(ctx, paused); err != nil {
			logger.Error("write gateway pause", zap.Error(err))
			return err
		}
		logger.Info("app-chain gateway pause set", zap.Bool("paused", paused))
	default:
		return fmt.Errorf("target must be identity|group|app-chain-gateway")
	}

	return nil
}

// --- bootstrapper ---

func appBootstrapperCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "bootstrapper",
		Short:        "Get/Set payload bootstrapper (identity & group)",
		SilenceUsage: true,
	}
	cmd.AddCommand(appBootstrapperGetCmd(), appBootstrapperSetCmd())
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
	admin, err := setupAppChainAdmin(ctx, logger)
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

func appBootstrapperSetCmd() *cobra.Command {
	var addr options.AddressFlag

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set bootstrapper address for BOTH identity & group",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appBootstrapperSetHandler(addr.Address)
		},
	}

	cmd.Flags().Var(&addr, "address", "bootstrapper address (checksummed hex)")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func appBootstrapperSetHandler(addr common.Address) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	if err := admin.SetIdentityUpdateBootstrapper(ctx, addr); err != nil {
		logger.Error("set identity bootstrapper", zap.Error(err))
		return err
	}
	if err := admin.SetGroupMessageBootstrapper(ctx, addr); err != nil {
		logger.Error("set group bootstrapper", zap.Error(err))
		return err
	}
	logger.Info("bootstrapper set", zap.String("address", addr.String()))
	return nil
}

// --- payload-size ---

func appPayloadSizeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:          "payload-size",
		Short:        "Get/Set payload size bounds for broadcasters",
		SilenceUsage: true,
	}
	cmd.AddCommand(appPayloadSizeGetCmd(), appPayloadSizeSetCmd())
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

	admin, err := setupAppChainAdmin(ctx, logger)
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
				zap.Uint64("bytes", v),
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
				zap.Uint64("bytes", v),
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
				zap.Uint64("bytes", v),
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
				zap.Uint64("bytes", v),
			)
		}
	default:
		return fmt.Errorf("target must be identity|group")
	}
	return nil
}

func appPayloadSizeSetCmd() *cobra.Command {
	var target options.Target
	var bound options.PayloadBound
	var size uint64

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set payload size for --target identity|group and --bound min|max",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appPayloadSizeSetHandler(target, bound, size)
		},
	}
	cmd.Flags().Var(&target, "target", "identity|group")
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().Var(&bound, "bound", "min|max")
	_ = cmd.MarkFlagRequired("bound")
	cmd.Flags().Uint64Var(&size, "size", 0, "size in bytes")
	_ = cmd.MarkFlagRequired("size")
	return cmd
}

func appPayloadSizeSetHandler(
	target options.Target,
	bound options.PayloadBound,
	size uint64,
) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()
	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("setup admin", zap.Error(err))
		return fmt.Errorf("could not setup appchain admin: %w", err)
	}

	switch target {
	case options.TargetIdentity:
		if bound == options.PayloadMin {
			if err := admin.SetIdentityUpdateMinPayloadSize(ctx, size); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		} else {
			if err := admin.SetIdentityUpdateMaxPayloadSize(ctx, size); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		}
	case options.TargetGroup:
		if bound == options.PayloadMin {
			if err := admin.SetGroupMessageMinPayloadSize(ctx, size); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		} else {
			if err := admin.SetGroupMessageMaxPayloadSize(ctx, size); err != nil {
				logger.Error("write", zap.Error(err))
				return err
			}
		}
	default:
		return fmt.Errorf("target must be identity|group")
	}

	logger.Info("payload size set",
		zap.String("target", string(target)),
		zap.String("bound", string(bound)),
		zap.Uint64("bytes", size),
	)
	return nil
}

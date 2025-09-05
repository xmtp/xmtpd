package commands

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func appChainCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "appchain",
		Short: "Manage App Chain (broadcasters & app-chain gateway)",
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
		Use:   "pause",
		Short: "Get/Set app-chain pause statuses",
	}
	cmd.AddCommand(appPauseGetCmd(), appPauseSetCmd())
	return &cmd
}

func appPauseGetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get pause status for target: identity|group|gateway",
		Run:   appPauseGetHandler,
	}
	cmd.Flags().String("target", "", "identity|group|gateway")
	_ = cmd.MarkFlagRequired("target")
	return &cmd
}

func appPauseGetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	ctx := context.Background()

	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}

	switch target {
	case "identity":
		p, e := admin.GetIdentityUpdatePauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("identity broadcaster pause", zap.Bool("paused", p))
	case "group":
		p, e := admin.GetGroupMessagePauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("group broadcaster pause", zap.Bool("paused", p))
	case "gateway":
		p, e := admin.GetAppChainGatewayPauseStatus(ctx)
		if e != nil {
			logger.Fatal("read", zap.Error(e))
		}
		logger.Info("app-chain gateway pause", zap.Bool("paused", p))
	default:
		logger.Fatal("target must be identity|group|gateway")
	}
}

func appPauseSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set pause status for target: identity|group|gateway",
		Run:   appPauseSetHandler,
	}
	cmd.Flags().String("target", "", "identity|group|gateway")
	cmd.Flags().Bool("paused", false, "pause status")
	_ = cmd.MarkFlagRequired("target")
	return &cmd
}

func appPauseSetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	paused, _ := cmd.Flags().GetBool("paused")
	ctx := context.Background()

	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}

	switch target {
	case "identity":
		if err := admin.SetIdentityUpdatePauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("identity broadcaster pause set", zap.Bool("paused", paused))
	case "group":
		if err := admin.SetGroupMessagePauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("group broadcaster pause set", zap.Bool("paused", paused))
	case "gateway":
		if err := admin.SetAppChainGatewayPauseStatus(ctx, paused); err != nil {
			logger.Fatal("write", zap.Error(err))
		}
		logger.Info("app-chain gateway pause set", zap.Bool("paused", paused))
	default:
		logger.Fatal("target must be identity|group|gateway")
	}
}

// --- bootstrapper ---

func appBootstrapperCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "bootstrapper",
		Short: "Get/Set payload bootstrapper (identity & group)",
	}
	cmd.AddCommand(appBootstrapperGetCmd(), appBootstrapperSetCmd())
	return &cmd
}

func appBootstrapperGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get bootstrapper addresses for identity & group",
		Run: func(cmd *cobra.Command, _ []string) {
			logger, err := cliLogger()
			if err != nil {
				log.Fatalf("could not build logger: %s", err)
			}
			ctx := context.Background()
			admin, err := setupAppChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}

			iu, err := admin.GetIdentityUpdateBootstrapper(ctx)
			if err != nil {
				logger.Fatal("read IU", zap.Error(err))
			}
			gm, err := admin.GetGroupMessageBootstrapper(ctx)
			if err != nil {
				logger.Fatal("read GM", zap.Error(err))
			}

			logger.Info("bootstrapper",
				zap.String("identity", iu.Hex()),
				zap.String("group", gm.Hex()),
			)
			if iu != gm {
				logger.Warn("identity and group bootstrappers differ")
			}
		},
	}
}

func appBootstrapperSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set bootstrapper address for BOTH identity & group",
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
			admin, err := setupAppChainAdmin(ctx, logger)
			if err != nil {
				logger.Fatal("setup admin", zap.Error(err))
			}

			if err := admin.SetIdentityUpdateBootstrapper(ctx, addr); err != nil {
				logger.Fatal("set identity bootstrapper", zap.Error(err))
			}
			if err := admin.SetGroupMessageBootstrapper(ctx, addr); err != nil {
				logger.Fatal("set group bootstrapper", zap.Error(err))
			}
			logger.Info("bootstrapper set", zap.String("address", addr.Hex()))
		},
	}
	cmd.Flags().String("address", "", "bootstrapper address")
	_ = cmd.MarkFlagRequired("address")
	return &cmd
}

// --- payload-size ---

func appPayloadSizeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "payload-size",
		Short: "Get/Set payload size bounds for broadcasters",
	}
	cmd.AddCommand(appPayloadSizeGetCmd(), appPayloadSizeSetCmd())
	return &cmd
}

func appPayloadSizeGetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get payload size for --target identity|group and --bound min|max",
		Run:   appPayloadSizeGetHandler,
	}
	cmd.Flags().String("target", "", "identity|group")
	cmd.Flags().String("bound", "", "min|max")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("bound")
	return &cmd
}

func appPayloadSizeGetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	bound, _ := cmd.Flags().GetString("bound")
	ctx := context.Background()

	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("setup admin", zap.Error(err))
	}

	switch target {
	case "identity":
		if bound == "min" {
			v, e := admin.GetIdentityUpdateMinPayloadSize(ctx)
			if e != nil {
				logger.Fatal("read", zap.Error(e))
			}
			logger.Info(
				"payload size",
				zap.String("target", "identity"),
				zap.String("bound", "min"),
				zap.Uint64("bytes", v),
			)
		} else if bound == "max" {
			v, e := admin.GetIdentityUpdateMaxPayloadSize(ctx)
			if e != nil {
				logger.Fatal("read", zap.Error(e))
			}
			logger.Info("payload size", zap.String("target", "identity"), zap.String("bound", "max"), zap.Uint64("bytes", v))
		} else {
			logger.Fatal("bound must be min|max")
		}
	case "group":
		if bound == "min" {
			v, e := admin.GetGroupMessageMinPayloadSize(ctx)
			if e != nil {
				logger.Fatal("read", zap.Error(e))
			}
			logger.Info(
				"payload size",
				zap.String("target", "group"),
				zap.String("bound", "min"),
				zap.Uint64("bytes", v),
			)
		} else if bound == "max" {
			v, e := admin.GetGroupMessageMaxPayloadSize(ctx)
			if e != nil {
				logger.Fatal("read", zap.Error(e))
			}
			logger.Info("payload size", zap.String("target", "group"), zap.String("bound", "max"), zap.Uint64("bytes", v))
		} else {
			logger.Fatal("bound must be min|max")
		}
	default:
		logger.Fatal("target must be identity|group")
	}
}

func appPayloadSizeSetCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set",
		Short: "Set payload size for --target identity|group and --bound min|max",
		Run:   appPayloadSizeSetHandler,
	}
	cmd.Flags().String("target", "", "identity|group")
	cmd.Flags().String("bound", "", "min|max")
	cmd.Flags().Uint32("size", 0, "size in bytes")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("bound")
	_ = cmd.MarkFlagRequired("size")
	return &cmd
}

func appPayloadSizeSetHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	target, _ := cmd.Flags().GetString("target")
	bound, _ := cmd.Flags().GetString("bound")
	size32, _ := cmd.Flags().GetUint32("size")
	size := uint64(size32)

	ctx := context.Background()
	admin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("setup admin", zap.Error(err))
	}

	switch target {
	case "identity":
		if bound == "min" {
			if err := admin.SetIdentityUpdateMinPayloadSize(ctx, size); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
		} else if bound == "max" {
			if err := admin.SetIdentityUpdateMaxPayloadSize(ctx, size); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
		} else {
			logger.Fatal("bound must be min|max")
		}
	case "group":
		if bound == "min" {
			if err := admin.SetGroupMessageMinPayloadSize(ctx, size); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
		} else if bound == "max" {
			if err := admin.SetGroupMessageMaxPayloadSize(ctx, size); err != nil {
				logger.Fatal("write", zap.Error(err))
			}
		} else {
			logger.Fatal("bound must be min|max")
		}
	default:
		logger.Fatal("target must be identity|group")
	}

	logger.Info(
		"payload size set",
		zap.String("target", target),
		zap.String("bound", bound),
		zap.Uint64("bytes", size),
	)
}

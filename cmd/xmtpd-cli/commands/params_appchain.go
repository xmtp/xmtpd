package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ---------- opts ----------

type AppGetOpts struct {
	Keys       []string
	Raw        bool
	TimeoutSec int
}

// ---------- root (params settlement) ----------

func paramsAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "app",
		Short:        "Operate on App chain parameters",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		appGetCmd(),
	)
	return cmd
}

// ---------- get ----------

func appGetCmd() *cobra.Command {
	var opts AppGetOpts

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get parameter(s) from the AppChain Parameter Registry (generic)",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return appGetHandler(opts)
		},
		Example: `
xmtpd-cli params app get \
  --key xmtp.rateRegistry.messageFee \
  --key xmtp.appChainGateway.paused`,
	}

	cmd.Flags().StringArrayVar(&opts.Keys, "key", nil, "parameter key (repeatable)")
	cmd.Flags().BoolVar(&opts.Raw, "raw", false, "treat value as 0x-prefixed 32-byte hex")
	cmd.Flags().IntVar(&opts.TimeoutSec, "timeout", 60, "timeout (seconds)")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

func appGetHandler(opts AppGetOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}

	if len(opts.Keys) == 0 {
		return errors.New("at least one --key is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.TimeoutSec)*time.Second,
	)
	defer cancel()

	paramAdmin, _, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}

	for _, k := range opts.Keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}

		if opts.Raw {
			val, gerr := paramAdmin.GetRawParameter(ctx, k)
			if gerr != nil {
				logger.Error("get parameter failed", zap.String("key", k), zap.Error(gerr))
				return gerr
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("bytes32", "0x"+common.Bytes2Hex(val[:])),
			)
			continue
		}

		switch paramType(k) {
		case ParamBool:
			v, err := paramAdmin.GetParameterBool(ctx, k)
			if err != nil {
				logger.Error("get bool failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.Bool("bool", v),
			)
		case ParamAddress:
			v, err := paramAdmin.GetParameterAddress(ctx, k)
			if err != nil {
				logger.Error("get address failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("address", v.Hex()),
			)
		case ParamUint8:
			v, err := paramAdmin.GetParameterUint8(ctx, k)
			if err != nil {
				logger.Error("get uint8 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint8("uint8", v))
		case ParamUint16:
			v, err := paramAdmin.GetParameterUint16(ctx, k)
			if err != nil {
				logger.Error("get uint16 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint16("uint16", v))
		case ParamUint32:
			v, err := paramAdmin.GetParameterUint32(ctx, k)
			if err != nil {
				logger.Error("get uint32 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint32("uint32", v))
		case ParamUint64:
			v, err := paramAdmin.GetParameterUint64(ctx, k)
			if err != nil {
				logger.Error("get uint64 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint64("uint64", v))
		case ParamUint96:
			v, err := paramAdmin.GetParameterUint96(ctx, k)
			if err != nil {
				logger.Error("get uint96 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("uint96", v.String()),
			)
		default:
			// Fallback: raw
			val, gerr := paramAdmin.GetRawParameter(ctx, k)
			if gerr != nil {
				logger.Error("get parameter failed", zap.String("key", k), zap.Error(gerr))
				return gerr
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("bytes32", "0x"+common.Bytes2Hex(val[:])),
			)
		}
	}

	return nil
}

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

type AppSetOpts struct {
	KVs        []string // each "key=value" where value is 0x + 64 hex chars
	NoWait     bool     // reserved (if your ParameterAdmin batches as a tx)
	TimeoutSec int
}

type AppGetOpts struct {
	Keys       []string
	TimeoutSec int
}

// ---------- root (params settlement) ----------

func paramsAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Operate on App chain parameters",
	}
	cmd.AddCommand(
		appSetCmd(),
		appGetCmd(),
	)
	return cmd
}

// ---------- set ----------

func appSetCmd() *cobra.Command {
	var opts AppSetOpts

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set parameter(s) in the AppChain Parameter Registry (generic key/value)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return appSetHandler(opts)
		},
		Example: `
xmtpd-cli params app set \
  --kv xmtp.rateRegistry.messageFee=0x00000000000000000000000000000000000000000000000000000000000003e8 \
  --kv xmtp.settlementChainGateway.paused=0x0000000000000000000000000000000000000000000000000000000000000001`,
	}

	cmd.Flags().StringArrayVar(&opts.KVs, "kv", nil, "key=value (value: 0x-prefixed 32-byte hex)")
	cmd.Flags().
		BoolVar(&opts.NoWait, "no-wait", false, "do not wait for confirmation (if applicable)")
	cmd.Flags().IntVar(&opts.TimeoutSec, "timeout", 120, "timeout (seconds)")
	_ = cmd.MarkFlagRequired("kv")

	return cmd
}

func appSetHandler(opts AppSetOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}

	if len(opts.KVs) == 0 {
		return errors.New("at least one --kv is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.TimeoutSec)*time.Second,
	)
	defer cancel()

	paramAdmin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}

	type kv struct {
		key string
		val [32]byte
	}

	var items []kv
	for _, kvs := range opts.KVs {
		key, valHex, perr := splitKV(kvs)
		if perr != nil {
			return perr
		}
		b32, perr := parseBytes32(valHex)
		if perr != nil {
			return fmt.Errorf("invalid value for key %s: %w", key, perr)
		}
		items = append(items, kv{key: key, val: b32})
	}

	for _, it := range items {
		if err := paramAdmin.SetRawParameter(ctx, it.key, it.val); err != nil {
			// If you wrap "no change" in a typed error, check it here and log a friendly line.
			logger.Error("set parameter failed", zap.String("key", it.key), zap.Error(err))
			return err
		}
		logger.Info("parameter set", zap.String("key", it.key))
	}

	logger.Info("all parameters set successfully")
	return nil
}

// ---------- get ----------

func appGetCmd() *cobra.Command {
	var opts AppGetOpts

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get parameter(s) from the AppChain Parameter Registry (generic)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return appGetHandler(opts)
		},
		Example: `
xmtpd-cli params app get \
  --key xmtp.rateRegistry.messageFee \
  --key xmtp.appChainGateway.paused`,
	}

	cmd.Flags().StringArrayVar(&opts.Keys, "key", nil, "parameter key (repeatable)")
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

	paramAdmin, err := setupAppChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}

	for _, k := range opts.Keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}

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

	return nil
}

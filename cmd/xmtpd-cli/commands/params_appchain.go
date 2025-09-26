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
		appGetCmd(),
	)
	return cmd
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

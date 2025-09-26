package commands

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ---------- opts ----------

type BridgeSendOpts struct {
	Keys       []string
	NoWait     bool
	TimeoutSec int
}

// ---------- root ----------

func paramsBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge",
		Short: "Bridge parameters Settlement → App",
	}
	cmd.AddCommand(bridgeSendCmd())
	return cmd
}

func bridgeSendCmd() *cobra.Command {
	var opts BridgeSendOpts

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send parameters via SettlementChainGateway.sendParameters(keys)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return bridgeSendHandler(opts)
		},
		Example: `
xmtpd-cli params bridge send \
  --key xmtp.nodeRegistry.maxCanonicalNodes \
  --key xmtp.groupMessageBroadcaster.maxPayloadSize`,
	}

	cmd.Flags().StringArrayVar(&opts.Keys, "key", nil, "parameter key to bridge (repeatable)")
	cmd.Flags().BoolVar(&opts.NoWait, "no-wait", false, "do not wait for confirmation")
	cmd.Flags().IntVar(&opts.TimeoutSec, "timeout", 180, "wait timeout seconds")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

func bridgeSendHandler(opts BridgeSendOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}
	if len(opts.Keys) == 0 {
		return fmt.Errorf("at least one --key is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.TimeoutSec)*time.Second,
	)
	defer cancel()

	paramAdmin, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}
	for _, k := range opts.Keys {
		raw, perr := paramAdmin.GetRawParameter(ctx, k)
		if perr != nil {
			logger.Error("fetch for preview failed", zap.String("key", k), zap.Error(perr))
			return perr
		}
		logger.Info("bridge preview",
			zap.String("key", k),
			zap.String("bytes32", "0x"+hex.EncodeToString(raw[:])),
		)
	}

	err = paramAdmin.BridgeParameters(ctx, opts.Keys)
	if err != nil {
		logger.Error("bridge failed", zap.Error(err))
		return err
	}
	logger.Info("bridge sent successfully")

	return nil
}

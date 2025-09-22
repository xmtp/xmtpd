package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
	"github.com/xmtp/xmtpd/pkg/currency"

	"github.com/ethereum/go-ethereum/common"

	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// ---- Options ----

type DepositOpts struct {
	Amount   string
	GasLimit int64
	GasPrice int64
}

type WithdrawOpts struct {
	PrivateKey   string
	Recipient    string
	Amount       string
	WithdrawType string
}

type BalancesOpts struct {
	Address string
}

// ---- Root ----

func fundsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "funds",
		Short:        "Manage deposits, withdrawals, and balances",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		depositCmd(),
		withdrawCmd(),
		checkWithdrawalsCmd(),
		balancesCmd(),
		mintCmd(),
	)

	return cmd
}

// ---- deposit ----

func depositCmd() *cobra.Command {
	var opts DepositOpts

	cmd := &cobra.Command{
		Use:          "deposit",
		Short:        "Deposit funds to XMTP",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return depositHandler(opts)
		},
		Example: `
Usage:
  xmtpd-cli funds deposit --amount <amount> [--gas-limit <gas-limit>] [--gas-price <gas-price>]

Example:
  xmtpd-cli funds deposit --amount 1000000000000000000
`,
	}
	cmd.Flags().
		StringVar(&opts.Amount, "amount", "", "amount to deposit (wei-scale or token base units)")
	_ = cmd.MarkFlagRequired("amount")

	cmd.Flags().Int64Var(&opts.GasLimit, "gas-limit", 3000000, "gas limit")
	cmd.Flags().Int64Var(&opts.GasPrice, "gas-price", 2000000000, "gas price")

	return cmd
}

func depositHandler(opts DepositOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupFundsAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup settlement chain admin", zap.Error(err))
		return err
	}

	amount, ok := new(big.Int).SetString(opts.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid --amount (raw uint256) %q", opts.Amount)
	}

	if amount.Sign() == -1 {
		return fmt.Errorf("invalid --amount %d; must be non-negative", amount)
	}
	if opts.GasLimit < 0 {
		return fmt.Errorf("invalid --gas-limit %d; must be non-negative", opts.GasLimit)
	}
	if opts.GasPrice < 0 {
		return fmt.Errorf("invalid --gas-price %d; must be non-negative", opts.GasPrice)
	}

	gasLimit := big.NewInt(opts.GasLimit)
	gasPrice := big.NewInt(opts.GasPrice)

	err = admin.Deposit(ctx, amount, gasLimit, gasPrice)
	if err != nil {
		logger.Error("could not deposit funds", zap.Error(err))
		return err
	}

	return nil
}

// ---- withdraw ----

func withdrawCmd() *cobra.Command {
	var opts WithdrawOpts

	cmd := &cobra.Command{
		Use:          "withdraw",
		Short:        "Withdraw funds from XMTP",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return withdrawHandler(opts)
		},
		Example: `
Usage:
  xmtpd-cli funds withdraw --private-key <private_key> --recipient <recipient> --amount <amount> [--withdraw-type normal]

Example:
  xmtpd-cli funds withdraw --private-key 0xabc... --recipient 0xdef... --amount 1000000000000000000 --withdraw-type normal
`,
	}

	cmd.Flags().
		StringVar(&opts.PrivateKey, "private-key", "", "private key to use for signing the withdrawal")
	cmd.Flags().StringVar(&opts.Recipient, "recipient", "", "recipient address")
	cmd.Flags().
		StringVar(&opts.Amount, "amount", "", "amount to withdraw (wei-scale or token base units)")
	cmd.Flags().
		StringVar(&opts.WithdrawType, "withdraw-type", "normal", "withdrawal type (e.g., normal)")

	_ = cmd.MarkFlagRequired("private-key")
	_ = cmd.MarkFlagRequired("recipient")
	_ = cmd.MarkFlagRequired("amount")

	return cmd
}

func withdrawHandler(_ WithdrawOpts) error {
	// TODO: implement withdraw_from_xmtp(privateKey, recipient, amount, withdrawType)
	return fmt.Errorf("withdraw not implemented yet")
}

// ---- check-withdrawals ----

func checkWithdrawalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "check-withdrawals",
		Short:        "Check pending/processed withdrawals",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return checkWithdrawalsHandler()
		},
		Example: `
Usage:
  xmtpd-cli funds check-withdrawals

Example:
  xmtpd-cli funds check-withdrawals
`,
	}
	return cmd
}

func checkWithdrawalsHandler() error {
	// TODO: implement check_withdrawals()
	return fmt.Errorf("check-withdrawals not implemented yet")
}

// ---- balances ----

func balancesCmd() *cobra.Command {
	var opts BalancesOpts

	cmd := &cobra.Command{
		Use:          "balances",
		Short:        "Check balances for an address",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return balancesHandler(opts)
		},
		Example: `
Usage:
  xmtpd-cli funds balances --address <address>

Example:
  xmtpd-cli funds balances --address 0xabc...
`,
	}

	cmd.Flags().StringVar(&opts.Address, "address", "", "address to query balances for")
	_ = cmd.MarkFlagRequired("address")

	return cmd
}

func balancesHandler(opts BalancesOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupFundsAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup funds admin", zap.Error(err))
		return err
	}

	return admin.Balances(ctx, common.HexToAddress(opts.Address))
}

func mintCmd() *cobra.Command {
	var toAddr options.AddressFlag
	var amountHuman string
	var raw bool
	cmd := &cobra.Command{
		Use:          "mint",
		Hidden:       true,
		Short:        "Mint mock underlying fee token to an address (max 10000 tokens)",
		SilenceUsage: true,
		Example: `
# Mint 1000 tokens to 0xabc... using the token's own decimals
xmtpd-cli funds mint \
  --to 0xRecipient \
  --amount 1000

# If you already have the raw uint256 (no scaling), use --raw
xmtpd-cli funds mint \
  --to 0xRecipient \
  --amount 1000000000 --raw
`,
		RunE: func(*cobra.Command, []string) error {
			return mintHandler(toAddr.Address, amountHuman, raw)
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

func mintHandler(to common.Address, amountStr string, raw bool) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	ctx := context.Background()

	admin, err := setupFundsAdmin(ctx, logger)
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

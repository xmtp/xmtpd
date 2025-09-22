// file: cmd/funds.go
package commands

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// ---- Options ----

type DepositOpts struct {
	PrivateKey string
	Recipient  string
	Amount     string
	TokenType  string
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
  xmtpd-cli funds deposit --private-key <private_key> --recipient <recipient> --amount <amount> [--token-type xusd]

Example:
  xmtpd-cli funds deposit --private-key 0xabc... --recipient 0xdef... --amount 1000000000000000000 --token-type xusd
`,
	}

	cmd.Flags().
		StringVar(&opts.PrivateKey, "private-key", "", "private key to use for signing the deposit")
	cmd.Flags().StringVar(&opts.Recipient, "recipient", "", "recipient address")
	cmd.Flags().
		StringVar(&opts.Amount, "amount", "", "amount to deposit (wei-scale or token base units)")
	cmd.Flags().StringVar(&opts.TokenType, "token-type", "xusd", "token type (default: xusd)")

	_ = cmd.MarkFlagRequired("private-key")
	_ = cmd.MarkFlagRequired("recipient")
	_ = cmd.MarkFlagRequired("amount")

	return cmd
}

func depositHandler(_ DepositOpts) error {
	// TODO: implement deposit_to_xmtp(privateKey, recipient, amount, tokenType)
	return fmt.Errorf("deposit not implemented yet")
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

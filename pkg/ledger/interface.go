package ledger

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/currency"
)

// The Ledger interface handles settled balances from the blockchain
type ILedger interface {
	// Register a deposit event to a payer's balance
	Deposit(
		ctx context.Context,
		payerID int32,
		amount currency.PicoDollar,
		eventID EventID,
	) error
	// Register a withdrawal event, which immediately reduces the payer's balance
	InitiateWithdrawal(
		ctx context.Context,
		payerID int32,
		amount currency.PicoDollar,
		eventID EventID,
	) error
	// Cancel a previous withdrawal
	CancelWithdrawal(ctx context.Context, payerID int32, eventID EventID) error
	// Decrement a payer's balance when usage is settled
	SettleUsage(
		ctx context.Context,
		payerID int32,
		amount currency.PicoDollar,
		eventID EventID,
	) error
	// Get the balance for a payer from the settled usage ledger
	GetBalance(ctx context.Context, payerID int32) (currency.PicoDollar, error)
}

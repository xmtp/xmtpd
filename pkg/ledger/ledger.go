package ledger

import (
	"bytes"
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

type Ledger struct {
	queries *queries.Queries
	logger  *zap.Logger
}

func NewLedger(logger *zap.Logger, querier *queries.Queries) *Ledger {
	return &Ledger{
		queries: querier,
		logger:  logger.Named("ledger"),
	}
}

func (l *Ledger) GetBalance(ctx context.Context, payerID int32) (currency.PicoDollar, error) {
	balance, err := l.queries.GetPayerBalance(ctx, payerID)
	if err != nil {
		return 0, err
	}
	return currency.PicoDollar(balance), nil
}

func (l *Ledger) Deposit(
	ctx context.Context,
	payerID int32,
	amount currency.PicoDollar,
	eventID EventID,
) error {
	var err error
	if err = validateAmount(amount); err != nil {
		return err
	}
	if err = validateEventID(eventID); err != nil {
		return err
	}

	return l.queries.InsertPayerLedgerEvent(ctx, queries.InsertPayerLedgerEventParams{
		EventID:           eventID[:],
		PayerID:           payerID,
		AmountPicodollars: int64(amount),
		EventType:         int16(EVENT_TYPE_DEPOSIT),
	})
}

func (l *Ledger) InitiateWithdrawal(
	ctx context.Context,
	payerID int32,
	amount currency.PicoDollar,
	eventID EventID,
) error {
	var err error
	if err = validateAmount(amount); err != nil {
		return err
	}
	if err = validateEventID(eventID); err != nil {
		return err
	}

	return l.queries.InsertPayerLedgerEvent(ctx, queries.InsertPayerLedgerEventParams{
		EventID:           eventID[:],
		PayerID:           payerID,
		AmountPicodollars: int64(amount) * -1,
		EventType:         int16(EVENT_TYPE_WITHDRAWAL),
	})
}

func (l *Ledger) CancelWithdrawal(ctx context.Context, payerID int32, eventID EventID) error {
	lastWithdrawal, err := l.queries.GetLastEvent(ctx, queries.GetLastEventParams{
		PayerID:   payerID,
		EventType: int16(EVENT_TYPE_WITHDRAWAL),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrWithdrawalNotFound
		}
		return err
	}

	// For additional safety, check if the last withdrawal was canceled by another event.
	// The smart contract should protect against this, but we are double checking.
	lastCancel, err := l.queries.GetLastEvent(ctx, queries.GetLastEventParams{
		PayerID:   payerID,
		EventType: int16(EVENT_TYPE_CANCELED_WITHDRAWAL),
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil && lastCancel.CreatedAt.After(lastWithdrawal.CreatedAt) {
		if bytes.Equal(lastCancel.EventID, eventID[:]) {
			// Cancelation already complete
			return nil
		}
		l.logger.Warn(
			"multiple cancelation events for a single withdrawal",
			zap.String("event_id", eventID.String()),
			zap.String("last_cancel_event_id", EventID(lastCancel.EventID).String()),
			zap.String("last_withdrawal_event_id", EventID(lastWithdrawal.EventID).String()),
		)

		return ErrWithdrawalAlreadyCanceled
	}

	return l.queries.InsertPayerLedgerEvent(ctx, queries.InsertPayerLedgerEventParams{
		EventID:           eventID[:],
		PayerID:           payerID,
		AmountPicodollars: int64(lastWithdrawal.AmountPicodollars) * -1,
		EventType:         int16(EVENT_TYPE_CANCELED_WITHDRAWAL),
	})
}

func (l *Ledger) SettleUsage(
	ctx context.Context,
	payerID int32,
	amount currency.PicoDollar,
	eventID EventID,
) error {
	var err error
	if err = validateAmount(amount); err != nil {
		return err
	}
	if err = validateEventID(eventID); err != nil {
		return err
	}

	return l.queries.InsertPayerLedgerEvent(ctx, queries.InsertPayerLedgerEventParams{
		EventID:           eventID[:],
		PayerID:           payerID,
		AmountPicodollars: int64(amount) * -1,
		EventType:         int16(EVENT_TYPE_SETTLEMENT),
	})
}

func (l *Ledger) FindOrCreatePayer(
	ctx context.Context,
	payerAddress common.Address,
) (int32, error) {
	return l.queries.FindOrCreatePayer(ctx, payerAddress.Hex())
}

func validateAmount(amount currency.PicoDollar) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func validateEventID(eventID EventID) error {
	if eventID == (EventID{}) {
		return ErrInvalidEventID
	}
	return nil
}

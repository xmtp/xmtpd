package contracts

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/currency"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/ledger"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

const (
	ErrParsePayerRegistryLog       = "error parsing payer registry log"
	ErrPayerRegistryUnhandledEvent = "unknown payer registry event"
	ErrFindOrCreatePayer           = "error finding or creating payer"
	ErrLedgerDeposit               = "error depositing to ledger"
	ErrLedgerInitiateWithdrawal    = "error initiating withdrawal from ledger"
	ErrInvalidEvent                = "invalid event"
	ErrLedgerSettleUsage           = "error settling usage in ledger"

	// WithdrawalFinalized is not handled, as it might be redundant with WithdrawalRequested.
	payerRegistryDepositEvent             = "Deposit"
	payerRegistryWithdrawalRequestedEvent = "WithdrawalRequested"
	payerRegistryWithdrawalCancelledEvent = "WithdrawalCancelled"
	payerRegistryUsageSettledEvent        = "UsageSettled"
)

type PayerRegistryStorer struct {
	abi      *abi.ABI
	logger   *zap.Logger
	ledger   ledger.ILedger
	contract *pr.PayerRegistry
}

var _ c.ILogStorer = &PayerRegistryStorer{}

func NewPayerRegistryStorer(
	logger *zap.Logger,
	contract *pr.PayerRegistry,
	payerLedger ledger.ILedger,
) (*PayerRegistryStorer, error) {
	abi, err := pr.PayerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &PayerRegistryStorer{
		abi:      abi,
		ledger:   payerLedger,
		logger:   logger.Named(utils.StorerLoggerName),
		contract: contract,
	}, nil
}

func (s *PayerRegistryStorer) StoreLog(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if len(log.Topics) == 0 {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, errors.New("no topics"))
	}

	event, err := s.abi.EventByID(log.Topics[0])
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	switch event.Name {
	case payerRegistryDepositEvent:
		return s.handleDeposit(ctx, log)
	case payerRegistryWithdrawalRequestedEvent:
		return s.handleWithdrawalRequested(ctx, log)
	case payerRegistryWithdrawalCancelledEvent:
		return s.handleWithdrawalCanceled(ctx, log)
	case payerRegistryUsageSettledEvent:
		return s.handleUsageSettled(ctx, log)
	default:
		s.logger.Info("unknown event", utils.EventField(event.Name))
		return re.NewNonRecoverableError(ErrPayerRegistryUnhandledEvent, errors.New(event.Name))
	}
}

func (s *PayerRegistryStorer) handleDeposit(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received event Deposit", zap.Any("log", log))
	}

	var err error
	var parsedEvent *pr.PayerRegistryDeposit
	parsedEvent, err = s.contract.ParseDeposit(log)
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	payerID, err := s.ledger.FindOrCreatePayer(ctx, parsedEvent.Payer)
	if err != nil {
		return re.NewRecoverableError(ErrFindOrCreatePayer, err)
	}

	amount := currency.FromMicrodollars(currency.MicroDollar(parsedEvent.Amount.Int64()))
	eventID := ledger.BuildEventID(log)

	if err = s.ledger.Deposit(
		ctx,
		payerID,
		amount,
		eventID,
	); err != nil {
		return wrapLedgerError(err, ErrLedgerDeposit)
	}

	s.logger.Debug(
		"deposit successful",
		utils.PayerAddressField(parsedEvent.Payer.Hex()),
		utils.AmountField(amount.String()),
		utils.EventIDField(eventID.String()),
	)

	return nil
}

func (s *PayerRegistryStorer) handleWithdrawalRequested(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received event WithdrawalRequested", zap.Any("log", log))
	}

	var err error
	var parsedEvent *pr.PayerRegistryWithdrawalRequested
	parsedEvent, err = s.contract.ParseWithdrawalRequested(log)
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	payerID, err := s.ledger.FindOrCreatePayer(ctx, parsedEvent.Payer)
	if err != nil {
		return re.NewRecoverableError(ErrFindOrCreatePayer, err)
	}

	amount := currency.FromMicrodollars(currency.MicroDollar(parsedEvent.Amount.Int64()))
	eventID := ledger.BuildEventID(log)

	if err = s.ledger.InitiateWithdrawal(
		ctx,
		payerID,
		amount,
		eventID,
	); err != nil {
		return wrapLedgerError(err, ErrLedgerInitiateWithdrawal)
	}

	s.logger.Debug(
		"withdrawal requested successful",
		utils.PayerAddressField(parsedEvent.Payer.Hex()),
		utils.AmountField(amount.String()),
		utils.EventIDField(eventID.String()),
	)

	return nil
}

func (s *PayerRegistryStorer) handleUsageSettled(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received event UsageSettled", zap.Any("log", log))
	}

	var err error
	var parsedEvent *pr.PayerRegistryUsageSettled
	parsedEvent, err = s.contract.ParseUsageSettled(log)
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	payerID, err := s.ledger.FindOrCreatePayer(ctx, parsedEvent.Payer)
	if err != nil {
		return re.NewRecoverableError(ErrFindOrCreatePayer, err)
	}

	amount := currency.FromMicrodollars(currency.MicroDollar(parsedEvent.Amount.Int64()))
	eventID := ledger.BuildEventID(log)

	if err = s.ledger.SettleUsage(ctx, payerID, amount, eventID); err != nil {
		return wrapLedgerError(err, ErrLedgerSettleUsage)
	}

	s.logger.Debug(
		"usage settled",
		utils.PayerAddressField(parsedEvent.Payer.Hex()),
		utils.AmountField(amount.String()),
		utils.EventIDField(eventID.String()),
	)

	return nil
}

func (s *PayerRegistryStorer) handleWithdrawalCanceled(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received event WithdrawalCancelled", zap.Any("log", log))
	}

	var err error
	var parsedEvent *pr.PayerRegistryWithdrawalCancelled
	parsedEvent, err = s.contract.ParseWithdrawalCancelled(log)
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	payerID, err := s.ledger.FindOrCreatePayer(ctx, parsedEvent.Payer)
	if err != nil {
		return re.NewRecoverableError(ErrFindOrCreatePayer, err)
	}

	eventID := ledger.BuildEventID(log)

	if err = s.ledger.CancelWithdrawal(ctx, payerID, eventID); err != nil {
		return wrapLedgerError(err, "error canceling withdrawal")
	}

	return nil
}

func wrapLedgerError(err error, msg string) re.RetryableError {
	if errors.Is(err, ledger.ErrInvalidAmount) || errors.Is(err, ledger.ErrInvalidEventID) {
		return re.NewNonRecoverableError(ErrInvalidEvent, err)
	}
	return re.NewRecoverableError(msg, err)
}

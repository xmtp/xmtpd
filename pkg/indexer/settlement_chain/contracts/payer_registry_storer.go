package contracts

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	re "github.com/xmtp/xmtpd/pkg/errors"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"go.uber.org/zap"
)

const (
	ErrParsePayerRegistryLog       = "error parsing payer registry log"
	ErrPayerRegistryUnhandledEvent = "unknown payer registry event"

	// WithdrawalFinalized is not handled, as it might be redundant with WithdrawalRequested.
	payerRegistryDepositEvent             = "Deposit"
	payerRegistryWithdrawalRequestedEvent = "WithdrawalRequested"
	payerRegistryWithdrawalCancelledEvent = "WithdrawalCancelled"
	payerRegistryUsageSettledEvent        = "UsageSettled"
)

type PayerRegistryStorer struct {
	abi     *abi.ABI
	queries *queries.Queries
	logger  *zap.Logger
}

var _ c.ILogStorer = &PayerRegistryStorer{}

func NewPayerRegistryStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *pr.PayerRegistry,
) (*PayerRegistryStorer, error) {
	abi, err := pr.PayerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &PayerRegistryStorer{
		abi:     abi,
		queries: queries,
		logger:  logger.Named("storer"),
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
		fmt.Println("error", err)
		return re.NewNonRecoverableError(ErrParsePayerRegistryLog, err)
	}

	switch event.Name {
	case payerRegistryDepositEvent:
		s.logger.Info("Deposit", zap.Any("log", log))
	case payerRegistryWithdrawalRequestedEvent:
		s.logger.Info("WithdrawalRequested", zap.Any("log", log))
	case payerRegistryWithdrawalCancelledEvent:
		s.logger.Info("WithdrawalCancelled", zap.Any("log", log))
	case payerRegistryUsageSettledEvent:
		s.logger.Info("UsageSettled", zap.Any("log", log))
	default:
		s.logger.Info("Unknown event", zap.String("event", event.Name))
		return re.NewNonRecoverableError(ErrPayerRegistryUnhandledEvent, errors.New(event.Name))
	}

	return nil
}

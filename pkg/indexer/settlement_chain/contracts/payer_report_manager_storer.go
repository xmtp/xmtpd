package contracts

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	re "github.com/xmtp/xmtpd/pkg/errors"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"go.uber.org/zap"
)

const (
	ErrParsePayerReportManagerLog       = "error parsing payer report manager log"
	ErrPayerReportManagerUnhandledEvent = "unknown payer report manager event"

	payerReportManagerPayerReportSubmittedEvent     = "PayerReportSubmitted"
	payerReportManagerPayerReportSubsetSettledEvent = "PayerReportSubsetSettled"
)

type PayerReportManagerStorer struct {
	abi     *abi.ABI
	queries *queries.Queries
	logger  *zap.Logger
}

var _ c.ILogStorer = &PayerReportManagerStorer{}

func NewPayerReportManagerStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *p.PayerReportManager,
) (*PayerReportManagerStorer, error) {
	abi, err := p.PayerReportManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &PayerReportManagerStorer{
		abi:     abi,
		queries: queries,
		logger:  logger.Named("storer"),
	}, nil
}

func (s *PayerReportManagerStorer) StoreLog(
	ctx context.Context,
	log types.Log,
) re.RetryableError {
	if len(log.Topics) == 0 {
		return re.NewNonRecoverableError(ErrParsePayerReportManagerLog, errors.New("no topics"))
	}

	event, err := s.abi.EventByID(log.Topics[0])
	if err != nil {
		return re.NewNonRecoverableError(ErrParsePayerReportManagerLog, err)
	}

	switch event.Name {
	case payerReportManagerPayerReportSubmittedEvent:
		s.logger.Info("PayerReportSubmitted", zap.Any("log", log))
	case payerReportManagerPayerReportSubsetSettledEvent:
		s.logger.Info("PayerReportSubsetSettled", zap.Any("log", log))
	default:
		s.logger.Info("Unknown event", zap.String("event", event.Name))
		return re.NewNonRecoverableError(
			ErrPayerReportManagerUnhandledEvent,
			errors.New(event.Name),
		)
	}

	return nil
}

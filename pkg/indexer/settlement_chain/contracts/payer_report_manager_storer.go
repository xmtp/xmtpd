package contracts

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

const (
	ErrReportAlreadyExists              = "report already present in database"
	ErrParsePayerReportManagerLog       = "error parsing payer report manager log"
	ErrPayerReportManagerUnhandledEvent = "unknown payer report manager event"
	ErrBuildPayerReportID               = "error building payer report id"
	ErrStoreReport                      = "error storing report"
	ErrSetReportSubmissionStatus        = "error setting report submission status"
	ErrSetReportAttestationStatus       = "error setting report attestation status"
	ErrLoadReportByIndex                = "error loading report by index"

	payerReportManagerPayerReportSubmittedEvent     = "PayerReportSubmitted"
	payerReportManagerPayerReportSubsetSettledEvent = "PayerReportSubsetSettled"
)

type PayerReportManagerStorer struct {
	abi             *abi.ABI
	store           payerreport.IPayerReportStore
	logger          *zap.Logger
	contract        PayerReportManagerContract
	domainSeparator common.Hash
}

var _ c.ILogStorer = &PayerReportManagerStorer{}

func NewPayerReportManagerStorer(
	db *sql.DB,
	logger *zap.Logger,
	contract PayerReportManagerContract,
) (*PayerReportManagerStorer, error) {
	abi, err := p.PayerReportManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	domainSeparator, err := contract.DOMAINSEPARATOR(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	store := payerreport.NewStore(db, logger)

	return &PayerReportManagerStorer{
		abi:             abi,
		store:           store,
		logger:          logger.Named("storer"),
		contract:        contract,
		domainSeparator: common.Hash(domainSeparator),
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
		var parsedEvent *p.PayerReportManagerPayerReportSubmitted
		if parsedEvent, err = s.contract.ParsePayerReportSubmitted(log); err != nil {
			return re.NewNonRecoverableError(ErrParsePayerReportManagerLog, err)
		}

		err := s.setReportSubmitted(ctx, parsedEvent)
		if err == nil {
			s.logger.Info(
				"Successfully stored payer report submitted event",
				zap.Uint32("originatorNodeID", parsedEvent.OriginatorNodeId),
				zap.Uint64("startSequenceID", parsedEvent.StartSequenceId),
				zap.Uint64("endSequenceID", parsedEvent.EndSequenceId),
				zap.String("payersMerkleRoot", hex.EncodeToString(parsedEvent.PayersMerkleRoot[:])),
				zap.Uint32s("activeNodeIDs", parsedEvent.NodeIds),
			)
		} else {
			if strings.Contains(err.Error(), ErrReportAlreadyExists) {
				s.logger.Debug("skipping store report, report already exists",
					zap.Uint32("originatorNodeID", parsedEvent.OriginatorNodeId),
					zap.Uint64("startSequenceID", parsedEvent.StartSequenceId),
					zap.Uint64("endSequenceID", parsedEvent.EndSequenceId),
					zap.String("payersMerkleRoot", hex.EncodeToString(parsedEvent.PayersMerkleRoot[:])),
					zap.Uint32s("activeNodeIDs", parsedEvent.NodeIds),
				)
			} else {
				return err
			}
		}

	case payerReportManagerPayerReportSubsetSettledEvent:
		s.logger.Info("PayerReportSubsetSettled", zap.Any("log", log))
		var parsedEvent *p.PayerReportManagerPayerReportSubsetSettled
		if parsedEvent, err = s.contract.ParsePayerReportSubsetSettled(log); err != nil {
			return re.NewNonRecoverableError(ErrParsePayerReportManagerLog, err)
		}

		// Only mark report as settled when all payers have been processed (remaining == 0)
		if parsedEvent.Remaining == 0 {
			if err := s.setReportSettled(ctx, parsedEvent); err != nil {
				return err
			}
			s.logger.Info(
				"Payer report fully settled",
				zap.Uint32("originatorNodeID", parsedEvent.OriginatorNodeId),
				zap.String("payerReportIndex", parsedEvent.PayerReportIndex.String()),
				zap.Uint32("count", parsedEvent.Count),
				zap.String("feesSettled", parsedEvent.FeesSettled.String()),
			)
		} else {
			s.logger.Info(
				"Payer report partially settled",
				zap.Uint32("originatorNodeID", parsedEvent.OriginatorNodeId),
				zap.String("payerReportIndex", parsedEvent.PayerReportIndex.String()),
				zap.Uint32("count", parsedEvent.Count),
				zap.Uint32("remaining", parsedEvent.Remaining),
				zap.String("feesSettled", parsedEvent.FeesSettled.String()),
			)
		}

	default:
		s.logger.Info("Unknown event", zap.String("event", event.Name))
		return re.NewNonRecoverableError(
			ErrPayerReportManagerUnhandledEvent,
			errors.New(event.Name),
		)
	}

	return nil
}

func (s *PayerReportManagerStorer) getReportIDFromIndex(
	ctx context.Context,
	nodeID uint32,
	index *big.Int,
) (*payerreport.ReportID, error) {
	result, err := s.contract.GetPayerReport(&bind.CallOpts{
		Context: ctx,
	}, nodeID, index)
	if err != nil {
		return nil, err
	}

	reportID, err := payerreport.BuildPayerReportID(
		nodeID,
		result.StartSequenceId,
		result.EndSequenceId,
		result.EndMinuteSinceEpoch,
		result.PayersMerkleRoot,
		result.NodeIds,
		s.domainSeparator,
	)
	if err != nil {
		return nil, err
	}

	return reportID, nil
}

func (s *PayerReportManagerStorer) setReportSettled(
	ctx context.Context,
	event *p.PayerReportManagerPayerReportSubsetSettled,
) re.RetryableError {
	reportID, err := s.getReportIDFromIndex(ctx, event.OriginatorNodeId, event.PayerReportIndex)
	if err != nil {
		return re.NewRecoverableError(ErrLoadReportByIndex, err)
	}

	if err = s.store.SetReportSettled(ctx, *reportID); err != nil {
		return re.NewRecoverableError(ErrSetReportSubmissionStatus, err)
	}

	return nil
}

func (s *PayerReportManagerStorer) setReportSubmitted(
	ctx context.Context,
	event *p.PayerReportManagerPayerReportSubmitted,
) re.RetryableError {
	var reportID *payerreport.ReportID
	var err error

	if reportID, err = payerreport.BuildPayerReportID(
		event.OriginatorNodeId,
		event.StartSequenceId,
		event.EndSequenceId,
		event.EndMinuteSinceEpoch,
		event.PayersMerkleRoot,
		event.NodeIds,
		s.domainSeparator,
	); err != nil {
		return re.NewNonRecoverableError(ErrBuildPayerReportID, err)
	}

	report := &payerreport.PayerReport{
		ID:                  *reportID,
		OriginatorNodeID:    event.OriginatorNodeId,
		StartSequenceID:     event.StartSequenceId,
		EndSequenceID:       event.EndSequenceId,
		EndMinuteSinceEpoch: event.EndMinuteSinceEpoch,
		PayersMerkleRoot:    event.PayersMerkleRoot,
		ActiveNodeIDs:       event.NodeIds,
	}
	// NOTE: this is the first time we might be hearing about this
	// depends on the race between indexer and sync
	numRows, err := s.store.StoreReport(ctx, report)
	if err != nil {
		return re.NewRecoverableError(ErrStoreReport, err)
	}

	// Will only set the status to Submitted if it was previously Pending.
	// If it is already settled, this is a no-op.
	if err = s.store.SetReportSubmitted(ctx, *reportID); err != nil {
		return re.NewRecoverableError(ErrSetReportSubmissionStatus, err)
	}

	// Sets attestation status to Approved, if it was previously Pending.
	// If it is already approved or rejected, this is a no-op.
	if err = s.store.SetReportAttestationApproved(ctx, *reportID); err != nil {
		return re.NewRecoverableError(ErrSetReportAttestationStatus, err)
	}

	if numRows == 0 {
		return re.NewNonRecoverableError(ErrReportAlreadyExists, nil)
	}

	return nil
}

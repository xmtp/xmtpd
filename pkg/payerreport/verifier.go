package payerreport

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var (
	ErrMismatchOriginator = errors.New(
		"originator id mismatch between old and new report",
	)
	ErrInvalidReportStart = errors.New(
		"report does not start where the previous report ended",
	)
	ErrInvalidSequenceID       = errors.New("invalid sequence id")
	ErrInvalidOriginatorID     = errors.New("originator id is 0")
	ErrNoNodes                 = errors.New("no nodes in report")
	ErrInvalidPayersMerkleRoot = errors.New("payers merkle root is invalid")
	ErrMessageNotAtMinuteEnd   = errors.New(
		"sequence id is not the last message in the minute",
	)
	ErrMessageAtStartSequenceIDNotFound = errors.New("message at start sequence id not found")
	ErrMessageAtEndSequenceIDNotFound   = errors.New("message at end sequence id not found")
	ErrMerkleRootMismatch               = errors.New("payers merkle root mismatch")
)

type PayerReportVerifier struct {
	log   *zap.Logger
	store IPayerReportStore
}

func NewPayerReportVerifier(log *zap.Logger, store IPayerReportStore) *PayerReportVerifier {
	return &PayerReportVerifier{
		log:   log,
		store: store,
	}
}

/*
IsValidReport validates a payer report.

This function checks that the new report is valid and that it is a valid
transition from the previous report.

- The previous report is assumed to be valid, and does not get validated again.
  - Will regenerate the payer map and verify that the merkle root is correct.
    *
  - @param prevReport The previous report.
  - @param newReport The new report.
*/
func (p *PayerReportVerifier) IsValidReport(
	ctx context.Context,
	prevReport *PayerReport,
	newReport *PayerReport,
) (bool, error) {
	var err error

	log := AddReportLogFields(p.log, newReport)

	if err = validateReportTransition(prevReport, newReport); err != nil {
		log.Warn("invalid report transition", zap.Error(err))
		return false, nil
	}

	if err = validateReportStructure(newReport); err != nil {
		log.Warn("invalid report content", zap.Error(err))
		return false, nil
	}

	var isValidMerkleRoot bool
	// If the start and end sequence IDs are the same, the report is empty and the merkle root must always be the hash of an empty set
	if newReport.StartSequenceID == newReport.EndSequenceID {
		// TODO:nm validate that the merkle root is the hash of an empty set
		isValidMerkleRoot = true
	} else {
		if isValidMerkleRoot, err = p.verifyMerkleRoot(ctx, newReport); err != nil {
			return false, err
		}
	}
	if !isValidMerkleRoot {
		log.Warn("invalid merkle root")
	}

	return isValidMerkleRoot, nil
}

// Re-generates the payer map and verifies that the merkle root in the report matches the newly generated one
func (p *PayerReportVerifier) verifyMerkleRoot(
	ctx context.Context,
	report *PayerReport,
) (bool, error) {
	startEnvelope, endEnvelope, err := p.getStartAndEndMessages(ctx, report)
	if err != nil {
		return false, err
	}
	// If the start sequence ID is 0, it is the first report and we should start from minute 0 since there are no preceding reports
	var startMinute int32
	if startEnvelope == nil {
		startMinute = 0
	} else {
		startMinute, err = getMinuteFromEnvelope(startEnvelope)
		if err != nil {
			return false, err
		}
	}

	endMinute, err := getMinuteFromEnvelope(endEnvelope)
	if err != nil {
		return false, err
	}

	// Validate the report ends after it starts
	if startMinute > endMinute {
		return false, ErrInvalidReportStart
	}

	originatorID, err := utils.Uint32ToInt32(report.OriginatorNodeID)
	if err != nil {
		return false, ErrInvalidOriginatorID
	}

	isAtMinuteEnd, err := p.isAtMinuteEnd(
		ctx,
		originatorID,
		endMinute,
		int64(report.EndSequenceID),
	)
	if err != nil {
		return false, err
	}

	if !isAtMinuteEnd {
		return false, ErrMessageNotAtMinuteEnd
	}

	// TODO:nm validate that the start sequence ID is the last message in the start minute and create a misbehavior report if it's not
	querier := p.store.Queries()
	reportData, err := querier.BuildPayerReport(ctx, queries.BuildPayerReportParams{
		OriginatorID:           originatorID,
		StartMinutesSinceEpoch: startMinute,
		EndMinutesSinceEpoch:   endMinute,
	})
	if err != nil {
		return false, err
	}

	payerMap := buildPayersMap(reportData)
	merkleTree, err := generateMerkleTree(payerMap)
	if err != nil {
		return false, err
	}

	merkleRoot := common.BytesToHash(merkleTree.Root())
	if report.PayersMerkleRoot != merkleRoot {
		return false, ErrMerkleRootMismatch
	}

	return true, nil
}

// Check if a given sequence ID is the last message in a minute
func (p *PayerReportVerifier) isAtMinuteEnd(
	ctx context.Context,
	originatorID int32,
	minute int32,
	expectedSequenceID int64,
) (bool, error) {
	querier := p.store.Queries()

	lastSequenceID, err := querier.GetLastSequenceIDForOriginatorMinute(
		ctx,
		queries.GetLastSequenceIDForOriginatorMinuteParams{
			OriginatorID:      originatorID,
			MinutesSinceEpoch: minute,
		},
	)
	if err != nil {
		return false, err
	}

	isAtMinuteEnd := lastSequenceID == expectedSequenceID
	if !isAtMinuteEnd {
		p.log.Error(
			"sequence id is not the last message in the minute",
			zap.Int64("last_sequence_id", lastSequenceID),
			zap.Int32("minute", minute),
			zap.Int64("expected_sequence_id", expectedSequenceID),
		)
	}

	return isAtMinuteEnd, nil
}

func (p *PayerReportVerifier) getStartAndEndMessages(
	ctx context.Context,
	report *PayerReport,
) (*envelopes.OriginatorEnvelope, *envelopes.OriginatorEnvelope, error) {
	querier := p.store.Queries()
	startSequenceID, err := utils.Uint64ToInt64(report.StartSequenceID)
	if err != nil {
		return nil, nil, ErrInvalidSequenceID
	}

	originatorNodeID, err := utils.Uint32ToInt32(report.OriginatorNodeID)
	if err != nil {
		return nil, nil, ErrInvalidOriginatorID
	}

	var startEnvelope *envelopes.OriginatorEnvelope

	if startSequenceID != 0 {
		startMessage, err := querier.GetGatewayEnvelopeByID(
			ctx,
			queries.GetGatewayEnvelopeByIDParams{
				OriginatorSequenceID: startSequenceID,
				OriginatorNodeID:     originatorNodeID,
			},
		)
		if err != nil {
			return nil, nil, ErrMessageAtStartSequenceIDNotFound
		}

		startEnvelope, err = envelopes.NewOriginatorEnvelopeFromBytes(
			startMessage.OriginatorEnvelope,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	endSequenceID, err := utils.Uint64ToInt64(report.EndSequenceID)
	if err != nil {
		return nil, nil, ErrInvalidSequenceID
	}

	endMessage, err := querier.GetGatewayEnvelopeByID(ctx, queries.GetGatewayEnvelopeByIDParams{
		OriginatorSequenceID: endSequenceID,
		OriginatorNodeID:     originatorNodeID,
	})
	if err != nil {
		return nil, nil, ErrMessageAtEndSequenceIDNotFound
	}

	endEnvelope, err := envelopes.NewOriginatorEnvelopeFromBytes(endMessage.OriginatorEnvelope)
	if err != nil {
		return nil, nil, err
	}

	return startEnvelope, endEnvelope, nil
}

// Static validations on the report transition
func validateReportTransition(prevReport *PayerReport, newReport *PayerReport) error {
	// Special validations for the first report
	if prevReport == nil {
		if newReport.StartSequenceID != 0 {
			return ErrInvalidReportStart
		}

		return nil
	}
	// Check if the reports are referring to the same originator.
	// This is a sanity check. Mismatched reports should never make it this far.
	if prevReport.OriginatorNodeID != newReport.OriginatorNodeID {
		return ErrMismatchOriginator
	}

	// Check if the new report starts where the previous report ended.
	// This is a sanity check. These should be filtered out first
	if prevReport.EndSequenceID != newReport.StartSequenceID {
		return ErrInvalidReportStart
	}

	return nil
}

// Validates that the report is well-formed and doesn't have any logical
// errors or invalid fields that can be detected without further processing.
func validateReportStructure(report *PayerReport) error {
	// The Originator Node ID is required
	if report.OriginatorNodeID == 0 {
		return ErrInvalidOriginatorID
	}

	if len(report.ActiveNodeIDs) == 0 {
		return ErrNoNodes
	}

	// The payers merkle root is required. It may be set to the hash of an empty set
	// if there are no payers in the report.
	if len(report.PayersMerkleRoot) != 32 {
		return ErrInvalidPayersMerkleRoot
	}

	// Check if the new report ends after it starts
	if report.StartSequenceID > report.EndSequenceID {
		return ErrInvalidReportStart
	}

	return nil
}

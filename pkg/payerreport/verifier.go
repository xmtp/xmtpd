package payerreport

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrMismatchOriginator      = errors.New("originator id mismatch between old and new report")
	ErrInvalidReportStart      = errors.New("report does not start where the previous report ended")
	ErrInvalidOriginatorID     = errors.New("originator id is 0")
	ErrNoNodes                 = errors.New("no nodes in report")
	ErrInvalidNodesHash        = errors.New("nodes hash is invalid")
	ErrInvalidPayersMerkleRoot = errors.New("payers merkle root is invalid")
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
  - Validate a report transition. Returns a bool if the report can be conclusively validated/rejected.
  - Otherwise returns an error.
    *
  - This function checks that the new report is valid and that it is a valid
  - transition from the previous report.
  - The previous report is assumed to be valid, and does not get validated again.
    *
  - @param prevReport The previous report.
  - @param newReport The new report.
*/
func (p *PayerReportVerifier) IsValidReport(
	ctx context.Context,
	prevReport *PayerReport,
	newReport *PayerReport,
) (bool, error) {
	newReportID, err := newReport.ID()
	if err != nil {
		p.log.Error("failed to get report id", zap.Error(err))
		return false, nil
	}

	log := p.log.With(
		zap.String("new_report_id", newReportID.String()),
	)

	if err := validateReportTransition(prevReport, newReport); err != nil {
		log.Warn("invalid report transition", zap.Error(err))
		return false, nil
	}

	if err := validateReportStructure(newReport); err != nil {
		log.Warn("invalid report content", zap.Error(err))
		return false, nil
	}

	return true, nil
}

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

package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const submitterWorkerID = 2

type SubmitterWorker struct {
	logger           *zap.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               *sync.WaitGroup
	payerReportStore payerreport.IPayerReportStore
	registry         registry.NodeRegistry
	reportsAdmin     blockchain.PayerReportsManager
	myNodeID         uint32
}

var _ stoppable = &SubmitterWorker{}

func NewSubmitterWorker(
	ctx context.Context,
	logger *zap.Logger,
	payerReportStore payerreport.IPayerReportStore,
	registry registry.NodeRegistry,
	reportsManager blockchain.PayerReportsManager,
	myNodeID uint32,
) *SubmitterWorker {
	ctx, cancel := context.WithCancel(ctx)

	return &SubmitterWorker{
		logger:           logger.Named(utils.PayerReportSubmitterWorkerLoggerName),
		ctx:              ctx,
		cancel:           cancel,
		wg:               &sync.WaitGroup{},
		payerReportStore: payerReportStore,
		registry:         registry,
		myNodeID:         myNodeID,
		reportsAdmin:     reportsManager,
	}
}

func (w *SubmitterWorker) Start() {
	tracing.GoPanicWrap(w.ctx, w.wg, "payer-report-submitter-worker", func(ctx context.Context) {
		for {
			nextRun := findNextRunTime(w.myNodeID, submitterWorkerID)
			wait := time.Until(nextRun)
			select {
			case <-w.ctx.Done():
				return
			case <-time.After(wait):
				if err := w.SubmitReports(ctx); err != nil {
					w.logger.Error("error submitting reports", zap.Error(err))
				}
			}
		}
	})
}

func (w *SubmitterWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}

// SubmitReports is the main loop of the submitter worker.
// The submitter worker fetches all reports that are pending submission and approved attestation.
// Note:
// All reports are fetched independently of the originator node ID. This means the submitter:
//   - will try to submit reports for other nodes if they are pending submission and approved attestation.
//   - works on a loop that gets activated every `findNextRunTime(originatorNodeID, submitterWorkerID)` minutes.
//     this distribution guarantees that no two nodes will submit reports for the same originator node at the same time.
//     the blockchain guarantees deduplication of report submissions.
//   - even with `findNextRunTime` the system has to guarantee that no duplicates are submitted.
func (w *SubmitterWorker) SubmitReports(ctx context.Context) error {
	haLock, err := w.payerReportStore.GetAdvisoryLocker(w.ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = haLock.Release()
	}()

	err = haLock.LockSubmitterWorker()
	if err != nil {
		return err
	}

	// Fetch all reports that are pending submission and approved attestation.
	w.logger.Debug("fetching reports to submit")
	reports, err := w.payerReportStore.FetchReports(
		ctx,
		payerreport.NewFetchReportsQuery().
			WithSubmissionStatus(payerreport.SubmissionPending).
			WithAttestationStatus(payerreport.AttestationApproved),
	)
	if err != nil {
		return err
	}

	if len(reports) == 0 {
		w.logger.Debug("no reports to submit, skipping")
		return nil
	}

	var latestErr error

	for _, report := range reports {
		reportLogger := payerreport.AddReportLogFields(w.logger, &report.PayerReport)

		requiredAttestations := (len(report.ActiveNodeIDs) / 2) + 1

		// Only handle reports that have the required number of approved attestations.
		if len(report.AttestationSignatures) < requiredAttestations {
			continue
		}

		reportLogger.Info("submitting report")
		// Submit the report to the blockchain
		reportIndex, submitErr := w.submitReport(report)
		if submitErr != nil {
			// If the on-chain protocol throws a PayerReportAlreadySubmitted error,
			// it means the same report was already submitted by another node.
			// We set the report submission status to submitted and continue.
			if submitErr.IsErrPayerReportAlreadySubmitted() {
				reportLogger.Info("report already submitted, skipping")

				err = w.payerReportStore.SetReportSubmitted(ctx, report.ID, reportIndex)
				if err != nil {
					reportLogger.Error(
						"failed to set report submitted",
						utils.PayerReportIDField(report.ID.String()),
						zap.Error(err),
					)
					latestErr = err
				}

				continue
			}

			// If the on-chain protocol throws an InvalidSequenceIDs or InvalidStartSequenceID error,
			// it means the report has invalid sequence IDs or start sequence ID.
			// Most likely another node submitted a valid report with the same start sequence ID before the node.
			// We set the report submission status to rejected and continue.
			if submitErr.IsErrInvalidSequenceIDs() {
				reportLogger.Info("report has invalid sequence IDs, submission rejected")

				err = w.payerReportStore.SetReportSubmissionRejected(ctx, report.ID)
				if err != nil {
					reportLogger.Error(
						"failed to set report submission rejected",
						utils.PayerReportIDField(report.ID.String()),
						zap.Error(err),
					)
					latestErr = err
				}

				continue
			}

			reportLogger.Error(
				"failed to submit report",
				utils.PayerReportIDField(report.ID.String()),
				zap.Error(submitErr),
			)

			latestErr = submitErr
			continue
		}

		// NOTE: there is a possible race when the indexer hears about the event before we get a confirmation from the chain
		// Since we are not holding a lock, the report might already end up being submitted by the time we get here
		// SetReportSubmitted should be able to handle that
		err = w.payerReportStore.SetReportSubmitted(ctx, report.ID, reportIndex)
		if err != nil {
			reportLogger.Warn(
				"failed to set report submitted",
				utils.PayerReportIDField(report.ID.String()),
			)
		}
	}

	return latestErr
}

func (w *SubmitterWorker) submitReport(
	report *payerreport.PayerReportWithStatus,
) (int32, blockchain.ProtocolError) {
	return w.reportsAdmin.SubmitPayerReport(w.ctx, report)
}

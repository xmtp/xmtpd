package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const submitterWorkerID = 2

type SubmitterWorker struct {
	log              *zap.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               *sync.WaitGroup
	payerReportStore payerreport.IPayerReportStore
	registry         registry.NodeRegistry
	reportsAdmin     blockchain.PayerReportsManager
	myNodeID         uint32
}

func NewSubmitterWorker(
	ctx context.Context,
	log *zap.Logger,
	payerReportStore payerreport.IPayerReportStore,
	registry registry.NodeRegistry,
	reportsManager blockchain.PayerReportsManager,
	myNodeID uint32,
) *SubmitterWorker {
	ctx, cancel := context.WithCancel(ctx)
	return &SubmitterWorker{
		log:              log.Named("reportsubmitter"),
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
	tracing.GoPanicWrap(w.ctx, w.wg, "payerreport-submitter", func(ctx context.Context) {
		for {
			nextRun := findNextRunTime(w.myNodeID, submitterWorkerID)
			wait := time.Until(nextRun)
			select {
			case <-w.ctx.Done():
				return
			case <-time.After(wait):
				if err := w.SubmitReports(ctx); err != nil {
					w.log.Error("error submitting reports", zap.Error(err))
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

	reports, err := w.payerReportStore.FetchReports(
		ctx,
		payerreport.NewFetchReportsQuery().WithSubmissionStatus(payerreport.SubmissionPending),
	)
	if err != nil {
		return err
	}

	var latestErr error

	for _, report := range reports {
		reportLogger := payerreport.AddReportLogFields(w.log, &report.PayerReport)

		requiredAttestations := (len(report.ActiveNodeIDs) / 2) + 1
		if len(report.AttestationSignatures) < requiredAttestations {
			continue
		}

		reportLogger.Info("submitting report")
		submitErr := w.submitReport(report)
		if submitErr != nil {
			reportLogger.Error(
				"failed to submit report",
				zap.String("report_id", report.ID.String()),
				zap.Error(submitErr),
			)

			latestErr = submitErr

			if submitErr.IsErrInvalidSequenceIDs() {
				reportLogger.Info("report has invalid sequence IDs, submission rejected")
				w.payerReportStore.SetReportSubmissionRejected(ctx, report.ID)
				continue
			}

			continue
		}

		// NOTE: there is a possible race when the indexer hears about the event before we get a confirmation from the chain
		// Since we are not holding a lock, the report might already end up being submitted by the time we get here
		// SetReportSubmitted should be able to handle that
		err = w.payerReportStore.SetReportSubmitted(ctx, report.ID)
		if err != nil {
			reportLogger.Warn(
				"failed to set report submitted",
				zap.String("report_id", report.ID.String()),
			)
		}
	}

	return latestErr
}

func (w *SubmitterWorker) submitReport(
	report *payerreport.PayerReportWithStatus,
) blockchain.ProtocolError {
	return w.reportsAdmin.SubmitPayerReport(w.ctx, report)
}

package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const settlementWorkerID = 3

type SettlementWorker struct {
	log              *zap.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               *sync.WaitGroup
	payerReportStore payerreport.IPayerReportStore
	myNodeID         uint32
}

func NewSettlementWorker(
	ctx context.Context,
	log *zap.Logger,
	payerReportStore payerreport.IPayerReportStore,
	myNodeID uint32,
) *SettlementWorker {
	ctx, cancel := context.WithCancel(ctx)
	return &SettlementWorker{
		log:              log.Named("reportsettlement"),
		ctx:              ctx,
		cancel:           cancel,
		wg:               &sync.WaitGroup{},
		payerReportStore: payerReportStore,
		myNodeID:         myNodeID,
	}
}

func (w *SettlementWorker) Start() {
	tracing.GoPanicWrap(w.ctx, w.wg, "payerreport-settlement", func(ctx context.Context) {
		for {
			nextRun := findNextRunTime(w.myNodeID, settlementWorkerID)
			wait := time.Until(nextRun)
			select {
			case <-w.ctx.Done():
				return
			case <-time.After(wait):
				if err := w.SettleReports(ctx); err != nil {
					w.log.Error("error settling reports", zap.Error(err))
				}
			}
		}
	})
}

func (w *SettlementWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}

func (w *SettlementWorker) SettleReports(ctx context.Context) error {
	haLock, err := w.payerReportStore.GetAdvisoryLocker(w.ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = haLock.Release()
	}()

	err = haLock.LockSettlementWorker()
	if err != nil {
		return err
	}

	// SettlementWorker fetches all reports that have been submitted but not yet settled.
	w.log.Debug("fetching reports to settle")
	reports, err := w.payerReportStore.FetchReports(
		ctx,
		payerreport.NewFetchReportsQuery().
			WithSubmissionStatus(payerreport.SubmissionSubmitted),
	)
	if err != nil {
		return err
	}

	var latestErr error

	for _, report := range reports {
		reportLogger := payerreport.AddReportLogFields(w.log, &report.PayerReport)

		reportLogger.Info("settling report")
		settleErr := w.settleReport(report)
		if settleErr != nil {
			reportLogger.Error(
				"failed to settle report",
				zap.String("report_id", report.ID.String()),
				zap.Error(settleErr),
			)

			latestErr = settleErr
			continue
		}

		err = w.payerReportStore.SetReportSettled(ctx, report.ID)
		if err != nil {
			reportLogger.Warn(
				"failed to set report settled",
				zap.String("report_id", report.ID.String()),
			)
		}
	}

	return latestErr
}

func (w *SettlementWorker) settleReport(
	report *payerreport.PayerReportWithStatus,
) error {
	// TODO: Implement settlement logic
	return nil
}

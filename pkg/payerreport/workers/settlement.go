package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	settlementWorkerID = 3
	MAX_PROOF_ELEMENTS = 100
)

type SettlementWorker struct {
	log                *zap.Logger
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 *sync.WaitGroup
	payerReportStore   payerreport.IPayerReportStore
	verifier           payerreport.IPayerReportVerifier
	reportManager      blockchain.PayerReportsManager
	myNodeID           uint32
	submissionNotifyCh <-chan struct{}
}

func NewSettlementWorker(
	ctx context.Context,
	log *zap.Logger,
	payerReportStore payerreport.IPayerReportStore,
	verifier payerreport.IPayerReportVerifier,
	reportManager blockchain.PayerReportsManager,
	myNodeID uint32,
	submissionNotifyCh <-chan struct{},
) *SettlementWorker {
	ctx, cancel := context.WithCancel(ctx)
	return &SettlementWorker{
		log:                log.Named("reportsettlement"),
		ctx:                ctx,
		cancel:             cancel,
		wg:                 &sync.WaitGroup{},
		payerReportStore:   payerReportStore,
		myNodeID:           myNodeID,
		verifier:           verifier,
		reportManager:      reportManager,
		submissionNotifyCh: submissionNotifyCh,
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
			case <-w.submissionNotifyCh:
				w.log.Debug("received submission notification, settling reports")
				if err := w.SettleReports(ctx); err != nil {
					w.log.Error("error settling reports after submission", zap.Error(err))
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
		settleErr := w.settleReport(ctx, report)
		if settleErr != nil {
			reportLogger.Error(
				"failed to settle report",
				zap.String("report_id", report.ID.String()),
				zap.Error(settleErr),
			)

			latestErr = settleErr
			continue
		} else {
			reportLogger.Info("report settled")
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
	ctx context.Context,
	report *payerreport.PayerReportWithStatus,
) error {
	reportLogger := payerreport.AddReportLogFields(w.log, &report.PayerReport)

	if report.SubmittedReportIndex == nil {
		return errors.New("report index is nil")
	}

	payerMap, err := w.verifier.GetPayerMap(ctx, &report.PayerReport)
	if err != nil {
		return err
	}

	merkleTree, err := payerreport.GenerateMerkleTree(payerMap)
	if err != nil {
		return err
	}

	summary, err := w.reportManager.SettlementSummary(
		ctx,
		report.OriginatorNodeID,
		uint64(*report.SubmittedReportIndex),
	)
	if err != nil {
		return err
	}

	// Return early if the report is already settled
	if summary.IsSettled {
		return nil
	}

	index := uint64(*report.SubmittedReportIndex)

	offset := int(summary.Offset)
	leafCount := merkleTree.LeafCount()
	remaining := leafCount - offset

	// This should never happen, but check just to be safe.
	if offset > leafCount {
		return errors.New("offset exceeds leaf count")
	}

	if remaining == 0 && !summary.IsSettled {
		reportLogger.Warn(
			"something is fishy. No items left to settle but report is not settled",
		)
	}

	for remaining > 0 {
		numElements := utils.MinInt(remaining, MAX_PROOF_ELEMENTS)
		proof, err := merkleTree.GenerateMultiProofSequential(offset, numElements)
		if err != nil {
			return err
		}

		err = w.reportManager.SettleReport(
			ctx,
			report.OriginatorNodeID,
			index,
			proof,
		)
		if err != nil {
			return err
		}

		offset += numElements
		remaining -= numElements

		reportLogger.Info(
			"settled subset of report",
			zap.Int("remaining_items", remaining),
			zap.Int("num_settled", numElements),
		)
	}

	return nil
}

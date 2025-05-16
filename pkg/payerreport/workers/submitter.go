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
	myNodeID         uint32
	reportsAdmin     blockchain.PayerReportsAdmin
}

func NewSubmitterWorker(
	ctx context.Context,
	log *zap.Logger,
	payerReportStore payerreport.IPayerReportStore,
	registry registry.NodeRegistry,
	myNodeID uint32,
	reportsAdmin blockchain.PayerReportsAdmin,
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
		reportsAdmin:     reportsAdmin,
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
				if err := w.submitReports(ctx); err != nil {
					w.log.Error("error submitting reports", zap.Error(err))
				}
			}
		}
	})
}

func (w *SubmitterWorker) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}

func (w *SubmitterWorker) submitReports(ctx context.Context) error {
	allNodes, err := w.registry.GetNodes()
	if err != nil {
		return err
	}

	submissionThreshold := int32((len(allNodes) / 2) + 1)

	reports, err := w.payerReportStore.FetchReports(
		ctx,
		payerreport.NewFetchReportsQuery().WithSubmissionStatus(payerreport.SubmissionPending).
			WithMinAttestations(submissionThreshold),
	)
	if err != nil {
		return err
	}

	for _, report := range reports {
		w.log.Info("submitting report", zap.String("reportID", report.ID.String()))
		if err = w.submitReport(report); err != nil {
			w.log.Error(
				"failed to submit report",
				zap.String("report_id", report.ID.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (w *SubmitterWorker) submitReport(report *payerreport.PayerReportWithStatus) error {
	return w.reportsAdmin.SubmitPayerReport(w.ctx, report)
}

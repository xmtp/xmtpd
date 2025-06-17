package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

type GeneratorWorker struct {
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	log               *zap.Logger
	store             payerreport.IPayerReportStore
	generator         payerreport.IPayerReportGenerator
	registry          registry.NodeRegistry
	registrant        registrant.IRegistrant
	myNodeID          uint32
	minReportInterval time.Duration
}

func NewGeneratorWorker(
	ctx context.Context,
	log *zap.Logger,
	store payerreport.IPayerReportStore,
	registry registry.NodeRegistry,
	registrant registrant.IRegistrant,
	minReportInterval time.Duration,
) *GeneratorWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &GeneratorWorker{
		ctx:               ctx,
		cancel:            cancel,
		log:               log.Named("generatorworker"),
		store:             store,
		generator:         payerreport.NewPayerReportGenerator(log, store.Queries(), registry),
		registry:          registry,
		registrant:        registrant,
		myNodeID:          registrant.NodeID(),
		minReportInterval: minReportInterval,
	}
	return worker
}

func (w *GeneratorWorker) Start() {
	tracing.GoPanicWrap(
		w.ctx,
		&w.wg,
		"generator-worker",
		func(ctx context.Context) {
			w.log.Info("Starting generator worker")

			for {
				nextRun := findNextRunTime(w.myNodeID)
				wait := time.Until(nextRun)
				select {
				case <-time.After(wait):
					if err := w.GenerateReports(); err != nil {
						w.log.Error("generating reports", zap.Error(err))
					}
				case <-ctx.Done():
					return
				}
			}
		},
	)
}

func (w *GeneratorWorker) Stop() {
	w.cancel()
}

func (w *GeneratorWorker) GenerateReports() error {
	w.log.Info("getting nodes from registry", zap.Any("registry", w.registry))
	allNodes, err := w.registry.GetNodes()
	if err != nil {
		return err
	}

	for _, node := range allNodes {
		if err = w.maybeGenerateReport(node.NodeID); err != nil {
			w.log.Warn(
				"error generating report for node",
				zap.Uint32("node_id", node.NodeID),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (w *GeneratorWorker) maybeGenerateReport(nodeID uint32) error {
	lastSubmittedReport, err := w.getLastSubmittedReport(nodeID)
	if err != nil {
		return err
	}

	// Only continue if the last submitted report doesn't exist (there have been no reports) or if it is older than the minimum report interval
	if lastSubmittedReport != nil && !w.isOlderThanReportInterval(lastSubmittedReport) {
		w.log.Debug("skipping report generation for node", zap.Uint32("node_id", nodeID))
		return nil
	}

	existingReportEndSequenceID := uint64(0)
	if lastSubmittedReport != nil {
		existingReportEndSequenceID = lastSubmittedReport.EndSequenceID
	}

	// Fetch all reports for the originator that are pending or
	existingReports, err := w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
			WithOriginatorNodeID(nodeID).
			// Ignore existing reports that were rejected or not yet attested
			WithAttestationStatus(payerreport.AttestationApproved).
			WithSubmissionStatus(payerreport.SubmissionPending).
			// We are looking for reports that start at the end of the last submitted report
			WithStartSequenceID(existingReportEndSequenceID),
	)
	if err != nil {
		return err
	}

	if len(existingReports) > 0 {
		w.log.Debug(
			"skipping report generation for node because there are existing reports pending",
			zap.Uint32("node_id", nodeID),
			zap.Int("num_existing_reports", len(existingReports)),
		)
		return nil
	}

	return w.generateReport(nodeID, existingReportEndSequenceID)
}

func (w *GeneratorWorker) generateReport(nodeID uint32, lastReportEndSequenceID uint64) error {
	report, err := w.generator.GenerateReport(w.ctx, payerreport.PayerReportGenerationParams{
		OriginatorID:            nodeID,
		LastReportEndSequenceID: lastReportEndSequenceID,
	})
	if err != nil {
		return err
	}

	clientEnvelope, err := report.ToClientEnvelope()
	if err != nil {
		return err
	}

	clientEnvelopeBytes, err := clientEnvelope.Bytes()
	if err != nil {
		return err
	}

	payerSignature, err := w.registrant.SignClientEnvelopeToSelf(clientEnvelopeBytes)
	if err != nil {
		return err
	}

	payerEnvelope, err := envelopes.NewPayerEnvelope(&envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     w.myNodeID,
		MessageRetentionDays: constants.DEFAULT_STORAGE_DURATION_DAYS,
	})
	if err != nil {
		return err
	}

	reportID, err := w.store.CreatePayerReport(w.ctx, &report.PayerReport, payerEnvelope)
	if err != nil {
		return err
	}

	w.log.Info("generated report", zap.String("report_id", reportID.String()))

	return nil
}

func (w *GeneratorWorker) getLastSubmittedReport(
	nodeID uint32,
) (*payerreport.PayerReportWithStatus, error) {
	reports, err := w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
			WithOriginatorNodeID(nodeID).
			WithSubmissionStatus(payerreport.SubmissionSubmitted, payerreport.SubmissionSettled),
	)
	if err != nil {
		return nil, err
	}

	var latestReport *payerreport.PayerReportWithStatus
	for _, report := range reports {
		if latestReport == nil || report.EndSequenceID > latestReport.EndSequenceID {
			latestReport = report
		}
	}

	return latestReport, nil
}

func (w *GeneratorWorker) isOlderThanReportInterval(
	report *payerreport.PayerReportWithStatus,
) bool {
	// Convert the report's end minute since epoch to a time.Time
	reportEndTime := time.Unix(int64(report.EndMinuteSinceEpoch)*60, 0).UTC()

	// Check if the report is older than the minimum report interval
	return time.Now().UTC().Sub(reportEndTime) > w.minReportInterval
}

package workers

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const generatorWorkerID = 1

type GeneratorWorker struct {
	ctx                  context.Context
	cancel               context.CancelFunc
	wg                   sync.WaitGroup
	logger               *zap.Logger
	store                payerreport.IPayerReportStore
	generator            payerreport.IPayerReportGenerator
	registry             registry.NodeRegistry
	registrant           registrant.IRegistrant
	myNodeID             uint32
	generateSelfPeriod   time.Duration
	generateOthersPeriod time.Duration
}

func NewGeneratorWorker(
	ctx context.Context,
	logger *zap.Logger,
	store payerreport.IPayerReportStore,
	registry registry.NodeRegistry,
	registrant registrant.IRegistrant,
	domainSeparator common.Hash,
	generateSelfPeriod time.Duration,
	generateOthersPeriod time.Duration,
) *GeneratorWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &GeneratorWorker{
		ctx:    ctx,
		cancel: cancel,
		logger: logger.Named(utils.PayerReportGeneratorWorkerLoggerName),
		store:  store,
		generator: payerreport.NewPayerReportGenerator(
			logger,
			store.Queries(),
			registry,
			domainSeparator,
		),
		registry:             registry,
		registrant:           registrant,
		myNodeID:             registrant.NodeID(),
		generateSelfPeriod:   generateSelfPeriod,
		generateOthersPeriod: generateOthersPeriod,
	}
	return worker
}

func (w *GeneratorWorker) Start() {
	tracing.GoPanicWrap(
		w.ctx,
		&w.wg,
		"payer-report-generator-worker",
		func(ctx context.Context) {
			w.logger.Info("starting")

			for {
				nextRun := findNextRunTime(w.myNodeID, generatorWorkerID)
				wait := time.Until(nextRun)
				select {
				case <-time.After(wait):
					if err := w.GenerateReports(); err != nil {
						w.logger.Error("generating reports", zap.Error(err))
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
	haLock, err := w.store.GetAdvisoryLocker(w.ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = haLock.Release()
	}()

	err = haLock.LockGeneratorWorker()
	if err != nil {
		return err
	}

	w.logger.Info("getting nodes from registry")
	allNodes, err := w.registry.GetNodes()
	if err != nil {
		return err
	}

	for _, node := range allNodes {
		if err = w.maybeGenerateReport(node.NodeID); err != nil {
			w.logger.Warn(
				"error generating report for node",
				utils.OriginatorIDField(node.NodeID),
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
		w.logger.Debug(
			"skipping report generation for node because the last submitted report is not too old",
			utils.OriginatorIDField(nodeID),
		)
		return nil
	}

	existingReportEndSequenceID := uint64(0)
	if lastSubmittedReport != nil {
		existingReportEndSequenceID = lastSubmittedReport.EndSequenceID
	}

	// TODO(mkysel) there is a timing hazard here
	// A concurrent submitter might transition a report from pending to submitted.
	// getLastSubmittedReport checks the latest submitted report, FetchReports looks for a pending/unsubmitted one
	// A concurrent submitter can transition from pending to submitted between these two calls
	// This will result in us missing the new report and generating a duplicate

	// Fetch all reports for the originator that are pending and approved
	w.logger.Debug(
		"maybe generating report, fetching existing reports",
		utils.OriginatorIDField(nodeID),
		utils.LastSequenceIDField(int64(existingReportEndSequenceID)),
	)

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

	validReports := make([]*payerreport.PayerReportWithStatus, 0, len(existingReports))

	if len(existingReports) > 0 {
		// Validate existing reports and expire old ones.
		for _, report := range existingReports {
			if w.isOlderThanReportInterval(report) {
				w.logger.Debug(
					"expiring old report",
					utils.OriginatorIDField(nodeID),
				)

				err = w.store.SetReportSubmissionRejected(w.ctx, report.ID)
				if err != nil {
					return err
				}
			} else {
				validReports = append(validReports, report)
			}
		}
	}

	if len(validReports) > 0 {
		w.logger.Debug(
			"skipping report generation for node because there are existing valid reports pending",
			utils.OriginatorIDField(nodeID),
			utils.CountField(int64(len(validReports))),
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
		MessageRetentionDays: constants.DefaultStorageDurationDays,
	})
	if err != nil {
		return err
	}

	reportID, err := w.store.CreatePayerReport(w.ctx, &report.PayerReport, payerEnvelope)
	if err != nil {
		return err
	}

	w.logger.Info("generated report", utils.PayerReportIDField(reportID.String()))

	return nil
}

func (w *GeneratorWorker) getLastSubmittedReport(
	nodeID uint32,
) (*payerreport.PayerReportWithStatus, error) {
	w.logger.Debug("fetching last submitted report",
		utils.OriginatorIDField(nodeID),
	)
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

	if report.OriginatorNodeID == w.registrant.NodeID() {
		return time.Now().UTC().Sub(reportEndTime) > w.generateSelfPeriod
	} else {
		return time.Now().UTC().Sub(reportEndTime) > w.generateOthersPeriod
	}
}

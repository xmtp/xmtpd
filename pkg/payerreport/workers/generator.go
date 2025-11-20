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
	expirySelfPeriod     time.Duration
	expiryOthersPeriod   time.Duration
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
	expirySelfPeriod time.Duration,
	expiryOthersPeriod time.Duration,
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
		expirySelfPeriod:     expirySelfPeriod,
		expiryOthersPeriod:   expiryOthersPeriod,
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
	w.logger.Debug(
		"maybe generating report, fetching existing reports snapshot",
		utils.OriginatorIDField(nodeID),
	)

	// Snapshot: all non-rejected reports for this originator in active states.
	reports, err := w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
			WithOriginatorNodeID(nodeID).
			WithSubmissionStatus(
				payerreport.SubmissionPending,
				payerreport.SubmissionSubmitted,
				payerreport.SubmissionSettled,
			),
	)
	if err != nil {
		return err
	}

	// Derive the "last submitted" (Submitted/Settled) report from the same snapshot.
	var lastSubmittedReport *payerreport.PayerReportWithStatus
	for _, report := range reports {
		// Only consider submitted and settled reports.
		if report.SubmissionStatus != payerreport.SubmissionSubmitted &&
			report.SubmissionStatus != payerreport.SubmissionSettled {
			continue
		}

		// Update the last submitted report if it's the highest end sequence ID.
		if lastSubmittedReport == nil || report.EndSequenceID > lastSubmittedReport.EndSequenceID {
			lastSubmittedReport = report
		}
	}

	// Only continue and generate a new report if there have been no reports yet OR the last submitted
	// report is past the minimum generation threshold.
	// Otherwise, wait until on going report resolution is complete.
	if lastSubmittedReport != nil && !w.isPastGenerationThreshold(lastSubmittedReport) {
		w.logger.Debug(
			"skipping report generation for node because the last submitted report is within the generation interval",
			utils.OriginatorIDField(nodeID),
		)
		return nil
	}

	// The generation threshold has passed, so we can generate a new report.
	// Get the last sequence ID of the last submitted report.
	// If there is no last submitted report, use 0 as the existing report end sequence ID.
	existingReportEndSequenceID := uint64(0)
	if lastSubmittedReport != nil {
		existingReportEndSequenceID = lastSubmittedReport.EndSequenceID
	}

	w.logger.Debug(
		"maybe generating report, checking for existing report at boundary",
		utils.OriginatorIDField(nodeID),
		utils.LastSequenceIDField(int64(existingReportEndSequenceID)),
	)

	// From the same snapshot:
	//   - Expire reports that are too old: consider only pending reports for this.
	//   - Collect valid (non-expired) reports that start exactly at the boundary.
	//   - A valid report is one that is submitted or settled, and:
	//     - Our local attestation status does not matter.
	//     - Even if we rejected the report, others might submit and settle it if there is sufficient consensus.
	//     - All other attestation status states are managed by the Attestation Worker and not relevant to generation.
	validReports := make([]*payerreport.PayerReportWithStatus, 0, len(reports))
	for _, report := range reports {
		if report.SubmissionStatus == payerreport.SubmissionPending && w.isReportExpired(report) {
			w.logger.Debug(
				"expiring old report",
				utils.OriginatorIDField(nodeID),
				utils.PayerReportIDField(report.ID.String()),
			)

			if err := w.store.SetReportSubmissionRejected(w.ctx, report.ID); err != nil {
				return err
			}

			continue
		}

		// If there's a submitted, settled or non-expired pending report at the boundary, skip generation.
		if report.StartSequenceID == existingReportEndSequenceID {
			validReports = append(validReports, report)
		}
	}

	if len(validReports) > 0 {
		w.logger.Debug(
			"skipping report generation for node because there are existing valid reports pending",
			utils.OriginatorIDField(nodeID),
			utils.LastSequenceIDField(int64(existingReportEndSequenceID)),
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

	if report == nil {
		return nil
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

	w.logger.Info("generated report",
		utils.PayerReportIDField(reportID.String()),
		utils.StartSequenceIDField(int64(report.StartSequenceID)),
		utils.LastSequenceIDField(int64(report.EndSequenceID)),
		utils.OriginatorIDField(report.OriginatorNodeID),
	)

	return nil
}

// isPastGenerationThreshold reports whether enough time has passed since the
// end of the given report to allow generating a new one.
//
// Each originator node has a minimum report generation interval:
//   - For the local node (w.registrant.NodeID), the interval is w.generateSelfPeriod.
//   - For all other nodes, the interval is w.generateOthersPeriod.
//
// A report is considered "past the generation threshold" when the duration since
// its EndMinuteSinceEpoch exceeds the applicable interval. When this function
// returns true, the caller may proceed to generate the next report for that
// originator node.
func (w *GeneratorWorker) isPastGenerationThreshold(
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

// isReportExpired determines whether the given report is considered expired
// based on how long ago it ended.
//
// Expiration serves two purposes:
//  1. Healing stalled consensus: If the network fails to reach consensus on a
//     report for an extended period, older reports should eventually be ignored
//     so the system can move forward and recover.
//  2. Preventing accumulation of stale work: Nodes should not continue acting
//     on or submitting reports that are far outside the normal operational
//     window.
//
// The expiration threshold depends on whether the report was produced by:
//   - the local node       → w.expirySelfPeriod
//   - another originator   → w.expiryOthersPeriod
//
// A report is considered expired once the elapsed time since its
// EndMinuteSinceEpoch exceeds the applicable threshold.
func (w *GeneratorWorker) isReportExpired(
	report *payerreport.PayerReportWithStatus,
) bool {
	reportEndTime := time.Unix(int64(report.EndMinuteSinceEpoch)*60, 0).UTC()

	if report.OriginatorNodeID == w.registrant.NodeID() {
		return time.Now().UTC().Sub(reportEndTime) > w.expirySelfPeriod
	} else {
		return time.Now().UTC().Sub(reportEndTime) > w.expiryOthersPeriod
	}
}

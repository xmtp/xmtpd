// Package workers implements the attestation worker for the payer report.
package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// AttestationWorker is responsible for periodically checking for reports that need attestation
// and signing them with the node's private key.
type AttestationWorker struct {
	ctx             context.Context
	cancel          context.CancelFunc
	logger          *zap.Logger
	registrant      registrant.IRegistrant
	store           payerreport.IPayerReportStore
	verifier        payerreport.IPayerReportVerifier
	wg              sync.WaitGroup
	pollInterval    time.Duration
	domainSeparator common.Hash
}

// NewAttestationWorker creates and starts a new attestation worker that will periodically
// check for reports that need attestation.
// It takes a context, logger, registrant for signing, store for accessing reports,
// and a poll interval that determines how often to check for reports.
func NewAttestationWorker(
	ctx context.Context,
	logger *zap.Logger,
	registrant registrant.IRegistrant,
	store payerreport.IPayerReportStore,
	pollInterval time.Duration,
	domainSeparator common.Hash,
) *AttestationWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &AttestationWorker{
		ctx:             ctx,
		logger:          logger.Named(utils.PayerReportAttestationWorkerLoggerName),
		registrant:      registrant,
		store:           store,
		verifier:        payerreport.NewPayerReportVerifier(logger, store),
		wg:              sync.WaitGroup{},
		cancel:          cancel,
		pollInterval:    pollInterval,
		domainSeparator: domainSeparator,
	}

	return worker
}

// Start launches the worker's main loop in a separate goroutine.
// The loop periodically checks for reports that need attestation.
func (w *AttestationWorker) Start() {
	tracing.GoPanicWrap(
		w.ctx,
		&w.wg,
		"payer-report-attestation-worker",
		func(ctx context.Context) {
			w.logger.Info("starting")
			var err error
			ticker := time.NewTicker(w.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err = w.AttestReports(); err != nil {
						w.logger.Error("attesting reports", zap.Error(err))
					}

				}
			}
		},
	)
}

func (w *AttestationWorker) Stop() {
	w.cancel()
	w.wg.Wait()
}

// AttestReports fetches all reports with pending attestation status
// and attempts to attest each one.
// Returns an error if there was a problem fetching the reports.
func (w *AttestationWorker) AttestReports() error {
	haLock, err := w.store.GetAdvisoryLocker(w.ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = haLock.Release()
	}()

	err = haLock.LockAttestationWorker()
	if err != nil {
		return err
	}

	uncheckedReports, err := w.findReportsNeedingAttestation()
	if err != nil {
		return err
	}

	for _, report := range uncheckedReports {
		if err := w.attestReport(report); err != nil {
			payerreport.AddReportLogFields(w.logger, &report.PayerReport).
				Error("attesting report", zap.Error(err))
		}
	}

	return nil
}

// findReportsNeedingAttestation fetches all reports that are pending attestation and pending submission.
// The only possible state where a report is needing attestation is when it's pending submission and attestation.
func (w *AttestationWorker) findReportsNeedingAttestation() ([]*payerreport.PayerReportWithStatus, error) {
	w.logger.Debug("fetching reports needing attestation")
	return w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
			WithAttestationStatus(payerreport.AttestationPending).
			WithSubmissionStatus(payerreport.SubmissionPending),
	)
}

// attestReport validates a single report by checking its consistency with previous reports
// and, if valid, signs it with the node's private key.
// Returns an error if the report is invalid or if there was a problem during attestation.
func (w *AttestationWorker) attestReport(report *payerreport.PayerReportWithStatus) error {
	log := payerreport.AddReportLogFields(w.logger, &report.PayerReport)

	var (
		prevReport *payerreport.PayerReport
		err        error
	)

	if report.StartSequenceID > 0 {
		var prevReportWithStatus *payerreport.PayerReportWithStatus
		if prevReportWithStatus, err = w.getPreviousReport(report); err != nil {
			return err
		}

		prevReport = &prevReportWithStatus.PayerReport
	}

	verifyResult, err := w.verifier.VerifyReport(w.ctx, prevReport, &report.PayerReport)
	if err != nil {
		return err
	}

	if verifyResult.IsValid {
		log.Info(
			"report is valid, submitting attestation",
			utils.ReasonField(verifyResult.Reason),
		)
		return w.submitAttestation(report)
	}

	log.Warn("report is invalid, not attesting", utils.ReasonField(verifyResult.Reason))
	return w.rejectAttestation(report)
}

// getPreviousReport retrieves the previous report for a given current report.
// The previous report should have been submitted or settled and should end
// at the start sequence ID of the current report.
// Returns an error if no previous report is found or if multiple previous reports are found.
func (w *AttestationWorker) getPreviousReport(
	currentReport *payerreport.PayerReportWithStatus,
) (*payerreport.PayerReportWithStatus, error) {
	w.logger.Debug("fetching previous report",
		utils.OriginatorIDField(currentReport.OriginatorNodeID),
		utils.StartSequenceIDField(int64(currentReport.StartSequenceID)),
	)
	prevReports, err := w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
			WithOriginatorNodeID(currentReport.OriginatorNodeID).
			// Look up reports that have been submitted or settled
			WithSubmissionStatus(payerreport.SubmissionSubmitted, payerreport.SubmissionSettled).
			// The previous report must end at exactly the start of this report
			WithEndSequenceID(currentReport.StartSequenceID),
	)
	if err != nil {
		return nil, err
	}

	if len(prevReports) != 1 {
		return nil, fmt.Errorf("expected 1 previous report, got %d", len(prevReports))
	}

	return prevReports[0], nil
}

// Save the attestation to the database and set the status
func (w *AttestationWorker) submitAttestation(
	report *payerreport.PayerReportWithStatus,
) error {
	nodeSignature, err := w.registrant.SignPayerReportAttestation(report.ID)
	if err != nil {
		return err
	}

	attestation := payerreport.NewPayerReportAttestation(
		&report.PayerReport,
		*nodeSignature,
	)

	clientEnvelope, err := attestation.ToClientEnvelope()
	if err != nil {
		return err
	}

	// Get a signed Payer Envelope, using the node's private key as the payer
	payerEnvelope, err := w.signClientEnvelope(clientEnvelope)
	if err != nil {
		return err
	}

	return w.store.CreateAttestation(w.ctx, attestation, payerEnvelope)
}

func (w *AttestationWorker) signClientEnvelope(
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopes.PayerEnvelope, error) {
	envelopeBytes, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, err
	}

	originatorID := w.registrant.NodeID()

	payerSignature, err := w.registrant.SignClientEnvelopeToSelf(envelopeBytes)
	if err != nil {
		return nil, err
	}

	protoEnvelope := envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: envelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     originatorID,
		MessageRetentionDays: constants.DefaultStorageDurationDays,
	}

	return envelopes.NewPayerEnvelope(&protoEnvelope)
}

func (w *AttestationWorker) rejectAttestation(report *payerreport.PayerReportWithStatus) error {
	return w.store.SetReportAttestationRejected(w.ctx, report.ID)
}

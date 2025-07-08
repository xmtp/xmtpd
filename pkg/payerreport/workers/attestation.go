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
	"go.uber.org/zap"
)

// AttestationWorker is responsible for periodically checking for reports that need attestation
// and signing them with the node's private key.
type AttestationWorker struct {
	ctx             context.Context
	cancel          context.CancelFunc
	log             *zap.Logger
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
	log *zap.Logger,
	registrant registrant.IRegistrant,
	store payerreport.IPayerReportStore,
	pollInterval time.Duration,
	domainSeparator common.Hash,
) *AttestationWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &AttestationWorker{
		ctx:             ctx,
		log:             log.Named("attestationworker"),
		registrant:      registrant,
		store:           store,
		verifier:        payerreport.NewPayerReportVerifier(log, store),
		wg:              sync.WaitGroup{},
		cancel:          cancel,
		pollInterval:    pollInterval,
		domainSeparator: domainSeparator,
	}

	return worker
}

// start launches the worker's main loop in a separate goroutine.
// The loop periodically checks for reports that need attestation.
func (w *AttestationWorker) Start() {
	tracing.GoPanicWrap(
		w.ctx,
		&w.wg,
		"attestation-worker",
		func(ctx context.Context) {
			w.log.Info("Starting attestation worker")
			var err error
			ticker := time.NewTicker(w.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err = w.AttestReports(); err != nil {
						w.log.Error("attesting reports", zap.Error(err))
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
	uncheckedReports, err := w.findReportsNeedingAttestation()
	if err != nil {
		return err
	}

	for _, report := range uncheckedReports {
		if err := w.attestReport(report); err != nil {
			w.log.Error("attesting report", zap.Error(err))
		}
	}

	return nil
}

func (w *AttestationWorker) findReportsNeedingAttestation() ([]*payerreport.PayerReportWithStatus, error) {
	return w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().WithAttestationStatus(payerreport.AttestationPending),
	)
}

// attestReport validates a single report by checking its consistency with previous reports
// and, if valid, signs it with the node's private key.
// Returns an error if the report is invalid or if there was a problem during attestation.
func (w *AttestationWorker) attestReport(report *payerreport.PayerReportWithStatus) error {
	var prevReport *payerreport.PayerReport
	var err error
	if report.StartSequenceID > 0 {
		var prevReportWithStatus *payerreport.PayerReportWithStatus
		if prevReportWithStatus, err = w.getPreviousReport(report); err != nil {
			return err
		}

		prevReport = &prevReportWithStatus.PayerReport
	}

	isValid, err := w.verifier.IsValidReport(w.ctx, prevReport, &report.PayerReport)
	if err != nil {
		return err
	}

	if isValid {
		return w.submitAttestation(report)
	}

	return w.rejectAttestation(report)
}

// getPreviousReport retrieves the previous report for a given current report.
// The previous report should have been submitted or settled and should end
// at the start sequence ID of the current report.
// Returns an error if no previous report is found or if multiple previous reports are found.
func (w *AttestationWorker) getPreviousReport(
	currentReport *payerreport.PayerReportWithStatus,
) (*payerreport.PayerReportWithStatus, error) {
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
		MessageRetentionDays: constants.DEFAULT_STORAGE_DURATION_DAYS,
	}

	return envelopes.NewPayerEnvelope(&protoEnvelope)
}

func (w *AttestationWorker) rejectAttestation(report *payerreport.PayerReportWithStatus) error {
	return w.store.SetReportAttestationStatus(
		w.ctx,
		report.ID,
		[]payerreport.AttestationStatus{payerreport.AttestationPending},
		payerreport.AttestationRejected,
	)
}

package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

// attestationWorker is responsible for periodically checking for reports that need attestation
// and signing them with the node's private key.
type attestationWorker struct {
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   registrant.IRegistrant
	store        payerreport.IPayerReportStore
	verifier     payerreport.PayerReportVerifier
	wg           sync.WaitGroup
	pollInterval time.Duration
}

// NewAttestationWorker creates and starts a new attestation worker that will periodically
// check for reports that need attestation.
// It takes a context, logger, registrant for signing, store for accessing reports,
// and a poll interval that determines how often to check for reports.
func NewAttestationWorker(
	ctx context.Context,
	log *zap.Logger,
	registrant registrant.IRegistrant,
	verifier payerreport.PayerReportVerifier,
	store payerreport.IPayerReportStore,
	pollInterval time.Duration,
) *attestationWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &attestationWorker{
		ctx:          ctx,
		log:          log.Named("attestationworker"),
		registrant:   registrant,
		store:        store,
		verifier:     verifier,
		wg:           sync.WaitGroup{},
		cancel:       cancel,
		pollInterval: pollInterval,
	}

	worker.start()

	return worker
}

// start launches the worker's main loop in a separate goroutine.
// The loop periodically checks for reports that need attestation.
func (w *attestationWorker) start() {
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
					if err = w.attestReports(); err != nil {
						w.log.Error("attesting reports", zap.Error(err))
					}

				}
			}
		},
	)
}

// attestReports fetches all reports with pending attestation status
// and attempts to attest each one.
// Returns an error if there was a problem fetching the reports.
func (w *attestationWorker) attestReports() error {
	uncheckedReports, err := w.findReportsNeedingAttestation()
	if err != nil {
		return err
	}

	for _, report := range uncheckedReports {
		if err := w.attestReport(report); err != nil {
			w.log.Error("attesting report", zap.Error(err))
		}
	}

	return err
}

func (w *attestationWorker) findReportsNeedingAttestation() ([]*payerreport.PayerReportWithStatus, error) {
	return w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().WithAttestationStatus(payerreport.AttestationPending),
	)
}

// attestReport validates a single report by checking its consistency with previous reports
// and, if valid, signs it with the node's private key.
// Returns an error if the report is invalid or if there was a problem during attestation.
func (w *attestationWorker) attestReport(report *payerreport.PayerReportWithStatus) error {
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

	return w.submitAttestation(report, isValid)
}

// getPreviousReport retrieves the previous report for a given current report.
// The previous report should have been submitted or settled and should end
// at the start sequence ID of the current report.
// Returns an error if no previous report is found or if multiple previous reports are found.
func (w *attestationWorker) getPreviousReport(
	currentReport *payerreport.PayerReportWithStatus,
) (*payerreport.PayerReportWithStatus, error) {
	prevReports, err := w.store.FetchReports(
		w.ctx,
		payerreport.NewFetchReportsQuery().
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
func (w *attestationWorker) submitAttestation(
	report *payerreport.PayerReportWithStatus,
	isValid bool,
) error {
	// TODO:nm save the attestation to the DB so that it can be replicated to other nodes
	var newStatus payerreport.AttestationStatus
	if isValid {
		newStatus = payerreport.AttestationApproved
	} else {
		newStatus = payerreport.AttestationRejected
	}

	return w.store.SetReportAttestationStatus(
		w.ctx,
		report.ID,
		[]payerreport.AttestationStatus{payerreport.AttestationPending},
		newStatus,
	)
}

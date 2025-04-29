package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

type attestationWorker struct {
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   *registrant.Registrant
	store        payerreport.Store
	wg           sync.WaitGroup
	pollInterval time.Duration
}

func NewAttestationWorker(
	ctx context.Context,
	log *zap.Logger,
	registrant *registrant.Registrant,
	store payerreport.Store,
	pollInterval time.Duration,
) *attestationWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &attestationWorker{
		ctx:          ctx,
		log:          log.Named("attestationworker"),
		registrant:   registrant,
		store:        store,
		wg:           sync.WaitGroup{},
		cancel:       cancel,
		pollInterval: pollInterval,
	}

	worker.start()

	return worker
}

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

func (w *attestationWorker) attestReports() error {
	pendingStatus := payerreport.AttestationStatus(payerreport.AttestationPending)
	uncheckedReports, err := w.store.FetchReports(w.ctx, payerreport.FetchReportsQuery{
		AttestationStatus: &pendingStatus,
		CreatedAfter:      time.Unix(0, 0),
	})
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

func (w *attestationWorker) attestReport(report *payerreport.PayerReportWithStatus) error {
	return nil
}

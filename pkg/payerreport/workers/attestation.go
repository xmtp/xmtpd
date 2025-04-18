package workers

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

type attestationWorker struct {
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   *registrant.Registrant
	queries      *queries.Queries
	wg           sync.WaitGroup
	pollInterval time.Duration
}

func NewAttestationWorker(
	ctx context.Context,
	log *zap.Logger,
	registrant *registrant.Registrant,
	queries *queries.Queries,
	pollInterval time.Duration,
) *attestationWorker {
	ctx, cancel := context.WithCancel(ctx)

	worker := &attestationWorker{
		ctx:          ctx,
		log:          log.Named("attestationworker"),
		registrant:   registrant,
		queries:      queries,
		wg:           sync.WaitGroup{},
		cancel:       cancel,
		pollInterval: pollInterval,
	}

	tracing.GoPanicWrap(ctx, &worker.wg, "attestation-worker", func(ctx context.Context) {
		worker.start()
	})

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
	return nil
}

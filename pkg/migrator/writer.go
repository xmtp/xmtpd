package migrator

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

// TODO: Run in db tx.
// TODO: Insert gateway envelope.
// TODO: Increment unsettled usage.
// If InboxLog (IdentityUpdates), derive AddressLog and insert both.

func (s *dbMigrator) insertOriginatorEnvelope(
	ctx context.Context,
	env *envelopes.OriginatorEnvelope,
) re.RetryableError {

	return nil
}

func retry(
	ctx context.Context,
	logger *zap.Logger,
	sleep time.Duration,
	fn func() re.RetryableError,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := fn(); err != nil {
				logger.Error("error storing log", zap.Error(err))

				if err.ShouldRetry() {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(sleep):
						continue
					}
				}

				return err
			}

			return nil
		}
	}
}

package db_migrator

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

const (
	groupMessageOriginatorID   = 10
	welcomeMessageOriginatorID = 11
	inboxLogOriginatorID       = 12 // IdentityUpdates in xmtpd.
	installationOriginatorID   = 13 // KeyPackages in xmtpd.
)

func (s *dbMigrator) insertOriginatorEnvelope(
	_ context.Context,
	_ *envelopes.OriginatorEnvelope,
) re.RetryableError {
	// TODO: Insert gateway envelope.
	// If InboxLog (IdentityUpdates), derive AddressLog and insert both?

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

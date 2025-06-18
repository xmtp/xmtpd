package db_migrator

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/envelopes"
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
) error {
	// TODO: Insert gateway envelope.
	// TODO: Process envelopes with retry.
	// TODO: Usage of retryable errors.

	return nil
}

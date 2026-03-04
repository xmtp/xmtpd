package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const DefaultFindOrCreatePayerMaxRetries = 3

// FindOrCreatePayerWithRetry wraps FindOrCreatePayer with retry logic to handle
// the PostgreSQL race condition where INSERT ... ON CONFLICT DO NOTHING + SELECT
// in a CTE returns no rows when concurrent transactions insert the same address.
func FindOrCreatePayerWithRetry(
	ctx context.Context,
	querier *queries.Queries,
	address string,
	maxRetries int,
) (int32, error) {
	id, err := querier.FindOrCreatePayer(ctx, address)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Millisecond):
		}

		id, err = querier.FindOrCreatePayer(ctx, address)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}

	return 0, err
}

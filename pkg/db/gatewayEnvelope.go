package db

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// InsertGatewayEnvelopeAndIncrementUnsettledUsage inserts a gateway envelope and increments the unsettled usage for the payer.
// It returns the number of rows inserted.
func InsertGatewayEnvelopeAndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	insertParams queries.InsertGatewayEnvelopeParams,
	incrementParams queries.IncrementUnsettledUsageParams,
) (int64, error) {
	return RunInTxWithResult(
		ctx,
		db,
		&sql.TxOptions{},
		func(ctx context.Context, txQueries *queries.Queries) (int64, error) {
			numInserted, err := txQueries.InsertGatewayEnvelope(ctx, insertParams)
			if err != nil {
				return 0, err
			}
			// If the numInserted is 0 it means the envelope already exists
			// and we don't need to increment the unsettled usage
			if numInserted == 0 {
				return 0, nil
			}

			err = txQueries.IncrementUnsettledUsage(ctx, incrementParams)
			if err != nil {
				return 0, err
			}

			return numInserted, nil
		},
	)
}

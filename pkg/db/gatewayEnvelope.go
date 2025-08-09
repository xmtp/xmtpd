package db

import (
	"context"
	"database/sql"
	"sync"

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

			var wg sync.WaitGroup
			var incrementErr, congestionErr error
			// Use the sequence ID from the envelope to set the last sequence ID value
			if incrementParams.SequenceID == 0 {
				incrementParams.SequenceID = insertParams.OriginatorSequenceID
			}
			// In this case, the message count is always 1
			if incrementParams.MessageCount == 0 {
				incrementParams.MessageCount = 1
			}

			wg.Add(2)

			go func() {
				defer wg.Done()
				incrementErr = txQueries.IncrementUnsettledUsage(ctx, incrementParams)
			}()

			go func() {
				defer wg.Done()
				congestionErr = txQueries.IncrementOriginatorCongestion(
					ctx,
					queries.IncrementOriginatorCongestionParams{
						OriginatorID:      incrementParams.OriginatorID,
						MinutesSinceEpoch: incrementParams.MinutesSinceEpoch,
					},
				)
			}()

			wg.Wait()

			if incrementErr != nil {
				return 0, incrementErr
			}

			if congestionErr != nil {
				return 0, congestionErr
			}

			return numInserted, nil
		},
	)
}

package db

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// InsertGatewayEnvelopeAndIncrementUnsettledUsage inserts a gateway envelope and
// updates unsettled usage and congestion counters within a single database transaction.
//
// This function runs inside a managed transaction created by RunInTxWithResult().
//
// Steps:
//  1. Calls InsertGatewayEnvelopeWithChecksTransactional() to insert the envelope,
//     automatically creating any missing partitions if needed.
//  2. If a new envelope is inserted, increments unsettled usage and congestion
//     counters for the originator within the same transaction.
//  3. If the envelope already exists (duplicate insert), metrics are not updated.
//
// The function ensures atomicity between envelope insertion and usage updates.
// Safe for high-throughput ingest paths where message persistence and usage tracking
// must succeed or fail together.
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
			numInserted, err := InsertGatewayEnvelopeWithChecksTransactional(
				ctx,
				txQueries,
				insertParams,
			)
			if err != nil {
				return 0, err
			}
			// If the numInserted is 0 it means the envelope already exists
			// and we don't need to increment the unsettled usage
			if numInserted.InsertedMetaRows == 0 {
				return 0, nil
			}

			var wg sync.WaitGroup
			// Use the sequence ID from the envelope to set the last sequence ID value
			if incrementParams.SequenceID == 0 {
				incrementParams.SequenceID = insertParams.OriginatorSequenceID
			}
			// In this case, the message count is always 1
			if incrementParams.MessageCount == 0 {
				incrementParams.MessageCount = 1
			}

			wg.Add(2)

			err = txQueries.IncrementUnsettledUsage(ctx, incrementParams)
			if err != nil {
				return 0, err
			}

			err = txQueries.IncrementOriginatorCongestion(
				ctx,
				queries.IncrementOriginatorCongestionParams{
					OriginatorID:      incrementParams.OriginatorID,
					MinutesSinceEpoch: incrementParams.MinutesSinceEpoch,
				},
			)
			if err != nil {
				return 0, err
			}

			return numInserted.InsertedMetaRows, nil
		},
	)
}

// InsertGatewayEnvelopeWithChecksTransactional attempts to insert a gateway envelope
// inside the current SQL transaction, with automatic partition creation and retry.
//
// Behavior:
//   - Creates a SAVEPOINT before the insert so that a failure does not abort the entire tx.
//   - On “no partition of relation …” errors, it rolls back to the savepoint,
//     creates the missing partition(s) using the same connection (within the tx),
//     and retries the insert once.
//   - On success, the savepoint is released.
//
// This variant must be called within an active transaction. It avoids full tx rollbacks
// and ensures inserts can proceed even when new partitions are required.
// Use for transactional ingestion flows where atomicity must be preserved.
func InsertGatewayEnvelopeWithChecksTransactional(
	ctx context.Context,
	q *queries.Queries,
	row queries.InsertGatewayEnvelopeParams,
) (queries.InsertGatewayEnvelopeRow, error) {
	err := q.InsertSavePoint(ctx)
	if err != nil {
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	inserted, err := q.InsertGatewayEnvelope(ctx, row)

	if err == nil {
		_ = q.InsertSavePointRelease(ctx)
		return inserted, nil
	}

	if !strings.Contains(err.Error(), "no partition of relation") {
		// leave tx in aborted state; caller will handle rollback
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	err = q.InsertSavePointRollback(ctx)
	if err != nil {
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	err = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     row.OriginatorNodeID,
		OriginatorSequenceID: row.OriginatorSequenceID,
		BandWidth:            1_000_000,
	})
	if err != nil {
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	return q.InsertGatewayEnvelope(ctx, row)
}

// InsertGatewayEnvelopeWithChecksStandalone inserts a gateway envelope in a non-transactional context,
// automatically creating missing partitions and retrying once.
//
// Behavior:
//   - Performs an insert into the v2 tables.
//   - On “no partition of relation …” errors, creates the necessary partitions
//     in the same connection, and retries the insert once.
//
// This function does not use SAVEPOINTs and does not depend on an explicit transaction.
// It is ideal for standalone operations such as backfills, batch imports, or
// ingestion workers where each insert is independent of others.
func InsertGatewayEnvelopeWithChecksStandalone(
	ctx context.Context,
	q *queries.Queries,
	row queries.InsertGatewayEnvelopeParams,
) (queries.InsertGatewayEnvelopeRow, error) {
	inserted, err := q.InsertGatewayEnvelope(ctx, row)

	if err == nil {
		return inserted, nil
	}

	if !strings.Contains(err.Error(), "no partition of relation") {
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	err = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     row.OriginatorNodeID,
		OriginatorSequenceID: row.OriginatorSequenceID,
		BandWidth:            1_000_000,
	})
	if err != nil {
		return queries.InsertGatewayEnvelopeRow{}, err
	}

	// retry insert
	return q.InsertGatewayEnvelope(ctx, row)
}

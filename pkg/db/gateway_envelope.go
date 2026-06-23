package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// ErrGatewayPartitionMissing signals that an insert could not proceed because the target
// gateway-envelope partition does not exist yet. Callers that own the surrounding transaction
// should roll it back, call EnsureGatewayPartitions (which creates the partition in its own
// transaction under the exclusive partition-creation lock), and retry the insert.
var ErrGatewayPartitionMissing = errors.New("gateway partition missing")

func isNoPartitionErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no partition of relation")
}

// EnsureGatewayPartitions creates the meta/blob partitions for the given originator and
// sequence id in a dedicated transaction, serialized by the exclusive partition-creation
// advisory lock.
//
// It runs on its own transaction (not the caller's insert transaction): the partition-creation
// DDL takes ShareRowExclusiveLock on the shared meta/blob parents (via the blob->meta and
// meta->payers foreign keys), which conflicts with the RowExclusiveLock of concurrent inserts.
// Holding the exclusive advisory lock here, while inserts hold the shared lock, guarantees DDL
// never overlaps DML and removes the cross-originator deadlock. A separate transaction also
// avoids a shared->exclusive advisory-lock upgrade, which would itself deadlock.
func EnsureGatewayPartitions(
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
	originatorSequenceID int64,
) error {
	return RunInTx(ctx, db, &sql.TxOptions{}, func(ctx context.Context, q *queries.Queries) error {
		if err := NewAdvisoryLocker().LockPartitionCreation(ctx, q); err != nil {
			return err
		}
		return q.EnsureGatewayPartsV3(ctx, queries.EnsureGatewayPartsV3Params{
			OriginatorNodeID:     originatorNodeID,
			OriginatorSequenceID: originatorSequenceID,
			BandWidth:            GatewayEnvelopeBandWidth,
		})
	})
}

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
	insertParams queries.InsertGatewayEnvelopeV3Params,
	incrementParams queries.IncrementUnsettledUsageParams,
	incrementCongestion bool,
) (int64, error) {
	insertTx := func(ctx context.Context, txQueries *queries.Queries) (int64, error) {
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

		// Use the sequence ID from the envelope to set the last sequence ID value
		if incrementParams.SequenceID == 0 {
			incrementParams.SequenceID = insertParams.OriginatorSequenceID
		}
		// In this case, the message count is always 1
		if incrementParams.MessageCount == 0 {
			incrementParams.MessageCount = 1
		}

		err = txQueries.IncrementUnsettledUsage(ctx, incrementParams)
		if err != nil {
			return 0, err
		}

		if !incrementCongestion {
			return numInserted.InsertedMetaRows, nil
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
	}

	result, err := RunInTxWithResult(ctx, db, &sql.TxOptions{}, insertTx)
	if errors.Is(err, ErrGatewayPartitionMissing) {
		// Create the missing partition out-of-band under the exclusive lock, then retry once.
		if ensErr := EnsureGatewayPartitions(
			ctx, db, insertParams.OriginatorNodeID, insertParams.OriginatorSequenceID,
		); ensErr != nil {
			return 0, ensErr
		}
		result, err = RunInTxWithResult(ctx, db, &sql.TxOptions{}, insertTx)
	}
	return result, err
}

// InsertGatewayEnvelopeWithChecksTransactional attempts to insert a gateway envelope inside the
// current SQL transaction.
//
// Behavior:
//   - Takes the shared partition-creation advisory lock so the insert is a reader: it runs
//     concurrently with other inserts but never overlaps exclusive partition creation.
//   - Creates a SAVEPOINT before the insert so that a failure does not abort the entire tx.
//   - On a missing-partition error, it rolls back to the savepoint and returns
//     ErrGatewayPartitionMissing. The caller (which owns the transaction) must roll back, call
//     EnsureGatewayPartitions to create the partition in its own transaction under the exclusive
//     lock, and retry. Partition creation is NOT done inline here: that would take
//     ShareRowExclusiveLock on the shared meta/blob parents while a concurrent insert holds a
//     conflicting lock, deadlocking across originators.
//   - On success, the savepoint is released.
//
// This variant must be called within an active transaction. Use
// InsertGatewayEnvelopeAndIncrementUnsettledUsage for the caller-side ensure-and-retry loop.
func InsertGatewayEnvelopeWithChecksTransactional(
	ctx context.Context,
	q *queries.Queries,
	row queries.InsertGatewayEnvelopeV3Params,
) (queries.InsertGatewayEnvelopeV3Row, error) {
	if err := NewAdvisoryLocker().SharedLockPartitionCreation(ctx, q); err != nil {
		return queries.InsertGatewayEnvelopeV3Row{}, err
	}

	err := q.InsertSavePoint(ctx)
	if err != nil {
		return queries.InsertGatewayEnvelopeV3Row{}, err
	}

	inserted, err := q.InsertGatewayEnvelopeV3(ctx, row)

	if err == nil {
		_ = q.InsertSavePointRelease(ctx)
		return inserted, nil
	}

	if !isNoPartitionErr(err) {
		// leave tx in aborted state; caller will handle rollback
		return queries.InsertGatewayEnvelopeV3Row{}, err
	}

	if rbErr := q.InsertSavePointRollback(ctx); rbErr != nil {
		return queries.InsertGatewayEnvelopeV3Row{}, rbErr
	}

	// Partition creation is done out-of-band by the caller (see EnsureGatewayPartitions); doing
	// it inline here would take ShareRowExclusiveLock on the shared meta/blob parents while a
	// concurrent insert holds a conflicting lock, deadlocking across originators.
	return queries.InsertGatewayEnvelopeV3Row{}, ErrGatewayPartitionMissing
}

// InsertGatewayEnvelopeWithChecksStandalone inserts a gateway envelope in a non-transactional context,
// automatically creating missing partitions and retrying once.
//
// Behavior:
//   - Performs an insert into the v3 tables.
//   - On “no partition of relation …” errors, creates the necessary partitions
//     in the same connection, and retries the insert once.
//
// This function does not use SAVEPOINTs and does not depend on an explicit transaction.
// It is ideal for standalone operations such as backfills, batch imports, or
// ingestion workers where each insert is independent of others.
func InsertGatewayEnvelopeWithChecksStandalone(
	ctx context.Context,
	q *queries.Queries,
	row queries.InsertGatewayEnvelopeV3Params,
) (queries.InsertGatewayEnvelopeV3Row, error) {
	inserted, err := q.InsertGatewayEnvelopeV3(ctx, row)

	if err == nil {
		return inserted, nil
	}

	if !strings.Contains(err.Error(), "no partition of relation") {
		return queries.InsertGatewayEnvelopeV3Row{}, err
	}

	err = q.EnsureGatewayPartsV3(ctx, queries.EnsureGatewayPartsV3Params{
		OriginatorNodeID:     row.OriginatorNodeID,
		OriginatorSequenceID: row.OriginatorSequenceID,
		BandWidth:            GatewayEnvelopeBandWidth,
	})
	if err != nil {
		return queries.InsertGatewayEnvelopeV3Row{}, err
	}

	// retry insert
	return q.InsertGatewayEnvelopeV3(ctx, row)
}

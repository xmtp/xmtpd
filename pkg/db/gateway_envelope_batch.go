package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"go.uber.org/zap"
)

// InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage inserts a batch of gateway envelopes,
// updates unsettled usage, and tracks originator congestion within a single database transaction.
//
// This is a convenience wrapper that creates its own transaction and, on a missing partition,
// creates it out-of-band under the exclusive partition-creation lock before retrying. Use
// InsertGatewayEnvelopeBatchV2Transactional when you need to participate in an existing
// transaction (the caller is then responsible for the ensure-and-retry loop).
func InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	logger *zap.Logger,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	insertTx := func(ctx context.Context, q *queries.Queries) (int64, error) {
		return InsertGatewayEnvelopeBatchV2Transactional(ctx, q, logger, input)
	}

	result, err := RunInTxWithResult(ctx, db, &sql.TxOptions{}, insertTx)
	if errors.Is(err, ErrGatewayPartitionMissing) {
		logger.Info("creating partitions for batch insert")
		if ensErr := EnsureGatewayPartitionsForBatch(ctx, db, input); ensErr != nil {
			return 0, ensErr
		}
		result, err = RunInTxWithResult(ctx, db, &sql.TxOptions{}, insertTx)
	}
	return result, err
}

// InsertGatewayEnvelopeBatchV2Transactional inserts a batch of gateway envelopes within an
// existing transaction, using the V2 SQL function that also tracks originator congestion.
//
// It takes the shared partition-creation advisory lock (running as a reader, concurrent with
// other inserts but never overlapping exclusive partition creation). On a missing partition it
// rolls back to the savepoint and returns ErrGatewayPartitionMissing; the caller that owns the
// transaction must roll back, call EnsureGatewayPartitionsForBatch, and retry. Partition
// creation is not done inline (see InsertGatewayEnvelopeWithChecksTransactional for why).
func InsertGatewayEnvelopeBatchV2Transactional(
	ctx context.Context,
	q *queries.Queries,
	logger *zap.Logger,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	if input.Len() == 0 {
		return 0, errors.New("empty input")
	}

	params := input.ToParamsV3()

	if err := NewAdvisoryLocker().SharedLockPartitionCreation(ctx, q); err != nil {
		return 0, err
	}

	err := q.InsertSavePoint(ctx)
	if err != nil {
		return 0, err
	}

	result, err := q.InsertGatewayEnvelopeBatchV3(ctx, params)
	if err == nil {
		_ = q.InsertSavePointRelease(ctx)
		return result.InsertedMetaRows, nil
	}

	if !isNoPartitionErr(err) {
		return 0, fmt.Errorf("insert batch v2: %w", err)
	}

	if rbErr := q.InsertSavePointRollback(ctx); rbErr != nil {
		return 0, rbErr
	}

	return 0, ErrGatewayPartitionMissing
}

// EnsureGatewayPartitionsForBatch creates the meta/blob partitions for every envelope in the
// batch, each in its own transaction under the exclusive partition-creation lock.
func EnsureGatewayPartitionsForBatch(
	ctx context.Context,
	db *sql.DB,
	input *types.GatewayEnvelopeBatch,
) error {
	for _, envelope := range input.Envelopes {
		if err := EnsureGatewayPartitions(
			ctx, db, envelope.OriginatorNodeID, envelope.OriginatorSequenceID,
		); err != nil {
			return fmt.Errorf("ensure gateway parts: %w", err)
		}
	}
	return nil
}

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"go.uber.org/zap"
)

// InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage inserts a batch of gateway envelopes,
// updates unsettled usage, and tracks originator congestion within a single database transaction.
//
// This is a convenience wrapper that creates its own transaction. Use
// InsertGatewayEnvelopeBatchV2Transactional when you need to participate in an existing transaction.
func InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	logger *zap.Logger,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	return RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, q *queries.Queries) (int64, error) {
			return InsertGatewayEnvelopeBatchV2Transactional(ctx, q, logger, input)
		})
}

// InsertGatewayEnvelopeBatchV2Transactional inserts a batch of gateway envelopes within an
// existing transaction, using the V2 SQL function that also tracks originator congestion.
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

	err := q.InsertSavePoint(ctx)
	if err != nil {
		return 0, err
	}

	result, err := q.InsertGatewayEnvelopeBatchV3(ctx, params)
	if err == nil {
		_ = q.InsertSavePointRelease(ctx)
		return result.InsertedMetaRows, nil
	}

	if !strings.Contains(err.Error(), "no partition of relation") {
		return 0, fmt.Errorf("insert batch v2: %w", err)
	}

	err = q.InsertSavePointRollback(ctx)
	if err != nil {
		return 0, err
	}

	logger.Info("creating partitions for batch insert")

	for _, envelope := range input.Envelopes {
		err = q.EnsureGatewayPartsV4(ctx, queries.EnsureGatewayPartsV4Params{
			OriginatorNodeID:     envelope.OriginatorNodeID,
			OriginatorSequenceID: envelope.OriginatorSequenceID,
			BandWidth:            GatewayEnvelopeBandWidth,
		})
		if err != nil {
			return 0, fmt.Errorf("ensure gateway parts: %w", err)
		}
	}

	result, err = q.InsertGatewayEnvelopeBatchV3(ctx, params)
	if err != nil {
		return 0, fmt.Errorf(
			"insert gateway envelope batch v2 and increment unsettled usage: %w",
			err,
		)
	}

	return result.InsertedMetaRows, nil
}

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
)

// InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage inserts a batch of gateway envelopes and
// updates unsettled usage within a single database transaction.
//
// This is a convenience wrapper that creates its own transaction. Use
// InsertGatewayEnvelopeBatchTransactional when you need to participate in an existing transaction.
func InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	return RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, q *queries.Queries) (int64, error) {
			return InsertGatewayEnvelopeBatchTransactional(ctx, q, input)
		})
}

// InsertGatewayEnvelopeBatchTransactional inserts a batch of gateway envelopes within an existing transaction.
//
// The input is an array of originator node IDs, sequence IDs, topics, payer IDs, gateway times,
// expiries, originator envelopes, and spend picodollars.
//
// The sequenceIDs are expected to be strictly ascending per originator node ID.
//
// Payer IDs considerations:
//   - if not 0, they must exist.
//   - if 0, they are treated as null, as it's nullable in gateway_envelopes_meta.
//   - if 0, no unsettled usage is incremented.
func InsertGatewayEnvelopeBatchTransactional(
	ctx context.Context,
	q *queries.Queries,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	if input.Len() == 0 {
		return 0, errors.New("empty input")
	}

	params := input.ToParams()

	// Create a save point to rollback to if the insert fails.
	err := q.InsertSavePoint(ctx)
	if err != nil {
		return 0, err
	}

	// Optimistically insert the envelopes and increment the unsettled usage.
	result, err := q.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(ctx, params)
	if err == nil {
		_ = q.InsertSavePointRelease(ctx)
		return result.InsertedMetaRows, nil
	}

	// Only retry for missing partition errors.
	if !strings.Contains(err.Error(), "no partition of relation") {
		return 0, fmt.Errorf("insert batch: %w", err)
	}

	// On error, rollback the save point and ensure the gateway parts.
	err = q.InsertSavePointRollback(ctx)
	if err != nil {
		return 0, err
	}

	// Ensure the gateway parts for the originator nodes.
	for _, envelope := range input.Envelopes {
		err = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
			OriginatorNodeID:     envelope.OriginatorNodeID,
			OriginatorSequenceID: envelope.OriginatorSequenceID,
			BandWidth:            GatewayEnvelopeBandWidth,
		})
		if err != nil {
			return 0, fmt.Errorf("ensure gateway parts: %w", err)
		}
	}

	// Retry the insert.
	result, err = q.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(ctx, params)
	if err != nil {
		return 0, fmt.Errorf(
			"insert gateway envelope batch and increment unsettled usage: %w",
			err,
		)
	}

	return result.InsertedMetaRows, nil
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage inserts a batch of gateway envelopes and
// updates unsettled usage and congestion counters within a single database transaction.
func InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	input queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams,
) (int64, error) {
	inputLength := len(input.OriginatorNodeIds)

	if inputLength == 0 {
		return 0, nil
	}

	if len(input.OriginatorSequenceIds) != inputLength ||
		len(input.Topics) != inputLength ||
		len(input.PayerIds) != inputLength ||
		len(input.GatewayTimes) != inputLength ||
		len(input.Expiries) != inputLength ||
		len(input.OriginatorEnvelopes) != inputLength ||
		len(input.SpendPicodollars) != inputLength {
		return 0, fmt.Errorf(
			"input array length mismatch: all arrays must have length %d",
			inputLength,
		)
	}

	return RunInTxWithResult(
		ctx,
		db,
		&sql.TxOptions{},
		func(ctx context.Context, txQueries *queries.Queries) (int64, error) {
			// Create a save point to rollback to if the insert fails.
			err := txQueries.InsertSavePoint(ctx)
			if err != nil {
				return 0, err
			}

			// Insert the envelopes and increment the unsettled usage.
			result, err := txQueries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				input,
			)
			if err == nil {
				_ = txQueries.InsertSavePointRelease(ctx)
				return result.InsertedMetaRows, nil
			}

			// Only retry for missing partition errors.
			if !strings.Contains(err.Error(), "no partition of relation") {
				return 0, fmt.Errorf("insert batch: %w", err)
			}

			// On error, rollback the save point and ensure the gateway parts.
			err = txQueries.InsertSavePointRollback(ctx)
			if err != nil {
				return 0, err
			}

			// Deduplicate originator node IDs.
			seen := make(map[int32]int64)
			for i, nodeID := range input.OriginatorNodeIds {
				if input.OriginatorSequenceIds[i] > seen[nodeID] {
					seen[nodeID] = input.OriginatorSequenceIds[i]
				}
			}

			// Ensure the gateway parts for the originator nodes.
			for nodeID, seqID := range seen {
				err = txQueries.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
					OriginatorNodeID:     nodeID,
					OriginatorSequenceID: seqID,
					BandWidth:            1_000_000,
				})
				if err != nil {
					return 0, fmt.Errorf("ensure gateway parts: %w", err)
				}
			}

			// Retry the insert.
			result, err = txQueries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				input,
			)
			if err != nil {
				return 0, fmt.Errorf(
					"insert gateway envelope batch and increment unsettled usage: %w",
					err,
				)
			}

			return result.InsertedMetaRows, nil
		},
	)
}

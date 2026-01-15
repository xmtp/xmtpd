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
//
// The input is an array of originator node IDs, sequence IDs, topics, payer IDs, gateway times,
// expiries, originator envelopes, and spend picodollars.
//
// The sequenceIDs are expected to be strictly ascending per originator node ID.
// Payers:
//   - if not 0, they must exist.
//   - if 0, they are treated as null, as it's nullable in gateway_envelopes_meta.
//   - if 0, no unsettled usage is incremented.
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

	// Deduplicate originator node IDs.
	// Check that sequence IDs are sorted in ascending order.
	// Save last sequence ID for each originator node.
	seen := make(map[int32]int64)
	for i, nodeID := range input.OriginatorNodeIds {
		seqID := input.OriginatorSequenceIds[i]
		if lastSeq, exists := seen[nodeID]; exists && seqID <= lastSeq {
			return 0, fmt.Errorf(
				"originator %d: sequence IDs must be strictly ascending (got %d after %d)",
				nodeID, seqID, lastSeq,
			)
		}
		seen[nodeID] = seqID
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

			// Optimistically insert the envelopes and increment the unsettled usage.
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

			// Ensure the gateway parts for the originator nodes.
			for i, nodeID := range input.OriginatorNodeIds {
				err = txQueries.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
					OriginatorNodeID:     nodeID,
					OriginatorSequenceID: input.OriginatorSequenceIds[i],
					BandWidth:            GatewayEnvelopeBandWidth,
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

package db

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
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
	input []types.GatewayEnvelopeRow,
) (int64, error) {
	if len(input) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	// Order by originator sequence ID ascending for each originator node.
	slices.SortFunc(input, func(a, b types.GatewayEnvelopeRow) int {
		if a.OriginatorSequenceID < b.OriginatorSequenceID {
			return -1
		}
		if a.OriginatorSequenceID > b.OriginatorSequenceID {
			return 1
		}
		return 0
	})

	params := toParallelArrays(input)

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
				params,
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
			for _, envelope := range input {
				err = txQueries.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
					OriginatorNodeID:     envelope.OriginatorNodeID,
					OriginatorSequenceID: envelope.OriginatorSequenceID,
					BandWidth:            GatewayEnvelopeBandWidth,
				})
				if err != nil {
					return 0, fmt.Errorf("ensure gateway parts: %w", err)
				}
			}

			// Retry the insert.
			result, err = txQueries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				params,
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

// toParallelArrays converts a slice of GatewayEnvelopeRow to parallel arrays.
func toParallelArrays(
	input []types.GatewayEnvelopeRow,
) queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams {
	n := len(input)

	params := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     make([]int32, n),
		OriginatorSequenceIds: make([]int64, n),
		Topics:                make([][]byte, n),
		PayerIds:              make([]int32, n),
		GatewayTimes:          make([]time.Time, n),
		Expiries:              make([]int64, n),
		OriginatorEnvelopes:   make([][]byte, n),
		SpendPicodollars:      make([]int64, n),
	}

	for i, row := range input {
		params.OriginatorNodeIds[i] = row.OriginatorNodeID
		params.OriginatorSequenceIds[i] = row.OriginatorSequenceID
		params.Topics[i] = row.Topic
		params.PayerIds[i] = row.PayerID
		params.GatewayTimes[i] = row.GatewayTime
		params.Expiries[i] = row.Expiry
		params.OriginatorEnvelopes[i] = row.OriginatorEnvelope
		params.SpendPicodollars[i] = row.SpendPicodollars
	}

	return params
}

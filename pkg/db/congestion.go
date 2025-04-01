package db

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// Gets the congestion for the minute specified by `endMinute` and the previous 4 minutes
// returned in descending order with missing values filled with 0
func Get5MinutesOfCongestion(
	ctx context.Context,
	querier *queries.Queries,
	originatorID, endMinute int32,
) (out [5]int32, err error) {
	var congestion []queries.GetRecentOriginatorCongestionRow
	congestion, err = querier.GetRecentOriginatorCongestion(
		ctx,
		queries.GetRecentOriginatorCongestionParams{
			OriginatorID: originatorID,
			EndMinute:    endMinute,
			NumMinutes:   5,
		},
	)
	if err != nil {
		return out, err
	}

	for _, congestion := range congestion {
		idx := endMinute - congestion.MinutesSinceEpoch
		out[idx] = congestion.NumMessages
	}

	return out, nil
}

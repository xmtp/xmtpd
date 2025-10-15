package payerreport

import (
	"context"
	"math"
	"math/big"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

func getMinuteFromSequenceID(
	ctx context.Context,
	querier *queries.Queries,
	originatorID int32,
	sequenceID int64,
) (int32, error) {
	envelope, err := querier.GetGatewayEnvelopeByID(ctx, queries.GetGatewayEnvelopeByIDParams{
		OriginatorSequenceID: sequenceID,
		OriginatorNodeID:     originatorID,
	})
	if err != nil {
		return 0, err
	}

	parsedEnvelope, err := envelopes.NewOriginatorEnvelopeFromBytes(envelope.OriginatorEnvelope)
	if err != nil {
		return 0, err
	}

	return getMinuteFromEnvelope(parsedEnvelope), nil
}

func getMinuteFromEnvelope(envelope *envelopes.OriginatorEnvelope) int32 {
	return utils.MinutesSinceEpoch(envelope.OriginatorTime())
}

func AddReportLogFields(logger *zap.Logger, report *PayerReport) *zap.Logger {
	return logger.With(
		zap.String("report_id", report.ID.String()),
		zap.Uint64("start_sequence_id", report.StartSequenceID),
		zap.Uint64("end_sequence_id", report.EndSequenceID),
		zap.Uint32("originator_node_id", report.OriginatorNodeID),
	)
}

// ValidateReportIndex checks if a big.Int PayerReportIndex fits within int32 bounds.
// Returns the validated int32 value or ErrReportIndexTooLarge if the value is too large.
func ValidateReportIndex(index *big.Int) (int32, error) {
	if !index.IsInt64() || index.Int64() > math.MaxInt32 {
		return 0, ErrReportIndexTooLarge
	}
	return int32(index.Int64()), nil
}

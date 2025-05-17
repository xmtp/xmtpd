package payerreport

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
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

	return utils.MinutesSinceEpoch(parsedEnvelope.OriginatorTime()), nil
}

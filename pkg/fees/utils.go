package fees

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

func CalculateStagedOriginatorEnvelopeFees(
	ctx context.Context,
	stagedEnv *queries.StagedOriginatorEnvelope,
	feeCalculator IFeeCalculator,
	querier *queries.Queries,
	nodeID uint32,
	retentionDays uint32,
) (currency.PicoDollar, currency.PicoDollar, error) {
	baseFee, err := feeCalculator.CalculateBaseFee(
		stagedEnv.OriginatorTime,
		int64(len(stagedEnv.PayerEnvelope)),
		retentionDays,
	)
	if err != nil {
		return 0, 0, err
	}

	congestionFee, err := feeCalculator.CalculateCongestionFee(
		ctx,
		querier,
		stagedEnv.OriginatorTime,
		nodeID,
	)
	if err != nil {
		return 0, 0, err
	}

	return baseFee, congestionFee, nil
}

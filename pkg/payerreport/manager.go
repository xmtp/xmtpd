package payerreport

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type PayerReportManager struct {
	log     *zap.Logger
	queries *queries.Queries
}

func NewPayerReportManager(
	log *zap.Logger,
	queries *queries.Queries,
) *PayerReportManager {
	return &PayerReportManager{
		log:     log,
		queries: queries,
	}
}

func (p *PayerReportManager) GenerateReport(
	ctx context.Context,
	params PayerReportGenerationParams,
) (*PayerReport, error) {
	originatorID := int32(params.OriginatorID)
	startMinute, err := p.getStartMinute(
		ctx,
		int64(params.LastReportEndSequenceID),
		originatorID,
	)
	if err != nil {
		return nil, err
	}

	endMinute, endSequenceID, err := p.getEndMinute(ctx, originatorID, startMinute)
	if err != nil {
		return nil, err
	}

	// If the end sequence ID is 0, we don't have enough envelopes to generate a report.
	// Returns an empty report rather than an error here
	if endSequenceID == 0 {
		return &PayerReport{
			OriginatorNodeID: uint32(originatorID),
			Payers:           make(map[common.Address]currency.PicoDollar),
			StartSequenceID:  params.LastReportEndSequenceID,
			EndSequenceID:    params.LastReportEndSequenceID,
			// TODO: Implement merkle calculation
			PayersMerkleRoot: []byte("fix me"),
			PayersLeafCount:  uint32(0),
		}, nil
	}

	payers, err := p.queries.BuildPayerReport(
		ctx,
		queries.BuildPayerReportParams{
			OriginatorID:           originatorID,
			StartMinutesSinceEpoch: startMinute,
			EndMinutesSinceEpoch:   endMinute,
		},
	)
	if err != nil {
		return nil, err
	}

	return &PayerReport{
		OriginatorNodeID: uint32(originatorID),
		Payers:           buildPayersMap(payers),
		StartSequenceID:  params.LastReportEndSequenceID,
		EndSequenceID:    uint64(endSequenceID),
		// TODO: Implement merkle calculation
		PayersMerkleRoot: []byte("fix me"),
		PayersLeafCount:  uint32(len(payers)),
	}, nil
}

/*
*  Returns the start minute to use for the report.
*
*  It does this by getting the envelope from the database with the given sequence ID.
*
*  It then parses the envelope and returns the minute.
 */
func (p *PayerReportManager) getStartMinute(
	ctx context.Context,
	sequenceID int64,
	originatorID int32,
) (int32, error) {
	// If the sequence ID is 0, we're starting from the first envelope
	if sequenceID == 0 {
		return 0, nil
	}

	envelope, err := p.queries.GetGatewayEnvelopeByID(ctx, queries.GetGatewayEnvelopeByIDParams{
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

/*
* Returns the end minute to use for the report.
* It is looking for the second last minute with an envelope for the originator
 */
func (p *PayerReportManager) getEndMinute(
	ctx context.Context,
	originatorID int32,
	startMinute int32,
) (int32, int64, error) {
	result, err := p.queries.GetSecondNewestMinute(
		ctx,
		queries.GetSecondNewestMinuteParams{
			OriginatorID:             originatorID,
			MinimumMinutesSinceEpoch: startMinute,
		},
	)

	if err != nil {
		return 0, 0, err
	}

	return result.MinutesSinceEpoch, result.MaxSequenceID, nil
}

func buildPayersMap(rows []queries.BuildPayerReportRow) map[common.Address]currency.PicoDollar {
	payersMap := make(map[common.Address]currency.PicoDollar)
	for _, row := range rows {
		payersMap[common.HexToAddress(row.PayerAddress)] = currency.PicoDollar(
			row.TotalSpendPicodollars,
		)
	}
	return payersMap
}

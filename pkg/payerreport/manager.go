package payerreport

import (
	"context"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

type payerMap map[common.Address]currency.PicoDollar

type PayerReportManager struct {
	log     *zap.Logger
	queries *queries.Queries
}

func NewPayerReportManager(
	log *zap.Logger,
	queries *queries.Queries,
) *PayerReportManager {
	return &PayerReportManager{
		log:     log.Named("reportmanager"),
		queries: queries,
	}
}

func (p *PayerReportManager) GenerateReport(
	ctx context.Context,
	params PayerReportGenerationParams,
) (*PayerReportWithInputs, error) {
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
		payers := make(map[common.Address]currency.PicoDollar)
		return &PayerReportWithInputs{
			PayerReport: PayerReport{
				OriginatorNodeID: uint32(originatorID),
				StartSequenceID:  params.LastReportEndSequenceID,
				EndSequenceID:    params.LastReportEndSequenceID,
				// TODO: Implement merkle calculation
				PayersMerkleRoot: buildMerkleRoot(payers),
				PayersLeafCount:  uint32(0),
			},
			Payers: payers,
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
	mappedPayers := buildPayersMap(payers)
	return &PayerReportWithInputs{
		PayerReport: PayerReport{
			OriginatorNodeID: uint32(originatorID),
			StartSequenceID:  params.LastReportEndSequenceID,
			EndSequenceID:    uint64(endSequenceID),
			// TODO: Implement merkle calculation
			PayersMerkleRoot: buildMerkleRoot(mappedPayers),
			PayersLeafCount:  uint32(len(payers)),
		},
		Payers: mappedPayers,
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

	return getMinuteFromSequenceID(ctx, p.queries, originatorID, sequenceID)
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

func buildPayersMap(rows []queries.BuildPayerReportRow) payerMap {
	payersMap := make(map[common.Address]currency.PicoDollar)
	for _, row := range rows {
		payersMap[common.HexToAddress(row.PayerAddress)] = currency.PicoDollar(
			row.TotalSpendPicodollars,
		)
	}
	return payersMap
}

// Totally fake function to get a merkle root from a payer map
func buildMerkleRoot(payers payerMap) [32]byte {
	keys := []common.Address{}
	for payerAddress := range payers {
		keys = append(keys, payerAddress)
	}
	sort.SliceStable(keys, func(i int, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	var out [32]byte
	d := ethcrypto.NewKeccakState()
	for _, key := range keys {
		d.Write(key[:])
	}
	//nolint:errcheck
	d.Read(out[:])

	return out
}

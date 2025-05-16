package payerreport

import (
	"context"
	"encoding/binary"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type payerMap map[common.Address]currency.PicoDollar

type PayerReportGenerator struct {
	log             *zap.Logger
	queries         *queries.Queries
	registry        registry.NodeRegistry
	domainSeparator common.Hash
}

func NewPayerReportGenerator(
	log *zap.Logger,
	queries *queries.Queries,
	registry registry.NodeRegistry,
	domainSeparator common.Hash,
) *PayerReportGenerator {
	return &PayerReportGenerator{
		log:             log.Named("reportgenerator"),
		queries:         queries,
		registry:        registry,
		domainSeparator: domainSeparator,
	}
}

func (p *PayerReportGenerator) GenerateReport(
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

	nodes, err := p.registry.GetNodes()
	if err != nil {
		return nil, err
	}

	activeNodeIDs := extractActiveNodeIDs(nodes)

	endMinute, endSequenceID, err := p.getEndMinute(ctx, originatorID, startMinute)
	if err != nil {
		return nil, err
	}

	// If the end sequence ID is 0, we don't have enough envelopes to generate a report.
	// Returns an empty report rather than an error here
	if endSequenceID == 0 {
		payers := make(map[common.Address]currency.PicoDollar)
		return BuildPayerReport(BuildPayerReportParams{
			OriginatorNodeID:    uint32(originatorID),
			StartSequenceID:     params.LastReportEndSequenceID,
			EndSequenceID:       params.LastReportEndSequenceID,
			EndMinuteSinceEpoch: 0,
			Payers:              payers,
			NodeIDs:             activeNodeIDs,
			DomainSeparator:     p.domainSeparator,
		})
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

	return BuildPayerReport(BuildPayerReportParams{
		OriginatorNodeID:    uint32(originatorID),
		StartSequenceID:     params.LastReportEndSequenceID,
		EndSequenceID:       uint64(endSequenceID),
		EndMinuteSinceEpoch: uint32(endMinute),
		NodeIDs:             activeNodeIDs,
		Payers:              mappedPayers,
		DomainSeparator:     p.domainSeparator,
	})
}

/*
*  Returns the start minute to use for the report.
*
*  It does this by getting the envelope from the database with the given sequence ID.
*
*  It then parses the envelope and returns the minute.
 */
func (p *PayerReportGenerator) getStartMinute(
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
func (p *PayerReportGenerator) getEndMinute(
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
func buildMerkleRoot(payers payerMap) common.Hash {
	keys := []common.Address{}
	for payerAddress := range payers {
		keys = append(keys, payerAddress)
	}
	sort.SliceStable(keys, func(i int, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	var out common.Hash
	d := ethcrypto.NewKeccakState()
	for _, key := range keys {
		d.Write(key[:])
		// Convert spend (uint64) to bytes and write it to the hash
		spendBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(spendBytes, uint64(payers[key]))
		d.Write(spendBytes)
	}
	//nolint:errcheck
	d.Read(out[:])

	return out
}

func extractActiveNodeIDs(nodes []registry.Node) []uint32 {
	activeNodeIDs := make([]uint32, len(nodes))
	for i, node := range nodes {
		activeNodeIDs[i] = node.NodeID
	}
	return activeNodeIDs
}

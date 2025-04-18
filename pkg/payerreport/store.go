package payerreport

import (
	"context"
	"errors"
	"math"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var (
	ErrOriginatorNodeIDTooLarge = errors.New("originator node ID is > max int32")
	ErrNodesCountTooLarge       = errors.New("nodes count is > max int32")
)

type Store struct {
	queries *queries.Queries
	log     *zap.Logger
}

func NewStore(queries *queries.Queries, log *zap.Logger) *Store {
	return &Store{
		queries: queries,
		log:     log.Named("payerreportstore"),
	}
}

// Store a report in the database. No validations have been performed.
func (s *Store) StoreReport(ctx context.Context, report *PayerReport) (ReportID, error) {
	id, err := report.ID()
	if err != nil {
		return nil, err
	}

	// The originator node ID is stored as an int32 in the database, but
	// a uint32 on the network. Do not allow anything larger than max int32
	if report.OriginatorNodeID > math.MaxInt32 {
		return nil, ErrOriginatorNodeIDTooLarge
	}

	if report.NodesCount > math.MaxInt32 {
		return nil, ErrNodesCountTooLarge
	}

	err = s.queries.InsertOrIgnorePayerReport(ctx, queries.InsertOrIgnorePayerReportParams{
		ID:               id,
		OriginatorNodeID: int32(report.OriginatorNodeID),
		StartSequenceID:  int64(report.StartSequenceID),
		EndSequenceID:    int64(report.EndSequenceID),
		PayersMerkleRoot: report.PayersMerkleRoot[:],
		PayersLeafCount:  int64(report.PayersLeafCount),
		NodesHash:        report.NodesHash[:],
		NodesCount:       int32(report.NodesCount),
	})
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *Store) StoreAttestation(ctx context.Context, attestation *PayerReportAttestation) error {
	reportID, err := attestation.Report.ID()
	if err != nil {
		return err
	}
	return s.queries.InsertOrIgnorePayerReportAttestation(
		ctx,
		queries.InsertOrIgnorePayerReportAttestationParams{
			PayerReportID: reportID,
			NodeID:        int64(attestation.NodeSignature.NodeID),
			Signature:     attestation.NodeSignature.Signature,
		},
	)
}

func (s *Store) FetchReport(ctx context.Context, id ReportID) (*PayerReport, error) {
	report, err := s.queries.FetchPayerReport(ctx, id)
	if err != nil {
		return nil, err
	}

	var payersMerkleRoot [32]byte
	var nodesHash [32]byte
	if payersMerkleRoot, err = utils.SliceToArray32(report.PayersMerkleRoot); err != nil {
		return nil, err
	}
	if nodesHash, err = utils.SliceToArray32(report.NodesHash); err != nil {
		return nil, err
	}

	return &PayerReport{
		OriginatorNodeID: uint32(report.OriginatorNodeID),
		StartSequenceID:  uint64(report.StartSequenceID),
		EndSequenceID:    uint64(report.EndSequenceID),
		PayersMerkleRoot: payersMerkleRoot,
		PayersLeafCount:  uint32(report.PayersLeafCount),
		NodesHash:        nodesHash,
		NodesCount:       uint32(report.NodesCount),
	}, nil
}

package payerreport

import (
	"context"
	"errors"
	"math"
	"time"

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

func (s *Store) FetchReport(ctx context.Context, id ReportID) (*PayerReportWithStatus, error) {
	report, err := s.queries.FetchPayerReport(ctx, id)
	if err != nil {
		return nil, err
	}

	return convertPayerReport(report)
}

type FetchReportsQuery struct {
	SubmissionStatus  *SubmissionStatus
	AttestationStatus *AttestationStatus
	StartSequenceID   *uint64
	EndSequenceID     *uint64
	CreatedAfter      time.Time
}

func (f *FetchReportsQuery) toParams() queries.FetchPayerReportsParams {
	return queries.FetchPayerReportsParams{
		CreatedAfter:      utils.NewNullTime(f.CreatedAfter),
		SubmissionStatus:  utils.NewNullInt16(f.SubmissionStatus),
		AttestationStatus: utils.NewNullInt16(f.AttestationStatus),
		StartSequenceID:   utils.NewNullInt64(f.StartSequenceID),
		EndSequenceID:     utils.NewNullInt64(f.EndSequenceID),
	}
}

func NewFetchReportsQuery() *FetchReportsQuery {
	return &FetchReportsQuery{}
}

func (f *FetchReportsQuery) WithSubmissionStatus(status SubmissionStatus) *FetchReportsQuery {
	f.SubmissionStatus = &status
	return f
}

func (f *FetchReportsQuery) WithAttestationStatus(status AttestationStatus) *FetchReportsQuery {
	f.AttestationStatus = &status
	return f
}

func (f *FetchReportsQuery) WithCreatedAfter(createdAfter time.Time) *FetchReportsQuery {
	f.CreatedAfter = createdAfter
	return f
}

func (f *FetchReportsQuery) WithStartSequenceID(startSequenceID uint64) *FetchReportsQuery {
	f.StartSequenceID = &startSequenceID
	return f
}

func (f *FetchReportsQuery) WithEndSequenceID(endSequenceID uint64) *FetchReportsQuery {
	f.EndSequenceID = &endSequenceID
	return f
}

func (s *Store) FetchReports(
	ctx context.Context,
	query *FetchReportsQuery,
) ([]*PayerReportWithStatus, error) {
	rows, err := s.queries.FetchPayerReports(ctx, query.toParams())
	if err != nil {
		return nil, err
	}

	return convertPayerReports(rows)
}

func convertPayerReports(rows []queries.PayerReport) ([]*PayerReportWithStatus, error) {
	out := make([]*PayerReportWithStatus, len(rows))
	for idx, row := range rows {
		converted, err := convertPayerReport(row)
		if err != nil {
			return nil, err
		}
		out[idx] = converted
	}
	return out, nil
}

func convertPayerReport(report queries.PayerReport) (*PayerReportWithStatus, error) {
	var (
		err              error
		payersMerkleRoot [32]byte
		nodesHash        [32]byte
		id               [32]byte
	)

	if payersMerkleRoot, err = utils.SliceToArray32(report.PayersMerkleRoot); err != nil {
		return nil, err
	}
	if nodesHash, err = utils.SliceToArray32(report.NodesHash); err != nil {
		return nil, err
	}
	if id, err = utils.SliceToArray32(report.ID); err != nil {
		return nil, err
	}

	return &PayerReportWithStatus{
		SubmissionStatus:  SubmissionStatus(report.SubmissionStatus),
		AttestationStatus: AttestationStatus(report.AttestationStatus),
		CreatedAt:         report.CreatedAt.Time,
		ID:                id,
		PayerReport: PayerReport{
			OriginatorNodeID: uint32(report.OriginatorNodeID),
			StartSequenceID:  uint64(report.StartSequenceID),
			EndSequenceID:    uint64(report.EndSequenceID),
			PayersMerkleRoot: payersMerkleRoot,
			PayersLeafCount:  uint32(report.PayersLeafCount),
			NodesHash:        nodesHash,
			NodesCount:       uint32(report.NodesCount),
		},
	}, nil
}

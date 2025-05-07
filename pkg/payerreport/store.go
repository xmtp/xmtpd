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
	ErrStartSequenceIDTooLarge  = errors.New("start sequence ID is > max int64")
	ErrEndSequenceIDTooLarge    = errors.New("end sequence ID is > max int64")
	ErrNodesCountTooLarge       = errors.New("nodes count is > max int32")
	ErrActiveNodeIDTooLarge     = errors.New("active node ID is > max int32")
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

	var (
		originatorNodeID int32
		startSequenceID  int64
		endSequenceID    int64
		activeNodeIDs    []int32
	)

	// The originator node ID is stored as an int32 in the database, but
	// a uint32 on the network. Do not allow anything larger than max int32
	if originatorNodeID, err = utils.Uint32ToInt32(report.OriginatorNodeID); err != nil {
		return nil, ErrOriginatorNodeIDTooLarge
	}

	if startSequenceID, err = utils.Uint64ToInt64(report.StartSequenceID); err != nil {
		return nil, ErrStartSequenceIDTooLarge
	}

	if endSequenceID, err = utils.Uint64ToInt64(report.EndSequenceID); err != nil {
		return nil, ErrEndSequenceIDTooLarge
	}

	if activeNodeIDs, err = utils.Uint32SliceToInt32Slice(report.ActiveNodeIDs); err != nil {
		return nil, ErrActiveNodeIDTooLarge
	}

	err = s.queries.InsertOrIgnorePayerReport(ctx, queries.InsertOrIgnorePayerReportParams{
		ID:               id,
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  startSequenceID,
		EndSequenceID:    endSequenceID,
		PayersMerkleRoot: report.PayersMerkleRoot[:],
		ActiveNodeIds:    activeNodeIDs,
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
	// Validate NodeID (assuming it should fit within int32 range for consistency with other node IDs)
	if attestation.NodeSignature.NodeID > math.MaxInt32 {
		return ErrOriginatorNodeIDTooLarge
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
	SubmissionStatusIn  []SubmissionStatus
	AttestationStatusIn []AttestationStatus
	StartSequenceID     *uint64
	EndSequenceID       *uint64
	CreatedAfter        time.Time
	OriginatorNodeID    *uint32
}

func (f *FetchReportsQuery) toParams() queries.FetchPayerReportsParams {
	return queries.FetchPayerReportsParams{
		CreatedAfter:        utils.NewNullTime(f.CreatedAfter),
		SubmissionStatusIn:  utils.NewNullInt16Slice(f.SubmissionStatusIn),
		AttestationStatusIn: utils.NewNullInt16Slice(f.AttestationStatusIn),
		StartSequenceID:     utils.NewNullInt64(f.StartSequenceID),
		EndSequenceID:       utils.NewNullInt64(f.EndSequenceID),
		OriginatorNodeID:    utils.NewNullInt32(f.OriginatorNodeID),
	}
}

func NewFetchReportsQuery() *FetchReportsQuery {
	return &FetchReportsQuery{}
}

func (f *FetchReportsQuery) WithSubmissionStatus(statuses ...SubmissionStatus) *FetchReportsQuery {
	f.SubmissionStatusIn = append(f.SubmissionStatusIn, statuses...)
	return f
}

func (f *FetchReportsQuery) WithAttestationStatus(
	statuses ...AttestationStatus,
) *FetchReportsQuery {
	f.AttestationStatusIn = append(f.AttestationStatusIn, statuses...)
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

func (f *FetchReportsQuery) WithOriginatorNodeID(originatorNodeID uint32) *FetchReportsQuery {
	f.OriginatorNodeID = &originatorNodeID
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
		id               [32]byte
	)

	if payersMerkleRoot, err = utils.SliceToArray32(report.PayersMerkleRoot); err != nil {
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
			ActiveNodeIDs:    utils.Int32SliceToUint32Slice(report.ActiveNodeIds),
		},
	}, nil
}

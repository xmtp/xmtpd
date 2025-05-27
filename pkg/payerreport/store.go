package payerreport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var (
	ErrNoActiveNodeIDs             = errors.New("no active node IDs")
	ErrInvalidReportID             = errors.New("invalid report ID")
	ErrOriginatorNodeIDTooLarge    = errors.New("originator node ID is > max int32")
	ErrStartSequenceIDTooLarge     = errors.New("start sequence ID is > max int64")
	ErrEndSequenceIDTooLarge       = errors.New("end sequence ID is > max int64")
	ErrEndMinuteSinceEpochTooLarge = errors.New("end minute since epoch is > max int32")
	ErrNodesCountTooLarge          = errors.New("nodes count is > max int32")
	ErrActiveNodeIDTooLarge        = errors.New("active node ID is > max int32")
	ErrReportNotFound              = errors.New("report not found")
	ErrReportNil                   = errors.New("report is nil")
)

type Store struct {
	queries *queries.Queries
	db      *sql.DB
	log     *zap.Logger
}

func NewStore(db *sql.DB, log *zap.Logger) *Store {
	return &Store{
		queries: queries.New(db),
		db:      db,
		log:     log.Named("payerreportstore"),
	}
}

// Store a report in the database. No validations have been performed, and no originator envelope is stored.
// This function is primarily used for testing
func (s *Store) StoreReport(ctx context.Context, report *PayerReport) (*ReportID, error) {
	params, err := prepareStoreReportParams(report)
	if err != nil {
		return nil, err
	}

	_, err = s.queries.InsertOrIgnorePayerReport(ctx, *params)
	if err != nil {
		return nil, err
	}

	return &report.ID, nil
}

// Store an attestation in the database
func (s *Store) StoreAttestation(ctx context.Context, attestation *PayerReportAttestation) error {
	if attestation.Report == nil {
		return ErrReportNil
	}

	reportID := attestation.Report.ID[:]

	// Validate NodeID (assuming it should fit within int32 range for consistency with other node IDs)
	if attestation.NodeSignature.NodeID > math.MaxInt32 {
		return ErrOriginatorNodeIDTooLarge
	}

	return s.queries.InsertOrIgnorePayerReportAttestation(
		ctx,
		queries.InsertOrIgnorePayerReportAttestationParams{
			PayerReportID: reportID[:],
			NodeID:        int64(attestation.NodeSignature.NodeID),
			Signature:     attestation.NodeSignature.Signature,
		},
	)
}

func (s *Store) FetchReport(ctx context.Context, id ReportID) (*PayerReportWithStatus, error) {
	reports, err := s.FetchReports(ctx, NewFetchReportsQuery().WithReportID(id))
	if err != nil {
		return nil, err
	}

	if len(reports) == 0 {
		return nil, ErrReportNotFound
	}

	return reports[0], nil
}

func (s *Store) FetchReports(
	ctx context.Context,
	query *FetchReportsQuery,
) ([]*PayerReportWithStatus, error) {
	params := query.toParams()
	rows, err := s.queries.FetchPayerReports(ctx, params)
	if err != nil {
		return nil, err
	}
	s.log.Info("Fetched reports", zap.Any("rows", rows))

	return convertPayerReports(rows)
}

func (s *Store) SetReportAttestationStatus(
	ctx context.Context,
	id ReportID,
	fromStatus []AttestationStatus,
	toStatus AttestationStatus,
) error {
	allowedPrevStatuses := make([]int16, len(fromStatus))
	for idx, status := range fromStatus {
		allowedPrevStatuses[idx] = int16(status)
	}

	return s.queries.SetReportAttestationStatus(ctx, queries.SetReportAttestationStatusParams{
		ReportID:   id[:],
		NewStatus:  int16(toStatus),
		PrevStatus: allowedPrevStatuses,
	})
}

func (s *Store) CreatePayerReport(
	ctx context.Context,
	report *PayerReport,
	payerEnvelope *envelopes.PayerEnvelope,
) (*ReportID, error) {
	payerEnvelopeBytes, err := payerEnvelope.Bytes()
	if err != nil {
		return nil, err
	}

	payerReportParams, err := prepareStoreReportParams(report)
	if err != nil {
		return nil, err
	}

	if err = db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{},
		func(ctx context.Context, tx *queries.Queries) error {
			var (
				numRows int64
				err     error
			)
			if numRows, err = tx.InsertOrIgnorePayerReport(ctx, *payerReportParams); err != nil {
				return err
			}

			if numRows == 0 {
				return nil
			}

			if _, err = tx.InsertStagedOriginatorEnvelope(
				ctx,
				queries.InsertStagedOriginatorEnvelopeParams{
					Topic:         payerEnvelope.TargetTopic().Bytes(),
					PayerEnvelope: payerEnvelopeBytes,
				},
			); err != nil {
				return err
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return &report.ID, nil
}

// To be called by the attester of a report. Updates the attestation status of the report and
// writes a staged originator envelope to the database that can be synced to other nodes.
func (s *Store) CreateAttestation(
	ctx context.Context,
	attestation *PayerReportAttestation,
	payerEnvelope *envelopes.PayerEnvelope,
) error {
	allowedPrevStatuses := []int16{int16(AttestationPending)}
	toStatus := AttestationApproved
	targetTopic := payerEnvelope.TargetTopic().Bytes()
	reportID := attestation.Report.ID[:]

	envelopeBytes, err := payerEnvelope.Bytes()
	if err != nil {
		return err
	}

	if attestation.Report == nil {
		return ErrReportNil
	}

	return db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{},
		func(ctx context.Context, tx *queries.Queries) error {
			if err = tx.InsertOrIgnorePayerReportAttestation(
				ctx,
				queries.InsertOrIgnorePayerReportAttestationParams{
					PayerReportID: reportID,
					NodeID:        int64(attestation.NodeSignature.NodeID),
					Signature:     attestation.NodeSignature.Signature,
				},
			); err != nil {
				return err
			}

			err := tx.SetReportAttestationStatus(ctx, queries.SetReportAttestationStatusParams{
				ReportID:   reportID,
				NewStatus:  int16(toStatus),
				PrevStatus: allowedPrevStatuses,
			})
			if err != nil {
				return err
			}

			if _, err = tx.InsertStagedOriginatorEnvelope(ctx, queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         targetTopic,
				PayerEnvelope: envelopeBytes,
			}); err != nil {
				return err
			}

			return nil
		},
	)
}

// Store a report that has been received through a stream from another node
func (s *Store) StoreSyncedReport(
	ctx context.Context,
	envelope *envelopes.OriginatorEnvelope,
	payerID int32,
	domainSeparator common.Hash,
) error {
	originatorEnvelopeBytes, err := envelope.Bytes()
	if err != nil {
		return err
	}

	clientEnvelope := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope
	payerReportProtoWrapper, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_PayerReport)
	if !ok {
		return errors.New("payload is not a payer report")
	}

	payerReportProto := payerReportProtoWrapper.PayerReport

	reportID, err := BuildPayerReportID(
		payerReportProto.OriginatorNodeId,
		payerReportProto.StartSequenceId,
		payerReportProto.EndSequenceId,
		common.BytesToHash(payerReportProto.PayersMerkleRoot),
		payerReportProto.ActiveNodeIds,
		domainSeparator,
	)
	if err != nil {
		return err
	}

	payerReport := &PayerReport{
		ID:                  *reportID,
		OriginatorNodeID:    payerReportProto.OriginatorNodeId,
		StartSequenceID:     payerReportProto.StartSequenceId,
		EndSequenceID:       payerReportProto.EndSequenceId,
		EndMinuteSinceEpoch: payerReportProto.EndMinuteSinceEpoch,
		PayersMerkleRoot:    [32]byte(payerReportProto.PayersMerkleRoot),
		ActiveNodeIDs:       payerReportProto.ActiveNodeIds,
	}

	storeReportParams, err := prepareStoreReportParams(payerReport)
	if err != nil {
		return err
	}

	return db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{},
		func(ctx context.Context, txQueries *queries.Queries) error {
			_, err := txQueries.InsertGatewayEnvelope(
				ctx,
				queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     int32(envelope.OriginatorNodeID()),
					OriginatorSequenceID: int64(envelope.OriginatorSequenceID()),
					Topic:                envelope.TargetTopic().Bytes(),
					OriginatorEnvelope:   originatorEnvelopeBytes,
					PayerID:              db.NullInt32(payerID),
					Expiry: db.NullInt64(int64(
						envelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
					)),
				},
			)
			if err != nil {
				return err
			}

			_, err = txQueries.InsertOrIgnorePayerReport(ctx, *storeReportParams)

			return err
		},
	)
}

func (s *Store) StoreSyncedAttestation(
	ctx context.Context,
	envelope *envelopes.OriginatorEnvelope,
	payerID int32,
) error {
	originatorEnvelopeBytes, err := envelope.Bytes()
	if err != nil {
		return err
	}

	clientEnvelope := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope
	attestationProtoWrapper, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_PayerReportAttestation)
	if !ok {
		return errors.New("payload is not a payer report")
	}

	attestationProto := attestationProtoWrapper.PayerReportAttestation

	return db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{},
		func(ctx context.Context, txQueries *queries.Queries) error {
			_, err := txQueries.InsertGatewayEnvelope(
				ctx,
				queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     int32(envelope.OriginatorNodeID()),
					OriginatorSequenceID: int64(envelope.OriginatorSequenceID()),
					Topic:                envelope.TargetTopic().Bytes(),
					OriginatorEnvelope:   originatorEnvelopeBytes,
					PayerID:              db.NullInt32(payerID),
					Expiry: db.NullInt64(
						int64(envelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()),
					),
				},
			)
			if err != nil {
				return err
			}

			return txQueries.InsertOrIgnorePayerReportAttestation(
				ctx,
				queries.InsertOrIgnorePayerReportAttestationParams{
					PayerReportID: attestationProto.ReportId,
					NodeID:        int64(attestationProto.Signature.NodeId),
					Signature:     attestationProto.Signature.Signature.Bytes,
				},
			)
		},
	)
}

func (s *Store) Queries() *queries.Queries {
	return s.queries
}

func convertPayerReports(rows []queries.FetchPayerReportsRow) ([]*PayerReportWithStatus, error) {
	results := make(map[string]*PayerReportWithStatus)
	var err error
	for _, row := range rows {
		key := fmt.Sprintf("%x", row.ID)
		_, hasExisting := results[key]
		if !hasExisting {
			results[key], err = convertPayerReport(&queries.PayerReport{
				ID:                  row.ID,
				OriginatorNodeID:    row.OriginatorNodeID,
				StartSequenceID:     row.StartSequenceID,
				EndSequenceID:       row.EndSequenceID,
				EndMinuteSinceEpoch: row.EndMinuteSinceEpoch,
				PayersMerkleRoot:    row.PayersMerkleRoot,
				ActiveNodeIds:       row.ActiveNodeIds,
				CreatedAt:           row.CreatedAt,
				SubmissionStatus:    row.SubmissionStatus,
				AttestationStatus:   row.AttestationStatus,
			})
			if err != nil {
				return nil, err
			}
		}
		if row.NodeID.Valid && len(row.Signature) > 0 {
			results[key].AttestationSignatures = append(
				results[key].AttestationSignatures,
				NodeSignature{
					NodeID:    uint32(row.NodeID.Int64),
					Signature: row.Signature,
				},
			)
		}
	}

	out := make([]*PayerReportWithStatus, 0, len(results))
	for _, result := range results {
		out = append(out, result)
	}
	return out, nil
}

func convertPayerReport(report *queries.PayerReport) (*PayerReportWithStatus, error) {
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
		PayerReport: PayerReport{
			ID:                  ReportID(id),
			OriginatorNodeID:    uint32(report.OriginatorNodeID),
			StartSequenceID:     uint64(report.StartSequenceID),
			EndSequenceID:       uint64(report.EndSequenceID),
			EndMinuteSinceEpoch: uint32(report.EndMinuteSinceEpoch),
			PayersMerkleRoot:    payersMerkleRoot,
			ActiveNodeIDs:       utils.Int32SliceToUint32Slice(report.ActiveNodeIds),
		},
	}, nil
}

func prepareStoreReportParams(
	report *PayerReport,
) (*queries.InsertOrIgnorePayerReportParams, error) {
	var (
		err                 error
		originatorNodeID    int32
		endMinuteSinceEpoch int32
		startSequenceID     int64
		endSequenceID       int64
		activeNodeIDs       []int32
	)

	if len(report.ID) != 32 {
		return nil, ErrInvalidReportID
	}

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

	if endMinuteSinceEpoch, err = utils.Uint32ToInt32(report.EndMinuteSinceEpoch); err != nil {
		return nil, ErrEndMinuteSinceEpochTooLarge
	}

	if len(activeNodeIDs) == 0 {
		return nil, ErrNoActiveNodeIDs
	}

	return &queries.InsertOrIgnorePayerReportParams{
		ID:                  report.ID[:],
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     startSequenceID,
		EndSequenceID:       endSequenceID,
		EndMinuteSinceEpoch: endMinuteSinceEpoch,
		PayersMerkleRoot:    report.PayersMerkleRoot[:],
		ActiveNodeIds:       activeNodeIDs,
	}, nil
}

// Package payerreport implements the payer report generator and verifier interfaces.
package payerreport

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
)

type NodeSignature struct {
	NodeID    uint32
	Signature []byte
}

type PayerReportGenerationParams struct {
	OriginatorID            uint32
	LastReportEndSequenceID uint64
	NumHours                int
}

type IPayerReportGenerator interface {
	GenerateReport(
		ctx context.Context,
		params PayerReportGenerationParams,
	) (*PayerReportWithInputs, error)
}

type VerifyReportResult struct {
	IsValid bool
	Reason  string
}

type IPayerReportVerifier interface {
	VerifyReport(
		ctx context.Context,
		prevReport *PayerReport,
		newReport *PayerReport,
	) (VerifyReportResult, error)
	GetPayerMap(
		ctx context.Context,
		report *PayerReport,
	) (PayerMap, error)
}

type IPayerReportStore interface {
	StoreReport(ctx context.Context, report *PayerReport) (int64, error)
	CreatePayerReport(
		ctx context.Context,
		report *PayerReport,
		payerEnvelope *envelopes.PayerEnvelope,
	) (*ReportID, error)
	CreateAttestation(
		ctx context.Context,
		attestation *PayerReportAttestation,
		payerEnvelope *envelopes.PayerEnvelope,
	) error
	FetchReport(ctx context.Context, id ReportID) (*PayerReportWithStatus, error)
	FetchReports(ctx context.Context, query *FetchReportsQuery) ([]*PayerReportWithStatus, error)
	StoreSyncedReport(
		ctx context.Context,
		envelope *envelopes.OriginatorEnvelope,
		payerID int32,
		domainSeparator common.Hash,
	) error
	StoreSyncedAttestation(
		ctx context.Context,
		envelope *envelopes.OriginatorEnvelope,
		payerID int32,
	) error
	SetReportSubmitted(ctx context.Context, id ReportID, reportIndex int32) error
	ForceSetReportSubmitted(ctx context.Context, id ReportID, reportIndex int32) error
	SetReportSettled(ctx context.Context, id ReportID) error
	SetReportSubmissionRejected(ctx context.Context, id ReportID) error
	SetReportAttestationApproved(ctx context.Context, id ReportID) error
	SetReportAttestationRejected(ctx context.Context, id ReportID) error
	Queries() *queries.Queries
	GetAdvisoryLocker(
		ctx context.Context,
	) (db.ITransactionScopedAdvisoryLocker, error)
	GetLatestSequenceID(ctx context.Context, originatorNodeID int32) (int64, error)
}

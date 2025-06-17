package payerreport

import (
	"context"

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

type IPayerReportVerifier interface {
	IsValidReport(
		ctx context.Context,
		prevReport *PayerReport,
		newReport *PayerReport,
	) (bool, error)
}

type IPayerReportStore interface {
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
	SetReportAttestationStatus(
		ctx context.Context,
		id ReportID,
		fromStatus []AttestationStatus,
		toStatus AttestationStatus,
	) error
	Queries() *queries.Queries
}

package payerreport

import (
	"context"
)

type NodeSignature struct {
	NodeID    uint32
	Signature []byte
}

type PayerReportAttestation struct {
	Report        *PayerReport
	NodeSignature NodeSignature
}

type PayerReportGenerationParams struct {
	OriginatorID            uint32
	LastReportEndSequenceID uint64
	NumHours                int
}

type PayerReportGenerator interface {
	GenerateReport(
		ctx context.Context,
		params PayerReportGenerationParams,
	) (*PayerReportWithInputs, error)
}

type PayerReportVerifier interface {
	IsValidReport(
		ctx context.Context,
		prevReport *PayerReport,
		newReport *PayerReport,
	) (bool, error)
}

type IPayerReportManager interface {
	PayerReportGenerator
	PayerReportVerifier
}

type IPayerReportStore interface {
	StoreReport(ctx context.Context, report *PayerReport) (ReportID, error)
	FetchReport(ctx context.Context, id ReportID) (*PayerReportWithStatus, error)
	FetchReports(ctx context.Context, query *FetchReportsQuery) ([]*PayerReportWithStatus, error)
	StoreAttestation(ctx context.Context, attestation *PayerReportAttestation) error
	SetReportAttestationStatus(
		ctx context.Context,
		id ReportID,
		fromStatus []AttestationStatus,
		toStatus AttestationStatus,
	) error
}

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
	) (*FullPayerReport, error)
}

type PayerReportAttester interface {
	AttestReport(
		ctx context.Context,
		prevReport *PayerReport,
		newReport *PayerReport,
	) (*PayerReportAttestation, error)
}

type IPayerReportManager interface {
	PayerReportGenerator
	PayerReportAttester
}

type IPayerReportStore interface {
	StoreReport(ctx context.Context, report *PayerReport) (ReportID, error)
	FetchReport(ctx context.Context, id ReportID) (*PayerReport, error)
	StoreAttestation(ctx context.Context, attestation *PayerReportAttestation) error
}

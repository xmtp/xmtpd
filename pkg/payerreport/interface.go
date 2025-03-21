package payerreport

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
)

type PayerReport struct {
	// The Originator Node that the report is about
	OriginatorNodeID uint32
	// The report applies to messages with sequence IDs > StartSequenceID
	StartSequenceID uint64
	// The report applies to messages with sequence IDs <= EndSequenceID
	EndSequenceID uint64
	// The payers in the report and the number of messages they paid for
	Payers map[common.Address]currency.PicoDollar
	// The merkle root of the Payers mapping
	PayersMerkleRoot []byte
	// The number of leaves in the Payers merkle tree
	PayersLeafCount uint32
}

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

type IPayerReportManager interface {
	GenerateReport(ctx context.Context, params PayerReportGenerationParams) (*PayerReport, error)
	AttestReport(
		ctx context.Context,
		prevReport *PayerReport,
		newReport *PayerReport,
	) (*PayerReportAttestation, error)
}

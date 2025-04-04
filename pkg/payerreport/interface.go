package payerreport

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
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
	// The merkle root of the Nodes included in the report
	NodesHash []byte
	// The number of leaves in the Nodes merkle tree
	NodesCount uint32
}

func (p *PayerReport) ToProto() *proto.PayerReport {
	return &proto.PayerReport{
		OriginatorNodeId: p.OriginatorNodeID,
		StartSequenceId:  p.StartSequenceID,
		EndSequenceId:    p.EndSequenceID,
		PayersMerkleRoot: p.PayersMerkleRoot,
		NodesHash:        p.NodesHash,
		PayersLeafCount:  p.PayersLeafCount,
		NodesCount:       p.NodesCount,
	}
}

func (p *PayerReport) ID() ([]byte, error) {
	packedBytes, err := payerReportMessageHash.Pack(
		p.OriginatorNodeID,
		p.StartSequenceID,
		p.EndSequenceID,
		utils.SliceToArray32(p.PayersMerkleRoot),
		p.PayersLeafCount,
		utils.SliceToArray32(p.NodesHash),
		p.NodesCount,
	)
	if err != nil {
		return nil, err
	}
	// Return the keccak256 hash
	return utils.HashPayerReportInput(packedBytes), nil
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

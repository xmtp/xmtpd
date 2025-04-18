package payerreport

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type ReportID []byte

type PayerReport struct {
	// The Originator Node that the report is about
	OriginatorNodeID uint32
	// The report applies to messages with sequence IDs > StartSequenceID
	StartSequenceID uint64
	// The report applies to messages with sequence IDs <= EndSequenceID
	EndSequenceID uint64
	// The merkle root of the Payers mapping
	PayersMerkleRoot [32]byte
	// The number of leaves in the Payers merkle tree
	PayersLeafCount uint32
	// The hash of all the nodes included in the report, sorted
	NodesHash [32]byte
	// The number of nodes included in the report
	NodesCount uint32
}

type FullPayerReport struct {
	PayerReport
	// The payers in the report and the number of messages they paid for
	Payers  map[common.Address]currency.PicoDollar
	NodeIDs []uint32
}

func (p *PayerReport) ToProto() *proto.PayerReport {
	return &proto.PayerReport{
		OriginatorNodeId: p.OriginatorNodeID,
		StartSequenceId:  p.StartSequenceID,
		EndSequenceId:    p.EndSequenceID,
		PayersMerkleRoot: p.PayersMerkleRoot[:],
		PayersLeafCount:  p.PayersLeafCount,
	}
}

func (p *PayerReport) ID() (ReportID, error) {
	packedBytes, err := payerReportMessageHash.Pack(
		p.OriginatorNodeID,
		p.StartSequenceID,
		p.EndSequenceID,
		p.PayersMerkleRoot,
		p.PayersLeafCount,
		p.NodesHash,
		p.NodesCount,
	)
	if err != nil {
		return nil, err
	}
	// Return the keccak256 hash
	return utils.HashPayerReportInput(packedBytes), nil
}

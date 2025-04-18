package payerreport

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// The arguments to use for hashing the payer report ID both on and off chain
var payerReportMessageHash = abi.Arguments{
	{
		Name: "originatorNodeID",
		Type: abi.Type{T: abi.UintTy, Size: 32},
	},
	{
		Name: "startSequenceID",
		Type: abi.Type{T: abi.UintTy, Size: 64},
	},
	{
		Name: "endSequenceID",
		Type: abi.Type{T: abi.UintTy, Size: 64},
	},
	{
		Name: "payersMerkleRoot",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
	{
		Name: "payersLeafCount",
		Type: abi.Type{T: abi.UintTy, Size: 32},
	},
	{
		Name: "nodesHash",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
	{
		Name: "nodesCount",
		Type: abi.Type{T: abi.UintTy, Size: 32},
	},
}

type SubmissionStatus int16

const (
	SubmissionPending   SubmissionStatus = iota
	SubmissionSubmitted                  = 1
	SubmissionSettled                    = 2
)

type AttestationStatus int16

const (
	AttestationPending  AttestationStatus = iota
	AttestationApproved                   = 1
	AttestationRejected                   = 2
)

type ReportID []byte

func (r ReportID) String() string {
	return hex.EncodeToString(r)
}

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

// A FullPayerReport is a superset of a PayerReport that includes the payers and node IDs
type PayerReportWithInputs struct {
	PayerReport
	// The payers in the report and the number of messages they paid for
	Payers  map[common.Address]currency.PicoDollar
	NodeIDs []uint32
}

type PayerReportWithStatus struct {
	PayerReport
	SubmissionStatus  SubmissionStatus
	AttestationStatus AttestationStatus
	CreatedAt         time.Time
	ID                [32]byte
}

func (p *PayerReport) ToProto() *proto.PayerReport {
	return &proto.PayerReport{
		OriginatorNodeId: p.OriginatorNodeID,
		StartSequenceId:  p.StartSequenceID,
		EndSequenceId:    p.EndSequenceID,
		PayersMerkleRoot: p.PayersMerkleRoot[:],
		PayersLeafCount:  p.PayersLeafCount,
		NodesHash:        p.NodesHash[:],
		NodesCount:       p.NodesCount,
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

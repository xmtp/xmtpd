package payerreport

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/merkle"
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
		Name: "activeNodeIds",
		Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.UintTy, Size: 32}},
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

type ReportID [32]byte

func (r ReportID) String() string {
	return hex.EncodeToString(r[:])
}

type PayerReport struct {
	// The Originator Node that the report is about
	OriginatorNodeID uint32
	// The report applies to messages with sequence IDs > StartSequenceID
	StartSequenceID uint64
	// The report applies to messages with sequence IDs <= EndSequenceID
	EndSequenceID uint64
	// The timestamp of the message at EndSequenceID
	EndMinuteSinceEpoch uint32
	// The merkle root of the Payers mapping
	PayersMerkleRoot [32]byte
	// The active node IDs in the report
	ActiveNodeIDs []uint32
}

type NewPayerReportParams struct {
	OriginatorNodeID uint32
	StartSequenceID  uint64
	EndSequenceID    uint64
	Payers           map[common.Address]currency.PicoDollar
	NodeIDs          []uint32
}

// A FullPayerReport is a superset of a PayerReport that includes the payers and node IDs
type PayerReportWithInputs struct {
	PayerReport
	// The payers in the report and the number of messages they paid for
	Payers     map[common.Address]currency.PicoDollar
	MerkleTree *merkle.MerkleTree
	NodeIDs    []uint32
}

type PayerReportWithStatus struct {
	PayerReport
	AttestationSignatures []NodeSignature
	// Whether the report has been submitted to the blockchain or not
	SubmissionStatus SubmissionStatus
	// Status of the current node's attestation of the report
	AttestationStatus AttestationStatus
	// The timestamp of when the report was inserted into the node's database
	CreatedAt time.Time
	// The ID of the report
	ID ReportID
}

func (p *PayerReport) ToProto() *proto.PayerReport {
	return &proto.PayerReport{
		OriginatorNodeId:    p.OriginatorNodeID,
		StartSequenceId:     p.StartSequenceID,
		EndSequenceId:       p.EndSequenceID,
		PayersMerkleRoot:    p.PayersMerkleRoot[:],
		ActiveNodeIds:       p.ActiveNodeIDs,
		EndMinuteSinceEpoch: p.EndMinuteSinceEpoch,
	}
}

func (p *PayerReport) ID() (ReportID, error) {
	packedBytes, err := payerReportMessageHash.Pack(
		p.OriginatorNodeID,
		p.StartSequenceID,
		p.EndSequenceID,
		p.PayersMerkleRoot,
		p.ActiveNodeIDs,
	)
	if err != nil {
		return ReportID{}, err
	}
	hash, err := utils.SliceToArray32(utils.HashPayerReportInput(packedBytes))
	if err != nil {
		return ReportID{}, err
	}
	return ReportID(hash), nil
}

func NewPayerReport(params NewPayerReportParams) (*PayerReportWithInputs, error) {
	merkleTree, err := NewPayerMerkleTree(params.Payers)
	if err != nil {
		return nil, err
	}

	merkleRoot, err := utils.SliceToArray32(merkleTree.Root())
	if err != nil {
		return nil, err
	}

	return &PayerReportWithInputs{
		PayerReport: PayerReport{
			OriginatorNodeID: params.OriginatorNodeID,
			StartSequenceID:  params.StartSequenceID,
			EndSequenceID:    params.EndSequenceID,
			PayersMerkleRoot: merkleRoot,
			PayersLeafCount:  uint32(len(params.Payers)),
			NodesCount:       uint32(len(params.NodeIDs)),
		},
		Payers:     params.Payers,
		NodeIDs:    params.NodeIDs,
		MerkleTree: merkleTree,
	}, nil
}

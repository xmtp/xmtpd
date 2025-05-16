package payerreport

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/merkle"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
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
	ID ReportID
	// The Originator Node that the report is about
	OriginatorNodeID uint32
	// The report applies to messages with sequence IDs > StartSequenceID
	StartSequenceID uint64
	// The report applies to messages with sequence IDs <= EndSequenceID
	EndSequenceID uint64
	// The timestamp of the message at EndSequenceID
	EndMinuteSinceEpoch uint32
	// The merkle root of the Payers mapping
	PayersMerkleRoot common.Hash
	// The active node IDs in the report
	ActiveNodeIDs []uint32
}

type BuildPayerReportParams struct {
	OriginatorNodeID    uint32
	StartSequenceID     uint64
	EndSequenceID       uint64
	EndMinuteSinceEpoch uint32
	Payers              map[common.Address]currency.PicoDollar
	NodeIDs             []uint32
	DomainSeparator     common.Hash
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

func (p *PayerReport) ToClientEnvelope() (*envelopes.ClientEnvelope, error) {
	payload := p.ToProto()
	targetTopic := topic.NewTopic(topic.TOPIC_KIND_PAYER_REPORTS_V1, utils.Uint32ToBytes(p.OriginatorNodeID)).
		Bytes()

	return envelopes.NewClientEnvelope(&proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_PayerReport{
			PayerReport: payload,
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: targetTopic,
		},
	})
}

func buildPayerReportID(
	originatorNodeID uint32,
	startSequenceID uint64,
	endSequenceID uint64,
	payersMerkleRoot common.Hash,
	activeNodeIDs []uint32,
	domainSeparator common.Hash,
) (*ReportID, error) {
	if domainSeparator == (common.Hash{}) {
		return nil, errors.New("domain separator is required")
	}

	packedBytes, err := payerReportMessageHash.Pack(
		originatorNodeID,
		startSequenceID,
		endSequenceID,
		payersMerkleRoot,
		activeNodeIDs,
	)
	if err != nil {
		return nil, err
	}
	hash := utils.HashPayerReportInput(packedBytes, domainSeparator)
	reportID := ReportID(hash)
	return &reportID, nil
}

func BuildPayerReport(params BuildPayerReportParams) (*PayerReportWithInputs, error) {
	merkleRoot := buildMerkleRoot(params.Payers)

	reportID, err := buildPayerReportID(
		params.OriginatorNodeID,
		params.StartSequenceID,
		params.EndSequenceID,
		merkleRoot,
		params.NodeIDs,
		params.DomainSeparator,
	)
	if err != nil {
		return nil, err
	}

	return &PayerReportWithInputs{
		PayerReport: PayerReport{
			ID:                  *reportID,
			OriginatorNodeID:    params.OriginatorNodeID,
			StartSequenceID:     params.StartSequenceID,
			EndSequenceID:       params.EndSequenceID,
			EndMinuteSinceEpoch: params.EndMinuteSinceEpoch,
			PayersMerkleRoot:    merkleRoot,
			ActiveNodeIDs:       params.NodeIDs,
		},
		Payers:     params.Payers,
		NodeIDs:    params.NodeIDs,
		MerkleTree: nil,
	}, nil
}

package payerreport

import (
	"encoding/hex"
	"errors"
	"log"
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
		Name: "typeHash",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
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
		Name: "endMinuteSinceEpoch",
		Type: abi.Type{T: abi.UintTy, Size: 32},
	},
	{
		Name: "payersMerkleRoot",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
	{
		Name: "activeNodeIdsHash",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
}

type SubmissionStatus int16

const (
	SubmissionPending   SubmissionStatus = iota
	SubmissionSubmitted                  = 1
	SubmissionSettled                    = 2
	SubmissionRejected                   = 3
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
	// The index of the report on the blockchain (null if not yet submitted or submission status is not tracked)
	SubmittedReportIndex *uint32
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

// PayerReportWithInputs is a superset of a PayerReport that includes the payers and node IDs.
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
	targetTopic := topic.NewTopic(topic.TopicKindPayerReportsV1, utils.Uint32ToBytes(p.OriginatorNodeID)).
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

func BuildPayerReportID(
	originatorNodeID uint32,
	startSequenceID uint64,
	endSequenceID uint64,
	endMinuteSinceEpoch uint32,
	payersMerkleRoot common.Hash,
	activeNodeIDs []uint32,
	domainSeparator common.Hash,
) (*ReportID, error) {
	if domainSeparator == (common.Hash{}) {
		return nil, errors.New("domain separator is required")
	}

	nodeIdsHash, err := utils.PackSortAndHashNodeIDs(activeNodeIDs)
	if err != nil {
		log.Printf("error packing node IDs: %v\n", err)
		return nil, err
	}

	packedBytes, err := payerReportMessageHash.Pack(
		payerReportDigestTypeHash,
		originatorNodeID,
		startSequenceID,
		endSequenceID,
		endMinuteSinceEpoch,
		payersMerkleRoot,
		nodeIdsHash,
	)
	if err != nil {
		log.Printf("error packing payer report message hash: %v\n", err)
		return nil, err
	}
	hash := utils.HashPayerReportInput(packedBytes, domainSeparator)
	reportID := ReportID(hash)
	return &reportID, nil
}

func BuildPayerReport(params BuildPayerReportParams) (*PayerReportWithInputs, error) {
	tree, err := GenerateMerkleTree(params.Payers)
	if err != nil {
		return nil, err
	}
	merkleRoot := common.BytesToHash(tree.Root())

	reportID, err := BuildPayerReportID(
		params.OriginatorNodeID,
		params.StartSequenceID,
		params.EndSequenceID,
		params.EndMinuteSinceEpoch,
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
		MerkleTree: tree,
	}, nil
}

package payerreport

import (
	"github.com/xmtp/xmtpd/pkg/envelopes"
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type PayerReportAttestation struct {
	Report        *PayerReport
	NodeSignature NodeSignature
}

func NewPayerReportAttestation(
	report *PayerReport,
	nodeSignature NodeSignature,
) *PayerReportAttestation {
	return &PayerReportAttestation{
		Report:        report,
		NodeSignature: nodeSignature,
	}
}

func (a *PayerReportAttestation) ToProto() *proto.PayerReportAttestation {
	return &proto.PayerReportAttestation{
		ReportId: a.Report.ID[:],
		Signature: &proto.NodeSignature{
			NodeId: a.NodeSignature.NodeID,
			Signature: &associations.RecoverableEcdsaSignature{
				Bytes: a.NodeSignature.Signature,
			},
		},
	}
}

func (a *PayerReportAttestation) ToClientEnvelope() (*envelopes.ClientEnvelope, error) {
	attestationProto := a.ToProto()

	targetTopic := topic.NewTopic(
		topic.TopicKindPayerReportAttestationsV1,
		utils.Uint32ToBytes(a.Report.OriginatorNodeID),
	)

	return envelopes.NewClientEnvelope(
		&proto.ClientEnvelope{
			Aad: &proto.AuthenticatedData{
				TargetTopic: targetTopic.Bytes(),
			},
			Payload: &proto.ClientEnvelope_PayerReportAttestation{
				PayerReportAttestation: attestationProto,
			},
		},
	)
}

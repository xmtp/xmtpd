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

func (a *PayerReportAttestation) ToProto() (*proto.PayerReportAttestation, error) {
	reportID, err := a.Report.ID()
	if err != nil {
		return nil, err
	}

	return &proto.PayerReportAttestation{
		ReportId: reportID[:],
		Signature: &proto.NodeSignature{
			NodeId: a.NodeSignature.NodeID,
			Signature: &associations.RecoverableEcdsaSignature{
				Bytes: a.NodeSignature.Signature,
			},
		},
	}, nil
}

func (a *PayerReportAttestation) ToClientEnvelope() (*envelopes.ClientEnvelope, error) {
	attestationProto, err := a.ToProto()
	if err != nil {
		return nil, err
	}

	targetTopic := topic.NewTopic(
		topic.TOPIC_KIND_PAYER_REPORT_ATTESTATIONS_V1,
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

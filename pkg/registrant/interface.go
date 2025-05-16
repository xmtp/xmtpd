package registrant

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
)

type IRegistrant interface {
	NodeID() uint32
	SignStagedEnvelope(
		stagedEnv queries.StagedOriginatorEnvelope,
		baseFee currency.PicoDollar,
		congestionFee currency.PicoDollar,
		retentionDays uint32,
	) (*envelopes.OriginatorEnvelope, error)
	SignPayerReportAttestation(
		reportID payerreport.ReportID,
		domainSeparator common.Hash,
	) (*payerreport.NodeSignature, error)
	SignClientEnvelopeToSelf(unsignedClientEnvelope []byte) ([]byte, error)
	TokenFactory() authn.TokenFactory
}

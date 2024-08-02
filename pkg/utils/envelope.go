package utils

import (
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/protobuf/proto"
)

func SignStagedEnvelope(
	stagedEnv queries.StagedOriginatorEnvelope,
) (*message_api.OriginatorEnvelope, error) {
	payerEnv := &message_api.PayerEnvelope{}
	if err := proto.Unmarshal(stagedEnv.PayerEnvelope, payerEnv); err != nil {
		return nil, err
	}
	unsignedEnv := message_api.UnsignedOriginatorEnvelope{
		OriginatorSid: SID(stagedEnv.ID),
		OriginatorNs:  stagedEnv.OriginatorTime.UnixNano(),
		PayerEnvelope: payerEnv,
	}
	unsignedBytes, err := proto.Marshal(&unsignedEnv)
	if err != nil {
		return nil, err
	}
	// TODO(rich): Plumb through public key and properly sign envelope
	signedEnv := message_api.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof:                      nil,
	}
	return &signedEnv, nil
}

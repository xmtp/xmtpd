package utils

import (
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"google.golang.org/protobuf/proto"
)

func UnmarshalOriginatorEnvelope(envelope []byte) (*envelopes.OriginatorEnvelope, error) {
	originatorEnvelope := &envelopes.OriginatorEnvelope{}
	err := proto.Unmarshal(envelope, originatorEnvelope)
	if err != nil {
		return nil, err
	}
	return originatorEnvelope, nil
}

func UnmarshalUnsignedEnvelope(
	unsignedEnvelopeBytes []byte,
) (*envelopes.UnsignedOriginatorEnvelope, error) {
	unsignedEnvelope := &envelopes.UnsignedOriginatorEnvelope{}
	err := proto.Unmarshal(unsignedEnvelopeBytes, unsignedEnvelope)
	if err != nil {
		return nil, err
	}
	return unsignedEnvelope, nil
}

func UnmarshalClientEnvelope(envelope []byte) (*envelopes.ClientEnvelope, error) {
	clientEnvelope := &envelopes.ClientEnvelope{}
	err := proto.Unmarshal(envelope, clientEnvelope)
	if err != nil {
		return nil, err
	}
	return clientEnvelope, nil
}

package utils

import (
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/protobuf/proto"
)

func UnmarshalOriginatorEnvelope(envelope []byte) (*message_api.OriginatorEnvelope, error) {
	originatorEnvelope := &message_api.OriginatorEnvelope{}
	err := proto.Unmarshal(envelope, originatorEnvelope)
	if err != nil {
		return nil, err
	}
	return originatorEnvelope, nil
}

func UnmarshalUnsignedEnvelope(
	unsignedEnvelopeBytes []byte,
) (*message_api.UnsignedOriginatorEnvelope, error) {
	unsignedEnvelope := &message_api.UnsignedOriginatorEnvelope{}
	err := proto.Unmarshal(unsignedEnvelopeBytes, unsignedEnvelope)
	if err != nil {
		return nil, err
	}
	return unsignedEnvelope, nil
}

func UnmarshalClientEnvelope(envelope []byte) (*message_api.ClientEnvelope, error) {
	clientEnvelope := &message_api.ClientEnvelope{}
	err := proto.Unmarshal(envelope, clientEnvelope)
	if err != nil {
		return nil, err
	}
	return clientEnvelope, nil
}

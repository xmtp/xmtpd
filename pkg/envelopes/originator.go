package envelopes

import (
	"errors"

	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type OriginatorEnvelope struct {
	proto                      *envelopesProto.OriginatorEnvelope
	UnsignedOriginatorEnvelope UnsignedOriginatorEnvelope
}

func NewOriginatorEnvelope(proto *envelopesProto.OriginatorEnvelope) (*OriginatorEnvelope, error) {
	if proto == nil {
		return nil, errors.New("originator envelope proto is nil")
	}

	unsigned, err := NewUnsignedOriginatorEnvelopeFromBytes(proto.UnsignedOriginatorEnvelope)
	if err != nil {
		return nil, err
	}

	return &OriginatorEnvelope{
		proto:                      proto,
		UnsignedOriginatorEnvelope: *unsigned,
	}, nil
}

func NewOriginatorEnvelopeFromBytes(bytes []byte) (*OriginatorEnvelope, error) {
	proto, err := utils.UnmarshalOriginatorEnvelope(bytes)
	if err != nil {
		return nil, err
	}
	return NewOriginatorEnvelope(proto)
}

func (o *OriginatorEnvelope) Bytes() ([]byte, error) {
	bytes, err := proto.Marshal(o.proto)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (o *OriginatorEnvelope) Proto() *envelopesProto.OriginatorEnvelope {
	return o.proto
}

func (o *OriginatorEnvelope) OriginatorNodeID() uint32 {
	return o.UnsignedOriginatorEnvelope.OriginatorNodeID()
}

func (o *OriginatorEnvelope) OriginatorSequenceID() uint64 {
	return o.UnsignedOriginatorEnvelope.OriginatorSequenceID()
}

func (o *OriginatorEnvelope) TargetTopic() topic.Topic {
	return o.UnsignedOriginatorEnvelope.TargetTopic()
}

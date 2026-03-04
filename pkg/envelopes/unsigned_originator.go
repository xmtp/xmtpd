package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/currency"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type UnsignedOriginatorEnvelope struct {
	proto         *envelopesProto.UnsignedOriginatorEnvelope
	PayerEnvelope PayerEnvelope
}

// NewUnsignedOriginatorEnvelope constructs an UnsignedOriginatorEnvelope and perform validations on any child fields.
// Does not verify signatures.
func NewUnsignedOriginatorEnvelope(
	proto *envelopesProto.UnsignedOriginatorEnvelope,
) (*UnsignedOriginatorEnvelope, error) {
	if proto == nil {
		return nil, errors.New("unsigned originator envelopeproto is nil")
	}

	payer, err := NewPayerEnvelopeFromBytes(proto.GetPayerEnvelopeBytes())
	if err != nil {
		return nil, err
	}

	return &UnsignedOriginatorEnvelope{proto: proto, PayerEnvelope: *payer}, nil
}

func (u *UnsignedOriginatorEnvelope) PayerEnvelopeBytes() []byte {
	return u.proto.GetPayerEnvelopeBytes()
}

func (u *UnsignedOriginatorEnvelope) OriginatorNodeID() uint32 {
	// Skip nil check because it is in the constructor
	return u.proto.GetOriginatorNodeId()
}

func (u *UnsignedOriginatorEnvelope) OriginatorSequenceID() uint64 {
	// Skip nil check because it is in the constructor
	return u.proto.GetOriginatorSequenceId()
}

func (u *UnsignedOriginatorEnvelope) OriginatorNs() int64 {
	// Skip nil check because it is in the constructor
	return u.proto.GetOriginatorNs()
}

func (u *UnsignedOriginatorEnvelope) BaseFee() currency.PicoDollar {
	// Skip nil check because it is in the constructor
	return currency.PicoDollar(u.proto.GetBaseFeePicodollars())
}

func (u *UnsignedOriginatorEnvelope) CongestionFee() currency.PicoDollar {
	// Skip nil check because it is in the constructor
	return currency.PicoDollar(u.proto.GetCongestionFeePicodollars())
}

func NewUnsignedOriginatorEnvelopeFromBytes(bytes []byte) (*UnsignedOriginatorEnvelope, error) {
	message, err := utils.UnmarshalUnsignedEnvelope(bytes)
	if err != nil {
		return nil, err
	}
	return NewUnsignedOriginatorEnvelope(message)
}

func (u *UnsignedOriginatorEnvelope) Proto() *envelopesProto.UnsignedOriginatorEnvelope {
	return u.proto
}

func (u *UnsignedOriginatorEnvelope) TargetTopic() topic.Topic {
	return u.PayerEnvelope.TargetTopic()
}

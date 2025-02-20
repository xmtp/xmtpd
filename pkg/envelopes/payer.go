package envelopes

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type PayerEnvelope struct {
	proto            *envelopesProto.PayerEnvelope
	ClientEnvelope   ClientEnvelope
	TargetOriginator uint32
}

func NewPayerEnvelope(proto *envelopesProto.PayerEnvelope) (*PayerEnvelope, error) {
	if proto == nil {
		return nil, errors.New("payer envelope proto is nil")
	}

	clientEnv, err := NewClientEnvelopeFromBytes(proto.UnsignedClientEnvelope)
	if err != nil {
		return nil, err
	}
	return &PayerEnvelope{
		proto:            proto,
		ClientEnvelope:   *clientEnv,
		TargetOriginator: proto.GetTargetOriginator(),
	}, nil
}

func (p *PayerEnvelope) Proto() *envelopesProto.PayerEnvelope {
	return p.proto
}

func (p *PayerEnvelope) Bytes() ([]byte, error) {
	bytes, err := proto.Marshal(p.proto)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (p *PayerEnvelope) RecoverSigner() (*common.Address, error) {
	payerSignature := p.proto.PayerSignature
	if payerSignature == nil {
		return nil, errors.New("payer signature is missing")
	}

	hash := utils.HashPayerSignatureInput(p.proto.TargetOriginator, p.proto.UnsignedClientEnvelope)
	signer, err := ethcrypto.SigToPub(hash, payerSignature.Bytes)
	if err != nil {
		return nil, err
	}

	address := ethcrypto.PubkeyToAddress(*signer)

	return &address, nil
}

func (p *PayerEnvelope) TargetTopic() topic.Topic {
	return p.ClientEnvelope.TargetTopic()
}

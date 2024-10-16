package envelopes

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type PayerEnvelope struct {
	proto          *message_api.PayerEnvelope
	ClientEnvelope ClientEnvelope
}

func NewPayerEnvelope(proto *message_api.PayerEnvelope) (*PayerEnvelope, error) {
	if proto == nil {
		return nil, errors.New("proto is nil")
	}

	clientEnv, err := NewClientEnvelopeFromBytes(proto.UnsignedClientEnvelope)
	if err != nil {
		return nil, err
	}
	return &PayerEnvelope{proto: proto, ClientEnvelope: *clientEnv}, nil
}

func (p *PayerEnvelope) Proto() *message_api.PayerEnvelope {
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

	hash := utils.HashPayerSignatureInput(p.proto.UnsignedClientEnvelope)
	signer, err := ethcrypto.SigToPub(hash, payerSignature.Bytes)
	if err != nil {
		return nil, err
	}

	address := ethcrypto.PubkeyToAddress(*signer)

	return &address, nil
}

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

func NewPayerEnvelopeFromBytes(bytes []byte) (*PayerEnvelope, error) {
	msg := &envelopesProto.PayerEnvelope{}
	if err := proto.Unmarshal(bytes, msg); err != nil {
		return nil, err
	}
	return NewPayerEnvelope(msg)
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

// RecoverSigner recovers the address of the entity that signed the payer envelope.
// This is the delegate (gateway) address when delegation is used.
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

// IsDelegated returns true if this envelope uses delegated signing.
// When delegated, the signer is the gateway but the payer is the user.
func (p *PayerEnvelope) IsDelegated() bool {
	return len(p.proto.DelegatedPayerAddress) > 0
}

// GetDelegatedPayerAddress returns the delegated payer address if set.
// Returns nil if this envelope is not using delegated signing.
func (p *PayerEnvelope) GetDelegatedPayerAddress() *common.Address {
	if !p.IsDelegated() {
		return nil
	}
	addr := common.BytesToAddress(p.proto.DelegatedPayerAddress)
	return &addr
}

// GetActualPayer returns the address that should be charged for this message.
// If delegation is used and valid, returns the delegated payer address.
// Otherwise, returns the signer address (gateway).
// NOTE: Caller should verify delegation validity on-chain before trusting this.
func (p *PayerEnvelope) GetActualPayer() (*common.Address, error) {
	if p.IsDelegated() {
		addr := p.GetDelegatedPayerAddress()
		if addr != nil {
			return addr, nil
		}
	}
	// Fall back to signer (legacy behavior)
	return p.RecoverSigner()
}

func (p *PayerEnvelope) TargetTopic() topic.Topic {
	return p.ClientEnvelope.TargetTopic()
}

func (p *PayerEnvelope) RetentionDays() uint32 {
	return p.proto.MessageRetentionDays
}

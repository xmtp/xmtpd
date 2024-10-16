package utils

import (
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
)

func HashPayerSignatureInput(unsignedClientEnvelope []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.PAYER_DOMAIN_SEPARATION_LABEL),
		unsignedClientEnvelope,
	)
}

func HashJWTSignatureInput(textToSign []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.JWT_DOMAIN_SEPARATION_LABEL),
		textToSign,
	)
}

func HashOriginatorSignatureInput(unsignedOriginatorEnvelope []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.ORIGINATOR_DOMAIN_SEPARATION_LABEL),
		unsignedOriginatorEnvelope,
	)
}

package utils

import (
	"encoding/binary"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
)

func HashPayerSignatureInput(originatorID uint32, unsignedClientEnvelope []byte) []byte {
	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, originatorID)
	return ethcrypto.Keccak256(
		[]byte(constants.TARGET_ORIGINATOR_SEPARATION_LABEL),
		targetBytes,
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

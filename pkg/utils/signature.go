package utils

import (
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func SignClientEnvelope(
	originatorID uint32,
	unsignedClientEnvelope []byte,
	payerPrivateKey *ecdsa.PrivateKey,
) ([]byte, error) {
	hash := HashPayerSignatureInput(originatorID, unsignedClientEnvelope)
	signature, err := ethcrypto.Sign(hash, payerPrivateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

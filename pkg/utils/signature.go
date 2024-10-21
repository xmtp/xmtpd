package utils

import (
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func SignClientEnvelope(
	unsignedClientEnvelope []byte,
	payerPrivateKey *ecdsa.PrivateKey,
) ([]byte, error) {
	hash := HashPayerSignatureInput(unsignedClientEnvelope)
	signature, err := ethcrypto.Sign(hash, payerPrivateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

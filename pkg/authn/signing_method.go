package authn

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	algorithm = "ES256K"
	sigLength = 65
	rLength   = 32
	sLength   = 32
)

var (
	ErrWrongKeyFormat = errors.New("wrong key type")
	ErrBadSignature   = errors.New("bad signature")
	ErrVerification   = errors.New("signature verification failed")
)

// SigningMethodSecp256k1 is the JWT signing method for secp256k1. Inspired by https://github.com/ureeves/jwt-go-secp256k1/blob/master/secp256k1.go
// but updated to work with the latest serverVersion of jwt-go.
type SigningMethodSecp256k1 struct{}

func (sm *SigningMethodSecp256k1) Verify(signingString string, sig []byte, key interface{}) error {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return ErrWrongKeyFormat
	}

	hashedString := utils.HashJWTSignatureInput([]byte(signingString))

	if len(sig) != sigLength {
		return ErrBadSignature
	}

	r := new(big.Int).SetBytes(sig[:rLength])                  // R
	s := new(big.Int).SetBytes(sig[rLength : rLength+sLength]) // S

	if !ecdsa.Verify(pub, hashedString, r, s) {
		return ErrVerification
	}

	return nil
}

func (sm *SigningMethodSecp256k1) Sign(signingString string, key interface{}) ([]byte, error) {
	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrWrongKeyFormat
	}

	hashedString := utils.HashJWTSignatureInput([]byte(signingString))

	sig, err := ethcrypto.Sign(hashedString, priv)
	if err != nil {
		return nil, err
	}

	if len(sig) != sigLength {
		return nil, ErrBadSignature
	}

	return sig, nil
}

func (sm *SigningMethodSecp256k1) Alg() string {
	return algorithm
}

func init() {
	method := &SigningMethodSecp256k1{}
	jwt.RegisterSigningMethod(algorithm, func() jwt.SigningMethod {
		return method
	})
}

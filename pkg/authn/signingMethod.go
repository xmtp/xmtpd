package authn

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ALGORITHM               = "ES256K"
	SIG_LENGTH              = 65
	R_LENGTH                = 32
	S_LENGTH                = 32
	DOMAIN_SEPARATION_LABEL = "jwt"
)

var (
	ErrWrongKeyFormat = errors.New("wrong key type")
	ErrBadSignature   = errors.New("bad signature")
	ErrVerification   = errors.New("signature verification failed")
)

/*
*
The JWT signing method for secp256k1. Inspired by https://github.com/ureeves/jwt-go-secp256k1/blob/master/secp256k1.go
but updated to work with the latest version of jwt-go.
*/
type SigningMethodSecp256k1 struct{}

func (sm *SigningMethodSecp256k1) Verify(signingString string, sig []byte, key interface{}) error {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return ErrWrongKeyFormat
	}

	hashedString := hashStringWithDomainSeparation(signingString)

	if len(sig) != SIG_LENGTH {
		return ErrBadSignature
	}

	r := new(big.Int).SetBytes(sig[:R_LENGTH])                    // R
	s := new(big.Int).SetBytes(sig[R_LENGTH : R_LENGTH+S_LENGTH]) // S

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

	hashedString := hashStringWithDomainSeparation(signingString)

	sig, err := ethcrypto.Sign(hashedString, priv)
	if err != nil {
		return nil, err
	}

	if len(sig) != SIG_LENGTH {
		return nil, ErrBadSignature
	}

	return sig, nil
}

func (sm *SigningMethodSecp256k1) Alg() string {
	return ALGORITHM
}

func hashStringWithDomainSeparation(signingString string) []byte {
	return ethcrypto.Keccak256([]byte(DOMAIN_SEPARATION_LABEL + signingString))
}

func init() {
	method := &SigningMethodSecp256k1{}
	jwt.RegisterSigningMethod(ALGORITHM, func() jwt.SigningMethod {
		return method
	})
}

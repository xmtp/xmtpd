package authn

import (
	"crypto/ecdsa"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_DURATION = time.Hour
)

type TokenFactory struct {
	privateKey *ecdsa.PrivateKey
	nodeID     int32
}

func NewTokenFactory(privateKey *ecdsa.PrivateKey, nodeID int32) *TokenFactory {
	return &TokenFactory{
		privateKey: privateKey,
		nodeID:     nodeID,
	}
}

func (f *TokenFactory) CreateToken(forNodeID int32) (*Token, error) {
	now := time.Now()
	expiresAt := now.Add(TOKEN_DURATION)

	token := jwt.NewWithClaims(&SigningMethodSecp256k1{}, &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(f.nodeID)),
		Audience:  []string{strconv.Itoa(int(forNodeID))},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
	})

	signedString, err := token.SignedString(f.privateKey)
	if err != nil {
		return nil, err
	}

	return NewToken(signedString, expiresAt), nil
}

type Token struct {
	SignedString string
	ExpiresAt    time.Time
}

func NewToken(signedString string, expiresAt time.Time) *Token {
	return &Token{
		SignedString: signedString,
		ExpiresAt:    expiresAt,
	}
}

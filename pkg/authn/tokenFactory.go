package authn

import (
	"crypto/ecdsa"
	"github.com/Masterminds/semver/v3"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_DURATION = time.Hour
)

type TokenFactory struct {
	privateKey    *ecdsa.PrivateKey
	nodeID        uint32
	serverVersion *semver.Version
}

func NewTokenFactory(
	privateKey *ecdsa.PrivateKey,
	nodeID uint32,
	serverVersion *semver.Version,
) TokenFactory {
	return TokenFactory{
		privateKey:    privateKey,
		nodeID:        nodeID,
		serverVersion: serverVersion,
	}
}

func (f *TokenFactory) CreateToken(forNodeID uint32) (*Token, error) {
	now := time.Now()
	expiresAt := now.Add(TOKEN_DURATION)

	claims := &XmtpdClaims{
		Version: f.serverVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(f.nodeID)),
			Audience:  []string{strconv.Itoa(int(forNodeID))},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Create a new token with custom claims
	token := jwt.NewWithClaims(&SigningMethodSecp256k1{}, claims)

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

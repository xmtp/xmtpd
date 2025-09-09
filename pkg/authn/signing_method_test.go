package authn_test

import (
	"testing"

	"github.com/xmtp/xmtpd/pkg/authn"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestSign(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	publicKey := privateKey.Public()

	method := &authn.SigningMethodSecp256k1{}

	signingString := "test"
	signature, err := method.Sign(signingString, privateKey)
	require.NoError(t, err)

	err = method.Verify(signingString, signature, publicKey)
	require.NoError(t, err)
}

func TestWrongSigner(t *testing.T) {
	goodPrivateKey := testutils.RandomPrivateKey(t)

	badPrivateKey := testutils.RandomPrivateKey(t)
	badPublicKey := badPrivateKey.Public()

	method := &authn.SigningMethodSecp256k1{}

	signingString := "test"
	signature, err := method.Sign(signingString, goodPrivateKey)
	require.NoError(t, err)

	err = method.Verify(signingString, signature, badPublicKey)
	require.Error(t, err)
}

func TestWrongSigningString(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	publicKey := privateKey.Public()

	method := &authn.SigningMethodSecp256k1{}

	signingString := "test"
	signature, err := method.Sign(signingString, privateKey)
	require.NoError(t, err)

	err = method.Verify("wrong signing string", signature, publicKey)
	require.Error(t, err)
}

func TestFullJWT(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	claims := &jwt.RegisteredClaims{
		Issuer: "test",
	}
	token := jwt.NewWithClaims(&authn.SigningMethodSecp256k1{}, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return privateKey.Public(), nil
	})
	require.NoError(t, err)

	issuer, err := parsedToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, issuer, claims.Issuer)
}

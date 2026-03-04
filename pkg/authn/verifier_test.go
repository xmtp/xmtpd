package authn_test

import (
	"crypto/ecdsa"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/xmtp/xmtpd/pkg/authn"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	registryMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registry"
)

const (
	VERIFIER_NODE_ID = 100
	SIGNER_NODE_ID   = 200
)

func buildVerifier(
	t *testing.T,
	verifierNodeID uint32,
	version *semver.Version,
) (*authn.RegistryVerifier, *registryMocks.MockNodeRegistry) {
	mockRegistry := registryMocks.NewMockNodeRegistry(t)
	verifier, err := authn.NewRegistryVerifier(
		testutils.NewLog(t),
		mockRegistry,
		verifierNodeID,
		version,
	)
	require.NoError(t, err)

	return verifier, mockRegistry
}

func buildJwt(
	t *testing.T,
	signerPrivateKey *ecdsa.PrivateKey,
	signerNodeID int,
	verifierNodeID int,
	issuedAt time.Time,
	expiresAt time.Time,
) string {
	token := jwt.NewWithClaims(&authn.SigningMethodSecp256k1{}, &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(signerNodeID)),
		Audience:  []string{strconv.Itoa(int(verifierNodeID))},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	})

	signedString, err := token.SignedString(signerPrivateKey)
	require.NoError(t, err)

	return signedString
}

func TestVerifier(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	tokenFactory := authn.NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), nil)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &signerPrivateKey.PublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)

	// Create a token targeting the verifier's node as the audience
	token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
	require.NoError(t, err)
	// This should verify correctly
	_, cancel, verificationError := verifier.Verify(token.SignedString)
	defer cancel()
	require.NoError(t, verificationError)

	// Create a token targeting a different node as the audience
	tokenForWrongNode, err := tokenFactory.CreateToken(uint32(300))
	require.NoError(t, err)
	// This should not verify correctly
	_, cancel, verificationError = verifier.Verify(tokenForWrongNode.SignedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestWrongAudience(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	tokenFactory := authn.NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), nil)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &signerPrivateKey.PublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)
	// Create a token targeting a different node as the audience
	tokenForWrongNode, err := tokenFactory.CreateToken(uint32(300))
	require.NoError(t, err)
	// This should not verify correctly
	_, cancel, verificationError := verifier.Verify(tokenForWrongNode.SignedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestUnknownNode(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	tokenFactory := authn.NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), nil)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(nil, errors.New("node not found"))

	token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
	require.NoError(t, err)

	_, cancel, verificationError := verifier.Verify(token.SignedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestWrongPublicKey(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	tokenFactory := authn.NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), nil)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)

	wrongPublicKey := testutils.RandomPrivateKey(t).PublicKey
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &wrongPublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)

	token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
	require.NoError(t, err)

	_, cancel, verificationError := verifier.Verify(token.SignedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestExpiredToken(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &signerPrivateKey.PublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)

	signedString := buildJwt(
		t,
		signerPrivateKey,
		SIGNER_NODE_ID,
		VERIFIER_NODE_ID,
		time.Now().Add(-2*time.Hour),
		time.Now().Add(-time.Hour),
	)

	_, cancel, verificationError := verifier.Verify(signedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestTokenDurationTooLong(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &signerPrivateKey.PublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)

	signedString := buildJwt(
		t,
		signerPrivateKey,
		SIGNER_NODE_ID,
		VERIFIER_NODE_ID,
		time.Now(),
		time.Now().Add(5*time.Hour),
	)

	_, cancel, verificationError := verifier.Verify(signedString)
	defer cancel()
	require.Error(t, verificationError)
}

func TestTokenClockSkew(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	verifier, nodeRegistry := buildVerifier(
		t,
		uint32(VERIFIER_NODE_ID),
		testutils.GetLatestVersion(t),
	)
	nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
		SigningKey: &signerPrivateKey.PublicKey,
		NodeID:     uint32(SIGNER_NODE_ID),
	}, nil)

	// Tokens issued 1 minute in the future are OK
	validToken := buildJwt(
		t,
		signerPrivateKey,
		SIGNER_NODE_ID,
		VERIFIER_NODE_ID,
		time.Now().Add(1*time.Minute),
		time.Now().Add(1*time.Hour),
	)

	_, cancel, verificationError := verifier.Verify(validToken)
	defer cancel()
	require.NoError(t, verificationError)

	invalidToken := buildJwt(
		t,
		signerPrivateKey,
		SIGNER_NODE_ID,
		VERIFIER_NODE_ID,
		time.Now().Add(10*time.Minute),
		time.Now().Add(1*time.Hour),
	)

	_, cancel, verificationError = verifier.Verify(invalidToken)
	defer cancel()
	require.Error(t, verificationError)
}

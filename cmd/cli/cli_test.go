package main

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestRegisterNodeArgParse(t *testing.T) {
	httpAddress := "http://localhost:8545"
	ownerAddress := testutils.RandomAddress()
	adminPrivateKey := utils.EcdsaPrivateKeyToString(testutils.RandomPrivateKey(t))
	signingKeyPub := utils.EcdsaPublicKeyToString(
		testutils.RandomPrivateKey(t).Public().(*ecdsa.PublicKey),
	)
	options, err := parseOptions(
		[]string{
			"register-node",
			"--http-address",
			httpAddress,
			"--admin-private-key",
			adminPrivateKey,
			"--owner-address",
			ownerAddress.Hex(),
			"--signing-key",
			signingKeyPub,
		},
	)
	require.NoError(t, err)
	require.Equal(t, options.Command, "register-node")
	require.Equal(t, options.RegisterNode.AdminPrivateKey, adminPrivateKey)
	require.Equal(t, options.RegisterNode.OwnerAddress, ownerAddress.Hex())
	require.Equal(t, options.RegisterNode.SigningKey, signingKeyPub)

	// Test missing options
	_, err = parseOptions([]string{"register-node"})
	require.Error(t, err)
	require.Equal(
		t,
		err.Error(),
		"Could not parse options: the required flags `--admin-private-key', `--http-address', `--owner-address' and `--signing-key' were not specified",
	)
}

func TestGenerateKeyArgParse(t *testing.T) {
	options, err := parseOptions([]string{"generate-key"})
	require.NoError(t, err)
	require.Equal(t, options.Command, "generate-key")
}

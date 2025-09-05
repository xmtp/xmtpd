package main

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestRegisterNodeArgParse(t *testing.T) {
	httpAddress := "ws://localhost:8545"
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
			"--admin.private-key",
			adminPrivateKey,
			"--node-owner-address",
			ownerAddress.Hex(),
			"--node-signing-key-pub",
			signingKeyPub,
		},
	)
	require.NoError(t, err)
	require.Equal(t, options.Command, "register-node")
	require.Equal(t, options.RegisterNode.AdminOptions.AdminPrivateKey, adminPrivateKey)
	require.Equal(t, options.RegisterNode.OwnerAddress.String(), ownerAddress.Hex())
	require.Equal(t, options.RegisterNode.SigningKeyPub, signingKeyPub)

	// Test missing options
	_, err = parseOptions([]string{"register-node"})
	require.Error(t, err)
	require.Equal(
		t,
		err.Error(),
		"could not parse options: the required flags `--admin.private-key', `--http-address', `--node-owner-address' and `--node-signing-key-pub' were not specified",
	)
}

func TestGenerateKeyArgParse(t *testing.T) {
	options, err := parseOptions([]string{"generate-key"})
	require.NoError(t, err)
	require.Equal(t, options.Command, "generate-key")
}

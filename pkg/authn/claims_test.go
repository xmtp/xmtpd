package authn

import (
	"bytes"
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"os/exec"
	"strings"
	"testing"
)

func getLatestTag(t *testing.T) string {
	// Prepare the command
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	require.NoError(t, err)
	return strings.TrimSpace(out.String())
}

func getLatestVersion(t *testing.T) semver.Version {
	tag := getLatestTag(t)
	v, err := semver.NewVersion(tag)
	require.NoError(t, err)

	return *v
}

func newVersionNoError(t *testing.T, version string, pre string, meta string) semver.Version {
	v, err := semver.NewVersion(version)
	require.NoError(t, err)

	vextras, err := v.SetPrerelease(pre)
	require.NoError(t, err)

	vmeta, err := vextras.SetMetadata(meta)
	require.NoError(t, err)

	return vmeta
}

func TestClaimsVerifierNoVersion(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	tests := []struct {
		name    string
		version *semver.Version
		wantErr bool
	}{
		{"no-version", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenFactory := NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), tt.version)

			verifier, nodeRegistry := buildVerifier(t, uint32(VERIFIER_NODE_ID))
			nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
				SigningKey: &signerPrivateKey.PublicKey,
				NodeID:     uint32(SIGNER_NODE_ID),
			}, nil)

			token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
			require.NoError(t, err)
			verificationError := verifier.Verify(token.SignedString)
			if tt.wantErr {
				require.Error(t, verificationError)
			} else {
				require.NoError(t, verificationError)
			}
		})
	}
}

func TestClaimsVerifier(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	currentVersion := getLatestVersion(t)

	tests := []struct {
		name    string
		version semver.Version
		wantErr bool
	}{
		{"current-version", currentVersion, false},
		{"next-patch-version", currentVersion.IncPatch(), false},
		{"next-minor-version", currentVersion.IncMinor(), true},
		{"next-major-version", currentVersion.IncMajor(), true},
		{"last-supported-version", newVersionNoError(t, "0.1.3", "", ""), false},
		{"with-prerelease-version", newVersionNoError(t, "0.1.3", "17-gdeadbeef", ""), false},
		{
			"with-metadata-version",
			newVersionNoError(t, "0.1.3", "", "branch-dev"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenFactory := NewTokenFactory(signerPrivateKey, uint32(SIGNER_NODE_ID), &tt.version)

			verifier, nodeRegistry := buildVerifier(t, uint32(VERIFIER_NODE_ID))
			nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
				SigningKey: &signerPrivateKey.PublicKey,
				NodeID:     uint32(SIGNER_NODE_ID),
			}, nil)

			token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
			require.NoError(t, err)
			verificationError := verifier.Verify(token.SignedString)
			if tt.wantErr {
				require.Error(t, verificationError)
			} else {
				require.NoError(t, verificationError)
			}
		})
	}
}

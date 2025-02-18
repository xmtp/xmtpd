package authn_test

import (
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"testing"
)

func newVersionNoError(t *testing.T, version string, pre string, meta string) semver.Version {
	v, err := semver.NewVersion(version)
	require.NoError(t, err)

	vextras, err := v.SetPrerelease(pre)
	require.NoError(t, err)

	vmeta, err := vextras.SetMetadata(meta)
	require.NoError(t, err)

	return vmeta
}

func TestClaimsNoVersion(t *testing.T) {
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
			tokenFactory := authn.NewTokenFactory(
				signerPrivateKey,
				uint32(SIGNER_NODE_ID),
				tt.version,
			)

			verifier, nodeRegistry := buildVerifier(
				t,
				uint32(VERIFIER_NODE_ID),
				testutils.GetLatestVersion(t),
			)
			nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
				SigningKey: &signerPrivateKey.PublicKey,
				NodeID:     uint32(SIGNER_NODE_ID),
			}, nil)

			token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
			require.NoError(t, err)
			_, verificationError := verifier.Verify(token.SignedString)
			if tt.wantErr {
				require.Error(t, verificationError)
			} else {
				require.NoError(t, verificationError)
			}
		})
	}
}

func TestClaimsVariousVersions(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	currentVersion := *testutils.GetLatestVersion(t)
	version013, err := semver.NewVersion("0.1.3")
	require.NoError(t, err)
	version014, err := semver.NewVersion("0.1.4")
	require.NoError(t, err)

	tests := []struct {
		name    string
		version semver.Version
		wantErr bool
	}{
		{"current-version", currentVersion, false},
		{"next-patch-version", currentVersion.IncPatch(), false},
		{"next-minor-version", currentVersion.IncMinor(), true},
		{"next-major-version", currentVersion.IncMajor(), true},
		{
			"with-prerelease-version",
			newVersionNoError(t, currentVersion.String(), "17-gdeadbeef", ""),
			false,
		},
		{
			"with-metadata-version",
			newVersionNoError(t, currentVersion.String(), "", "branch-dev"),
			false,
		},
		{"known-0.1.3", *version013, true},
		{"known-0.1.4", *version014, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenFactory := authn.NewTokenFactory(
				signerPrivateKey,
				uint32(SIGNER_NODE_ID),
				&tt.version,
			)

			verifier, nodeRegistry := buildVerifier(
				t,
				uint32(VERIFIER_NODE_ID),
				testutils.GetLatestVersion(t),
			)
			nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
				SigningKey: &signerPrivateKey.PublicKey,
				NodeID:     uint32(SIGNER_NODE_ID),
			}, nil)

			token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
			require.NoError(t, err)
			_, verificationError := verifier.Verify(token.SignedString)
			if tt.wantErr {
				require.Error(t, verificationError)
			} else {
				require.NoError(t, verificationError)
			}
		})
	}
}

func TestClaimsValidator(t *testing.T) {
	signerPrivateKey := testutils.RandomPrivateKey(t)

	currentVersion := *testutils.GetLatestVersion(t)

	tests := []struct {
		name          string
		version       semver.Version
		serverVersion semver.Version
		wantErr       bool
	}{
		{"current-version", currentVersion, currentVersion, false},
		{
			"with-prerelease-version",
			currentVersion,
			newVersionNoError(t, currentVersion.String(), "17-gdeadbeef", ""),
			false,
		},
		{
			"with-metadata-version",
			currentVersion,
			newVersionNoError(t, currentVersion.String(), "", "branch-dev"),
			false,
		},
		{
			"future-major-rejects-us",
			currentVersion,
			currentVersion.IncMajor(),
			true,
		},
		{
			"future-patch-accepts-us",
			currentVersion,
			currentVersion.IncPatch(),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenFactory := authn.NewTokenFactory(
				signerPrivateKey,
				uint32(SIGNER_NODE_ID),
				&tt.version,
			)

			verifier, nodeRegistry := buildVerifier(t, uint32(VERIFIER_NODE_ID), &tt.serverVersion)
			nodeRegistry.EXPECT().GetNode(uint32(SIGNER_NODE_ID)).Return(&registry.Node{
				SigningKey: &signerPrivateKey.PublicKey,
				NodeID:     uint32(SIGNER_NODE_ID),
			}, nil)

			token, err := tokenFactory.CreateToken(uint32(VERIFIER_NODE_ID))
			require.NoError(t, err)
			_, verificationError := verifier.Verify(token.SignedString)
			if tt.wantErr {
				require.Error(t, verificationError)
			} else {
				require.NoError(t, verificationError)
			}
		})
	}
}

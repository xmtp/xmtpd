package authn_test

import (
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/xmtp/xmtpd/pkg/authn"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestTokenFactory(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	factory := authn.NewTokenFactory(privateKey, 100, nil)

	token, err := factory.CreateToken(200)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEmpty(t, token.SignedString)
	require.NotZero(t, token.ExpiresAt)
	require.True(t, token.ExpiresAt.After(time.Now().Add(59*time.Minute)))
	require.True(t, token.ExpiresAt.Before(time.Now().Add(61*time.Minute)))
}

func TestTokenFactoryWithVersion(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)

	tests := []struct {
		name    string
		version string
	}{
		{"current-ish", "0.1.3"},
		{"future-ish", "11.7.3"},
		{"with-git-describe", "0.1.0-15-gdeadbeef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := semver.NewVersion(tt.version)
			require.NoError(t, err)
			factory := authn.NewTokenFactory(privateKey, 100, version)

			token, err := factory.CreateToken(200)
			require.NoError(t, err)
			require.NotNil(t, token)
		})
	}
}

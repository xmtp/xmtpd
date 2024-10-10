package authn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestTokenFactory(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	factory := NewTokenFactory(privateKey, 100)

	token, err := factory.CreateToken(200)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEmpty(t, token.SignedString)
	require.NotZero(t, token.ExpiresAt)
	require.True(t, token.ExpiresAt.After(time.Now().Add(59*time.Minute)))
	require.True(t, token.ExpiresAt.Before(time.Now().Add(61*time.Minute)))
}

package stress

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestEnvelopesGenerator(t *testing.T) {
	generator, err := NewEnvelopesGenerator(
		"http://localhost:5050",
		testutils.TEST_PRIVATE_KEY,
		100,
	)
	require.NoError(t, err)

	resp, err := generator.PublishWelcomeMessageEnvelopes(context.Background(), 1, 100)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

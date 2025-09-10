package stress

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvelopesGenerator(t *testing.T) {
	generator, err := NewEnvelopesGenerator(
		"http://localhost:5050",
		"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		100,
	)
	require.NoError(t, err)

	generator.PublishWelcomeMessageEnvelopes(context.Background(), 10, 100)
}

package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func Marshal(t *testing.T, msg proto.Message) []byte {
	bytes, err := proto.Marshal(msg)
	require.NoError(t, err)
	return bytes
}

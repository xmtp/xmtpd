package testing

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

func RequireProtoEqual(t *testing.T, expected, actual any) {
	require.Equal(t, "", cmp.Diff(expected, actual, protocmp.Transform()))
}

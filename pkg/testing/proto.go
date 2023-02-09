package testing

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

func RequireProtoEqual[S ~[]E, E any](t *testing.T, expected, actual S) {
	t.Helper()
	require.Equal(t, "", cmp.Diff(expected, actual, protocmp.Transform()))
}

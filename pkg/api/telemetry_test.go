package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI_Telemetry_parseVersionHeaderValue(t *testing.T) {
	tcs := []struct {
		name            string
		value           []string
		expectedName    string
		expectedVersion string
		expectedFull    string
	}{
		{
			name:            "empty",
			value:           []string{""},
			expectedName:    "",
			expectedVersion: "",
			expectedFull:    "",
		},
		{
			name:            "name and version",
			value:           []string{"test/0.0.0"},
			expectedName:    "test",
			expectedVersion: "0.0.0",
			expectedFull:    "test/0.0.0",
		},
		{
			name:            "name only with slash",
			value:           []string{"test/"},
			expectedName:    "test",
			expectedVersion: "",
			expectedFull:    "test/",
		},
		{
			name:            "name only without slash",
			value:           []string{"test"},
			expectedName:    "test",
			expectedVersion: "",
			expectedFull:    "test",
		},
		{
			name:            "version only with slash",
			value:           []string{"/0.0.0"},
			expectedName:    "",
			expectedVersion: "0.0.0",
			expectedFull:    "/0.0.0",
		},
		{
			name:            "multiple values",
			value:           []string{"test/0.0.0", "other/0.0.1"},
			expectedName:    "test",
			expectedVersion: "0.0.0",
			expectedFull:    "test/0.0.0",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			name, version, full := parseVersionHeaderValue(tc.value)
			require.Equal(t, tc.expectedName, name)
			require.Equal(t, tc.expectedVersion, version)
			require.Equal(t, tc.expectedFull, full)
		})
	}
}

// Package flags provides control over flags that are set at build time.
package flags

import "testing"

var raceTestEnabled bool

func SkipOnRaceTest(t *testing.T) {
	if raceTestEnabled {
		t.Skip("Skipping go test -race")
	}
}

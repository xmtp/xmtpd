//go:build race

package flags

func init() {
	// raceTestEnabled is only enabled when running `go test -race`.
	raceTestEnabled = true
}

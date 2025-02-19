package testutils

import (
	"bytes"
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
	"os/exec"
	"strings"
	"testing"
)

func GetLatestTag(t *testing.T) string {
	// Prepare the command
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	require.NoError(t, err, out.String())
	return strings.TrimSpace(out.String())
}

func GetLatestVersion(t *testing.T) *semver.Version {
	tag := GetLatestTag(t)
	v, err := semver.NewVersion(tag)
	require.NoError(t, err)

	return v
}

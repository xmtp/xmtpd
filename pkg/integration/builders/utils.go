// Package builders implements a collection of container builders for the integration tests.
package builders

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const (
	xdbgImage             = "ghcr.io/xmtp/xdbg:sha-ac533c4"
	cliImage              = "ghcr.io/xmtp/xmtpd-cli:sha-06cb109"
	adminPrivateKey       = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	anvilJSONRelativePath = "../../../pkg/config/environments/anvil.json"
)

func loadEnvFromShell() (map[string]string, error) {
	scriptPath := testutils.GetScriptPath("../scripts/load_env.sh")
	cmd := exec.Command(scriptPath)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(
			"error loading env via shell script: %w\nError: %s",
			err,
			errBuf.String(),
		)
	}

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(&outBuf)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap, nil
}

func expandVars(vars map[string]string) {
	vars["XMTPD_API_ENABLE"] = "true"
	vars["XMTPD_INDEXER_ENABLE"] = "true"

	dbName := testutils.GetCallerName(3) + "_" + testutils.RandomStringLower(6)

	vars["XMTPD_DB_NAME_OVERRIDE"] = dbName
}

func convertLocalhost(vars map[string]string) {
	for varKey, varValue := range vars {
		if strings.Contains(varValue, "localhost") {
			vars[varKey] = strings.ReplaceAll(varValue, "localhost", "host.docker.internal")
		}
	}
}

func constructVariables(t *testing.T) map[string]string {
	envVars, err := loadEnvFromShell()
	require.NoError(t, err)
	expandVars(envVars)
	convertLocalhost(envVars)

	return envVars
}

func handleExitedContainer(
	context context.Context,
	exitedContainer testcontainers.Container,
) error {
	state, err := exitedContainer.State(context)
	if err != nil {
		return err
	}

	if state.ExitCode != 0 {
		logs, logErr := exitedContainer.Logs(context)
		if logErr != nil {
			return fmt.Errorf(
				"container exited with code %d, but failed to get logs: %w",
				state.ExitCode,
				logErr,
			)
		}
		defer func() {
			_ = logs.Close()
		}()

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, logs)

		return fmt.Errorf("container exited with code %d\nLogs:\n%s", state.ExitCode, buf.String())
	}

	return nil
}

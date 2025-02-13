package upgrade_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const testFlag = "ENABLE_UPGRADE_TESTS"

func skipIfNotEnabled(t *testing.T) {
	if _, isSet := os.LookupEnv(testFlag); !isSet {
		t.Skip("Skipping upgrade test")
	}
}

func getScriptPath(scriptName string) string {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	return filepath.Join(baseDir, scriptName)
}

func loadEnvFromShell() (map[string]string, error) {
	scriptPath := getScriptPath("./scripts/load_env.sh")
	cmd := exec.Command(scriptPath)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(
			"error loading env via shell script: %v\nError: %s",
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
	vars["XMTPD_REPLICATION_ENABLE"] = "true"
	vars["XMTPD_INDEXER_ENABLE"] = "true"

	dbName := testutils.GetCallerName(3) + "_" + testutils.RandomStringLower(6)

	vars["XMTPD_DB_NAME_OVERRIDE"] = dbName
}

func convertLocalhost(vars map[string]string) {
	for varKey, varValue := range vars {
		if strings.Contains(varValue, "localhost") {
			vars[varKey] = strings.Replace(varValue, "localhost", "host.docker.internal", -1)
		}
	}
}

func dockerRmc(containerName string) error {
	killCmd := exec.Command("docker", "rm", containerName)
	return killCmd.Run()
}

func dockerKill(containerName string) error {
	killCmd := exec.Command("docker", "kill", containerName)
	return killCmd.Run()
}

func dockerLogs(containerName string) (string, error) {
	logsCmd := exec.Command("docker", "logs", containerName)
	var outBuf bytes.Buffer
	logsCmd.Stdout = &outBuf
	err := logsCmd.Run()
	if err != nil {
		return "", err
	}
	return outBuf.String(), nil
}

func constructVariables(t *testing.T) map[string]string {
	envVars, err := loadEnvFromShell()
	require.NoError(t, err)
	expandVars(envVars)
	convertLocalhost(envVars)

	return envVars
}

func streamDockerLogs(containerName string) (chan string, func(), error) {
	logsCmd := exec.Command("docker", "logs", "-f", containerName)
	stdoutPipe, err := logsCmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	err = logsCmd.Start()
	if err != nil {
		return nil, nil, err
	}

	logChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			logChan <- scanner.Text()
		}
		close(logChan)
	}()

	cancelFunc := func() {
		_ = logsCmd.Process.Kill()
	}

	return logChan, cancelFunc, nil
}

func runContainer(
	t *testing.T,
	containerName string,
	imageName string,
	envVars map[string]string,
) {
	var dockerEnvArgs []string
	for key, value := range envVars {
		dockerEnvArgs = append(dockerEnvArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	_ = dockerRmc(containerName)

	dockerCmd := append([]string{"run", "-d"}, dockerEnvArgs...)
	dockerCmd = append(dockerCmd, "--name", containerName, imageName)

	cmd := exec.Command("docker", dockerCmd...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	require.NoError(t, err, "Error: %s", errBuf.String())

	defer func() {
		_ = dockerKill(containerName)
	}()

	logChan, cancel, err := streamDockerLogs(containerName)
	require.NoError(t, err, "Failed to start log streaming")
	defer cancel()

	timeout := time.After(5 * time.Second)

	for {
		select {
		case line, ok := <-logChan:
			if !ok {
				t.Fatalf("Log stream closed before finding target log")
			}
			t.Logf(line)
			if strings.Contains(line, "replication.api\tserving grpc") {
				t.Logf("Service started successfully")
				return
			}
		case <-timeout:
			t.Fatalf("Timeout: 'replication.api\tserving grpc' not found in logs within 5 seconds")
		}
	}
}

func buildDevImage(t *testing.T) {
	scriptPath := getScriptPath("../../dev/docker/build")
	cmd := exec.Command(scriptPath)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	require.NoError(t, err, "build failed:\n%s\n%s", errBuf.String(), outBuf.String())

}

func TestUpgradeFrom014(t *testing.T) {
	//skipIfNotEnabled(t)
	envVars := constructVariables(t)

	buildDevImage(t)

	t.Logf("Starting old container")
	runContainer(t, "xmtpd_test_014", "ghcr.io/xmtp/xmtpd:0.1.4", envVars)

	t.Logf("Starting new container")
	runContainer(t, "xmtpd_test_dev", "ghcr.io/xmtp/xmtpd:dev", envVars)
}

package upgrade_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

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
	vars["XMTPD_REFLECTION_ENABLE"] = "true"
	vars["XMTPD_PAYER_ENABLE"] = "true"
	vars["XMTPD_REPLICATION_ENABLE"] = "true"
	vars["XMTPD_SYNC_ENABLE"] = "true"
	vars["XMTPD_INDEXER_ENABLE"] = "true"
}

func convertLocalhost(vars map[string]string) {
	for varKey, varValue := range vars {
		if strings.Contains(varValue, "localhost") {
			vars[varKey] = strings.Replace(varValue, "localhost", "host.docker.internal", -1)
		}
	}
}

func TestLoadEnvAndPassToDocker(t *testing.T) {
	envVars, err := loadEnvFromShell()
	if err != nil {
		t.Fatalf("Failed to load environment variables: %v", err)
	}
	expandVars(envVars)
	convertLocalhost(envVars)

	var dockerEnvArgs []string
	for key, value := range envVars {
		dockerEnvArgs = append(dockerEnvArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}
	containerName := "xmtpd_test_env"

	dockerCmd := append([]string{"run"}, dockerEnvArgs...)
	dockerCmd = append(dockerCmd, "--name", containerName, "ghcr.io/xmtp/xmtpd:0.1.3")

	t.Logf("Running docker command: %v", dockerCmd)
	cmd := exec.Command("docker", dockerCmd...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Start()
	if err != nil {
		t.Fatalf(
			"Docker run failed: %v\nOutput: %s\nError: %s",
			err,
			outBuf.String(),
			errBuf.String(),
		)
	}

	// Wait for 5 seconds
	time.Sleep(5 * time.Second)

	// Collect logs
	logsCmd := exec.Command("docker", "logs", containerName)
	var logsBuf bytes.Buffer
	logsCmd.Stdout = &logsBuf
	logsCmd.Stderr = &errBuf

	err = logsCmd.Run()
	if err != nil {
		t.Fatalf("Failed to collect logs: %v\nError: %s", err, errBuf.String())
	}

	// Kill the container
	killCmd := exec.Command("docker", "kill", containerName)
	err = killCmd.Run()
	if err != nil {
		t.Fatalf("Failed to kill Docker container: %v", err)
	}

	t.Logf("Logs:\n%s", logsBuf.String())

}

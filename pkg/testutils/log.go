package testutils

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	// Set log level for zap debugger used in tests. Default is 'debug' and can be quite noisy.
	envLogLevel = "XMTP_TEST_LOG_LEVEL"

	// Disable stack trace in Zap for warn and above log levels. Since we test negative cases errors
	// are expected and should not clutter test output.
	envDisableStackTrace = "XMTP_TEST_DISABLE_STACK_TRACE"

	// Write logs to a per-test JSON log file.
	envFileLogger = "XMTP_TEST_FILE_LOGGER"
)

var (
	logLevel          = os.Getenv(envLogLevel)
	disableStackTrace = parseBoolConfig(os.Getenv(envDisableStackTrace))
	logToFile         = parseBoolConfig(os.Getenv(envFileLogger))
)

func NewLog(t testing.TB) *zap.Logger {
	if logToFile {
		return NewJSONLog(t)
	}

	level, err := zap.ParseAtomicLevel(strings.ToLower(logLevel))
	if err != nil {
		// Default to debug log level.
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = level
	cfg.DisableStacktrace = disableStackTrace

	log, err := cfg.Build()
	require.NoError(t, err)

	return log
}

func NewJSONLog(t testing.TB) *zap.Logger {
	level, err := zap.ParseAtomicLevel(strings.ToLower(logLevel))
	if err != nil {
		// Default to debug log level.
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// test_log_testname_hhmmss.json
	logName := fmt.Sprintf("test_log_%v_%v.json", t.Name(), time.Now().Format("150405"))

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = level
	cfg.DisableCaller = disableStackTrace

	cfg.Encoding = "json"
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.OutputPaths = []string{logName}

	log, err := cfg.Build()
	require.NoError(t, err)

	return log
}

func parseBoolConfig(str string) bool {
	switch strings.ToLower(str) {
	case "1", "true", "y", "yes":
		return true
	default:
		return false
	}
}

// CapturingWriteSyncer is a WriteSyncer that stores logs in memory.
type CapturingWriteSyncer struct {
	logs []string
}

func (c *CapturingWriteSyncer) Write(p []byte) (n int, err error) {
	c.logs = append(c.logs, string(p))
	return len(p), nil
}

func (c *CapturingWriteSyncer) Sync() error {
	return nil
}

// CapturingLogger wraps a zap.Logger and stores emitted logs.
type CapturingLogger struct {
	*zap.Logger
	cws *CapturingWriteSyncer
}

// NewCapturingLogger creates a zap.Logger that records all logs in memory.
func NewCapturingLogger(level zapcore.Level) *CapturingLogger {
	cws := &CapturingWriteSyncer{}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		cws,
		level,
	)

	return &CapturingLogger{
		Logger: zap.New(core),
		cws:    cws,
	}
}

// Logs returns all captured logs as a single string.
func (cl *CapturingLogger) Logs() string {
	return strings.Join(cl.cws.logs, "\n")
}

// Contains checks if any captured log contains the given substring.
func (cl *CapturingLogger) Contains(substr string) bool {
	for _, l := range cl.cws.logs {
		if strings.Contains(l, substr) {
			return true
		}
	}
	return false
}

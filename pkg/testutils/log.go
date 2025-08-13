package testutils

import (
	"flag"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug level logging in tests")
}

func NewLog(t testing.TB) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	if !debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	log, err := cfg.Build()
	require.NoError(t, err)
	return log
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

package testutils

import (
	"flag"
	"testing"

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
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	log, err := cfg.Build()
	require.NoError(t, err)
	return log
}

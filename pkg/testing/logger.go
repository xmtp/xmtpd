package testing

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var debug, json bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug level logging in tests")
	flag.BoolVar(&json, "json-logs", false, "log in json format")
}

func NewLogger(t testing.TB) *zap.Logger {
	log, err := zap.NewDevelopmentLogger(debug, json)
	require.NoError(t, err)
	return log
}

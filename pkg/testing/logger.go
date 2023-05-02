package testing

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var debug, json bool

func init() {
	flag.BoolVar(&debug, "log-debug", false, "debug level logging in tests")
	flag.BoolVar(&json, "log-json", false, "log in json format in tests")
}

func NewLogger(t testing.TB) *zap.Logger {
	log, err := zap.NewDevelopmentLogger(debug, json)
	require.NoError(t, err)
	return log
}

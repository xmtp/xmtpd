package testing

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug level logging in tests")
}

func NewLogger(t *testing.T) *zap.Logger {
	log, err := zap.NewDevelopmentLogger(debug)
	require.NoError(t, err)
	return log
}

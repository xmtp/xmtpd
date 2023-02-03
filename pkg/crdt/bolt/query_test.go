package bolt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/tests"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func Test_Query(t *testing.T) {
	fn, storeMaker := tempStoreMaker(t)
	defer os.Remove(fn)
	tests.QueryTests(t, tests.WithStore(storeMaker))
}

// helpers

func tempStoreMaker(t *testing.T) (string, func(*zap.Logger) crdt.NodeStore) {
	f, err := os.CreateTemp("", "crdt-test")
	require.NoError(t, err)
	return f.Name(), func(l *zap.Logger) crdt.NodeStore {
		l.Debug("creating db", zap.String("fn", f.Name()))
		s, err := NewStore(f.Name())
		require.NoError(t, err)
		return s
	}
}

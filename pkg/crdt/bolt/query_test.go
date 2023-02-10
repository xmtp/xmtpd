package bolt

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/tests"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func Test_Query(t *testing.T) {
	tests.QueryTests(t, tests.WithStore(tempStoreMaker(t)))
}

func randomTestModifiers(t *testing.T) []tests.ConfigModifier {
	return []tests.ConfigModifier{
		tests.WithStore(tempStoreMaker(t)),
		tests.WithPerMessageTimeout(5 * time.Millisecond),
	}
}

func Test_RandomMessages(t *testing.T) {
	t.Run("3n/1t/100m", func(t *testing.T) {
		tests.RandomMsgTest(t, 3, 1, 100, randomTestModifiers(t)...).Close()
	})
	t.Run("3n/3t/100m", func(t *testing.T) {
		tests.RandomMsgTest(t, 3, 1, 100, randomTestModifiers(t)...).Close()
	})
}

// helpers

type tempStore struct {
	fn string
	*Store
}

func (s *tempStore) Close() error {
	return os.Remove(s.fn)
}

func tempStoreMaker(t *testing.T) func(*zap.Logger) crdt.NodeStore {
	return func(l *zap.Logger) crdt.NodeStore {
		f, err := os.CreateTemp("", "crdt-test")
		require.NoError(t, err)
		l.Debug("creating db", zap.String("fn", f.Name()))
		s, err := NewStore(f.Name())
		require.NoError(t, err)
		return &tempStore{f.Name(), s}
	}
}

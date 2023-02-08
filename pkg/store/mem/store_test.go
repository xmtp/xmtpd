package memstore_test

import (
	"testing"

	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	storetest "github.com/xmtp/xmtpd/pkg/store/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryStore_QueryEnvelopes(t *testing.T) {
	storetest.TestQueryEnvelopes(t, func(t *testing.T) *storetest.TestStore {
		s := memstore.New(test.NewLogger(t))
		return storetest.NewTestStore(s)
	})
}

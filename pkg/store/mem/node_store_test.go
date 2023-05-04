package memstore_test

import (
	"testing"

	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_TopicLifecyle(t *testing.T) {
	ctx := test.NewContext(t)
	ntest.RunTopicLifecycleTest(t, memstore.NewNodeStore(ctx))
}

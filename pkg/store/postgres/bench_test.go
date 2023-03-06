package postgresstore_test

import (
	"testing"

	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

/* This was used to benchmark the EVENTS bucket arrangements by CID and ByTimeKey:
   To produce benchmark result run something like
		go test -run XXX -bench=. -count=10 | tee events-by-time.bench
   To compare multiple benchmark results use benchstat
		benchstat events-by-cid.bench events-by-time.bench
*/

func BenchmarkQuery(b *testing.B) {
	ctx := test.NewContext(b)
	topic := newTestStore(b, "topic")
	defer topic.Close()
	testStore := crdttest.NewTestStore(ctx, topic)
	testStore.Seed(b, "topic", 100000)
	b.Log("seeded")

	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 10000, 10)
	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 99000, 10)
	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 10000, 100)
	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 99000, 100)
	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 10000, 1000)
	ntest.BenchmarkQuery(b, ctx, testStore, 1000, 99000, 1000)
}
